# The run

`an-answered-question-parks`, `opus` at `high` effort, round 1, against base `36a0968`, recorded on 2026-08-21.
It parked with 4 findings at `medium` risk, and it is recorded whatever it found.
`an-answered-question-second-look` is the same run's round 2, recorded the same day, and it is recorded whatever it found as well.
The corpus now holds eight answers and `internal/corpus/outcomes/distribution.txt` was rewritten to say so.

That is the whole of what this task may require of a live run.
The person decided both questions, the patchable finding was answered, round 2 ran against the patches, and the run ended `escalated` at the round cap.

Three things about the routing are worth stating because they are what this change was for, and none of them is a requirement on what either reviewer said.
Round 1 parked rather than ended, which is the behaviour this change adds, applied to the change that adds it.
Every round-1 answer in the corpus now parks, so the six that used to read `escalated` moved together and no first look reaches `escalated` any more.
And the distribution holds a real `escalated` entry for the first time, reached by round 2 at the cap rather than by a fixture.
The remaining `escalated` path, a finding that holds the change while its own action asks for nothing, is still reached by fixtures alone, and whether a real reviewer ever returns one is not this repository's to arrange.

## What the run cost to get

The first reviewer process died at once: `You cannot use --strict-mcp-config when an enterprise MCP config is present`.
The retry ran without the flag and the run recorded that the machine's tool servers were loaded, which is the documented behaviour and the reason the seal rather than the flag is what holds the reviewer read-only.
Nothing else about the run needed a person.

## What the dialogue was like from the inside

This is the first run where the branch this change adds was walked rather than tested.
Two findings arrived actioned `ask-user`, `scrat status` said `of those, needing a person rather than a patch: 2`, and `scrat continue` refused while either was unanswered and named them.
The session did not answer them: it put both in front of the person, and the person decided both should be fixed inside this change rather than deferred.
`--fixed` would have been refused against either, which is the whole of what enforcement is left after the terminal state went away.

Under the old rule this run would have exited `1` at round 1 and the two findings below would have been handed over as questions nobody could answer inside the run.
The patches for them would have been new code that nobody read.

That is the part worth keeping from this run, because round 2 read exactly those patches and found the two tests written for round 1's findings could not fail on the counts they named.
The round the cap allows is not a formality here: it caught vacuous evidence for the fixes the round before it prescribed, which is the case the design argues the round exists for.

## Where the change was unfinished

The reviewer found the reporting around the two states the rule moved, which is exactly the seam this change cuts and exactly where its author stopped looking.

**The parked event says only how many findings can be patched.**
`presentation.roundParkedSentence` is "round %d is waiting on you: %d findings can be patched", and that count came from `classified.Patchable`.
Before this change a round could park on nothing else, so the count was a complete account of what the round waited on.
After it, a round can park on a question alone, and then `scrat watch` and `scrat logs` print zero and say nothing more.
`scrat status` gained `awaiting_a_person` in this change; the event a person watching the run reads did not.
The concrete path is a scenario already in the suite: `internal/e2e/readonly_test.go` returned one `warning`/`ask-user` finding, and that run parked with nothing patchable.

Fixed: the event carries `needing_a_person` beside `patchable`, the renderer names it in a second sentence where that count is above zero, and `patchable` is now the blocking count the round actually waits on rather than every patchable finding, informational ones included.

**The escalated surface lost its account of why it escalated.**
The skill's rewritten branch says `scrat status` prints which of the two reasons applied.
It prints a positive marker for one of them: `EscalatedByTheCap` appends the cap sentence.
For the other, the verdict reads "a reviewer reported 1 findings, 0 of them for the person", the counts line reads "0 for the person, 0 to patch, 1 informational", and `report.Audience` labels the finding that caused the escalation "informational".
So a session is told to hand everything to the person while every number says nothing is for them.
The sentence held before this change because the non-cap escalation reason used to be an `ask-user` finding, which the person-count did report.
It does not hold now, and `whatAReviewerFound` counting `NeedingAPerson` is a count of the findings that no longer escalate anything.

Fixed: `verdictReason` names the other reason too, before the cap clause, so each escalation carries a positive account of why it escalated.
The counts line still mirrors the classification, which is action-based on purpose, so an `error` finding actioned `no-op` still reads as `error` and `informational` side by side.
That contradiction is the reviewer's own and socrates reports it rather than resolving it.

