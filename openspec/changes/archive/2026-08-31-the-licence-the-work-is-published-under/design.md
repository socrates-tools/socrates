## Context

See `proposal.md` for why.

Three things about the current shape decide most of this design.

`web/build.mjs` runs esbuild over `web/src/surface.js` and `web/src/surface.css` with `legalComments: 'none'`, writes the result into `internal/assets/bundle`, and writes a digest of its inputs to `internal/assets/bundle.sum`.
`internal/assets` embeds that directory, and `make test` fails while the committed bundle and its source disagree.
So a change to the bundler is a change to a generated file, and both have to move together.

`highlight.js` is the only runtime dependency in `web/package.json` and it is BSD-3-Clause.
It ships no legal comment of its own: nothing under its `lib` carries `/*!`, `//!`, `@license`, `@preserve` or the word copyright, so building with `legalComments: 'eof'` produces the same bytes as `none` and carries no notice either way.
This was measured rather than assumed, and it is why moving the setting is necessary and not sufficient.

The four `web/fonts/*.woff2` files are IBM Plex, under OFL-1.1.
They are binary, so no comment can live inside them, and the licence's requirement is that the notice accompanies the copy rather than that it be embedded in it.

## Goals / Non-Goals

**Goals:**

- One file states the terms, at the root, where every tool that looks for terms looks.
- One file carries every borrowed notice, so a person auditing what a socrates binary contains has one place to read.
- The bundle that ships carries what its own licence asks it to carry, proved by a test against the shipped bytes rather than against the source.
- Adding a fifth font or a second library cannot silently skip its notice.

**Non-Goals:**

- Per-file SPDX headers. See the decision below.
- A command that prints the licences. It is a real idea and it belongs to `distribution`, where the question of what a release archive contains is already being answered.
- Auditing the licences of Go module dependencies mechanically. The module graph is eleven direct and nine indirect entries, all permissive, and a generated dependency report is a tool this repository does not otherwise need.

## Decisions

### The terms are stated once, at the root, and never per file

The convention this repository holds is that source carries no comments.
A licence header is a comment, on every file, saying the same thing the root file says.

The counter-argument is real: a single file extracted from the repository loses its terms.
It is rejected because the file that gets extracted from a Go project is nearly always extracted with its module path, which resolves to the repository, which carries `LICENSE`.
The cost of the header is 262 files each opening with a comment block in a project whose stated rule is that a name that needs a sentence beside it is the wrong name, and whose reviewer prompt teaches the same.

`internal/convention` gains a check that no source file opens with a licence header, so the decision is enforced in the direction it can drift.

### `legalComments: 'eof'` rather than `'linked'` or `'external'`

esbuild's options put the comments inline where they were, at the end of the file, in a separate file, or nowhere.

`eof` is chosen because the bundle is embedded in a Go binary and served as one asset.
`external` and `linked` produce a second file the surface would have to serve and the notice would then be one request away from the code it belongs to, which is the failure mode the licence exists to prevent.
`inline` keeps every comment where it was, which is correct and larger; `eof` collects them at the end and is what a reader auditing the bundle would look for.

The CSS bundle gets the same treatment, and carries nothing today, which is the point of asserting on the shipped bytes rather than on the setting.

### The notice the bundle carries is written in the client source, not in the bundler

Since `highlight.js` ships no legal comment, something has to write the one the bundler carries.

esbuild offers two places.
A `banner` is inserted verbatim into the output and survives `legalComments: 'none'`, which was measured.
A `/*!` comment in `web/src/surface.js` is collected at the end of the bundle under `eof` and deleted under `none`, which was also measured.

The comment in the client source is chosen because it is the only one of the two the setting acts on.
With a banner, `eof` would be a setting nothing in the suite defends: somebody could set it back to `none` and every test would still pass, which is the check nobody has broken on purpose.
With the comment, the scenario that a bundler deleting notices fails the suite is a scenario a person can actually run, and it fails naming the notice that went missing.

The comment sits at the top of the file that imports `highlight.js`, which is where a reader looking for what the client pulls in already is.

### A borrowed work's notice is not the header the repository forbids

The decision above puts a comment carrying somebody else's copyright line at the top of a source file, which the check from the requirement above would otherwise report.

The two requirements are about different things.
One forbids restating socrates' own terms on every file, because the root file states them once.
The other requires a borrowed work's notice to travel with the copy, which means being in the copy.

The check tells them apart by whether the comment cites the notice file.
A header that states terms and cites nothing is the header that is forbidden; a comment that names a borrowed work and points at `NOTICE` for the text is the notice that is required.
That rule is mechanical, and it fails in the direction that matters: a licence header pasted on a file cites no notice file and is still reported.

### The notice file is `NOTICE` and it carries full text, not references

A file of links is not a notice.
OFL-1.1 requires the licence text to accompany the fonts, and a URL is not the text.

`NOTICE` rather than `THIRD_PARTY_LICENSES.md` because it is the name tooling recognises, it sits beside `LICENSE` alphabetically and visually, and it is not Markdown pretending a licence is prose.

### The check is on what is embedded, not on what is in `web/`

`internal/convention` gains a check that walks what `internal/assets` embeds and fails for any file whose name is not accounted for in `NOTICE`.
Checking `web/` instead would pass for a font that was added to the source tree and never bundled, and fail for one that was bundled from somewhere else.
What ships is the thing the licence attaches to.

The check names files rather than parsing licences, because a check that tried to decide whether a notice is legally sufficient would be a check nobody can trust and everybody would work around.
It answers one question, which is whether somebody wrote an entry, and leaves the entry's correctness to the person who wrote it.

## Risks / Trade-offs

**The bundle grows.**
Legal comments at the end of a minified bundle are a few kilobytes on an asset already measured in tens.
The alternative is shipping without them, which is the defect.

**`NOTICE` goes stale against the module graph.**
The convention check covers embedded assets, which is where the binary's obligations actually are, and does not cover Go dependencies.
A dependency whose licence changes under a version bump is not caught.
That residual is stated here rather than covered, and Dependabot's arrival in `what-continuous-integration-actually-proves` is what puts a person in front of each bump.

**Regenerating the bundle is a diff nobody reads.**
`make assets` rewrites minified output, so the review of this change is the review of `build.mjs` and the test, not of the bundle.
That is why the test asserts on the shipped bytes: it is the only part of the diff a reader can trust without reading minified JavaScript.
