# 10. A repeat review of the same pull request is not deduplicated

- Status: accepted
- Date: 2026-09-09
- Part of: reviewing a branch, and reviewing a change that is not mine

## Decision

Findings have no identity across runs, so matching a finding from one run to a finding from the next needs a model judging whether two descriptions mean the same thing and then withholding one, which is the thing the design refuses hardest.
The session instead reads the existing comment threads, matches on file and line only, and puts what it finds on the card as evidence: there is already a comment on this line, from this date, saying this.
The reviewer still reads the change cold, nothing is suppressed, and the person drops the duplicate themselves with the reason a drop now requires.
Repeat reviews are therefore in scope for the first version.

## Rejected

- fuzzy matching of finding text, and treating repeat review as out of scope.
