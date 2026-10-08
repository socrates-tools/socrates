# 11. socrates-review says nothing about how far the named ref has moved

- Status: accepted
- Date: 2026-09-09
- Part of: reviewing a branch, and reviewing a change that is not mine

## Decision

Tell nobody.
The reviewer is not told, because staleness is a property of what was asked for rather than of the code, and a reviewer told its information is stale has a standing excuse for every miss.
The person is not told either, because the intended workflow rebases onto the named ref first, which makes the gap zero and any report of it noise.
Known cost: nothing enforces the rebase, so a review against a months-old fork point looks exactly like a review against a current one, and can report clean.
`status` prints the base commit it resolved, as a bare fact and not a warning.
Decided after I described the "tell nobody" option as leaving the base discoverable on `status`, which was wrong: `report.Status` has no such field and `Status.Lines()` prints no such line.
So nobody is warned about the gap, and anybody can see which commit was reviewed against.
