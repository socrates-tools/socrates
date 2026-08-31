## Context

See `proposal.md` for why.

`.github/workflows/ci.yml` has one job on `ubuntu-latest`.
It checks out, sets up Go from `go.mod`, resolves a browser by looking for `google-chrome` or `chromium` on `PATH` and exporting `AGORA_BROWSER`, runs `make lint`, `make test`, `make scenario-coverage`, `make scenario-drop`, `make corpus-outcomes` and `make fixtures-diff`, then sets up Node and checks that `make assets` reproduces the committed bundle.

`make test` runs `go test -race ./...`, which includes `internal/e2e`, which drives real binaries, real git and a real browser.
So the suite is genuinely end to end, on one platform.

`make fixtures` and `make live-smoke` are behind build tags and are not in the workflow, which `agent-seam` requires.

## Goals / Non-Goals

**Goals:**

- A green run means the suite passed on every platform the project claims, not on one of them.
- What a workflow does changes only when this repository changes.
- A contributor and the workflow run the same thing, so "it passes locally" and "it passes in CI" stop being different claims.
- A dependency with a published advisory becomes visible without anybody looking.

**Non-Goals:**

- Windows in the matrix. `platform-support` states what is supported, and this change runs the suite there and nowhere else.
- Reaching the real CLI. That boundary exists, is held by build tags, and is inherited rather than restated.
- Line coverage as a gate. The release metric here is scenario coverage, which `make scenario-coverage` already enforces, and line coverage cannot see the bug class this project actually has.
- A merge queue, required reviews, or branch protection. Those are repository settings rather than specifications, and they are raised as an issue.

## Decisions

### The matrix is the supported set, read rather than written

The workflow's matrix comes from the supported platforms `platform-support` states.

Writing the platforms into the workflow would be a third copy of a set that already has to agree in two places, and `platform-support` requires those to agree in both directions.
A platform added to the matrix and not to the statement, or the reverse, fails the suite rather than passing quietly.

### Actions are pinned to commits, and a check holds it

Every `uses:` names a commit, with the human-readable version in a trailing comment.

A tag is a mutable pointer somebody else controls.
Pinning to a commit means the only way a workflow's behaviour changes is a commit in this repository, which is the same property the project already insists on for the reviewer's tool list and for the committed asset bundle.

`internal/convention` gains a check that no `uses:` names a tag, because a pin that is only a habit is a pin that lapses on the first hurried change.

Dependabot updates the pins, which is what makes pinning sustainable rather than a decision to fall behind.

### The token is read-only at the top

`permissions: contents: read` at workflow level, escalated per job only where a job needs more.

Nothing in verification needs to write.
The release workflow does, and it asks for what it needs in the job that needs it.
The default token permission is a repository setting, which means today the workflow's authority depends on a setting nobody stated in a file anybody reviews.

### The browser is pinned, and its absence is a failure rather than a skip

The current step fails when no browser is found, which is correct.
What it does not control is which browser.

A pinned browser is installed by the workflow rather than found on the runner, so a runner-image change cannot alter what the page half of the suite runs against, and cannot remove it.
`AGORA_BROWSER` already exists as the seam for this, so nothing in the suite changes.

### `make verify` is the single entry point

The Makefile gains `verify`, which runs the sequence `AGENTS.md` already lists, and the workflow calls it.

Today the workflow lists the targets and `AGENTS.md` lists them again, so a seventh gate is added in two places or one.
A contributor running one command and a workflow running the same command is also what makes `CONTRIBUTING.md` in `the-door-a-stranger-arrives-at` a short document instead of a list somebody has to keep current.

`make lint` stays the first thing `verify` runs, and `gofmt` stays the first thing `lint` runs, so the cheapest failure is still the fastest.

### golangci-lint with a committed configuration, starting narrow

The configuration is committed and starts with the analysers whose findings this project would act on, rather than with everything enabled and a list of suppressions.

An inherited default that fires on 37,000 lines of existing code produces a change whose diff is noise, and this repository's own rule is that a legacy file lightly touched does not get a broad refactor riding along.

## Risks / Trade-offs

**A matrix multiplies runtime and flakiness.**
`internal/e2e` spawns real processes and drives a real browser, so running it on a second platform doubles the surface on which a flake can appear.
That is the point: a flake that only appears on macOS is a bug this project cannot currently see.

**Pinned actions plus Dependabot is a steady stream of small pull requests.**
That is the intended cost of the guarantee, and `CONTRIBUTING.md` says so, so that a long queue reads as the mechanism working rather than as neglect.

**A wider linter will disagree with existing code.**
The configuration starts narrow for that reason, and widening it later is its own change rather than a rider on this one.

**The vulnerability check can fail on something with no fix available.**
A failing gate nobody can fix gets disabled, which is worse than no gate.
It runs on a schedule as well as on a run so that the finding arrives as a notification rather than only as a blocked pull request.
