# 7. One round is derived rather than capped

- Status: accepted
- Date: 2026-09-09
- Part of: reviewing a branch, and reviewing a change that is not mine

## Decision

The round cap stays two and stays unconfigurable.
A handover round parks at `awaiting-response` because forwarding addresses the findings, and `socrates-review continue` ends the run rather than starting a round, on the existing path where it already does that for a run with nothing left to review.
`dialogue.After` is untouched.
`dialogue.AwaitsAPerson` gains the handover fact, because every handover round waits on a person where a round of `auto-fix` findings does not.
