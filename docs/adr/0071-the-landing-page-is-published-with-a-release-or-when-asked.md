# 71. The landing page is published with a release or when asked

- Status: accepted
- Date: 2026-10-15
- Part of: site

## Decision

`site/` holds the landing page GitHub Pages serves: one HTML page, its stylesheet and a script that copies the install command, written by hand with no framework and no build step.
`make site` assembles it into `bin/site` beside the fonts under `agora/web/fonts`, the mark under `agora/web/brand` and `NOTICE`, so the page paints with the bytes the decisions page ships and no second copy of them is committed.

`.github/workflows/pages.yml` publishes what `make site` assembled, when somebody asks for it and as a job of every release that needs the verification matrix and the code scan, like every other job that publishes.
A release refreshes the page so it describes what a person installing gets, and a change to the page alone is published by asking for a run.

The page installs with the command `README.md` states and states the toolchain floor `go.mod` requires, and the suite fails while either disagrees.

## Rejected

- Serving the page from a branch or from `docs/` with Pages' own build.
  It publishes on every push, which is the trigger ADR 64 keeps to Scorecard alone, and runs Jekyll over a folder that holds decisions rather than a page.
- Publishing on every push to `main` that touches `site/`.
  The same trigger by another route, and the page would then describe commits nobody has released.
- Shipping the design's bundled export as the page.
  It carries a React runtime and its assets inline to draw one static page, and nobody can review a change to it.
