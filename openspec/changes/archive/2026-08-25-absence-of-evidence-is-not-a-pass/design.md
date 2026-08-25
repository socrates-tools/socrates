## Context

Everything here is a change to a check rather than to behaviour, and all of it is the same mistake made three times: a check whose empty answer is indistinguishable from its clean answer.

## Goals / Non-Goals

**Goals:**

- A run says whether socrates could read the session's tool surface.
- A run says what the session reported having connected, not only what it reported holding.
- A gate that cannot check what it was asked to check says so instead of passing.

**Non-Goals:**

- Failing a review because the tool list is absent, or because a source was named. Rejected below.
- Reading the opening event's inventory of the machine. Its `skills`, `plugins` and `agents` fields name what is installed rather than what the session has, so a check over them would report a customisation on every run of every machine that has one. The field naming the sources the session connected is not one of those, which is the decision below.
- Making the drop gate read prose. It compares headings, and what a delta's sentences say is still a reader's job.

## Decisions

### An absent tool list is recorded, and does not fail the round

The tool list bounds nothing. It is an observation socrates records for a person, and the requirement it belongs to already decides what to do when the observation is unwelcome: run the review, record what was found.

The recorded reasoning for the case next door settles this one. A vendor release that adds a built-in outside the set socrates asked for makes every run record an unexpected tool until somebody looks, and that was chosen because it is one line, true when it appears, and the direction to be wrong in. A release that stops emitting the field is the same shape, and the same answer follows: one line, true when it appears.

The alternative fails every review on a field rename, which is the thing the same requirement already refuses by name: a rule that stops working by failing every review rather than by reporting the truth.

The event is its own kind rather than an empty tool list on the existing one, because the two say different things to a person: one is a session that held something socrates did not ask for, and the other is a session that did not say what it held.

It is recorded only where the round otherwise completed. A round that failed before the reviewer's first byte has no tool list either, and the failure is the thing to read there; recording the absence beside it would put a second line on every failure that says less than the first.

### The sources the session named are read beside the tools it held

The change before this one decided the opposite, and the sentence it decided it in is the one to answer: reading the source list instead "would swap one incomplete signal for another, and would keep the vocabulary this change removes".

The first half is true of a swap and this is not one. Both readings come from the same event and both are incomplete, in different places: the tool list misses a source whose tools have not arrived, and the source list misses a built-in the vendor added in a release nobody read. Reading both leaves neither gap. Reading one leaves the gap this repository measured: a probe on this machine named four sources in an opening event whose tool list held none of their tools, because sources connect asynchronously.

The second half is a real constraint and it is met. The vendor's field name appears once, in the JSON tag of the wire struct, which is where every other vendor field name already lives beside `structured_output` and `cache_read_input_tokens`. No requirement, event kind, recorded field, rendered line or exported name says anything but source.

What makes this field worth reading where the event's other fields are not is that it is a statement about the session rather than an inventory of the machine. That is measured, not assumed: under the one mechanism the sources named went from four to none, while the `plugins` inventory beside it did not move. So a run on a working machine records nothing, and a run that records a source is a run where the mechanism did not hold. That is the same shape as the tool-list check and it is there for the same reason.

Recording a source does not fail the round, for the reason the absent tool list does not: the observation bounds nothing, and a machine whose policy is set by somebody entitled to set it should still get reviews.

Unlike the absence, it is recorded whether or not the round completed. An absence on a failed round says less than the failure beside it; a source the session named is positive evidence, it is true whichever way the round went, and a reviewer that died holding a connected source is exactly the run somebody will want to read.

### The gate reports a rename it cannot check

The signature is exact: a delta declares a scenario rename, the promoted requirement holds the new name, and it does not hold the old one. There is nothing for the rename to follow, so the comparison it was taught for is not happening.

It fails the gate rather than warning, for the reason the gate exists: a check nobody can tell apart from a passing one is worse than no check, because it gets run and ticked.

A requirement rename in the same state already fails, through the condition that says the requirement it modifies is in no promoted specification. Only the scenario case was silent.

The identical `FROM`/`TO` pair is left alone. It declares nothing and follows nothing, and a gate that refused it would be refusing a no-op rather than protecting anything.

## Risks / Trade-offs

- **A machine whose vendor release drops the field records a line on every run.** → Deliberate, stated above, and it is the direction the requirement next door already chose.
- **The new gate condition fires on a legitimate workflow somebody uses.** → The workflow it fires on is editing the promoted specifications while the change that promotes them is still open, which is what disarmed the gate. If somebody wants that, the change is ready to archive and archiving it is the way to say so.
- **A machine that loads a source anyway now records a line on every round of every run.** → That machine is one where the one mechanism did not hold, which is the thing worth a line. It is the same trade the tool-list check already takes, and the run still completes.
- **A vendor release that renames the source field makes the reading silent again.** → The guard in the ordinary suite pins the field against the checked-in recordings, so its disappearance is a diff somebody reads rather than a check that quietly stops. That is what the tool list already has and what this borrows.
- **Three unrelated checks in one change.** → They are the same defect, two of them were found by the same reviewer in the same rounds and the third is the first one one field along, and separating them would put the leftovers back on the list nobody gets to.

## Migration Plan

No stored run changes shape, and no promoted scenario moves.

## Open Questions

None.
