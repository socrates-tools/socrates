## MODIFIED Requirements

### Requirement: The binary is `scrat` and it does one thing

socrates SHALL ship exactly one binary, named `scrat`.
It SHALL provide `review`, `status`, `findings`, `respond`, `continue` and `cancel`, and it SHALL NOT provide a command that commits, pushes, opens a pull request, watches continuous integration, or runs the repository's tests.

The scope rule is mechanical rather than a judgement call about what is close enough to review, because judgement calls erode.
`respond` and `continue` record and advance; neither writes in the repository, and adding them does not move the boundary.

#### Scenario: The command surface is what it says

- **WHEN** `scrat --help` is run
- **THEN** it lists `review`, `status`, `findings`, `respond`, `continue` and `cancel`, and nothing that lands or verifies a change

#### Scenario: Recording a response writes nothing in the repository

- **WHEN** `scrat respond` and `scrat continue` are run
- **THEN** nothing inside the repository has changed except the git objects the new round's snapshot needed

## ADDED Requirements

### Requirement: `scrat status` says what the run is waiting for

Where a run is `awaiting-response`, `scrat status` SHALL say so, say which round it is on, and say how many findings are still unanswered.

That count is the only thing the waiting session needs in order to know whether it has more work to do, and it is one number.

#### Scenario: A parked run reports what it wants

- **WHEN** `scrat status` is run against a run awaiting responses
- **THEN** it says the run is awaiting responses, names the round, and gives the count of unanswered findings

#### Scenario: The count reaches zero

- **WHEN** every blocking finding has been responded to
- **THEN** `scrat status` reports zero unanswered and says the run is ready to continue

### Requirement: `scrat respond` takes one finding at a time

`scrat respond` SHALL act on one finding per invocation, named by its identifier, with either a fix or a rejection and its reason.

One finding per invocation is deliberate.
A batch interface invites a session to answer six findings in one call with one reason that covers none of them well, and the reason is the entire value of a rejection.

The reason SHALL be accepted as an argument or on standard input, so that a long one does not have to survive a shell.
Standard input SHALL be read only where the invocation asked for it, because a command that reads it on its own hangs where nothing is written and stores whatever is where a reason should be.

What the author said SHALL never be dropped without a word.
A flag belonging to the other resolution SHALL be refused, naming the one that applies, rather than ignored: the author's claim is the whole value of a response, and one silently discarded reaches the next reviewer as an author who said nothing.

#### Scenario: One finding per invocation

- **WHEN** `scrat respond` is examined
- **THEN** it takes exactly one finding identifier and offers no way to name several

#### Scenario: A reason from standard input

- **WHEN** a rejection reason is supplied on standard input
- **THEN** it is stored verbatim, including its line breaks

#### Scenario: A flag that does not apply to the answer

- **WHEN** a fix is recorded with a rejection's reason, or a rejection with a fix's note
- **THEN** it is refused, naming the flag that applies to the chosen answer, and nothing is stored

#### Scenario: A reason is read only where it is asked for

- **WHEN** a rejection names no reason and standard input holds something nobody meant as one
- **THEN** it is refused for having no reason, and what was on standard input is not stored as the author's claim

#### Scenario: Responding to a run that is not awaiting responses

- **WHEN** `scrat respond` is run against a run that is clean, escalated, failed or still reviewing
- **THEN** it refuses naming the run's state, and stores nothing
