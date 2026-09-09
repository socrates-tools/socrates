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

`socrates-review continue` is unchanged: a round did happen there, the person answered it, and a patched tree that holds nothing its base does not means the change was reverted rather than never read.

## Rejected

- ending the run `failed`, which writes a run somebody has to read to learn that nothing happened; and leaving `clean`, which is the false pass.
