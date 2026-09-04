# run-presentation Specification

## Purpose
TBD - created by archiving change one-window-three-needs. Update Purpose after archive.

## Requirements

### Requirement: Rendering is a pure function of the events

Turning the event log into something a person reads SHALL be a function from a sequence of events to a sequence of lines, reaching no filesystem, no git, no clock and no terminal.

Nothing SHALL write pre-rendered text into the log.
The same bytes are summarised by one command, rendered by another and parsed by whatever comes after, and a log that carries display text has already chosen for all three.

#### Scenario: Rendering needs no terminal

- **WHEN** the renderer is given a sequence of events
- **THEN** it produces lines with no terminal, no pseudo-terminal and no environment attached

#### Scenario: The renderer decides nothing about the run

- **WHEN** the package holding the renderer is examined
- **THEN** it imports nothing that reads the filesystem, the network, the clock or a git repository

#### Scenario: An unknown event kind is still shown

- **WHEN** the log carries an event of a kind the renderer does not know
- **THEN** it is rendered legibly rather than dropped, because a reader looking at a failure needs everything that was recorded

### Requirement: A follower started too early waits rather than failing

Where the event log does not exist yet, `socrates-review watch` SHALL wait for it, bounded, reporting the run's state while it waits.

`socrates-review review` returns as soon as the run exists, so a person who runs `watch` immediately arrives before the reviewer has written a byte.
Failing in that window punishes them for being quick, in exactly the window where they are most likely to try.

Where the wait expires and no reviewer is alive, the follower SHALL report the run as failed rather than continuing to wait.

#### Scenario: Watching immediately after starting

- **WHEN** `socrates-review watch` is run in the moment between the run directory appearing and the reviewer's first event
- **THEN** it waits, says the run is starting, and then follows normally

#### Scenario: Waiting is bounded

- **WHEN** no event log appears and no reviewer process is alive
- **THEN** the follower stops within its bound, reports the run as failed, and exits `2`

#### Scenario: A run identifier that does not exist

- **WHEN** `socrates-review watch` is given a run identifier no run has
- **THEN** it refuses at once, naming the identifier, and exits with the usage code rather than waiting for a run that will never appear

### Requirement: A rotated log says so where it is read

Where a run's event log was rotated, both reading commands SHALL say that earlier events are gone.

A reader looking for the moment something went wrong, in a log that silently lost its beginning, concludes the wrong thing about what socrates did.

#### Scenario: A rotated log is marked

- **WHEN** a run's log rotated during the review and is then rendered
- **THEN** the rendering says that earlier events were discarded, before the first event it does show

### Requirement: `socrates-review watch` follows a run and stops when it is over

`socrates-review watch` SHALL render events as they arrive, stop when the run reaches a terminal state, print the verdict, and exit with the verdict's code.

It SHALL treat a terminal event in the log and a terminal state in `run.json` as equally sufficient to stop.
Waiting only for a terminal event is the version that hangs: a run that failed before its reviewer wrote anything has a terminal state and an empty log.

It SHALL also stop where the round parked on the author, print what the round is waiting for, and exit with the code for a run that is not terminal yet.
A parked round has no process of its own running and ends only when the author answers it.
A follower that treated parking as work in progress would wait for events that arrive only after the person acts.

#### Scenario: Following a review to a clean verdict

- **WHEN** `socrates-review watch` is run against a review in progress that goes on to find nothing blocking
- **THEN** it renders the events as they arrive, prints the clean verdict, and exits `0`

#### Scenario: Following a review to an escalation

- **WHEN** the run escalates
- **THEN** the follower stops, prints the verdict and what needs the person, and exits `1`

#### Scenario: A run that failed before writing anything

- **WHEN** `socrates-review watch` is run against a run whose state is `failed` and whose event log is empty
- **THEN** it stops at once, prints the failure, and exits `2` rather than waiting

#### Scenario: A follower does not outlive the run

- **WHEN** a run reaches any terminal state while a follower is attached
- **THEN** the follower exits within a bounded time, and never waits for further events

#### Scenario: A follower stops where the round stops

- **WHEN** the round parks on the author while a follower is attached
- **THEN** the follower stops, says what the round is waiting for, and exits `4`

#### Scenario: Cancelling while following

- **WHEN** the run is cancelled from another terminal while a follower is attached
- **THEN** the follower stops, says the run was cancelled, and exits with the code for a cancelled run

### Requirement: `socrates-review logs` reads a finished run and does not follow

`socrates-review logs` SHALL render a run's events and exit, without following and without waiting for anything.

Where its output is a terminal and a pager is configured, it SHALL page.
Where its output is not a terminal, it SHALL write plainly, so that the output can be piped, redirected and searched.
Where a configured pager cannot be started, it SHALL fall back to writing plainly rather than failing.

#### Scenario: Reading a finished run

- **WHEN** `socrates-review logs` is run against a terminal run
- **THEN** it renders every event the log holds and exits with the verdict's code

#### Scenario: Reading a run still in progress

- **WHEN** `socrates-review logs` is run against a run that is still reviewing
- **THEN** it renders what the log holds so far, does not wait, and exits with the code meaning still working

#### Scenario: Piping works

- **WHEN** the output of `socrates-review logs` is a pipe
- **THEN** no pager is started and the output is plain

#### Scenario: A pager that cannot be started

- **WHEN** the configured pager does not exist
- **THEN** the output is written plainly and the command still exits with the verdict's code
