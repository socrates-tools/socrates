# 81. A primitive is a custom element in the light DOM

- Status: superseded by ADR 83
- Date: 2026-10-21
- Part of: typescript-client-on-primitives

## Decision

Every primitive is an autonomous custom element, such as `<ui-dialog>`, defined by the library with `customElements.define` and drawn without a shadow root.
Its parts are children marked with `data-part`.
The contract ADR 80 describes is the set of element names and the parts each one holds.

A revision replaces the whole surface with a freshly rendered one.
Today the client follows that with `bind()` and undoes the bindings that came before with `releaseEverythingBehind()`, and both are code that has to be right on every swap.
A custom element is connected and disconnected by the browser itself, so attaching behaviour and releasing it again become a primitive's `connectedCallback` and `disconnectedCallback`, and no registry of ours walks the page.

The light DOM is deliberate.
Without a shadow root, `tokens.css` and the page's stylesheet reach a primitive as they reach everything else, the Content-Security-Policy applies unchanged, and a domain class such as `.confirmation` can sit on the element so the browser suite's selectors keep finding what they find today.

## Rejected

- Plain elements carrying `data-ui` and `data-ui-part`, mounted by a registry whose `mountAll(root)` returns a disposer for the revision swap to call. It is the closest to today's `bind()`, but releasing what was mounted would stay code of ours to get right.
