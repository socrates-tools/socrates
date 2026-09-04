## RENAMED Requirements

- FROM: `### Requirement: `scrat skill` installs it and says where`
- TO: `### Requirement: `socrates-review skill` installs it and says where`

- FROM: `#### Scenario: The skill names the full command`
- TO: `#### Scenario: The skill names the one name the tool has`

## MODIFIED Requirements

### Requirement: The skill teaches polling rather than following

The skill SHALL name `socrates-review status` as the command to wait with, and SHALL say that everything the session reads lands in its own context.

It SHALL name `socrates-review watch` only as something to suggest to the person, in another terminal, and SHALL NOT put it in the session's own loop.

It SHALL name how often to poll, and SHALL say that the wait is bounded by socrates rather than by the session.

A session given no cadence and no bound invents both, and each way of inventing them is wrong in its own direction.
Polling several times a second spends the review re-reading a document nobody has written yet, and a timeout of the session's own reports a live review as stuck.
Neither is needed, because a reviewer that answers nothing before its deadline ends the run as `failed`, so every run reaches a state to act on.

#### Scenario: The waiting command is status

- **WHEN** the skill describes waiting for a review
- **THEN** it names `socrates-review status` and gives the reason

#### Scenario: The waiting step says how often and for how long

- **WHEN** the skill describes waiting for a review
- **THEN** it names a polling interval and says that the run ends on its own, so the session invents neither

#### Scenario: Watch is for the person

- **WHEN** the skill mentions `socrates-review watch`
- **THEN** it is as a suggestion for the person in another terminal, and it appears in no loop the session is told to run

### Requirement: The response rules come from the same source the tool prints

The four response rules the skill states SHALL be read from the same source `socrates-review findings` prints them from, and a test SHALL compare the two renderings.

Two copies of four sentences drift, and the direction is predictable: the copy nobody reads while implementing stays as it was, and the copy in front of the session gets softened.

#### Scenario: One source, two surfaces

- **WHEN** the rules the skill states are compared with the rules `socrates-review findings` prints
- **THEN** they are the same text

#### Scenario: The rules are stated before the findings are read

- **WHEN** the skill is rendered
- **THEN** the response rules appear before the step that reads the findings

### Requirement: `socrates-review skill` installs it and says where

`socrates-review skill` SHALL write the skill at user level and print the path it wrote.
`--to` SHALL write it to a named directory instead, which is what the repository's own generation target uses.

Running it again with the same target SHALL succeed and produce the same result, so that whether it has been run is not a question anybody has to answer.

#### Scenario: Installing names the path

- **WHEN** `socrates-review skill` is run
- **THEN** the skill exists at user level and the command printed the path

#### Scenario: Installing twice

- **WHEN** `socrates-review skill` is run twice
- **THEN** both succeed, the second says the same thing, and the file is the same

#### Scenario: Installing elsewhere

- **WHEN** `socrates-review skill --to` names a directory
- **THEN** the skill is written there and nothing is written at user level

#### Scenario: A directory that cannot be written

- **WHEN** the target directory cannot be written
- **THEN** the command refuses naming the path and the reason, and exits with the usage code

### Requirement: A further look is a new run, and only where the person asks

The skill SHALL say that a run has at most two rounds, that a further look is a fresh `socrates-review review` rather than a third round, and that one is started only where the person asks for it.

It SHALL say that the reviewer of a new run is told nothing about the run before it, so the person knows what they are getting.
That is a cold reading of the code as it now stands, not a continuation of the argument.

A session that believes a third round exists will look for the flag that gives it one, and there is none.
A session that starts a fresh review on its own initiative spends a full review on a change nobody asked it to re-read, which is the cost this whole design is careful with.

#### Scenario: There is no third round

- **WHEN** the skill describes a run that escalated at the round cap
- **THEN** it says that no further round is possible and names a fresh review as the only way to look again

#### Scenario: A fresh review is the person's to ask for

- **WHEN** the skill describes starting another review after a run ended
- **THEN** it says to offer it and to wait for the person to ask, rather than starting one unprompted

#### Scenario: A new run knows nothing of the old one

- **WHEN** the skill describes a fresh review of patched code
- **THEN** it says the reviewer is told nothing about the earlier run, its findings or its responses

### Requirement: The skill is named for the product

The skill SHALL be named `socrates-review`, SHALL be installed at `.claude/skills/socrates-review/`, and SHALL name `socrates-review` in every invocation its steps give.

Exactly one constant in the source SHALL carry that name, and it SHALL be the one the installer writes the directory from and the one the surface check reads invocations with.
Two constants would let the skill teach a command name that the check never looks for, and the check that proves every named command exists would go on passing while naming nothing.

#### Scenario: The skill is installed under the product's name

- **WHEN** `socrates-review skill` is run
- **THEN** the skill is written at `.claude/skills/socrates-review/SKILL.md`

#### Scenario: One constant names the skill and reads its invocations

- **WHEN** the constant that names the skill is changed
- **THEN** the installer's directory and the invocations the surface check reads both follow it, and no second copy of the name remains

#### Scenario: The skill names the one name the tool has

- **WHEN** the skill's steps are read
- **THEN** every invocation in them names `socrates-review`, and the tool has no other name for one to use instead

### Requirement: The skill writes an invocation as code

The skill SHALL write every invocation it names as code, inline or as an indented block, and the check that reads its invocations SHALL read code alone.

A sentence whose subject is the tool and an invocation of it are distinguishable by no rule that reads prose.
Read as prose, `socrates patches nothing` names a command called `patches`, and the check that exists to prove every named command is real would report five commands the binary has never had.

#### Scenario: Prose naming the product is not an invocation

- **WHEN** the skill states a sentence whose subject is the product rather than a command
- **THEN** the check reads no command from it, and every invocation written as code is still read
