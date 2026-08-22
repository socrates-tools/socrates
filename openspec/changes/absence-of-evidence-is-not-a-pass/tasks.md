**Both guards replace ones that read as passing while checking nothing**, so each of them is mutated on purpose before it is trusted: the test that proves the new behaviour is run against the old code first.

**The agent half asserts on a delivered stream**, read back from the round directory or the simulator's invocation, never on socrates' own source text.

**The gate half asserts over a fabricated tree**, because what the rule does must not depend on whatever change happens to be proposed today.

**`make lint test` after every task, not at the end of a group.**

## 1. A session that never said what it held

- [x] 1.1 Teach the simulator to emit an opening event with no tool list, and assert from the delivered stream that the field is absent rather than empty
- [x] 1.2 Record it under its own event kind, on a round that otherwise completed. Covers `A session that never said what it held`
- [x] 1.3 Confirm the round is not failed for it, and that nothing records the session as holding only what socrates asked for
- [x] 1.4 Confirm a round that failed before answering records the failure and not the absence
- [x] 1.5 Render it as one sentence in the one place that renders an event

## 2. A rename the gate cannot check

- [x] 2.1 Report a declared scenario rename whose new name the promoted specification already holds and whose old name it does not
- [x] 2.2 Confirm a legitimate rename is still followed, and that an identical FROM/TO pair is left alone
- [x] 2.3 Fail the gate on it, and give the printed scenario line the wording its condition needs

## 3. What the documents say

- [x] 3.1 State both guards as invariants, and correct nothing else

## 4. Verification

- [x] 4.1 `make fmt && make lint test` clean
- [x] 4.2 `make scenario-coverage` naming the scenario this change adds
- [x] 4.3 `make scenario-drop` before archiving, on a tree where it has a before-state to read
- [x] 4.4 `make corpus-outcomes` read, because the routing rule is untouched
- [ ] 4.5 A live run confirming from the recorded stream what the session reported and what the run recorded about it. Require the run, the recording and the routing, and require no severity, action or count
