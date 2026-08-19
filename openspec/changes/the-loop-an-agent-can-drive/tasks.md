**Nothing here adds a command or a flag for the skill's benefit.** The skill drives the surface that already exists. A task that finds the skill needing something the binary does not offer has found a gap in an earlier layer, and the fix belongs there.

**This change is all text, and text fails quietly.** Every assertion below is about generation, drift, the shared source, or a journey driven end to end. None of them is about whether the wording persuades a model, which is evaluated by using it.

**`make lint test` after every task, not at the end of a group.**

## 1. One source

- [ ] 1.1 Move the four response rules into one constant, and make `scrat findings` read it rather than holding its own copy. Covers `One source, two surfaces`
- [ ] 1.2 Add the skill's text as a constant, and the generator that writes the committed copy
- [ ] 1.3 Add the drift test and the regeneration target. Covers `Drift fails the suite`, `Regenerating is one command`

## 2. What the skill says

- [ ] 2.1 Write the loop, with every terminal state and `failed` before `escalated`. Covers `Every terminal state is covered`, `A failed run is stated not to be a pass`, `The loop closes`
- [ ] 2.2 Write the waiting step around `scrat status`, with its reason, and mention `scrat watch` once as a suggestion for the person. Covers `The waiting command is status`, `Watch is for the person`
- [ ] 2.3 Put the response rules ahead of the step that reads the findings. Covers `The rules are stated before the findings are read`
- [ ] 2.4 Write the three prohibitions. Covers `The three prohibitions are present`, `No channel is implied for the author's account`

## 3. The skill cannot name something that is not there

- [ ] 3.1 Add the check that extracts every command and flag the skill names and asserts each exists on the built binary. Covers `Every command named is real`
- [ ] 3.2 Add the test that removing a named flag fails the check. Covers `A removed flag fails the suite`. Without it the extraction can silently match nothing and pass forever

## 4. Installing

- [ ] 4.1 Add `scrat skill`, its user-level default, and the printed path. Covers `Installing names the path`
- [ ] 4.2 Add idempotence and `--to`. Covers `Installing twice`, `Installing elsewhere`
- [ ] 4.3 Add the refusal on an unwritable target. Covers `A directory that cannot be written`
- [ ] 4.4 Update the command surface and its test. Covers `The command surface is what it says`, `Installing the skill writes nothing in the repository`

## 5. End to end

- [ ] 5.1 Drive the loop the skill describes, step by step, against the simulator, and assert it reaches `clean`. The skill is only correct if the sequence it names actually works
- [ ] 5.2 Drive the same sequence through the escalated branch, the failed branch and the parked branch, and assert each stops where the skill says to stop
- [ ] 5.3 Assert across every journey that nothing was written in the repository

## 6. Documentation

- [ ] 6.1 Update `README.md`: install the skill, and what the session does with it
- [ ] 6.2 Update `AGENTS.md` with where the skill's source is and what regenerates it

## 7. Verification

- [ ] 7.1 `make fmt && make lint test` clean
- [ ] 7.2 `make scenario-coverage` naming no scenario this change adds
- [ ] 7.3 `make scenario-drop` before archiving, for the command surface requirement this change substitutes
- [ ] 7.4 A live run: install the skill, hand a real session a real task, and let it drive a real review to a real verdict without being told the sequence. Record in `retrospective.md` every place it deviated, because each one is a sentence the skill is missing rather than a fault of the session
