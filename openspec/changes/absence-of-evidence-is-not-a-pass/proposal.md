## Why

Two of socrates' checks report nothing when they are handed nothing, and reporting nothing is what a passing check looks like.

Four consecutive live rounds named the first one.
socrates takes the session's tool list from the stream's opening event, subtracts the tools it asked for, and records the remainder.
An opening event that carries no tool list at all leaves that remainder empty, so the run records nothing, and a person reading the run cannot tell a session that held exactly what socrates asked for from a session that never said.
The suite already knows the distinction matters: two tests fail with `the session reported no tools at all, so this test would pass on an empty reading` before they assert anything, which is exactly the guard the program does not have.
A vendor release that renames the field is all it takes.

The second one disarmed a gate this repository ran and read.
The drop gate compares a change's delta against the promoted specifications, and it was taught to follow a scenario rename the delta declares.
Where the promoted specification already holds the new name, there is nothing to follow: the rename matches nothing, the comparison sees every scenario carried, and the gate reports that the promotion would lose nothing.
That is what happened on the change before this one, whose promoted specifications were written at the same time as its delta.
`make scenario-drop` was run, ticked and reported nothing, and the five renames it was taught for matched no promoted name.

Neither of these is a wrong answer. Both are the absence of one, printed in the shape of a right one.

## What Changes

- **A tool list socrates could not read is recorded.** An opening event that reports no tools at all is its own event, distinct from a session holding a tool socrates did not ask for, and it does not fail the round: a release that stops emitting the field should report the truth rather than fail every review.
- **A rename the gate cannot check is reported.** Where a delta declares a scenario rename and the promoted specification already holds the new name, the gate says so and fails, because the alternative is a gate that passes while checking nothing.
- **The simulator can report no tool list**, so the first case is driven through the real binary rather than reasoned about.

## Capabilities

### Modified Capabilities

- `agent-seam`: what socrates records when the session does not say what it held.

The drop gate is tooling this repository holds with tests rather than a promoted capability, so it changes here without a specification delta.

## Impact

- `internal/agent`: one more event kind, recorded where the tool list is absent on a round that otherwise completed.
- `internal/presentation`: one more sentence.
- `internal/claudesim`: a scenario that omits the tool list.
- `internal/promotion`: one more condition, and the rename comparison that detects it.
- `cmd/scenario-drop`: the scenario line depends on the condition it belongs to.
- No stored run changes shape.
