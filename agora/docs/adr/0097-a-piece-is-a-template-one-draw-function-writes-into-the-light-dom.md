# 97. A piece is a template one draw function writes into the light DOM

- Status: accepted
- Date: 2026-11-19
- Part of: the-page-is-drawn-in-the-browser

## Decision

`agora/web/packages/surface` draws the page with one exported function, `draw(page)`, which returns a new `.surface` element for the first load and for every revision.
Each piece the glossary names, such as the decision, the composer, the receipt, the index and what was settled, is a private function in that module returning a Lit template built from the library's primitives.
No piece is a custom element, and none has a shadow root.
Each piece keeps the one class that names it, such as `.decision` or `.composer`, so the page's script and the browser suite find it as they do today.

The next engineer will see every primitive as a Lit element with a shadow root and expect a piece to be one too.
It is not, because a piece behind a shadow root hides its controls from the browser suite and from the page's own links, and the line between a primitive and a piece is still the glossary's: a primitive names nothing in agora's domain.

The class names stay a convention the type checker cannot see, and a piece cannot be drawn alone without a whole page.

## Rejected

- Each piece its own Lit element with a shadow root, such as `<agora-decision>` and `<agora-composer>`, taking properties and sending events.
  The card's cost was that the suite could no longer reach a button two shadow roots deep, so each piece would hand out its controls, and links to a decision, the confirmation's jump and labels across pieces would need hand-written routing.
  ADR 85's card had offered the same and it was declined then too.
