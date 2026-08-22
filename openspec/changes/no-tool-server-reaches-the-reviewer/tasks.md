**The invocation is the interface, so every task asserts on the invocation that was delivered** rather than on socrates' own source text, and the simulator's invocation log is where that is read.

**The `init` event is the arbiter.** No task may reintroduce a decision that reads the text of a vendor message.

**Mutate what you rely on.** Two guards here replace guesses, and a guard nobody has broken on purpose is a guard nobody knows works.

**`make lint test` after every task, not at the end of a group.**

## 1. The allowlist

- [x] 1.1 Add the empty tool-server allowlist to every invocation's settings. Covers `Nothing asks to replace a managed set`
- [x] 1.2 Teach the simulator the settings argument and have it honour one, and assert the delivered invocation carries it. Covers `No tool server is loaded`
- [x] 1.3 Confirm socrates supplies the allowlist and nothing else on every invocation, and record honestly that a user's own allowlist can widen it, which is what reading the opening event catches

## 2. Reading what loaded

- [x] 2.1 Parse the stream's opening event for the tool servers and the session's tool list, and carry both out of the parser
- [x] 2.2 Have the simulator report both in its opening event, so the recording path is testable at all
- [x] 2.3 Record what that event reported rather than what an error message said. Covers `What loaded is read rather than inferred`
- [x] 2.4 Mutate the recording: report a server in the opening event and confirm the run records it, then report none and confirm it records none

## 3. The flag, and everything it needed

- [x] 3.1 Remove `--strict-mcp-config` and its empty `--mcp-config` from the invocation, because the allowlist does the same job and no machine refuses it. Covers `Nothing asks to replace a managed set`
- [x] 3.2 Delete the retry, the second invocation and the event that reported a refusal, none of which has anything left to do
- [x] 3.3 Delete the simulator's refusal behaviour, because nothing sends the flag that triggered it, and a scenario nothing can reach rots
- [x] 3.4 Confirm one invocation per round: a rejected startup fails with no fallback, and a binary that is not there is not retried

## 4. What a person reads

- [x] 4.1 Confirm the line reading `the reviewer failed` is gone from an ordinary run, because there is no rejected attempt left to report
- [x] 4.2 Say in the loaded-servers line what is actually true, rather than offering the seal as the reassurance this change retracts
- [x] 4.3 Confirm a review that really did fail still records and renders as one

## 5. What the mechanisms cover

- [x] 5.1 State in the specification that all three bound the worktree, and assert the seal and the comparison are unchanged by this change. Covers `What the mechanisms do not cover`, `A tool server that writes cannot reach the change under review`, `A reviewer that wrote fails the round`

## 6. Fidelity

- [x] 6.1 Confirm the recorded shapes already carry the opening event's tool servers and tool list, and guard them in the ordinary suite, because the drift check is one-directional and would not notice the vendor removing one
- [x] 6.2 Confirm the simulator satisfies the recording, and mutate a recording to prove the guard fails

## 6b. What the documents said

- [x] 6b.1 Rewrite the tool-server sections of `docs/gotchas.md`, `AGENTS.md` and `README.md`, which all say the refusal is answered by the seal
- [x] 6b.2 Say in each of them what the seal bounds and what it does not

## 7. Verification

- [x] 7.1 `make fmt && make lint test` clean
- [x] 7.2 `make scenario-coverage` naming no scenario this change adds
- [x] 7.3 `make scenario-drop` before archiving
- [x] 7.4 `make corpus-outcomes` read, because the routing rule is untouched and its numbers should not move
- [x] 7.5 A live run on the machine that has a managed tool server configuration, confirming from the delivered invocation and the recorded stream that no tool server loaded. Require the run, the recording and what the opening event reported, and require no severity, action or count
