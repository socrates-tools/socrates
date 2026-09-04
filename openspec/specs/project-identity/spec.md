# project-identity Specification

## Purpose
The name and the account the project answers to, where that name is allowed to be written, and what the normative context is required to describe.

## Requirements

### Requirement: The module path names the organisation

The module path SHALL be the project's repository under the organisation that publishes it, and every import within the project SHALL use it.

The path SHALL be written once as data where a build reads it, so that what a build stamps into itself and what the compiler resolves cannot disagree.

Renaming the module SHALL change where the source lives and nothing about what the tool is called: the binaries, the commands, the home directory and the configuration file are unaffected by it.

#### Scenario: Every import resolves under the organisation

- **WHEN** the module is built and its imports are read
- **THEN** the module path names the organisation's repository, and no import names any other module of this project

#### Scenario: The build reads the module path rather than repeating it

- **WHEN** the build stamps what a binary says about itself
- **THEN** it composes that from the module path stated once, and no second copy of the path exists for it to disagree with

### Requirement: Nothing in the source names another account

A check SHALL fail while any source file or live specification names a code-hosting account other than the organisation's.

The check SHALL be on what account is named rather than on which account was left behind, because a check for one departed name passes on the day a third name arrives.

Archived changes SHALL be exempt.
An archived change records a decision as it was taken, and the account it names is the account that existed then.
Rewriting it would make a past document claim something its author did not write, which a reader cannot distinguish from the truth.

#### Scenario: A source file naming another account fails the check

- **WHEN** a source file or a live specification names a code-hosting account that is not the organisation's
- **THEN** the check fails and names the file

#### Scenario: An archived change naming the former account passes

- **WHEN** an archived change names the account the project was published under before the move
- **THEN** the check passes, because the archive is the record of what was decided at the time

### Requirement: The normative context describes the project that exists

The context in `openspec/config.yaml` SHALL name the binaries, commands and directories the repository actually has.

It SHALL cite no source a reader cannot reach.
A normative document is one that can be checked against.
A citation nobody outside the author's machine can open cannot be checked, which is worse than none at all because it looks like evidence.

`openspec/` SHALL be stated as the sole normative source for the project's design.

#### Scenario: The context names the binary the repository builds

- **WHEN** the normative context is read
- **THEN** the binary it names and the command directory it names are the ones the repository builds

#### Scenario: The context cites nothing a reader cannot reach

- **WHEN** the normative context is read
- **THEN** it cites no private or unreachable source as the origin of the design, and states that `openspec/` is normative
