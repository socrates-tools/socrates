## MODIFIED Requirements

### Requirement: A repository socrates cannot snapshot is refused

socrates SHALL refuse to start a run outside a git repository, and in a repository whose `HEAD` resolves to nothing.
Each refusal SHALL name the condition and SHALL exit with the usage code rather than a verdict, because nothing was reviewed and no run exists to report on.

#### Scenario: Not a repository

- **WHEN** `socrates-review review` is run in a directory that is not inside a git working tree
- **THEN** it refuses, saying so, creates no run directory, and exits with the usage code

#### Scenario: A repository with no commits

- **WHEN** `socrates-review review` is run in a repository whose `HEAD` resolves to no commit
- **THEN** it refuses, saying that there is no base to review against, creates no run directory, and exits with the usage code

#### Scenario: A bare repository

- **WHEN** `socrates-review review` is run with a bare repository as its working directory
- **THEN** it refuses, saying there is no working tree to review, and exits with the usage code

### Requirement: The snapshot worktree cannot be written to

socrates SHALL seal the snapshot worktree after checking it out: every directory and every file in it SHALL be left without write permission for as long as the round runs.
A write attempted by a reviewer, a hook, or anything else a vendor's session starts SHALL therefore be refused by the filesystem rather than trusted, restricted by a flag, or caught afterwards.

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
- **THEN** `socrates-review review` refuses, naming the worktree and both conditions, exits with the usage code, spawns no reviewer, and leaves nothing behind

#### Scenario: The seal is lifted when the round ends

- **WHEN** a run reaches any terminal state, whether its worktree is removed or kept for inspection
- **THEN** nothing is left sealed, so a person can read and delete what remains without changing permissions first
