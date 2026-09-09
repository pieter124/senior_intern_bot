# senior_intern_bot

## 1. What it is

A Discord bot that watches job boards for SWE/CS internships relevant to EMEA
students and posts them to Discord with no human curation step.

Two properties define the product. Everything else in this spec is downstream
of them:

- **No missed opportunity.** A posting that a user would have wanted to see must
  not be silently dropped.
- **No human in the loop.** No one reviews, approves, or hand-tags postings
  before they go out.

## 2. Non-goals (v1)

- Per-user preference matching, filtering, or subscriptions.
- Application tracking, deadline reminders, or CV tooling.
- Non-EMEA coverage.
- Scraping HTML careers pages. Structured sources only.
- Ranking, scoring, or "best match" ordering.

## 3. Latency budget

| Metric | Target |
|---|---|
| Poll interval | 20 minutes per source (per-source override where a rate limit demands it) |
| Delivery latency, p50 | < 15 minutes |
| Delivery latency, p99 | < 45 minutes |
| Hard ceiling (SLO) | 6 hours |

## 4. Classification

Keyword matching is a false-negative generator by construction. The filter is
three-way:

| Verdict | Condition | Destination |
|---|---|---|
| `ACCEPT` | Confident internship AND EMEA-eligible, with no disqualifying discipline | `#internships` |
| `REVIEW` | Any signal ambiguous, missing, or confusing | `#unsorted` |
| `REJECT` | Confidently out of scope on an explicit signal | Logged only, never shown |

Absence of evidence routes to `REVIEW`, never `REJECT`.

### 4.1 Axes

Classification answers three independent questions, each against a specific
field:

| Axis | Question | Field read |
|---|---|---|
| Role | Is this an internship or student-level position? | `title` |
| Region | Is this posting in EMEA? | `location` |
| Discipline | Is this a software engineering or CS role? | `title` |

Each axis returns one of three values: in-scope, out-of-scope, or **unknown**.
Unknown is the zero value of the type, so an axis that was never evaluated can
never read as a confident answer.

An axis commits to a direction only when exactly one side of the evidence
fires. No evidence is unknown, and so is *contradictory* evidence: a title
matching both an in-scope and an out-of-scope phrase is the "confusing signal"
case above, and routes to `REVIEW`.

### 4.2 Combination rule

```
role == out || region == out || discipline == out  ->  REJECT
role == in  && region == in                        ->  ACCEPT
otherwise                                          ->  REVIEW
```

`REVIEW` is the fallthrough case, so there is no code path that reaches
`REJECT` without an explicit out-of-scope signal on some axis.

