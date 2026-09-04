## MODIFIED Requirements

### Requirement: A run holds its state, its rounds and its log

A run directory SHALL hold:

- `run.json`: the repository root, the base commit, the snapshot commit, the stated intent, when it started, its current state, and its verdict once it has one.
- `events.jsonl`: the append-only event log.
- one directory per round, named `r<n>` from `r1`, holding that round's prompt, the reviewer's output stream and its findings.

Nothing else in a run directory SHALL be load-bearing, so that a run is legible with `ls` and `cat` at the moment socrates itself is what is broken.

#### Scenario: A started run is on disk immediately

- **WHEN** `socrates-review review` returns
- **THEN** the run directory exists, holds a `run.json` naming the repository and the snapshot commit, and holds `r1`

#### Scenario: Round output is under the round

- **WHEN** a reviewer has run
- **THEN** `r1` holds the exact prompt that was sent, the reviewer's output stream, and its findings, each as its own file

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
The findings are on disk, `--wait` would poll a run that never leaves `reviewing`, and `socrates-review cancel` would be the only way out of a review that succeeded.

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
