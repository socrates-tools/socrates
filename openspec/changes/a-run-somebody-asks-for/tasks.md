**The suite must be green after every task, not only at the end of a group.**
The check on what starts a run is written before the triggers it describes are changed, so that changing them is what makes it pass.

**Nothing here touches what a run proves.**
The matrix, the entry point, the pins, the tokens, the browser and the build-tag boundary are out of scope. A task that edits one of them is in the wrong change.

**No platform leaves the supported set.**
`platform-support` supports a platform because the suite runs there. A task that drops the macOS leg to make a run cheaper is taking macOS out of the supported set, which is a different proposal.

## 1. What starts a run

- [x] 1.1 Add a convention check reading what starts each workflow, and verify it reports the workflows that start on a push as the tree stands

  It reported five: `ci.yml` on push and on pull request, `vulnerabilities.yml` on push, on pull request and on a schedule.
- [x] 1.2 Take `push` and `pull_request` off the verification workflow, leaving it startable on request and by a release. Covers `A run is started by a person or by a release`
- [x] 1.3 Take `push`, `pull_request` and the schedule off the check on the module graph, and amend the check that held it to the footing the specification now states. Covers `A known vulnerability in a dependency is reported`
- [x] 1.4 Verify the release still gates on the full matrix, by reading what it calls rather than by running one. Covers `A release still gates on the full matrix`
- [x] 1.5 Mutate the guard: put a push trigger back on a fabricated workflow and verify the check reports it, then restore

## 2. What still holds

- [x] 2.1 Verify the checks that read the verification workflow still hold: one entry point, the matrix read from the statement, a browser at a stated version, a read-only token, a superseded run cancelled
- [x] 2.2 Run `make verify` and record that it is green with nothing starting on a push
