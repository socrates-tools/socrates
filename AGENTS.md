# socrates

socrates is the one binary the suite publishes, and it answers to that one name.
Each namespace of it is a tool, and each tool has a directory, a glossary, a directory of decisions and an instruction file of its own, named by `CONTEXT-MAP.md`.

- `socrates review` puts an independent reviewer in front of unreviewed local changes. See `review/AGENTS.md`.
- `socrates agora` serves one document as one page on a loopback address, waits for a person, and prints what they answered. See `agora/AGENTS.md`.
- `socrates design` interviews a person over agora and leaves a spec and its decisions behind, and cuts that spec into tickets. See `design/AGENTS.md`.
- `socrates implement` builds a change's tickets through implementers. See `implement/AGENTS.md`.

The dispatcher that routes `socrates` to a namespace holds wiring and nothing else: `internal/dispatcher` mounts review's command tree, agora's command tree, design's command tree and implement's command tree, and is the only thing that imports all four.
Each tool reads one configuration file of its own, per ADR 42: review reads `~/.socrates/review/config.yaml`, agora reads `~/.socrates/agora/config.yaml` and implement reads `~/.socrates/implement/config.yaml`.
review writes under `~/.socrates/review`, agora under `~/.socrates/agora` and implement under `~/.socrates/implement`.
design reads and writes none of these: nothing about a design outlives the session that ran it beside the spec, the ADRs, the glossary entries and the tickets it wrote to the repository being designed for.
`~/.socrates` is the brand's root, holding one directory per tool published under it, and `SOCRATES_HOME` moves the root.

`CONTEXT-MAP.md` says where each tool's glossary, decisions and instructions are, and `CONTEXT.md` defines the words the suite shares.
`docs/adr/` holds the decisions about the suite and each tool's own `docs/adr/` holds its own, each with what it rejected and why, under one number space.
This file is the map: the suite's layout, the verification sequence, and the invariants that hold across tools.
A tool's own file holds its layout and its invariants, and is read when working inside that tool.
A change starts as a spec, committed to `docs/changes/<slug>/spec.md` with its tickets as files beside it, and `CONTRIBUTING.md` says how.
Its folder is removed in the change that lands it, so `docs/changes/` holds only work still under way.
An issue is only where somebody using socrates reports a bug or asks for something.

## Layout

What the suite is:

- The module root - one package, which embeds `LICENSE` and `NOTICE` for `socrates licences` and holds nothing else.
- `cmd/socrates` - the binary: builds `internal/dispatcher`'s tree, runs it, and turns an error into an exit code.
- `internal/dispatcher` - the root command: routes `review`, `agora`, `design` and `implement` to their own trees, owns the error-to-exit-code mapping and the one failure prefix, and holds the top-level `version`, the top-level `licences`, which prints the `LICENSE` and `NOTICE` the root package embeds, per ADR 63, and the top-level `skill`, which runs every mounted tree's own `skill` without naming one, per ADR 61. It is not a tool, and nothing above the trees it mounts names a vendor, a namespace's own flag, or anything any tree does.
- `review/`, `agora/`, `design/`, `implement/` - one directory per tool, each with its own `AGENTS.md`, `CONTEXT.md` and `docs/adr/`.

What holds them all:

- `internal/convention` - the repository's conventions as executable checks.
- `internal/e2e` - real binaries, a real repository, real git, real processes, each vendor's simulator on `PATH`, a real browser.
- `internal/testing` - harness, temporary repository, simulator installer, source helpers.
- `internal/tty` - whether a stream is a terminal at all, for every tool that draws something only a person can answer.
- `internal/skill` - the machinery no tool owns: the ceiling a skill is held to, installing one under a directory, and reading which invocations a document names so one the binary does not offer fails the suite.
- `internal/settingsmenu` - the menu a settings file is changed from, for any tool: the rows, the view, the write that leaves the rest of the file alone, the write offered one key at a time to a writer that is not the menu, and the reading of a line the loader could not read.

## Verification

Run this before reporting work done:

```sh
make fmt
make verify
```

`make help` lists every target.
`make verify` is the one entry point, and it runs `make lint`, `make test`, `make corpus-outcomes` and `make fixtures-diff`, in that order.
Continuous integration runs that same command on every supported platform, so passing locally and passing there are one claim rather than two.

