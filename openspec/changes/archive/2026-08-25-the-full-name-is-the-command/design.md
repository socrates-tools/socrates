# Two names, one binary

## How the second name is installed

`go install ./cmd/socrates` puts one binary wherever the toolchain installs to, and `make install` then links `scrat` to it in that same directory.

The directory is asked of the toolchain rather than guessed: `go env GOBIN` when it is set, and `$(go env GOPATH)/bin` when it is not.
A Makefile that hardcoded `~/go/bin` would be wrong on the first machine that sets `GOBIN`, and wrong silently, because the install would succeed and the short name would land somewhere nothing reads.

The link is relative, so the pair survives the directory being moved, and `-f` makes reinstalling idempotent rather than a failure the second time.

## Why not two binaries

`go install ./...` would produce `socrates` and `scrat` from two mains, and both would be real programs with no symlink involved.
It was rejected for what it does to the question the previous change just made answerable: two binaries on one machine can be at different commits, so `scrat version` and `socrates version` could disagree, and the whole point of that line is that it can be trusted.

It also costs the command layer its shape.
Two mains cannot share a `main` package, so every command file would move to an internal package to be called from both, which is a large rearrangement bought entirely by avoiding one `ln -s`.

## Why the text names one of them

The binary knows which name it was invoked as, and could report that name back.
It does not, because then `socrates` and `scrat` would print different text, the skill and the specifications would still have to pick one, and a person reading a refusal could not tell whether the sentence is about the thing they typed or the thing the documents describe.

The full name is what the text says, under either invocation.
The short name exists to be typed, not to be quoted back.

## What is not renamed

Specifications already promoted, and the deltas of changes still open, keep the `scrat <command>` they were written with.
`scrat` is still a command name, so those sentences are true, and the alternative is a rename declaration for every scenario heading that names one while two changes are open against the same capabilities.
The gates compare headings by normalised text, so a sweep of that prose would be the kind of edit that reports a clean promotion because nothing is left to compare.

What is renamed is the source that produces text a person reads, and the documents this repository owns about itself.
