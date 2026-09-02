## Context

See `proposal.md` for why.

What was measured rather than reasoned about, with a real browser against the real binary:

- A page is opened, a decision committed, and `Send and finish` clicked.
  The caller then writes a second document for the same session.
  The open tab does not take it up: after eight seconds it still heads with the first document's title and still reports `data-revision` of 1.
- Reloading that same address shows the second document, with `data-over` of `no` and the send control enabled.
  So the server is already right and the page is the half that is wrong.
- Removing the client's `live.close()` and rebuilding does **not** fix it.
  The tab still fails to take up the second document.

That last measurement is what the design turns on, because the obvious one-line fix is the one that does not work.

The mechanism, read out of the source afterwards:

- `Exchange.Revise` sets `over` to false and bumps the revision, so a caller writing again reopens the exchange (`internal/surface/exchange.go:70-79`).
- `Server.live` writes `event: closed` and returns as soon as the exchange is over (`internal/surface/server.go:324-328`).
- The client's `closed` listener calls `live.close()`, which suppresses the browser's automatic reconnection, and `close()`, which disables every control (`web/src/surface.js:320-324`).
- `Server.live` takes `Changed()` before its loop and writes a revision only after that channel fires (`internal/surface/server.go:320,335-344`).
  It states nothing on connection.
  So a stream that connects after a revision it did not see is waiting for the next one, and where the caller writes once and stops there is no next one.

The two are independent, and each alone leaves the behaviour broken.
Keeping the stream open reconnects onto a server that will not say what it missed; stating the revision on connection helps nobody if the client has closed the stream for good.

## Goals / Non-Goals

**Goals:**

- The page the person is already looking at takes up a later document over an ended exchange, without a reload.
- A stream can be rejoined: what a client needs to be current is what it is told on connecting.
- The scenario that names this is proved without a reload, and reloading keeps a scenario of its own.

**Non-Goals:**

- Changing when an exchange is over, what `Send and finish` delivers, or what `ends` means to a caller.
  Finishing still closes the exchange, and the page still refuses to answer until a caller asks again.
- A second address for a later round. One address serves a whole review; that rule is what makes this reachable and is not what is wrong.
- Reconnection on any other footing: no `Last-Event-ID`, no replay of what a stream missed beyond the current revision, no history.
  The page holds one document at a time, so the current revision is the whole of what a client has to catch up to.
- Any change to what a caller writes, to the document contract, or to what is delivered.

## Decisions

**The stream states where the exchange stands on connection, and the client stops closing the stream.**
Both halves, because each was measured to be insufficient alone.
The server half is the one that generalises: a stream nobody can rejoin is a stream every client has to hold perfectly for the life of the page, which is not a property a browser offers.
Stating the current revision on connection is also what the neighbouring requirement already does for the closed case, so the two connect on the same footing rather than one being special.

**The closed state stops the controls and not the listening.**
`close()` today does both, and only the first is what being over means.
A page that is over is one that cannot be answered; it is not one that has stopped caring what the caller does next, because the caller asking again is the only thing that can make it answerable and the person has no other way to learn it happened.

**A revision arriving over a closed page reopens it by the path a re-render already takes.**
`take()` already replaces the content, rebinds, and reads `data-over` back off the new document, which is `no` after `Revise`.
So the reopening is not new machinery; it is machinery the closed page currently cannot reach because it dropped the stream.
Alternative rejected: a separate reopen path on the client, which would be a second way to do what `take()` does and would drift from it.

**The test for the scenario stops reloading, and reloading gets a scenario of its own.**
`internal/e2e/closing_test.go:185` reloads before asserting, so the stated requirement has been unproved since it was written while a test reported it green.
Deleting the reload would lose a real guarantee, so the reload becomes `A reloaded page still reopens` and the original scenario is proved on the open page.
This is the part of the change with the most value per line: the behaviour is already required, and what was missing was anything that would notice it was absent.

**The idle accounting is left alone.**
`Reading(1)` and `Reading(-1)` bracket a live stream, and a stream that stays open over a closed exchange keeps the session out of the idle count for as long as the tab is open.
That is what already happens for an exchange that never ended, and the collection rule exists for surfaces nobody is looking at rather than for surfaces nobody is answering.

## Risks / Trade-offs

**A page that reopens itself could take an answer the person thought they had finished giving.**
Mitigation: it can only reopen because the caller wrote a new document, which is the caller asking again, and the page shows that new document rather than the one they answered.
The existing rule that an uncommitted answer is not carried into content that no longer holds the decision is what keeps the previous answer out of it.

**A stream that never closes is a connection held for the life of a tab.**
Mitigation: it is one server-sent-events connection per open page on a loopback address, which is what an unfinished page already holds, and the surface stops on its own once no session has been read for the idle period.

**Stating the revision on connection makes every connect write a message.**
Mitigation: it is one line per connect, and the existing closed case already writes one.
The alternative, a client that cannot be caught up, was measured to be the defect.

**A test that stops reloading could become flaky, waiting on something that arrives over a network.**
Mitigation: it waits on a condition rather than a duration, as the neighbouring re-render test already does, and it runs against a loopback address in the same process tree.

## Migration Plan

1. The stream states the current revision on connection, with a test that opens one and asserts it is told without anything changing.
2. The client keeps the stream open across the end, and takes up a later revision.
3. The scenario is proved without a reload, and the reload becomes its own scenario.
4. The committed bundle is regenerated from `web/`.

Nothing persists across the change and nothing a caller wrote changes shape, so there is no rollback story to write.
A page open across the change is a page that reloads onto the new behaviour.
