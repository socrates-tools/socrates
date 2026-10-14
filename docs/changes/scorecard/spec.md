# OpenSSF Scorecard runs weekly and the README carries its badge

Decisions: ADR 64.

## Problem Statement

The repository holds most of what OpenSSF Scorecard measures: actions pinned by commit, read-only tokens, a vulnerability check, code scanning, signed releases and reviewed dependency updates.
Nobody outside the repository can see that without reading the workflows, and a setting that lives on the forge rather than in the tree, such as branch protection, can change without anything in the suite noticing.

## Solution

A workflow runs Scorecard once a week and on request, publishes the results, and uploads them to code scanning.
The README carries the Scorecard badge beside its title.
It is the one workflow that starts without anybody asking, and the suite fails where any other workflow does, or where this one starts on anything but its schedule and a request.

## User Stories

1. As a person deciding whether to install socrates, I want to see an independent score of how the repository is hardened, so that I do not have to read its workflows to judge it.
2. As the maintainer, I want a weekly reading, so that a setting changed on the forge shows up within a week rather than never.
3. As the maintainer, I want the results in code scanning, so that each finding is something I can read and dismiss in one place.
4. As the maintainer, I want to run it on request, so that I can see a fix take effect without waiting a week.
5. As a contributor, I want the suite to fail where a second workflow starts without being asked, or where Scorecard starts on a push, so that the exception stays one workflow wide.

## Implementation Decisions

- `.github/workflows/scorecard.yml`, holding what Scorecard allows a publishing workflow to hold: no top-level `env` or `defaults`, a read-only token at the top, and one job escalating to `security-events: write` and `id-token: write`, whose steps are checkout without persisted credentials, the Scorecard action and the SARIF upload.
- Every action is pinned by commit with its version beside it, like every other workflow.
- `on:` states a weekly `schedule` and `workflow_dispatch`, and a superseded run is cancelled.
- `internal/convention` states the one workflow allowed to start unasked and the one trigger it may start on.
- `AGENTS.md`'s invariant about when a run happens names the exception and ADR 64.
- `README.md` carries the badge, linking to the Scorecard viewer.

## Testing Decisions

- A convention test reads every workflow and fails where one other than Scorecard starts on `push`, `pull_request` or `schedule`, or where Scorecard starts on anything but `schedule` beside `workflow_dispatch`.
  It is mutated by adding a schedule to `ci.yml` and a push to `scorecard.yml`.
- The existing checks on pinning, read-only tokens and superseded runs read the new workflow without changing.
- That the workflow runs and publishes is checked live once, by starting it after it lands and recording the run.

## Out of Scope

- The OpenSSF Best Practices badge, which is a questionnaire answered by the maintainer at bestpractices.dev rather than something this repository builds.
- Acting on what the first reading finds, which is a change of its own for each finding worth acting on.
