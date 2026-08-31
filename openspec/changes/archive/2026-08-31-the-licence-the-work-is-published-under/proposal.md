## Why

socrates is about to be published, and nothing in the repository says what anybody may do with it.
A work with no stated terms is exclusively copyrighted by default, so a reader who wants to use it, fork it or vendor it has to ask, and most will not ask.
The terms are MIT, and the only reason nobody can act on that is that it is written down nowhere.

The harder half is what socrates does not own.
The four IBM Plex faces under `web/fonts` are compiled into the binary through `internal/assets`, and their licence requires the notice to travel with the copy.
`web/build.mjs` sets `legalComments: 'none'`, which deletes highlight.js's notice from the bundle that ships, so the binary today distributes a BSD-3-Clause work with the one thing that licence asks for stripped out.
Neither is a documentation gap.
Both are the distribution failing a condition of the licence it distributes under, and publishing multiplies the number of people it happens to.

The stated terms and the borrowed notices are the same requirement seen from two sides: what socrates grants, and what socrates was granted.
Stating one without the other publishes a permission the project does not entirely hold.

## What Changes

- The repository SHALL state the terms it is published under, in a file at its root, under the MIT licence.
- Everything embedded in a shipped binary SHALL carry the notice its own licence requires, in a file at the root that names each borrowed work, its licence and its full text.
- `web/build.mjs` SHALL stop deleting legal comments, so the bundle that is embedded and served keeps the notices of what it bundles, and a test SHALL assert those notices are in the shipped bundle rather than only in its source.
- The terms SHALL be stated once, at the root, and SHALL NOT be repeated as a per-file header, because a header on every file is a comment on every file and the project writes none.
- A new embedded asset SHALL NOT be addable without an entry naming its licence, held by a check in `internal/convention` rather than by care.

## Capabilities

### New Capabilities

- `project-licence`: the terms the work is published under, and the notices the work carries for everything it embeds.

### Modified Capabilities

None.

## Impact

- `LICENSE`: new, MIT, `Copyright (c) 2026 Rui Cheng`.
- `NOTICE`: new, naming IBM Plex (OFL-1.1), highlight.js (BSD-3-Clause) and every module dependency whose licence requires attribution, each with its full text.
- `web/build.mjs`: `legalComments` moves off `none`, which rewrites `internal/assets/bundle` and `internal/assets/bundle.sum`.
- `internal/assets/bundle_test.go`: a scenario that the shipped bundle carries the notice of what it bundles.
- `internal/convention`: a check that every file embedded through `internal/assets` is named in `NOTICE`, and a check that no source file carries a licence header.
- `README.md`: a line naming the licence, and a line pointing at `NOTICE`.
- Nothing reaches the reviewer, the run, the dialogue, the surface or the page.
