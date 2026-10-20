# agora's client is TypeScript, built on a library of primitives

Decisions: ADR 79, ADR 80, ADR 81.

## Problem Statement

agora's page is getting harder to change safely.
Go renders it from `agora/render/page.gohtml`, and one untyped script, `agora/web/src/surface.js`, gives it its behaviour in 879 lines.
That one script holds the form, the receipt, the confirmation, the live revision, the code blocks, the tooltip and the navigation.
Its stylesheet holds about 2,000 lines named after agora's features (`.decision`, `.composer`, `.hunk`, `.receipt`) rather than after the pieces they are made of.
The design system exists only as `tokens.css`.
Every button, dialog and tooltip is styled and wired where it is used, so a second use of one is a copy of it.
Nothing types the markup the script reaches into, so a renamed class in a template is found by the browser suite or not at all.
When a revision arrives, the script swaps the whole surface, binds everything again with `bind()`, and undoes what it bound before with `releaseEverythingBehind()`, and that code has to be right on every swap.
Changing the client means editing, running `make assets`, rebuilding the binary, serving again and reloading, with nothing that watches the source.
The landing page under `site/` carries its own copy of the palette under other names (`--aegean-deep`, `--sea-wash`) beside the same three fonts.

## Solution

The client becomes strict TypeScript on the architecture it already has: Go still renders the page and esbuild still bundles it.
React and a component framework are not adopted.

The client is built on an internal library, `@agora/ui`, under `agora/web/packages/ui`, inside a pnpm workspace whose other package is the surface itself.
The library holds the tokens, the fonts, the mark and every primitive.
A primitive is any piece of the page agora's glossary does not name, and every piece the glossary does name is built from primitives.

A primitive's markup is a Go partial and its style and behaviour are the library's.
The library declares every primitive in one contract, the build writes it out as JSON, and a Go test fails where the partials and the contract disagree, per ADR 80.
Every primitive is a custom element drawn without a shadow root, so the browser attaches and releases its behaviour on every revision, per ADR 81.

`build.mjs` type-checks the client before bundling it, so `make assets` refuses a type error and a contributor who never touches `agora/web` still needs no Node, per ADR 79.
A build tag makes a binary serve the bundle from disk, and esbuild in watch mode rebuilds it on save, so a change to the client is seen with a reload.
The catalog draws every primitive in every state for the browser suite, and is served by the test suite and by that development build, never by a released binary.
The landing page loads the library's tokens and fonts and drops its own copy of the palette.

## User Stories

1. As a contributor changing agora's page, I want the client written in strict TypeScript, so that a mistake in what the script expects of the page is a build error rather than something the browser suite may or may not catch.
2. As a contributor who only changes Go, I want to build, lint and test socrates without Node or pnpm installed, so that the client's toolchain costs me nothing while I never touch it.
3. As a contributor who introduces a type error in the client, I want `make assets` to refuse it and name the file and line, so that no bundle is ever built from source that does not type-check.
4. As a contributor reviewing a pull request, I want a client source that does not type-check to be impossible to commit a bundle for, so that I never have to ask whether the committed bundle came from source that compiles.
5. As a contributor working on the page, I want to save a TypeScript or CSS file and see the change on a reload of a page I am already serving, so that I am not rebuilding and reinstalling the binary for every edit.
6. As somebody who installed a released socrates, I want my binary unable to serve assets from anywhere but its own embedded bundle, so that nothing on my machine can change what the page I answer looks like or does.
7. As a contributor adding a feature to the page, I want a button, a field, a dialog, a tooltip, a code block or a copy control to already exist as a primitive, so that I build with it rather than copying the nearest one.
8. As a contributor, I want to know whether a piece belongs in the library or in the surface by reading agora's glossary, so that the boundary is a rule rather than a judgement made again for every piece.
9. As a contributor changing a primitive's markup in a Go partial, I want a test to fail when the markup names an element or a part the library does not declare, or leaves out one the library declares, so that markup and behaviour cannot drift apart unnoticed.
10. As a contributor changing a primitive's markup, I want to do it in Go without running the Node build, so that markup stays where the page is rendered.
11. As a contributor, I want a primitive's behaviour attached when its element enters the page and released when it leaves, so that a revision swapping the surface never leaves a listener behind and never needs a registry of ours to get that right.
12. As a contributor, I want the browser suite's selectors on agora's own pieces, such as the confirmation, the send and the choices, to keep finding what they find today, so that moving a piece onto a primitive is proved by tests that already exist.
13. As a contributor building or changing a primitive, I want a catalog that draws it in every state through the real partials, so that I see and test it without composing a document that happens to reach each state.
14. As a contributor, I want the browser suite to prove each primitive's behaviour against the catalog, so that a primitive is covered once rather than through every feature that uses it.
15. As a contributor, I want every colour and every font the client paints to come from the library's tokens, so that the palette is stated once and a literal anywhere else fails the build.
16. As somebody writing an illustration, I want `/assets/tokens.css` to keep its path and its custom property names, so that a drawing I wrote against them keeps looking like the page.
17. As somebody visiting the landing page, I want it painted in the same palette and fonts as the page socrates serves, so that the two read as one product.
18. As a contributor changing the palette, I want to change it in the library once and have both the page and the landing page follow, so that the two cannot disagree.
19. As a contributor reading the client, I want it to carry no comments, as the Go source carries none, so that one convention holds across the repository.
20. As somebody using socrates, I want the page to look and behave exactly as it did, so that a change to how the client is built is invisible to me.

