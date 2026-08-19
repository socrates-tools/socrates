## Why

socrates cannot complete a review on the machine it was written on.

The reviewer's tools are restricted to reading, and socrates asks the vendor's CLI to load no external tool server.
Both mechanisms are refused by the machine: an enterprise tool server configuration exists at a path outside the person's control, the CLI declines `--strict-mcp-config` while it does, and one of the servers named there writes a project directory into its working directory the moment it starts.
The working directory of a reviewer is the snapshot worktree, so the round ends with the worktree differing from the snapshot commit, and the comparison correctly fails the run for a write the reviewer never asked to make.

The comparison is not wrong and the person's configuration is not wrong.
What is wrong is that socrates has two mechanisms that can both be taken away from it, and a third that only reports the damage afterwards.
`no-mistakes` does not have this problem because it never asks the question: its reviewer runs with the full tool set in a worktree nothing checks, so infrastructure droppings are invisible there and are committed by the next `git add -A`.
Being unable to notice is not a property worth copying.

The remaining mechanism is the filesystem.
A worktree whose directories are not writable cannot be written to by a reviewer, by a tool server, by a hook, or by any part of a session socrates does not control and cannot enumerate.
It was verified against the real CLI on the affected machine, with the enterprise servers loaded: the reviewer read the change, returned a successful result, and the worktree afterwards held exactly the snapshot's contents.

## What Changes

- The snapshot worktree is checked out and then sealed: every directory and file in it is made unwritable for the duration of the round.
- Sealing is proved rather than assumed. socrates attempts one write into the sealed worktree before any reviewer is spawned, and refuses the review if that write succeeds, naming the filesystem as the reason.
- The seal is lifted when the round ends, so a worktree kept for inspection is an ordinary directory a person can read and delete.
- The comparison after the round stays exactly as it is. It is now the check for what defeated the seal rather than the only thing standing between a reviewer and the change under review.
- The specifications stop resting the read-only property on the vendor's flags. Three mechanisms are named in the order they act: the vendor's tool restriction, the seal, and the comparison.

## Capabilities

### New Capabilities

<!-- None. This change modifies requirements that already exist. -->

### Modified Capabilities

- `snapshot`: the worktree the reviewer is given is sealed against writing, the seal is proved before a reviewer is spawned, and it is lifted when the round ends.
- `reviewer-run`: the reviewer never writes because it cannot, and the comparison after the round becomes the check on what defeated that.
- `agent-seam`: where the CLI refuses to ignore the machine's tool servers, what holds the reviewer read-only is the seal rather than the after-the-fact comparison.

## Impact

- `internal/snapshot`: sealing, the write probe and its refusal, unsealing on every removal and failure path.
- `internal/reviewer`: lifting the seal when the round ends, whether the worktree is removed or kept.
- `internal/claudesim`: a scenario knob that defeats the seal on demand, so the comparison stays testable now that an ordinary reviewer cannot dirty the worktree.
- `internal/e2e`, `internal/snapshot` tests: the write-detection tests change how they dirty the worktree.
- `AGENTS.md`, `README.md`, `docs/gotchas.md`: the third mechanism, and why the first two are not enough.
- Windows is not affected in this change beyond best effort: the probe is what decides whether socrates will run, so a filesystem that cannot express the seal is a refusal rather than a quiet loss of the property.
