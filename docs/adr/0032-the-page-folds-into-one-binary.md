# 32. The page folds into one binary

- Status: accepted
- Date: 2026-09-16
- Part of: one command

## Decision

`socrates` is one binary, and `agora`'s page server is linked into it beside the reviewer.

This supersedes `docs/monorepo/spec.md`'s decision that socrates-review links none of agora's Go code, which was written down specifically to stop itself being undone by accident.
That decision's real content was never that the binary omits an HTTP server.
It was that no review code reaches the page and no review command names one.
Both halves are restated here as an import-graph rule rather than as a fact about what a binary happens to link: no shipped file under `review/` imports any package under `agora/`, and the only thing that imports both is `internal/dispatcher`, which is not a tool.
`internal/convention`'s cross-tool-import check holds this over the whole source tree, the same way the network-reach check already holds that review reaches no network outside a build tag.

The residual is stated rather than hidden.
The binary that runs reviews now contains an HTTP server, where it once provably did not: invoking `socrates agora serve` binds a loopback address on purpose, and running a review starts no server and reaches no network.
`SECURITY.md` says so.

One exit table is a consequence of this decision rather than a decision of its own, and gets no record here.
The two programs arrived independently at the same four slots with the same meanings, because both are read by a person or a script the same way: `0` is nothing asked of you, `1` is something asked of you, `2` is the command failed, `3` is the command was not usable, `4` is the work is still running.
Folding into one binary makes that agreement load-bearing rather than coincidental, and `internal/exit` is its one home.

## Rejected

- Leaving the review binary without a server and shipping `agora` as a second binary, which is what `docs/monorepo/spec.md` chose and what this decision undoes.
  A second binary makes the one-name promise a promise about one of two things a person finds on their `PATH`, and the next workflow after review would face the same choice with no better answer.
- Loosening the import-graph rule to "the review command tree imports nothing of agora's" while leaving some other review package free to import agora, which would let the boundary erode one package at a time with nothing to catch it.
- Recording the exit table in its own decision record.
  It is not a choice this change makes; it is what falls out of one binary sharing one process's worth of exit codes, and both tools had already agreed on it independently.
