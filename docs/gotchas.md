# Gotchas

Things this repository had to re-derive, each of them silent when it went wrong.
Add to this file when something costs you an hour.

## git

**`git stash create` ignores untracked files.**
It was the obvious way to freeze a working tree and it is the wrong one.
A freshly implemented feature is mostly untracked files, so run against a real repository it produced nothing at all, and the review that followed was a confident review of an empty diff.
The recipe that works is a temporary index plus `add -A`, which honours `.gitignore` and includes untracked files.

**`GIT_INDEX_FILE` must name a path that does not exist yet.**
`read-tree` creates the temporary index from `HEAD`, and `add -A` stages into that.
Point it at an existing file and you are staging into somebody's index.
Point it at a directory that does not exist and git fails.
socrates makes a temporary directory and names a file inside it that git then creates.

**A snapshot commit is only reproducible if you pin the dates.**
`commit-tree` folds the author and committer timestamps into the commit, so the same tree produced a new SHA every second.
Pinning `GIT_AUTHOR_DATE` and `GIT_COMMITTER_DATE` to `@0 +0000` makes the commit content-addressed: the same worktree and base always yield the same commit, and therefore a byte-identical prompt.
That is what makes "the prompt is stable across two reviews" testable at all, and it stops a fresh object accumulating per review.

**`git worktree` leaves administrative state behind.**
Removing the directory is not removing the worktree.
`git worktree prune` after the removal is what finishes it, and every failure path needs it too.

## Go

**`testing/synctest` cannot test a deadline that a real subprocess is on the other side of.**
The bubble fakes time for goroutines inside it, not for a process the kernel is scheduling.
Use it for the deadline logic around `exec` and drive the real-subprocess deadline from an end-to-end test with a short configured deadline instead.
A `synctest` bubble also fails if a goroutine inside it blocks on something outside the bubble, which a pipe from a real process is.

**A hermetic `HOME` breaks `go build` inside a test.**
Tests that set `HOME` to a temporary directory lose the module cache and the build cache with it, and the build then tries to reach the network.
`internal/testing/source` captures `GOPATH`, `GOMODCACHE` and `GOCACHE` once at package initialisation and passes them into every build it runs.

**`go.yaml.in/yaml/v3` decodes `time.Duration` from a string.**
`agent_deadline: 45m` lands in a `time.Duration` field with no custom unmarshaller.
`KnownFields(true)` is what makes a misspelled setting a refusal rather than a silent default.

**An empty YAML document decodes as `io.EOF`.**
An empty configuration file is the defaults, not an error, so that case is handled explicitly.

**cobra has no typed usage error.**
An unknown flag reaches you through `SetFlagErrorFunc`, and an unknown command comes back from `Execute` as a plain error indistinguishable by type.
socrates gives the root command `ArbitraryArgs` and a `RunE` that refuses the first argument itself, so a usage error is a refusal and exits `3` rather than looking like a failed review.

**`go vet` is not the whitespace gate.**
`gofmt -l .` is, and `make lint` runs it first, so formatting is never something to remember separately.

## The vendor CLI

**Restricting the reviewer's tools does not restrict its tool servers.**
`--tools Read,Grep,Glob` says nothing about the servers the machine loads, and a server writes wherever it likes the moment it starts.
The first live run of socrates failed because the machine's own configured server wrote a `.serena` directory into the snapshot worktree, and the round correctly refused a write the reviewer never asked to make.
socrates now asks for `--mcp-config {"mcpServers":{}} --strict-mcp-config`.

**`--strict-mcp-config` is refused outright where an enterprise configuration exists.**
`You cannot use --strict-mcp-config when an enterprise MCP config is present`, and the exit status is 1.
An enterprise configuration lives outside the person's control, at `/etc/claude-code/managed-mcp.json` on Linux, so there is nothing socrates can do about it from inside the invocation.
socrates asks for strictness, and on that refusal runs the review again without it and records `agent-tool-servers-loaded` in the run's event log.
Where a server named there writes into the working directory, every review on that machine fails the worktree comparison, and the remedy is to take the server out of that configuration rather than to weaken the comparison.

## Testing

**Never revert a mutation with `git checkout <file>` while the file is new.**
A file that is tracked at an older content reverts to that content, not to what you had a minute ago.
This destroyed a rewrite of `internal/prompt/prompt.go` once.
Copy the file into the scratchpad before mutating and copy it back afterwards.

**A green suite against a simulator manufactures confidence.**
The drift check earned its place on its first run: it found that the simulator emitted `structured_output` even when no schema was supplied, which the real CLI does not do.
No test that used the simulator could have found that, because every one of them supplied a schema.

**Record shapes, not streams.**
A raw recording of the real CLI carries the recording machine's hook output, memory paths, socket paths and session identifiers.
The drift check only needs the field paths per event type, and shapes are what a reviewer of the diff can actually read.

**Only invention is drift.**
The real CLI emits event types and fields the simulator does not, and that is fine.
The check is one-directional: every field the simulator produces must exist in the recording.
The other direction would fail permanently and be turned off within a week.

**`testscript` scripts assert exit codes through `sh`.**
`exec` fails a script on any non-zero status and there is no assertion for a specific code, so `exec sh -c 'scrat status; echo exit=$?'` plus `stdout 'exit=3'` is how the contract gets checked.
`cp stdout file` is how you keep output for a later assertion, and `$WORK` is expanded inside single quotes.
