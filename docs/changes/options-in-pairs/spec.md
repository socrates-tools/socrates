# Options sit two to a row on a wide window

Decisions: none.

## Problem Statement

On the agora page every option of a decision stacks in one column, however wide the window is.
Since ADR 75 the decisions take the full width between the navigation and the reply panel, up to 1600px, so on a wide screen each option is a long, thin bar with most of its width empty.
A decision with four or five options plus *None of these* takes a lot of height, and a person scrolls further than the content needs.

The design system has already moved.
Its `OptionGroup` in `docs/ui-design/components/forms/Option.jsx`, changed in the working tree and not yet committed, lays options out two to a row on a wide screen.
The page does not follow it.

## Solution

The page lays out options the way the design system's `OptionGroup` does.
Options stack in one column on phones, tablets and laptops.
When the window is 1280px or wider and the option group is at least 600px wide, the options sit two to a row, in equal columns.

On the real page at exactly 1280px the option group is 590px wide: 1280 less the 272px navigation, the 304px reply panel, 64px of main padding and 48px of card padding.
So pairs start at a window of about 1290px, and a 1920px window still shows two to a row, never more.

*None of these* is laid out like any other option: it takes the next cell and may sit beside the last option.
Options read across, then down, in the order the document lists them, so the first row holds options 1 and 2.
With an odd number of cells, the last one keeps half width and the cell beside it stays empty.

Multi-select decisions pair the same way.
The design system's uncommitted `Option` change lands in the same change, so the page and the design system agree.

## User Stories

1. As a person answering a round on a wide monitor, I want options laid out two to a row, so that a decision with many options takes about half the height and I scroll less.
2. As a person answering a round on a laptop, I want options to keep stacking in one column, so that each option keeps enough width for its label and meaning to read comfortably.
3. As a person answering on a phone or tablet, I want the page to look as it does today, so that nothing gets squeezed on a narrow screen.
4. As a person scanning a decision, I want options to read across then down in the order they were written, so that the first option, often the recommended one, is still the first thing I see, top left.
5. As a person using the keyboard, I want focus to move through options in the same order I read them, so that tabbing never jumps around the grid.
6. As a person comparing options, I want every option the same width, including the last one when the count is odd, so that no option looks more important because of its size.
7. As a person who does not like any option, I want *None of these* to still be there, in the grid with the others, so that I can always decline and write my own answer.
8. As a person answering a multi-select decision on a wide screen, I want its checkboxes paired the same way, so that both kinds of decision look alike.
9. As a person who resizes the window, I want the layout to switch between one and two columns at a stable point, so that options do not jump between three and four columns as the window moves.
10. As a contributor changing the page, I want the design system and the page to state the same rule, so that I can follow the design system without first finding where the page differs.

## Implementation Decisions

- `agora/render/page.gohtml`: inside `<fieldset class="options">`, after the `legend`, the options and the *None of these* label are wrapped in one element of their own, the list the grid applies to.
  The legend stays outside it, as the design system's `so-options-list` does, because a fieldset's rendered legend is not a grid item.
  Nothing about each option's markup changes, and their order in the source stays the listed order with *None of these* last.
- `agora/web/src/surface.css`:
  - `.options` becomes a size container (`container-type: inline-size`) and stops being the flex column.
  - The new list element carries today's column: `display: flex; flex-direction: column; gap: 8px`.
  - Under `@media (min-width: 1280px)` and `@container (min-width: 600px)`, the list becomes `display: grid; grid-template-columns: repeat(2, minmax(0, 1fr))`, with the same 8px gap.
  - The grid uses the default row flow, so options fill across then down, and options in one row stretch to the same height.
  - `.none` gets no column span: it is one cell like any option.
- `agora/assets/bundle/surface.css` is regenerated with `make assets`, never edited by hand.
- `docs/ui-design/components/forms/Option.jsx`, `Option.d.ts` and `Option.prompt.md`: the working-tree change that introduced the rule is committed as part of this change.
  The design system's `DecisionCard.jsx` already puts *None of these* inside `OptionGroup` as one more cell, which matches what was settled, so it does not change.
- The two thresholds are the design system's, 1280px for the window and 600px for the group, and are not tuned to the page's own widths.
  At 1280px exactly the page therefore still stacks, and that is accepted.
- No document field, rendering of an option, or answer the page sends changes.
  What is posted back is the same as before, because only layout changes.

## Testing Decisions

The layout is tested in `internal/e2e`, in the real browser, at stated window widths, the same seam as `internal/e2e/shell_test.go` and its `openingAt(t, document, width, height)`.
A test reads the bounding box of every `.option` and `.none` in a decision and asserts which of them share a row.

Scenarios, each a test named for it:

- A window 1440px wide pairs the options of a single-choice decision: with three options and *None of these*, options 1 and 2 share the first row, option 3 and *None of these* share the second, and both columns are the same width.
- A window 1440px wide pairs the checkboxes of a multi-select decision the same way.
- With an odd number of cells at 1440px, the last one is as wide as the others and starts in the first column.
- A window 1920px wide still holds two to a row, never three.
- A window 1200px wide, below 1280px, stacks every option in one column even though the group is wider than 600px there, because the window rule holds.
- A window exactly 1280px wide stacks every option, because the group is 590px there, below 600px.
- On a narrow window, below 960px, options stack as they do today.
- Tabbing through a paired decision moves focus in the listed order: option 1, option 2, option 3 and on.

Each test must fail against today's page for the right reason (options stacked where they should pair) before the CSS changes, as the repository's test-first rule asks.
`make test` already checks that the committed bundle matches its source, so a forgotten `make assets` fails the suite.

## Out of Scope

- More than two options to a row on very wide windows.
- Laying *None of these* out across the full row, or giving it any other special place in the grid.
- Column-first ordering (down the first column, then the second).
- Stretching a lone last option across its row.
- Screenshot comparison against committed images.
- Any change to the shell's own breakpoints, the navigation, the reply panel or the card padding.
- Any change to how a document states options, or to what an answer carries.

## Further Notes

The thresholds come from the design system, so pairs start just past 1280px on the real page rather than at it.
If that gap ever matters, change the design system's rule first and let the page follow it, rather than tuning the page alone.

No ADR was written.
The rule is a layout choice that is easy to reverse, so it does not meet all three of the bar's conditions.
No glossary entry changed, because no new term was settled.
