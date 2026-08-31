## Why

The module is published under the author's personal account, and the project is moving to an organisation, `github.com/socrates-devtool`.
For a Go project the module path is not metadata: it is the import path in 448 places across 181 files, it is the address `go install` resolves, and it is what `Makefile` stamps into the build through `MODULE`.
A project published under one account and importable only under another is a project whose first instruction to a stranger is wrong.

`openspec/config.yaml` is the second half of the same problem, and the worse half.
It is normative, it says the binary is `scrat` and the command directory is `cmd/scrat`, and both were renamed long ago.
It also names the design the project implements as a claude.ai artifact, by identifier, which nobody outside this machine can open.
A normative document that cites a source its readers cannot reach is a document that cannot be checked, and it is about to acquire readers who are not the author.

Doing the rename late is worse than doing it now.
The module path is the one thing about a Go project that other people's files start containing, and every day the repository is public is a day somebody could have written it down.

## What Changes

- The module path SHALL become `github.com/socrates-devtool/socrates`, and every import of it SHALL move with it, in one change that carries nothing else.
- `MODULE` in the `Makefile` SHALL move with it, so what `socrates version` reports keeps coming from the linked build rather than from a path written twice.
- `openspec/config.yaml` SHALL say that the binary is `socrates` and that `scrat` is a second name for it, that the command directory is `cmd/socrates`, and that `openspec/` is the sole normative source, replacing the reference to a private artifact.
- Nothing SHALL name an account other than the organisation, held by a check in `internal/convention` rather than by care.
- The archived changes SHALL be left as they are. An archived change is the record of what was decided when it was decided, and rewriting one makes the past claim something it did not say. The two that name the old path name it as it was, correctly.
- The commit history SHALL be rewritten to a personal address before the first public push, because authorship is the one thing that cannot be corrected once clones exist.

## Capabilities

### New Capabilities

- `project-identity`: the name and the account the project answers to, and where that name is allowed to be written.

### Modified Capabilities

None.

## Impact

- `go.mod`: the module path.
- 181 Go files: 448 import references.
- `Makefile`: `MODULE`.
- `openspec/config.yaml`: the stale binary and directory names, and the private artifact reference.
- `README.md`: the module path wherever an install or an import names it.
- `internal/convention`: a check that no source file names an account other than the organisation.
- The git remote, and the repository under the organisation.
- `openspec/changes/archive/`: untouched, deliberately.
- No behaviour changes. `make test` proves that by passing unchanged, and a diff that touches anything but a path is this change going wrong.
