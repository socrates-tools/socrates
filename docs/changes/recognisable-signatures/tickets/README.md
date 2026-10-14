# A release's signature and provenance carry the names verifiers look for

The signature is published as `checksums.txt.sigstore.json`, and the provenance is attached as `socrates.intoto.jsonl`.

Spec: [../spec.md](../spec.md). Decisions: [ADR 65](../../../adr/0065-a-release-names-its-signature-and-provenance-the-way-verifiers-look-for-them.md).

## Working through the tickets

- Commit this folder before handing out any ticket.
- A ticket may start once every ticket it names as a blocker is done.
- Mark a ticket in progress before handing it out.
- Mark it done once its acceptance criteria hold and its work is committed.
- Mark it abandoned, with a reason, where building it turns out to be wrong, and in the same change every ticket it blocked takes its blockers in its place.

## Tickets

| Ticket | Blocked by | Status |
| --- | --- | --- |
| [1: The signature and the provenance are published under names verifiers recognise](1-names.md) | None | in progress |
