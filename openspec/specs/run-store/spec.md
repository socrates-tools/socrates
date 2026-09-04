# run-store Specification

## Purpose
TBD - created by archiving change a-review-nobody-authored. Update Purpose after archive.
## Requirements
### Requirement: A run is a directory whose name answers the common question

Every run SHALL own one directory under socrates' home, beside the configuration file rather than under a directory of its own, named `<UTC instant>-<repository key>`, where the instant is second-resolution and lexically sortable and the repository key is derived from the absolute path of the repository root.

Naming SHALL carry both, so that the most recent run for a given repository is answerable by listing directory names and reading no files.
The command that asks that question is called in a loop by a session waiting on a review, and a design that opens every run to answer it pays that cost on every poll.

#### Scenario: Runs sort by when they started

- **WHEN** three runs are started in the same repository in order
- **THEN** their directory names sort lexically into the order they were started

#### Scenario: Two repositories do not collide

- **WHEN** runs are started in two different repositories within the same second
- **THEN** their directory names differ, and each names its own repository's key

#### Scenario: The same path always gets the same key

- **WHEN** a repository key is derived twice for one repository root
- **THEN** the two are equal, and a different repository root produces a different key

#### Scenario: A run is found beside the configuration

- **WHEN** a review is run and socrates' home is examined
- **THEN** the run's directory and the configuration file are under that one home, and no run exists outside it

### Requirement: A run holds its state, its rounds and its log

A run directory SHALL hold:

- `run.json`: the repository root, the base commit, the snapshot commit, the stated intent, when it started, its current state, and its verdict once it has one.
- `events.jsonl`: the append-only event log.
- one directory per round, named `r<n>` from `r1`, holding that round's prompt, the reviewer's output stream and its findings.

Nothing else in a run directory SHALL be load-bearing, so that a run is legible with `ls` and `cat` at the moment socrates itself is what is broken.

#### Scenario: A started run is on disk immediately

- **WHEN** `scrat review` returns
- **THEN** the run directory exists, holds a `run.json` naming the repository and the snapshot commit, and holds `r1`

#### Scenario: Round output is under the round

- **WHEN** a reviewer has run
- **THEN** `r1` holds the exact prompt that was sent, the reviewer's output stream, and its findings, each as its own file

### Requirement: Every state write is atomic

`run.json` SHALL be written by writing a temporary file in the same directory and renaming it over the target.
A reader SHALL therefore observe either the previous complete state or the new complete state, never a partial one.

The writer is a detached process and the reader is a session polling it, and they share no other synchronisation.

#### Scenario: A reader never sees half a state

- **WHEN** `run.json` is read repeatedly while the run advances through every state
- **THEN** every read parses, and every read yields a state the run actually held

#### Scenario: A crash mid-write leaves the previous state

- **WHEN** the process writing `run.json` dies after writing the temporary file and before renaming it
- **THEN** the run's recorded state is the one before the write, and no temporary file is mistaken for the state

### Requirement: The event log is structured and bounded

The event log SHALL be JSONL, one object per line, each carrying at least an instant, a kind and a subject.
It SHALL NOT hold pre-rendered display text, because the same bytes are summarised by one command, rendered by another, and parsed by whatever comes after.

The log SHALL be bounded by a byte ceiling.
On reaching it, socrates SHALL rotate by truncating in place and retaining the tail rather than by renaming, because the reviewer holds an open descriptor and a rename leaves it writing to a file nobody reads.
A rotation SHALL record that it happened, so that a reader knows the log is not the whole story.

#### Scenario: Events are one JSON object per line

- **WHEN** a run completes
- **THEN** every line of `events.jsonl` parses as one JSON object carrying an instant, a kind and a subject

#### Scenario: A chatty reviewer does not fill the disk

- **WHEN** a reviewer emits more output than the log's ceiling
- **THEN** the log stays at or below the ceiling, holds the most recent events, and carries a record that a rotation occurred

#### Scenario: The writer keeps writing across a rotation

- **WHEN** the log rotates while the reviewer is still emitting
- **THEN** subsequent events appear in the same file, and no event is written to a file no reader will open

### Requirement: A run that nothing is advancing is not still working

A run recorded as `starting` or `reviewing` whose reviewer process no longer exists SHALL be reported as `failed`, naming that the reviewer went away.
A run whose state says it is working while nothing is working is a stall that looks alive, and a session polling it waits forever.

The recorded process identity SHALL be what this is judged on.
The judgement SHALL be made when the run is read rather than by anything polling in the background, because socrates has no process that outlives a run.

#### Scenario: A killed reviewer stops the run

- **WHEN** a reviewer process is killed and the run is then read
- **THEN** the run reports `failed`, naming that the reviewer went away, rather than `reviewing`

#### Scenario: A live reviewer is still working

