# 53. The page is light whatever the reader prefers

- Status: accepted
- Date: 2026-10-04
- Part of: the-aegean-page

## Decision

The page paints one palette, a whitewashed limestone ground with navy ink, Aegean blue for what the person chooses or presses, and terracotta clay for the recommendation and for the audience that needs a person, and it states `color-scheme: light`.
It does not follow `prefers-color-scheme`: a reader who prefers dark is shown the same page as one who prefers light.

This supersedes ADR 45.
What ADR 45 held beyond the hue still holds: one palette stated once, every colour a token, the contrast floor measured under both preferences, and 12px as the smallest text the page renders.

Aegean marks what is chosen and what to press, the open status, links and the keywords of code, and no body text is Aegean.
Clay marks only the recommendation, the audience that needs a person, and the calls in code.

The page carries the brand rather than the name of the server: the wordmark says socrates, set in Cinzel, and Cinzel appears nowhere else.
The title, the questions and everything a person reads or presses are set in Manrope.
The wordmark and the status sit in a bar that stays at the top of the window, so pressing an entry lands its card below that bar rather than under it.

On a card, a severity or audience mark also says its name in a word beside its icon, so the meaning ADR 47 shows on hover and on focus is also read at a glance.
In the index, where there is room only for the icon, the mark is the icon alone.
The tone ADR 47 drew in gold is drawn in clay.

When the exchange is over, the page says so under the title and stops drawing the fields nothing will read: each decision's own words, and the message with its presses unless they hold what was just sent.

## Rejected

- Keeping the dark palette beside the light one and following the reader's preference.
  ADR 45 rejected this for the same reason it still fails: every token stated twice and every contrast held twice, for a design drawn once.
- Keeping the dark palette.
  The light redesign is the one drawn, and the dark one reads as a terminal rather than as a page a person settles something on.
- Disabling the fields of a closed exchange rather than removing them.
  A disabled field still asks to be read, and a person who meets one on a page that is over looks for the way to answer.
