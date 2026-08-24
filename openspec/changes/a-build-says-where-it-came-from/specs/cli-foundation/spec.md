## MODIFIED Requirements

### Requirement: The binary is `scrat` and it does one thing

socrates SHALL ship exactly one binary, named `scrat`.
It SHALL provide `review`, `status`, `findings`, `respond`, `continue`, `watch`, `logs`, `cancel`, `skill` and `version`, and it SHALL NOT provide a command that commits, pushes, opens a pull request, watches continuous integration, or runs the repository's tests.

The scope rule is mechanical rather than a judgement call about what is close enough to review, because judgement calls erode.
`respond` and `continue` record and advance; `skill` writes at user level; `version` reads what the binary records about itself and reaches nothing.
None of them writes in the repository, and adding them does not move the boundary.

#### Scenario: The command surface is what it says

- **WHEN** `scrat --help` is run
- **THEN** it lists `review`, `status`, `findings`, `respond`, `continue`, `watch`, `logs`, `cancel`, `skill` and `version`, and nothing that lands or verifies a change

#### Scenario: Recording a response writes nothing in the repository

- **WHEN** `scrat respond` and `scrat continue` are run
- **THEN** nothing inside the repository has changed except the git objects the new round's snapshot needed

#### Scenario: Installing the skill writes nothing in the repository

- **WHEN** `scrat skill` is run with its default target
- **THEN** the skill is written at user level and nothing inside the repository has changed

## ADDED Requirements

### Requirement: `scrat version` says which build this is

`scrat version` SHALL print the short commit the binary was built from, and SHALL say when the working tree it was built from held uncommitted changes.

The commit SHALL be recorded in the binary when it is linked, and SHALL NOT be read at run time from the repository the command is invoked in.
A binary that reported the tree it is standing in would answer the question this command exists for with the one answer that cannot be checked, since `scrat` is run against repositories it was not built from.

A binary carrying no recorded commit SHALL say so and SHALL say what to build it with.
It SHALL NOT print an empty version and SHALL NOT infer one, because a version line nobody can trust is worse than one that admits it is missing.

`scrat version` SHALL take no arguments, and SHALL refuse one with the usage code.

#### Scenario: The build says which commit it came from

- **WHEN** `scrat` is built through the Makefile and `scrat version` is run
- **THEN** it prints the short commit of the checkout it was built from

#### Scenario: A build from a tree with uncommitted changes says so

- **WHEN** `scrat version` is run on a binary linked from a working tree that held uncommitted changes
- **THEN** the line says the tree it was built from held them, beside the commit

#### Scenario: A build that recorded no commit says so

- **WHEN** `scrat version` is run on a binary that was linked without the stamp
- **THEN** it says no commit was recorded and says what to build it with, and prints no commit of its own

#### Scenario: `scrat version` takes no arguments

- **WHEN** `scrat version` is given an argument
- **THEN** it refuses, names the argument, and exits with the usage code
