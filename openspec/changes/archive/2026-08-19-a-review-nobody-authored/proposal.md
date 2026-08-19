## Why

A coding session that implements a change is the worst possible judge of it.
It has every reason to believe the code is right, it remembers what it meant rather than what it wrote, and asking it to check its own work returns a confident summary of its own intentions.

The fix that works is an independent reviewer, and the reason people do not have one is that arranging it is tedious.
The reviewer has to look at a tree that is not moving while the author keeps working, it has to be a process with no memory of having written anything, it has to return something a machine can route rather than an essay, and none of that is worth setting up by hand for one change.

`no-mistakes` solves this and a great deal more, and the "more" is where its complexity lives: a gate repository, a daemon, a landing phase, a forge, a multi-user trust boundary.
socrates takes the review and leaves the rest.

This change is the whole of the core: freeze the tree, ask somebody who did not write it, come back with findings a machine can act on, and say honestly when no review happened.
Everything after it is a layer on this.

## What Changes

- **`scrat review` becomes a thing that exists.** It freezes the working tree into a commit nothing can move, spawns a detached reviewer against a worktree of that commit, and returns the run identifier immediately. The person's session keeps working; the reviewer reads a tree that cannot change under it.
- **The snapshot includes untracked files and touches nothing.** It is built through a temporary index, so the person's index and working tree are byte-identical afterwards, and a feature whose new files are all untracked is reviewed rather than half-reviewed.
- **Findings are a schema, not prose.** Severity says whether the run blocks, action says who resolves it, and the two are orthogonal. An unclassified finding defaults to needing a person, and the classification is one helper rather than two rules that can disagree.
- **Every model call goes through one seam, with one implementation.** The reviewer role resolves to its own model and effort, because the reviewer wants the strongest model available and nothing else here does.
- **The reviewer never writes, and socrates proves it rather than asking.** The snapshot worktree is compared against the snapshot commit when the round ends, and a difference fails the round.
- **A review that did not happen is never a pass.** A missing findings file, a reviewer that exited non-zero, a schema violation and an empty stream are all `failed`, and `failed` is a different exit code from `clean` and from `escalated`.
- **`scrat status` exists and is deliberately tiny**, because the session waiting on the review is what calls it most and everything it prints lands in that session's context.
- **The one third party outside our control gets a simulator.** `cmd/claudesim` speaks the real `claude` CLI's wire contract, is driven by scenario files, and records every invocation. It is what makes an end-to-end suite possible at all, and its fidelity to the real tool is itself checked against recordings.

## Capabilities

### New Capabilities

- `snapshot`: freezing a working tree, untracked files included, into a commit and a detached worktree without touching the person's index, working tree or refs.
- `run-store`: the run directory, run identity, run state, the bounded append-only event log, and the atomic writes that keep a half-written run unreadable rather than misleading.
- `agent-seam`: the single narrow interface every model call passes through, the per-role model and effort resolution, and the Claude Code implementation behind it.
- `reviewer-run`: spawning the reviewer detached, the process group it is reaped by, the re-entrancy marker, the read-only enforcement, and the positive evidence a completed review must produce.
- `review-findings`: the findings contract, the orthogonality of severity and action, the fail-closed defaulting, and the single classification helper both consumers share.
- `reviewer-prompt`: the wording the reviewer is given, which was arrived at through incidents rather than design and is the cheapest high-value thing in the system.
- `cli-foundation`: the command surface, the run selection rule, the exit codes, and the configuration file.

### Modified Capabilities

None. This is the first change.

## Impact

- New Go module `github.com/rui-ktei/socrates`, Go 1.25 or later, one binary at `cmd/scrat`.
- New simulator binary at `cmd/claudesim`, built only by the test suite and never shipped.
- Dependencies: `spf13/cobra`, `go.yaml.in/yaml/v3`, `google/go-cmp`, `rogpeppe/go-internal`. Nothing else.
- Writes under the user's state directory only. Inside a repository it writes nothing at all: no gate remote, no hook, no config file it creates on its own, no commit on any branch.
- `git` is the only external system, reached as a subprocess. The `claude` binary is the only other process started, through the agent seam.
