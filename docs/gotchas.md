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
socrates-review makes a temporary directory and names a file inside it that git then creates.

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
socrates-review gives the root command `ArbitraryArgs` and a `RunE` that refuses the first argument itself, so a usage error is a refusal and exits `3` rather than looking like a failed review.

**`go vet` is not the whitespace gate.**
`gofmt -l .` is, and `make lint` runs it first, so formatting is never something to remember separately.

## The vendor CLI

**Restricting the reviewer's tools does not restrict where else tools come from.**
`--tools` is an allowlist over the vendor's built-in set and nothing else.
It says nothing about the customisations the machine loads, and one of those writes wherever it likes the moment it starts: socrates-review's first live run failed because a configured tool source wrote a `.serena` directory into the snapshot worktree, and the round correctly refused a write the reviewer never asked to make.
A later run found the same session holding script execution, file upload and browser control, none of which `--tools` had any opinion about.

**`--safe-mode` is the one flag that answers all of it.**
It disables the customisations the machine brings on its own account, and the vendor's own help lists them: memory files and auto memory, skills, plugins, hooks, external tool sources, custom commands and agents, output styles, workflows, themes, keybindings, the status line, the file-suggestion commands and LSP servers.
Auth, model selection, built-in tools and permissions keep working, which is the whole point: the reviewer loses the machine's additions and keeps the CLI.
The same help says `Admin-managed (policy) settings still apply`, and the binary says which part of them that is: `Hooks from settings files are suspended and will not run this session (managed policy hooks still apply)`, beside a `Hooks Restricted by Policy` notice.
So the residual is a hook an administrator's policy configures, along with a policy-configured status line and policy-configured file-suggestion commands: only the policy-configured ones survive, and the machine's own do not load.
A policy-configured tool source is not in that residual, which the measurement below shows going away, and the vendor's documentation puts managed plugins, managed skills and managed memory with it.
That half is documentation rather than measurement, because the tool list does not say what contributed a tool and its `plugins` field is an inventory of the machine.
A policy hook contributes no tool for the list check to see, but its command is a child of the session, whose working directory is the sealed worktree, so a write it makes there is refused and the comparison afterwards catches one the seal did not.
What nothing socrates-review holds bounds is where that command reaches outside the worktree, so the residual is socrates-review's to state rather than to close: the flag that would ask to replace a managed set is the one refused at startup.
Measured on this machine, which deploys a managed configuration: without it the opening event named four external sources, with it none, and the plugin-contributed `LSP` tool disappeared from the session's tool list at the same time.
It is not refused by a managed configuration, which is the trap the flags below fall into.

**`--safe-mode` arrived in 2.1.169, and socrates-review keeps no fallback.**
An older CLI dies at startup on an unknown option, every round of every run.
Recognising that text to say something better is the pattern this repository deleted for being fragile, so the floor is stated in the README prerequisites instead and the CLI's own message is what a person reads.

**The flags that ask to replace a managed set are refused at startup.**
`You cannot use --strict-mcp-config when an enterprise MCP config is present`, exit status 1, and the refusal is deliberate rather than a bug: a managed configuration is a managed set and that flag asks to replace it.
The file lives outside the person's control, at `/etc/claude-code/managed-mcp.json` on Linux.
socrates-review sent it anyway once and fell back on the refusal, which is where a second invocation, a retry and a log line reading `the reviewer failed` on every round of every run all came from.
Nothing sends it now, and the simulator no longer accepts it, so putting it back fails every test that runs a reviewer rather than quietly reintroducing the fallback.

**A settings allowlist also works, and was what socrates-review used before `--safe-mode`.**
An empty `allowedMcpServers` supplied through `--settings` is evaluated before a source loads, including one from a managed configuration, and an empty array is documented as distinct from unset: unset permits everything, empty permits nothing.
It took a reviewer's session from 74 externally contributed tools to none.
`--safe-mode` replaced it because it covers hooks and plugins too, and because it needs no vendor feature named in socrates-review's own vocabulary.