`make lint` fails on anything `gofmt` would rewrite, then vets, then runs the analysers `.golangci.yaml` names at the version the `Makefile` pins, so the cheapest gate is the first to fail.
`make test` drives the surface page with a real browser, so a Chrome or Chromium must be on `PATH` or named by `AGORA_BROWSER`.
`make skill`, `make assets` and `make platforms` regenerate committed files from their sources, and `make test` fails while a generated file and its source disagree.
`make corpus-outcomes` states which outcomes only fixtures reach and never fails for one; rewrite its recorded distribution with `scripts/corpus-outcomes -write` once a move is deliberate.
`make fixtures-diff` compares the simulator against the recordings and reaches nothing.

`make release` publishes what a tag names, gated on the verification matrix, and `make release-snapshot` produces the same archive locally without publishing anything.

These three reach outside this machine, and nothing in continuous integration runs them:

```sh
make fixtures     # record what the real claude CLI answers, as shapes
make live-smoke   # one real review against the real binary
make forge-smoke  # post what one review forwarded to a real pull request, and delete it again
```

## Invariants

Each is held by a test that fails when it stops holding.
None is a matter of care, and the reasoning for each is in `docs/adr/`, or in the tool's own `docs/adr/` where a tool's invariant is in its own file.

**The skill**

- Every skill any namespace ships carries `disable-model-invocation: true`, stated once in `internal/skill` and rendered into every skill's front matter, and the suite fails over every committed `socrates-*` skill and over everything `socrates <namespace> skill` installs, for each namespace in `convention.Tools` whose tree holds `skill`, rather than over a list, per ADR 49.

**The settings menu**

- `internal/settingsmenu` states every label the menu draws in its own `messages.go`, each tool's configure package states its own, and the wording check reads the plain string literals of all of them rather than only their `fmt` call sites.
- `internal/settingsmenu` names no tool: what differs between one tool's menu and another's is carried by the `Definition` a tool hands it.
- A setting whose written key is not its displayed key carries the segments it is written at, so a `default_branch` row is written under `default_branches.<repository top level>` whatever dots the path holds, and a row whose written key is its displayed key is written there.
- Every user-facing string lives in its package's `messages.go`, names what the person can do, and carries no em dash.

**The binaries**

- `socrates` answers to one name and nothing is installed beside it, and `agora` is a namespace inside it rather than a program of its own.
- What a binary says about itself was recorded when it was linked and is never read where it is run.
- review reaches no network outside a build tag, and what is measured is `review/command` and the packages it reaches rather than this repository.

**Project identity**

- The module path is the project's repository under the organisation that publishes it, and every import of the project resolves under it.
- The build states the path once and is checked against `go.mod`, so what a binary stamps into itself cannot disagree with what the compiler resolved.
- No source file and no document names a code-hosting account other than the organisation's.

**Verification**

- One entry point runs every gate, a run invokes that entry point rather than listing the gates, and formatting fails before the suite does.
- The tool list the cross-tool rule reads agrees with the tree: a namespace the dispatcher mounts and a top-level directory holding a `command` package are each on it, and a tool on it is mounted and holds a `command` package, so none is invisible to the rule that no tool imports another.
- The matrix is read from the platform statement, so the suite runs on every platform the set names and on no runner it does not.
- Every action a workflow uses is named by a commit with its version beside it, and a check fails while any is named by a tag.
- Every workflow states a read-only token at the top and escalates only in the job that needs more, and a superseded run is cancelled rather than left to report.
- A run happens when somebody asks for one and when a release gates on one, and nothing starts on a push, with one exception per ADR 64: Scorecard also runs weekly on a schedule, and no other workflow starts without being asked.
- The module graph is checked for known vulnerabilities on the same footing, and dependency updates arrive as changes somebody reviews.
- The browser the page half drives is installed by the run at a stated version rather than found on the runner, and an absent browser fails rather than skips.
- Nothing a workflow runs reaches the real CLI or a forge, by target or by build tag.

**Platform support**

