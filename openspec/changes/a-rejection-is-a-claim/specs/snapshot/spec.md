## ADDED Requirements

### Requirement: Every round takes a fresh snapshot

Each round SHALL snapshot the working tree again at the moment it starts, and the reviewer for that round SHALL read a worktree of that round's snapshot.

The reviewer of a later round must see the patches as commits it can diff, not as a description of what changed.
A description of a patch is a claim about a patch, and the whole point of the round is to judge the claim against the code.

Each round's snapshot commit SHALL be recorded against that round, so that what a round reviewed is answerable afterwards.

#### Scenario: The second round sees the patches

- **WHEN** the author patches a file and the next round starts
- **THEN** that round's snapshot contains the patched contents, and the reviewer's worktree is a worktree of it

#### Scenario: Each round records what it reviewed

- **WHEN** a run completes two rounds
- **THEN** each round records its own snapshot commit, and the two differ where the working tree changed between them

#### Scenario: A round with nothing patched still snapshots

- **WHEN** every finding of a round was rejected and nothing was changed
- **THEN** the next round still takes a snapshot, and it is the reviewer's worktree even though its tree matches the previous round's

### Requirement: Every round is reviewed against the run's original base

The base commit a round is reviewed against SHALL be the run's original base, not the previous round's snapshot.

Reviewing a later round against the previous round's snapshot shows only the patches, which reads them out of the context they were applied to.
A patch that is wrong because of what it sits next to is exactly the kind this layer exists to catch, and it is invisible in a diff of the patch alone.

#### Scenario: The base does not move between rounds

- **WHEN** a run completes two rounds
- **THEN** both rounds record the same base commit, and it is the commit `HEAD` resolved to when the run started

#### Scenario: The second round sees the whole change

- **WHEN** the second round's prompt carries a diff
- **THEN** it is the diff between the original base and that round's snapshot, containing the author's original change as well as the patches

#### Scenario: The base still exists when the person has moved on

- **WHEN** the person commits or amends on their branch while a run is awaiting a response
- **THEN** the run's recorded base commit is still what later rounds are diffed against, and a base that has become unreachable fails the round naming it rather than silently reviewing against something else
