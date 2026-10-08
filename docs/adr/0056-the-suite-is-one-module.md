# 56. The suite is one module

- Status: accepted
- Date: 2026-10-08
- Part of: monorepo

## Decision

The suite is one Go module with one directory per tool, and every tool is released together at one version, in one archive.
The shared layer lives in the root `internal/`, so Go's own rule keeps each tool's `internal/` invisible to the others at no cost.

This supersedes ADR 13, which decided that a review of somebody else's change needs no program of its own "the way `agora` did", and so took a tool to be a program with a build and a release of its own.
That premise no longer holds: the suite publishes one binary, per ADR 32 and ADR 33, and a tool is a directory and a namespace inside it.

## Given up

An independent release cadence for the tools.
It was the goal the first draft of the monorepo pursued, and it is deliberately given up.
A change to the page cannot ship without a release of every other tool, and a change to review cannot be held back while the page is fixed.
If that becomes intolerable, the way back is a multi-module split, and its cost is the one written under Rejected below.

## Rejected

- One module per tool, joined by a workspace.
  A workspace only fixes local development, so the layer the tools share would have to become a module of its own, published and versioned, and every change to a convention would become a release of four modules in an order somebody has to remember.
- One repository per tool.
  It has the same cost as separate modules and adds a second place to keep the documents, the verification and the release in step.
- Per-tool archives.
  agora serves one document and waits, which is useless without something writing documents, so it is a component and not a product.
