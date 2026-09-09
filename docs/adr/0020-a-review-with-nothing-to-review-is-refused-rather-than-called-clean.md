# 20. A review with nothing to review is refused rather than called clean

- Status: accepted
- Date: 2026-09-09
- Part of: reviewing a branch, and reviewing a change that is not mine

## Decision

`socrates-review review` used to end a run `clean` where the working tree held nothing its base commit did not, with a recorded reason saying so.
The reason was honest and the state was not: no reviewer was spawned, nothing was read, and a session acting on the state reported that an independent reviewer had found nothing blocking.
A finished branch with everything committed reached that path every time, which is what the problem statement behind ADR 2 called a false pass on a whole feature.

So `review` refuses instead, before a run is created, naming the base it resolved and naming `--base` as the way to review a branch's own work.
It exits with the usage code, so nothing about it can be read as a verdict.
No run directory is written, so there is no state for a session to act on and nothing for `status` to report.

`socrates-review continue` reaches the same condition where the answers reverted the whole change, and it ends the run `failed` rather than `clean`.
A run cannot be refused there: it exists, it is parked, and refusing would leave it parked with nothing that ends it.
`failed` is what the state means, because no reviewer read the round and nothing about the change was checked, and the recorded reason says both.
Round one is kept and readable either way.

One rule covers both paths: socrates-review reports `clean` for no round a reviewer did not read.

## Rejected

- refusing the review with a run that records the refusal, which writes a run somebody has to read to learn that nothing happened; leaving `clean` in either place, which is the false pass; and refusing `continue`, which leaves the run parked with nothing left that ends it.
