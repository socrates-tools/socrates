# snapshot Specification

## Purpose
TBD - created by archiving change a-review-nobody-authored. Update Purpose after archive.
## Requirements
### Requirement: A snapshot captures the whole working tree

socrates SHALL freeze the repository's current state into a commit before any reviewer is spawned, and that commit SHALL contain every change the person has made, whether or not the change is tracked.

The snapshot SHALL be built through a temporary index created from `HEAD`, so that a path is staged into that index and never into the person's own.
It SHALL honour `.gitignore`, so a snapshot is what the repository would contain rather than what the directory happens to hold.

`git stash create` SHALL NOT be used.
It silently ignores untracked files, and a freshly implemented feature is mostly untracked files, so a snapshot taken that way is a confident review of half a change.

#### Scenario: An untracked file is in the snapshot

- **WHEN** a repository has a new file that has never been added, and a review is started
- **THEN** the snapshot commit's tree contains that file with its current contents

#### Scenario: A modified tracked file is in the snapshot

- **WHEN** a tracked file has uncommitted modifications, and a review is started
- **THEN** the snapshot commit's tree contains the modified contents, not the contents at `HEAD`

#### Scenario: A staged change and an unstaged change are both in the snapshot

- **WHEN** one file is staged and another is modified without being staged, and a review is started
- **THEN** the snapshot commit's tree contains both, because the snapshot is of the working tree rather than of either index

#### Scenario: A deleted file is deleted in the snapshot

- **WHEN** a tracked file has been removed from the working tree, and a review is started
- **THEN** the snapshot commit's tree does not contain it

#### Scenario: An ignored file is not in the snapshot

- **WHEN** the working tree holds a file matched by `.gitignore`, and a review is started
- **THEN** the snapshot commit's tree does not contain it

#### Scenario: A change with nothing to review

- **WHEN** a review is started in a repository whose working tree is identical to `HEAD`
- **THEN** the run ends immediately with the `clean` verdict, records that there was nothing to review, and spawns no reviewer

### Requirement: A snapshot changes nothing the person can see

Taking a snapshot SHALL leave the person's index, working tree, `HEAD`, current branch, stash and reflog exactly as they were.
It SHALL create no ref, no tag and no branch pointing at the snapshot commit, so that ordinary garbage collection reclaims it once the worktree is gone.

#### Scenario: Git status is unchanged

- **WHEN** a snapshot is taken
- **THEN** the output of `git status --porcelain=v2 --branch --untracked-files=all` is byte-identical before and after

#### Scenario: The index is unchanged

- **WHEN** a snapshot is taken in a repository with a partially staged file
- **THEN** the person's index file is byte-identical afterwards, and the staged and unstaged halves of that file are still separately staged and unstaged

#### Scenario: No ref names the snapshot

- **WHEN** a snapshot is taken
- **THEN** no ref under `refs/` resolves to the snapshot commit, and the reflog has no new entry

#### Scenario: The snapshot's parent is the current head

- **WHEN** a snapshot is taken
- **THEN** the snapshot commit has exactly one parent and it is the commit `HEAD` resolved to when the snapshot was taken

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

### Requirement: A repository socrates cannot snapshot is refused

socrates SHALL refuse to start a run outside a git repository, and in a repository whose `HEAD` resolves to nothing.
Each refusal SHALL name the condition and SHALL exit with the usage code rather than a verdict, because nothing was reviewed and no run exists to report on.

#### Scenario: Not a repository

- **WHEN** `scrat review` is run in a directory that is not inside a git working tree
- **THEN** it refuses, saying so, creates no run directory, and exits with the usage code

#### Scenario: A repository with no commits

- **WHEN** `scrat review` is run in a repository whose `HEAD` resolves to no commit
- **THEN** it refuses, saying that there is no base to review against, creates no run directory, and exits with the usage code

#### Scenario: A bare repository

- **WHEN** `scrat review` is run with a bare repository as its working directory
- **THEN** it refuses, saying there is no working tree to review, and exits with the usage code

### Requirement: A snapshot in progress leaves nothing behind when it fails

Where any step of taking a snapshot fails, socrates SHALL remove the temporary index it created and any worktree it registered, and SHALL report the failing git invocation.
A partially registered worktree left in `.git/worktrees` is a repository the person has to repair by hand, which is a worse outcome than the review not happening.

#### Scenario: Worktree creation fails

- **WHEN** checking out the snapshot worktree fails
- **THEN** the run ends as `failed`, `git worktree list` names no worktree for this run, and the temporary index file does not exist

#### Scenario: The temporary index never joins the repository

- **WHEN** a snapshot is taken, successfully or not
- **THEN** no file was written inside the repository's `.git` directory other than the objects the snapshot commit needed and the worktree administrative entry

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
- **THEN** `scrat review` refuses, naming the worktree and both conditions, exits with the usage code, spawns no reviewer, and leaves nothing behind

#### Scenario: The seal is lifted when the round ends

- **WHEN** a run reaches any terminal state, whether its worktree is removed or kept for inspection
- **THEN** nothing is left sealed, so a person can read and delete what remains without changing permissions first

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

