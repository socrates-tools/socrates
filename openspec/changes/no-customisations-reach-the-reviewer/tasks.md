**The invocation is the interface**, so every task asserts on the invocation that was delivered rather than on socrates' own source text, and the simulator's invocation log is where that is read.

**The expected set is derived from the invocation**, never written twice.

**Mutate what you rely on.** Both new guards replace ones a live run proved could be false while reading true.

**No vendor feature name survives** in a requirement, an event, a rendered line or socrates' own source. The gotchas file keeps what it cost, because that is what a gotchas file is for.

**`make lint test` after every task, not at the end of a group.**

## 1. The flag

- [x] 1.1 Carry `--safe-mode` on every invocation and delete the settings allowlist it replaces. Covers `Nothing asks for the machine's own customisations`
- [x] 1.2 Teach the simulator the flag, and assert the delivered invocation carries it. Covers `No tool socrates did not ask for reaches the reviewer`
- [x] 1.3 Measure it against a managed configuration rather than trusting the documentation, and record what the opening event reported with and without it

## 2. The check

- [x] 2.1 Read the session's tool list from the stream's opening event and stop reading anything else from it. Covers `What the reviewer held is read rather than assumed`
- [x] 2.2 Subtract the tools socrates asked for, deriving the expected set from what the invocation carries
- [x] 2.3 Record what is left, under an event kind named for tools rather than for a vendor feature. Covers `A machine that hands the reviewer more than it asked for`
- [x] 2.4 Confirm an ordinary run records nothing, and that a session holding one extra tool records it

## 3. What the reviewer gets

- [x] 3.1 Add the tools that let a reviewer check a claim, and assert every one of them reaches the invocation
- [x] 3.2 Assert the tools that write are absent, naming each of them
- [x] 3.3 Confirm the seal, the comparison and the write refusal are untouched. Covers `A tool that writes cannot reach the change under review`, `What the mechanisms do not cover`

## 4. What goes

- [x] 4.1 Delete the tool-server parsing from the stream, the server list from the response, and the prefix check that read a tool name as evidence of a source
- [x] 4.2 Delete the loaded-servers event and its rendering, and confirm a run recorded before this change still renders legibly
- [x] 4.3 Delete the simulator's allowlist modelling and server emission, and give it a way to hand the reviewer a tool socrates did not ask for
- [x] 4.4 Stop the simulator accepting the two refused flags, and confirm it rejects them

## 5. The gate that blocked the rename

- [x] 5.1 Teach the drop gate a scenario rename declared in the delta
- [x] 5.2 Confirm a scenario renamed without declaring it is still reported

## 6. Fidelity

- [x] 6.1 Pin the opening event's tool list, and stop pinning the field whose presence socrates' own behaviour prevents

## 7. What the documents said

- [x] 7.1 Rewrite the invariant, the README bullets and the gotchas section as the tool surface and the one flag
- [x] 7.2 Correct the README claim that a reviewer can see the installed skill, which the flag makes false
- [x] 7.3 Confirm no invariant, no rendered line, no requirement and no source file names the vendor's feature

## 8. Verification

- [x] 8.1 `make fmt && make lint test` clean
- [x] 8.2 `make scenario-coverage` naming every scenario this change renames
- [x] 8.3 `make scenario-drop` before archiving
- [x] 8.4 `make corpus-outcomes` read, because the routing rule is untouched and its numbers should not move
- [ ] 8.5 A live run confirming from the delivered invocation and the recorded stream that the reviewer held the tools socrates asked for and no others. Require the run, the recording and what the opening event reported, and require no severity, action or count
