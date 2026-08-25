## Why

One real handover lost a person's answer twice, and then hid that it had.
The sequence, in the order it happened:

- A page held one decision.
  The person clicked `Send and finish` before `Commit this answer`.
  The exchange ended carrying nothing: `poll` printed `{"decisions": null, "ended": true}`, and the answer already selected on screen was neither delivered nor mentioned.
- They reloaded the page.
  It rendered exactly as it had the first time, with every control live.
  They chose again, committed, and sent.
  The send was refused, and the page said `That was not accepted: 409. Nothing was sent.`
- The caller re-served the same document to ask again.
  `serve` printed the same address it had before, and that address could only refuse.
  Nothing said so, at the page or at the caller.

These are three views of one failure: what the page shows and what the exchange will accept are allowed to disagree, and the person is the one who finds out.
The first loses an answer that was on screen.
The second invites an answer that cannot be accepted.
The third sends the person back to a page that was never going to work, which is exactly what a caller does after the first two.

Choosing sends nothing is deliberate and stays.
A control that sent on every click takes back nothing when somebody changes their mind, which is why the commit step exists.
What is wrong is not that an uncommitted selection is undelivered, it is that finishing discards one in silence.

## What Changes

- **A finish that would leave an answer undelivered asks first.**
  Where a decision holds a selection or written words that were never committed, `Send and finish` states which decisions those are and what will not reach the caller, and proceeds only when the person confirms.
  An uncommitted selection is still never delivered.
  A person who means to end without deciding still can, in one further click.
- **A page that cannot be answered says so, however it was loaded.**
  The closed state stops being something only a page that was open at the time is told about.
  A page loaded fresh onto an ended exchange renders closed, and a live stream opened onto one is told immediately rather than after the next change that never comes.
- **A caller who writes a new document is asking again, so the exchange is no longer over.**
  Writing a new document to an ended surface reopens it and the page becomes answerable.
  A poll still never reopens an exchange, and delivery still happens once.
  `serve` therefore stops handing back an address whose only possible reply is a refusal.

## Capabilities

### Modified Capabilities

- `surface-page`: finishing asks before it discards an uncommitted answer; a page renders its closed state however it was loaded, and takes up a new document by becoming answerable again.
- `review-surface`: a new document reopens an ended exchange, stated beside the rule that a poll does not.

## Impact

- `internal/surface/exchange.go`: `Revise` clears `over`; `Over` becomes something the page renderer reads rather than only the stream.
- `internal/surface/server.go`: `page` consults `Over` and passes it into the frame; `live` reports a closed exchange on connection rather than waiting on the next change first.
- `internal/render`: the frame carries whether the exchange is over, and the template renders the closed state.
- `web/src/surface.js`: the finish path asks before discarding, and a page served closed starts closed rather than waiting for a stream event.
- `internal/e2e`: the browser journeys for finishing with an uncommitted answer, reloading onto an ended exchange, and re-serving to ask again.
- **No change to** `internal/document`, `internal/sanitise`, `cmd/agora`'s command surface, or anything socrates does.
  `serve` already revises an open surface and prints its address; what changes is that the exchange it revises is answerable.

## Dependencies

None.
`2026-08-25-the-page-speaks-to-the-person` is the most recent change to both capabilities and is unaffected: nothing here renders a command, names a decision, or touches the palette.
