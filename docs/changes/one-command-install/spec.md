# One command installs socrates and every skill

Decisions: ADR 61.

## Problem Statement

Installing socrates from a release takes four commands to place the binary and three more to install the skills, and the four do not work.

The download pattern names the operating system and not the architecture, so it fetches both the `amd64` and the `arm64` archive.
`tar -xzf socrates_*.tar.gz` then reads the second archive as the name of a member to extract, and exits 2.
On an `amd64` machine the right binary lands only because `amd64` sorts first.
On an Apple Silicon Mac or an `arm64` Linux machine the `amd64` binary lands, beside an error.

The commands also need `gh` signed in, though the repository is public.
They assume `~/.local/bin` exists, and they assume `sha256sum`, which a Mac may not have.

After the binary, the person has to know to run `socrates review skill`, `socrates design skill` and `socrates implement skill`.
`make install` runs all three, but nothing a person installing a release runs does, and `go install` installs none.

## Solution

Each release carries `install.sh`, and installing socrates is one line:

```sh
curl -fsSL https://github.com/socrates-tools/socrates/releases/latest/download/install.sh | sh
```

The script:

- reads the operating system and the architecture of the machine it runs on, mapping `x86_64` to `amd64` and `aarch64` to `arm64`.
- downloads `checksums.txt` from the release and finds the archive for that pair in it, so the release, rather than the script, says which pairs exist.
- refuses a pair the release carries no archive for, naming the pair and the pairs it does carry, and installs nothing.
- downloads that one archive with `curl`, or `wget` where `curl` is absent.
- checks the archive against `checksums.txt` with `sha256sum`, or `shasum -a 256` where that is absent, and refuses a mismatch, installing nothing.
- where `cosign` is on `PATH`, verifies the signature over `checksums.txt` before trusting it, and refuses a signature that does not verify.
- places `socrates` in `~/.local/bin`, or in the directory `SOCRATES_INSTALL_DIR` names, creating it where it does not exist.
- runs `socrates skill`.
- says where the binary went, and, where that directory is not on `PATH`, says so and how to add it.

`SOCRATES_VERSION` names a release tag to install instead of the latest.

`socrates skill` installs every skill every namespace offers, at each one's default, by running each namespace's own `skill`.
`make install` runs it once in place of the three namespace operations.

The README's manual route names the architecture, so it downloads one archive.

## User Stories

1. As a person installing socrates for the first time, I want one command that leaves the binary on my machine and every skill where my coding session reads it, so that I do not have to read the README to find the other six.
2. As a person on an Apple Silicon Mac, I want the installer to pick the `arm64` archive, so that the binary I get runs on my machine.
3. As a person on an `arm64` Linux machine, I want the same, for the same reason.
4. As a person on a platform the release carries no archive for, I want the installer to refuse, naming my platform and the ones it does carry, so that I know why nothing was installed rather than being handed a binary that cannot run.
5. As a person without `gh`, I want to install socrates with `curl` or `wget`, so that I do not need a GitHub account to install a public tool.
6. As a person on a Mac without `sha256sum`, I want the installer to check the archive with `shasum`, so that the check happens on my machine too.
7. As a person whose download was corrupted or tampered with, I want the installer to refuse an archive whose checksum does not match, and to leave nothing installed, so that I never run a binary nobody built.
8. As a person with `cosign` installed, I want the installer to verify the release's signature before it trusts the checksums, so that I get the stronger check without running it myself.
9. As a person with no `~/.local/bin`, I want the installer to create it, so that the install does not fail on a directory.
10. As a person whose `~/.local/bin` is not on `PATH`, I want the installer to tell me, and how to add it, so that `socrates: command not found` is not the next thing I see.
11. As a person who keeps binaries elsewhere, I want to name the directory with `SOCRATES_INSTALL_DIR`, so that socrates goes where my other tools are.
12. As a person who needs a particular release, I want to name it with `SOCRATES_VERSION`, so that I can install the version my team agreed on.
13. As a person upgrading, I want running the installer again to replace the binary and the skills with the new release's, so that upgrading is the same one command.
14. As a person who already has the binary, I want `socrates skill` to install every skill at once, so that I do not have to know which namespaces ship one.
15. As a person who wants one tool's skills somewhere else, I want `socrates <namespace> skill --to` to keep working as before, so that nothing I scripted breaks.
16. As a person installing by hand from the README, I want the commands there to download exactly one archive and to work, so that the manual route is not a trap.
17. As a contributor, I want the suite to fail when a release would not carry the installer, or when the README's one-liner names an asset the release does not carry, so that the two cannot drift.
18. As a contributor adding a tool with a skill, I want `socrates skill` to install it without the dispatcher changing, so that a new tool's skill is never left out.

