# cli-foundation Specification

## Purpose
TBD - created by archiving change a-review-nobody-authored. Update Purpose after archive.

## Requirements

### Requirement: Naming no run means the current repository's most recent

Every command that takes a run SHALL accept a run identifier, and SHALL default to the most recent run for the repository the working directory is in.

Where that repository has no run, the command SHALL say so and exit with the usage code rather than reporting on some other repository's run.

#### Scenario: The bare command finds this repository's run

- **WHEN** runs exist for two repositories and `socrates-review status` is run inside one of them
- **THEN** it reports that repository's most recent run

#### Scenario: A repository with no runs

- **WHEN** `socrates-review status` is run in a repository that has never been reviewed
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
A later layer wires `socrates-review` into a hook whose whole behaviour is the exit code, and a code that shifts meaning between releases silently changes what that hook does.

#### Scenario: Each verdict has its own code

- **WHEN** commands are run against runs that are clean, escalated, failed and still working
- **THEN** the exit codes are `0`, `1`, `2` and `4` respectively

#### Scenario: A usage error is not a verdict

- **WHEN** a command is given a run identifier that does not exist
- **THEN** it exits `3`, and prints no verdict

#### Scenario: Failure is not escalation

- **WHEN** a run failed because no review completed
- **THEN** the exit code is `2` and not `1`, because nothing was checked and the two call for different actions

### Requirement: Configuration is one optional user file

socrates SHALL read one configuration file under its own home directory, and SHALL run with defaults where it is absent.
It SHALL hold the per-role model and effort, the inlined-diff ceiling, the agent deadline, the log ceiling, the run retention age, and one block holding what an operator states about the surface, and nothing else.

That block SHALL hold the address the surface binds and the address it is reached at.
The surface SHALL read it from the one file rather than parsing a file of its own, so that an operator configuring this tool opens one file.

Everything socrates reads and everything it writes SHALL live under that one home directory, and so SHALL everything the surface reads and writes.
A stated home SHALL be refused unless it is absolute, rather than ignored in favour of the default.

That home SHALL be a directory named for this tool under the brand's root, and one environment variable SHALL move the brand's root.
Where none is stated the root SHALL be a directory named for the brand under the person's own home.
The tool's directory SHALL be named for the part of its name the brand does not supply, so that the brand's root is not spelled twice on one path.

The brand publishes more than one tool, and a person who runs two of them asks where each keeps things exactly once.
A root per brand and a directory per tool answers that for the tool that exists and for the tool that does not exist yet, without either one having to move.
A tool SHALL NOT read or write outside its own directory under that root, so that removing a tool is removing a directory.

One directory rather than two is a question with one answer: where socrates keeps things.
The answer is given to a stranger reading an issue that will not reproduce.
A second directory belonging to a second binary grows that answer a second half, for a program the stranger has never heard of.
Refusing a relative home rather than ignoring it is the same argument: a stated home that silently does nothing is a person who believes their state moved and is looking in the wrong place.

socrates SHALL NOT read configuration from the repository under review, SHALL NOT create a configuration file on its own, and SHALL NOT write anything at all inside a repository.

#### Scenario: No configuration file

- **WHEN** `socrates-review review` is run with no configuration file present
- **THEN** it runs with defaults and creates no configuration file

#### Scenario: One variable moves everything

- **WHEN** the brand's root is stated and a review is run
- **THEN** both the configuration file it reads and the runs it writes are under the tool's own directory in the stated root, and nothing is written anywhere else

#### Scenario: The tool writes under its own directory in the brand's root

- **WHEN** no root is stated, a review is run and the person's home is examined
- **THEN** the brand's root is a directory named for the brand, the tool's configuration file and runs are in a directory named for the tool inside it, and neither is directly in the root

#### Scenario: A relative home is refused

- **WHEN** socrates' home is stated as a relative path
- **THEN** the command refuses, names the variable and the path, and does not fall back to the default

#### Scenario: The repository is never written to

- **WHEN** a review runs to completion
- **THEN** nothing inside the repository has changed except the git objects the snapshot needed, and no socrates file exists anywhere under it

#### Scenario: An unreadable configuration file is a refusal

- **WHEN** the configuration file exists and does not parse
- **THEN** the command refuses naming the file and the parse error, and exits with the usage code, rather than silently running on defaults

#### Scenario: What an operator states about the surface is read from the one file

- **WHEN** the configuration file states the address the surface binds and the address it is reached at, and a surface is served
- **THEN** the surface binds what was stated and prints what was stated, having read no file of its own

#### Scenario: A surface block that names a field socrates does not know is refused

