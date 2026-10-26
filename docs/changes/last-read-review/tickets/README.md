# The last read puts socrates review in front of the whole change

This change has the overseer of `socrates-implement` put socrates review in front of the whole change once every ticket has landed, with nobody asked anything.
socrates review gains a run with nobody to ask, on which a session settles a finding that asks a person itself and the record says so, and the overseer's last read uses it, hands every fix to one implementer and says at hand-back how the review ended.

Spec: [../spec.md](../spec.md). Decisions: [ADR 90](../../../../implement/docs/adr/0090-the-last-read-puts-socrates-review-in-front-of-the-whole-change.md), [ADR 91](../../../../review/docs/adr/0091-a-decision-the-session-made-is-settled-and-only-a-run-with-nobody-to-ask-takes-one.md).

## Working through the tickets

- Commit this folder before handing out any ticket.
- A ticket may start once every ticket it names as a blocker is done.
- Mark a ticket in progress before handing it out.
- Mark it done once its acceptance criteria hold and its work is committed.
- Mark it abandoned, with a reason, where building it turns out to be wrong, and in the same change every ticket it blocked takes its blockers in its place.

## Tickets

| Ticket | Blocked by | Status |
| --- | --- | --- |
| [1: A run with nobody to ask takes a settled answer](1-settled-on-a-run-with-nobody-to-ask.md) | None | in progress |
| [2.1: The second round reads a settled answer as the session's decision](2.1-second-round-reads-settled-as-the-sessions-decision.md) | 1 | pending |
| [2.2: The last read reviews the whole change](2.2-last-read-reviews-the-whole-change.md) | 1 | pending |
