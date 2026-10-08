# agora binds one port, and asks it who stands there

agora stops walking a window of ten ports and binds one, 24672, refusing it when something else holds it.
`serve` finds the surface already standing by asking that port rather than by trusting `surface.json`, refuses a surface of another root, and replaces a surface of another build once it holds no session.

Spec: [../spec.md](../spec.md). Decisions: [ADR 57](../../../../agora/docs/adr/0057-the-surface-binds-one-port-and-asks-it-who-stands-there.md), [ADR 58](../../../../agora/docs/adr/0058-a-surface-of-another-build-is-replaced-once-it-holds-no-session.md).

## Working through the tickets

- Commit this folder before handing out any ticket.
- A ticket may start once every ticket it names as a blocker is done.
- Mark a ticket in progress before handing it out.
- Mark it done once its acceptance criteria hold and its work is committed.
- Mark it abandoned, with a reason, where building it turns out to be wrong, and in the same change every ticket it blocked takes its blockers in its place.

## Tickets

| Ticket | Blocked by | Status |
| --- | --- | --- |
| [1: The surface binds one port, and a held port is refused](1-one-default-port.md) | None | done |
| [2: serve finds the surface by asking the port who stands there](2-surface-found-by-asking-the-port.md) | 1 | done |
| [3: A surface of another build is replaced once it holds no session](3-another-build-replaced-once-idle.md) | 2 | in progress |
