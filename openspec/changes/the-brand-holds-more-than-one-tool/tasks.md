**The suite must be green after every task, not only at the end of a group.**
The ordering below exists for that reason.
A rename is one substitution applied at once or a tree that does not compile, so each group renames one thing completely and the suite answers before the next begins.

**Nothing here migrates anything.**
No task reads the old home, warns about the old name, or accepts it.
`design.md` records why, and a task that adds one is in the wrong change.

**The archive is not edited.**
`openspec/changes/archive` names the account and the binary that existed when each decision was taken.
A task that rewrites it is in the wrong change.

**The prose subject is not renamed.**
Where a specification names the tool as the subject of a sentence rather than as a command, it is left alone.
A task that substitutes it is doing the 2952-line rename `design.md` declines.

## 1. The identity

The account check matches a hosted account against the repository's name, so it reads nothing until the repository is named `socrates-review`.
The constants and the module path therefore move together, and the guard is mutated afterwards rather than relied on to fail first.

- [ ] 1.1 Change `Organisation` to `socrates-tools` and `Repository` to `socrates-review` in `internal/convention/account.go`, change the module path in `go.mod` to `github.com/socrates-tools/socrates-review`, rewrite every import in the module, and verify the account check passes and the tree builds
- [ ] 1.2 Change `MODULE` in the `Makefile` and the linker flags and attestation repository in `.goreleaser.yaml` to the new path, and verify the version stamp still composes from the one statement
- [ ] 1.3 Mutate the account guard by naming another account in a source file, verify it fails naming that file, and revert
- [ ] 1.4 Exempt the module path from the check that the surface names no word of the other tool, and verify the check still reports a word the surface itself names

  The tool's name carries `review`, so every import in the surface half tripped a check that reads the surface's vocabulary.
      An import path is where the code lives rather than a word the surface uses, so the path is removed from a line before the line is read.
      `internal/convention/ignorance.go` reads the path from `go.mod` rather than holding a second copy of it.

## 2. The home

- [ ] 2.1 Add a test that the brand's root holds the tool's configuration file and runs in a directory named for the tool, and that neither is directly in the root.
      Verify it fails against the flat layout. Covers `The tool writes under its own directory in the brand's root`
- [ ] 2.2 Resolve the tool's home as `review` beneath the root in `internal/state`, leaving `SOCRATES_HOME` naming the root, and verify 2.1 passes
- [ ] 2.3 Update the scenario that one variable moves everything to assert the tool's directory rather than the root, and verify `agora` still resolves its directory through the same package and reads no configuration file of its own. Covers `One variable moves everything`
- [ ] 2.4 Verify the refusal for a relative home still names the variable and the path, and that the seal message naming `SOCRATES_HOME` still names a variable that exists

## 3. The one name

- [ ] 3.1 Add a test that installing puts `socrates-review` on the path and no second name for that same binary beside it.
      Verify it fails while `install` links the short name. Covers `One name is installed and nothing is put beside it`
- [ ] 3.2 Rename `cmd/socrates` to `cmd/socrates-review`, update `BINARY`, `build` and `install` in the `Makefile`, remove `SHORT_NAME` and the symlink, and verify 3.1 passes
- [ ] 3.3 Rename every string in a `messages.go` that names the tool or a command of it, and verify the em dash check and the user-facing string check still pass. Covers `What it prints names the tool`
- [ ] 3.4 Remove `convention.ShortName` and the check that no emitted string carries it, and verify nothing else reads either
- [ ] 3.5 Rename the tests named for the two scenarios whose command changed. Covers `` `socrates-review version` takes no arguments `` and `` `socrates-review status` does not carry the wait ``
- [ ] 3.6 Update every txtar script and end-to-end helper that invokes the binary by name, and verify the suite is green
- [ ] 3.7 Update `openspec/config.yaml` to name the binary, the command directory and the project the repository now has, and verify the check that the context names a command directory that exists still passes
- [ ] 3.8 Add a test that the last segment of the module path, the command directory and the command the tool installs are one name, and that the segment before it is the organisation rather than a repeat of it. Covers `The repository is named for the tool`

## 4. The skill

- [ ] 4.1 Change `skill.Name` to `socrates-review`, regenerate the committed copy with `make skill`, and verify the generated file and its source agree
- [ ] 4.2 Rename the test that the skill names the full command to name the one name the tool has, asserting that every invocation names `socrates-review` and that no second name exists for one to use instead. Covers `The skill names the one name the tool has`
- [ ] 4.3 Verify the skill still names every command the binary offers, that the loop is driven end to end through the real binary, and that `.claude/skills/socrates` is replaced rather than left beside the new one

## 5. What a release carries

- [ ] 5.1 Remove the `scrat` build from `.goreleaser.yaml`, and update `project_name`, the archive template, the signing identity and the attestation repository to the new path
- [ ] 5.2 Add a test that an archive carries the tool under one name, the second binary, `LICENSE` and `NOTICE`, and no second name for either.
      Verify it against `make release-snapshot`. Covers `An archive carries the tool and the second binary`
- [ ] 5.3 Verify the release matrix still reads the platform statement and that nothing in the release path reaches the real CLI

## 6. What a person reads

- [ ] 6.1 Rename the tool, the commands and the paths in `README.md`, and verify the requirements block still comes first and that both directions of the reach are still stated
- [ ] 6.2 Rename them in `AGENTS.md`, correcting the layout, the invariants and the map to the repository that exists
- [ ] 6.3 Rename them in `CONTRIBUTING.md`, `SECURITY.md`, `SUPPORT.md`, `CODE_OF_CONDUCT.md`, `docs/`, `.github/` and `.github/CODEOWNERS`
- [ ] 6.4 Add an entry to `CHANGELOG.md` under `Unreleased` saying what the move and the rename mean to somebody who has `v0.1.0`, naming the command, the skill and the home, and saying to remove the old home rather than move it
- [ ] 6.5 Verify `make help` still names every target a document tells a contributor to run, and that no live document names the former account

## 7. The repository

- [ ] 7.1 Run `make fmt` and `make verify`, and verify every gate passes on the renamed tree
- [ ] 7.2 Transfer the repository to the `socrates-tools` organisation and rename it to `socrates-review`
- [ ] 7.3 Point the local remote at the new path, and verify `gh repo view` and a fetch both resolve
- [ ] 7.4 Verify the account check passes against the transferred repository, and that the archive is still exempt from it

  `the-door-a-stranger-arrives-at` task 3.3 enables private vulnerability reporting on this repository.
  It is done after 7.2 or it is done twice.
