**The rule comes first and it is tested over its whole input space.** `dialogue.After` is the most load-bearing rule here, and every task below either changes it, feeds it, or reads what it decided. A task that finds itself special-casing a state has found the rule in the wrong place.

**Nothing here lets socrates write in a reviewed repository.** The author still applies every patch, and no task adds a second writer.

**The enforcement this change gives up is stated, not hidden.** A terminal state used to stop a session answering an intent question. After this, what stops it is which resolution the finding accepts, a required reason a fresh reviewer reads, and a sentence in the skill. Task 5.3 is the test that the sentence is there.

**`make lint test` after every task, not at the end of a group.**

## 1. The rule

- [x] 1.1 Add the third resolution to `internal/dialogue` with its required reason, and refuse a decision for a finding that asks for a patch and a patch for a finding that asks a person. Covers `Recording a decision`, `A decision with no reason`, `A patch offered where a decision is required`, `A decision offered where none was asked for`
- [x] 1.2 Change `Answerable` to admit findings whose action asks a person, and keep it the one owner of what a round waits on. Covers `A question alone parks the run`
- [x] 1.3 Change `After` so a round holding a finding that asks a person parks rather than escalates, with the whole-input-space test rewritten for the new rule. Covers `A question about intent stops the loop at once`, `Mechanical findings park for the author`
- [x] 1.4 Keep the cap escalating, including a question left unanswered when the cap ran out. Covers `A question left unanswered at the cap escalates`, `The cap escalates`
- [x] 1.5 Confirm a blocking finding actioned `no-op` still escalates and is waited on by nothing. Covers `A blocking finding nobody can patch goes to the person`, `A clean round finishes the run`, `An informational round finishes the run`, `A clean second round does not name the cap`, `The rule reaches nothing outward`
- [x] 1.6 Mutate the rule: make a question escalate again and confirm the parking tests fail, and make it park at the cap and confirm 1.4's test fails

## 2. Recording a decision

- [x] 2.1 Add `--decided` to `scrat respond`, refusing the flags of the other resolutions with a message naming the one that applies. Covers `A flag that does not apply to the answer`, `One finding per invocation`
- [x] 2.2 Read a decision's reason from standard input on the same terms as a rejection's. Covers `A decision's reason from standard input`, `A reason from standard input`, `A reason is read only where it is asked for`
- [x] 2.3 Confirm the append-only record and the replacement rule hold for a decision. Covers `Recording a fix`, `Recording a rejection`, `A rejection with no reason`, `A response to a finding that does not exist`, `Changing a response`, `Responding to a run that is not awaiting responses`
- [x] 2.4 Make `scrat continue` refuse while a question has no decision, naming which findings need a person. Covers `Continuing with a question unanswered`, `Continuing too early`, `Continuing when everything is answered`, `Non-blocking findings need no response`, `Nothing left to review ends the run`

## 3. What a parked run says

- [x] 3.1 Add the count of unanswered findings needing a person to what `status` prints and to the JSON object. Covers `A parked run says what needs a person`, `A parked run reports what it wants`, `The count reaches zero`
- [x] 3.2 Add the fifth response rule to the one constant both surfaces read, and confirm `scrat findings` prints it where a round parks and nowhere else. Covers `The rules arrive with the findings`, `The rules do not arrive where they do not apply`
- [x] 3.3 Update what `awaiting-response` means in the run store's own words, and confirm a run parked on a question reports as still working. Covers `A run parked on a question is not a verdict either`, `A parked run is not finished`, `A terminal run does not move`, `A failed run reports no findings as its verdict`, `A parked run can still be cancelled`

## 4. What the next reviewer is told

- [x] 4.1 Carry a decision into the prior-round block labelled as a decision, through the same per-field neutralisation every other model-facing field goes through. Covers `A decision arrives as a decision`
- [x] 4.2 Add the decision sentence to the provenance clause, asserted on the delivered prompt rather than on socrates' source. Covers `The clause reaches a re-review`, `The clause is absent from a first review`, `A claimed fix is not taken on trust`

## 5. The skill

- [x] 5.1 Split the `awaiting-response` branch into what the session answers and what the person answers, and regenerate the committed copy
- [x] 5.2 Say that the session does not answer a finding that asks a person, and that it records what the person decided. Covers `A question is the person's`
- [x] 5.3 Assert the sentence is present, and mutate it away to confirm the test fails, because this sentence is what replaced a terminal state
- [x] 5.4 Say that the session continues once everything answerable is answered, without asking whether to spend the round, and that a finished run is reported with how many rounds it took and what the last one found. Covers `A round is spent without asking`, `What the person is told at the end`
- [x] 5.5 Say that there is no third round, that a further look is a fresh review, that the person asks for it, and that its reviewer knows nothing of the run before it. Covers `There is no third round`, `A fresh review is the person's to ask for`, `A new run knows nothing of the old one`
- [x] 5.6 Confirm the skill still names every state a run rests in and every command and flag it names exists, now that `--decided` is one of them. Covers `Every terminal state is covered`, `A failed run is stated not to be a pass`, `A cancelled run is not a verdict`, `Every state a run rests in is named`, `The loop closes`

## 6. End to end

- [x] 6.1 Drive a journey where a round asks a question, the answer is recorded as a decision, the next round judges it, and the run ends
- [x] 6.2 Drive a journey where a round holds a question and a patchable finding, and confirm the run waits on both and refuses to continue until both are answered
- [x] 6.3 Drive a journey where a question reaches the cap unanswered and confirm the run escalates with everything kept
- [x] 6.4 Assert across every journey that nothing was written in the repository under review

## 7. What moved because the rule moved

- [x] 7.1 Rewrite `internal/corpus/outcomes/distribution.txt` with `scripts/corpus-outcomes -write`, and read the diff: every recorded answer holding a question moves from `escalated` to `awaiting-response`
- [x] 7.2 Update `README.md`: the exit code a question now returns, what a parked round can be waiting for, and that a further look after the cap is a new run rather than a third round
- [x] 7.3 Update `AGENTS.md`: the routing invariant no longer says `ask-user` dominates by ending the run, and the invariant that a decision is a claim rather than a shortcut

## 8. Verification

- [x] 8.1 `make fmt && make lint test` clean
- [x] 8.2 `make scenario-coverage` naming no scenario this change adds
- [x] 8.3 `make scenario-drop` before archiving, for the nine requirements this change substitutes
- [x] 8.4 `make corpus-outcomes` read rather than skimmed, because this change is the reason its numbers moved
- [x] 8.5 A live run: review this change with the real CLI, record whatever it answers, and write `retrospective.md` from it. Require the run, the recording and the routing, and require no severity, action or count
