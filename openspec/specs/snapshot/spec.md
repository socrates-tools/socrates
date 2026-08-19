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
- **THEN** the worktree is left in place, the run records where it is, and `git worktree list` still names it

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

