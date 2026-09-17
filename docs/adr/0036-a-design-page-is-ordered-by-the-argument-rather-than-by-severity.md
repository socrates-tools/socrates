# 36. A design page is ordered by the argument rather than by severity

- Status: accepted
- Date: 2026-09-17
- Part of: socrates-design

## Decision

A design document states no severity on any decision or source entry.

The surface orders cards by the severity the source carries, most severe first and informational last, and breaks ties in the order the document listed them.
A document stating no severity anywhere is therefore one tie, and every card renders where the session put it.

That is what a design page needs.
The order that matters in an interview is what was settled, then what the person skipped, then what is newly asked, and none of that is a severity.
A design question is not more or less severe than another one, it is earlier or later in the argument.

The surface is unchanged, and the invariant that no field of a document reaches the order is unchanged: the order the document listed them in is already the stated tiebreak.

The cost accepted is that no design question can be marked as the one that matters most.
Against it, the preamble carries one line per settled question, so the person reads the whole argument without scrolling to it.

## Rejected

- Ordering by severity, with open questions first and settled ones last as informational.
  The page then opens mid-argument, and a person arriving at round four has to scroll down to recover what was agreed before they can read what is being asked.
- Giving the surface an order a document can state.
  The surface deciding order, and no field of a document reaching it, is a rule agora has because a caller that can reorder the page can bury a decision.
  Nothing about a design needs it, since the tiebreak already gives the caller its order in exchange for stating no severity at all.
- Putting the settled questions in the source and only the open ones as cards.
  The source carries a name, a severity, an audience and a location, and no body, so a settled card could not show the person the words they answered with.
