# 21. Each skill carries its own ceiling

- Status: accepted
- Date: 2026-09-09
- Part of: reviewing a branch, and reviewing a change that is not mine

## Decision

The handover loop is a second generated skill rather than a branch inside the existing one, so a session reviewing its own work does not read the forwarding half in the session where its context is the scarce resource.
Two documents therefore need two budgets.

`skill.Ceiling` bounds the skill a session reads about its own change, and `skill.HandoverCeiling` bounds the one it reads about somebody else's.
Each is a number stated once beside the document it bounds, and `skill.All()` is what the suite reads, so a skill added without a ceiling fails.
The handover document is larger because it teaches strictly more: the same review loop, the same page, and the discipline for posting a comment in a repository nobody here owns.

One ceiling covering both would have been one number, and the number would have had to be the larger of the two, which lets the document read by every reviewing session grow by the difference without anybody noticing.

## Rejected

- one ceiling for both documents, and no ceiling on the second.
