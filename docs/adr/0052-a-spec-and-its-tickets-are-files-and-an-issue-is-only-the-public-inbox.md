# 52. A spec and its tickets are files, and an issue is only the public inbox

- Status: accepted
- Date: 2026-10-01
- Part of: specs-and-tickets-are-files

## Decision

A spec lives at `docs/changes/<slug>/spec.md` and its tickets under `docs/changes/<slug>/tickets/`, committed, and neither is ever published as an issue.
`socrates-design` offers no publishing in any repository: its last card offers write it, whatever a repository states about a tracker, and nothing that starts another skill, per ADR 49.
In this repository an issue is where somebody using socrates reports a bug or asks for something, and nothing else; a request that is taken up becomes a spec in a change folder.
`docs/agents/issue-tracker.md` is gone, and `CONTRIBUTING.md` names the spec's path where it named a GitHub issue.

The person stated it after a spec was published as #71, in their own words: "We don't rely on github issue anymore because not every project uses github. We already write spec, adr, everything is file based, commited and pushed to git."

The surprise is a project hosted on GitHub whose own specs and tickets are not issues, while its issues stay open to everybody else.
What is bought is one place a change lives, which a worktree reads and a project on any forge, or on none, has.
This supersedes the part of ADR 39 that gave the tracker the tickets, which ADR 43 had already contradicted for the tickets `socrates-tickets` writes.

## Rejected

- Publishing stays offered on the last card where a repository states a tracker, but never recommended.
  Declined without a word.
- `docs/agents/issue-tracker.md` kept and rewritten to say specs and tickets are files.
  Declined without a word.
- Issues turned off, with `SUPPORT.md`, the issue templates and every link to an issue gone.
  Declined without a word.