**Discipline is a veto, not a requirement.** Note the asymmetry: an
out-of-scope discipline rejects, but the *absence* of a discipline signal does
not block `ACCEPT`. Many genuine engineering internships carry no discipline
word in the title, especially non-English ones ("Praktikum
Softwareentwicklung"). Requiring an explicit signal would bury them in
`#unsorted` and would silently undo the multilingual coverage of the role
axis.

### 4.3 Keyword list asymmetry

The in-scope and out-of-scope lists on each axis are deliberately unequal in
length and in standard of evidence.

Only the out-of-scope lists can produce `REJECT`, and under §7 a `REJECT` is
never reconsidered. So those lists are the only place where a wrong entry
causes a permanent silent drop. Every other error costs one junk line in
Discord.

Consequences:

- In-scope lists are generous and multilingual.
- Out-of-scope lists are short, unambiguous, and English-only. The test for
  adding a phrase: it must be wrong for *every* posting it can appear in. One
  counterexample means it belongs in neither list, and the posting lands in
  `REVIEW`.
- No two-letter tokens anywhere. "CA" is California or Canada, "IN" is Indiana
  or India, "MA" is Massachusetts or Morocco. This is why "Remote - US"
  evaluates to unknown rather than out-of-scope, which is the accepted cost of
  the rule.
- The role axis's out-of-scope list holds seniority markers only. Nothing about
  function or team: "Backend" says nothing about whether a role is an
  internship.

### 4.4 Field discipline

Signals are read from the `title` and `location` fields, never from the job
description body. The distinction matters: a title is a curated summary, so a
word appearing in it means the role **is** that thing. A description is prose,
so a word appearing in it means the role merely **touches** that thing.

An engineering internship description will mention "sales", "marketing", and
"recruitment"; a marketing internship description will mention "engineering"
and "technology". Evaluating an axis against the body causes both sides to fire
on almost every posting, collapsing the axis to unknown and disabling it in
both directions.

If description text is used later, it may only contribute *positive* evidence,
never negative.

### 4.5 Interface

Verdict-assignment logic sits behind a single filter interface, kept separate
from polling, persistence, and delivery, so it can evolve (from the v1 keyword
heuristic to an LLM classifier later) without those pieces changing.

## 5. Sources

Multiple ATS job-board APIs, all unauthenticated JSON GETs, one board per
watched company. The watch list is a configured set of tokens per source, not
discovered automatically.

| Source | Endpoint | Status |
|---|---|---|
| Greenhouse | `boards-api.greenhouse.io/v1/boards/{token}/jobs` | v1 |
| Ashby | `api.ashbyhq.com/posting-api/job-board/{name}` | v1 |
| Lever | `api.lever.co/v0/postings/{slug}?mode=json` (also `api.eu.lever.co` for EU-provisioned accounts) | planned |

Each source is a separate implementation of one `Poller` interface. A partial
failure of one board must not discard results from the others: `Poll` returns
the postings it did retrieve alongside a joined error, and callers proceed with
the partial result. Discarding good data on one board's 500 would be a silent
drop, which §1 forbids.

### 5.1 Wire formats do not reach the domain

Each source has its own DTO with its own JSON tags, plus a mapper into
`domain.Posting`. The domain type carries no JSON tags at all.

This is not stylistic. The sources disagree on more than field names: Lever's
IDs are UUIDs while Greenhouse's are integers, and Lever's `createdAt` is epoch
milliseconds while Greenhouse's `updated_at` is an ISO string. No shared struct
tag could reconcile those. It also keeps the filter and storer from being
implicitly coupled to one vendor's payload shape.

### 5.2 Descriptions are not fetched

Greenhouse's `?content=true` inlines full HTML descriptions. Measured across 33
boards, that raises a poll cycle from a few MB to **81 MB** (6,678 postings,
~12 KB each), which at 20-minute intervals is 5.7 GB/day. The largest boards
(Databricks 9 MB, Anthropic 8.4 MB) also exceed a 10-second request timeout on
an ordinary connection, meaning the biggest sources would be the ones that fail
most reliably.

Since §4.4 restricts classification to `title` and `location`, descriptions buy
nothing. They are not requested.

If description-based signals are wanted later, the correct shape is to fetch
descriptions only for the handful of `REVIEW` postings per cycle via the
single-job endpoint, not for all postings on every cycle. Note that Greenhouse
double-encodes the content field (`&lt;p&gt;`), so entity decoding and tag
stripping would be required before any matching.

## 6. Identity and persistence

A posting's identity is the triple `(source, company, job_id)`. The source
component is required because the same company could appear on two ATSs with
different IDs, and because IDs are only unique within a source. `job_id` is a
string, since Lever uses UUIDs.

Seen postings are persisted in SQLite so a restart can't cause either a replay
of everything already seen (duplicate spam) or a silent skip of whatever
arrived while the bot was down (a missed opportunity, which §1 forbids).

Table `postings`, keyed on `(source, company, job_id)`:

| Column | Purpose |
|---|---|
| `source`, `company`, `job_id` | identity |
| `title`, `location` | the fields classification is based on |
| `url` | the link included in the Discord message; not read for classification, kept for possible future audit/history features |
| `posted_at` | the source's own posting timestamp |
| `updated_at` | the source's own last-modified timestamp; intended to drive reclassification, see §7.1 |
| `verdict` | `ACCEPT` / `REVIEW` / `REJECT` |

### 6.1 Driver

`modernc.org/sqlite` (pure Go), not `mattn/go-sqlite3` (requires cgo). The
container builds with `CGO_ENABLED=0` onto a static base image, which cgo would
not link against.

The database file must be mounted on a volume. An ephemeral database in a
container gives none of the restart safety this section exists to provide.

## 7. Reclassification rules (v1)

- Not yet in the table → classify.
- Already in the table → skip. `Dedupe` does a plain existence check on
  `(source, company, job_id)`; no verdict currently reaches `Filter` or
  `Sender` twice.

`REJECT` postings are never written to the table (§8.2), so a `REJECT`ed
posting is reclassified from scratch on every poll rather than being
permanently excluded. That's redundant work, not a correctness problem — the
same out-of-scope signal is still present, so the outcome doesn't change.

### 7.1 Deferred: one-way `REJECT`, reclassifiable `REVIEW`

The intended long-term behavior is stricter than v1's plain existence check:

- `REJECT` as a one-way door: persisted immediately at classification time,
  independent of `Send` (since `REJECT` postings are never sent), so it is
  never reconsidered even if the posting is edited later. An employer whose
  posting carried an explicit out-of-scope signal would be responsible for
  correcting it; the bot wouldn't re-check.
- `REVIEW` as the only verdict that gets re-evaluated: if `updated_at` moves on
  a later poll, the posting is reclassified from scratch. This is how a
  posting that was ambiguous (missing signal) gets a second chance once the
  employer adds the missing information.
- `ACCEPT` stays terminal either way — already delivered, so a later edit
  changes nothing actionable.

`updated_at` is the source's timestamp, not ours. It can bump for edits that
don't affect classification, and in principle could fail to bump on an edit
that does. The first costs one unnecessary reclassification; the second costs
a missed second chance on a posting already sitting in `#unsorted` where a
human can see it. Both would be cheap relative to the alternative, which was
fetching every description on every cycle (§5.2).

**Not implemented in v1.** Shipping the simpler existence-check `Dedupe` and
deferring this was a deliberate choice to get a working pipeline out first
rather than block on it. Building it later needs two things: a second
`Storer.Store` call site right after `Filter` (to persist `REJECT`s
immediately, since they never reach `Send`), and an `updated_at` comparison
inside `Dedupe`'s query for postings previously verdicted `REVIEW`.

## 8. Delivery and write ordering

Each poll classifies every newly-seen posting, partitions them into the three
disjoint buckets, and composes messages per destination channel (`#internships`
for `ACCEPT`, `#unsorted` for `REVIEW`). `REJECT` postings are never sent
anywhere.

### 8.1 Batching

Postings are packed into as few messages as possible without exceeding
Discord's 2,000-character limit. Batching operates on **postings**, not on
pre-formatted strings, so each batch retains the set of postings that produced
it. Without that, a successful send cannot be attributed back to specific
postings and §8.2 is unimplementable.

Two invariants, both testable on a pure function:

- No batch is empty.
- Every posting handed to the batcher appears in exactly one batch.

Titles are truncated (on runes, not bytes) at format time. This makes an
over-limit single message unreachable by construction rather than something the
batcher has to handle after the fact.

### 8.2 Write ordering

`Storer.Store` is called once per poll cycle, after `Send`, with exactly the
postings that were actually delivered — `ACCEPT` and `REVIEW` postings from
batches that sent successfully (§8.1). If a send fails, the affected postings
are simply absent from that call and are naturally retried (reclassified and
resent) on the next poll. The absence of a row *is* the retry signal; no
separate failure tracking exists.

`REJECT` postings have no send step and are never passed to `Store` at all in
v1 — see §7.1 for why that's a deferred gap rather than a design choice.

All postings in one `Store` call are written in a single transaction, using
`INSERT OR IGNORE`. `Dedupe` already filters out anything previously seen, so
a primary-key collision reaching `Store` should be rare — a race, a bug, a
manual re-run — and `OR IGNORE` keeps that edge case from failing the whole
write rather than being the expected path. One consequence worth naming: this
commits per poll cycle, not per Discord message, so a hard write failure
partway through rolls back postings from multiple already-successfully-sent
batches together, not just the one that triggered it — those postings would
be resent next poll even though Discord already delivered them once.

Reversing send/persist order entirely would be worse regardless of commit
granularity. Persist-then-send leaves a row marking a posting "seen" that
nobody ever received, making it permanently invisible — the silent drop §1
forbids, and unrecoverable without manual intervention. Persist-after-send
cannot produce that state; the cost is a possible duplicate post instead of a
possible silent loss, which is the accepted trade per §1.

## 9. Decision log

| Date | Decision | Reason |
|---|---|---|
| 2026-09 | `REVIEW` is the zero value of `Verdict` | A posting that never reaches the filter must default to the safe bucket, not to a confident match |
| 2026-09 | Added a discipline axis | Two axes accepted `Marketing Intern, London`, which satisfies the spec as written but not the product |
| 2026-09 | Discipline is a veto, not a requirement | Requiring an explicit engineering word buried non-English titles such as `Praktikum Softwareentwicklung` |
| 2026-09 | Classification reads `title` and `location` only | Evaluating an axis against the description body caused both sides to fire on nearly every posting, disabling the axis |
| 2026-09 | Dropped `?content=true` | 81 MB per cycle for a field nothing reads after the field-discipline decision |
| 2026-09 | Dropped `content_hash`, replaced with `updated_at` | Cannot hash content that isn't fetched; `updated_at` answers the same question for free |
| 2026-09 | Dropped `raw_json` | Without descriptions the payload is a near-copy of the structured columns, and three sources would mean three incompatible schemas in one column |
| 2026-09 | Identity is `(source, company, job_id)`, `job_id` a string | IDs are unique per source, not globally; Lever uses UUIDs |
| 2026-09 | Batching operates on postings, not strings | §8.2 requires attributing a successful send back to specific postings |
| 2026-09 | Added Ashby as a second source | 19 of 27 checked companies not on Greenhouse were reachable on Ashby, including Deliveroo, Trainline, and Thought Machine |
| 2026-09 | Added a separate Dedupe service | Primarily to prevent duplicate messaging on jobs and helps to reduce load on filtering service. |
| 2026-09 | `Dedupe` runs right after `Poll`, sharing `Storer`'s DB connection | Skips wasted `Filter`/`Send` work on postings already seen; `Storer` stays sole writer, `Dedupe` sole reader, avoiding two independent DB owners |
| 2026-09 | `postings` table gains a `url` column beyond §6's original list | Cheap to store, kept for possible future audit/history features even though it's not read for classification |
| 2026-09 | Deferred `REJECT`-as-one-way-door and `REVIEW` reclassification (§7.1) | v1 ships a simpler existence-check `Dedupe` to get a working pipeline out first; the fuller behavior is documented, not dropped |
| 2026-09 | `Storer.Store` commits one transaction per call, `INSERT OR IGNORE` on conflict | `Dedupe` already filters upstream, so a collision reaching `Store` is a rare edge case, not the normal path — no need for per-posting atomicity, and `IGNORE` keeps that edge case from crashing the batch |