- **WHEN** a run is read while its reviewer is running
- **THEN** it reports `reviewing`

### Requirement: Old runs are removed and the current one never is

socrates SHALL remove run directories older than a retention age, and SHALL never remove a run that is not terminal.
Removal SHALL happen when a new run starts, so that no background process is needed for it.

#### Scenario: An old run is swept

- **WHEN** a run older than the retention age exists and a new run starts
- **THEN** the old run's directory is gone and the new run's is not

#### Scenario: A live run is never swept

- **WHEN** a run older than the retention age is still `reviewing` and a new run starts
- **THEN** its directory remains

### Requirement: A run holds one of seven states and four of them are terminal

A run SHALL be in exactly one of `starting`, `reviewing`, `awaiting-response`, `clean`, `escalated`, `failed`, `cancelled`.
`clean`, `escalated`, `failed` and `cancelled` SHALL be terminal.

- `awaiting-response`: a review completed, what it found can be answered, and socrates is waiting for the author to fix, to argue, or to bring a question to the person and record what they decided.
- `clean`: a review completed and reported nothing blocking.
- `escalated`: a review completed and socrates has nothing further it can do, because what it found can be answered by nobody or because the round cap ran out.
- `failed`: no review completed.
- `cancelled`: the person stopped it.

`failed` SHALL NOT be merged with `escalated`.
They call for different actions: escalated means read the findings, failed means nothing has been checked.

`awaiting-response` SHALL NOT be terminal and SHALL NOT be reported as a verdict.
A run parked on the author is a run that has not concluded, and reporting it as one would let a session stop at the point where the dialogue was about to happen.
A run parked on a question the person must answer is parked on the same terms: the dialogue has not happened yet, and what is waited on is a person rather than a patch.

#### Scenario: A terminal run does not move

- **WHEN** a run has reached a terminal state
- **THEN** no later write changes its state or its verdict

#### Scenario: A failed run reports no findings as its verdict

- **WHEN** a run ends as `failed`
- **THEN** its recorded verdict is the failure and its reason, and it does not present an empty findings set as a result

#### Scenario: A parked run is not finished

- **WHEN** a run is `awaiting-response`
- **THEN** commands reporting on it exit with the code meaning still working, and no verdict is printed

#### Scenario: A run parked on a question is not a verdict either

- **WHEN** a run is `awaiting-response` because a finding asks the person a question
- **THEN** it is reported as still working rather than as escalated, and no verdict is printed

#### Scenario: A parked run can still be cancelled

- **WHEN** a run is `awaiting-response` and is cancelled
- **THEN** it becomes `cancelled`, its snapshot worktree is removed, and its rounds' artifacts are kept

### Requirement: A parked run is not a stalled one

The rule that a run recorded as working with no live reviewer is reported as failed SHALL NOT apply to `awaiting-response`.

A parked run is waiting on a person, and no process should exist for it.
Applying the liveness rule to it would report every parked run as failed the moment its reviewer exited, which is immediately.

A run that has been parked longer than the retention age SHALL still be swept only under the rule that never removes a run which is not terminal, so a parked run is not removed while it waits.

A parked run SHALL record no process identifiers, because none of its processes is alive and a run can wait indefinitely.
Cancelling signals what the run records, and a number that outlived its process names whatever the machine gave it next.

A round's worktree SHALL be torn down before the run is published as parked.
The next round checks out the same path, so a run that reads parked while its worktree is still being removed hands the author a command that fails on the state socrates itself left behind.

A failure to tear it down SHALL NOT stop the run being published as parked, and SHALL be recorded against the round.
A review that completed and cannot be reached by any command is worse than a directory left behind.
The findings are on disk, `--wait` would poll a run that never leaves `reviewing`, and `scrat cancel` would be the only way out of a review that succeeded.

The next round SHALL clear whatever stands where its worktree goes before it checks one out.
Otherwise the dead end is only deferred: the round parks, the author answers, and the next round refuses on the leftover and ends the run failed, which is the same completed review lost one command later.

#### Scenario: A leftover worktree does not stop the next round

- **WHEN** the next round starts and the previous round's worktree is still on disk
- **THEN** it is cleared, the round runs, and the run is not failed for it

#### Scenario: A parked run has nothing left to tear down

- **WHEN** a run first reads `awaiting-response`
- **THEN** the round's worktree is already gone

#### Scenario: A parked run holds no process identifiers

- **WHEN** a run parks on the author
- **THEN** its record holds no reviewer and no agent process identifier

#### Scenario: A parked run with no processes is still parked

- **WHEN** a run is `awaiting-response` and no process from it is running
- **THEN** it is reported as awaiting a response, not as failed

#### Scenario: A parked run survives the sweep

- **WHEN** a run has been parked for longer than the retention age and a new run starts
- **THEN** its directory remains

