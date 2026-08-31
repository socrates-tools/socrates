## Why

The verification workflow runs eight steps and every one of them is worth running.
What it does not do is run them anywhere but one Linux runner, on one Go version, with actions resolved from mutable tags, under a token whose permissions nobody stated.

A green badge is read as "this works".
What this one currently proves is "this worked once, on ubuntu-latest, on whatever Go `go.mod` asked for, against whatever Chrome that runner image happened to ship".
macOS is a first-class platform for a developer tool that spawns processes and manipulates permission bits, and it is exactly where the two unix mechanisms socrates depends on behave differently enough to matter.
Nothing in this repository has ever run there.

Three of the gaps are supply-chain rather than correctness.
Actions pinned to `@v7` resolve to whatever that tag points at today, so the workflow's behaviour can change without a commit in this repository.
No `permissions:` block means the workflow token gets the repository default rather than the least it needs.
Nothing checks the module graph for known vulnerabilities, and nothing updates it, so a dependency with a published advisory sits until somebody happens to look.

`make lint` is `gofmt -l` and `go vet`.
That is a real gate and a narrow one, and this project's own reviewer prompt holds contributors to a standard the linter does not check.

None of this is urgent while the only contributor is the author.
All of it is load-bearing the moment a pull request arrives from somebody whose machine nobody has seen.

## What Changes

- Verification SHALL run on every supported platform, rather than on one, so a platform is supported because the suite passed there and not because it compiled.
- Every action a workflow uses SHALL be pinned to an immutable commit, so that what a workflow does changes only when this repository changes.
- Every workflow SHALL state the least permission its token needs, at the top, and SHALL escalate only in the job that needs more.
- Superseded runs SHALL be cancelled rather than left to finish, so what is reported is what is on the branch now.
- The module graph SHALL be checked for known vulnerabilities on every run and on a schedule, and SHALL be kept current by something that opens a change rather than by somebody remembering.
- The linter SHALL be widened beyond formatting and vet, with its configuration committed, keeping the formatting gate as the first thing that fails.
- The browser the surface tests drive SHALL be resolved to a known version rather than to whatever the runner image ships, so that a runner-image change cannot silently remove the half of the suite that drives a real page.
- Nothing in continuous integration SHALL reach the real CLI. This is already true and this change must not breach it.

## Capabilities

### New Capabilities

- `continuous-verification`: what an automated verification run proves, on what platforms, under what permissions, and what it is not allowed to reach.

### Modified Capabilities

None.
`agent-seam` already requires that the paths reaching the real CLI are behind build tags so the ordinary suite compiles neither. This change depends on that and restates none of it.

## Impact

- `.github/workflows/ci.yml`: a matrix over the supported platforms, pinned actions, a permissions block, a concurrency group, a pinned browser, and a vulnerability check.
- `.github/workflows/codeql.yml`: new, on a schedule.
- `.github/dependabot.yml`: new, for Go modules, npm and Actions.
- `.golangci.yaml`: new, committed configuration.
- `Makefile`: `lint` widens; `verify` arrives as the single entry point the workflow and a contributor both use.
- `internal/convention`: a check that every action reference in a workflow is a commit rather than a tag.
- Nothing reaches the reviewer, the run, the dialogue, the surface or the page.

## Dependencies

This change SHALL follow `the-platforms-the-seal-holds-on`.
The verification matrix is the supported set, and running the suite on a platform nobody has decided to support proves something about a platform the project does not claim.
