# 13. This stays one tool

- Status: accepted
- Date: 2026-09-09
- Part of: reviewing a branch, and reviewing a change that is not mine

## Decision

The premises behind the core design survive a review of somebody else's change.
The tree can still move under the reviewer, so the snapshot earns its place, and a vendor session still writes, so the seal does.
The no-inherited-session rule and the round cap's reason are moot where there is no second round rather than violated.
No forge, no auth and no network cross the boundary, so nothing about a handover review needs a program of its own the way `agora` did.