## Implementation Decisions

**The workspace**

- `agora/web` becomes a pnpm workspace with two packages: `packages/ui`, published inside the workspace as `@agora/ui`, and `packages/surface`, the page.
  `packages/surface` imports `@agora/ui`, and `@agora/ui` imports nothing from the surface.
- `typescript` joins the devDependencies.
  It ships nothing into the bundle, so `NOTICE` gains no entry for it.
- `tsconfig.json` states `strict` and the DOM library, emits nothing, and covers both packages.
- `fonts/`, `brand/`, `tokens.css` and `fonts.css` move into `packages/ui/tokens/`.
  The bundle still holds `tokens.css` at `/assets/tokens.css`, with every custom property it holds today under the same name, because illustrations load it by that path and both skills tell them to.

**The build**

- `build.mjs` runs `tsc` over the workspace before esbuild and stops with `tsc`'s diagnostics when it reports any, per ADR 79.
- `build.mjs` also refuses any comment in a `.ts` source, read with the TypeScript scanner rather than a pattern, except the `/*!` notice comments that carry a bundled library's licence.
- The suite's palette check reads the whole workspace and fails on a colour literal or a font family named in any stylesheet outside `packages/ui/tokens/`.
- esbuild bundles `packages/surface/main.ts`, the library's stylesheets and `height.ts` into the same files the bundle holds today: `surface.js`, `surface.css`, `tokens.css`, `height.js`, the fonts and the mark.
  `internal/convention.ClientBundle` and the server's asset routes are unchanged.
- `build.mjs` writes `agora/assets/bundle/contract.json` from `packages/ui/contract.ts`, as it already writes `agora/document/languages.json`.
- `bundle.sum` stamps `packages/`, `build.mjs`, `tsconfig.json`, `package.json` and `pnpm-lock.yaml`, in place of `src`, `fonts` and `brand`.
- `build.mjs --watch` rebuilds the bundle on every save and runs the same type check, printing its diagnostics without exiting.

**The development build**

- A file under the `agoradev` build tag in `agora/assets` makes `assets.Files()` an `os.DirFS` over `agora/assets/bundle` on disk, and the file without the tag keeps the embed.
  A binary built without the tag has no way to serve assets from disk, per the decision recorded in this spec.
- The `agoradev` build also serves the catalog under the surface, so a person can look at it.
- `make dev` builds `socrates` with `-tags agoradev` and runs `build.mjs --watch`.
  `make help` names it.

**Primitives**

- Every piece of the page that agora's glossary does not name is a primitive and moves into `@agora/ui`.
  From today's page, that is at least: the button, the text field, the choice that an option is drawn as, the badge a mark is drawn on, the progress bar, the status dot, the dialog, the tooltip, the copy control, the expand control, the code block, the diff block, the card, the overline and the app shell.
- Every piece the glossary names stays in `packages/surface` and is built only from primitives: the decision, the index, the mark, the settled answer, the option, the evidence, the illustration's frame, the composer, the receipt and the confirmation.
- Each primitive is an autonomous custom element named `ui-<name>`, defined with `customElements.define` and drawn without a shadow root, per ADR 81.
  Its parts are children marked with `data-part`.
  Its behaviour is attached in `connectedCallback` and released in `disconnectedCallback`.
