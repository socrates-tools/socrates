## ADDED Requirements

### Requirement: The skill has one source and the committed copy is generated from it

The skill's text SHALL live in exactly one place in the source, and the copy committed to the repository SHALL be generated from it.
A test SHALL regenerate and compare, and SHALL fail while the two differ.

A committed file that is read rather than generated has one failure mode and it is the one that matters: somebody edits the file, the generator is never run, and the installed skill and the repository's copy differ with nothing failing.

#### Scenario: Drift fails the suite

- **WHEN** the committed skill file differs from what the source generates
- **THEN** the test fails, naming the file and the command that regenerates it

#### Scenario: Regenerating is one command

- **WHEN** the regeneration target is run
- **THEN** the committed copy matches the source and the test passes

### Requirement: The skill teaches the loop, including every terminal branch

The skill SHALL describe the sequence a session follows: implement, start a review carrying the person's stated intent, poll until the run is not working, and act on each terminal state.

It SHALL name what to do for each of `clean`, `failed`, `escalated` and `awaiting-response`, and `failed` SHALL be described before `escalated`.

A session scanning for what to do about findings reads past a state that has none, and a `failed` run has none.

#### Scenario: Every terminal state is covered

- **WHEN** the skill is rendered
- **THEN** it names `clean`, `failed`, `escalated` and `awaiting-response`, and says what to do for each

#### Scenario: A failed run is stated not to be a pass

- **WHEN** the skill describes `failed`
- **THEN** it says that no review completed, that nothing has been checked, and that this is not a clean result

#### Scenario: The loop closes

- **WHEN** the skill describes `awaiting-response`
- **THEN** it says to respond to every blocking finding, then continue, then return to polling

### Requirement: The skill teaches polling rather than following

The skill SHALL name `scrat status` as the command to wait with, and SHALL say that everything the session reads lands in its own context.

It SHALL name `scrat watch` only as something to suggest to the person, in another terminal, and SHALL NOT put it in the session's own loop.

#### Scenario: The waiting command is status

- **WHEN** the skill describes waiting for a review
- **THEN** it names `scrat status` and gives the reason

#### Scenario: Watch is for the person

- **WHEN** the skill mentions `scrat watch`
- **THEN** it is as a suggestion for the person in another terminal, and it appears in no loop the session is told to run

### Requirement: The response rules come from the same source the tool prints

The four response rules the skill states SHALL be read from the same source `scrat findings` prints them from, and a test SHALL compare the two renderings.

Two copies of four sentences drift, and the direction is predictable: the copy nobody reads while implementing stays as it was, and the copy in front of the session gets softened.

#### Scenario: One source, two surfaces

- **WHEN** the rules the skill states are compared with the rules `scrat findings` prints
- **THEN** they are the same text

#### Scenario: The rules are stated before the findings are read

- **WHEN** the skill is rendered
- **THEN** the response rules appear before the step that reads the findings

### Requirement: The skill states what socrates will not do

The skill SHALL state each of the following, because each prevents a specific attempt a session makes:

- socrates never patches anything. The session patches, in its own working tree.
- There is no way to explain the change to the reviewer. The reviewer receives the diff, the repository and the stated intent, and that is the design rather than a limitation to work around.
- A rejection is a claim that goes to the next reviewer, not a way to dismiss a finding, and a rejection with no reason is refused.

#### Scenario: The three prohibitions are present

- **WHEN** the skill is rendered
- **THEN** it states each of the three

#### Scenario: No channel is implied for the author's account

- **WHEN** the skill is rendered
- **THEN** it names no command, flag or file through which a session could send the reviewer its own account of the change

### Requirement: The skill drives only the surface that exists

Every command and flag the skill names SHALL exist, and a test SHALL check each against the built binary.

A skill naming a flag that does not exist is a session doing something plausible and wrong, and it fails at the point where the person has stopped watching.

#### Scenario: Every command named is real

- **WHEN** the commands and flags the skill names are extracted and checked against the binary
- **THEN** every one of them exists

#### Scenario: A removed flag fails the suite

- **WHEN** a flag the skill names is removed from the binary
- **THEN** the check fails, naming the flag and the skill

### Requirement: `scrat skill` installs it and says where

`scrat skill` SHALL write the skill at user level and print the path it wrote.
`--to` SHALL write it to a named directory instead, which is what the repository's own generation target uses.

Running it again with the same target SHALL succeed and produce the same result, so that whether it has been run is not a question anybody has to answer.

#### Scenario: Installing names the path

- **WHEN** `scrat skill` is run
- **THEN** the skill exists at user level and the command printed the path

#### Scenario: Installing twice

- **WHEN** `scrat skill` is run twice
- **THEN** both succeed, the second says the same thing, and the file is the same

#### Scenario: Installing elsewhere

- **WHEN** `scrat skill --to` names a directory
- **THEN** the skill is written there and nothing is written at user level

#### Scenario: A directory that cannot be written

- **WHEN** the target directory cannot be written
- **THEN** the command refuses naming the path and the reason, and exits with the usage code
