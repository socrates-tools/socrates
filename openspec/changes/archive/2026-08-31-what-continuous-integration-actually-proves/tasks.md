**The suite must be green after every task, not only at the end of a group.**
The ordering below exists for that reason: the single entry point arrives before the workflow calls it, and the pin check is written before the pins it describes.

**This change follows `the-platforms-the-seal-holds-on`.**
The matrix is the supported set. A task here that writes a platform list into a workflow is this change making the third copy the dependency exists to prevent.

**A task about a workflow run says record what it proved.**
No task asserts a run duration, a runner image version, or a particular advisory. What is recorded is that the run happened, on which platforms, and that its gates passed.

## 1. One entry point

- [x] 1.1 Add `verify` to the `Makefile`, running the sequence `AGENTS.md` lists, with `lint` first and `gofmt` first within it, and verify it passes on a clean tree
- [x] 1.2 Add a test that verification fails on formatting before the end-to-end suite runs, by running it against an unformatted tree. Covers `Formatting fails before the suite does`
- [x] 1.3 Point `.github/workflows/ci.yml` at `make verify` instead of listing the targets, and verify the run is green. Covers `One command runs every gate`

## 2. The matrix is the supported set

- [x] 2.1 Derive the workflow matrix from the stated supported set rather than writing platforms into the workflow, and verify a run happens on each
- [x] 2.2 Add a test that a platform in the matrix and not in the supported set fails, and verify it fails when one is added. Covers `A platform in the matrix but not in the supported set fails`
- [x] 2.3 Run the full suite on each supported platform and record which platforms ran and that every gate passed. Covers `The suite runs on each supported platform`
- [x] 2.4 Fix whatever the second platform finds, or record it as a defect with the failing gate named, before this group is done

## 3. What a workflow may do, and what it is

- [x] 3.1 Add a convention check that no `uses:` in any workflow names an action by a tag, and verify it fails against the workflows as they stand. Covers `An action named by a tag fails the check`
- [x] 3.2 Pin every action to a commit with its version in a trailing comment, and verify 3.1 passes
- [x] 3.3 Add `permissions: contents: read` at the top of every workflow, escalating only in a job that needs more, and verify the run is green. Covers `Verification runs read-only`
- [x] 3.4 Add a concurrency group so a superseded run is cancelled, and verify a second push cancels the first run. Covers `A superseded run does not report`
- [x] 3.5 Add `.github/dependabot.yml` for Go modules, npm and Actions, and record the first change it opens. Covers `A dependency update arrives as a change somebody reviews`

## 4. What the run checks and what it must keep

- [x] 4.1 Add a vulnerability check over the module graph to the run and to a schedule, and record what it reported on the graph as it stands. Covers `A known vulnerability in a dependency is reported`
- [x] 4.2 Add `.github/workflows/codeql.yml` on a schedule, and record its first result
- [x] 4.3 Install the browser at a known version in the run and point `AGORA_BROWSER` at it, rather than resolving one from the runner image, and verify the page half of the suite still passes. Covers `The page half runs against a known browser`
- [x] 4.4 Mutate the guard: make no browser available and verify the run fails rather than skipping, then restore
- [x] 4.5 Verify no workflow compiles or runs anything behind the `record`, `drift` or `livesmoke` tags that reaches the real CLI, and record what was checked. Covers `A verification run reaches no reviewer CLI`

## 5. The linter widens

- [x] 5.1 Add `.golangci.yaml` enabling the analysers whose findings this project would act on, keeping `gofmt` as the first gate, and verify it passes against the tree as it stands
- [x] 5.2 Widen `make lint` to run it, and verify `make verify` is green on every supported platform
