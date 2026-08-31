## Why

Verification runs on every push, and the account cannot afford it.

The repository is private, and a private repository bills a macOS minute as ten.
Measured on a real run, one push to `main` costs about 123 billed minutes: roughly 110 for the macOS leg, 12 for the Linux leg, one for the vulnerability check.
The free allowance is 2,000 minutes a month, so sixteen pushes spend the month, and a day of ordinary work spent it.
What that looks like is not a slow queue: jobs stop starting altogether, with the release among them, so the mechanism that publishes nothing unless verification passed publishes nothing at all.

Nothing about what a run proves is wrong.
The matrix is read from the platform statement, the actions are pinned, the tokens are read-only, the browser is installed at a stated version, and the run invokes the one entry point a contributor invokes.
The only thing this change touches is when a run happens.

The alternative that keeps the pushes is to stop running the suite on macOS, and that is worse.
`platform-support` supports macOS on the basis that the suite runs there, so trimming the matrix would take macOS out of the supported set, and a supported platform nobody runs is exactly what that specification exists to prevent.
Running the whole suite less often keeps every platform honest; running less of the suite does not.

## What Changes

- A verification run SHALL happen when somebody asks for one and when a release gates on one, rather than on every push and every pull request.
- The check on the module graph SHALL run on the same footing, rather than on a schedule of its own.
- Nothing about what a run proves SHALL change: the same entry point, the same matrix read from the same statement, the same pins, the same tokens, the same browser and the same boundary.
- What is given up SHALL be stated rather than absorbed: nothing tells anybody that `main` is red, and nothing tells anybody that an advisory has been published, until somebody asks.

## Capabilities

### Modified Capabilities

- `continuous-verification`: when a run happens, and on what footing the check on the module graph runs.

### New Capabilities

None.
This change removes a trigger and states what that costs. It adds no mechanism.

## Impact

- `.github/workflows/ci.yml`: started on request and by a release, rather than on push and on pull request.
- `.github/workflows/vulnerabilities.yml`: the same, and no schedule.
- `internal/convention`: the check on when the module graph is read.
- `.github/workflows/release.yml`: unchanged. It calls the verification workflow, and calling it is unaffected by what else starts it.
- Nothing reaches the reviewer, the run, the dialogue, the surface or the page.
- No platform leaves the supported set, and `platforms.json` is untouched.

## Dependencies

This change SHALL follow `what-continuous-integration-actually-proves`, whose requirements it amends.
