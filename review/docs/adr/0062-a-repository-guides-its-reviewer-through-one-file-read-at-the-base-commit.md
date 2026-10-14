# 62. A repository guides its reviewer through one file, read at the base commit

- Status: accepted
- Date: 2026-10-14
- Part of: the conventions the project brings, recorded after it shipped

## Decision

A repository tells its reviewer what it would otherwise get wrong by committing `.socrates/review-guide.md`.
review reads that one file, at the run's base commit, and carries it into one section of the prompt, "What this project says about its own code", after review's own rules and before the reporting contract.

The base commit is the trust anchor.
A forge reads such a file from the default branch, so that a pushed branch cannot write the rules that gate it.
review has no forge and no default branch, but every run has a base, and the base is by construction the last state no part of the change under review authored.
So a guide the change writes governs the next run rather than this one, and every round of a run is reviewed against the guide its first round had.
The threat is not a stranger.
It is the session that wrote the change, which review already keeps away from the reviewer's ear, and which is one keystroke from writing "this is fine" into a guide read from the working tree.

The guide adds and removes nothing.
It cannot switch a rule off, exempt a file or change what `clean` means, because an independent reviewer the reviewed party can quiet is not one.
That is held by what review does not build: no code reads a word of the guide, so no severity, action, verdict, state or exit code is derived from it.
What is asked of the model is only the framing sentence in the section.

One fixed path, which nothing configures, so a repository cannot point review at a file its own tooling generates.
The guide is neutralised on the way into the prompt like any other quoted text, so it cannot open a section of its own.
It is held to a ceiling, 16 KiB by default, which is `guidance_ceiling` in the person's own file and never the repository's, and which a run resolves once when it begins.
Over the ceiling the guide is not carried and the review still runs, and the person is told the size, the ceiling and that the next review carries none of it.
A run records whether a guide was carried, whether there was none, whether it was over the ceiling, and whether the change carries a guide the base does not, so the base-commit rule teaches itself the first time it applies.

## Rejected

- Reading the guide from the working tree or the snapshot.
  It hands the change a channel to the reviewer that gates it, for no gain beyond one run's latency.
- A guide that overrides review's rules.
  The first rule written under a deadline is the one removing the finding somebody did not want.
  If a waiver is ever genuinely needed, it is a named list of rules a repository may waive, argued for on its own, and not prose that outranks the prompt.
- Guidance scoped to globs and selected against the changed files.
  Elsewhere it defeats a contributor adding a rule's glob to an ignore list, and review has no ignore list, so what is left is attention, which the ceiling already bounds.
- Truncating a guide over the ceiling.
  It would deliver part of a maintainer's rules while telling the reviewer they are the rules.
- Refusing the review over an oversized guide.
  The guide is read at the base, so every review in that repository would be refused, including the one that shrinks it.
- Guidance kept in the person's own file, keyed by repository.
  A coherent feature that answers another question, because it does not arrive with a clone, and nothing here forecloses it.
