# 5. A resolution that leaves no artifact carries a reason

- Status: accepted
- Date: 2026-09-09
- Part of: reviewing a branch, and reviewing a change that is not mine

## Decision

`fixed` needs no reason because the patch is the evidence.
`forwarded` needs none because the comment carries the finding's own text.
`dropped` and `rejected` leave nothing behind, so each carries a reason and is refused without one.
This generalises the existing rule at `internal/dialogue/responses.go` rather than adding a second one.
A reason is required for every drop, not only above a severity, because a rule keyed on a severity a model filled in is a rule the model can waive.
