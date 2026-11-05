# 93. A stream is guarded by its origin and carries no token

- Status: accepted
- Date: 2026-11-05
- Part of: agora-live-over-websocket

## Decision

The surface opens a stream only on an upgrade whose `Origin` names the same host and port the request was sent to, and refuses every other.
The stream carries no token, unlike every route that changes state.

A WebSocket is not covered by CORS, so without a check any website open in the same browser could open a stream on a session it can name.
The browser always sends `Origin` on a WebSocket and a page cannot forge it, so checking it keeps every other website out.
A token would keep out nothing more: the page's own `GET` embeds the session's token without asking for one, so anything that can read the page, which is any local process, can read the token too.
A WebSocket also cannot set the `x-agora-token` header, so a token would ride in a subprotocol or the URL, where the browser logs it, for a stream that changes nothing.

## Rejected

- The `Origin` check and the session's token in a subprotocol, so the stream is guarded like the routes that change state: a second mechanism, and a token in a header the browser logs, that keeps out nothing the `Origin` check does not.
