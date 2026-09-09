# 4. A finding on a change that is not mine is forwarded, dropped or rejected

- Status: accepted
- Date: 2026-09-09
- Part of: reviewing a branch, and reviewing a change that is not mine

## Decision

`dialogue.Accepts` at `internal/dialogue/responses.go` returns `[Fixed, Rejected]` or `[Decided]` by the finding's own action.
Where the run has no authority to patch the head, it returns `[Forwarded, Dropped, Rejected]` for every finding regardless of action, because nothing here can be fixed and no question about intent is mine to answer.
`forwarded` is new and means the finding is handed to whoever wrote the change.
`dropped` is new and means the finding is true and I am not raising it.
`rejected` is unchanged and still means the finding is false.
