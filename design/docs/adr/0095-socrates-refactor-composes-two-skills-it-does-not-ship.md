# 95. socrates-refactor composes two skills it does not ship

- Status: accepted
- Date: 2026-11-19
- Part of: socrates-refactor

## Decision

`socrates-refactor` is a third skill of `design/skill`, installed by `socrates design skill` beside `socrates-design` and `socrates-tickets`.
It runs `improve-codebase-architecture`, from `mattpocock/skills`, and puts its report and the pick of a candidate on a page `lavish`, from `kunchenguid/lavish-axi`, serves for the person to annotate.
The person asked for it in their own words: "i want to create an aggregated skill that invokes improve-codebase-architecture with lavish, so that i can get the value out of the both."

Two things are the surprise.
The skill installs either of them where it is missing, with `npx -y skills add` for the host it was rendered for, because the person asked that it "auto install them if any of them is missing".
And it reads `improve-codebase-architecture` from where it is installed rather than invoking it, because that skill is a person's to start and no session can invoke it by name.

It reaches for no other socrates skill, per ADR 49, and uses lavish rather than agora because the report is a free-form page of diagrams rather than a document of decisions.

## Rejected

- A namespace of its own, `socrates refactor`, with a command tree that installs one skill.
  A namespace that holds only a skill is what design already is, and an architecture review is a design activity.
- Copying both skills into socrates' own source.
  Each would drift from what its author publishes, and lavish says its own guidance lives in its CLI for that reason.
