# 84. Tailwind paints only the tokens

- Status: accepted
- Date: 2026-10-21
- Part of: lit-and-tailwind-primitives
- Amends: ADR 71

## Decision

The primitives, the surface and the landing page are styled with Tailwind utilities, and Tailwind's own theme is not part of the design.
`tokens.css` stays where every token is stated, and illustrations keep loading it at `/assets/tokens.css`.
The library's Tailwind stylesheet clears every default with `@theme { --*: initial; }` and then names each token as a utility in an `@theme inline` block, so `bg-accent`, `text-fg-2`, `rounded-pill` and `font-body` exist and `bg-red-500` and `rounded-lg` do not.
`tokens.css` gains one `--spacing` step, so `p-4` is 16px and the raw px values the stylesheets write today become steps of one scale.

No class string carries an arbitrary value.
A check reads every class string in the client's TypeScript, the partials and the landing page, and fails on a utility written with `[ ]`, naming the file and the line.
The palette check reads stylesheets and nothing else, so an arbitrary colour in a class string is the one way a value would reach the page without being a token.
The next engineer reaching for `p-[13px]` is meant to add a token or use a step instead.

## Rejected

- Moving the tokens into `@theme`, with `tokens.css` generated from it for illustrations and the landing page.
  It was declined on the card.
- Clearing only colour and type and keeping Tailwind's other scales.
  The card said values such as `rounded-lg` and `shadow-md` would then be Tailwind's rather than the design's, and it was declined.
- Allowing arbitrary sizes and refusing only arbitrary colours and fonts.
  It was declined on the card.
