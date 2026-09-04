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

### Requirement: The reviewer never writes, and socrates proves it

socrates SHALL hold this with three mechanisms, named here in the order they act, each weaker than the one after it:

1. The vendor's own tool restriction, which refuses a write the model asks for. It covers the built-in tools the model may call and nothing else, so it does not cover a tool from any other source, nor a hook or any other process the vendor's session starts.
2. The seal on the snapshot worktree, which refuses a write from anything in that session, whatever asked for it. It is what makes the property true rather than requested.
3. The comparison of the snapshot worktree against the snapshot commit when the round ends. It is what catches a write that defeated the seal.

All three bound the snapshot worktree and none of them bounds anything else.
A tool that reaches outside the worktree is therefore not covered by any of them, and what covers it is that the reviewer holds the tools socrates asked for and no others.
That is why the tool surface is a mechanism of its own rather than left to the seal, why what the session actually held is read rather than assumed, and why a run that held more records that it did.

When a round ends, socrates SHALL compare the snapshot worktree against the snapshot commit.
Any difference SHALL fail the round, naming the paths that differ, and SHALL leave the worktree in place.

A reviewer with write access will eventually decide to just fix something, after which the next round's diff contains changes nobody attributes.
The first two mechanisms can be taken away from socrates by the machine it runs on, and were.
A customisation the machine brings can put a tool in the session that the vendor's restriction was never about, and a filesystem can decline to enforce permissions.
The comparison cannot be taken away, which is why it stays even though the seal makes it unreachable in ordinary use.

#### Scenario: A reviewer that wrote fails the round

- **WHEN** a file in the snapshot worktree has been modified by the time the round ends, whatever defeated the seal
- **THEN** the run ends as `failed`, names the modified path, keeps the worktree, and reports no findings as a result

#### Scenario: What the mechanisms do not cover

- **WHEN** the three mechanisms are examined for what they bound
- **THEN** each of them bounds the snapshot worktree, and what keeps a tool that reaches past it away from the reviewer is that the reviewer holds the tools socrates asked for and no others

#### Scenario: A reviewer that created a file fails the round

- **WHEN** a file has been created in the snapshot worktree by the time the round ends, whatever defeated the seal
- **THEN** the run ends as `failed` and names the created path

#### Scenario: A clean reviewer passes the check

- **WHEN** the reviewer reads freely and writes nothing
- **THEN** the check passes and the round's verdict is whatever the findings say

#### Scenario: A write the seal refused is not a finding against the reviewer

- **WHEN** something in the reviewer's session attempted a write, the seal refused it, and the reviewer returned findings anyway
- **THEN** the worktree still matches the snapshot commit, the comparison passes, and the verdict is whatever the findings say

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

