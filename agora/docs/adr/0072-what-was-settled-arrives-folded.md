# 72. What was settled arrives folded

- Status: accepted
- Date: 2026-10-16
- Part of: the-kylix-page

## Decision

The block of what was settled arrives folded while the page holds a question still open.
Its heading still says what it is and how many it holds, and pressing it, or the index entry that leads to it, unfolds it.
A page holding nothing but what was settled arrives with it unfolded, because folding it would leave the page showing nothing.

By a later round a person has already answered everything the block lists.
Shown open, it puts every earlier answer between them and the first question still owed one, and the list grows with every round.
The preamble already says what the earlier rounds settled in a sentence, so the block is there to check against rather than to read first.

This refines ADR 46, which placed the block ahead of every open question and left it open.

## Rejected

- Keeping the block open.
  The person scrolls past answers they gave themselves before reaching anything to answer, and more of them with every round.
- A document field choosing whether the block arrives folded.
  Every tool writing a document would have to decide it, and the right default does not depend on the tool.
- Showing only the latest round's settled answers and dropping the rest.
  The surface counts no rounds, per ADR 46, and a folded block costs nothing to keep whole.
