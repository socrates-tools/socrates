# What the seed already says

The corpus opens with four entries: two looks at the dialogue layer and two at the fixes it needed, each one round, `opus` at `high` effort, recorded on 2026-08-20.
Between them they hold 28 findings, and they are the first answers from a real reviewer this project has ever kept.

Every number below is computed from `internal/corpus/answers`, not remembered.

## The 28 findings

| severity | action | findings |
| --- | --- | --- |
| warning | ask-user | 12 |
| warning | auto-fix | 9 |
| info | no-op | 4 |
| info | ask-user | 2 |
| info | auto-fix | 1 |

By severity: 21 `warning`, 7 `info`, and **no `error` at all**.
By action: 14 `ask-user`, 10 `auto-fix`, 4 `no-op`.
Every finding named an action; not one left the field out, so `EffectiveAction`'s fallback has never been needed by a real answer.
All four documents rated the change `medium` risk: no reviewer has yet said `low` or `high` about socrates.

## What no real reviewer has produced

- **`error` severity.** Not once in 28 findings. The end-to-end journey that drives a parked round opens with an `error` (`internal/e2e/dialogue_test.go`, `patchableFindings`), and so does the journey that escalates, and so does the rejection that a fresh reviewer finds again.
- **`{warning, no-op}`.** The shape that `dialogue.After` was rewritten for twice in one afternoon appears in no recorded answer, and in no scenario in the suite either. Only the exhaustive sweep in `internal/dialogue/dialogue_test.go` covers it, which is exactly the coverage that a rule needs and a journey cannot give.
- **A parked round.** Every one of the four rounds carried at least one `ask-user` finding, so every one escalated at round one. `awaiting-response` has never been reached by a real answer.
- **`clean`.** No real reviewer has yet found nothing blocking in socrates. The clean journey rests entirely on an invented answer.

`make corpus-outcomes` now says all of that in one place, beside the number of scenarios in the suite that assert each outcome.

## Two `{info, ask-user}` findings, and nobody decided that

Twice in 28 findings a reviewer marked a finding `info` and actioned it `ask-user`.
`dialogue.After` escalates a run for either one, whatever else is in the set, because `ask-user` dominates the rule.

So a reviewer can end a run as escalated over a finding whose own severity says it holds nothing.
That may be right: the reviewer chose to put a question to a person, and severity is about the code while action is about who answers.
It may also be a reviewer using `info` to mean "not blocking" and `ask-user` to mean "your call", which the rule then reads as blocking.

Nobody has decided this on purpose, and this change does not decide it.
It is recorded here as an open question, with the evidence attached, for whoever changes the rule next.
The distribution file is what will show whether the answer moved anything.

## The task that could not be satisfied

`a-rejection-is-a-claim` shipped with task 9.4 open: a live dialogue including one genuine rejection.
Four live reviews, four escalations at round one, and the same cost each time.
The task was unsatisfiable as written, because it asked a model to return a particular shape of answer.

That is now a requirement rather than a lesson: a task about a live run may require the run, the recording and the routing, and may not require a severity, an action, a count or a shape.
Task 9.4 was reworded to what a live run can satisfy, and it closed on the rewording, not on a dialogue having happened.

## What this change had to reword about itself

Task 2.3 asked for the no-repair proof to use "an entry whose findings would all be renumbered by a helpful implementation".
No such entry exists: all four recorded answers number their findings `f1` upward with no gaps, so a renumbering simulator would be invisible.

Inventing an entry to make the task satisfiable is the exact move this change exists to refuse, so the proof is the one the corpus can carry: the emitted document is byte-identical to the recorded one, and it is byte-identical for two prompts that share no vocabulary, which no adapting or renumbering implementation manages.
The task now says that.
