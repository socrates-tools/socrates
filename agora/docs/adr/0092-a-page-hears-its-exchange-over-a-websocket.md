# 92. A page hears its exchange over a WebSocket

- Status: accepted
- Date: 2026-11-05
- Part of: agora-live-over-websocket

## Decision

Each page holds its stream on `/sessions/<id>/live` as one WebSocket, and the surface serves no `text/event-stream` route.

Every session shares the surface's one address, and Chromium allows six HTTP/1.1 connections per host and port.
An `EventSource` holds one of them for as long as its tab is open, so with six pages open the `POST /answer` of any of them waits in the browser's queue with no time limit and never reaches the surface.
A WebSocket does not count against that pool, so any number of pages stay live and a send always gets a connection.

The stream says what it said before, `revision`, `closed` and `uncollected`, in the same order and only on change, and is told where the exchange stands when it connects.
It counts as a reader of its session for exactly as long as the socket is open.

No route is kept for an `EventSource` page from before an upgrade, because a page only ever talks to the surface of the build that served it: per ADR 58, a surface of another build is replaced only once it holds no session, so the page it served has nothing left to hear.

A WebSocket does not reconnect on its own, so the page reconnects with a backoff and says while it is trying.

## Rejected

- One shared stream per browser, held by a `SharedWorker` for every agora page of the origin: it keeps SSE, but `/live` becomes one stream across many sessions, the reader count has to be fed per session through the worker, and a browser without `SharedWorker` falls back to the limit.
- Short polling of `/live` every second or two: it holds no connection, but a revision appears up to a poll late, the surface answers a request per open page every interval, and being a reader becomes having asked recently rather than being connected.
- Closing the stream while a tab is hidden: it lowers the count without removing the limit, so six visible windows still hang a send, and a hidden page stops counting as a reader.
