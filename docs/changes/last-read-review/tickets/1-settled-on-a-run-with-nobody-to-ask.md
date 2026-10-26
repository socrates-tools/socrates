# 1: A run with nobody to ask takes a settled answer

## What to build

A session can begin a review with nobody present to answer a finding, and answer a finding that asks a person by settling it itself, with a reason, recorded as the session's decision and never as a person's.
A run begun with `socrates review begin --nobody-to-ask` records that once, and every later command reads it.
On such a run a finding that asks a person accepts `socrates review respond --settled --reason` and refuses `--decided`; on every other run that is mine it accepts `--decided` and refuses `--settled`.
`--nobody-to-ask` beside `--not-my-change` is refused before anything is snapshotted.
A round left waiting only on findings that ask a person does not count as waiting on a person on such a run, so `begin --wait` returns rather than reporting that it parked on somebody.
`status` and `findings --json` say the run has nobody to ask and which responses were settled.
The response rules the socrates-review skill carries say `settled` belongs only to a run begun with `--nobody-to-ask`, which that skill never begins.
review's glossary already defines settled and nobody to ask, and review's invariants gain what ADR 91 states.

## Acceptance criteria

- [ ] `socrates review begin --nobody-to-ask` is accepted, and `socrates review status` on that run says nobody is there to ask; today the flag is refused as unknown.
- [ ] On that run, `socrates review respond --settled --reason "..."` answering a finding that asks a person is recorded, and `findings --json` reports its resolution as `settled`; today `--settled` is refused as unknown.
- [ ] On that run, `respond --decided` answering the same finding is refused, naming `--settled` as what it takes.
- [ ] On a run begun without the flag, `respond --settled` is refused, naming `--decided` as what it takes, and `--decided` is still recorded as before.
- [ ] `respond --settled` with no `--reason` is refused, saying a settlement needs a reason.
- [ ] `begin --nobody-to-ask --not-my-change` is refused before a snapshot is taken or a reviewer spawned, and no run is created.
- [ ] On a run with nobody to ask whose first round holds only findings that ask a person, `begin --wait` exits with the awaiting-response status without saying the run waits on a person, and `continue` after every finding is settled starts the second round.
- [ ] The socrates-review skill, as installed, says `settled` belongs to a run begun with `--nobody-to-ask` and never tells its session to begin one.

## User stories

- 8, 9, 15

## Left to other tickets

- 2.1-second-round-reads-settled-as-the-sessions-decision.md: how the second round's reviewer is told about a settled response.
- 2.2-last-read-reviews-the-whole-change.md: the overseer using any of this.

## Blocked by

None
