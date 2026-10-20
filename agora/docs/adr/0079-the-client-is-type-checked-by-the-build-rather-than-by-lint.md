# 79. The client is type-checked by the build rather than by lint

- Status: accepted
- Date: 2026-10-21
- Part of: typescript-client-on-primitives

## Decision

agora's client is TypeScript, checked by a strict `tsc` that `build.mjs` runs before esbuild bundles anything.
`make assets` refuses a type error, and `make lint` does not run `tsc` at all.

A contributor who never touches `agora/web` needs no Node today: `make verify` runs Go alone, and `CONTRIBUTING.md` asks for Node and `pnpm` only of somebody who touches `web/`.
Putting the type check in the build keeps that promise.
A type error still cannot reach a commit: the committed bundle is only ever what `make assets` wrote, `bundle.sum` stamps the source it was built from and fails the suite when the two disagree, and continuous integration rebuilds the bundle and diffs it against what was committed.
A source that does not type-check produces no bundle, so it cannot produce a committed one.

The next engineer will expect `tsc` beside `go vet` in `make lint`, and this record is why it is not there.

## Rejected

- Running `tsc` in `make lint`, so that `make verify` fails on a type error. Every contributor would then need Node and `pnpm`, and the toolchain `CONTRIBUTING.md` states would change for people who never touch the client.
