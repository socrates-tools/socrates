# 85. The glossary's pieces stay Go markup laid out by utilities

- Status: accepted
- Date: 2026-10-21
- Part of: lit-and-tailwind-primitives

## Decision

Lit renders the primitives and nothing else.
The pieces the glossary names, such as the decision, the composer, the receipt and the index, are still written by Go in `page.gohtml` and its partials, composing primitives and laid out with Tailwind utilities that Tailwind reads from the `.gohtml` files.
The surface's stylesheet empties to what utilities cannot say, and the check that held every class of it to a piece the glossary names, `convention.PiecesOfTheSurface`, is retired.
A piece keeps the one class that names it, such as `.decision` or `.composer`, for the script and the browser suite to find it by.

The next engineer will see Lit elements for a button and Go markup for a decision and expect one of them to be the other.
The line is the glossary's own line between a primitive and a piece: a primitive names nothing in agora's domain and is drawn by the library, and a piece is agora's domain and is drawn by the page Go renders.

## Rejected

- Each piece as a Lit element too, such as `<agora-decision>` and `<agora-composer>`, rendering itself from attributes.
  The card said Go would then write a page of domain elements, which is most of the way to Lit rendering the whole page, and it was declined.
