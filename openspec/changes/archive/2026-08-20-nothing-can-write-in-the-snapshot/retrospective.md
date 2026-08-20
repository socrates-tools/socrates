# Retrospective: nothing-can-write-in-the-snapshot

## The live run

Task 3.2 asks for one review on the machine whose enterprise configuration loads a tool server that writes, which is the machine that made this change necessary.
The change reviewed was this one: the seal, the probe, the refusal, the lift, the simulator knob, the two journeys and the documentation, as a working tree against its parent commit.
The reviewer ran for about eleven minutes at `opus` and `high` with the machine's tool servers loaded, and returned four findings and a `medium` risk rating.

The previous change could not complete a review on this machine at all.
This one did, and the observation that matters was taken while it was still running:

```
dr-xr-xr-x 9 user user 4096 .../worktree
```

No `.serena` directory, with the same enterprise configuration loaded and the same server started.
`agent-tool-servers-loaded` is in the run's event log, so strictness was refused exactly as before, and the seal held anyway.
Afterwards the worktree was removed, `git worktree list` in the reviewed repository named only the repository itself, and `git status` was identical to before the review, line for line.

That is the whole claim of this change, observed rather than asserted.

## What one live review found that a green suite did not

`make fmt`, `make lint`, `make test` and `make scenario-coverage` were clean before the run, and the guarded journey had been run against a real unenforcing filesystem.
The review still found two defects and two overclaims.

### Every review would refuse when socrates runs as root

The probe treats a successful write as proof that the seal does not hold.
Root holds `CAP_DAC_OVERRIDE`, so that write succeeds on an ordinary filesystem, and socrates would refuse every review in a root container or a root CI image.
Worse, the refusal blamed the filesystem and told the person to point `SOCRATES_STATE_HOME` somewhere else, which cannot help a process that is root.

Refusing under root is the right answer and stays: as root the seal is not a mechanism at all.
The refusal now names both conditions and the next action for each.
There is one message rather than a branch on the effective user identity, because a branch that only root can enter is a branch no test here can enter.

### An abandoned run left a directory nothing could remove

A supervisor killed mid-round leaves `run.json` non-terminal on disk.
The retention sweep skips non-terminal runs on purpose, so it never reclaimed it, and `scrat cancel` returned early on a run that read terminal only because the reviewer had been observed to be gone.
Nothing lifted the seal, and `rm -rf` on the directory is refused.
Before the seal that was untidy; with the seal it is a directory a person cannot delete without `chmod -R` first.

`scrat cancel` now finishes what the dead supervisor could not: it writes the observed failure to disk, so the sweep can reclaim the run, and it lifts the seal, so a person can read or delete what is left.
`TestCancellingARunWhoseReviewerWentAwayLiftsTheSeal` failed for both reasons before the fix.

### The seal was described as something it is not

`AGENTS.md` and the specification both said a write is refused by the filesystem, full stop.
The session owns those files, so something that sets out to defeat the seal can clear the bit back and write, which is precisely what the simulator's own `defeats_the_seal` knob does.
The rest of the change was honest about this, because the comparison after the round exists to catch exactly that, but the invariant list is what a later change reads as the contract.
Both now say the seal is not a sandbox, and what it removes is every write nothing went out of its way to make, which is the whole of what has actually written in a snapshot worktree so far.

### Sealing and lifting are not inverses

`docs/gotchas.md` claimed the write bits round-trip.
They do for `0o755`, and they do not for `0o664` or for a deliberately read-only `0o444` file, which both lift to `0o644`.
Nothing depends on the modes a lifted worktree carries, so the fix is the sentence rather than the code, and the sentence now says which way it is not an inverse.

## What the design got right

The probe earned its place before a person ever saw it.
It is the reason a mutation that removed the seal failed nine tests instead of quietly reviewing an unsealed worktree, and it is the reason the root case is a refusal with a message rather than a silent loss of the property.

Keeping the comparison was also right, and for a reason the proposal only half stated.
It is not just insurance against a future mechanism: it is what makes the honest description of the seal possible.
Because the comparison is there, the seal does not have to be a sandbox, and nobody has to pretend it is one.
