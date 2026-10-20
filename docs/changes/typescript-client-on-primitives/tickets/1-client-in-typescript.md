# 1: The client is strict TypeScript, type-checked by the build

## What to build

agora's client is rewritten in strict TypeScript, without changing anything a person sees on the page.
The client source becomes a pnpm workspace, and the page's script is its surface package.
The one script is split into modules along the lines it already has: sending and the receipt, the confirmation, the live revision, code blocks, the tooltip and the navigation.
Building the bundle type-checks the workspace first, and a type error stops the build and names the file and line.
Building also refuses a comment in the TypeScript, other than the notice comment that carries a bundled library's licence.
The suite's conventions read TypeScript wherever they read the client's JavaScript today, and the stamp that proves the committed bundle came from its source covers the new layout.

## Acceptance criteria

- [ ] The client's source holds no JavaScript file other than the build script, and the type checker runs in strict mode over all of it.
- [ ] A type error introduced into the client makes `make assets` exit non-zero, naming the file and line, and writes no bundle.
- [ ] A comment added to a TypeScript source makes `make assets` exit non-zero, naming the file, while the highlight.js notice comment still reaches the bundle.
- [ ] A licence header added to a TypeScript source fails the suite, as one added to the client's JavaScript does today.
- [ ] Changing any TypeScript source without rebuilding fails the suite's check that the committed bundle was built from its source.
- [ ] The bundle still serves `surface.js`, `surface.css`, `tokens.css`, `height.js`, the fonts and the mark under `/assets/`, and the browser suite passes unchanged.
- [ ] `make verify` passes on a machine with no Node installed.

## User stories

- 1, 2, 3, 4, 19, 20

## Blocked by

- None
