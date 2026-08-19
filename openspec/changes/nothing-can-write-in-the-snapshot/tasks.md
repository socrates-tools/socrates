**The comparison after the round is not touched.** No task weakens it, adds an exception to it, or makes it conditional. The seal is a mechanism in front of it, not a replacement for it.

**Test first.** Every task names the scenario it covers, and is not done until a test named for that scenario failed for the right reason and then passed.

**`make lint test` after every task, not at the end of a group.**

## 1. The seal

- [ ] 1.1 Seal the worktree after checkout: clear the write bit on every directory and file in it, and unseal on every path that removes it, including the failure paths inside `Checkout`. Covers `A reviewer cannot create a file`, `A reviewer cannot modify a file`, `The reviewer can still read everything`
- [ ] 1.2 Prove the seal with one write into the sealed worktree before any reviewer is spawned, remove what it wrote where the write succeeded, and refuse the review naming the worktree and its filesystem. Covers `A seal that does not hold is a refusal`. Write the test that fakes an unenforcing filesystem first, because a probe that cannot fail is a probe nobody has seen work
- [ ] 1.3 Lift the seal when the round ends, whether the worktree is removed or kept. Covers `The seal is lifted when the round ends`, `The worktree is cleaned up`, `A dirtied worktree is kept`
- [ ] 1.4 Add the simulator knob that defeats the seal on demand: make the path writable, then write. Covers `A reviewer that wrote fails the round`, `A reviewer that created a file fails the round`. This is what keeps the comparison tested once an ordinary reviewer can no longer dirty the worktree
- [ ] 1.5 Add the end-to-end journey where a write is attempted and refused: the round completes and the verdict is the findings'. Covers `A write the seal refused is not a finding against the reviewer`, `A clean reviewer passes the check`

## 2. The machine that loads what it likes

- [ ] 2.1 Add the end-to-end journey for a reviewer whose session writes into its working directory as it starts, with the tool server restriction refused: the write does not land, the review completes, and the run's verdict is the findings'. Covers `A tool server that writes cannot reach the change under review`
- [ ] 2.2 Restate the read-only property in `AGENTS.md`, `README.md` and `docs/gotchas.md` as three mechanisms in the order they act, and say which of them the machine can take away

## 3. Verification

- [ ] 3.1 `make fmt && make lint test` clean, and `make scenario-coverage` naming no scenario this change adds
- [ ] 3.2 A live run on the machine whose enterprise configuration loads a tool server that writes: the review completes, the worktree is byte-identical to the snapshot commit afterwards, and the person's repository is untouched. Record what it found in this change's `retrospective.md`
