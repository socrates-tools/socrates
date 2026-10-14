# One command installs socrates and every skill

Each release carries `install.sh`, which picks the archive for the machine's operating system and architecture, checks it, places `socrates` and runs `socrates skill`.
`socrates skill` installs every skill every namespace offers, and the README's manual route downloads one archive.

Spec: [../spec.md](../spec.md). Decisions: [ADR 61](../../../adr/0061-the-top-level-installs-every-namespaces-skills.md).

## Working through the tickets

- Commit this folder before handing out any ticket.
- A ticket may start once every ticket it names as a blocker is done.
- Mark a ticket in progress before handing it out.
- Mark it done once its acceptance criteria hold and its work is committed.
- Mark it abandoned, with a reason, where building it turns out to be wrong, and in the same change every ticket it blocked takes its blockers in its place.

## Tickets

| Ticket | Blocked by | Status |
| --- | --- | --- |
| [1: socrates skill installs every namespace's skills](1-top-level-skill.md) | None | done |
| [2: The release carries an installer that picks the right archive](2-install-script.md) | 1 | done |
| [3: The README installs in one line, and its manual route downloads one archive](3-readme.md) | 2 | done |