**The session's tool list is read rather than trusted.**
socrates-review takes `tools` from the stream's opening event, subtracts what it asked for, and records the remainder.
Anything a customisation contributes lands in that list, so one check covers every source rather than one of them.
An externally contributed tool is recognisable in that list by an `mcp__` prefix, which is how the live run that found this was diagnosed.

**Two of the opening event's fields are about the session and the rest are an inventory of the machine.**
Under `--safe-mode` its `skills`, `plugins`, `slash_commands` and `agents` fields name everything installed on the machine, and a check that read them would report a customisation on every run of every machine that has one.
`tools` and `mcp_servers` are the two that answer for the session: the tool list holds only what socrates-review asked for, and the source list is empty.
The measurement is what separates them, not the field names: with the flag the sources named went from four to none while the plugin inventory beside them did not move.

**The opening event is not always the first line of the stream.**
Without `--safe-mode` a `SessionStart` hook fires first on this machine and its `hook_started` event arrives before `system`/`init`; with the flag no hook of this machine's own runs and the init event is first, though a machine whose policy configures one would still see it fire.
socrates-review matches the event by type and subtype wherever it appears, so this is one more thing the flag makes true rather than something to rely on.

**A tool list at the opening event can be incomplete, so the source list is read beside it.**
External sources connect asynchronously, so a probe on this machine reported four of them in the opening event while the tool list held none of their tools yet.
Reading only the list would leave a machine that loaded a source anyway recording nothing at all, which is indistinguishable from a clean session.
Both fields are read, and they are blind in different places: the list misses a source whose tools are still arriving, the source list misses a built-in the vendor added.
Neither reading fails a round, and what covers a source that connects after socrates-review has read the event is `--safe-mode` not loading it and the comparison after the round.

**Pinning `mcp_servers` against the recordings survives a re-record, and pinning `mcp_servers[].name` would not.**
Under socrates-review's own invocation the field comes back as an empty array, so no `[].name` is derived from it and a guard requiring one fails reading as though the vendor had removed a field.
`fidelity.ShapesOf` records the key before it walks the value, so an empty array still yields the field itself.
That is the same trap the earlier guard fell into, one field along.

**A sealed worktree is fully usable by the vendor's CLI.**
Clearing every write bit in the worktree does not disturb a reviewer: `git status`, `git grep` and every read work, because a worktree's index and administrative state live under the main repository's `.git/worktrees/<name>/` rather than in the worktree.
A tool that writes at startup has its write refused and the session carries on, verified against the real CLI with a managed configuration loaded.

**`chmod` returning no error is not evidence that permissions are enforced.**
A filesystem can accept the call and ignore the bits, and `root` ignores them by design.
socrates-review attempts one write in the sealed worktree before spawning a reviewer, and refuses the review where that write succeeds.

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

**A page cannot work out its own base from the URL a person typed.**
The surface serves each session under `/sessions/<session>`, and the client builds `/answer`, `/live` and the page re-read from that base.
Deriving it from `document.URL` works only while the path carries a trailing slash, which is a property of what somebody typed rather than one the page holds.
The server states the base in a `data-` attribute beside the token, and the client reads it from there.

**`{$}` and a wildcard segment are different patterns to Go's `ServeMux`.**
`POST /sessions/{$}` matches `/sessions/` and nothing else, `GET /sessions/{session}` matches one non-empty segment, and neither shadows the other.
A page reachable both with and without a trailing slash needs `/sessions/{session}` and `/sessions/{session}/{$}` registered, because the first does not match a path ending in a slash.

## Testing

**One journey needs a filesystem that ignores permission bits, and skips without one.**
`TestASealThatDoesNotHoldIsARefusal` in `internal/e2e` is the only test of what `socrates-review review` does where the seal cannot be enforced, and no test can manufacture such a filesystem.
Name one and it runs:

```sh
SOCRATES_TEST_UNENFORCING_DIRECTORY=<a directory that ignores permission bits> go test ./internal/e2e/ -run TestASealThatDoesNotHoldIsARefusal
```

