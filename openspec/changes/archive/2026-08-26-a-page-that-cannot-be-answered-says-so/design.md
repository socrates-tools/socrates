# A page that cannot be answered says so

## The failure, stated once

`over` is authoritative for writes and invisible to every read.

`Exchange.Receive` refuses once `over` is set, and returns the error the server turns into a 409.
`Exchange.Over` is consulted in exactly one place, the `live` stream, and only after that stream has already waited for a change.
`Server.page` never consults it at all.
So the one state that decides whether an answer can be accepted is absent from the page that invites the answer, and a person can only discover it by losing something.

Everything below follows from closing that gap in the three places it is open.

## Decision: an uncommitted answer is named before it is discarded, not delivered

`send` builds its decisions from the `committed` map alone.
A selection the person made sits in the DOM, where `chosenIn` would find it, and is dropped without a word.

The tempting fix is to harvest it.
It is wrong, and the spec says why: *Choosing SHALL change only what the page holds, and SHALL send nothing*, because *a control that sent on every click sends an answer the person was still deciding, and takes back nothing when they change their mind*.
Harvesting at finish time is that control with a delay.
A person who selected an option to read it against the others, then decided to end without answering, would have their reading delivered as their answer.

So the selection stays undelivered and the silence goes.
The finish states which decisions hold something uncommitted and what will not reach the caller, and waits for one confirmation.
This keeps both properties that matter: nothing is delivered that was not committed, and nothing is discarded that the person was not told about.

Ending without answering stays possible, because a person who is not deciding is a real answer to a caller and `poll` already reports it.
It costs one further click, which is the price of the click that lost an answer.

## Decision: `Revise` clears `over`

Two candidates were considered for the re-serve path.

**`serve` refuses onto an ended exchange.**
Terminal means terminal, and the caller runs `agora stop` before serving again.
Rejected: it makes every caller learn a stop-then-serve dance to recover from a mistake the person made, and the failure it replaces a 409 with is still a failure at the moment somebody is trying to fix something.

**`Revise` clears `over`.**
A caller writing a new document is asking again.
That is what the re-render path already means when the exchange is open, and there is no reading of it that means something else when the exchange is closed.
Chosen, because it is the smaller change and it only alters behaviour in the case that is broken: while an exchange is open, `over` is already false and clearing it is a no-op.

This does not disturb `Delivery happens once` or `Ending is distinguishable`.
Both are about what a poll gets: an answer is handed to one poll and not a second, and a poll never reopens an exchange.
Reopening is something only a caller writing a new document does, and the spec delta states it beside the rule it could be confused with.

A page that closed itself when the person ended stays closed until it is reloaded or the stream tells it otherwise.
That is acceptable and is why the revision is a stream event: the page it reaches takes up the new document, and the closed state it takes up with it is the one the exchange is now in.

## Decision: the closed state is rendered, not only streamed

The stream already emits `closed`, so the page has a closed state and knows how to enter it.
What is missing is that a page arriving after the fact is never told.
Two edits, both narrow:

- `page` reads `Over` and passes it into the frame, and the template renders the closed state directly.
  The client's `close` already does the rest, so a served-closed page starts closed rather than waiting for an event.
- `live` checks `Over` before it waits, rather than after.
  As it stands the loop takes `Changed`, blocks on it, and only then asks whether the exchange is over, so a stream opened onto an already ended exchange waits for a change that will never come.

The second is worth fixing even though the first makes it invisible on a reload, because a stream that reports a terminal state only on the way past it is wrong on its own terms.

## What is deliberately not changed

- **The commit step.** It is the reason an answer can be taken back, and this change makes it safer rather than optional.
- **`over` as a real state.** It still refuses writes, `poll` still reports it, and the page still says the exchange is over. It stops being permanent only against a caller writing a new document.
- **The command surface.** `serve`, `poll`, `end`, `stop`, `version`, and nothing added. Recovering from an ended exchange is serving a document, which is a thing `serve` already does.
