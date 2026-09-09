# 16. The author's own account of their change reaches the reviewer, relayed by the session

- Status: accepted
- Date: 2026-09-09
- Part of: reviewing a branch, and reviewing a change that is not mine

## Decision

The reviewer cannot read a pull request: it is spawned into a sealed worktree of a frozen commit and sees the prompt and the tree.
The session reads the title and description and passes them to `--intent`.
One flag serves both cases, because `prompt.IntentHeading` is already "What the author says they were doing", which is author-relative and correct for somebody else's change.
The attribution sentence in the section body is not: it says the person who started the review stated the intent, which under `--not-my-change` is false and changes what a disagreement means.
Your own intent against the code says the implementation missed the goal.
The author's description against the code says their account of their own work is incomplete or misleading, which is the more serious finding and is only available if the reviewer knows who is speaking.
So the section names the party the run record says stated it.
`--intent` is neutralised through `internal/hygiene` unconditionally rather than only under the flag, because a channel that is sometimes neutralised is a channel nobody can reason about, and neutralising a typed one-liner costs nothing.
