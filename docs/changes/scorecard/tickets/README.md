# OpenSSF Scorecard runs weekly and the README carries its badge

A weekly Scorecard workflow publishes its results and uploads them to code scanning, and it is the one workflow that starts without being asked.

Spec: [../spec.md](../spec.md). Decisions: [ADR 64](../../../adr/0064-scorecard-is-the-one-run-nobody-asks-for.md).

## Working through the tickets

- Commit this folder before handing out any ticket.
- A ticket may start once every ticket it names as a blocker is done.
- Mark a ticket in progress before handing it out.
- Mark it done once its acceptance criteria hold and its work is committed.
- Mark it abandoned, with a reason, where building it turns out to be wrong, and in the same change every ticket it blocked takes its blockers in its place.

## Tickets

| Ticket | Blocked by | Status |
| --- | --- | --- |
| [1: Scorecard runs weekly, publishes, and is the one run nobody asks for](1-scorecard.md) | None | todo |