- **WHEN** the surface block of the configuration file names a field that is not one of its own
- **THEN** the command refuses naming the field, exactly as it does for a field of any other block

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

### Requirement: The three reading commands are for three consumers

`status`, `watch` and `logs` SHALL stay distinct, and none SHALL grow into another.

- `status` answers one shot, in a few lines, for the session waiting on the review. It never streams and never follows.
- `watch` follows live, for the person in a second tab. It always follows and always stops at a terminal state.
- `logs` reads what happened, for the person afterwards. It never follows.

`status` SHALL NOT gain a follow option and `watch` SHALL NOT gain a one-shot option.
Either would be the same command reached two ways.
They are separate because one prints into a language model's context and two print into a person's terminal, which is an order of magnitude difference in what verbosity costs.

#### Scenario: Status stays small when the reviewer is loud

- **WHEN** a reviewer has emitted a large volume of output and `socrates-review status` is run
- **THEN** its output is the same handful of lines it prints for a quiet run

#### Scenario: Neither command grows the other's flag

- **WHEN** the flags of `status` and `watch` are examined
- **THEN** `status` offers no way to follow and `watch` offers no way to return after one read

#### Scenario: All three agree about the run

- **WHEN** the three commands are run against one terminal run
- **THEN** all three report the same verdict and exit with the same code

### Requirement: The binary is `socrates-review` and it does one thing

socrates SHALL ship exactly one binary, named `socrates-review` after the tool.
Installing it SHALL put no second name beside it, and the binary SHALL answer to that name alone.

One name rather than two.
A second name is a second answer to what the tool is called, and it is given to a stranger reading an issue, a person searching for the documentation, and a session that has to write an invocation down.
The brand the tool is published under supplies the shorter word already, and it names the organisation rather than this tool.

It SHALL provide `review`, `status`, `findings`, `respond`, `continue`, `watch`, `logs`, `cancel`, `skill` and `version`, and it SHALL NOT provide a command that commits, pushes, opens a pull request, watches continuous integration, or runs the repository's tests.

The scope rule is mechanical rather than a judgement call about what is close enough to review, because judgement calls erode.
`respond` and `continue` record and advance; `skill` writes at user level; `version` reads what the binary records about itself and reaches nothing.
None of them writes in the repository, and adding them does not move the boundary.

#### Scenario: The command surface is what it says

- **WHEN** `socrates-review --help` is run
- **THEN** it lists `review`, `status`, `findings`, `respond`, `continue`, `watch`, `logs`, `cancel`, `skill` and `version`, and nothing that lands or verifies a change

#### Scenario: One name is installed and nothing is put beside it

- **WHEN** socrates is installed and what it put on the path is examined
- **THEN** `socrates-review` is there, and no second name for that same binary is

#### Scenario: Recording a response writes nothing in the repository

- **WHEN** `socrates-review respond` and `socrates-review continue` are run
- **THEN** nothing inside the repository has changed except the git objects the new round's snapshot needed

#### Scenario: Installing the skill writes nothing in the repository

- **WHEN** `socrates-review skill` is run with its default target
- **THEN** the skill is written at user level and nothing inside the repository has changed

### Requirement: What socrates-review prints names socrates-review

Every string socrates prints SHALL name `socrates-review`, including the prefix on a failure, the command a refusal tells the person to run next, and the step a report says to take.

No string it prints SHALL name a second name for the binary, and no check SHALL exist that reconciles two.
A person copying a command out of a message runs the command they were given, so the name in the message is the name that has to work.

#### Scenario: What it prints names the tool

- **WHEN** the binary refuses something
- **THEN** what it prints names `socrates-review`, and the command it tells the person to run next is one they can type

### Requirement: `socrates-review version` says which build this is

`socrates-review version` SHALL print the short commit the binary was built from, and SHALL say when the working tree it was built from held uncommitted changes.

The commit SHALL be recorded in the binary when it is linked, and SHALL NOT be read at run time from the repository the command is invoked in.
A binary that reported the tree it is standing in would answer the question this command exists for with the one answer that cannot be checked, since `socrates-review` is run against repositories it was not built from.

A binary carrying no recorded commit SHALL say so and SHALL say what to build it with.
It SHALL NOT print an empty version and SHALL NOT infer one, because a version line nobody can trust is worse than one that admits it is missing.

`socrates-review version` SHALL take no arguments, and SHALL refuse one with the usage code.

#### Scenario: The build says which commit it came from

- **WHEN** `socrates-review` is built through the Makefile and `socrates-review version` is run
- **THEN** it prints the short commit of the checkout it was built from

#### Scenario: A build from a tree with uncommitted changes says so

- **WHEN** `socrates-review version` is run on a binary linked from a working tree that held uncommitted changes
- **THEN** the line says the tree it was built from held them, beside the commit

