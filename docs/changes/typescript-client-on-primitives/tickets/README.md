# agora's client is TypeScript, built on a library of primitives

This change rewrites agora's client in strict TypeScript on the architecture it already has, and builds the page on an internal library of primitives, `@agora/ui`.
Nothing a person sees on the page changes.
Ticket 1 is the prefactor everything builds on, 3.1 and 4 carry the first primitive end to end, and 5.1 to 5.4 migrate the rest in parallel.

Spec: [../spec.md](../spec.md). Decisions: [ADR 79](../../../../agora/docs/adr/0079-the-client-is-type-checked-by-the-build-rather-than-by-lint.md), [ADR 80](../../../../agora/docs/adr/0080-the-library-declares-what-a-primitive-is-and-go-checks-its-markup-against-it.md), [ADR 81](../../../../agora/docs/adr/0081-a-primitive-is-a-custom-element-in-the-light-dom.md).

## Working through the tickets

- Commit this folder before handing out any ticket.
- A ticket may start once every ticket it names as a blocker is done.
- Mark a ticket in progress before handing it out.
- Mark it done once its acceptance criteria hold and its work is committed.
- Mark it abandoned, with a reason, where building it turns out to be wrong, and in the same change every ticket it blocked takes its blockers in its place.

## Tickets

| Ticket | Blocked by | Status |
| --- | --- | --- |
| [1: The client is strict TypeScript, type-checked by the build](1-client-in-typescript.md) | None | done |
| [2.1: A development build serves the bundle from disk and rebuilds it on save](2.1-development-build.md) | 1-client-in-typescript | done |
| [2.2: The tokens, fonts and mark live in the library](2.2-tokens-in-the-library.md) | 1-client-in-typescript | done |
| [3.1: The dialog is the first primitive, and the contract holds it to its markup](3.1-the-dialog-as-the-first-primitive.md) | 2.2-tokens-in-the-library | done |
| [3.2: The landing page is painted by the library's tokens](3.2-landing-page-on-the-tokens.md) | 2.2-tokens-in-the-library | done |
| [4: The catalog draws every primitive in every state](4-the-catalog.md) | 3.1-the-dialog-as-the-first-primitive, 2.1-development-build | done |
| [5.1: The page's controls are primitives](5.1-controls-as-primitives.md) | 4-the-catalog | done |
| [5.2: Badges, progress, status and tooltips are primitives](5.2-marks-and-status-as-primitives.md) | 4-the-catalog | done |
| [5.3: Code blocks and diff blocks are primitives](5.3-code-blocks-as-primitives.md) | 4-the-catalog | in progress |
| [5.4: The app shell is a primitive](5.4-app-shell-as-a-primitive.md) | 4-the-catalog | done |
| [6: The surface is built on primitives only](6-the-surface-built-on-primitives-only.md) | 5.1-controls-as-primitives, 5.2-marks-and-status-as-primitives, 5.3-code-blocks-as-primitives, 5.4-app-shell-as-a-primitive | pending |
