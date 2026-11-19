# 3: Go writes no page markup

## What to build

Go no longer holds the page or the catalog as templates, the partials, the views behind them, its own reading of a patch or of a mark, the page's words, or the Go tests of the rendered page.
The asset build writes no contract and nothing reads one.
Tailwind reads its classes from the client and the landing page alone.
agora's instructions, the suite's map and the contributing guide say the page is drawn by the client, and say who needs Node.

## Acceptance criteria

- [ ] Adding a primitive takes a Lit element and its catalog states, and no Go file.
- [ ] The asset build writes no contract, and nothing in the suite reads one.
- [ ] A class written with an arbitrary value in the client still fails the suite, and no Go template is read for classes.
- [ ] agora's instructions name no partial and no page template, and describe the page as drawn by the client.
- [ ] The full verification passes on a machine with no Node.

## User stories

- 1, 7, 14

## Blocked by

- 1.3-the-catalog-drawn-by-the-client.md
- 2-the-page-drawn-in-the-browser.md
