**The suite must be green after every task, not only at the end of a group.**
The ordering below exists for that reason: the notice file exists before any check reads it, and the bundler moves only once a test can say what moving it must produce.

**A test asserts on what ships.**
No task proves a notice is present by reading `web/src` or `web/build.mjs`. Every test that needs the notice reads the bytes under `internal/assets/bundle`, as `internal/assets/bundle_test.go` already reads them.

## 1. The terms are stated

- [x] 1.1 Add `LICENSE` at the root: the MIT licence in full, with `Copyright (c) 2026 Rui Cheng`, and verify `make test` passes unchanged
- [x] 1.2 Add a convention check that no source file opens with a licence header, and verify a unit test that it fails for a file that does and passes for the tree as it stands. Covers `No source file carries a licence header`
- [x] 1.3 Add an end-to-end test that the root states the MIT licence and names a copyright holder. Covers `The root states the terms`

## 2. What the work borrows is written down

- [x] 2.1 Establish what each embedded asset is published under: the four `web/fonts/*.woff2` faces, and every third-party library `web/build.mjs` bundles into `surface.js`
- [x] 2.2 Add `NOTICE` at the root, naming each borrowed work, its licence and that licence's full text, and verify `make lint test` is green
- [x] 2.3 Add an end-to-end test that the notice names each embedded font family and carries its licence text. Covers `An embedded font's licence travels with the work`
- [x] 2.4 Add an end-to-end test that the notice names each bundled library and carries its licence text. Covers `An embedded library's licence travels with the work`

## 3. The bundle keeps what it is obliged to keep

`highlight.js` ships no legal comment of its own, so moving the bundler's setting is necessary and not sufficient: there has to be a notice for it to carry.
See the two decisions in `design.md` for why the notice is written in the client source rather than in the bundler, and why it is not the header the requirement above forbids.

- [x] 3.1 Add a test in `internal/assets` that the shipped bundle carries the copyright notice of every library bundled into it, and verify it fails against the bundle as committed today, for the right reason. Covers `The shipped bundle carries the notice of what it bundles`
- [x] 3.2 Exempt a comment citing `NOTICE` from the licence-header check, and verify a unit test that a borrowed notice passes while a header stating the work's own terms still fails. Covers `A borrowed work's notice is not the header this forbids`
- [x] 3.3 Open `web/src/surface.js` with a `/*!` comment naming `highlight.js`, its licence, its copyright line and `NOTICE` for the text, move `web/build.mjs` off `legalComments: 'none'` so notices are collected at the end of each output, run `make assets` to regenerate `internal/assets/bundle` and `internal/assets/bundle.sum`, and verify 3.1 passes and `make test` is green
- [x] 3.4 Mutate the guard: set `legalComments` back to `'none'`, regenerate, and verify the suite fails naming the missing notice, then restore. Covers `A bundler that deletes notices fails the suite`

## 4. A notice cannot be skipped by the next person

- [x] 4.1 Add a convention check that every file embedded through `internal/assets` is accounted for by an entry in `NOTICE`, reading what is embedded rather than what is in `web/`
- [x] 4.2 Add a unit test that the check fails for an embedded file with no entry and names it. Covers `An embedded asset with no entry fails the check`
- [x] 4.3 Add a unit test that a file present in `web/` and embedded in nothing needs no entry. Covers `An asset in the source but not in the binary needs no entry`

## 5. What a person reads

- [x] 5.1 Say in `README.md` which licence the work is published under and where the borrowed notices are, and verify the convention check on user-facing strings and the em dash check still pass
