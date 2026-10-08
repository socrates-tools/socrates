# 44. A ticket names its blockers, and its number is derived from them

- Status: accepted
- Date: 2026-09-30
- Part of: socrates-tickets

## Decision

Every ticket names the tickets that block it, and that list is the whole dependency graph.
The number in a ticket's name, `{N}-{task}.md`, is its wave: one more than the highest wave among its blockers, and 1 for a ticket with none.
A wave holding one ticket names it `N`, and a wave holding several names them `N.1`, `N.2` and on.

The name orders the folder and gates nothing.
A ticket starts once its own blockers are done, so `3.1` blocked only by `2.1` starts while `2.2` is still running.
That is the surprise a reader needs warning of: `3.1` looks as though it waits on every ticket of wave 2, and it does not.

Because the wave is computed from the blockers rather than chosen, the name and the blocker list agree when the tickets are written.
The name is fixed from then on.
Abandoning a ticket during implementation passes its own blockers to every ticket it blocked, so nothing starts earlier than it would have, and renames nothing, because renaming would break every link to the renamed tickets and every Blocked by entry naming them.
So a name states the wave its ticket was cut into, not the wave its blockers would give it after something was abandoned, and the person decided that in review of this change rather than renaming on every abandon.

## Rejected

- The number is the whole graph: every ticket of wave N waits on every ticket of wave N-1, and no ticket states an edge.
  Simpler to read, at the cost of parallelism, because a ticket waits on work it does not need.
  Declined without a word.
