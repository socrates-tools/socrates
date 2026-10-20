# 80. The library declares what a primitive is, and Go checks its markup against it

- Status: accepted
- Date: 2026-10-21
- Part of: typescript-client-on-primitives

## Decision

A primitive has two halves in two languages.
Go writes its markup, as a partial in `agora/render`, because Go renders the page.
The library under `agora/web/packages/ui` gives it its style and its behaviour.

The library is where a primitive is declared.
`contract.ts` names every primitive's element and every part it holds, and `make assets` writes the same contract to `agora/assets/bundle/contract.json`, as it already writes `agora/document/languages.json` from the bundled highlighter.
A Go test renders every partial and fails on an element or a part the contract does not declare, and on a declared one no partial writes.

The markup stays in Go, so changing a partial needs no Node, which is the promise ADR 79 keeps.
The contract lives with the behaviour that depends on it, so a primitive's TypeScript and its declaration change in one place, and Go is told about the change through a generated file it already knows how to read.

## Rejected

- Each primitive's folder holding its own markup as a `.gohtml` that `make assets` copies into the bundle for `agora/render` to parse from the embed. A primitive would be one folder, but editing markup would need the Node build.
- Go declaring each primitive and `go generate` writing the library's `contract.ts`. Go would be the source of truth for a contract whose only reader is the TypeScript.
