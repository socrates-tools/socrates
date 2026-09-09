# 8. The run ends escalated, with a third reason

- Status: accepted
- Date: 2026-09-09
- Part of: reviewing a branch, and reviewing a change that is not mine

## Decision

`escalated` already means the findings are the person's to decide and socrates-review has nothing further it can do, which is what a handover run reaches.
No new state joins, so the skill's promise to list every state a run rests in holds unchanged, and the exit-code table, `internal/presentation`, `report` and `follow` are untouched.
Known cost: `escalated` reads as alarm, and a round where every finding was rejected as false ends in the same state as one where four were forwarded.

## Rejected

- a state of its own, and a verdict derived from the dispositions.
