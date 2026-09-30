# 46. What was settled is carried by the document rather than as a card

- Status: accepted
- Date: 2026-09-30
- Part of: the-kylix-page

## Decision

A document may carry a `settled` list beside its decisions.
Each entry states an identifier, a name and a location where it has them, the question, the labels the person chose and their own words.
The page reads the list ahead of every open question, as one quiet block that offers no control, is counted by nothing and is never sent back.

A settled answer used to be a decision with no options, whose meaning restated the answer.
Such a card still carried the words field every decision carries, so a page opening a later design round asked for nothing once per settled question before it asked for anything.
A settled answer now offers no field of its own, and a person reopening one says so in the message the page carries for the whole round.

An identifier is either settled or still open, and a document naming one in both is refused before anything is served.

The block says what it is in agora's own words: *Settled earlier* above the list and *Still open* above the cards.
The surface counts no rounds and names none, because a round is a word of the tools that use it rather than of the page.

## Rejected

- Keeping settled answers as cards with no options.
  The words field on each is how the page invited reopening, and the invariant that no document field can remove it is right for a decision: a decision is something to answer.
  A settled answer is not one, so it gets a shape of its own rather than an exception to that rule.
- A document field marking a card read-only.
  It would be the one field that removes the words field, and every tool writing a document would have to know when to set it.
- A `round` on the document and on each settled answer, with the page showing a round badge and which round settled what.
  It gives the surface a noun it does not have, and a design and a review each count rounds for reasons of their own.
