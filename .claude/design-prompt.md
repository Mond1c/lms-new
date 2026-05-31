# Design prompt — LMS web frontend

> Paste this whole block into a fresh Claude conversation (or design tool like
> Claude artifacts / v0). It's self-contained — no repo needed. Replace the
> **TASK** line at the bottom with the specific screen(s) you want.

---

You are a senior product designer + frontend engineer. Design (and, if asked,
build) the web frontend for the product below. Ask clarifying questions before
finalizing if anything is ambiguous. Keep it clean, fast, and developer-friendly.

## Product

A **learning-management platform for programming courses**. Students do
assignments in **their own Git repositories**; work is graded by a **combination
of automated tests, human code review, and (optionally) an oral defence**. Think
*GitHub Classroom + an autograder + a review/defence workflow*. Provider-agnostic
(Gitea / GitHub / GitLab). There is **no frontend yet** — greenfield. The backend
exposes a **gateway / BFF** the browser talks to.

### How submitting actually works (important — no "Submit" button)

A student does **not** click "submit". Instead:

1. They push to their repo and open a **pull request**.
2. Tests run automatically (external CI **or** our self-hosted runner) → a **test
   score** attaches to the commit/PR.
3. They **request review** from the course's **instructors' team** on the PR.
   *That request-review action is the submission.* It puts the work into the
   instructors' **review queue** as "needs review".
4. Reviewers comment **inside the PR**, then record grades in the LMS.

