# Handoff: Send hangs when six or more agora pages are open

Status: found and reproduced, not fixed.
Read against `main` at `5e54fce`.

## What a person sees

With about six agora pages open in one Chrome or Edge, pressing *Send and wait* on any of them seems to do nothing.
The page writes a small "Sending." in the side panel, the buttons stay enabled, the status stays "Open", and the console is quiet.
`socrates agora poll` never returns, because the answer never reaches the surface.
Closing one agora tab lets the queued send go through at once.

## Why

Every page opens an `EventSource` on `/sessions/<id>/live` (`agora/web/packages/surface/revision.ts:38`), served by `Server.live` in `agora/surface/server.go`.
Each one holds an HTTP/1.1 connection for as long as its tab is open.
Chromium allows six HTTP/1.1 connections per host and port, and every session shares the surface's one address, `127.0.0.1:24672`.
Once six streams are open, the `POST /answer` in `sending.ts` waits in the browser's queue with no time limit and never reaches the server.
Loopback is plain HTTP/1.1, so there is no HTTP/2 to multiplex the streams.

Evidence from the session that hit it: Edge held nine open connections to `127.0.0.1:24672`, and the session folder had no answer for the round the person sent.

## Reproduce

1. Serve any document and open its page in Chrome or Edge.
2. Open five more agora pages on the same surface, any sessions, or open five `EventSource`s to the same host from the page's console.
3. Press *Send and wait*: it stays on "Sending." and the surface records no answer.
4. Close one page or stream: the send goes through and the status turns to "Sent".

## The fix

### 1. Carry `/live` over a WebSocket

This is the fix that removes the cause.
WebSockets do not count against the browser's six-connection HTTP/1.1 pool, so any number of pages can stay live and a send still gets a connection.
lavish-axi met the same problem and made this move, and keeps the old SSE route only for pages from before an upgrade: see `src/chrome-client.js` (`connectLiveEvents`) and `src/server.js` (`handleEventUpgrade`) in <https://github.com/kunchenguid/lavish-axi>.

- Keep what the stream says: the `revision`, `closed` and `uncollected` events, in the same order and only on change, as `whereItStands` and `whatIsUncollected` do today.
- Keep the reader count: `Server.live` calls `exchange.Reading(1, ...)` and `Reading(-1, ...)`, and a session is collected after 45 minutes with no reader, so a WebSocket must count as a reader for exactly as long as it is open.
- Check the `Origin` on the upgrade: a WebSocket is not covered by CORS, so the surface must refuse any origin but its own page, or a website open in the same browser could connect.
  `/live` carries no token today; decide whether the socket should also carry the session's token, for example in the URL or a subprotocol, since a WebSocket cannot set the `x-agora-token` header.
- Add a heartbeat and reconnect with backoff, so a page that sleeps or loses the server says so rather than going quiet.
- The Go side has no WebSocket package yet; `github.com/coder/websocket` is the usual small one. Only `internal/e2e/illustration_test.go` opens a WebSocket today.
- Decide whether to keep SSE for one release for pages still open across an upgrade, as lavish does, or to drop it.

### 2. Make a stalled send visible

Even with the cause gone, a send that cannot get out should never look like nothing happened.
In `send()` in `agora/web/packages/surface/sending.ts`:

- disable both presses while a send is in flight, so a second press does not queue a second request;
- give the `fetch` a timeout, and on it say plainly that nothing was sent and why it might be;
- catch a rejected `fetch`, which nothing handles today, and say so on the page.

### Not chosen

Closing the stream while a tab is hidden and reopening it when shown lowers the count but does not remove the limit: six visible windows would still hang a send.

## Where it touches

- `agora/surface/server.go`: the `/live` route and `Server.live`.
- `agora/web/packages/surface/revision.ts`: `follow()`.
- `agora/web/packages/surface/sending.ts`: `send()`.
- `agora/AGENTS.md` and `agora/CONTEXT.md` where they describe the live revision.
- A decision in `agora/docs/adr/` for the transport, since it changes what every page holds open.
- `make assets` rebuilds the committed bundle after a change under `agora/web/`.

## How to check it

- An end-to-end test in `internal/e2e` that opens seven pages on one surface in the real browser, sends from the last, and sees the answer collected.
- The existing live-revision tests still pass: a revision served to an open page swaps in, `closed` closes it, and `uncollected` flips after a poll.
- `make fmt` and `make verify`, as `AGENTS.md` says.

This change starts as a spec under `docs/changes/<slug>/`, per `CONTRIBUTING.md`; this file is only the handoff, and goes once that spec exists.
