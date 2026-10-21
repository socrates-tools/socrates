# 3: Every primitive a Lit element, and the light-DOM primitive gone

## What to build

Every primitive is a Lit element, so adding one is a Lit element, a partial that writes its element, and its states in the catalog, and nothing else.
The contract declares attributes and slots only, a primitive's parts and `data-part` are gone, the per-primitive stylesheets, the stylesheet importing them and the parts lookups are removed, and the suite asserts that every primitive has a shadow root where it once asserted none did.
`agora/AGENTS.md` and `CONTRIBUTING.md` describe the client as it is now built.

## Acceptance criteria

- [ ] A primitive drawn without a shadow root fails the suite, naming it.
- [ ] The contract carries no parts, and a partial writing `data-part` fails the Go contract test.
- [ ] Adding a primitive takes a Lit element, a partial and its states in the catalog, and the suite fails on one the catalog does not draw, naming it.
- [ ] No hand-written stylesheet styles a primitive.
- [ ] `agora/AGENTS.md` states the invariants of the shadow root, the contract of attributes and slots and the theme, and none of the light DOM's, and `CONTRIBUTING.md` asks for Node and `pnpm` of somebody who touches `web/` or `site/`.
- [ ] `make verify` passes.

## User stories

- 3

## Blocked by

- 2.1-the-indicators-and-the-tooltip.md
- 2.2-the-controls.md
- 2.3-the-choice-group.md
- 2.4-the-card-and-the-dialog.md
- 2.5-the-app-shell.md
- 2.6-the-code-and-diff-blocks.md
- 2.7-the-pieces-and-the-catalog-on-utilities.md
