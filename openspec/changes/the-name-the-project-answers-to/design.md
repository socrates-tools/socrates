## Context

See `proposal.md` for why.

The module is entirely `internal/`.
Nothing outside this repository can import any of it, so no consumer breaks, no compatibility shim is needed, and the rename is safe in a way most module renames are not.
What makes it delicate is not risk but size: 448 references in one commit, where any unrelated edit riding along makes the diff unreviewable.

`Makefile` composes the version ldflags from `MODULE`, so the path appears there as data rather than as an import.
`openspec/config.yaml` names the old binary in prose that no compiler checks.

The history carries two author addresses, one personal and one an employer's.
The decision recorded in `proposal.md` is to rewrite to the personal address before the first public push.

## Goals / Non-Goals

**Goals:**

- One commit that changes paths and nothing else, verifiable by the suite passing unchanged before and after.
- Nothing in the source names an account, so the next rename is a check failing rather than a grep somebody remembers to run.
- The normative context file describes the project that exists.

**Non-Goals:**

- A compatibility path, an alias module, or a `retract` directive. Nothing external imports this, so there is nothing to be compatible with.
- Rewriting the archive. See the decision below.
- Renaming the binaries, the commands, the home directory or the configuration file. `socrates`, `scrat` and `agora` are unaffected: this changes where the source lives, not what the tool is called.

## Decisions

### The rename rides alone

The commit contains the module path, the imports, `MODULE`, `openspec/config.yaml`, the README's path references, and the new convention check.
Nothing else.

The reason is reviewability, and it is the reason socrates itself exists.
A 181-file diff in which 180 files changed one line and one file changed its behaviour is a diff where the behaviour change is invisible, to a human reader and to a reviewer alike.
`make verify` passing identically before and after is what makes the claim "this changed no behaviour" checkable rather than asserted.

### The archive is left alone

Two archived documents name `github.com/rui-ktei/socrates`, and one links an issue in the old repository.

They are left as they are.
An archived change is not documentation that has gone stale; it is the record of a decision, and the path it names is the path that existed when the decision was taken.
Rewriting it produces a document that claims a past author wrote something they did not write, which is worse than a document that is out of date, because a reader cannot tell the difference.

The consequence accepted is that the repository will contain the old account name, in the archive, forever.
The convention check is therefore scoped to source and to live specifications, and skips `openspec/changes/archive` explicitly, in the same way `internal/coverage` already skips it.

The issue link in `a-round-a-person-reads-becomes-a-page/design.md` follows the repository: transferring a repository moves its issues and redirects their URLs, so the link keeps resolving without being edited.

### The check is on the account, not on the path

`internal/convention` gains a check that no source file names a GitHub account other than the organisation.

Checking for the specific old string would pass the day somebody introduces a third account.
Checking that the only account named is the organisation's fails for anything unexpected, which is the property worth having, and it is the same shape as the check that every user-facing string lives in its package's `messages.go`.

### The transfer, not a fresh repository

The repository is transferred to the organisation rather than recreated there.

A transfer carries the issues, and issue 2 is an open question this project still owes an answer to.
It also redirects the old URLs, so the archived link keeps working.
A fresh repository would lose both, and would buy nothing, because the history is being rewritten either way for the author addresses.

The order is: rewrite the authorship, then rename the module, then transfer, then push.
Rewriting after the transfer would mean force-pushing over a repository somebody could already have cloned.

## Risks / Trade-offs

**A rewritten history invalidates every existing clone and worktree.**
This repository is developed across git worktrees that share a stash stack and an object store.
Every one of them has to be re-created from the rewritten history, and any work in flight has to be replayed rather than merged.
That is the price of fixing authorship, and it is only payable before publication.

**The rename and the authorship rewrite are two operations on the same history.**
They are kept separate and ordered: the rewrite is a history operation on the commits that exist, the rename is a new commit on top.
Combining them would put a 448-file change inside a filter pass, where nothing can review it.

**A path can re-enter through a generated file.**
`internal/skill` generates the committed skill and `web/build.mjs` generates the bundle.
Neither carries the module path today, and the convention check covers the generated copies as source, so one arriving later fails rather than passing quietly.
