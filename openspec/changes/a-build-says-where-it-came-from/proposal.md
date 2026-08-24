## Why

A `scrat` binary cannot say which commit it was built from, and the failure that follows is already recorded in this repository.
The live run that closed `a-round-a-person-reads-becomes-a-page` read a build two days old: `make install` writes to `$GOPATH/bin`, another directory earlier on `PATH` held an older `scrat`, and the run reported nothing about the new behaviour because the binary predated it.
Nothing the binary printed said which build it was, so the only way to notice was to guess.

That is the whole need: a person or a session holding a `scrat` should be able to ask what it was built from and check the answer against `git log`.

## What Changes

- `scrat version` prints one line naming the short commit the binary was built from, and says so when the working tree it was built from held uncommitted changes.
- The commit is recorded in the binary when it is linked, by the Makefile, and it is never read from the repository the command happens to be run in.
  A binary that reported the tree it is standing in would answer the one question this command exists for with the one answer that cannot be checked.
- Go's own build stamping is not the mechanism, because it records nothing in a linked git worktree, which is where this repository is developed.
  This is measured rather than assumed: `go build -buildvcs=true` in this worktree emits no `vcs.revision`, no `vcs.modified` and no `vcs` at all, while the same build in an ordinary clone emits all three.
- A binary carrying no recorded commit says so and says what to build it with.
  It never prints an empty version, and it never infers one.
- The command surface gains one command and nothing else: no flag, no JSON, and nothing in the skill, because a version line is not part of the loop a session drives.

## Capabilities

### Modified Capabilities
- `cli-foundation`: the surface gains `version`, which reads what the binary itself records and reaches nothing.

## Impact

- `internal/version`: the two fields the linker sets, and one pure rule from them to the line a person reads.
- `cmd/scrat`: one command, which prints that line and takes no arguments.
- `Makefile`: one stamp, computed once and passed by both `build` and `install`, so the two cannot disagree about what they produced.
- `internal/e2e`: the Makefile's own recipe is what the tests drive, with the stamp overridden to reach the branches a checkout cannot be put into on purpose.
- Nothing in `internal/skill`, nothing in any prompt, and nothing written in the repository under review.
