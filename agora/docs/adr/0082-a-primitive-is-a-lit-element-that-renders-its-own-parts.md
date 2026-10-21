# 82. A primitive is a Lit element that renders its own parts, and Go checks the attributes and slots it writes

- Status: accepted
- Date: 2026-10-21
- Part of: lit-and-tailwind-primitives

## Decision

Every primitive is a Lit element in `agora/web/packages/ui`, and its template, its utilities and its behaviour are one TypeScript file.
Go no longer writes a primitive's parts.
A partial in `agora/render/partials` writes the element, the attributes it declares, and the content it slots, and nothing inside it.

The contract stays generated, and stays what Go reads.
`make assets` reads every element's declared reactive properties and the names of the slots its template renders, and writes them to `agora/assets/bundle/contract.json`.
A Go test renders every partial and fails on an attribute or a slot a partial writes that no element declares, and on a declared one no partial writes, naming the element and the attribute or slot.
Changing a partial still needs no Node, which is the promise ADR 79 keeps.
Changing what a primitive draws now needs Node, because what it draws is TypeScript.

This supersedes ADR 80, which kept a primitive's markup in Go so that changing it needed no Node.
A primitive split across a Go partial and a TypeScript element was one of the four things this change is for, and the split is what ADR 80 chose.

## Rejected

- Go keeps the markup, now with Tailwind classes, and Lit adds behaviour only.
  The card offered it as the way ADR 79 and ADR 80 would both stand, and it was declined: a primitive would stay split across two languages.
- Lit renders the whole page, with Go handing it the document as data.
  The card said it reaches far past the primitives, and it was declined.
- Dropping the contract and relying on the catalog in a real browser alone.
  The card said a misspelt attribute would then be caught only where the browser suite drives that state, and it was declined.
