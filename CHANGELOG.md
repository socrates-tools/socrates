# Changelog

What each release means to somebody who has socrates installed.
Written by hand: the commit messages in this repository say what a change is, and this says what it is for.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and socrates is versioned by [semantic versioning](https://semver.org/spec/v2.0.0.html).
While the major version is `0`, a minor release may change the surface half in ways a script has to follow.

## [Unreleased]

## [0.8.0] - 2026-10-16

0.7.0 was tagged but never published, so this is the first release that carries its changes too: read its entry below as well if you have 0.6.0.

### Added

- **An illustration can be anything HTML can draw.**
  A decision names an HTML file with `"illustration": {"path": "q3.html"}`, and the card draws it in a sandboxed frame.
  A relative path resolves against the document's folder, or against the working directory when the document arrives on standard input.
  The images, stylesheets and fonts it references are copied with it, so real screenshots, charts and interactive prototypes appear on the card itself.
  A drawing may load scripts, styles and fonts from cdnjs, jsDelivr and unpkg, and it can send nothing anywhere.
- **A drawing can look like the page.**
  Linking `/assets/tokens.css` gives it the page's own colours and fonts.
- **A frame grows to its drawing, up to the window's height, and can be expanded over the whole window.**
  Escape, the close button or a press on the backdrop closes it, and the drawing keeps its state.
- **socrates review's skill shows a UI finding as a before/after pair.**
  The agent captures the screenshots with whatever browser it has, and describes the change in words where it has none.

### Changed

- **The page now sends nothing beyond this machine, rather than reaching nowhere beyond it.**
  Opening a page that carries an illustration may fetch a library from one of those three hosts, and that host sees the request for it.

### Removed

- **A decision's `fragment` and the document's `stylesheet` are gone.**
  A document still naming either is refused, naming the field.
  Write the drawing as an HTML file and name it with `illustration.path` instead.

### Fixed

- **A link in the index no longer lands under the top bar.**
  Below 960px, a link to what was settled or to a decision, scrolled to the top of the window, could sit under the bar, so pressing it did nothing.

## [0.7.0] - 2026-10-16

### Added

- **`socrates agora poll --session <session> --last` prints the last answer again.**
  It prints what the last poll of that session printed, in the same form, without waiting and even once the surface has stopped.
  If an agent's output was cut short, by a pipe through `tail` or by its own process ending, the person no longer has to answer the page again.

### Changed

- **A person's answer is kept in its session before `poll` prints it.**
  It is written to `answer.json` in the session's directory under `~/.socrates/agora/sessions/`, and it goes when the surface collects the session after 45 minutes with nobody reading it.
  If it cannot be written, `poll` fails and the answer waits for the next poll rather than being handed out unkept.
- **The page fills the window.**
  From 960px wide, the navigation and the decisions scroll on their own beside each other, the navigation marks the decision in view, and the top bar carries the page's title.
  From 1280px a reply panel on the right lists what will go as the person answers, and keeps the message and both buttons in reach however far down they are.
  Below 960px the page is one column, as before.

## [0.6.0] - 2026-10-16

### Added

- **socrates has a landing page at [socrates-tools.github.io/socrates](https://socrates-tools.github.io/socrates/).**
  It is published with each release, and the README links to it.

### Changed

- **What was settled in earlier rounds arrives folded.**
  A later design round used to open with every earlier answer listed above the first question still open, and the list grew with each round.
  It now shows one line saying how many were settled, and pressing it, or the entry in the index, unfolds them.
  A page holding nothing but what was settled still shows it unfolded.

### Removed

- **The installer no longer downloads with `wget`.**
  It needs `curl`, which macOS ships and every Linux and WSL distribution socrates supports can install, and on a machine without it the installer stops and says so before writing anything.

## [0.5.0] - 2026-10-15

### Changed

- **Claude Code's skills and agent definition are installed under `$CLAUDE_CONFIG_DIR`.**
  Where you set it, an install and `socrates skill` now write there, which is where Claude Code reads them, rather than to `~/.claude`, which it does not.
  Where it is unset, nothing moves.
- **`--to` and `--agents-to` need `--host` beside them.**
  A directory holds one host's rendering, so `socrates review skill --to <dir>` is refused until you say whose, as `--host claude` or `--host codex`.
  If a script of yours passes `--to` alone, add `--host claude` to keep what it installed before.
- **The page draws *Recommended* in sea and *Your call* in Aegean blue.**
  Both were terracotta, so on one card they read as one signal.
- **On a narrow screen an option's meaning uses the whole width of the option.**
  It used to stop short of the *Recommended* mark and wrap into a narrow column beside empty space.

### Notes

- An install also writes Codex's copies of every skill under `~/.agents/skills` and the implementer as a Codex role under the Codex agents directory.
  Codex as a host is being built and is not yet supported, so nothing in Claude Code changes because of them.
- `socrates implement config` and `socrates implement configure` show the Codex implementer's model and effort beside Claude Code's, and they apply only once Codex is supported.

## [0.4.0] - 2026-10-14

### Added

- **`socrates licences` prints the terms the binary carries.**
  It prints `LICENSE` and then `NOTICE` from copies built into the binary, so a copy installed with `go install`, which has neither file beside it, carries them too.
  `socrates licenses` is the same command.
- **Each release carries its build provenance as `socrates.intoto.jsonl`.**
  It is the same provenance `gh attestation verify` checks, attached to the release for a verifier that reads it from there rather than from GitHub.

### Changed

- **The signature over the checksums is published as `checksums.txt.sigstore.json`.**
  It was `checksums.txt.bundle`, a name verifiers such as OpenSSF Scorecard do not look for, so earlier releases read as unsigned though they were signed.
  If you verify by hand, pass the new name to `cosign verify-blob --bundle`.
  The installer verifies against the new name and falls back to the old one, so `SOCRATES_VERSION` still installs an earlier release with its signature checked.

### Notes

- `cursor-agent`, `copilot` and `pi` stay unsupported as reviewer CLIs, and supporting them is not planned.
  `claude`, `codex` and `opencode` are unchanged.
- `SECURITY.md` now says where each supported reviewer CLI sends your change: `claude` to Anthropic, `codex` to OpenAI, and `opencode` to the provider you configured.

## [0.3.0] - 2026-10-14

### Added

- **One command installs socrates and every skill.**
  `curl -fsSL https://github.com/socrates-tools/socrates/releases/latest/download/install.sh | sh` downloads the archive for your operating system and architecture, checks it against the release's checksums, puts `socrates` in `~/.local/bin`, and installs every skill.
  It needs `curl` or `wget`, and no `gh` and no GitHub account.
  Where `cosign` is on your `PATH`, it also verifies the release's signature before trusting the checksums.
  It refuses and installs nothing where a check fails or the release carries no archive for your machine, and it tells you where `~/.local/bin` is not on your `PATH`.
  `SOCRATES_INSTALL_DIR` puts the binary somewhere else, `SOCRATES_VERSION` installs another release, and running it again upgrades the binary and the skills.
- **`socrates skill` installs every skill every namespace offers.**
  It runs `socrates review skill`, `socrates design skill` and `socrates implement skill` at their defaults, so after `go install` or an upgrade by hand you run one command rather than three.
  Each namespace's own `skill` still works, with `--to` and `--agents-to` as before, for putting one tool's skills somewhere else.

### Fixed

- **The README's manual install downloads one archive.**
  The commands v0.2.0 documented fetched both the `amd64` and the `arm64` archive, so `tar` exited with an error, and on an Apple Silicon Mac or an `arm64` Linux machine it placed the `amd64` binary.
  The manual route now picks the archive by architecture, uses `curl` rather than `gh`, and names `shasum` for a Mac without `sha256sum`.

## [0.2.0] - 2026-10-14

### Changed

- **The page moved to port 24672, and a held port is refused rather than walked past.**
  v0.1.0 bound port 18018 and, where that was taken, counted down through ten ports and took the first free one, so the page could end up where nothing you forwarded reached it.
  Now `socrates agora serve` binds `127.0.0.1:24672` and nothing else, and where something else holds it, it refuses and names `surface.bind` and `--address`, which still move the page.
  Whoever forwarded 18018 forwards 24672 instead.
- **From this release on, the next page after an upgrade is drawn by the build you installed.**
  v0.1.0 kept serving the page of whichever build started the surface until nobody had touched it for 45 minutes.
  Now `socrates agora serve` replaces a surface of another build once every session on it has been ended, and leaves it alone while a page on it is still being answered.
  A page that was already ended goes with the old surface, and a `serve` naming a session never replaces the surface holding it.
  This starts with the surfaces this release starts, because a v0.1.0 surface cannot say which build it is.
  One still running is left alone on 18018 while the page moves to 24672, and where you pinned its port with `surface.bind` or `--address`, `serve` refuses that port until the old surface exits on its own after 45 minutes idle or you stop its process.
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
- `socrates design skill` also installs `socrates-tickets`, which cuts a spec under `docs/changes/<slug>/` into tickets beside it, with `README.md` holding their order and status.
- `socrates implement`, a third workflow: once the tickets are cut and committed, the session you invoke the `socrates-implement` skill in oversees implementers that build each ticket in a worktree of its own, and you stay free to talk to it while they work.
  It lands each ticket by fast-forward on a branch of its own, runs the repository's own verification after every landing, and reads the whole change against the spec before telling you the branch is ready.
  `socrates implement skill` installs the skill and the `socrates-implementer` agent definition, and `socrates implement config`, `configure` and `default-branch` hold how many implementers run at once, which model they run on and which branch a feature starts from.

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

[Unreleased]: https://github.com/socrates-tools/socrates/compare/v0.8.0...HEAD
[0.8.0]: https://github.com/socrates-tools/socrates/releases/tag/v0.8.0
[0.7.0]: https://github.com/socrates-tools/socrates/releases/tag/v0.7.0
[0.6.0]: https://github.com/socrates-tools/socrates/releases/tag/v0.6.0
[0.5.0]: https://github.com/socrates-tools/socrates/releases/tag/v0.5.0
[0.4.0]: https://github.com/socrates-tools/socrates/releases/tag/v0.4.0
[0.3.0]: https://github.com/socrates-tools/socrates/releases/tag/v0.3.0
[0.2.0]: https://github.com/socrates-tools/socrates/releases/tag/v0.2.0
[0.1.0]: https://github.com/socrates-tools/socrates/releases/tag/v0.1.0
