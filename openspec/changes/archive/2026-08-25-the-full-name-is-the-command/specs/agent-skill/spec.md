## ADDED Requirements

### Requirement: The skill is named for the product

The skill SHALL be named `socrates`, SHALL be installed at `.claude/skills/socrates/`, and SHALL name `socrates` in every invocation its steps give.

Exactly one constant in the source SHALL carry that name, and it SHALL be the one the installer writes the directory from and the one the surface check reads invocations with.
Two constants would let the skill teach a command name that the check never looks for, and the check that proves every named command exists would go on passing while naming nothing.

#### Scenario: The skill is installed under the product's name

- **WHEN** `socrates skill` is run
- **THEN** the skill is written at `.claude/skills/socrates/SKILL.md`

#### Scenario: One constant names the skill and reads its invocations

- **WHEN** the constant that names the skill is changed
- **THEN** the installer's directory and the invocations the surface check reads both follow it, and no second copy of the name remains

#### Scenario: The skill names the full command

- **WHEN** the skill's steps are read
- **THEN** every invocation in them names `socrates`, and none names the short name

### Requirement: The skill writes an invocation as code

The skill SHALL write every invocation it names as code, inline or as an indented block, and the check that reads its invocations SHALL read code alone.

The command and the product are now one word, so a sentence whose subject is the product is indistinguishable from an invocation by any rule that reads prose.
Read as prose, `socrates patches nothing` names a command called `patches`, and the check that exists to prove every named command is real would report five commands the binary has never had.

#### Scenario: Prose naming the product is not an invocation

- **WHEN** the skill states a sentence whose subject is the product rather than a command
- **THEN** the check reads no command from it, and every invocation written as code is still read