## Where the checks were weaker than they read

**Two seal assertions stopped being able to fail.**
`TestACleanReviewerPassesTheCheck` and `TestAWriteTheSealRefusedIsNotAFindingAgainstTheReviewer` were moved from `WaitForTerminal` to `WaitFor(name, AwaitingResponse)` when their `ask-user` finding stopped escalating.
A parked run carries no verdict, `reason(run)` returns the literal "no verdict", and both tests then check that the verdict does not name a write.
Neither check can be true for any input.
A regression that wrongly blamed the reviewer would end the run `failed` and be caught by `WaitFor`'s own `Fatalf`, never by the assertion whose message states what the test is for.

Fixed: both tests are back on `WaitForTerminal` against a scenario that escalates, so a verdict exists for the assertion to read, and both also assert that no recorded event names the write.
Round 2 found that the event checks were unreachable where they were placed, and the section below records what that took to finish.

**Two script assertions were kept alive by padding.**
The escalating scenario in `exit-codes.txtar` and `reading.txtar` was rewritten to escalate through a blocking `no-op`, and a second `info` finding with no stated action was added beside it so that `stdout 'needs a person'` and `stdout 'for the person'` still matched.
They match on the padding.
The finding that causes the escalation renders as "informational", the added one is non-blocking and is not why the run escalated, and `'for the person'` would match a count of zero.
Both scripts pass and neither would notice the escalated surface losing its account of why it escalated, which is the finding above.

Fixed: the padding finding is gone from both scenarios, and both scripts now assert on the sentence that says why the run escalated.
The check and the thing it is about are the same thing again.

## What the second round found, and what was done

It reported five findings and the run escalated at the cap.
Three of them were about the round before it: the evidence offered for round 1's fixes, rather than the fixes.

**The unit test for the new count could not fail on it.**
`presentation.Lines` prefixes every line with the timestamp of the test's fixed clock, `09:14:02`, so a rendered line already holds the characters `0`, `9`, `1`, `4` and `2`.
The subtests asserted `Contains(line, "2")` and `Contains(line, "1")` for counts of 2 and 1, and both were satisfied by the timestamp alone.
The reviewer named the counterexample: pass `patchableField` where `asked` goes, a plausible copy of the neighbouring argument, and the whole test stays green with the count wired to the wrong field.
Fixed by asserting the noun with the number, `"2 need a person"`, and that mutation was run: it now fails, twice.

**Two `person` assertions were satisfied by an unrelated line.**
`awaitingLine` reads "awaiting responses: socrates is waiting for the author, or for a person the author asks", so `Contains(printed, "person")` is true of every parked run.
The two e2e checks written for the new count would have passed with `awaitingAPersonLine` and `Status.AwaitingAPerson` deleted outright.
Both now assert `"needing a person rather than a patch: 1"`.

**The seal tests' event checks never executed.**
They sat after the `t.Fatalf` on the run's state, and a wrongful blame ends the run `failed`, so the state check always fired first.
The mutation that was claimed as proof was in fact satisfied by the state assertion rather than by the assertions whose messages say what the tests are for.
The event and verdict checks now come first, and the mutation fails all three lines in each test.

**The retrospective described fixed defects in the present tense.**
Round 2's snapshot was taken before the resolutions were written into this document, and the reviewer was right about the snapshot it read.
The archived retrospectives state each problem and then what was done, and this one now does too.

## What is left, and why it is left

One finding asks a person and the round cap ran out, so it is recorded here rather than answered.

**The escalated surface still contradicts itself on the counts.**
For the one non-cap escalation reason, `scrat status` prints a verdict saying the finding goes to the person immediately above a counts line saying zero findings are for the person, and `scrat findings` labels it `informational`.
The verdict now names the reason, which is what round 1's decision asked for, and the counts still mirror the classification, which is action-based by design and shared with the prompt's own definitions.
Making a blocking `no-op` read as the person's would mean `report.Counts` no longer mirroring `findings.Classification`, and the reviewer itself called it a product call.
It is the person's to make, there is no round left to judge it, and a fresh review is what a change for it would get.
