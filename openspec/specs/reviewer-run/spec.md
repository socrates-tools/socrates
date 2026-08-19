# reviewer-run Specification

## Purpose
TBD - created by archiving change a-review-nobody-authored. Update Purpose after archive.
## Requirements
### Requirement: The reviewer runs detached and the caller does not wait

`scrat review` SHALL spawn the reviewer as a process that survives the invoking command, record its identity in the run, and return.
The session that asked for the review goes back to work; that is the property the whole tool is arranged around.

The spawned process SHALL NOT hold the invoking terminal's standard streams, so that a session which exits does not close the reviewer's output.

#### Scenario: The command returns while the review runs

- **WHEN** `scrat review` is invoked
- **THEN** it returns before the review has finished, having printed the run identifier, and the reviewer process is still running

#### Scenario: The reviewer survives its caller

- **WHEN** the process that invoked `scrat review` exits
- **THEN** the reviewer keeps running and the run reaches a terminal state

### Requirement: The reviewer and everything it starts are reaped together

The reviewer SHALL be started in its own process group, and socrates SHALL terminate that whole group when the round ends, whether it ended cleanly, in error, or by cancellation.

Workers a reviewer spawned otherwise escape, hold the snapshot worktree open, and accumulate until the operating system kills something.
Reaping only on cancellation is the version of this that looks correct and is not: the ordinary exit is the common path.

#### Scenario: A child that outlives the agent is killed

- **WHEN** the reviewer spawns a child that would outlive it, and the round ends cleanly
- **THEN** no process from the reviewer's group is running afterwards

#### Scenario: Reaping happens on error too

- **WHEN** the round ends because the agent call failed
- **THEN** no process from the reviewer's group is running afterwards

#### Scenario: Reaping happens on cancellation

- **WHEN** the run is cancelled while the reviewer is working
- **THEN** no process from the reviewer's group is running afterwards, and the run is `cancelled`

#### Scenario: The worktree can be removed after the round

- **WHEN** a round ends and the worktree is removed
- **THEN** the removal succeeds, because no process is still holding a file in it

### Requirement: A reviewer cannot start a review

The reviewer's environment SHALL carry a marker naming the run it belongs to, and `scrat review` SHALL refuse when its own environment already carries one, naming the outer run.

The reviewer is a coding session reading a repository whose agent instructions may tell it to run socrates.
The failure without the marker is recursive process creation, bounded only by the machine.

#### Scenario: The reviewer's environment is marked

- **WHEN** a reviewer is spawned
- **THEN** its recorded environment carries the marker, and the marker names the run

#### Scenario: A review inside a review is refused

- **WHEN** `scrat review` is invoked in a process whose environment carries the marker
- **THEN** it refuses, names the outer run, creates no run directory, and exits with the usage code

### Requirement: The reviewer never writes, and socrates proves it

When a round ends, socrates SHALL compare the snapshot worktree against the snapshot commit.
Any difference SHALL fail the round, naming the paths that differ, and SHALL leave the worktree in place.

The vendor's own tool restrictions are the first line and are not what this requirement rests on.
A rule that depends on a third party's flag name is a rule that quietly stops holding, and a reviewer with write access will eventually decide to just fix something, after which the next round's diff contains changes nobody attributes.

#### Scenario: A reviewer that wrote fails the round

- **WHEN** the reviewer modifies a file in the snapshot worktree and returns findings
- **THEN** the run ends as `failed`, names the modified path, keeps the worktree, and reports no findings as a result

#### Scenario: A reviewer that created a file fails the round

- **WHEN** the reviewer creates a file in the snapshot worktree
- **THEN** the run ends as `failed` and names the created path

#### Scenario: A clean reviewer passes the check

- **WHEN** the reviewer reads freely and writes nothing
- **THEN** the check passes and the round's verdict is whatever the findings say

### Requirement: A review that did not happen is never a pass

`clean` SHALL require positive evidence that a review completed: a reviewer that exited successfully, a findings document that parsed against the schema, and a worktree that was not written to.

The absence of findings SHALL NOT by itself be treated as a clean review.
A findings file that never appeared, a reviewer that died, a schema violation and an empty output stream SHALL each end the run as `failed`, naming which.

#### Scenario: A missing findings document is a failure

- **WHEN** the reviewer exits successfully but writes no findings document
- **THEN** the run ends as `failed` naming the missing findings, not `clean`

#### Scenario: A dead reviewer is a failure

- **WHEN** the reviewer is killed part-way through
- **THEN** the run ends as `failed` naming that the reviewer went away, not `clean`

#### Scenario: An empty findings array with evidence is clean

- **WHEN** the reviewer exits successfully, writes a findings document with an empty findings array that satisfies the schema, and wrote nothing to the worktree
- **THEN** the run ends as `clean`

#### Scenario: An unparseable findings document is a failure

- **WHEN** the findings document does not parse, or parses but violates the schema
- **THEN** the run ends as `failed`, the document is kept for inspection, and the failure names the violation

### Requirement: The reviewer's stream is recorded as it arrives

Everything the reviewer emits SHALL be appended to the round's output file as it arrives, rather than buffered until the round ends.
A run that is still working is a run somebody may want to look at, and a run that failed is one whose last output is the evidence.

#### Scenario: Output is visible while the review is running

- **WHEN** the reviewer has emitted several events and has not finished
- **THEN** those events are already on disk in the round's output file

#### Scenario: A failed round keeps what it got

- **WHEN** a round fails part-way through
- **THEN** everything the reviewer emitted before the failure is on disk

