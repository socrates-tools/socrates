# Contributing to socrates

## A change starts as a proposal, not as code

Read this before you write anything.

`openspec/` is normative in this repository.
`openspec/specs` says what is required and why, and the design document of each archived change says what was decided and what was rejected.
Where the source and `openspec/` disagree, `openspec/` is right and the source is the bug.

So a change here starts as a proposal in `openspec/changes/<name>/`, holding:

- `proposal.md`: why, what changes, which capabilities, what the impact is, what it depends on.
- `design.md`: the decisions, and the goals and non-goals. A non-goal here is load-bearing: it is how a later change proposing the thing is recognised as proposing it rather than as an oversight.
- `specs/<capability>/spec.md`: the requirements, each with the scenarios that hold it.
- `tasks.md`: the order, and the constraints particular to that change.

A pull request that edits code with no proposal behind it will be turned away, and that is not a formality.
Two gates enforce it mechanically: the scenario-coverage gate fails while a scenario in an underway change is named by no test, and the scenario-drop gate fails while archiving a change would take a scenario out of the promoted specifications.
Both run under the one verification command below.

If you are not sure whether something needs a proposal, open an issue and say what you want to change.
A bug fix that restores stated behaviour does not need one. Anything that changes what socrates promises does.

## What socrates will not become

`README.md` has a section called "What it deliberately is not", and it is the shortest way to learn what a proposal should not ask for.
socrates reviews. It does not land, verify, render, serve or wait for a person.
Committing, pushing, opening a pull request and watching continuous integration are not socrates' and never become socrates'.

## The toolchain

- Go 1.26 or later.
- `git`.
- A Chrome or Chromium on your `PATH`, or named by `AGORA_BROWSER`. The page half of the suite drives a real browser.
- Node and `pnpm`, only if you touch `web/`. `make assets` rebuilds the committed bundle.
- The `claude` CLI, only if you re-record the simulator's fixtures or run the live smoke suite. Nothing else reaches it, and nothing in continuous integration does.

## One command checks a change

```sh
make fmt
make verify
```

`make verify` is the one entry point and it runs every gate.
Continuous integration runs that same command on every supported platform, so passing locally and passing there are one claim rather than two.
Do not add a gate to a workflow: add it to `make verify`, and a check fails while a workflow runs a gate of its own.

`make help` lists every target.

## The conventions

Each of these is checked by a test, so it is cheaper to follow than to argue with:

- **No comments in the source.** Not few, not only where necessary, none. A name that needs a sentence beside it is the wrong name. Reasoning belongs in the commit message, a constraint somebody would re-derive belongs in `docs/gotchas.md`, and a decision belongs in `openspec/`.
- **No em dash**, anywhere: prose, documents, commit messages, emitted strings.
- **One sentence per line** in long Markdown.
- **Every user-facing string lives in its package's `messages.go`**, names what the person can do next, and carries no em dash.
- **Tests name their scenario.** Prefer a test name containing the specification's scenario text, and use a `// scenario:` marker where the honest name would not.
- **Test first.** A task is not done until a test named for the scenario failed for the right reason and then passed.
- **Mutate a guard you rely on.** A check nobody has broken on purpose is a check nobody knows works.
- **Commit messages say what the change is.** No rationale essay, no summary of the conversation that produced it, no attribution trailer.

`AGENTS.md` is the map: the layout, the verification sequence and the invariants that hold the design up. Read it once before your first change.

## Dependencies arrive as a stream of small changes

Every action a workflow uses is pinned to a commit, the linter and the releaser are pinned to a version in the `Makefile`, and Dependabot watches the Go modules, the npm packages and the actions.

That produces a steady queue of small pull requests, each bumping one pin.
A long queue is the mechanism working rather than a sign of neglect, and reviewing one is usually reading a changelog and letting `make verify` answer the rest.

## Opening the change

Branch, commit, and open a pull request that names the proposal it implements.
The template asks for that first, because a reviewer who knows which proposal a change serves can read it against something.
