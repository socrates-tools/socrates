# 51. The default branch is remembered per repository, in socrates implement's own file

- Status: accepted
- Date: 2026-10-01
- Part of: socrates-implement

## Decision

socrates never assumes which branch a repository's features start from.
The overseer asks once per repository, proposing the branches named `main`, `master`, `dev` or `develop` and never one like `feature/*` or `hotfix/*`, and saves the answer with `socrates implement default-branch <branch>`.
It is remembered under `default_branches` in `~/.socrates/implement/config.yaml`, keyed by the top level of the repository's main worktree, so every worktree of one clone reads one entry, and a moved or second clone is asked again.

The person stated the problem in their own words: "socrates should not assume the default branch for the repo (because sometimes a default branch might be main, but usually developers branch off rom dev branch for their work). If there is no configured default branch (configured by scorates), socrates should ask user to pick one [...] once user tells AI the default branch, socrates should remember it/configure it so next time nothing needs to be asked to user anymore."

The default branch is what the change is measured against: the overseer's read of the whole change diffs from its fork point with the change's branch.
The change's branch is still cut from the current `HEAD`, where the spec and the tickets were committed.

The word is default branch, as socrates means it, which need not be the branch a forge calls its default.
The person kept it knowing that: "This is okay, in this case even they're different, dev branch is not the actual "default" branch, but it is the "base" branch to start a feature's implementation. This is for socrates and if this default confuses with the "default" in git's world, we can rename this concept."

## Rejected

- The clone's own git config, under a `socrates` key in `.git/config`, written by a socrates implement command.
  Declined without a word.
- A file committed in the repository, at `.socrates/implement.yaml`, so everyone who clones it inherits the answer.
  Declined without a word.
- The change's branch cut from the default branch's tip rather than from the current `HEAD`.
  Declined without a word.
- Base branch as the word, beside review's base, which is a commit rather than a branch.
  Declined, with the words quoted above.
- Trunk as the word.
  Declined without a word.
