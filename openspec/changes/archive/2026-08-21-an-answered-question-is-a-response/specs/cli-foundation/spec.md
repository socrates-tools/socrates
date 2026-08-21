## MODIFIED Requirements

### Requirement: `scrat status` says what the run is waiting for

Where a run is `awaiting-response`, `scrat status` SHALL say so, say which round it is on, say how many findings are still unanswered, and say how many of those need a person rather than a patch.

The first count is what a waiting session needs in order to know whether it has more work to do.
The second is what it needs in order to know whether the work is its own: a round waiting on a decision is a round where the session's next move is to stop and ask, and a caller that must branch on that SHALL be able to without reading prose.
The exit code stays the one for a run that is not terminal, because a parked run is not terminal whichever kind of answer it waits on, and one code cannot report a round that waits on both.

#### Scenario: A parked run reports what it wants

- **WHEN** `scrat status` is run against a run awaiting responses
- **THEN** it says the run is awaiting responses, names the round, and gives the count of unanswered findings

#### Scenario: A parked run says what needs a person

- **WHEN** `scrat status` is run against a run whose unanswered findings include one that asks a person
- **THEN** it gives the number needing a person alongside the total, and both are in the JSON object

#### Scenario: The count reaches zero

- **WHEN** every blocking finding has been responded to
- **THEN** `scrat status` reports zero unanswered and says the run is ready to continue

### Requirement: `scrat respond` takes one finding at a time

`scrat respond` SHALL act on one finding per invocation, named by its identifier, with a fix, a rejection and its reason, or a decision and its reason.

One finding per invocation is deliberate.
A batch interface invites a session to answer six findings in one call with one reason that covers none of them well, and the reason is the entire value of a rejection or a decision.

The reason SHALL be accepted as an argument or on standard input, so that a long one does not have to survive a shell.
Standard input SHALL be read only where the invocation asked for it, because a command that reads it on its own hangs where nothing is written and stores whatever is where a reason should be.

What the author said SHALL never be dropped without a word.
A flag belonging to another resolution SHALL be refused, naming the one that applies, rather than ignored: the author's claim is the whole value of a response, and one silently discarded reaches the next reviewer as an author who said nothing.

#### Scenario: One finding per invocation

- **WHEN** `scrat respond` is examined
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

- **WHEN** `scrat respond` is run against a run that is clean, escalated, failed or still reviewing
- **THEN** it is refused, naming the run's state, and nothing is stored
