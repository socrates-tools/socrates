**Nothing here changes what socrates does.** No verdict moves, no command gains a flag, and the routing rule is not touched. A task that finds itself editing `internal/dialogue` is a task in the wrong change.

**Task 1.1 is the seed, and it is first for a reason.** The four recorded answers were preserved out of a temporary directory to survive a reboot. Every other task in this change can be redone; those four cannot, because obtaining them cost four live reviews.

**The report must never fail for an unobserved outcome.** Getting that wrong produces a permanently red check, and a permanently red check teaches people to run the suite with it excluded, which is worse than not having written it.

## 1. The corpus

- [ ] 1.1 Seed the corpus from the four live rounds of `a-rejection-is-a-claim`, before anything else, from the four run directories preserved at `~/winhost/temp/socrates-live-answers`, each holding that run's `run.json`, `findings.json`, `events.jsonl` and delivered prompt. Each carries model `opus`, effort `high`, round 1, and its run's name, base and snapshot from its `run.json`. Covers `An entry holds the answer and its circumstances`
- [ ] 1.2 Add the entry type, its storage and its reader. Covers `An entry that names no model is refused`
- [ ] 1.3 Add recording from a completed run, and the refusal on any repository that is not socrates. Covers `Recording from another repository is refused`
- [ ] 1.4 Assert the stored document is byte-identical to what the reviewer returned, on an entry recorded through the recording path rather than on one written by the test

## 2. Replay

- [ ] 2.1 Teach the simulator to answer from a named entry, emitting the document verbatim. Covers `A scenario names a recorded answer`
- [ ] 2.2 Add the refusal on an entry the corpus does not hold. Covers `A named answer that does not exist is refused`
- [ ] 2.3 Assert the simulator neither adapts, renumbers nor completes a recorded document, using an entry whose findings would all be renumbered by a helpful implementation
- [ ] 2.4 Replace one hand-written end-to-end scenario with a replayed entry, and keep the hand-written one beside it. The point is that both run, not that recorded input is better

## 3. The report

- [ ] 3.1 Add the report, classifying every entry through `dialogue.After` and counting per outcome. Covers `The report counts what it classified`
- [ ] 3.2 Add the naming of outcomes no recorded answer reaches, and the test that it does not fail for one. Covers `An outcome no reviewer has reached is named`
- [ ] 3.3 Count the scenarios in the suite that assert each outcome, from the tests rather than from a list somebody maintains
- [ ] 3.4 Write the distribution to a checked-in file, and add the test comparing it. Covers `A removed entry is visible`
- [ ] 3.5 Mutate the routing rule and confirm the comparison fails. Covers `A routing change that moves the corpus is visible`
- [ ] 3.6 Add the report to the ordinary verification sequence, and to `AGENTS.md` beside `make scenario-coverage`

## 4. The task policy

- [ ] 4.1 Record the policy where a task author will meet it, which is `AGENTS.md`, not only in the specification. Covers `A live task is satisfied by recording whatever came back`
- [ ] 4.2 Reword the open task 9.4 of the archived `a-rejection-is-a-claim` to what a live run can satisfy, and say in that change's `retrospective.md` that the rewording is what closed it rather than a dialogue having happened

## 5. What the seed already says

- [ ] 5.1 Report the four seeded entries and record the numbers in this change's own `retrospective.md`: no `error` severity in 28 findings, no `{warning, no-op}` at all, two `{info, ask-user}`, and not one round that would have parked
- [ ] 5.2 Raise, as an open question rather than a change, whether a finding actioned `ask-user` at `info` severity should escalate a run. It happened twice in 28 findings, the rule says it escalates, and nobody has decided that on purpose
- [ ] 5.3 Note which hand-written scenarios rest on shapes no reviewer has produced, `error` severity first among them. Not to remove them: a rule must hold for input a reviewer has not produced yet, and knowing which is which is the whole point of this change

## 6. Verification

- [ ] 6.1 `make fmt && make lint test` clean
- [ ] 6.2 `make scenario-coverage` naming no scenario this change adds
- [ ] 6.3 The ordinary suite still compiles no network-reaching code, and replay works with no build tag
- [ ] 6.4 One live run, recorded into the corpus, satisfied by whatever it answered
