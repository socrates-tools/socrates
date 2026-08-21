## MODIFIED Requirements

### Requirement: The binary is `scrat` and it does one thing

socrates SHALL ship exactly one binary, named `scrat`.
It SHALL provide `review`, `status`, `findings`, `respond`, `continue`, `watch`, `logs`, `cancel` and `skill`, and it SHALL NOT provide a command that commits, pushes, opens a pull request, watches continuous integration, or runs the repository's tests.

The scope rule is mechanical rather than a judgement call about what is close enough to review, because judgement calls erode.
`respond` and `continue` record and advance; `skill` writes at user level. None of them writes in the repository, and adding them does not move the boundary.

#### Scenario: The command surface is what it says

- **WHEN** `scrat --help` is run
- **THEN** it lists `review`, `status`, `findings`, `respond`, `continue`, `watch`, `logs`, `cancel` and `skill`, and nothing that lands or verifies a change

#### Scenario: Recording a response writes nothing in the repository

- **WHEN** `scrat respond` and `scrat continue` are run
- **THEN** nothing inside the repository has changed except the git objects the new round's snapshot needed

#### Scenario: Installing the skill writes nothing in the repository

- **WHEN** `scrat skill` is run with its default target
- **THEN** the skill is written at user level and nothing inside the repository has changed
