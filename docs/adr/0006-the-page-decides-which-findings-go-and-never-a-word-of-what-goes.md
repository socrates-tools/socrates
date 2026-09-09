# 6. The page decides which findings go and never a word of what goes

- Status: accepted
- Date: 2026-09-09
- Part of: reviewing a branch, and reviewing a change that is not mine

## Decision

Every option on the page stays a `socrates-review respond` invocation, so the skill's guard is untouched.
The comment body is read back from `socrates-review findings --json`, which is socrates-review's own recorded text.
A tampered page can therefore forward a finding nobody picked, and cannot post model-authored prose to another person's repository under the reviewer's name.

## Rejected

- an option whose command is a forge invocation.
