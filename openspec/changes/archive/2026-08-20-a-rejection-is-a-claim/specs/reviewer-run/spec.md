## MODIFIED Requirements

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

- **WHEN** `scrat review` is invoked in a process whose environment carries the marker
- **THEN** it refuses, names the outer run, creates no run directory, and exits with the usage code

#### Scenario: A round inside a review is refused

- **WHEN** `scrat continue` is invoked on a parked run in a process whose environment carries the marker
- **THEN** it refuses, names the outer run, starts no round, leaves the run parked, and exits with the usage code
