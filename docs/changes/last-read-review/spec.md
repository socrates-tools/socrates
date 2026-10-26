# The last read puts socrates review in front of the whole change

Decisions: ADR 90, ADR 91.

## Problem Statement

`socrates-implement` reviews every ticket's diff as it lands, and once the last one has landed it reads the whole change against the spec's user stories.
That last read is the spec count only.
Nothing reads the whole change against the repository's standards, so what shows only across tickets has no reader: the same helper written twice by 2.1 and 2.2, or one concept spread over five packages, each half of which passed its own ticket's review.
On a long run the overseer's own context has compacted by the time it gets there, and the reader it lacks is one that wrote nothing and read nothing else.

The suite already ships that reader, socrates review, but ADR 49 kept it out of the overseer's run.
socrates review also cannot finish a run with no person present: a finding that asks a person accepts only `decided`, which says a person answered, so such a finding parks the run and the second round never happens.

This change grew out of a comparison with `mattpocock/skills`, whose `implement-spec` ends by running `code-review` over the integration branch and handing everything it raised to one implementer, asking the person nothing.
It is G4 of that comparison, and G1, G2, G3 and G7 landed already in `bf9a475`.

## Solution

Once the last ticket has landed and passed, the overseer runs `socrates review` over the whole change, from the default branch's fork point, per ADR 90.
She begins the run with `--nobody-to-ask`, an intent that points at the spec, and her own host, and walks the spec's user stories while the reviewer reads.
What the review asks to have fixed and the stories the change falls short of go out together to one fresh implementer, are landed and verified the same way a ticket is, and the second round judges the patched branch.

A finding that asks a person is settled by the overseer rather than put to the person.
socrates review gains a resolution for that, `settled`, which records that the session decided and never claims a person did, per ADR 91.
Only a run begun with `--nobody-to-ask` takes it, and such a run never takes `decided`.

Every skill stays a person's to start.
The overseer drives review's commands and never the `socrates-review` skill.

The hand-back says how the review ended.
Where it ended clean, it says so; where it ended escalated or failed, the branch is still handed back, with the verdict and each finding still open or the failure the reviewer gave, and no third run.
Each finding the overseer settled is listed with its reason, so the person sees what was decided while they were away.

## User Stories

1. As a person handing a change's tickets to `socrates-implement`, I want an independent reviewer to read the whole change once every ticket has landed, so that what shows only across tickets, like one helper written twice under two names, is caught before I publish the branch.
2. As a person running a long implementation, I want that whole-change review done by a reader that wrote none of the code and holds none of the run's history, so that its judgement does not depend on how much of the spec survived in the overseer's compacted context.
3. As a person who left a run alone, I want the review, the fixes and the second round to happen without asking me anything, so that I come back to a finished branch rather than a question that stalled it.
4. As a person reading the hand-back, I want to be told whether the whole-change review ended clean, escalated or failed, so that I know how far the branch was checked before I decide whether to publish it.
5. As a person whose run ended escalated, I want each finding still open listed at hand-back, so that I can judge it myself rather than discover it in a pull request.
6. As a person whose reviewer CLI is not installed, not signed in or failed, I want the branch handed back anyway with the failure the reviewer gave, so that a broken reviewer costs me the review and not the run.
7. As a person reading the hand-back, I want every finding the overseer settled in my absence listed with what it decided and why, so that I can reverse a decision I disagree with before publishing.
8. As a person reading a run's record later, I want a decision the session made to be recorded as the session's and never as mine, so that the record does not claim I decided something I never saw.
9. As a person running `/socrates-review` myself, I want my session never able to settle a question the reviewer put to me, so that what asks me is always answered by me.
10. As a person, I want the whole-change review and the story read to go to one implementer together, so that one hand that sees both halves of a cross-ticket problem fixes it.
11. As a person, I want the fixes from the last read landed and verified the same way every ticket was, so that the branch I am handed passed the repository's verification after its last change.
12. As a person, I want the reviewer told that the change builds the spec, so that it judges the code against what I asked for rather than only against the repository's standards.
13. As a person, I want no skill to start another skill, as ADR 49 promised, so that the overseer running review's commands never turns into the overseer starting `socrates-review`.
14. As a fresh reviewer in a second round, I want a settled response to read as the session's decision with its reason, so that I judge it as a claim rather than as a person's ruling.
15. As a contributor, I want a run's presence recorded once at `begin`, the same way whose change it is is recorded, so that which resolutions a finding accepts is decided in one place.
16. As a contributor, I want the suite to fail where the overseer's skill names a review command the binary does not offer, or names any skill, so that the skill cannot drift from the binary or from ADR 49.

## Implementation Decisions

**review: who is present**

- `socrates review begin` gains `--nobody-to-ask`, recorded on the run in `run.json` once, when the run begins, and read by every later command.
  It is refused together with `--not-my-change`, before a snapshot is taken and before a reviewer is spawned: a change nobody here can patch waits on a person by definition, so the two cannot both hold.
- `review/dialogue` gains the resolution `Settled` beside `Decided`, and the presence a run records beside its `Authorship`.
  `dialogue.Accepts` stays the one rule for which resolutions a finding takes: on a run with nobody to ask, a finding that asks a person accepts `Settled` alone, and on any other run that is mine it accepts `Decided` alone.
  `dialogue.Answering` refuses `Decided` on a run with nobody to ask and `Settled` on any other, naming the resolution the finding does take.
