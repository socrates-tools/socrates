# 75. The page fills the window as an app shell

- Status: accepted
- Date: 2026-10-16
- Part of: the design system's app shell

## Decision

The page was a 1200px column centred in the window: a sticky top bar, a 240px rail beside cards of up to 800px, and the message and the presses at the end of the cards.
On a wide screen most of the window stayed empty while code and diffs wrapped inside an 800px card, and the presses were as far away as the last card.

From 960px the page now fills the window as the design system's app shell, and the window itself never scrolls.
A full-width top bar carries the wordmark, the exchange's title as a breadcrumb, and the status.
Below it the navigation is a 272px column on the sunk surface, the decisions take the rest of the width up to 1600px, and each scrolls on its own.
The decisions open with a page header that counts what is open and what was settled, then the title and the preamble, divided from the cards by a hairline.
The navigation marks the decision nearest the top of the window as current.

From 1280px a 304px reply panel joins on the right.
It summarises what will go, one line per decision as the person answers, and keeps the message and both presses pinned at its bottom, so sending is always in reach.
It is padded 16px at its sides rather than 24px, the narrowest at which both presses still fit side by side, so it takes as little of the decisions' width as it can.
Once something is sent the receipt takes the summary's place, and the way back to amending is pinned where the presses were.
Below 1280px the same composer closes the decisions instead, and below 960px everything stacks into one column that scrolls with the window, navigation first.
There is one composer, which the page moves between the panel and the end of the decisions as the window crosses 1280px, so what a person typed goes with it.

The preamble keeps to 72 characters a line, as the design system sets it, and no longer runs to the cards' right edge.
The press that waits for a reply stays first in the form, so Enter in a field still sends rather than finishes, and the design's order is drawn rather than written.

This amends ADR 53 only where it placed the top bar over a scrolling page: from 960px the page does not scroll, so the bar stays by being outside what does.

## Rejected

- Keeping the centred column and only widening it.
  The presses would still be as far away as the last card, which is what the panel is for.
- Rendering the composer twice, once in the panel and once after the decisions, and showing one.
  Two message fields would hold two messages, and the form would send whichever one its script read.
- Marking the current decision by where cards cross a band of the window.
  On a short window the first card can start below the band, and then nothing is marked at all.
