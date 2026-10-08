# 45. The page is dark whatever the reader prefers

- Status: superseded by ADR 53
- Date: 2026-09-30
- Part of: the-kylix-page

## Decision

The page paints one palette, a warm black ground with gold for what the person chooses or presses, and states `color-scheme: dark`.
It no longer follows `prefers-color-scheme`: a reader who prefers light is shown the same page as one who prefers dark.

The contrast floor still holds, and it is measured under both preferences, because the reader's preference is exactly what the page now ignores.
Every colour is still a token stated once, in one block rather than in two.

The smallest text on the page is 12px: the overlines, the locations, the badges and the recommendation mark.
The floor the page states is lowered to that size and every text is still held to it.

What this costs is the reader who prefers light for a reason, whether a bright room or a need for dark text on a pale ground.
What it buys is one palette to hold to the floor instead of two, and a page whose one accent colour is spent sparingly.
Gold marks what is chosen and what to press, and beyond that only the wordmark, the open badge, links, the recommendation mark, the audience that needs a person, and the keywords and types of code, as the design draws them.
No body text is gold.

## Rejected

- Keeping a light palette beside the dark one.
  Every token would be stated twice and every contrast held twice, and the gold that marks a choice on black has no counterpart on white that stays both gold and legible.
- Deriving a light palette from the dark one automatically.
  A derived palette is one nobody drew, and the floor would be met by arithmetic rather than by design.
