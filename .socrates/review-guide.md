# What a reviewer should know about socrates-review

The deterministic gates own formatting, linting, comments in source, em dashes, user-facing strings living in a package's `messages.go`, tests that read implementation source, and whether every specification scenario is named by a test.
`make verify` fails on all of those, so a finding about one of them is a finding about something already caught.

What follows is what nothing here checks.

## The scope rule

socrates-review reviews.
It does not land, verify, render, serve, or wait for a person.
Committing, pushing, opening a pull request and watching continuous integration are not socrates-review's and never become socrates-review's.
Running the reviewed project's own tests is not in the core layer either, so `clean` means an independent reviewer found nothing blocking, not that the tests pass.

A change that reaches into any of those is out of scope however well it is written, and that is worth reporting.

## The reasoning does not live in the code

`CONTEXT.md` defines the vocabulary, `docs/adr/` holds the decisions with what each one rejected, and `AGENTS.md` states the invariants.
A change that alters behaviour and leaves those saying the old thing has left the record wrong, and the code is not where the correction belongs.

A constraint somebody would otherwise re-derive belongs in `docs/gotchas.md`.
Reasoning about a change belongs in its commit message.

## A guard nobody has broken on purpose is a guard nobody knows works

A test added in the same breath as the code it covers is a test nobody has seen fail for the right reason.
Where a change adds a check, the question worth asking is whether anything shows it failing without the code it guards, not whether it passes now.

## The asymmetry is the whole design

The author has the context and does the patching.
The reviewer has independence and never authored anything.
A change that gives the reviewer a way to write, or gives the author's account of the change a channel to the reviewer, spends the one thing the rest of the design is paying for.
