**The suite must be green after every task, not only at the end of a group.**
The ordering below exists for that reason: the version is stamped and proved before anything publishes it, and the archive's contents are asserted before a workflow builds one.

**This change follows `the-platforms-the-seal-holds-on`.**
The release matrix is the supported set. A task here that names a platform the supported set does not is this change reaching past its dependency.

**A task about a release says record what it produced.**
No task requires a particular archive size, a particular digest, or a particular signature value. What is verified is that each artifact was produced and that the checks over it pass.

## 1. A build knows which release it is

- [x] 1.1 Add a `version` variable to `internal/version` set at link time beside `commit` and `uncommitted`, defaulting to the development-build statement, and verify `make lint test` is green
- [x] 1.2 Add the version to the ldflags in `Makefile` for both `STAMP` and `SURFACE_STAMP`, taken from the tag, and verify a build from an untagged tree still passes the suite
- [x] 1.3 Add an end-to-end test that a binary built from a tagged commit reports that version and the commit stamp it reported before. Covers `A tagged build reports its version`
- [x] 1.4 Add an end-to-end test that a binary built off a tag says it is a development build and names no version. Covers `An untagged build does not claim a version`
- [x] 1.5 Extend the existing boundary test so the version is held to the same rule as the commit stamp: recorded when linked, read nowhere it runs. Covers `The version is read where the binary was linked, not where it runs`

## 2. What an archive has to contain

- [x] 2.1 Add a test asserting an archive's contents: the tool under both names, the second binary, `LICENSE` and `NOTICE`. Covers `An archive carries the terms and the notices`
- [x] 2.2 Extend it to assert that both names report the same version and the same commit, so they cannot be at different builds. Covers `An archive carries both names and the second binary`
- [x] 2.3 Add a test that what a release builds and the stated supported set agree, and verify it fails when a target outside the set is added. Covers `The release matrix is the supported set`

## 3. The release runs

- [x] 3.1 Add `.goreleaser.yaml` building the supported set and no others, producing archives that satisfy the tests from group 2, and verify those tests pass against a locally produced snapshot build
- [x] 3.2 Add `.github/workflows/release.yml`, on tag, that runs the verification matrix first and publishes nothing if it does not pass. Covers `A release that did not pass verification produces nothing`
- [x] 3.3 Add checksums, a keyless signature over them, a bill of materials and build provenance to what the release publishes. Covers `A release carries what is needed to verify it`

  Four of the four are in the release. GitHub stores the provenance attestation for a public repository only, so the step signs and is then refused, and it is left in place and left failing the way the code-scanning upload is. Recorded under "Continuous integration" in `docs/gotchas.md`.
- [x] 3.4 Verify that nothing the release compiles or runs reaches the real CLI, by the same build-tag boundary continuous integration relies on, and record what was checked. Covers `Nothing in the release path reaches the real reviewer CLI`
- [x] 3.5 Mutate the guard: point the release at a commit whose verification fails, and verify no artifact is published, then restore

  `v0.0.0-guard` was pushed on a commit that broke the formatting gate and is on no branch. Both verification jobs failed, `publication` was reported skipped, and the repository held no release, no archive and no attestation afterwards. The tag was deleted and the commit is unreachable.

## 4. What a person reads

- [x] 4.1 Add `CHANGELOG.md` with an entry for the first release, saying what it means to somebody installing it. Covers `A release is recorded in the changelog`
- [x] 4.2 Add an install section to `README.md`, naming the release archive and `go install`, above the from-source instructions that exist today, and verify the em dash check and the user-facing string check still pass

## 5. The first release

- [x] 5.1 Tag `v0.1.0`, run the release, and record which artifacts it produced and that the verification over them passes

  `v0.1.0` on `9b9a1fb`. The verification matrix passed on Linux and on macOS, and the release then published ten artifacts: four archives, one bill of materials beside each, `checksums.txt`, and `checksums.txt.bundle`. The provenance step signed and was refused, as recorded in `docs/gotchas.md`, because GitHub stores an attestation for a public repository only.

  It took three attempts to get there, and neither of the first two published anything. The first was withdrawn so that `CHANGELOG.md` could stop claiming provenance the release does not carry. The second could not start its publication job at all: the account had reached its Actions spending limit, which a private repository reaches quickly because a macOS minute is billed as ten.
- [x] 5.2 Verify a downloaded archive extracts to a working `socrates`, `scrat` and `agora` on a supported platform, and record it

  `socrates_0.1.0_linux_amd64.tar.gz` downloaded on Linux, matched against `checksums.txt`, and its checksums verified with `cosign verify-blob --bundle` against the release workflow's identity and the Actions issuer: `Verified OK`. It extracted to `socrates`, `scrat` and `agora` beside `LICENSE` and `NOTICE`, all three reporting `v0.1.0, built from 9b9a1fb`, and the tool answered `--help`.