The UI must make this model obvious to students (guide them to "open a PR and
request review", not look for a submit button). A web-upload fallback may exist
but is **not** the primary path.

## Users & roles (a course can have MANY instructors)

- **Student** — sees courses/assignments, watches test results, requests review
  via their PR, follows review/defence status, reads feedback and final grade.
- **Instructor / reviewer** — works a **shared review queue**: claims a
  submission (so other instructors see it's taken), opens the PR, sets a
  **code-quality score (0..1)**, may **override the auto test score**, requests
  changes or approves, and (if required) records a **defence score**.
- **Admin** — user/role administration.

Roles are **per course** (a user can be a student in one course, instructor in
another).

## Grading model (design the UI around this)

A submission's grade has up to **three normalized components, each 0..1**:

- **tests** — from automated tests; **can be overridden** by an instructor.
- **quality** — set by the reviewing instructor (code quality).
- **defence** — set at the oral defence; only if the assignment `requires_defense`.

The **final grade** is computed by a **per-assignment grading policy**:

- Default: `final = (w_tests·tests + w_quality·quality) · defence` (weighted sum
  of tests & quality, multiplied by the defence score; defence = 1 when not
  required). Default weights ~0.7 / 0.3.
- An instructor may set **custom weights** or even a **custom formula** over
  `tests`, `quality`, `defence`.
- `final` (0..1) × the assignment's `max_score` = points shown.

The UI should always show the **breakdown** (tests / quality / defence → final),
not just a number, and indicate when a component was **overridden** or is still
**pending**.

## Submission state — three independent tracks

Don't model one status enum; a submission has three tracks, and the UI derives an
overall badge + queue filters:

- **Test:** pending → running → passed / failed (+ score).
- **Review:** not requested → **requested (needs review)** → **claimed (under
  review by X)** → changes requested / approved.
- **Defence:** not required → awaiting → passed/failed (+ score).

Queue filters map to these: **need review**, **under review**, **defence**, and
**by assignment** / by student.

## Core screens

### Student
1. **Login** — "Continue with GitHub/GitLab/Gitea" (VCS OAuth) **and** an
   email/password form; handle the OAuth redirect round-trip.
2. **Dashboard** — assignment cards per course: deadline (soft vs hard), test
   status, review status, defence status, current grade/breakdown, attempt count.
   Empty/loading states. Guidance on how to submit (open PR + request review).
3. **Submission detail** — VCS attribution with a deep link to the **PR/commit**;
   test results (expandable per-test table; **hidden** test cases shown but
   redacted); review status & reviewer comments summary (full thread lives in the
   PR); defence status; **grade breakdown** (tests/quality/defence → final).

### Instructor (a primary focus — currently missing)
4. **Course review queue** — *the* instructor workspace. A table/board of **all
   submissions in the course**, with:
   - **Filters/sort:** need review · under review · defence · by assignment · by
     student · by deadline.
   - **Claiming:** a "Take for review" action that **locks** the submission to
     that instructor and is **visible to everyone** ("Reviewing: Alice", with
     claimed-at / release). Prevents two instructors grading the same work.
   - Per row: student, assignment, test score, review/defence state, who claimed.
5. **Review screen** (open from the queue) — deep link to the PR; a panel to set
   **code-quality 0..1**, **override the test score** (with the original shown),
   leave a summary, and **approve / request changes**. If `requires_defense`,
   move to / record the **defence score**. Show the live **grade breakdown** as
   inputs change.
6. **Assignment authoring** — create/edit an assignment: template repo, deadlines,
   max score, **`requires_defense` flag**, **grading policy** (weights / custom
   formula), **runner choice** (external CI vs self-hosted).
7. **(Later) Course roster & grade overview** — per-student grades across
   assignments; export.

## Things to be thoughtful about

- **Multi-instructor coordination:** claiming must feel instant and unambiguous;
  show stale/abandoned claims; allow release/steal with attribution.
- **VCS-driven submission legibility:** make clear a submission came from a PR
  review-request; always link back to the exact PR/commit; comments live in the
  PR, grades live in the LMS.
- **Override transparency:** when an instructor overrides the test score, show
  original vs overridden and who/when.
- **Live test state** while tests run; clear failed/compile-error states.
- **Hidden test cases** — convey coverage without leaking them.
- **Grade breakdown** everywhere a grade appears.
- Soft **deadline** vs **hard_deadline**; late submissions.

## API the frontend talks to (gateway/BFF)

Some of this exists, some is **(proposed)** — design around the capabilities:

- `Login` (VCS OAuth code *or* email+password) → tokens + user_id.
- `RefreshToken`, `Whoami` (→ course memberships + per-course role).
- `MyDashboard` → assignment cards (deadlines, states, grade, attempts).
- `MySubmission(id)` → submission + test result + review/defence + final grade.
- **(proposed)** `ListCourseSubmissions(course, filters)` → review queue rows incl.
  `claimed_by`.
- **(proposed)** `ClaimSubmission` / `ReleaseClaim` (visible lock).
- **(proposed)** `SubmitReview(quality 0..1, optional test_override, approve|changes)`.
- **(proposed)** `RecordDefence(score 0..1)`; `OverrideTestScore`.
- `SubmitFromWeb` → rare fallback upload (not the primary path).

## Constraints / preferences

Modern, accessible, responsive. TypeScript SPA (React) over a ConnectRPC/JSON
gateway — but the visual/UX design is what matters here. Calm, information-dense,
developer-tool aesthetic (think Linear / GitHub).

## TASK

> ⬇️ Replace this line with what you want, e.g.:
> - "Design the instructor **course review queue** — table/board, filters
>    (need review / under review / defence / by assignment), and the claiming UX
>    with multi-instructor visibility."
> - "Design the instructor **review screen**: PR deep link, code-quality 0..1,
>    test-score override, approve/request-changes, and live grade breakdown."
> - "Design the **assignment authoring** screen incl. requires_defense, grading
>    policy (weights / custom formula), and runner choice."
> - "Design the student **Dashboard** and **submission detail** with the
>    tests/quality/defence breakdown."
> - "Propose the overall IA, navigation, and a component inventory for both
>    student and instructor surfaces."

Deliver: (1) clarifying questions if needed; (2) the flow; (3) ASCII/low-fi
wireframes; (4) component breakdown and key states (loading / empty / error /
live-updating / claimed-by-someone-else); (5) how each piece maps to the gateway
API above.
