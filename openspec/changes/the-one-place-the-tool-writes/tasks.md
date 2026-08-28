**The suite must be green after every task, not only at the end of a group.**
One ordering below exists for that reason: `internal/config` holds the surface block before `internal/surface` stops parsing a file of its own.

**socrates is untouched.**
No task adds a command, a flag, a run-record field or a network reach to `cmd/socrates` or any package it reaches, and the check that measures what socrates reaches is what says so.

## 1. One home holds both

- [ ] 1.1 Give `internal/state` the surface's directory beside `runs/`, and verify a unit test that it sits under the home the one variable states
- [ ] 1.2 Resolve `internal/surface`'s home through `internal/state`, deleting its own variable and its own resolution rule, and verify every path it hands back is under socrates' home
- [ ] 1.3 Refuse `agora` at startup while the variable that once moved its home is set, naming it and naming the one that replaced it, and verify against a live invocation. Covers `The variable that moved the surface's own home is refused`
- [ ] 1.4 Drive a review and a surface under one stated home and verify the runs, the surface's state and the one configuration file are all under it, and that no run directory holds any of the surface's. Covers `One variable moves everything both binaries have` and `The surface keeps a directory of its own inside that home`
- [ ] 1.5 Verify a relative stated home is still refused for both binaries. Covers `A relative home is refused`

## 2. One configuration file

- [ ] 2.1 Add the surface block to `internal/config`, holding the bind and the address the surface is reached at, and verify an unknown field inside it is refused naming the field. Covers `A surface block that names a field socrates does not know is refused`
- [ ] 2.2 Read the block from `internal/config` in `internal/surface`, delete `internal/surface`'s own configuration file, loader and messages, and verify a surface binds and prints what the one file states. Covers `What an operator states about the surface is read from the one file`
- [ ] 2.3 Verify the one file still refuses when it does not parse, and that what it holds is its stated list and the one block

## 3. The skill keeps the session

- [ ] 3.1 Say, where the session is already stated once, that it is kept for the whole review and that a second round a person reads is written to the page already served. Covers `A second round reaches the page the person already has`
- [ ] 3.2 Verify neither branch restates it and that the existing statement is unchanged. Covers `The skill says to carry the session from serve to poll` and `The session is stated once rather than on every branch`
- [ ] 3.3 Regenerate the committed skill with `make skill`, and verify the check that every command and flag the skill names exists still passes
- [ ] 3.4 Drive the loop end to end through the real binaries across two rounds, serving one page, answering it twice and reaching a verdict, and verify the person was handed one address and that the second round reached the first page

## 4. What goes

- [ ] 4.1 Verify nothing names the surface's old home, its old variable or its own configuration file
- [ ] 4.2 Update `README.md` and `AGENTS.md`: one home for the tool, one configuration file, and the boundary that is no longer a boundary
- [ ] 4.3 Verify the surface still names no review, finding, verdict, round or run
- [ ] 4.4 Verify `socrates` gained no command, no flag, no run-record field and no network reach

## 5. Verification

- [ ] 5.1 Mutate the refusal of the old variable and confirm the suite fails, then restore it
- [ ] 5.2 Mutate the rule that the surface's state sits beside the runs rather than inside one and confirm the suite fails, then restore it
- [ ] 5.3 Run `make fmt`, `make lint`, `make test`, `make scenario-coverage`, `make scenario-drop` and `make corpus-outcomes`, and verify every scenario this change states is named by a test
- [ ] 5.4 Record one live run with the real CLI through `make live-smoke`, reaching a page whose state is under socrates' home, and record that it happened