- Each primitive's folder in `packages/ui` holds its stylesheet and, where it has behaviour, its element class.
  A primitive with behaviour exposes it as methods on its element, such as `open()` and `close()` on `ui-dialog`, which the surface calls rather than reaching into its parts.
- A domain piece carries its own class on the primitive's element, such as `<ui-dialog class="confirmation">`, so the browser suite's selectors stay as they are.
- `bind()` and `releaseEverythingBehind()` are removed once every piece they bind is a primitive or a domain piece built on one.
  What remains of the revision swap is fetching the page and replacing the surface.

**The contract**

- `packages/ui/contract.ts` names every primitive's element and the parts it holds, and the element classes are typed against it, per ADR 80.
- Each primitive's markup is a partial in `agora/render`, defined once and named after the primitive, and `page.gohtml` reaches every primitive through its partial.
- A Go test in `agora/render` renders every partial, reads `contract.json` from `agora/assets`, and fails on an element or a part the contract does not declare and on a declared element or part no partial writes.

**The catalog**

- The catalog is a document-free page drawn by the real partials, with every primitive in every state it has.
- A test helper in `internal/testing` serves it to chromedp.
  The released binary has no route to it, and the `agoradev` build serves it for a person.

**The landing page**

- `make site` copies the library's `tokens.css`, `fonts.css` and fonts into the site, and `site/index.html` loads `tokens.css` ahead of `site.css`.
- `site/site.css` drops its own `@font-face` declarations and its own palette, and every use of a site-only name moves to the library's name, such as `--aegean-deep` to `--aegean-700`.
  Its layout stays its own.

**The conventions**

- `internal/convention`'s source extensions gain `.ts`, so the licence-header check and the vocabulary check read the TypeScript as they read the JavaScript today.
- `CONTRIBUTING.md` keeps stating Node and pnpm as needed only by somebody who touches `web/`, and gains `make dev`.
- `agora/AGENTS.md`'s layout names `agora/web/packages/ui` and `agora/web/packages/surface`, and its invariants gain the contract, the light-DOM custom elements and the development build.

## Testing Decisions

- Behaviour is proved through chromedp and nowhere else: no JavaScript test runner is added.
- The page is proved by the browser suite as it stands, in `internal/e2e` and in agora's own browser tests, unchanged.
  Every extraction of a primitive keeps the page's behaviour, and the suite passing unchanged is the proof of it.
- Each primitive is proved against the catalog, following the shape of `internal/e2e/copying_test.go` and `expanding_test.go`, with a test named for each state the catalog draws.
- One test proves that a revision swap releases a primitive: it opens a dialog, delivers a revision, and asserts that the old element's behaviour no longer answers.
- The contract check is a Go test in `agora/render`, with a scenario for an undeclared element, an undeclared part, and a declared element no partial writes, and each is mutated on purpose once to show it fails.
- The type check, the comment refusal and the colour-literal refusal are proved by running `build.mjs` over a scratch copy of the workspace holding one violation each, and asserting that it exits non-zero and names the file.
  These run in the step continuous integration already runs after `make verify`, where Node is present.
- That a release binary cannot serve from disk is proved by a test in `agora/assets` built without the tag, asserting that `Files()` is the embed.
- The landing page is proved by `make site` and a check that `site/site.css` names no colour literal and no font family.

## Out of Scope

- React, a component framework, shadcn/ui, Tailwind and a change of bundler.
- A JavaScript test runner.
- A route to the catalog in the released binary.
- The landing page using primitives: it takes the tokens and fonts and no primitive's markup or style.
- Any change a person using socrates could see on the page.
- Renaming or removing any custom property `tokens.css` holds today.

## Further Notes

- On the brief page answered before this interview, a full move to React, TypeScript, shadcn/ui and Vite was weighed and kept for later.
  The triggers recorded there were state that spans cards and has to survive without a re-render, rich widgets such as comboboxes or nested dialogs, more than one page, and the typed client passing about 3,000 lines with its bugs coming mostly from keeping the DOM and state in step.
  The tokens and the primitives' stylesheets are what would carry over if one of those happens.
- The tickets follow the order the brief set, each leaving the suite green: TypeScript and its type check, the development build, the tokens moving into the library, the contract and its Go check, the primitives one at a time, the catalog, the surface built on the library only, and the landing page.
- Serving a page through a reverse proxy that keeps the original Host header, as `surface.reached_at` invites, is refused by the surface's host check unless the proxy rewrites the Host to the address the surface bound.
  That was found while running this design and is not part of this change.
