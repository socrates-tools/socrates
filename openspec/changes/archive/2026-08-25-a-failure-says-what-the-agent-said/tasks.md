**The guard being replaced reads as informative while carrying nothing**, so the test that proves the new reason is run against the old code first and seen to fail on the empty quotation it produces.

**Every assertion is on a delivered stream**, driven through the real binary against the simulator, never on socrates' own source text.

**`make lint test` after every task, not at the end of a group.**

## 1. A stream that reports its own failure

- [x] 1.1 Teach the simulator to report an error in its result event, carrying an account, keeping the subtype the live stream carried, writing nothing to the error stream, and exiting with what the scenario asks for
- [x] 1.2 Assert from the delivered stream that the error flag, the subtype and the account are where the live run recorded them

## 2. The reason a failed call carries

- [x] 2.1 Read the account the result event carries, and prefer it to the agent's error stream. Covers `A non-zero exit after the agent said why`
- [x] 2.2 Carry it where the agent exits zero, in place of the subtype beside it. Covers `A result that reports an error carries what it said`
- [x] 2.3 Keep the subtype as what is said where the event reports an error and carries no account. Covers `A result that reports an error and says nothing`
- [x] 2.4 Carry both accounts where both are there, and confirm the existing non-zero exit still carries what the agent wrote to its error stream
- [x] 2.5 Say that the agent said nothing where it said nothing, and never trail a colon onto an empty string. Covers `A failure the agent said nothing about`
- [x] 2.6 Bound the account on the ceiling and the mechanism the error stream already has. Covers `An account too long to carry is shortened`
- [x] 2.7 Confirm nothing branches on the account: no state, verdict or exit code moves with a word in it

## 3. The field this reading depends on

- [x] 3.1 Pin the account and the error flag on the result event against the checked-in recordings, and mutate the guard to prove it fails without them

## 4. What the documents say

- [x] 4.1 State the invariant, and correct nothing else

## 5. Verification

- [x] 5.1 `make fmt && make lint test` clean
- [x] 5.2 `make scenario-coverage` naming the scenarios this change adds
- [x] 5.3 `make scenario-drop` read, on a tree where it has a before-state
- [x] 5.4 `make corpus-outcomes` read, because the routing rule is untouched
- [x] 5.5 A live run whose recorded stream is read for the fields this reading depends on, and whose answer socrates routed. Require the run, the recording and the routing; require no severity, action, count, and no failure, because whether a live run fails is not this repository's to arrange

## What the live run recorded

Run `20260825T005107Z-7034b716`, against `charon-test` on a working tree nobody had reviewed, with `--intent "add an endpoint that returns the most recent items"`.

The reviewer answered, the round parked, and the run recorded `round-parked` after `agent-completed`.
The result event this reading depends on carried `is_error: false` and `subtype: "success"`, so `reportsAnError` read false and the account this change reads was correctly left unread on a call that did not fail.
`scrat findings` reported one finding needing a patch and one needing a person, and `scrat status` reported the run `awaiting-response` with the same two counts.
What the reviewer found is the reviewer's and is not what this task required.

The delivered invocation carried the four tools socrates asks for and no external source.
