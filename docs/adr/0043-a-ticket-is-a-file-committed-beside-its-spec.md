# 43. A ticket is a file committed beside its spec

- Status: accepted
- Date: 2026-09-30
- Part of: socrates-tickets

## Decision

`socrates-tickets` writes a change's tickets as Markdown files under `docs/changes/<slug>/tickets/`, committed beside the spec they were cut from, and publishes nothing to a tracker.

The change folder already holds what is only ever about one change, and a ticket is exactly that.
Committed is the part that matters: whoever implements a ticket may be working in a worktree of its own, and a worktree sees what is committed and nothing else.
Marking a ticket done is therefore a commit, and the history of a change records the order its tickets landed in.

This is surprising beside ADR 39, which says the tracker holds the tickets because a committed file cannot carry a state, an assignee or a blocking edge.
The tickets `socrates-tickets` writes answer that differently: the main ticket carries the state, and each ticket names its own blockers, per ADR 44.
Whether this repository's own earlier practice follows was put to the person, and they left it outside this change in their own words: "Don't worry about the existing state of this repo. [...] whatever is done as github issue or not, doesn't matter to us, we only need to make sure the work of socrates-tickets feature follows socrates-design flow".

## Rejected

- `docs/changes/<slug>/`, flat beside `spec.md`.
  Committed and with no subfolder, so the tickets and the spec share one listing.
  Declined without a word.
- `.scratch/<slug>/tickets/`, uncommitted, which is the shape `to-tickets` writes.
  Nothing lands in history, and an implementer in a worktree of its own cannot read the tickets without being handed a path outside the repository.
  Declined without a word.
