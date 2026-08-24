# Where the commit comes from

## What was measured

Go stamps `vcs.revision`, `vcs.time` and `vcs.modified` into a binary built from a git repository, and `runtime/debug.ReadBuildInfo` reads them back with no build flags and no Makefile.
That is the mechanism this change would have used, and it does not work here.

In a linked git worktree, whose `.git` is a file pointing into the main checkout's `.git/worktrees`, `go build` emits no VCS stamps at all.
`go build -buildvcs=true`, which is documented to fail rather than silently omit, exits zero and still emits none.
The same package in an ordinary clone of the same repository emits `vcs=git`, `vcs.revision`, `vcs.time` and `vcs.modified`.

This repository is developed in linked worktrees, so the branch that would carry the answer is the branch that never runs.
A mechanism that works everywhere except where the work happens is not a mechanism.

## What is used instead

The Makefile computes the short commit and whether the tree held uncommitted changes, and passes both to the linker with `-X`.
`build` and `install` take them from one variable, so a binary from either says the same thing about itself.

Two fields rather than one string, because the dirty marker is a fact about the build and not part of the commit's name, and a rule that had to split a suffix back off a stamped string could disagree with the Makefile that joined it.

Uncommitted means whatever `git status --porcelain` reports, which counts untracked files.
That is the same reading socrates already applies to the tree it snapshots, where the reason is written down: a freshly implemented feature is mostly untracked files, so a definition of clean that ignores them is a definition that calls the interesting case clean.

## What was rejected

**Reading the commit at run time.**
Shelling out to `git rev-parse` when `scrat version` runs would report the repository the command is standing in, which for a tool run against other people's repositories is almost never the one it was built from.
It would have answered the stale-binary question with the same confidence and the wrong commit.

**A version constant edited by hand.**
It goes stale silently, and silent staleness is the failure this change exists to end.

**`--version` on the root command.**
Every other question a person asks `scrat` is a subcommand, and a flag that answers rather than modifies would be the only one of its kind on the surface.
One spelling, because two spellings of one answer is the duplication this repository keeps failing tests to prevent.

## How the branches are reached in a test

Neither the dirty tree nor the unstamped build can be produced by a test that refuses to write in the checkout it is running from, and this one refuses.

The stamp is therefore a pair of Makefile variables, so a test overrides them on the command line and the real recipe does the rest.
`make build BINARY=... UNCOMMITTED=yes` reaches the second branch, `make build BINARY=... COMMIT=` reaches the third, and the default invocation reaches the first with whatever the checkout actually is, compared against the same two git commands the Makefile used.
What is driven is the recipe that ships, rather than a reconstruction of it in Go.
