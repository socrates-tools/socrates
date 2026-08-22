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
socrates asks for `--settings '{"allowedMcpServers":[]}'` and for nothing else.

**`--strict-mcp-config` is refused outright where an enterprise configuration exists, so socrates does not send it.**
`You cannot use --strict-mcp-config when an enterprise MCP config is present`, and the exit status is 1.
The refusal is deliberate rather than a bug: an enterprise configuration is a managed set, and that flag asks to replace it.
It lives outside the person's control, at `/etc/claude-code/managed-mcp.json` on Linux.
socrates used to send it anyway and fall back on the refusal, which is where the second invocation, the retry and a log line reading `the reviewer failed` on every round all came from. The allowlist does the same job and is never refused, so the flag is gone and so is all of that.

**An empty allowlist stops the servers where the flag cannot.**
`allowedMcpServers` may be set in any settings file, entries merge from every source, and the check runs before a server loads, including one from a managed configuration.
An empty array is documented as distinct from unset: unset permits every server, empty permits none.
Supplying it with `--settings` alongside an enterprise configuration took a reviewer's session from 74 tool-server tools to none, measured on this machine.
That leaves `allowManagedMcpServersOnly` as the case socrates cannot answer, because it makes managed allowlists the only ones that count.

**A person's own settings can widen what socrates supplied.**
Allowlists merge from every source, so an `allowedMcpServers` entry in `~/.claude/settings.json` permits that server in a reviewer too.
This is why socrates does not trust either mechanism to have worked and reads the stream's opening event instead, which names `mcp_servers` and the session's `tools`.
What the run records is what that event reported.

**The reviewer's tool list is where a loaded server is visible.**
`--tools` covers built-in tools only, so a tool a server contributes appears in the opening event's `tools` beside them, prefixed `mcp__`.
The live run that found this held `mcp__playwright__browser_run_code_unsafe` and `mcp__chrome-devtools__evaluate_script`, and the seal bounds the snapshot worktree rather than what those reach.

**Doing this right takes the evidence for `mcp_servers[].name` away.**
The fidelity recordings pin the fields socrates reads out of the opening event, and they are recorded through socrates' own invocation, which now permits no server.
So `make fixtures` on a machine like this one records `mcp_servers` as an empty array, and no recording afterwards holds `mcp_servers[].name` for the guard to pin.
The guard requires the name of any server a recording does list, and stays quiet where none does, because failing would read as the vendor removing a field when what happened is socrates keeping servers out.
The two recordings checked in were taken before the allowlist and still list servers, so the name is pinned until somebody re-records them.

**A sealed worktree is fully usable by the vendor's CLI.**
Clearing every write bit in the worktree does not disturb a reviewer: `git status`, `git grep` and every read work, because a worktree's index and administrative state live under the main repository's `.git/worktrees/<name>/` rather than in the worktree.
A tool server that writes at startup has its write refused and the session carries on, verified against the real CLI with an enterprise configuration loaded.

**`chmod` returning no error is not evidence that permissions are enforced.**
A filesystem can accept the call and ignore the bits, and `root` ignores them by design.
socrates attempts one write in the sealed worktree before spawning a reviewer, and refuses the review where that write succeeds.

**`os.Chmod` follows symlinks.**
Sealing a worktree by walking it would reach outside the worktree through any symlink in the change under review, which is a write in somebody else's tree.
`filepath.WalkDir` does not follow them, and the walk skips anything that is not a directory or a regular file.

**Clearing the write bit means `mode &^ 0o222`, not `0o444`.**
Assigning a flat mode makes a `0o600` file world-readable and takes the executable bit off a script.
Clearing the write bits keeps everything else: `0o755` seals to `0o555`, and lifting restores the owner's write bit, which is what deleting and inspecting the worktree need.
Lifting is not the inverse of sealing: a `0o664` file lifts to `0o644` and a `0o444` file lifts to `0o644`.
Nothing depends on the modes a lifted worktree carries, because a worktree is thrown away or handed to a person to read.

**A sealed directory cannot be removed.**
`os.RemoveAll` unlinks entries, which needs write permission on the directory that holds them, so every path that removes a worktree or a run directory lifts the seal first.
One abandoned run with a sealed worktree would otherwise make the retention sweep fail, and the sweep runs at the start of every review.

## Testing

**One journey needs a filesystem that ignores permission bits, and skips without one.**
`TestASealThatDoesNotHoldIsARefusal` in `internal/e2e` is the only test of what `scrat review` does where the seal cannot be enforced, and no test can manufacture such a filesystem.
Name one and it runs:

```sh
SOCRATES_TEST_UNENFORCING_DIRECTORY=~/winhost/temp go test ./internal/e2e/ -run TestASealThatDoesNotHoldIsARefusal
```

A Windows drive mounted into WSL2 is one: `chmod` there returns no error and changes nothing.
The test refuses to run where the directory named does enforce permissions, rather than passing while asserting nothing.

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
