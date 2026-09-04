## Context

See `proposal.md` for why.

`the-name-the-project-answers-to` moved the module path once already, and its non-goals said the binaries, the commands, the home directory and the configuration file were unaffected.
This change does the thing that change declined to do, and it does it for a reason that did not exist then: the account the project is published under now holds a brand rather than a project.

The module is entirely `internal/`, so nothing outside this repository imports any of it and no consumer breaks.
One release exists, `v0.1.0`, and its only user is the author.

What makes this delicate is size rather than risk.
The product's name appears 3284 times across 493 files, and `openspec/specs` names it in 90 of its 180 requirements.
A diff in which one file changes behaviour and 492 change a noun is a diff where the behaviour change is invisible, which is the failure socrates exists to prevent.

## Goals / Non-Goals

**Goals:**

- One name for the tool, reachable by typing it, and no second answer to what it is called.
- A home layout that the brand's second tool can adopt without this tool's directory moving.
- A delta that a reader can read: the requirements that changed are visible rather than buried among substitutions.
- Nothing in the source names an account or a repository twice, so the next move is a check failing rather than a grep somebody remembers to run.

**Non-Goals:**

- Renaming `agora`, or moving it to a repository of its own.
  It ships with this tool and is installed by it, and nothing about that changed.
- A second tool under the brand.
  This change makes room for one and builds none.
- A configuration file shared between tools.
  There is one tool, and a shared parser is a decision to take when there is something to share.
- Migrating an existing `~/.socrates`.
  See the decision below.
- Rewriting `openspec/changes/archive`.
- Renaming the tool where the specifications name it as the subject of a sentence rather than as a command.

## Decisions

### The tool is named for what it does, and the brand keeps the short word

`socrates` becomes the organisation and stays the brand.
`socrates-review` becomes the tool, the binary, the repository and the skill.

The alternative was to leave the binary `socrates` and rename only the repository, which is what the earlier module move did.
It fails on the second tool: two binaries published by `socrates-tools`, one of them called `socrates`, leaves a person unable to say which one they mean, and leaves the second tool's documentation explaining that the first tool is not the brand.

The cost is a longer command.
That cost is paid once per invocation by a session that writes it into a script, and the skill writes every invocation for the session anyway.

### `scrat` is removed rather than renamed

There is no short name.

`scrat` existed because `socrates` is long to type, and the machinery that came with it is not small: a second build in the releaser, a second name in the archive, `convention.ShortName` written as a rune slice so the check for it does not match itself, the rule that what the tool prints names the full name whichever name it was invoked as, and the rule that the skill writes the full command and never the short one.

Renaming it to something shorter than `socrates-review` reintroduces every one of those and reintroduces the question the brand collision is about: which of the two words is the tool called.

A person who wants a short name has `alias`, and an alias they wrote themselves is one they can spell however they like without the tool having to reconcile two.

### The home is a brand root holding one directory per tool

`SOCRATES_HOME` keeps its name and changes what it names: the brand's root rather than this tool's home.
The tool's own home is `review` beneath it, so what was `~/.socrates/config.yaml` is `~/.socrates/review/config.yaml`.

The tool's directory is the part of its name the brand does not supply.
`~/.socrates/socrates-review` spells the brand twice on one path, and every path a person reads or types carries that stutter for the life of the tool.

Everything below that is unchanged.
`internal/state` is still the one package that resolves the home, `agora` still asks it rather than resolving one of its own, and the runs, the surface's directory and the configuration file still sit under one home.
The requirements in `run-store`, `review-surface` and `agent-seam` say "under socrates' home" and stay true without being touched, which is the property that made stating the home in one package worth doing.

### Nothing is carried forward

No migration reads the old layout, no deprecation warns about it, and no compatibility path accepts the old name.

The tool has one user, who is told to remove `~/.socrates` and reinstall.
Code that migrates a layout is code that has to be correct about a state that will exist on exactly one machine for exactly as long as it takes to run one command, and it then stays in the tool forever, because nothing tells you when the last machine has moved.

### `.socrates/review-guide.md` in a repository under review is unchanged

The guide socrates reads out of a reviewed repository keeps its path.

The same rule applies there as applies to `~/.socrates`: `.socrates/` is the brand's directory, and `review-guide.md` names the tool within it.
It was already shaped that way, so it needs no change to be consistent, and changing it would break the guide in every repository that has one to express a rule it already expresses.

### The specifications rename commands, not the subject

A requirement is restated in a delta where it names a command a person types.
Where it names the tool as the subject of a sentence, it is left alone.

Renaming the subject as well reaches 90 of 180 requirements and 2952 lines, against 32 requirements and 948 lines for the commands alone.
The requirements that actually change, six of them, are readable in the smaller delta and are not in the larger one.

What this accepts is that `openspec/specs` goes on calling the tool `socrates` in prose.
Inside its own repository there is nothing for that to collide with: the ambiguity this change removes is between the brand and the tool as published, and the specifications describe one tool.
A later change that gives the brand a second tool's specifications in this repository is the change that has to revisit it, and this decision is where it will look.

### The move and the rename ride together

One change does both.

Doing them apart means publishing `socrates-tools/socrates` first, which is the exact collision the rename exists to remove, and it means two rounds of the same 493-file substitution.

Within the change they are still separable to read: the identity is three constants and a module path, and the rename is everything else.

### The archive is left alone

Two archived documents name the module under the previous account, and they keep naming it.

An archived change records a decision as it was taken.
The account and the binary it names are the ones that existed then, and rewriting it produces a document that claims a past author wrote something they did not.

`internal/convention` already skips `openspec/changes/archive` for the account check, and the same exemption covers this.

### The generic validator does not read scenario renames

Six requirements rename a scenario, because the scenario named a command that no longer exists or a behaviour that no longer happens.

`make scenario-drop`, which is this repository's gate and runs under `make verify`, reads `## RENAMED` entries for scenarios and passes.
`openspec validate` reads them for requirements only, and reports each renamed scenario as one the delta drops.

The gate is the one in `make verify`, and the archive step in this repository moves the change directory rather than running `openspec archive`, so nothing in the flow acts on the generic validator's report.
This is written down so the next person to see those six errors knows they were seen.
