## Context

See `proposal.md` for why.

`.github/workflows/ci.yml` runs on `push` to `main`, on `pull_request`, on `workflow_dispatch` and on `workflow_call`.
It reads the matrix from `platforms.json`, installs a browser at a stated version, and runs `make verify`.
`.github/workflows/release.yml` reaches it through `workflow_call`, and every publishing job needs that call to have passed.

`.github/workflows/vulnerabilities.yml` runs `govulncheck` on `push`, on `pull_request`, on a daily schedule and on `workflow_dispatch`.

`continuous-verification` requires the module graph to be checked "on every verification run and on a schedule".
Nothing in it requires a run to happen on a push, which is why only that one sentence is amended here.

## Goals / Non-Goals

**Goals:**

- Nothing starts a run by itself.
- A release still gates on the full matrix, because a release that nobody verified is the thing the release specification exists to prevent.
- What a run proves is untouched, so turning the triggers back on restores exactly what was there.

**Non-Goals:**

- Trimming the matrix. See `proposal.md`: a supported platform the suite does not run on is what `platform-support` exists to prevent, and dropping the macOS leg would take macOS out of the supported set.
- Deleting the workflows. The release calls one of them, and a trigger line is what "for now" means.
- Moving to a self-hosted runner. It replaces a bill with a machine to maintain and a new trust boundary on the thing that verifies the tool.
- Making the repository public to get free minutes. That is a decision about publishing, taken on its own terms, not a way to pay for continuous integration.

## Decisions

### The trigger goes, the workflow stays

`ci.yml` keeps `workflow_dispatch` and `workflow_call` and loses `push` and `pull_request`.

Deleting the file would take the release's gate with it, and would take the checks that read it: that a run invokes the one entry point, that the matrix is read from the statement, that the browser is installed at a stated version.
Those checks read the workflow rather than its schedule, so they keep holding, and they keep holding for the run a release performs.

Restoring the old behaviour is two lines in one file, which is what makes this reversible rather than a decision to unmake later.

### The vulnerability check moves onto the same footing

The daily schedule is the cheapest thing here, about one billed minute, and it is dropped anyway.

Leaving one workflow starting by itself while the reason for this change is that runs start by themselves would be a rule with an exception nobody would remember. The check is dispatchable, and a person who wants the answer can have it in a minute.

This is the requirement being amended rather than reinterpreted, because the sentence it replaces exists precisely to say that an advisory arrives without anybody looking, and after this change it does not.

### What is given up is written down, not absorbed

Three things stop being true, and each is stated in the specification rather than left for somebody to discover:

- A commit on `main` can be broken and nothing says so.
- An advisory published against the module graph waits for somebody to ask.
- A pull request arrives with no evidence attached to it.

The mitigation is the one this repository already relies on: `make verify` is one entry point, a contributor runs it before reporting work done, and it is the same command a run invokes. That was already the arrangement; what changes is that it is now the only one.

## Risks / Trade-offs

**A suite nobody runs rots.**
The mitigation is that a release runs it in full on every supported platform, so the longest anything can hide is one release.
That is weaker than a run per push and it is the trade being made.

**"For now" is the kind of thing that becomes permanent.**
Nothing here expires on its own.
What holds it is that the change is small, the specification says plainly what was given up, and turning it back on is two lines.

**A person may read a dispatchable workflow as a disabled one.**
It is not disabled: a release runs it, and it is the gate that decides whether a release publishes.
