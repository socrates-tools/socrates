# socrates: one repository, one directory per tool, agora beside socrates-review

The suite is one repository and one binary, with each tool in a directory of its own.
Most of that has landed.
What is left is that agora still writes inside review's directory, the documents are not yet split per tool, and nothing ties the tool list to the namespaces the dispatcher mounts.

Spec: [../spec.md](../spec.md). Decisions: [ADR 42](../../../adr/0042-each-tool-reads-a-configuration-file-of-its-own.md).

## Working through the tickets

- Commit this folder before handing out any ticket.
- A ticket may start once every ticket it names as a blocker is done.
- Mark a ticket in progress before handing it out.
- Mark it done once its acceptance criteria hold and its work is committed.
- Mark it abandoned, with a reason, where building it turns out to be wrong, and in the same change every ticket it blocked takes its blockers in its place.

## Tickets

| Ticket | Blocked by | Status |
| --- | --- | --- |
| [1.1: agora writes under its own directory](1.1-agora-writes-under-its-own-directory.md) | None | done |
| [1.2: a glossary and a decision directory per tool](1.2-a-glossary-and-a-decision-directory-per-tool.md) | None | in progress |
| [1.3: the tool list agrees with the tree](1.3-the-tool-list-agrees-with-the-tree.md) | None | done |
