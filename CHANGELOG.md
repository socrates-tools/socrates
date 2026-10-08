# Changelog

What each release means to somebody who has socrates installed.
Written by hand: the commit messages in this repository say what a change is, and this says what it is for.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and socrates is versioned by [semantic versioning](https://semver.org/spec/v2.0.0.html).
While the major version is `0`, a minor release may change the surface half in ways a script has to follow.

## [Unreleased]

### Changed

- **The suite publishes one binary, `socrates`, and every workflow is reached through it.**
  v0.1.0 put `socrates`, `scrat` and `agora` on your `PATH`.
  This puts one name there, and `socrates review` and `socrates agora` are namespaces inside it rather than programs of their own.
  Remove `scrat` and `agora` yourself; nothing removes them for you.
- **A review begins with `socrates review begin`.**
  It was `socrates review`, which named the workflow twice once the namespace named it.
- **Every other review command moved under `review`.**
  `socrates status` is `socrates review status`, and the same for `findings`, `respond`, `continue`, `watch`, `logs`, `cancel`, `config`, `configure` and `skill`.
- **The page's commands moved under `agora`.**
  `agora serve` is `socrates agora serve`, and the same for `poll`, `end` and `stop`.
- **`socrates version` answers for the whole suite**, so a version in a bug report identifies one set of behaviour.
  `agora version` is gone, because there is no second program to ask.
- **Nothing is carried forward.**
  There is no shim, no alias and no deprecation period: the old names are not commands, and typing one is refused as an unknown command.
- **What socrates reads and writes moved into a directory per tool.**
  `socrates review` reads and writes `~/.socrates/review`, and `socrates agora` writes its server record and its page sessions under `~/.socrates/agora`, beside its own configuration file.
  `~/.socrates` is now the brand's root, holding one directory per tool, and `SOCRATES_HOME` moves the root rather than one tool's home.
  Nothing migrates the old layout: remove `~/.socrates` and run `socrates review skill` again.
- **The page's settings moved to a file of their own, `~/.socrates/agora/config.yaml`.**
  v0.1.0 read the `surface` block from the same file as every review setting.
  `~/.socrates/review/config.yaml` now holds review's settings only, and a `surface` block left in it is refused, naming the file it belongs in.
  Move the block there as it stands; nothing moves it for you.
  `socrates review config` and `socrates review configure` no longer list `surface.bind` and `surface.reached_at`.
- **The skill is installed at `.claude/skills/socrates-review/`.**
  It keeps the brand-qualified name because that directory is a flat namespace with no dispatcher to route a bare word.
  Remove `.claude/skills/socrates/` yourself; nothing removes it for you.
- **One exit code means one thing whichever namespace produced it**, so a script can branch on it without knowing which ran.
  `0` is nothing is asked of you, `1` is something is asked of you, `2` is the command failed, `3` is the command was not usable, and `4` is the work is still running.
  A clean review and an answered page are both `0`, and an escalated review and an unanswered page are both `1`.
  No code changed value, so a script that branched on them branches the same way.
- **One archive carries one binary**, with `LICENSE` and `NOTICE` beside it.
  Extracting it gives you what installing from source gives you, as before.
- `.socrates/review-guide.md` in a repository under review is unchanged, and a guide you already committed keeps working.
- **The page is a form, answered with one press at the end of it.**
  Choosing an option or typing words is answering: there is no longer a per-decision *Commit this answer* to press as well.
  The bar after the last card carries *Send and wait for a reply* and *Send and finish*, and Enter reaches the first.
- **The option a document recommends arrives chosen**, so agreeing with every recommendation costs one press.
  A single-choice decision also offers *None of these*, which clears the choice and moves you to the words field.
- **Sending over a decision you left blank asks first, in a dialog** that names each blank decision and opens on *Send anyway*.
  *Go back and answer*, Escape and a press outside the dialog all leave the page as it was, and each name in the dialog takes you to that card.
- **`socrates agora poll` returns your last answer.**
  A second send before anything collected the first replaces it rather than queueing behind it, and an answer to a document that has since been replaced is never handed out as an answer to the new one.

### Added

- socrates now drives `codex` and `opencode` as well as `claude`.
  If you run one of those and have no Anthropic account, socrates works for you now, and the first thing `README.md` puts in front of you is no longer a coupling.
  Name the one you want in `reviewer_cli` in `~/.socrates/review/config.yaml`.
- Which CLIs socrates drives follows from a stated bar rather than from a list somebody maintains.
  A CLI qualifies where the reviewed repository's own instructions can be kept out of the reviewer's session, where a terminal event says whether the work happened and what it said about a failure, and where a model is selectable per invocation.
  `README.md` names the three that clear it, the three that do not, and which capability each of those lacks.
- A findings document now arrives one of two ways, and both reach the same check.
  Where the vendor constrains the output against socrates' schema it is used as it stands; where it does not, the reviewer is asked to end its message with the document and socrates extracts it.
  A document that does not parse or does not satisfy the schema fails the round either way, so `clean` still means what it meant.
- A run records which reviewer CLI answered and how its document arrived, and `socrates review status` prints both.
- `socrates review config` prints every setting, the value in force and whether you set it, and `socrates review configure` changes one from a menu without opening the file.
- `socrates design`, a second workflow: a coding session interviews you over the page, one round of questions at a time, and leaves a spec and its decision records in the repository.
  `socrates design skill` installs the `socrates-design` skill that drives it, and the namespace holds nothing else.
  It keeps no state of its own, has no round cap, and ends at any round you choose, on the page or with `socrates agora end`.
- A document's options may carry no command, for an answer that runs nothing.
  A decision offers a command on every option or on none, and a decision mixing the two is refused.

### Notes

- An absent `reviewer_cli` runs `claude`, so an existing `~/.socrates/review/config.yaml` needs no edit and nothing about an existing run changes.
- A run recorded before socrates named a vendor reads as `claude`, because that is what answered.
- A `model` you state under `roles` is passed to the CLI you named and nothing translates it, because each CLI names its models its own way.
  Leave it out to let that CLI choose.
- Naming a reviewer CLI socrates does not drive is refused before a snapshot is taken and before a reviewer is spawned, naming the CLI and the guarantee it does not offer.
- The binary that runs reviews now contains an HTTP server, because the page folded into it.
  Running a review still starts no server and reaches no network; `socrates agora serve` binds a loopback address on purpose.

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

[Unreleased]: https://github.com/socrates-tools/socrates/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/socrates-tools/socrates/releases/tag/v0.1.0
