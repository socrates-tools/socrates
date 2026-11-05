# A send goes through however many agora pages are open

With six or more agora pages open in one browser, a send waits forever because every page's stream holds one of the browser's six connections to the surface.
This change carries the stream over a WebSocket, makes a send that cannot get out say so, has the page say when it loses the surface, and drops a stream whose page has gone silent.

Spec: [../spec.md](../spec.md). Decisions: [ADR 92](../../../../agora/docs/adr/0092-a-page-hears-its-exchange-over-a-websocket.md), [ADR 93](../../../../agora/docs/adr/0093-a-stream-is-guarded-by-its-origin-and-carries-no-token.md).

## Working through the tickets

- Commit this folder before handing out any ticket.
- A ticket may start once every ticket it names as a blocker is done.
- Mark a ticket in progress before handing it out.
- Mark it done once its acceptance criteria hold and its work is committed.
- Mark it abandoned, with a reason, where building it turns out to be wrong, and in the same change every ticket it blocked takes its blockers in its place.

## Tickets

| Ticket | Blocked by | Status |
| --- | --- | --- |
| [1.1: The stream over a WebSocket](1.1-the-stream-over-a-websocket.md) | None | done |
| [1.2: A send that is never silent](1.2-a-send-that-is-never-silent.md) | None | done |
| [2.1: The page says it lost the surface](2.1-the-page-says-it-lost-the-surface.md) | 1.1 | pending |
| [2.2: A silent stream is dropped](2.2-a-silent-stream-is-dropped.md) | 1.1 | done |
