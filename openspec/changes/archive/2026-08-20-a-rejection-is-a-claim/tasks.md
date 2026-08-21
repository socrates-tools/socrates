**Nothing here patches anything.** socrates records what the author claims and starts the next round; the author does the patching, in their own session, in their own working tree. A task that finds itself wanting socrates to apply a fix is a task in a tool this is deliberately not.

**The inversion is task 1.1 for a reason.** Every other task in this change is safe to get wrong and easy to fix. Getting the response rule backwards produces a tool that reads as if it works and quietly lets an author dismiss every finding, and there is no test that fails.

**`make lint test` after every task, not at the end of a group.**

## 1. The routing rule, as a pure core

- [x] 1.1 Add the routing rule as a pure function of the classified findings and the round number, with `ask-user` dominating. Test it against its whole input space rather than a sample. Covers `A clean round finishes the run`, `A question about intent stops the loop at once`, `Mechanical findings park for the author`, `The cap escalates`
- [x] 1.2 Add the test that its package imports nothing outward. Covers `The rule reaches nothing outward`
- [x] 1.3 Add the cap as a constant with its reason recorded beside it, and the test that nothing configures it. Covers `A second round is permitted`, `A third round is not`, `No setting changes it`

## 2. The run parks

- [x] 2.1 Add `awaiting-response` to the run states and to terminality. Covers `A parked run is not finished`, `A terminal run does not move`, `A failed run reports no findings as its verdict`
- [x] 2.2 Exempt `awaiting-response` from the liveness rule. Write this test before the state is wired into anything: without the exemption every parked run reports as failed the instant its reviewer exits, which is immediately, and it looks like a race. Covers `A parked run with no processes is still parked`
- [x] 2.3 Exempt a parked run from the retention sweep, which the existing never-remove-a-live-run rule already does. Covers `A parked run survives the sweep`
- [x] 2.4 Add cancellation from parked. Covers `A parked run can still be cancelled`

## 3. Responses

- [x] 3.1 Add the response document, per round, and its storage. Covers `Recording a fix`, `Recording a rejection`
- [x] 3.2 Add `scrat respond`, one finding per invocation, with the reason accepted as an argument or on standard input. Covers `One finding per invocation`, `A reason from standard input`
- [x] 3.3 Add the refusals. Covers `A rejection with no reason`, `A response to a finding that does not exist`, `Responding to a run that is not awaiting responses`
- [x] 3.4 Add replacement of an earlier response and the record that it happened. Covers `Changing a response`
- [x] 3.5 Add the unanswered count to `scrat status`. Covers `A parked run reports what it wants`, `The count reaches zero`

## 4. Advancing

- [x] 4.1 Add `scrat continue` and the refusal on unanswered findings. Covers `Continuing too early`, `Non-blocking findings need no response`
- [x] 4.2 Add the per-round snapshot at the top of the round. Covers `Continuing when everything is answered`, `The second round sees the patches`, `Each round records what it reviewed`, `A round with nothing patched still snapshots`
- [x] 4.3 Keep the base fixed at the run's original base, and fail the round naming it where it has become unreachable. Covers `The base does not move between rounds`, `The second round sees the whole change`, `The base still exists when the person has moved on`
- [x] 4.4 Keep every round's artifacts. Covers `Two rounds, two directories`, `An escalated run keeps everything`
- [x] 4.5 Assert no session identifier, resume instruction or reference to a prior conversation reaches a later round's invocation. Covers `No session is carried`

## 5. What the second reviewer is told

- [x] 5.1 Add the prior-round block: findings and responses, attributed, labelled as metadata to be judged. Covers `The claims travel explicitly`, `The block says what it is`, `A later round carries the claims and nothing more`, `A response that reads like an instruction`
- [x] 5.2 Add the provenance clause, absent from the first round. Covers `The clause reaches a re-review`, `The clause is absent from a first review`, `A claimed fix is not taken on trust`, `The first round carries no claims`
- [x] 5.3 Add the test that no suppression instruction is ever emitted. Covers `No suppression instruction reaches the reviewer`, `A disputed finding may be found again`. It is the requirement most likely to be undone by somebody copying the equivalent prompt from a tool where it is correct

## 6. Hygiene, on both channels

- [x] 6.1 Extend neutralising to the response document, per field. Covers `Conflict markers do not survive`, `Line endings are normalised`, `Neutralising is per field`
- [x] 6.2 Assert the stored documents are verbatim and only the prompt copy is neutralised. Covers `The stored record is not neutralised`

## 7. What the author is told

- [x] 7.1 Add the four response rules to `scrat findings` where the run is parked, and only there. Covers `The rules arrive with the findings`, `The rules do not arrive where they do not apply`
- [x] 7.2 Update the command surface and its test. Covers `The command surface is what it says`, `Recording a response writes nothing in the repository`
- [x] 7.3 Any new sentence goes in its package's own registry, and the plain-language scan passes with nothing added to a permit list

## 8. End to end

- [x] 8.1 Drive the full dialogue: round one parks, the author patches and responds, round two finds nothing, the run is clean
- [x] 8.2 Drive the rejection accepted: round one parks, the author rejects with a reason, round two does not re-report it, the run is clean
- [x] 8.3 Drive the rejection refused: round one parks, the author rejects, round two re-reports it, the run escalates naming the cap. This is the journey the whole change exists for
- [x] 8.4 Drive the intent question: round one returns an `ask-user` finding among several `auto-fix` ones, the run escalates at once, and nothing is offered for a response
- [x] 8.5 Drive the author who edits between rounds, and assert the second round's snapshot has the edits and the base has not moved
- [x] 8.6 Assert across every journey that socrates wrote nothing in the repository beyond the snapshots' git objects
- [x] 8.7 Add the `testscript` scripts for `respond` and `continue`, including every refusal

## 9. Verification

- [x] 9.1 `make fmt && make lint test` clean
- [x] 9.2 `make scenario-coverage` naming no scenario this change adds
- [x] 9.3 `make scenario-drop` before archiving. Three requirements here are `MODIFIED` deltas that substitute a requirement wholesale, and a scenario a delta does not carry leaves the specification the moment the change is archived, with the suite still green
- [x] 9.4 A live run against the real agent on a real change, with whatever it answered recorded into the corpus and routed as this specification requires. Four runs on four real subjects, all four escalated at round one, all four kept in `internal/corpus/answers`. The rejection this task originally asked for is not something a task can ask a model for, which is what `what-a-reviewer-actually-answered` turned into a requirement
