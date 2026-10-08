# 9. socrates-review never learns what a pull request is

- Status: accepted
- Date: 2026-09-09
- Part of: reviewing a branch, and reviewing a change that is not mine

## Decision

No forge, no auth, no network, no `gh` in the requirements block.
The session checks the branch out in a throwaway worktree, captures the head commit, and posts the forwarded findings itself.
socrates-review takes a base, records a resolution, and reports it.
