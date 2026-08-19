# cli-foundation Specification

## Purpose
TBD - created by archiving change a-review-nobody-authored. Update Purpose after archive.
## Requirements
### Requirement: The binary is `scrat` and it does one thing

socrates SHALL ship exactly one binary, named `scrat`.
It SHALL provide `review`, `status`, `findings` and `cancel` in this change, and it SHALL NOT provide a command that commits, pushes, opens a pull request, watches continuous integration, or runs the repository's tests.

The scope rule is mechanical rather than a judgement call about what is close enough to review, because judgement calls erode.

#### Scenario: The command surface is what it says

- **WHEN** `scrat --help` is run
- **THEN** it lists `review`, `status`, `findings` and `cancel`, and nothing that lands or verifies a change

### Requirement: `scrat review` starts a run and returns

`scrat review` SHALL take a snapshot, spawn a reviewer, print the run identifier and how to follow it, and return without waiting.

`scrat review --wait` SHALL instead block until the run is terminal and exit with the verdict's code.
Both are needed: a session goes back to work, a person at a terminal wants an answer.

An optional stated intent SHALL be accepted, and it SHALL be the only free text a caller may supply.

#### Scenario: The default does not block

- **WHEN** `scrat review` is invoked
- **THEN** it prints the run identifier, exits zero for having started a run, and the review is still going

#### Scenario: Waiting returns the verdict

- **WHEN** `scrat review --wait` is invoked and the review finds nothing blocking
- **THEN** it blocks until the run is terminal and exits with the code for `clean`

#### Scenario: Starting is not a verdict

- **WHEN** `scrat review` starts a run that later escalates
- **THEN** the exit code of the starting invocation is zero, because it reports that a run started rather than what it found

### Requirement: `scrat status` is deliberately tiny

`scrat status` SHALL print a fixed small number of lines: the run identifier, its state, the round it is on, and once terminal, the verdict and the counts of findings by who they are for.

It SHALL NOT stream, SHALL NOT follow, and SHALL NOT print the reviewer's output.
It is called in a loop by the session waiting on the review, and everything it prints lands in that session's context, which is the scarcest resource in the whole system.

It SHALL exit with the verdict's code once the run is terminal, and with the code for a run still working while it is not, so that a caller can branch on the exit code without parsing the text.

#### Scenario: A running review is a few lines

- **WHEN** `scrat status` is run against a review in progress
- **THEN** it prints at most a handful of lines, none of them the reviewer's output, and exits with the code meaning still working

#### Scenario: A terminal run reports its verdict

- **WHEN** `scrat status` is run against a run that escalated
- **THEN** it prints the verdict and how many findings need the person, and exits with the code for `escalated`

#### Scenario: Machine-readable on request

- **WHEN** `scrat status --json` is run
- **THEN** it prints one JSON object carrying the same facts, and nothing else

### Requirement: `scrat findings` is how the findings are read

`scrat findings` SHALL print the findings of a run's most recent round: each one's identifier, severity, file, line, description and who it is for, followed by the risk level and its rationale.

`scrat findings --json` SHALL print the findings document as socrates classified it, with every finding's effective action resolved, so that a reader never has to reapply the defaulting rule and never has a chance to apply it differently.

This is the command the session that asked for the review reads once the run is terminal.
It is separate from `status` because `status` is polled and this is not, and separate from reading the findings file directly because the file holds what the reviewer said and this holds what socrates concluded.

#### Scenario: Findings are readable

- **WHEN** `scrat findings` is run against an escalated run
- **THEN** it prints every finding with its identifier, severity, file, line, description and who it is for, and the risk level and rationale

#### Scenario: The effective action is resolved for the reader

- **WHEN** a finding carried no action and `scrat findings --json` is run
- **THEN** the finding is reported as needing a person, and the reader is not given the unresolved field to interpret

#### Scenario: A clean run has findings to read too

- **WHEN** `scrat findings` is run against a clean run that returned informational findings
- **THEN** it prints them, and exits with the code for `clean`

