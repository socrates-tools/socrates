**The rule half is decided over its whole input space**, which is small enough that sampling it would be a choice rather than a necessity.

**The delivery half drives the Makefile's own recipe**, never a reconstruction of its flags in Go, because the recipe is what a person runs.

**Nothing here asserts a particular commit.** What is asserted is that the line agrees with the checkout the build was made from, read by the same two git commands the recipe reads.

**Each guard is mutated on purpose before it is trusted:** the test that proves the new behaviour is run against the old code first.

**`make lint test` after every task, not at the end of a group.**

## 1. The line a person reads

- [x] 1.1 Decide the line over every combination of a recorded commit and a recorded dirty state, including the absent commit. Covers `The build says which commit it came from` and `A build from a tree with uncommitted changes says so`
- [x] 1.2 The absent commit says so and names what to build it with, and carries no commit of its own. Covers `A build that recorded no commit says so`
- [x] 1.3 Confirm the package reaches no filesystem, no process and no clock, so nothing about the line can be read from the repository the command is run in

## 2. The command

- [x] 2.1 Print the line and nothing else, and exit zero
- [x] 2.2 Refuse an argument, naming it, with the usage code. Covers `scrat version takes no arguments`
- [x] 2.3 Confirm the surface lists it and still lists nothing that lands or verifies. Covers `The command surface is what it says`

## 3. The stamp

- [x] 3.1 Compute the commit and the dirty state once in the Makefile and pass them from `build` and from `install`
- [x] 3.2 Drive `make build` and confirm the line names the checkout's own short commit and agrees with what `git status --porcelain` says about it
- [x] 3.3 Drive the recipe with the dirty state overridden, and with the commit overridden to nothing, and confirm each branch is what the binary prints
- [x] 3.4 Mutate the stamp out of the recipe and confirm the driven tests fail rather than the rule's alone

## 4. What the documents say

- [x] 4.1 Record `internal/version` in the layout
- [x] 4.2 State the invariant that the commit is recorded at link time and never read from the repository the command is run in
- [x] 4.3 Correct nothing else

## 5. Verification

- [x] 5.1 `make fmt && make lint test` clean
- [x] 5.2 `make scenario-coverage` naming every scenario this change adds
- [x] 5.3 `make scenario-drop`, which reads one substituted requirement carrying everything it holds
- [x] 5.4 `make corpus-outcomes` read, because nothing about the routing rule is touched
- [x] 5.5 `make install` and `scrat version` against `git rev-parse --short=7 HEAD`, which is the check the whole change exists to make possible

## What the verification recorded

`make install` linked `-X internal/version.commit=fbc7bf8 -X internal/version.uncommitted=yes` from a checkout at `fbc7bf8` holding ten uncommitted paths, and the installed binary said `scrat fbc7bf8, built from a working tree with uncommitted changes`.

Removing the stamp from the recipe was tried on purpose: the driven tests failed naming the commit they expected, and the rule's own tests stayed green, which is the split the guard is there to prove.
