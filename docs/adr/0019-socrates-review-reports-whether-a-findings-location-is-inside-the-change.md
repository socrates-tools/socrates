# 19. socrates-review reports whether a finding's location is inside the change it reviewed

- Status: accepted
- Date: 2026-09-09
- Amends: ADR 18
- Part of: reviewing a branch, and reviewing a change that is not mine

## Decision

ADR 18 left the session to work out whether a forwarded finding can be a comment on a line, which reads as needing the pull request's diff and therefore a forge.
It does not.
A finding carries a file and a line, and `internal/reviewer` already computes the changed paths and the diff between the run's base and its snapshot, so whether a finding's line falls inside a hunk of the change is a local question about data socrates-review already holds.

ADR 2 is what makes the two agree.
A named base resolves to the fork point, which is what a forge computes a pull request's diff against, so the hunks socrates-review holds and the hunks a forge renders are the same hunks.

So `findings --json` reports, per finding, whether its location is inside the change.
The session maps that to where a comment goes: inside means a comment on the line, outside means the one general comment.

The placement rule therefore has one owner and is code rather than prose a session re-derives on every review.
The card's promise in ADR 18, that it says which kind of comment a finding will become, is sourced from socrates-review rather than from the session's own arithmetic.
A mistake in it fails a test in this repository instead of landing in somebody else's.

ADR 9 is unchanged.
socrates-review does not learn what a pull request is, does not authenticate anywhere, and posts nothing.
It reports a property of its own diff.