- `internal/platform` states the supported set once, on the basis that reaping by process group is available, the filesystem enforces the seal and the suite runs there rather than on the basis that socrates compiles.
- The refusal, what a release builds and what a person reads all read that set, and no second copy of it exists in the source.
- A review is refused where this build cannot reap by process group, before a snapshot is taken and before a reviewer is spawned, naming the missing guarantee and the same way forward the README names.
- A command that neither spawns nor reaps is not refused there, so a finished run stays readable.
- `platforms.json` is generated from the statement, and a platform on one side and not the other fails the suite in both directions.

**Releases**

- A binary reports the release it was linked from beside the commit stamp, the version is set by the linker from the tag rather than written in the source, and a build that is not at a tag says it is a development build instead of borrowing the nearest one.
- What a release builds is the supported set and no others, and a target on one side and not the other fails the suite.
- One archive carries the one binary the suite publishes, `LICENSE` and `NOTICE`, and the binary and the archive come from the same build so they cannot be at different commits.
- Every job that publishes needs the job that ran the verification matrix, so a tag on a commit whose verification did not pass produces nothing.
- Every job that publishes also needs the code scan, and the scan starts only when somebody asks for one or a release gates on one.
- A release carries a checksum, a keyless signature over the checksums, a bill of materials and build provenance, and nothing in the release path reaches the real CLI.
- Each release is recorded in `CHANGELOG.md`, written by hand, in terms of what it means to somebody who has the previous one.

**Documentation**

- `README.md` states what socrates review needs before it explains how socrates review works, and the reviewer CLI, its floor and the account behind it are part of that.
- The toolchain floor a document states is read against `go.mod`, so the two cannot drift.
- Every binary `make install` puts on a path is explained rather than mentioned.
- Where what socrates review reads goes is stated in both directions: socrates review sends nothing, and the reviewer CLI it drives sends the change to its vendor. Neither half is stated without the other.
- `SECURITY.md` names a private channel and says what socrates review does not defend against, including the residual `README.md` already names.
- `CONTRIBUTING.md` says a change starts as a spec before it says anything about the toolchain, and it names one verification command rather than the gates.
- `make help` names every target a document tells a contributor to run.
- The documents a first change needs are present, and `.editorconfig` and `.gitattributes` match the conventions the source is held to.

**Licensing**

- `LICENSE` states MIT in full and names the holder, and no source file opens with a licence header.
- A comment carrying a borrowed work's notice and citing `NOTICE` is not that header, and one that states terms and cites nothing still is.
- `NOTICE` names everything a shipped binary carries that socrates did not author, with each licence's full text rather than a link: the client bundle's libraries, the font faces embedded beside it, and the modules the one binary links in.
- Every file `agora/assets` embeds is accounted for in `NOTICE`, read from what is embedded rather than from `web/`.
- The bundle that ships carries the copyright notice of every library bundled into it, asserted against the shipped bytes, and a bundler set to delete legal comments fails the suite.

**The suite**

- A test may not prove behaviour by reading implementation source; one that legitimately reads a Go file is named in `convention.ContractsTestsRead`.
- A gate that cannot check what it was asked to check fails rather than reporting a pass.
- A verification task may require that a live run happened and was recorded, never a particular severity, action or count.
- The menu half is driven by the real binary in a real pseudo-terminal, typed at and read back from what was drawn; exercising the model in process would pass on a build where `socrates review configure` cannot start.
- Nothing in continuous integration reaches the real CLI or a forge.
- The source carries no commentary.

## Conventions

- **No comments in source.** Names carry the meaning.
- **One sentence per line in Markdown.**
- **No em dashes anywhere.**
- **Commit messages say what the change is**, with no rationale essay and no AI attribution trailer.
- **The maintainer's commits carry the GitHub noreply address.**
  Author and committer are `Rui Cheng <2195250+ktei@users.noreply.github.com>`, for the maintainer and for any agent committing on their behalf.
  Run `git config user.email` before the first commit of a session, set the address in the repository if it differs, and never commit under an employer's or any other real address.
- **Tests name their scenario.** Prefer a test name containing the scenario text; use a `// scenario:` marker where the honest name would not.
- **Test first.** A task is not done until a test named for the scenario failed for the right reason and then passed.
- **A task about a live run says record it, never return this.**
- **Mutate a guard you rely on.** A check nobody has broken on purpose is a check nobody knows works.
