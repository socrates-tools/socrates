# A send goes through however many agora pages are open

Decisions: ADR 92, ADR 93.

## Problem Statement

With about six agora pages open in one Chrome or Edge, pressing *Send and wait* on any of them seems to do nothing.
The page says "Sending.", the buttons stay enabled, the status stays "Open", and the console is quiet.
`socrates agora poll` never returns, because the answer never reaches the surface.
Closing one agora tab lets the queued send through at once.

Every page holds a stream on `/sessions/<id>/live` as an `EventSource`, and each holds an HTTP/1.1 connection to the surface for as long as its tab is open.
Chromium allows six such connections per host and port, and every session shares the surface's one address.
Once six streams are open, the `POST /answer` waits in the browser's queue with no time limit.
The handoff that found and reproduced this is [handoff.md](handoff.md).

Two smaller faults make it worse.
A send that cannot get out looks like nothing happened: `send()` awaits a request with no time limit, leaves both presses enabled so a second press queues a second request, and never catches a rejected request.
And a page that loses the surface says nothing: after `socrates agora stop`, or once the surface collects the session, the page still says "Open" and offers to send into nothing.

## Solution

Each page holds its stream as one WebSocket, which the browser does not count against its six connections, per ADR 92.
The stream says what it said before, `revision`, `closed` and `uncollected`, in the same order and only on change, and counts as a reader for exactly as long as it is open.
The surface opens a stream only for an upgrade whose `Origin` is the surface's own, and the stream carries no token, per ADR 93.
The surface pings each stream and drops one that does not answer, so a page that vanished without closing stops keeping its session alive.

When the page loses the surface, it says so and keeps trying with a backoff.
When it reconnects, it is told where the exchange stands and catches up.
When the surface answers that it holds no such session, the page closes the exchange, says the session is gone and stops trying.

A send is visible however it ends.
Both presses are disabled while a send is in flight.
A send that has not been answered within one minute is aborted, and the page says nothing was sent and why it might be.
A request the browser rejects says the same.
Either way the presses are enabled again and nothing the person wrote is lost.

## User Stories

1. As a person answering several agora pages at once, I want *Send and wait* to reach the surface however many pages I have open, so that I never have to close a tab to get an answer through.
2. As a person who keeps old agora pages open in tabs I have not closed, I want those pages to cost the page I am answering nothing, so that a forgotten tab never stops a send.
3. As a person who pressed *Send and wait*, I want both presses disabled while the answer is going, so that a second press does not queue a second request behind the first.
4. As a person whose send cannot get out, I want the page to say within a minute that nothing was sent and why it might be, so that I am never left looking at "Sending." wondering whether it went.
5. As a person whose send failed, I want my choices, words and attachments still on the page and the presses enabled again, so that I can send again without writing anything twice.
6. As a person whose browser rejects the request outright, I want the page to say nothing was sent, so that a failure is never silent.
7. As a person whose surface was stopped while I was reading, I want the page to say it lost the surface and is trying again, so that I know an answer sent now would go nowhere.
8. As a person whose surface comes back, I want the page to reconnect by itself and show where the exchange stands, including a revision I missed while it was gone, so that I do not have to reload.
9. As a person whose session the surface has collected or never held, I want the page to say the session is gone and stop offering to answer, so that I do not write an answer nobody will read.
10. As a person who opens a page through a forward named by `surface.reached_at`, I want a page I closed by putting my machine to sleep to stop keeping its session alive, so that the surface still collects the session and stops once nobody is reading.
11. As a person with any other website open in the same browser, I want that website unable to open a stream on my session, so that nothing outside agora learns when my exchange changes.
12. As a caller running `socrates agora serve` and `socrates agora poll`, I want nothing about serving, revising, ending or polling to change, so that the skills that drive agora keep working as written.
13. As a caller who revises a page the person already has open, I want the revision to reach that page as promptly as it does today, so that the person sees the next round without reloading.
14. As a caller who ends a session, I want the page to say the exchange is over as it does today, so that the person knows nothing further will be read.
15. As a maintainer, I want the stream's rules held at the surface by tests over a real WebSocket, and the person's experience held by the browser suite with seven pages open, so that the hang cannot come back unnoticed.

## Implementation Decisions

### The stream

- `/sessions/{session}/live` stays the route and becomes a WebSocket upgrade, built with `github.com/coder/websocket`, a new module that is linked into the one binary and so gets its licence into `NOTICE`.
- No `text/event-stream` route remains, per ADR 92.
- `Server.live` in `agora/surface/server.go` accepts the upgrade, calls `exchange.Reading(1, ...)` once it is accepted and `Reading(-1, ...)` when the socket closes for any reason, and otherwise keeps its loop: it writes where the exchange stands on connecting, then each change once, until the request ends, the surface stops or the socket closes.
- A session the surface does not hold is refused before the upgrade with the same not-found answer it gives today, so the page can tell a session that is gone from a surface that is gone.
- Each message is one text frame naming the event and its data, so the page reads `revision`, `closed` and `uncollected` as it does today.
  `whereItStands` and `whatIsUncollected` keep deciding what is said and when, and change only in how they write it.
