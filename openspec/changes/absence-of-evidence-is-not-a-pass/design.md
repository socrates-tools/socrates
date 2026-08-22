## Context

Both changes here are to checks rather than to behaviour, and both are the same mistake made twice: a check whose empty answer is indistinguishable from its clean answer.

## Goals / Non-Goals

**Goals:**

- A run says whether socrates could read the session's tool surface.
- A gate that cannot check what it was asked to check says so instead of passing.

**Non-Goals:**

- Failing a review because the tool list is absent. Rejected below.
- Reading the opening event's other fields to make up for it. They are an inventory of the machine rather than a statement about the session, which is already recorded.
- Making the drop gate read prose. It compares headings, and what a delta's sentences say is still a reader's job.

## Decisions

### An absent tool list is recorded, and does not fail the round

The tool list bounds nothing. It is an observation socrates records for a person, and the requirement it belongs to already decides what to do when the observation is unwelcome: run the review, record what was found.

The recorded reasoning for the case next door settles this one. A vendor release that adds a built-in outside the set socrates asked for makes every run record an unexpected tool until somebody looks, and that was chosen because it is one line, true when it appears, and the direction to be wrong in. A release that stops emitting the field is the same shape, and the same answer follows: one line, true when it appears.

The alternative fails every review on a field rename, which is the thing the same requirement already refuses by name: a rule that stops working by failing every review rather than by reporting the truth.

The event is its own kind rather than an empty tool list on the existing one, because the two say different things to a person: one is a session that held something socrates did not ask for, and the other is a session that did not say what it held.

It is recorded only where the round otherwise completed. A round that failed before the reviewer's first byte has no tool list either, and the failure is the thing to read there; recording the absence beside it would put a second line on every failure that says less than the first.

### The gate reports a rename it cannot check

The signature is exact: a delta declares a scenario rename, the promoted requirement holds the new name, and it does not hold the old one. There is nothing for the rename to follow, so the comparison it was taught for is not happening.

It fails the gate rather than warning, for the reason the gate exists: a check nobody can tell apart from a passing one is worse than no check, because it gets run and ticked.

A requirement rename in the same state already fails, through the condition that says the requirement it modifies is in no promoted specification. Only the scenario case was silent.

The identical `FROM`/`TO` pair is left alone. It declares nothing and follows nothing, and a gate that refused it would be refusing a no-op rather than protecting anything.

## Risks / Trade-offs

- **A machine whose vendor release drops the field records a line on every run.** → Deliberate, stated above, and it is the direction the requirement next door already chose.
- **The new gate condition fires on a legitimate workflow somebody uses.** → The workflow it fires on is editing the promoted specifications while the change that promotes them is still open, which is what disarmed the gate. If somebody wants that, the change is ready to archive and archiving it is the way to say so.
- **Two unrelated checks in one change.** → They are the same defect, they were found by the same reviewer in the same rounds, and separating them would put the second one back on the list nobody gets to.

## Migration Plan

No stored run changes shape, and no promoted scenario moves.

## Open Questions

None.
