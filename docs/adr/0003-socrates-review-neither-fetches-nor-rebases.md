# 3. socrates-review neither fetches nor rebases

- Status: accepted
- Date: 2026-09-09
- Part of: reviewing a branch, and reviewing a change that is not mine

## Decision

Fetching reaches the network, which the suite asserts socrates-review does not do outside a build tag.
Rebasing writes refs, rewrites commits, touches the reflog and can leave the repository mid-conflict, against the invariant that nothing writes in the repository under review.
The session fetches and the person rebases.
socrates-review resolves the base.
Whether it says anything about how far the named ref has moved past the fork point is ADR 11, and the answer there is that it does not.
