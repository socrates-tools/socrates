# 41. The confirmation opens on Send anyway

- Status: accepted
- Date: 2026-09-23
- Part of: blank-confirmation

## Decision

The dialog that confirms a send over a decision left blank opens on *Send anyway*.
That press carries the cursor when the dialog opens, so Enter takes it, and it is drawn as the primary of the two.
*Go back and answer* is the quiet one beside it.

The word is *Send anyway* whichever button raised the dialog.
The dialog names neither *Send and wait for a reply* nor *Send and finish*, because the person knows which of the two they pressed: pressing it is what put the dialog in front of them.
*Send anyway* means do what I told you anyway, which is true of both, so the dialog carries one word for the send and the finish alike and never has to say which it is standing in front of.

The usual shape of a confirmation is the other way round.
A dialog standing in front of something consequential conventionally rests on the safe press, so that a person who answers it without reading it is returned to where they were rather than carried forward.

The person here already pressed send.
The dialog is not asking whether they meant to send, it is telling them what their send will leave out, and the answer to it is usually yes.
Resting on the escape makes the common path a press and an aim, and makes deliberately leaving a card blank the expensive way through the page, which is the cost the arming mechanism already charged and this change exists to remove.

What that accepts is that a person holding Enter down, or pressing it twice in quick succession, passes both the send and its confirmation without reading either.
The brakes against that are the ones ADR 31 already names and this decision does not add to: the bar sits at the end of the page in the normal flow, so reaching it means having passed every card, and the dialog itself is modal, so it takes the cursor and the page behind it goes inert.

## Rejected

- *Go back and answer* carrying the cursor, with *Send anyway* the deliberate press every time.
  The cautious default, and what most confirmations over something irreversible do.
  It charges a press and an aim to everybody who left a card blank on purpose, which is the cost this change set out to remove rather than relocate.
- *Send anyway* drawn as the primary while the cursor rests on *Go back and answer*, so the eye is led to sending and the keyboard is not.
  This splits the two signals a dialog has for saying which press it expects, and reads as the dialog disagreeing with itself.
