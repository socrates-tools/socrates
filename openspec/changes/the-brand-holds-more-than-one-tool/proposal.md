## Why

The project moves to `socrates-tools`, an organisation that will publish more than one tool.
`socrates` currently names the organisation and this tool at once, so the second tool the organisation publishes has no name left to take, and a person reading `socrates` cannot tell whether they are being told about the brand or about the reviewer.

The tool takes the name of what it does within the brand.
The brand keeps the name it already answers to, and becomes the thing more than one tool sits under: the account that publishes them, and the directory on a person's machine that holds what they write.

Doing the move and the rename apart would publish `socrates-tools/socrates` first, which is the collision this change exists to remove.

## What Changes

- **BREAKING** The module path becomes `github.com/socrates-tools/socrates-review`, and the organisation the identity check holds becomes `socrates-tools`.
- **BREAKING** The binary becomes `socrates-review`.
  `scrat` is removed rather than renamed, and everything that existed because there were two names goes with it.
- **BREAKING** The skill becomes `socrates-review`, installed at `.claude/skills/socrates-review/`.
- **BREAKING** `SOCRATES_HOME` names the brand's root rather than this tool's, and this tool owns one directory beneath it.
  What was `~/.socrates/config.yaml` is `~/.socrates/review/config.yaml`, and what was `~/.socrates/runs` is `~/.socrates/review/runs`.
- Every string the tool prints names `socrates-review`, and no string it prints names a second name for it.
- A release archive carries one binary under one name beside `agora`, rather than one binary under two names beside `agora`.
- `agora` is unchanged: its name, its commands, and that it writes in a directory of its own under the tool's home.
- `.socrates/review-guide.md` in a repository under review is unchanged.
  `.socrates/` there is the brand's directory by the same rule as `~/.socrates`, and `review-guide.md` already names the tool within it.

## Capabilities

### New Capabilities

None.
The brand's root and the tool's directory beneath it are where the tool's home resolves to, which `cli-foundation` already owns.

### Modified Capabilities

- `project-identity`: the module path names the organisation and the repository the tool is named for, and the rule that renaming the module renames nothing about the tool is replaced by the rule that the two are one name stated in two places.
- `cli-foundation`: the binary has one name and no second one, what it prints names that name, and the home the one environment variable moves is the brand's rather than the tool's.
- `agent-skill`: the skill is named and installed for the tool under its new name, and the requirement that an invocation is written under the full name rather than the short one is removed with the short name.
- `distribution`: an archive carries the tool under one name, so the requirement that the two names come from one build has nothing left to hold.
- `review-dialogue`, `review-surface`, `reviewer-run`, `run-presentation`, `run-store`, `snapshot`: the invocations these requirements name are the commands a person types, and the command a person types has changed.

A requirement is restated here where it names a command.
Where it names the tool as the subject of a sentence it is left alone, and `design.md` records why.

## Impact

**Identity, stated once each:**

- `go.mod`, and every import in the module.
- `Makefile`'s `MODULE`, which the version stamp is composed from.
- `internal/convention/account.go`: `Organisation` and `Repository`.

**The binary and the short name:**

- `cmd/socrates` becomes `cmd/socrates-review`.
- `.goreleaser.yaml` loses the `scrat` build, and its `project_name`, archive template, signing identity and attestation repository follow the new path.
- `internal/convention/wording.go` loses `ShortName` and the check that no emitted string carries it.
- `internal/skill`: `Name`, and the tests that hold the skill to naming the full command rather than the short one.
- `Makefile`'s `SHORT_NAME` and the symlink `install` makes.

**The home:**

- `internal/state`: `homeDirectory` becomes the brand's root, and a directory named for the tool is resolved beneath it before runs, surface and configuration.
- `internal/surface` resolves its directory through the same package and is unaffected beyond what it resolves to.

**Documents:**

- `README.md`, `AGENTS.md`, `CONTRIBUTING.md`, `SECURITY.md`, `SUPPORT.md`, `CODE_OF_CONDUCT.md`, `docs/`, `.github/`, `CHANGELOG.md`.
- `openspec/config.yaml`, which names the binary, the command directory and the project.
- `openspec/specs`, in the requirement headings and prose that name `scrat` and `socrates`.

**Not touched:**

- `openspec/changes/archive`, which records decisions as they were taken and names the account and the binary that existed then.
- `cmd/agora` and `internal/surface`'s own naming.
- `.socrates/review-guide.md` in a repository under review.

**Depends on:** nothing.
`the-door-a-stranger-arrives-at` is underway and its remaining task enables private vulnerability reporting on the repository, which is a setting on the repository this change transfers.
That task is done after the transfer or it is done twice.