#### Scenario: A build that recorded no commit says so

- **WHEN** `socrates-review version` is run on a binary that was linked without the stamp
- **THEN** it says no commit was recorded and says what to build it with, and prints no commit of its own

#### Scenario: `socrates-review version` takes no arguments

- **WHEN** `socrates-review version` is given an argument
- **THEN** it refuses, names the argument, and exits with the usage code

### Requirement: `socrates-review review` starts a run and returns

`socrates-review review` SHALL take a snapshot, spawn a reviewer, print the run identifier and how to follow it, and return without waiting.

`socrates-review review --wait` SHALL instead block until the run is terminal and exit with the verdict's code.
Both are needed: a session goes back to work, a person at a terminal wants an answer.

An optional stated intent SHALL be accepted, and it SHALL be the only free text a caller may supply.

#### Scenario: The default does not block

- **WHEN** `socrates-review review` is invoked
- **THEN** it prints the run identifier, exits zero for having started a run, and the review is still going

#### Scenario: Waiting returns the verdict

- **WHEN** `socrates-review review --wait` is invoked and the review finds nothing blocking
- **THEN** it blocks until the run is terminal and exits with the code for `clean`

#### Scenario: Starting is not a verdict

- **WHEN** `socrates-review review` starts a run that later escalates
- **THEN** the exit code of the starting invocation is zero, because it reports that a run started rather than what it found

### Requirement: `socrates-review status` is deliberately tiny

`socrates-review status` SHALL print a fixed small number of lines: the run identifier, its state, the round it is on, and once terminal, the verdict and the counts of findings by who they are for.

It SHALL NOT stream, SHALL NOT follow, and SHALL NOT print the reviewer's output.
It is called in a loop by the session waiting on the review, and everything it prints lands in that session's context, which is the scarcest resource in the whole system.

It SHALL exit with the verdict's code once the run is terminal, and with the code for a run still working while it is not, so that a caller can branch on the exit code without parsing the text.

#### Scenario: A running review is a few lines

- **WHEN** `socrates-review status` is run against a review in progress
- **THEN** it prints at most a handful of lines, none of them the reviewer's output, and exits with the code meaning still working

#### Scenario: A terminal run reports its verdict

- **WHEN** `socrates-review status` is run against a run that escalated
- **THEN** it prints the verdict and how many findings need the person, and exits with the code for `escalated`

#### Scenario: Machine-readable on request

- **WHEN** `socrates-review status --json` is run
- **THEN** it prints one JSON object carrying the same facts, and nothing else

### Requirement: `socrates-review findings` is how the findings are read

`socrates-review findings` SHALL print the findings of a run's most recent round: each one's identifier, severity, file, line, description and who it is for, followed by the risk level and its rationale.

`socrates-review findings --json` SHALL print the findings document as socrates classified it, with every finding's effective action resolved, so that a reader never has to reapply the defaulting rule and never has a chance to apply it differently.

This is the command the session that asked for the review reads once the run is terminal.
It is separate from `status` because `status` is polled and this is not.
It is separate from reading the findings file because the file holds what the reviewer said and this holds what socrates concluded.

#### Scenario: Findings are readable

- **WHEN** `socrates-review findings` is run against an escalated run
- **THEN** it prints every finding with its identifier, severity, file, line, description and who it is for, and the risk level and rationale

#### Scenario: The effective action is resolved for the reader

- **WHEN** a finding carried no action and `socrates-review findings --json` is run
- **THEN** the finding is reported as needing a person, and the reader is not given the unresolved field to interpret

#### Scenario: A clean run has findings to read too

- **WHEN** `socrates-review findings` is run against a clean run that returned informational findings
- **THEN** it prints them, and exits with the code for `clean`

#### Scenario: A failed run has no findings to give

- **WHEN** `socrates-review findings` is run against a run that failed
- **THEN** it says no review completed, prints no findings, and exits `2`

### Requirement: `socrates-review cancel` stops a run and says it stopped

`socrates-review cancel` SHALL terminate the reviewer's process group, mark the run `cancelled`, and remove the snapshot worktree.

A cancelled run SHALL NOT be reported as clean, failed or escalated.
A person who stopped a review knows nothing was decided, and a state that says otherwise is a lie about what was checked.

#### Scenario: Cancelling stops everything

- **WHEN** `socrates-review cancel` is run against a review in progress
- **THEN** no process from the reviewer's group survives, the worktree is gone, and the run is `cancelled`

#### Scenario: Cancelling a terminal run changes nothing

- **WHEN** `socrates-review cancel` is run against a run that already finished
- **THEN** the run's state is unchanged and the command says the run was already finished

