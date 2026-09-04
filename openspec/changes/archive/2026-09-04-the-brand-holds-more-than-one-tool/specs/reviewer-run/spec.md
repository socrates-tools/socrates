## MODIFIED Requirements

### Requirement: The reviewer runs detached and the caller does not wait

`socrates-review review` SHALL spawn the reviewer as a process that survives the invoking command, record its identity in the run, and return.
The session that asked for the review goes back to work; that is the property the whole tool is arranged around.

The spawned process SHALL NOT hold the invoking terminal's standard streams, so that a session which exits does not close the reviewer's output.

#### Scenario: The command returns while the review runs

- **WHEN** `socrates-review review` is invoked
- **THEN** it returns before the review has finished, having printed the run identifier, and the reviewer process is still running

#### Scenario: The reviewer survives its caller

- **WHEN** the process that invoked `socrates-review review` exits
- **THEN** the reviewer keeps running and the run reaches a terminal state

### Requirement: A reviewer cannot start a review

The reviewer's environment SHALL carry a marker naming the run it belongs to, and every command that starts a reviewer SHALL refuse when its own environment already carries one, naming the outer run.

The reviewer is a coding session reading a repository whose agent instructions may tell it to run socrates.
The failure without the marker is recursive process creation, bounded only by the machine.

The refusal SHALL be a property of starting a reviewer rather than of one command.
A run can be advanced as well as begun, and a reviewer that can advance a parked run can seat a reviewer of its own with everything the marker exists to prevent.

#### Scenario: The reviewer's environment is marked

- **WHEN** a reviewer is spawned
- **THEN** its recorded environment carries the marker, and the marker names the run

#### Scenario: A review inside a review is refused

- **WHEN** `socrates-review review` is invoked in a process whose environment carries the marker
- **THEN** it refuses, names the outer run, creates no run directory, and exits with the usage code

#### Scenario: A round inside a review is refused

- **WHEN** `socrates-review continue` is invoked on a parked run in a process whose environment carries the marker
- **THEN** it refuses, names the outer run, starts no round, leaves the run parked, and exits with the usage code
