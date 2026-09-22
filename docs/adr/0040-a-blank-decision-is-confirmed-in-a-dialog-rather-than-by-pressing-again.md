# 40. A blank decision is confirmed in a dialog rather than by pressing again

- Status: accepted
- Date: 2026-09-23
- Part of: blank-confirmation

## Decision

A press over a decision left blank used to arm rather than send.
The first press opened an inline panel at the foot of the page naming the blank decisions, and the client recorded which of the two buttons had been pressed.
A second press of that same button sent.
Pressing the other button, or answering a card and pressing again, silently disarmed and started over.

A press now raises a modal dialog that asks.
The dialog names every decision left blank, states that a blank decision goes as nothing, and offers *Send anyway* and *Go back and answer*.
Nothing about the button changes meaning between presses: a press either sends, or raises a dialog whose own press sends.

The trigger does not change.
A page with nothing blank sends on one press, on either button, and somebody who agrees with every recommendation on a document meets no dialog at all.
Finishing is not confirmed on its own account, even though finishing closes the exchange, because the dialog is about the gaps in an answer rather than about the press that carries it.

Cancelling leaves nothing behind.
*Go back and answer*, Escape and a press on the backdrop are one path: the dialog closes, the page is as it was, and the cursor returns to the button that raised it.
The blank cards are found through the dialog's own names, each of which closes it and takes the person to that card, rather than through anything left marked on the page.

This supersedes the brake ADR 31 describes.
That decision let a document's recommendations arrive chosen, and named the notice as one of the two things standing between a person and a page sent without being read.
The brake is unchanged in effect and changed in kind: a card the caller left to the person is still blank, and still has to be confirmed or answered before it goes as nothing, but confirming it is now answering a question rather than repeating a press.

What the dialog costs is a press for every person who deliberately leaves a card blank, and one more piece of client behaviour than a panel needs, since a modal dialog has an open state, a cancel path and a cursor to put back.
What it buys is that the mechanism no longer has to be taught by a line of body text at the bottom of the page.

## Rejected

- Raising the dialog wherever a decision is blank and over every *Send and finish*, blank or not, on the reasoning that finishing closes the exchange and cannot be amended.
  This confirms a fully answered page, which is the one case where the person has demonstrably read it, and it puts a dialog in front of the ordinary end of every design and every review round.
- Raising the dialog over every press, blank or not, so that the page behaves the same way every time.
  Uniform, and it costs one extra press on every round for everybody, including the person who agreed with a whole document and expected the one press ADR 31 exists to give them.
