# 59. A settled answer is checked against what the session collected

- Status: accepted
- Date: 2026-10-14
- Part of: the exchange carries what a person said
- Amended by: ADR 87, whose attachments are held beside the question, the labels and the words: a settled answer states the names of the files collected for its identifier, in the order collected, or is refused

## Decision

A caller writes what was settled afresh on every revision, from what `agora poll` handed back, because nothing on the page carries across one.
The person's own words in a settled entry were therefore written by the caller, and a caller that restated them rather than copying them changed what the person said with nothing to notice.

An exchange now remembers, per identifier, the last answer a caller collected: the question it answered, the labels chosen and the words.
A revision whose `settled` list states an identifier the exchange collected an answer for must state that question, those labels and those words exactly, and one that does not is refused before the document is written or served, naming the identifier, the field and what was collected.
The page the person holds stays where it was.

Only a collected answer binds.
An answer sent but never collected was never handed to the caller, so the caller cannot be held to it, and an identifier the exchange holds no answer for is served as the caller states it.
A question asked again and answered again binds what was answered last.

What carries across a revision on the page is unchanged: no selection, typed word, message, receipt or confirmation.
What the exchange remembers is what it already handed to the caller, and it is never sent back to the page.

## Rejected

- The exchange filling in settled answers itself, from identifiers alone.
  It would change the document every tool writes, and a caller that settled something the exchange never saw, such as a decision taken before the session opened, would have no way to state it.
- Comparing words loosely, ignoring whitespace or case.
  The words become a spec and its decisions exactly as written, and a rule that admits a small change has to say how small, which nobody can.
- Refusing every settled entry the exchange holds no answer for.
  A first document has collected nothing, so it could carry nothing settled, and the case this guards against is a restated answer rather than an invented one.
