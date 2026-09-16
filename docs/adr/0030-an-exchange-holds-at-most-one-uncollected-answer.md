# 30. An exchange holds at most one uncollected answer

- Status: accepted
- Date: 2026-09-16
- Part of: the page becomes a form

## Decision

An exchange held every answer sent to it in a queue, and collecting one took the front of that queue.
Two answers to the same document therefore left one of them queued, and the next caller to poll was handed it: once a later document had already been served, that caller was waiting on the new document, and the queued answer arrived as if it answered that one instead of the one it actually answered.

Receiving an answer now replaces whatever the exchange is still holding uncollected, rather than appending to a queue, so what a person last sent is what a caller collects.
Collecting an answer notifies on the same channel a revision and a close already notify on, so a live stream learns of a collection without polling for one.
What the exchange reports as standing carries a third fact alongside the revision and whether the exchange is over: whether it is holding an answer nobody has collected.
A stream connecting is told all three at once, matching the catch-up rule a revision already followed, rather than waiting for the next change before saying where the exchange stands.

The held answer records the revision it was made against, and collecting only hands one back while that revision is still the exchange's current one.
Serving a new document therefore leaves an answer left over from the old one sitting unreachable rather than erasing it: a caller waiting on the new document is never handed it, and the standing the exchange reports stops naming an answer uncollected the moment the document moves past it, with no second write inside the method that serves a new document.

## Rejected

- Appending to a queue, which is the design this replaces: it hands a caller the wrong answer once two arrive before either is collected.
- Refusing a second submit while the exchange already holds an uncollected answer, which would turn an amended answer into an error a person has to work around, rather than the exchange's ordinary case.
- Draining the held answer at the point a new document is served, which duplicates the rule a submit already enforces about replacing an answer, in a second place it could drift from the first.
