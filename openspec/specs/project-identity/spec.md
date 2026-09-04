# project-identity Specification

## Purpose
The name and the account the project answers to, where that name is allowed to be written, and what the normative context is required to describe.

## Requirements

### Requirement: The module path names the organisation

The module path SHALL be the project's repository under the organisation that publishes it, and every import within the project SHALL use it.

The path SHALL be written once as data where a build reads it, so that what a build stamps into itself and what the compiler resolves cannot disagree.

The repository SHALL be named for the tool it holds, and the organisation SHALL be named for the brand that publishes it.
A person holding the command holds the repository's name, and a person holding the repository's name knows which of the organisation's tools it is.
An organisation that publishes more than one tool and a repository named after the organisation leave the second tool with no name to take.

Renaming the module on its own SHALL change where the source lives and nothing about what the tool is called: the binaries, the commands, the home directory and the configuration file are unaffected by it.
Where the tool itself is renamed the repository SHALL follow it, and that is a change to both stated as one rather than a consequence of either.

#### Scenario: Every import resolves under the organisation

- **WHEN** the module is built and its imports are read
- **THEN** the module path names the organisation's repository, and no import names any other module of this project

#### Scenario: The build reads the module path rather than repeating it

- **WHEN** the build stamps what a binary says about itself
- **THEN** it composes that from the module path stated once, and no second copy of the path exists for it to disagree with

#### Scenario: The repository is named for the tool

- **WHEN** the last segment of the module path is compared with the command the tool installs
- **THEN** they are the same name, and the segment before it is the organisation rather than a repeat of that name

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
