# senior_intern_bot 

## 1. What it is
A discord bot that watches job boards for SWE/CS internships relevant to EMEA students and posts them to Discord with no human curation step.

Two properties define the product. Everything else in this spec is downstream of them:

- No missed opportunity. A posting that a user would have wanted to see must not be silently dropped.
- No human in the loop. No one reviews, approves, or hand-tags postings before they go out.

## 2. Non-goals (v1)
- Per-user preference matching, filtering, or subscriptions.
- Application tracking, deadline reminders, or CV tooling.
- Non-EMEA coverage.
- Scraping HTML careers pages. Structured sources only. 
- Ranking, scoring, or "best match" ordering.

## 3. Latency budget
Poll interval - 20 minutes per source (per-source override where a rate limit demands it).
Delivery latency, p50 - < 15 minutes.
Delivery latency, p99 - < 45 minutes.
Hard ceiling (SLO) - 6 hours.

## 4. Classification
Keyword matching is a false-negative generator by construction. The filter is three-way:

Verdict - Condition - Destination
`ACCEPT` - Confident internship AND EMEA-eligible - `#internships`
`REVIEW` - Any signal ambiguous, missing, or confusing. - `#unsorted`
`REJECT` - Confidently out of scope on an explicit signal - Logged only, never shown.

Absence of evidence routes to `REVIEW`, never `REJECT`.

The actual verdict-assignment logic sits behind a single filter interface, kept separate from polling, persistence, and delivery, so it can evolve (from a simple heuristic in v1 to something more sophisticated later) without those other pieces changing.

## 5. Source (v1)
Single source: the Greenhouse job boards API, one board per watched company (`boards-api.greenhouse.io/v1/boards/{token}/jobs`). The watch list is a configured set of companies/board tokens, not discovered automatically. Because there is only one source, polling is a single global 20-minute interval (§3) — no per-source scheduling is needed yet; that only becomes necessary once a second source with its own rate limit exists.

## 6. Identity & persistence
A posting's identity is the pair `(company, job_id)`. Seen postings are persisted in SQLite so a restart can't cause either a replay of everything already seen (duplicate spam) or a silent skip of whatever arrived while the bot was down (a missed opportunity, which §1 forbids).

Table `postings`, keyed on `(company, job_id)`:
- `company`, `job_id` — identity
- `title`, `location`, `content` — the fields classification is based on
- `posted_at` — the source's own posting timestamp
- `raw_json` — full API payload, kept for archival/debugging
- `content_hash` — a hash of exactly `title + location + content`, i.e. precisely what the filter reads and nothing more. Hashing only the classifier's inputs means a change anywhere else in the payload can't cause a spurious reclassification, and any change the filter would actually care about is guaranteed to be noticed.
- `verdict` — `ACCEPT` / `REVIEW` / `REJECT`

## 7. Reclassification rules
- Not yet in the table → classify.
- `REJECT` is a one-way door: never reconsidered, even if the posting is edited later. An employer whose posting carried an explicit out-of-scope signal is responsible for correcting it; the bot doesn't re-check.
- `ACCEPT` is also terminal: it has already been delivered, so a later edit changes nothing actionable.
- `REVIEW` is the only verdict that gets re-evaluated: if `content_hash` changes on a later poll, the posting is reclassified from scratch. This is how a posting that was ambiguous (missing signal) gets a second chance once the employer adds the missing information.

## 8. Delivery & write ordering
- Each poll classifies every newly-seen posting, partitions them into the three disjoint buckets, and composes one message per destination channel (`#internships` for `ACCEPT`, `#unsorted` for `REVIEW`), listing title + link for each posting in the batch. `REJECT` postings are never sent anywhere.
- If a batch would exceed Discord's per-message size limit, it's split across multiple messages.
- A posting's row is written to SQLite only *after* its corresponding message has sent successfully — committed per message, not per whole batch. If a send fails, the affected postings are simply absent from the table and get naturally retried (reclassified and resent) on the next poll. `REJECT` postings have no send step, so they're persisted immediately on classification.
- This ordering guarantees the spec's core property (no silent drop) at the cost of a possible duplicate post, if a crash or write failure lands in the narrow window after a successful send but before the commit. Given §1's stated priority, an occasional duplicate is an acceptable trade for never silently losing a posting.



