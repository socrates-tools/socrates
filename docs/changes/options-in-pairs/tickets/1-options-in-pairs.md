# 1: Options sit two to a row on a wide window

## What to build

On a wide window, the options of a decision sit two to a row instead of one long column.
Pairs appear only when the window is 1280px or wider and the option group is at least 600px wide, which on the real page starts at about 1290px, and there are never more than two to a row.
Options fill across, then down, in the order the document lists them.
*None of these* is one more cell, like any option, and a lone last cell keeps half width.
Multi-select decisions pair the same way.
The design system's `OptionGroup` already states this rule, and the page follows it.

## Acceptance criteria

- [ ] In a window 1440px wide, a single-choice decision with three options shows options 1 and 2 side by side on the first row, and option 3 and *None of these* side by side on the second, in two columns of equal width.
- [ ] In a window 1440px wide, the checkboxes of a multi-select decision sit two to a row in the same way.
- [ ] In a window 1440px wide, a decision with an odd number of cells shows the last one in the first column, as wide as the others, with the cell beside it empty.
- [ ] In a window 1920px wide, options still sit two to a row and never three.
- [ ] The same decision that pairs at 1440px stacks in one column at 1200px, though its group is wider than 600px there, and at exactly 1280px, where its group is 590px wide.
- [ ] In a paired decision, Tab moves focus through the options in the listed order: option 1, option 2, option 3 and on.

## User stories

- 1, 2, 3, 4, 5, 6, 7, 8, 9, 10

## Needs a person

The design system's `Option.jsx`, `Option.d.ts` and `Option.prompt.md` under `docs/ui-design/components/forms` are changed in the maintainer's working tree and not committed.
An implementer in a worktree of its own sees only what is committed, so the maintainer commits that change before this ticket is handed out.

## Blocked by

None
