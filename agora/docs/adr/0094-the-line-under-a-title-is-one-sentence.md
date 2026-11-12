# 94. The line under a title is one sentence

- Status: accepted
- Date: 2026-11-12
- Amends: ADR 36

## Decision

A document whose preamble is longer than 200 characters is refused before anything is served.
The refusal names the length and the limit, and says to write one sentence on what the page asks.

The design skill no longer carries a sentence per settled question into the preamble.
It writes one sentence on what the round asks the person to decide, or nothing.
A fact found since the last round goes in the meaning of the decision it bears on.
The tickets skill asks its two questions of the breakdown in one sentence, and names a user story no ticket delivers on a card of its own.
The review skills hold the page to the same line: one sentence of at most 200 characters, or nothing.

ADR 36 put the settled answers in the preamble so a person could read the whole argument without scrolling to it.
Since ADR 72 the page shows what was settled as a folded block above the open questions, and its navigation links to it.
The preamble then said everything the block said, a third time, and by the second round of a real design it ran past 900 characters.
A person opening the page read a wall of text before the first thing they had to decide.

## Rejected

- Fixing the skill's wording and nothing else.
  The wording was what grew the preamble in the first place, and nothing would stop it growing again.
- A recap by name only, such as "Settled: Questions 1 to 7".
  It is shorter, and it still repeats the block and the count above the title.
- No preamble on a design page at all.
  One sentence on what the round asks is worth its line, and a round with nothing to say leaves it out already.
- A looser limit of 300 characters.
  It fits two sentences, and two sentences is where a recap starts again.