- The surface reads from the socket only to answer the protocol's own control frames, and treats anything a page sends as nothing.

### The guard

- The upgrade is refused unless `Origin` names the same host and port as the request, which is `coder/websocket`'s default check on `Accept`, so no origin is listed by hand, per ADR 93.
- The host check that already turns away a host the surface did not bind runs before it, as for every route.
- The stream carries no token.
- The page's policy keeps `connect-src 'self'`, which admits a WebSocket to the page's own host and port, and the browser suite records no violation of it.

### The ping

- The surface pings each stream about every thirty seconds and closes one whose pong does not arrive in time, which takes it off the session's readers.
- The interval and the wait are stated once in `agora/surface` beside the other bounds on waiting.

### The page

- `follow()` in `agora/web/packages/surface/revision.ts` opens a `WebSocket` on `${base()}/live` with the scheme of the page, and dispatches each message to what the three `EventSource` listeners do today.
- When the socket closes or never opens, the status says the page lost the surface and is trying again, and the page reconnects after a backoff that starts near a second and doubles to a cap near thirty seconds.
- On reconnecting, the status returns to where the exchange stands, and the first message brings any revision the page missed, as the stream already promises.
- When the surface answers that it holds no such session, the page closes the exchange through `closeTheExchange()`, says the session is gone and stops reconnecting.
- Every new sentence the page says lives in the surface's messages beside `exchangeIsOver`, carries no em dash and says what the person can do.

### The send

- `send()` in `agora/web/packages/surface/sending.ts` disables both presses before the request and enables them again on every way it ends but success.
- The request carries an `AbortSignal` that aborts it after one minute, stated once in the surface's script.
- A timeout and a rejected request each say plainly that nothing was sent and why it might be, with the `error` tone the page already uses for a refused answer.
- Nothing the person chose, wrote or attached is cleared by a failed send.

### Documents

- `agora/AGENTS.md` names the stream as a WebSocket where it describes the live revision, and gains invariants for the origin check, the reader count over the socket's life, the ping, the page that says it lost the surface, and the send that is never silent.
- `agora/CONTEXT.md` already defines **stream**, written during this design.
- `make assets` rebuilds the committed bundle after the change under `agora/web/`.
- This folder's `handoff.md` goes with the folder when the change lands.

## Testing Decisions

The change is tested at two seams that already exist: the browser suite in `internal/e2e`, and the surface's own tests in `agora/surface/stream_test.go`.

### The browser suite

`internal/e2e` drives the real binary and a real browser, as `closing_test.go` already does for the stream.

- Seven pages open on one surface, any sessions, and a send from the last is collected by `socrates agora poll`.
  This is the scenario the hang was found in, and it fails on today's build.
- After `socrates agora stop`, an open page says it lost the surface and is trying again.
- A page whose surface comes back with the session still held reconnects and shows a revision served while it was gone.
- A page whose session the surface no longer holds says the session is gone and offers nothing to answer.
- A send that runs out of time says nothing was sent and enables both presses again, with what the person wrote still there.
  The suite reaches the timeout by holding the request at a surface that never answers rather than by waiting a real minute, through whatever the suite's harness already uses to stand in for a surface.
- A second press while a send is in flight sends nothing more.
- The existing tests that a revision served to an open page swaps in, that `closed` closes it and that `uncollected` flips after a poll still pass unchanged.

### The surface's own tests

`agora/surface/stream_test.go` is rewritten over a real WebSocket against a real `httptest` server, keeping each of its five scenarios by name:

- a stream is told where the exchange stands when it connects;
- a stream that missed a revision catches up on connecting again;
- a stream opened after the end is told at once;
- an open stream reports each change once and then the close;
- a stream connecting while an answer is uncollected is told at once.

It gains:

- an upgrade from another origin is refused, and one from the surface's own is accepted;
- a stream counts as a reader for exactly as long as it is open, read through `IdleSince`;
- a stream that does not answer a ping is closed and stops counting as a reader;
- an upgrade naming a session the surface does not hold is refused with not found.

Each guard is mutated once on purpose, as the conventions ask: removing the origin check, the reader decrement and the ping each fails a named test.

## Out of Scope

- Any change to `serve`, `poll`, `end`, `stop` or the routes they reach.
- A token on the stream, per ADR 93.
- A route for an `EventSource` page from before an upgrade, per ADR 92.
- Closing or pausing a stream while its tab is hidden.
- Sharing one stream between pages through a worker.
- A limit on how many pages may be open.

## Further Notes

- The one-minute send limit was chosen as stated, and the ping of about thirty seconds was chosen as an approximate interval.
  The backoff from about a second to about thirty seconds is this spec's, since the design settled only that the page retries with one.
- The design asked nine questions over two rounds, and the person chose the recommended option on each without adding words.
