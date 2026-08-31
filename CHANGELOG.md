# Changelog

What each release means to somebody who has socrates installed.
Written by hand: the commit messages in this repository say what a change is, and this says what it is for.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and socrates is versioned by [semantic versioning](https://semver.org/spec/v2.0.0.html).
While the major version is `0`, a minor release may change the surface half in ways a script has to follow.

## [Unreleased]

## [0.1.0] - 2026-09-01

The first release anybody other than its author can install.

### Added

- `socrates`, and `scrat` beside it, put an independent reviewer in front of the changes you just wrote.
  `socrates review` freezes the working tree into a commit nothing can move, spawns a detached reviewer against that snapshot, and returns.
  `status`, `findings`, `respond`, `continue`, `watch`, `logs` and `cancel` are how you drive it to a verdict.
- `socrates skill` writes the skill that teaches a coding session the whole loop, so the session polls and patches and you are asked only where a person is needed.
- `agora`, a second binary, serves one document as one page on a loopback address, waits for a person, and prints what they answered.
- Downloadable archives for Linux and macOS, on `amd64` and `arm64`.
  Each carries `socrates`, `scrat`, `agora`, `LICENSE` and `NOTICE`, so extracting one archive gives you what installing from source gives you.
- A checksum for every archive, a keyless signature over the checksums, and a bill of materials.
  Together they let you check a download against this repository without extending trust to whoever pushed the tag.
  The release records build provenance too, which GitHub stores for a public repository only, so v0.1.0 carries none.
- `socrates version` and `agora version` now name the release they were built from as well as the commit.
  A build that is not at a tag says it is a development build rather than borrowing the nearest one, so a version in a bug report is a version somebody can install.

### Notes

- socrates is supported on Linux and macOS, because the two promises it makes are held by unix mechanisms and because the suite runs there.
  On Windows, run it under WSL.
  A review started anywhere else refuses before it takes a snapshot.
- socrates needs Go 1.26 or later to build from source, `git`, and the `claude` CLI 2.1.169 or later on your `PATH`.
- socrates reviews.
  It does not commit, push, open a pull request, or run your project's tests, so `clean` means an independent reviewer found nothing blocking rather than that your suite is green.

[Unreleased]: https://github.com/ktei/socrates/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/ktei/socrates/releases/tag/v0.1.0