- `socrates review respond` gains `--settled`, which needs `--reason`, as `--decided` does, since a settlement leaves no artifact behind, per ADR 5.
- `dialogue.AwaitsAPerson` reads the presence: on a run with nobody to ask, a round left `awaiting-response` by findings that ask a person does not wait on a person, because the session answers it, so `begin --wait` returns there rather than reporting that it parked on somebody.
- The second round's prompt names a settled response as the decision of the session that started the run, beside a decided one named as a person's.
- `status`, `findings --json` and the run's events report that a run has nobody to ask and that a response was settled.
- `dialogue.ResponseRules`, rendered into the socrates-review skill, says `settled` belongs to a run begun with `--nobody-to-ask`, which that skill never begins.
  The socrates-review and socrates-review-handover skills do not pass `--nobody-to-ask`.
- Nothing in agora changes: the overseer serves no page for her review, and review's own commands serve none.

**implement: the last read**

- `implement/skill`'s last read section is rewritten, in each host's rendering, to say:
  - Once the last ticket has landed and passed, begin the review on the change's branch, in the person's own checkout, with `socrates review begin --base <default branch> --intent "This change builds docs/changes/<slug>/spec.md" --nobody-to-ask --host <host>`, running in the background.
  - Walk the spec's user stories over the whole diff while the reviewer runs.
  - Read what the round found with `socrates review findings --json`.
  - Settle each finding that asks a person by her own judgement, reject each she judges false, and hand every fix, settled answers included, to one fresh implementer together with the stories the change falls short of, in a worktree of its own made the same way as a ticket's.
  - Land and verify its work the same way as a ticket's, then record a response to every finding with `socrates review respond`, using `--fixed`, `--rejected` or `--settled`, and run `socrates review continue`.
  - A round one that ended clean has no second round, and its story shortfalls still go to the one implementer.
  - Where the run ends escalated or failed, start no third run and hand back anyway.
  - The hand-back states the verdict, every finding still open or the failure the reviewer gave, and each settled finding with its reason, beside what it already lists.
- The overseer records review responses and runs review's commands herself, because that is reading and recording rather than writing code.
  The implementer's agent definition and Codex role stay free of every `socrates` command.
- `implement/skill/overseer_test.go`: `TestTheOverseerSkillNamesNoSocratesReviewInvocation` is replaced by a test that the rendered skill names review's commands at the last read and names no skill, `/socrates-review`, `$socrates-review` or otherwise.
  `TestTheOverseerSkillNamesNoNamespaceButImplement` widens to implement and review.
- `implement/AGENTS.md`'s invariants change to match: implement knows nothing about review's skills, agora or design, and its skill names `socrates implement` commands and, at the last read, `socrates review` commands.
  The last read's invariant gains the review, one implementer for both lists, the escalated and failed hand-back, and the settled list.
- `review/AGENTS.md` gains the invariants ADR 91 states: a run records whether anybody is there to ask once, at `begin`, `dialogue.Accepts` offers `settled` only where nobody is and `decided` only where somebody is, and `--nobody-to-ask` is refused beside `--not-my-change`.
- `OverseerCeilings` and review's skill ceilings rise where the rendered skills outgrow them, and only by what the new text needs.

## Testing Decisions

- The review half is tested at review's end-to-end dialogue seam, `internal/e2e/dialogue_test.go`, through the real binary against every vendor's simulator.
  One journey begins with `--nobody-to-ask`, gets a finding that asks a person, answers it with `respond --settled --reason`, runs `continue`, and reads `findings --json` reporting the response as settled and the run as having nobody to ask.
  Further scenarios refuse `--decided` on a run with nobody to ask, refuse `--settled` on a run without it, and refuse `--nobody-to-ask` beside `--not-my-change` before any snapshot is taken.
- The second round's prompt is asserted on as delivered, naming a settled response as the session's decision, per the suite's rule that tests read the delivered prompt rather than review's source.
- The implement half is tested at the skill's own seam: the rendered skill, in each host's rendering, states each promise of the last read above, and every `socrates review` command and flag it names exists in the real binary, read the way the skill's `socrates implement` commands already are.
- Each new or replaced guard is mutated once to watch it fail: a skill naming `/socrates-review`, and `dialogue.Accepts` offering `decided` on a run with nobody to ask.

## Out of Scope

- Reviewing each ticket in a fresh context rather than in the overseer's own, G5 of the comparison, which waits until a run shows a late landing reviewed against a forgotten ADR.
- A smell baseline for repositories that state no standards, G6, which belongs in review's prompt beside `review/testquality` if it ever comes.
- Wide-refactor batches that cannot land green alone, G8.
- A third review run after an escalated one, and any page served for the overseer's review.
- The socrates-review skill settling anything, or beginning a run with `--nobody-to-ask`.
- Publishing the branch, which stays the person's.

## Further Notes

The comparison this came from, `mattpocock/skills@5beb35b` against `socrates-tickets` and `socrates-implement`, was a page outside the repository, and this spec carries what it needed from it in the Problem Statement.

Every option the person chose in this design was the one marked as the recommendation, and nobody wrote words beside a choice.
The answers that stay inside this spec rather than in an ADR:

- Who fixes what the review finds: one fresh implementer takes every finding and every story shortfall together, rather than one implementer per finding up to `concurrency`.
- The order: the story read and the review both read the landed change, and their shortfalls go out as one hand-out, rather than stories first and the review over the result.
- An escalated or failed review: hand back anyway, saying the review did not pass and why, rather than stopping before the hand-back to ask in the terminal.
- The intent: a pointer to the spec, rather than the spec's Problem Statement verbatim or no intent at all.
- The tests: review's end-to-end dialogue seam plus the skill's statement and command checks, rather than unit seams alone.