#### Scenario: A failed run has no findings to give

- **WHEN** `scrat findings` is run against a run that failed
- **THEN** it says no review completed, prints no findings, and exits `2`

### Requirement: Naming no run means the current repository's most recent

Every command that takes a run SHALL accept a run identifier, and SHALL default to the most recent run for the repository the working directory is in.

Where that repository has no run, the command SHALL say so and exit with the usage code rather than reporting on some other repository's run.

#### Scenario: The bare command finds this repository's run

- **WHEN** runs exist for two repositories and `scrat status` is run inside one of them
- **THEN** it reports that repository's most recent run

#### Scenario: A repository with no runs

- **WHEN** `scrat status` is run in a repository that has never been reviewed
- **THEN** it says so and exits with the usage code

#### Scenario: An identifier wins over the default

- **WHEN** a run identifier is given
- **THEN** that run is reported on, whatever the working directory is

### Requirement: Exit codes name the outcome and are stable

Commands that report on a run SHALL exit:

- `0` for `clean`, and for a command whose job was to start something that started
- `1` for `escalated`
- `2` for `failed`
- `3` for a usage or precondition error
- `4` for a run that is not terminal yet

These are a contract.
A later layer wires `scrat` into a hook whose whole behaviour is the exit code, and a code that shifts meaning between releases silently changes what that hook does.

#### Scenario: Each verdict has its own code

- **WHEN** commands are run against runs that are clean, escalated, failed and still working
- **THEN** the exit codes are `0`, `1`, `2` and `4` respectively

#### Scenario: A usage error is not a verdict

- **WHEN** a command is given a run identifier that does not exist
- **THEN** it exits `3`, and prints no verdict

#### Scenario: Failure is not escalation

- **WHEN** a run failed because no review completed
- **THEN** the exit code is `2` and not `1`, because nothing was checked and the two call for different actions

### Requirement: `scrat cancel` stops a run and says it stopped

`scrat cancel` SHALL terminate the reviewer's process group, mark the run `cancelled`, and remove the snapshot worktree.

A cancelled run SHALL NOT be reported as clean, failed or escalated.
A person who stopped a review knows nothing was decided, and a state that says otherwise is a lie about what was checked.

#### Scenario: Cancelling stops everything

- **WHEN** `scrat cancel` is run against a review in progress
- **THEN** no process from the reviewer's group survives, the worktree is gone, and the run is `cancelled`

#### Scenario: Cancelling a terminal run changes nothing

- **WHEN** `scrat cancel` is run against a run that already finished
- **THEN** the run's state is unchanged and the command says the run was already finished

### Requirement: Configuration is one optional user file

socrates SHALL read one configuration file under the user's configuration directory, and SHALL run with defaults where it is absent.
It SHALL hold the per-role model and effort, the inlined-diff ceiling, the agent deadline, the log ceiling and the run retention age, and nothing else.

socrates SHALL NOT read configuration from the repository under review in this change, SHALL NOT create a configuration file on its own, and SHALL NOT write anything at all inside a repository.

#### Scenario: No configuration file

- **WHEN** `scrat review` is run with no configuration file present
- **THEN** it runs with defaults and creates no configuration file

#### Scenario: The repository is never written to

- **WHEN** a review runs to completion
- **THEN** nothing inside the repository has changed except the git objects the snapshot needed, and no socrates file exists anywhere under it

#### Scenario: An unreadable configuration file is a refusal

- **WHEN** the configuration file exists and does not parse
- **THEN** the command refuses naming the file and the parse error, and exits with the usage code, rather than silently running on defaults

### Requirement: What a person reads is plain and says what to do next

Every message a person or a session reads SHALL name what happened and what to do about it.
A refusal SHALL name the condition, not only that something went wrong.

No message SHALL contain an em dash.

#### Scenario: A refusal names the next action

- **WHEN** any refusal in this specification is triggered
- **THEN** the message names the condition and what the person can do about it

#### Scenario: No em dash reaches a person

- **WHEN** every string socrates can print is examined
- **THEN** none contains an em dash

