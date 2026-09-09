# 2. A named base resolves to the fork point, not to the named commit

- Status: accepted
- Date: 2026-09-09
- Part of: reviewing a branch, and reviewing a change that is not mine

## Decision

`internal/snapshot/snapshot.go` diffs with `git diff <base> <commit>`, which is two-dot.
A base of `origin/develop` when develop has moved renders every file develop changed and the branch did not as a deletion the branch never made.
`--base <ref>` therefore resolves to `merge-base(<ref>, HEAD)`, so the diff is the branch's own work.
After a rebase onto `<ref>` the two readings coincide, which is why rebasing first is the right instinct rather than a workaround.
