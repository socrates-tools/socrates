## ADDED Requirements

### Requirement: The snapshot worktree cannot be written to

socrates SHALL seal the snapshot worktree after checking it out: every directory and every file in it SHALL be left without write permission for as long as the round runs.
A write attempted by a reviewer, a tool server, a hook, or anything else a vendor's session starts SHALL therefore be refused by the filesystem rather than trusted, restricted by a flag, or caught afterwards.

The seal is not a sandbox and SHALL NOT be described as one.
The session runs as the owner of those files, so something that sets out to defeat the seal can clear the bit back and write, which is what the comparison after the round is for.
What the seal removes is every write nothing went out of its way to make, which is the whole of what has actually written in a snapshot worktree so far.

The seal SHALL cover the worktree and nothing else, because the run directory around it is where the prompt, the reviewer's stream and its findings are written while the round runs.

socrates SHALL prove the seal before a reviewer is spawned, by attempting one write inside the sealed worktree.
A `chmod` that returned no error is not evidence: a filesystem can ignore permission bits entirely.

Where that write succeeds, socrates SHALL remove what it wrote and refuse the review, naming the worktree and both ways the seal fails to hold: a filesystem that does not enforce permissions, and a process that is root and ignores them.
The refusal SHALL name the next action for each, because the remedy for one is not the remedy for the other.
A refusal is the only honest answer there, because socrates cannot otherwise promise that the change reviewed is the change that was frozen, and every verdict it produced would be a claim it could not support.

#### Scenario: A reviewer cannot create a file

- **WHEN** anything running in the reviewer's session attempts to create a file in the snapshot worktree
- **THEN** the attempt fails with a permission error and the worktree still matches the snapshot commit

#### Scenario: A reviewer cannot modify a file

- **WHEN** anything running in the reviewer's session attempts to modify a file that the snapshot contains
- **THEN** the attempt fails with a permission error and the file's contents are the snapshot's

#### Scenario: The reviewer can still read everything

- **WHEN** a reviewer reads, greps and lists the sealed worktree
- **THEN** every file the snapshot contains is readable and the review proceeds normally

#### Scenario: A seal that does not hold is a refusal

- **WHEN** a write into the sealed worktree succeeds, whether because the filesystem enforces nothing or because the process is root
- **THEN** `scrat review` refuses, naming the worktree and both conditions, exits with the usage code, spawns no reviewer, and leaves nothing behind

#### Scenario: The seal is lifted when the round ends

- **WHEN** a run reaches any terminal state, whether its worktree is removed or kept for inspection
- **THEN** nothing is left sealed, so a person can read and delete what remains without changing permissions first

## MODIFIED Requirements

### Requirement: The reviewer reads a worktree of the snapshot

socrates SHALL check the snapshot commit out as a detached worktree under the run directory, and that worktree SHALL be what a reviewer is given as its working directory.
The reviewer SHALL NOT be given the person's own repository directory, because a reviewer that can see a tree still moving cannot anchor a finding to a line.

That worktree SHALL be sealed against writing while the round runs, and the seal SHALL be lifted before it is removed or handed to a person to inspect.

The worktree SHALL be removed when the run reaches a terminal verdict, and SHALL be kept when the run failed because the worktree was written to, so that the difference can be inspected.

#### Scenario: The worktree is the snapshot

- **WHEN** a reviewer is spawned
- **THEN** its working directory is a worktree whose `HEAD` is the snapshot commit, detached, and whose contents match the snapshot tree

#### Scenario: Editing while the reviewer reads

- **WHEN** the person modifies a file in their own working tree after the reviewer has been spawned
- **THEN** the reviewer's worktree is unaffected and the findings it returns are anchored to the snapshot's contents

#### Scenario: The worktree is cleaned up

- **WHEN** a run reaches `clean` or `escalated`
- **THEN** the snapshot worktree is removed and `git worktree list` no longer names it

#### Scenario: A dirtied worktree is kept

- **WHEN** a run fails because the snapshot worktree no longer matches the snapshot commit
- **THEN** the worktree is left in place, unsealed, the run records where it is, and `git worktree list` still names it
