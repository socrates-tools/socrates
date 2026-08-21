# What the seed already says

The corpus opens with five entries: four seeded from `a-rejection-is-a-claim` and one from this change's own live run.

The four seeded ones are two looks at the dialogue layer and two at the fixes it needed, each one round, `opus` at `high` effort, recorded on 2026-08-20.
Between them they hold 28 findings, and they are the first answers from a real reviewer this project has ever kept, which is why the numbers below are about them.

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

`make corpus-outcomes` now says all of that in one place, beside the number of tests that name each outcome.

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

## What the fifth entry found, live, in this change

Task 6.4 ran one live review of this change against the real CLI: `the-reviewer-corpus`, `opus` at `high` effort, round 1, 36 files against `e88de0b`.
It escalated with 7 findings, 3 of them for the person, at `medium` risk, and it is recorded whatever it found, which is what the new requirement asks of a live task.

It also found a defect that this change's own headline invariant was already failing.

**Replay was not byte-for-byte for three of the five entries.** The simulator emitted the result event with `json.Marshal`, which runs a `json.RawMessage` through HTML escaping, so every `<`, `>` and `&` in a recorded document was rewritten as `\u003c`, `\u003e` and `\u0026`.
Two of the four seeded answers quote a comparison (`round 1 < Cap`), and the fifth holds twenty such characters.
The byte comparison passed only because both tests named `dialogue-layer-first-look`, the one entry with none.
The test now runs over every entry in the corpus, and the simulator writes its stream with an encoder that does not escape, which is also what the real CLI does: the raw `<` in the seed is how it arrived.

The other six were about the recording path and the report, and all six are fixed:

- The model an entry was recorded with came from the last `agent-started` event in the whole log, so a run whose second round used a different model recorded round one as that model's answer. The log is now walked in order and `round-started` tells which round each `agent-started` belongs to; a round the log no longer holds a routing for is refused instead of borrowing a later one's.
- Recording accepted a document socrates never routed, which would have poisoned `make corpus-outcomes` and two tests on the next unrelated change while the command that created it printed success. Both the run reader and the corpus now parse the document, so the refusal reaches the operator.
- A name opening with a dot or an underscore was accepted, written, and then never seen again, because `go:embed` drops it. Names are now a single path element starting with a letter, and a test compares what is on disk against what the corpus can see.
- The report's second number counted test functions that *mention* a state, which includes tests asserting that a run does not reach it. The number is now labelled as naming rather than asserting, in the report, in `AGENTS.md` and in the specification, and `Naming` is tested against a fabricated tree rather than against the suite that contains the test.
- One assertion in this change's own tests could not fail: it counted tests naming an outcome by walking the file it was written in, and that file names two outcomes. It is gone.
- The `// scenario: A live review's answer is kept` marker sat on an end-to-end test that runs against the simulator, so `make scenario-coverage` reported a live scenario as covered by a simulated run while task 6.4 was still open. The marker moved to the test over the recorded corpus, whose entries are live answers by construction.

Two of those, taken together, are the argument for this change existing.
The escaping defect was invisible to a suite that named one entry by hand and would have stayed invisible for as long as nobody recorded an answer quoting a comparison.
The mislabelled column was a number about coverage that read as stronger than it was, in the report whose whole purpose is to stop exactly that.
