# 18. A forwarded finding lands on its line where the pull request has that line, and otherwise in one general comment

- Status: accepted
- Date: 2026-09-09
- Part of: reviewing a branch, and reviewing a change that is not mine
- Amended by: ADR 19, which moves the question of whether a finding is on a changed line from the session to socrates-review

## Decision

A forge only pins a comment to a line that appears in the diff, and the reviewer reads the whole repository, so findings divide in two.
A finding whose location is in the diff becomes a comment on that line.
Every other forwarded finding joins a single general comment on the pull request, one rather than one each, so it does not read as a pile-on.
The card on the page says which of the two a finding will become before it is forwarded, because a precise note on a line and a paragraph in a general comment are different acts: the second is easier to miss and several of them read as a lecture about the author's architecture rather than a review of their change.

## Rejected

- everything in one general comment, which discards the position of the precise findings; and refusing to forward findings off the diff, which discards the class of finding that is often worth the most.