### Requirement: `socrates-review status` says what the run is waiting for

Where a run is `awaiting-response`, `socrates-review status` SHALL say so, say which round it is on, say how many findings are still unanswered, and say how many of those need a person rather than a patch.

The first count is what a waiting session needs in order to know whether it has more work to do.
The second is what it needs in order to know whether the work is its own.
A round waiting on a decision is a round where the session's next move is to stop and ask, and a caller that must branch on that SHALL be able to without reading prose.
The exit code stays the one for a run that is not terminal, because a parked run is not terminal whichever kind of answer it waits on, and one code cannot report a round that waits on both.

#### Scenario: A parked run reports what it wants

- **WHEN** `socrates-review status` is run against a run awaiting responses
- **THEN** it says the run is awaiting responses, names the round, and gives the count of unanswered findings

#### Scenario: A parked run says what needs a person

- **WHEN** `socrates-review status` is run against a run whose unanswered findings include one that asks a person
- **THEN** it gives the number needing a person alongside the total, and both are in the JSON object

#### Scenario: The count reaches zero

- **WHEN** every blocking finding has been responded to
- **THEN** `socrates-review status` reports zero unanswered and says the run is ready to continue

### Requirement: `socrates-review respond` takes one finding at a time

`socrates-review respond` SHALL act on one finding per invocation, named by its identifier, with a fix, a rejection and its reason, or a decision and its reason.

One finding per invocation is deliberate.
A batch interface invites a session to answer six findings in one call with one reason that covers none of them well, and the reason is the entire value of a rejection or a decision.

The reason SHALL be accepted as an argument or on standard input, so that a long one does not have to survive a shell.
Standard input SHALL be read only where the invocation asked for it, because a command that reads it on its own hangs where nothing is written and stores whatever is where a reason should be.

What the author said SHALL never be dropped without a word.
A flag belonging to another resolution SHALL be refused, naming the one that applies, rather than ignored.
The author's claim is the whole value of a response, and one silently discarded reaches the next reviewer as an author who said nothing.

#### Scenario: One finding per invocation

- **WHEN** `socrates-review respond` is examined
- **THEN** it takes exactly one finding identifier and offers no way to name several

#### Scenario: A reason from standard input

- **WHEN** a rejection reason is supplied on standard input
- **THEN** it is stored verbatim, including its line breaks

#### Scenario: A decision's reason from standard input

- **WHEN** a decision's reason is supplied on standard input
- **THEN** it is stored verbatim on the same terms as a rejection's

#### Scenario: A flag that does not apply to the answer

- **WHEN** a fix is recorded with a rejection's reason, or a rejection with a fix's note
- **THEN** it is refused, naming the flag that applies to the chosen answer, and nothing is stored

#### Scenario: A reason is read only where it is asked for

- **WHEN** a rejection names no reason and standard input holds something nobody meant as one
- **THEN** it is refused for having no reason, and what was on standard input is not stored as the author's claim

#### Scenario: Responding to a run that is not awaiting responses

- **WHEN** `socrates-review respond` is run against a run that is clean, escalated, failed or still reviewing
- **THEN** it is refused, naming the run's state, and nothing is stored

### Requirement: `socrates-review findings` says whether the round waits on a person

`socrates-review findings` SHALL report whether the round it is printing waits on a person, as one line for a reader and one field in the object `--json` prints.

It SHALL report it from the one rule that decides it, and SHALL NOT restate the condition.

The line and the field SHALL name the wait rather than the reading.
Every round that completed a review is handed over as a page, so a fact reported as whether a person reads the round would be true of every round it is printed for and would decide nothing.
What still needs deciding is whether the session stops at the page, and that is what is reported.

`socrates-review status` SHALL NOT carry it.
`status` is called in a loop and everything it prints lands in the waiting session's context, which is the scarcest resource in the system.
`findings` is read once, on exactly the branches where the fact is needed, and is already where socrates tells the author the rules.

#### Scenario: A round that waits on a person says so

- **WHEN** `socrates-review findings` is run against a run that escalated, or against a parked round holding a finding a person answers
- **THEN** it says the round waits on a person, and the field is in the object `--json` prints

#### Scenario: A round the author answers alone says nothing about waiting

- **WHEN** `socrates-review findings` is run against a clean run, or against a parked round whose findings the author answers alone
- **THEN** it does not say the round waits on a person, and it says nothing that would stop the session handing the round over as a page

#### Scenario: `socrates-review status` does not carry the wait

- **WHEN** the lines and the JSON object of `socrates-review status` are examined for a run whose round waits on a person
- **THEN** neither carries it, and what `status` prints is what it printed before
