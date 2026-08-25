## RENAMED Requirements

- FROM: `### Requirement: The binary is `scrat` and it does one thing`
- TO: `### Requirement: The binary is `socrates`, the short name is `scrat`, and it does one thing`

## MODIFIED Requirements

### Requirement: The binary is `socrates`, the short name is `scrat`, and it does one thing

socrates SHALL ship exactly one binary, named `socrates` after the product.
Installing it SHALL put `scrat` beside it as a second name for that same binary, so both are commands a person can type and neither is a shell alias they have to write themselves.

One binary under two names rather than two binaries, because two things built from the same source are two things that can be at different versions on one machine, and the question `socrates version` exists to answer is exactly that one.

It SHALL provide `review`, `status`, `findings`, `respond`, `continue`, `watch`, `logs`, `cancel`, `skill` and `version`, and it SHALL NOT provide a command that commits, pushes, opens a pull request, watches continuous integration, or runs the repository's tests.

The scope rule is mechanical rather than a judgement call about what is close enough to review, because judgement calls erode.
`respond` and `continue` record and advance; `skill` writes at user level; `version` reads what the binary records about itself and reaches nothing.
None of them writes in the repository, and adding them does not move the boundary.

#### Scenario: The command surface is what it says

- **WHEN** `socrates --help` is run
- **THEN** it lists `review`, `status`, `findings`, `respond`, `continue`, `watch`, `logs`, `cancel`, `skill` and `version`, and nothing that lands or verifies a change

#### Scenario: Both names are installed and reach the same binary

- **WHEN** socrates is installed and both `socrates version` and `scrat version` are run
- **THEN** each names the same commit, and the short name is a second name for the one binary rather than a second binary

#### Scenario: Recording a response writes nothing in the repository

- **WHEN** `socrates respond` and `socrates continue` are run
- **THEN** nothing inside the repository has changed except the git objects the new round's snapshot needed

#### Scenario: Installing the skill writes nothing in the repository

- **WHEN** `socrates skill` is run with its default target
- **THEN** the skill is written at user level and nothing inside the repository has changed

## ADDED Requirements

### Requirement: What socrates prints names socrates

Every string socrates prints SHALL name `socrates`, including the prefix on a failure, the command a refusal tells the person to run next, and the step a report says to take.

It SHALL name `socrates` when it was invoked as `scrat`, because a tool that reports itself by whichever name it was invoked as has no name, and a person reading `scrat: run scrat --help` cannot tell which of the two the documentation is about.

#### Scenario: The short name reports the full one

- **WHEN** the binary is invoked as `scrat` and refuses something
- **THEN** what it prints names `socrates`, and no string it prints names the short name
