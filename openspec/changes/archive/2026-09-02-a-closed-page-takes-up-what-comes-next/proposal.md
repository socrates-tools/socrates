## Why

`surface-page` already requires that a page takes up new content without the person reloading it, and that a new document makes an ended page answerable again.
The second half does not hold.
A person who sends and finishes is left on a tab that says the exchange is over, and a document written for that surface afterwards never reaches them: the page shows the old content at the old revision until somebody thinks to reload an address that looks finished.

The test named for the scenario reloads the page before asserting, so it proves the reloaded page reopens rather than the stated requirement, which is the reload itself.
The requirement's own reasoning is what the test walks past: a person who has to reload to see the reply does not know the reply arrived.

This is reachable on the ordinary path.
One address serves a whole review, so the second round of a review the person answered on the page is written to the surface they are already looking at, and that is exactly the surface that stopped listening when they finished.

## What Changes

- A page whose exchange the person ended takes up a later document while it is open, becomes answerable again, and says so, with no reload.
  The server already reopens the exchange when a caller writes again; the page is the half that does not follow.
- The live stream states the revision the exchange is at when a stream connects, rather than only when the next change arrives.
  Without this a client that reconnects after missing a revision waits for a change that has already happened and never catches up, so no client-side fix can work alone.
- The page keeps its stream open across the end of an exchange, so that a later document can reach it.
  Today it closes the stream permanently on being told the exchange is over.
- The scenario `A new document reopens an ended page` is proved without a reload, and a reloaded page keeps a scenario of its own, so the two are not one claim again.

## Capabilities

### New Capabilities

None.
The behaviour is already required; what is missing is that it happens and that anything checks it.

### Modified Capabilities

- `surface-page`: the page takes up a later document over an ended exchange while open rather than on reload, and the live stream states the current revision on connection.

## Impact

- `web/src/surface.js`: what the client does when it is told the exchange is over, and what it does with a revision that arrives afterwards.
- `internal/surface/server.go`: what the live stream writes when it connects.
- `internal/assets`: the committed bundle, regenerated from `web/`.
- `internal/e2e/closing_test.go`: the test named for the scenario, which reloads today.
- No change to what a caller writes, to the document contract, to what is delivered, or to when an exchange is over.
