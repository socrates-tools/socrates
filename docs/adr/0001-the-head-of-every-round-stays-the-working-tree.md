# 1. The head of every round stays the working tree

- Status: accepted
- Date: 2026-09-09
- Part of: reviewing a branch, and reviewing a change that is not mine

## Decision

The head is not the latest commit.
`internal/snapshot/snapshot.go` builds the tree with `git add --all`, so the tree is the working tree whether or not the author committed.
The skill tells the session to patch and record with `respond --fixed` and never tells it to commit, so a head of "latest commit" would hand round two a tree with no patches in it, and the run would escalate at the cap on defects that were fixed.

## Rejected

- head = latest commit.