## Implementation Decisions

**`internal/dispatcher`**

- A `skill` command at the top level, beside `version`.
  It finds the operation named `skill` in each namespace the root mounts, runs each with no flags set, and stops at the first that fails, returning that failure.
  It names no namespace, no flag and no skill: what each one prints is what that namespace's own `skill` prints.
- Its short help and its refusal of arguments live in `messages.go`.

**`install.sh`**

- At the repository root, POSIX `sh`, so it runs under `curl ... | sh` on Linux and macOS.
- Where it downloads from is `https://github.com/socrates-tools/socrates/releases/latest/download`, or `.../releases/download/$SOCRATES_VERSION` where that is set.
  `SOCRATES_RELEASES_URL` replaces `https://github.com/socrates-tools/socrates/releases`, which is the seam the suite downloads from a local server through, and which a mirror can use.
- It works in a temporary directory it removes on exit, and places the binary only after every check has passed.
- It carries no comments, like the rest of the source.
- It names no operating system or architecture as supported: which pairs exist is read from `checksums.txt`, so the platform statement in `internal/platform` stays the one copy.

**`.goreleaser.yaml`**

- `release.extra_files` attaches `install.sh` to every release.
- `checksum.extra_files` puts it in `checksums.txt`, so the signature covers it.

**`Makefile`**

- `install` runs `socrates skill` once in place of the three namespace operations.

**`README.md`**

- Installing leads with the one-liner, then names `SOCRATES_INSTALL_DIR` and `SOCRATES_VERSION`.
- The manual route downloads one archive, named by the operating system and the architecture, and says to use `shasum -a 256` where `sha256sum` is absent.
- The skills paragraph names `socrates skill` first, then the namespace operations for a person who wants one tool's skills elsewhere.

## Testing Decisions

- `internal/e2e` runs `socrates skill` against the built binary under the harness's home, and asserts that every skill every namespace's `skill` installs, and the implementer's agent definition, are where those operations put them, with what their sources generate.
  It follows `TestEverySkillAnyNamespaceInstallsIsInvokedByAPersonOnly` in walking `convention.Tools` rather than a list, so a new tool's skill is checked without the test changing.
- `internal/e2e` runs `install.sh` with `sh`, a real `curl`, a real `tar` and a real checksum tool, against a release served by `httptest` through `SOCRATES_RELEASES_URL`.
  The release holds an archive per pair: the one for the pair under test carries the built binary, every other carries a decoy that fails when run.
  A `uname` placed first on `PATH` states the pair, the way each vendor's simulator stands in for its CLI.
  The tests cover: `x86_64` Linux, `aarch64` Linux and `arm64` Darwin each installing the right archive and every skill; a pair with no archive refused by name with nothing installed; a mismatched checksum refused with nothing installed; a missing install directory created; a directory not on `PATH` reported; a second run replacing the first.
  The suite fails rather than skips where neither `curl` nor `wget` is on `PATH`.
- `internal/convention` asserts that `.goreleaser.yaml` attaches `install.sh` to the release and lists it in the checksums, and that the asset the README's one-liner downloads is the one the release attaches.
  Each check is mutated once to see it fail.
- `make help` and the documents check keep holding: `make install` still installs one binary, and the README still explains it.

## Out of Scope

- A Homebrew tap.
  It needs a tap repository in the organisation and a token that can push to it, which is a decision about the organisation rather than this repository, and the one-liner already serves macOS.
- Installing on Windows outside WSL.
  The release carries no Windows archive, and the installer refuses that pair by name like any other.
- `go install` installing the skills.
  Go installs a binary and runs nothing, so a person using it runs `socrates skill` afterwards, which the README says.
- Uninstalling.
