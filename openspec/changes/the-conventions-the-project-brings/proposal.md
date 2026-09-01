## Why

The reviewer knows nothing about the project it is reviewing.

That is deliberate, and the reason is good: what the reviewer is told must not be written by the thing under review.
`--safe-mode` takes the machine's own customisations out of the session, the vendor is refused the repository's project and local settings, and the prompt is socrates' own text delivered identically everywhere.
The reviewer is therefore independent, and it is also ignorant of every convention the project holds itself to, which is the price nobody agreed to pay explicitly.

A project's conventions are real review criteria.
This repository writes no comments in source, states one sentence per line in Markdown, uses no em dash, and holds tests to a rule about what they may read.
A reviewer told none of that reports a missing comment as a finding and passes a comment that should not be there, and neither is a judgement about the code.
The gap is recorded already: `openspec/config.yaml` lists repository-scoped review guidance under what is deliberately deferred, and names the answer as a repository file socrates itself reads.

That answer does not fit the specification as written.
`openspec/specs/agent-seam/spec.md` carries the scenario `Project settings are not loaded`, whose outcome is that the reviewer's prompt contains none of the repository's instructions.
A file socrates reads out of the repository is the repository's instructions under that wording.
The roadmap note and the promoted requirement point in opposite directions, and one of them is wrong.
This change settles that, and builds the channel.

The thing the roadmap note never answered is where the trusted copy lives.
`no-mistakes` has one because it reads its configuration from the default branch rather than from the branch under review.
socrates has no forge, no push and no default branch, so the note's answer stops at "a repository file" and does not say which version of it.
It has one anyway, and it has had it since the first commit: the base commit.
Every run is measured against a base, every round is snapshotted afresh against the run's original base, and the base is by construction the last state no part of the change under review authored.

Resolves ktei/socrates#2.

## What Changes

- **socrates reads one file from the repository under review**, `.socrates/review-guide.md`, and delivers it to the reviewer as one named prompt section.
  It is the only file socrates reads out of a reviewed repository for the reviewer's benefit, and its path is fixed rather than configured.

- **It is read at the run's base commit, never from the working tree or the snapshot.**
  The guidance that reviews a change is the guidance as it stood before the change.
  A guide the change itself adds or edits governs the next run rather than this one, so the session that wrote the code cannot also write the rules it is judged by.
  Because every round is snapshotted against the run's original base, the guidance is fixed for the whole of a run and cannot move between rounds of a dialogue.

- **The guidance is additive and can remove nothing.**
  It adds what the project knows about itself.
  It does not disable a rule socrates delivers, does not change what a severity or an action means, does not touch the verdict rule or the round cap, and cannot suppress a finding.
  socrates' own text says so in the section, so a guide instructing otherwise is instructing the reviewer against what the prompt already told it.
  Anything else would make `clean` mean a different thing in each repository, and the first rule anybody writes under a deadline is the one that switches off the finding they do not want.

- **It is neutralised on the way into the prompt**, on the mechanism model output already goes through, because it is untrusted text entering a prompt from outside socrates.
  That is what stops a guide forging a section heading and reading as socrates' own instruction.

- **It is bounded**, and a guide over the ceiling refuses the review naming the file and the ceiling rather than delivering part of one.
  Half a maintainer's rules, silently, is the failure socrates refuses everywhere else.

- **What happened is recorded in the run**: that guidance was carried, that there was none, or that the change adds a guide the base commit does not have and it did not reach this reviewer.
  The third is the one that would otherwise be a silent surprise, because it is what a person sees the first time they write the file.

- **The `agent-seam` scenario is reworded** to say what it always meant: the repository's own agent instructions are not loaded, by the vendor, on the vendor's account.
  One file socrates reads itself, at a commit it chose, neutralised, additive and bounded, is a different thing and is now stated as one.

The diff ceiling's requirement that it applies to every repository with no per-repository override is untouched.
A repository does not configure socrates here.
It contributes text to one section, on socrates' terms, and every knob remains user-level.

## Capabilities

### New Capabilities

- `review-guidance`: the file, its fixed path, the base commit as its anchor, the ceiling, what the guidance may and may not do, and what the run records about it.

### Modified Capabilities

- `reviewer-prompt`: a new named section carrying the project's guidance, absent where there is none, placed after socrates' own rules and before the reporting contract.
- `agent-seam`: the settings requirement and the `Project settings are not loaded` scenario name the repository's *agent* instructions, so that the one file socrates reads on its own account is not read as forbidden.

## Impact

- **New** `internal/guidance`: the path, the read at the base commit, the ceiling, the refusal, and the three outcomes a run records.
- `internal/snapshot`: reading a file at a commit becomes an exported operation, so that git invocation stays in the one package that does it.
- `internal/prompt`: one more section in the fixed order, with socrates' own framing around the guidance and the existing neutraliser applied to it.
- `internal/reviewer`: reads the guidance once per round from the run's original base and passes it to the prompt.
- `internal/runstore`, `internal/presentation`: the outcome recorded and the line a person reads for it.
- `internal/e2e`: a repository whose base carries a guide, one whose change adds a guide the base does not have, and one whose guide is over the ceiling.
- **No change to** `internal/findings`, `internal/dialogue`, `internal/hygiene`, `internal/agent`, the findings schema, the classification, the verdict rule, the round cap, the diff ceiling, or any command's arguments.
- **No new configuration.** The path is fixed, the ceiling is fixed, and `~/.socrates/config.yaml` gains nothing.
- **Nothing is written in the repository under review**, and `.socrates/` there is a directory socrates only ever reads.

## Dependencies

None.

The most recent decision this stands on is `2026-08-22-no-customisations-reach-the-reviewer`, which decided what the reviewer may be told.
This change does not re-open it: nothing the vendor loads on its own account is re-enabled, and the one channel added is socrates' own, anchored, bounded and additive.

`2026-08-27-the-rules-the-reviewer-was-not-given` carried three rules from the same audit that produced this one, and named this fourth mechanism as the open question it could not take.
This is that question answered.