A FAT or NTFS volume mounted on a Unix host is one, and so is a Windows drive mounted into WSL2.
Recognise one by clearing a write bit and reading it back: `chmod a-w` there returns no error and `ls -l` shows the bit still set.
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
`exec` fails a script on any non-zero status and there is no assertion for a specific code, so `exec sh -c 'socrates-review status; echo exit=$?'` plus `stdout 'exit=3'` is how the contract gets checked.
`cp stdout file` is how you keep output for a later assertion, and `$WORK` is expanded inside single quotes.

## Continuous integration

**Code scanning cannot report while the repository is private and unlicensed.**
The analysis runs to the end and the upload is refused with "Advanced Security must be enabled for this repository to use code scanning", and enabling it over the API answers "Advanced security has not been purchased", which is what a private repository on a plan without it answers.
So nothing here analyses the code: a result no reader can reach is a weekly failure that teaches nobody anything, and a failure nobody acts on is worse than an absence somebody can see.
Code scanning is free for a public repository, so a workflow that analyses the code belongs with making this one public.

**Go's default per-package timeout is 10 minutes, and the end-to-end package outgrew it.**
The package that drives real binaries also produces a release snapshot, which compiles the releaser and cross-compiles three binaries for four targets.
On a macOS runner with a cold build cache that lands either side of 10 minutes, so the run failed as `panic: test timed out` in whichever test happened to be under way rather than in the slow one.
`make test` states its own timeout, and a package that hangs still fails rather than running forever.

**The releaser the `Makefile` pins has to build under the toolchain the module asks for.**
A verification run sets `GOTOOLCHAIN=local`, so `go run` cannot fetch a newer toolchain for a tool whose own `go.mod` asks for one.
Locally the switch is silent and the release snapshot succeeds; on a runner the same command fails with "requires go >= 1.27.0 (running go 1.26.7; GOTOOLCHAIN=local)".
Read the candidate's `go` directive before bumping `RELEASER`, and pin the newest version that stays at or below what `go.mod` requires.

**cosign 3 writes one bundle, and the old pair of output flags is refused rather than deprecated.**
`--output-signature` and `--output-certificate` still parse, and then the run fails with "create bundle file: open : no such file or directory", because `--new-bundle-format` is on by default and wants `--bundle`.
The release signs the checksums into `checksums.txt.bundle`, verified with `cosign verify-blob --bundle`, and the workflow pins the cosign release so a new default cannot change that shape under it.

**A private repository bills macOS minutes at ten times the rate, and the release is what stops first.**
Each verification run spends about eleven macOS minutes, billed as a hundred and ten, so a day of pushes reaches the account's spending limit.
What that looks like is not a slow queue: the job simply never starts, with no steps and a four second duration, annotated "The job was not started because recent account payments have failed or your spending limit needs to be increased".
The verification jobs had already run, so the tag exists and nothing is published, and rerunning the failed job fails the same way until the limit is raised or the repository is public.

**Build provenance is refused while the repository is private.**
`actions/attest-build-provenance` signs the archives and then fails to store the attestation, saying the feature is not available for the account and that the way to it is to upgrade the billing plan or make the repository public.
The step is left in the release and left failing, because it starts working the moment the repository is public and its failure is the reminder until then.
Everything before it publishes: the archives, the checksums, the cosign bundle and the bills of materials are all on the release.

**Private vulnerability reporting cannot be enabled while the repository is private.**
`PUT /repos/{owner}/{repo}/private-vulnerability-reporting` answers 404 rather than 403, which reads as a wrong path rather than as a feature that is not there.
It is available for public repositories only, so `SECURITY.md` names the channel and does not claim it is enabled, and enabling it belongs with making the repository public.

**A CodeQL job on a private repository needs `actions: read`.**
Without it the analysis runs to the end and then fails with "Resource not accessible by integration" against the workflow-runs endpoint, which reads as a bug in the action rather than as a permission.
`security-events: write` alone is enough on a public repository, and it is not enough on a private one.

**A checked box in a change's `tasks.md` is what makes its scenarios a gate.**
`make scenario-coverage` counts a change as under way as soon as one box is checked, and then every scenario in that change needs a test naming it.
Check the boxes and then run `make verify`, never the other way round, or the run fails on scenarios whose tests are still ahead of you in the same change.
