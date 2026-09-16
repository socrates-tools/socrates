# Contributing to socrates

## A change starts as a proposal

Read this before you write anything.

A change here starts as a spec, published as a GitHub issue in this repository and labelled `ready-for-agent`.

The spec says, in the words of somebody using socrates review rather than the words of the code:

- the problem, and the solution to it.
- the user stories, at length rather than in summary.
- the implementation decisions, including which modules change and how their interfaces change.
- the testing decisions, including which seams are used and what prior art in the suite they follow.
- what is out of scope.

Two other places hold what a spec assumes rather than what it changes:

- `CONTEXT.md` holds the vocabulary. A spec names a domain concept in the words defined there, and a concept with no word yet is either language this project does not use or a gap worth recording.
- `docs/adr/` holds the decisions, each with what was rejected and why. A change that contradicts one says so rather than overriding it quietly.

A pull request that edits code with no spec behind it will be turned away, and that is not a formality.
Nothing enforces it mechanically. This repository held it with two gates that read a specification directory, and holds it by whoever reads the pull request instead.

If you are not sure whether something needs a spec, open an issue and say what you want to change.
A bug fix that restores stated behaviour does not need one. Anything that changes what socrates review promises does.

## Out of scope

`README.md` has a section called "Out of scope", and it is the shortest way to learn what a proposal should not ask for.
review reviews. It does not land, verify, render, serve or wait for a person about a review.
The one command that waits for a person is `socrates review configure`, and what it waits about is review's own settings.
Committing, pushing, opening a pull request and watching continuous integration are not review's and never become review's.

## The toolchain

- Go 1.26 or later.
- `git`.
- A Chrome or Chromium on your `PATH`, or named by `AGORA_BROWSER`. The page half of the suite drives a real browser.
- Node and `pnpm`, only if you touch `web/`. `make assets` rebuilds the committed bundle.
- The `claude` CLI, only if you re-record the simulator's fixtures or run the live smoke suite. Nothing else reaches it, and nothing in continuous integration does.

## Verifying a change

```sh
make fmt
make verify
```

`make verify` is the one entry point and it runs every gate.
Continuous integration runs that same command on every supported platform, so passing locally and passing there are one claim rather than two.
Do not add a gate to a workflow: add it to `make verify`, and a check fails while a workflow runs a gate of its own.

`make help` lists every target.

## Conventions

Each of these is checked by a test, so it is cheaper to follow than to argue with:

- **No comments in the source.** Not few, not only where necessary, none. A name that needs a sentence beside it is the wrong name. Reasoning belongs in the commit message, a constraint somebody would re-derive belongs in `docs/gotchas.md`, and a decision belongs in `docs/adr/`.
- **No em dash**, anywhere: prose, documents, commit messages, emitted strings.
- **One sentence per line** in long Markdown.
- **Every user-facing string lives in its package's `messages.go`**, names what the person can do next, and carries no em dash.
- **Tests name their scenario.** Prefer a test name containing the specification's scenario text, and use a `// scenario:` marker where the honest name would not.
- **Test first.** A task is not done until a test named for the scenario failed for the right reason and then passed.
- **Mutate a guard you rely on.** A check nobody has broken on purpose is a check nobody knows works.
- **Commit messages say what the change is.** No rationale essay, no summary of the conversation that produced it, no attribution trailer.

`AGENTS.md` is the map: the layout, the verification sequence and the invariants that hold the design up. Read it once before your first change.

## Dependencies

Every action a workflow uses is pinned to a commit, the linter and the releaser are pinned to a version in the `Makefile`, and Dependabot watches the Go modules, the npm packages and the actions.

That produces a steady queue of small pull requests, each bumping one pin.
A long queue is the mechanism working rather than a sign of neglect, and reviewing one is usually reading a changelog and letting `make verify` answer the rest.

## Opening a pull request

Branch, commit, and open a pull request that names the proposal it implements.
The template asks for that first, because a reviewer who knows which proposal a change serves can read it against something.
