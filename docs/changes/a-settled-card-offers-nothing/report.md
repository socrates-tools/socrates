# A settled card still offers a words field

Raised on 2026-09-30, in the message of the second round of the `socrates-tickets` design, by the person being designed with:

> For settled answers coming from previous round, you only need to show them as readonly visually, there's nothing to talk about them anymore so no textbox needs to follow them at all - nothing to answer for them!

## What the person sees

Round two of a design opens with every question round one settled, as cards that state the answer and offer no options.
Each of those cards still renders a words field under it.
In the round that raised this, nine settled cards each carried an empty text box above the six questions that were actually open, so the page asked for nothing nine times before it asked for anything.

## Why it is there today

It is stated, not accidental, in three places:

- The surface invariant in `AGENTS.md`: "Every decision carries a free-text field that no document field can remove, including one that offers no options."
- `socrates-design`, in `design/skill`: "The card still carries the field a person writes in, for whoever wants to reopen it."
- The `socrates-design` spec's user story 6: "As somebody being designed with, I want to reopen something I settled two rounds ago, so that a design is not a ratchet."

So the field on a settled card is how a design stays reopenable, and a document cannot remove it, because the surface's invariant says no document field can.

## What a fix has to decide

The request collides with that invariant and with the reopening story, so a fix is a change with a spec of its own rather than a patch:

- Whether a settled card renders read-only with no field, and reopening moves somewhere else: the page's own message, which every page already carries, or a control on the card that opens the field only when pressed.
- Whether the rule belongs to agora or to the design skill.
  agora's invariant is about any decision offering no options, which review uses for a finding a session already resolved, so dropping the field there changes review's page too.
  A narrower change states a way for a document to mark a card as read-only, which the invariant's "no document field can remove" forbids today.
- Whether the settled block needs cards at all, given the preamble already carries one sentence per settled question.

## Status

Fixed by `docs/changes/the-kylix-page/spec.md`.
A settled answer is carried by the document's `settled` list rather than as a card, offers no field, and is reopened through the page's message, per ADR 46.
