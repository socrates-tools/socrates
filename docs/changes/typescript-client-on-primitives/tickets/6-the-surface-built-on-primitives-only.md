# 6: The surface is built on primitives only

## What to build

The surface's own bindings are gone: a revision fetches the page and replaces the surface, and every primitive's behaviour is attached and released by the browser.
What remains in the surface package is exactly the pieces agora's glossary names, each built from primitives.
`agora/AGENTS.md` states the layout of the library and the surface and the invariants this change leaves: the contract, the primitives as custom elements in the light DOM, and the development build.

## Acceptance criteria

- [ ] The client holds no code that binds or releases listeners across a revision by hand, and the browser suite passes unchanged.
- [ ] A stylesheet rule in the surface package that styles a piece the glossary does not name fails the suite, naming the rule.
- [ ] After several revisions in a row, the page holds no listener left from an earlier surface.
- [ ] `agora/AGENTS.md` names the library and the surface in its layout and states the three invariants.

## User stories

- 7, 8, 11, 20

## Blocked by

- 5.1-controls-as-primitives.md
- 5.2-marks-and-status-as-primitives.md
- 5.3-code-blocks-as-primitives.md
- 5.4-app-shell-as-a-primitive.md
