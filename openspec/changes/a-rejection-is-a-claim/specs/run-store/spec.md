## MODIFIED Requirements

### Requirement: A run holds one of six states and three of them are terminal

A run SHALL be in exactly one of `starting`, `reviewing`, `awaiting-response`, `clean`, `escalated`, `failed`, `cancelled`.
`clean`, `escalated`, `failed` and `cancelled` SHALL be terminal.

- `awaiting-response`: a review completed, what it found is mechanical, and socrates is waiting for the author to fix or to argue.
- `clean`: a review completed and reported nothing blocking.
- `escalated`: a review completed and what it found needs the person.
- `failed`: no review completed.
- `cancelled`: the person stopped it.

`failed` SHALL NOT be merged with `escalated`.
They call for different actions: escalated means read the findings, failed means nothing has been checked.

`awaiting-response` SHALL NOT be terminal and SHALL NOT be reported as a verdict.
A run parked on the author is a run that has not concluded, and reporting it as one would let a session stop at the point where the dialogue was about to happen.

#### Scenario: A terminal run does not move

- **WHEN** a run has reached a terminal state
- **THEN** no later write changes its state or its verdict

#### Scenario: A failed run reports no findings as its verdict

- **WHEN** a run ends as `failed`
- **THEN** its recorded verdict is the failure and its reason, and it does not present an empty findings set as a result

#### Scenario: A parked run is not finished

- **WHEN** a run is `awaiting-response`
- **THEN** commands reporting on it exit with the code meaning still working, and no verdict is printed

#### Scenario: A parked run can still be cancelled

- **WHEN** a run is `awaiting-response` and is cancelled
- **THEN** it becomes `cancelled`, its snapshot worktree is removed, and its rounds' artifacts are kept

## ADDED Requirements

### Requirement: A parked run is not a stalled one

The rule that a run recorded as working with no live reviewer is reported as failed SHALL NOT apply to `awaiting-response`.

A parked run is waiting on a person, and no process should exist for it.
Applying the liveness rule to it would report every parked run as failed the moment its reviewer exited, which is immediately.

A run that has been parked longer than the retention age SHALL still be swept only under the rule that never removes a run which is not terminal, so a parked run is not removed while it waits.

#### Scenario: A parked run with no processes is still parked

- **WHEN** a run is `awaiting-response` and no process from it is running
- **THEN** it is reported as awaiting a response, not as failed

#### Scenario: A parked run survives the sweep

- **WHEN** a run has been parked for longer than the retention age and a new run starts
- **THEN** its directory remains
