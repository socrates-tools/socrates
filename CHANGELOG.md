# Changelog

What each release means to somebody who has socrates-review installed.
Written by hand: the commit messages in this repository say what a change is, and this says what it is for.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and socrates-review is versioned by [semantic versioning](https://semver.org/spec/v2.0.0.html).
While the major version is `0`, a minor release may change the surface half in ways a script has to follow.

## [Unreleased]

### Added

- socrates-review now drives `codex` and `opencode` as well as `claude`.
  If you run one of those and have no Anthropic account, socrates-review works for you now, and the first thing `README.md` puts in front of you is no longer a coupling.
  Name the one you want in `reviewer_cli` in `~/.socrates/review/config.yaml`.
- Which CLIs socrates-review drives follows from a stated bar rather than from a list somebody maintains.
  A CLI qualifies where the reviewed repository's own instructions can be kept out of the reviewer's session, where a terminal event says whether the work happened and what it said about a failure, and where a model is selectable per invocation.
  `README.md` names the three that clear it, the three that do not, and which capability each of those lacks.
- A findings document now arrives one of two ways, and both reach the same check.
  Where the vendor constrains the output against socrates-review's schema it is used as it stands; where it does not, the reviewer is asked to end its message with the document and socrates-review extracts it.
  A document that does not parse or does not satisfy the schema fails the round either way, so `clean` still means what it meant.
- A run records which reviewer CLI answered and how its document arrived, and `socrates-review status` prints both.

### Changed

- **The tool is now `socrates-review`, published under the `socrates-tools` organisation.**
  The command you type, the skill it installs and the repository it lives in all carry that name.
  `socrates` is the brand the organisation publishes under rather than this tool, so a second tool under it has a name to take.
- **`scrat` is gone.**
  There is no short name, and nothing is installed beside `socrates-review`.
  Write an alias if you want a shorter word, and spell it however you like.
- **What both binaries read and write moved to `~/.socrates/review`.**
  `~/.socrates` is now the brand's root, holding one directory per tool, and `SOCRATES_HOME` moves the root rather than this tool's home.
  Nothing migrates the old layout: remove `~/.socrates` and run `socrates-review skill` again.
- The skill is installed at `.claude/skills/socrates-review/`.
  Remove `.claude/skills/socrates/` yourself; nothing removes it for you.
- `.socrates/review-guide.md` in a repository under review is unchanged, and a guide you already committed keeps working.
- An absent `reviewer_cli` runs `claude`, so an existing `~/.socrates/review/config.yaml` needs no edit and nothing about an existing run changes.
- A run recorded before socrates-review named a vendor reads as `claude`, because that is what answered.
- A `model` you state under `roles` is passed to the CLI you named and nothing translates it, because each CLI names its models its own way.
  Leave it out to let that CLI choose.
- Naming a reviewer CLI socrates-review does not drive is refused before a snapshot is taken and before a reviewer is spawned, naming the CLI and the guarantee it does not offer.

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

[Unreleased]: https://github.com/socrates-tools/socrates-review/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/socrates-tools/socrates-review/releases/tag/v0.1.0
