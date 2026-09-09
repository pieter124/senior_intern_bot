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
| `posted_at` | the source's own posting timestamp |
| `updated_at` | the source's own last-modified timestamp; drives §7 |
| `verdict` | `ACCEPT` / `REVIEW` / `REJECT` |

### 6.1 Driver

`modernc.org/sqlite` (pure Go), not `mattn/go-sqlite3` (requires cgo). The
container builds with `CGO_ENABLED=0` onto a static base image, which cgo would
not link against.

The database file must be mounted on a volume. An ephemeral database in a
container gives none of the restart safety this section exists to provide.

## 7. Reclassification rules

- Not yet in the table → classify.
- `REJECT` is a one-way door: never reconsidered, even if the posting is edited
  later. An employer whose posting carried an explicit out-of-scope signal is
  responsible for correcting it; the bot doesn't re-check.
- `ACCEPT` is also terminal: it has already been delivered, so a later edit
  changes nothing actionable.
- `REVIEW` is the only verdict that gets re-evaluated. If `updated_at` moves on
  a later poll, the posting is reclassified from scratch. This is how a posting
  that was ambiguous (missing signal) gets a second chance once the employer
  adds the missing information.

`updated_at` is the source's timestamp, not ours. It can bump for edits that
don't affect classification, and in principle could fail to bump on an edit
that does. The first costs one unnecessary reclassification; the second costs a
missed second chance on a posting already sitting in `#unsorted` where a human
can see it. Both are cheap relative to the alternative, which was fetching
every description on every cycle (§5.2).

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

A posting's row is written to SQLite only **after** its message has sent
successfully, committed per message rather than per cycle. If a send fails, the
affected postings are simply absent from the table and are naturally retried
(reclassified and resent) on the next poll. The absence of a row *is* the retry
signal; no separate failure tracking exists.

`REJECT` postings have no send step, so they're persisted immediately on
classification. Persisting them is what makes §7's one-way door enforceable:
without the row, a rejected posting looks unseen on the next poll and gets
reclassified.

The unit of knowledge is the batch, not the posting: Discord accepts or rejects
a whole message, so all postings in a batch persist or none do.

Reversing the order would be worse. Persist-then-send leaves a row marking a
posting "seen" that nobody ever received, making it permanently invisible. That
is the silent drop §1 forbids, and it is unrecoverable without manual
intervention. The current ordering cannot produce that state.

The cost is a possible duplicate post if a crash lands in the narrow window
after a successful send but before the commit. Given §1's stated priority, an
occasional duplicate is an acceptable trade for never silently losing a
posting.

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