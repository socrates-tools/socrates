# socrates-review

socrates-review puts an independent reviewer in front of unreviewed local changes.
The binary is `socrates-review`, and it answers to that one name.
Everything socrates-review prints names `socrates-review`.

`agora` is a second and separate program in this repository, not a third name.
It serves one document as one page on a loopback address, waits for a person, and prints what they answered.
It knows nothing about code review, and nothing socrates-review does reaches it yet.
One tool, so one home and one configuration file: both binaries read `~/.socrates/review/config.yaml` and write under `~/.socrates/review`.
`~/.socrates` is the brand's root, holding one directory per tool published under it, and `SOCRATES_HOME` moves the root.

`CONTEXT.md` defines the vocabulary and `docs/adr/` holds the decisions, each with what it rejected and why.
This file is the map: the layout, the verification sequence, and the invariants the design rests on.
A change starts as a spec published to this repository's issues, and `CONTRIBUTING.md` says how.

## Scope

socrates-review reviews.
It does not land, verify, render, serve, or wait for a person.
Committing, pushing, opening a pull request, posting a comment and watching continuous integration are not socrates-review's and never become socrates-review's.
It resolves a base, records what a person decided about a finding, and reports whether a finding's location is inside the change; fetching, rebasing and posting are the session's and the person's.
Running the reviewed project's own tests is not in the core layer either, so `clean` means an independent reviewer found nothing blocking, not that the tests pass.

## Layout

What socrates-review is:

- `cmd/socrates-review` - the binary: `review`, `status`, `findings`, `respond`, `continue`, `watch`, `logs`, `cancel`, `skill`, `version`, and the hidden `supervise`.
- `internal/snapshot` - the freeze: a temporary index, a commit, a detached worktree, and the seal that makes it unwritable.
- `internal/runstore` - run directories, `run.json`, run selection, rounds, retention.
- `internal/runstate` - the states and which are terminal.
- `internal/findings` - the findings contract, the reviewer's schema, the classification and the verdict rule.
- `internal/dialogue` - what a round leads to, which rounds wait on a person, which resolutions a finding accepts, the round cap, and the response document.
- `internal/hunks` - the lines a change touches, read from its own diff, and whether a finding's location is one of them.
- `internal/hygiene` - model output neutralised before it re-enters a prompt.
- `internal/guidance` - the one file socrates-review reads out of the repository under review, the commit it is read at, and the ceiling it is held to.
- `internal/prompt` - the delivered prompt, assembled from named sections.
- `internal/testquality` - the test-quality rule, shared by every role that needs it.
- `internal/platform` - the one place the supported platforms are stated, and the basis on which one qualifies.
- `internal/agent` - the only package that starts a model process, the reviewer CLIs socrates-review drives and the basis on which one qualifies, and one file per vendor behind the seam.
- `internal/reviewer` - beginning, supervising, responding, advancing, cancelling.
- `internal/report` - what `status` and `findings` print.
- `internal/presentation` - the one place an event becomes a line a person reads.
- `internal/follow` - following a run's event log live.
- `internal/terminal` - the pager, and the plain fallback.
- `internal/skill` - the two skills a coding session reads, one for its own change and one for somebody else's, and installing them.
- `internal/exit`, `internal/version`, `internal/config`, `internal/runlog`, `internal/state`, `internal/refusal`, `internal/process` - the small owned things: exit codes, the linked build stamp, the one optional configuration file both binaries read, the JSONL log, the one home both binaries write under, refusals, process groups.
- `internal/corpus` - what real reviewers answered, verbatim and embedded, and what the routing rule does with it.
- `internal/fidelity` - the recorded shapes of each real CLI, what was measured about each one's criteria, the drift check and the live smoke check, all network behind a build tag.

What agora is:

- `cmd/agora` - the binary: `serve`, `poll`, `end`, `stop`, `version`, and nothing else.
- `internal/document` - what a caller writes for the surface to render, what it copies rather than accepts, and the languages the page can highlight.
- `internal/sanitise` - a caller's HTML fragment reduced to a fixed allowlist, and every name it carries prefixed.
- `internal/stylesheet` - a caller's stylesheet read as a parse, reduced to a fixed allowlist, and scoped to the illustration.
- `internal/render` - a document to one HTML page, and a patch to the rows of a split diff.
- `internal/assets` - the built client bundle and fonts, embedded, and the stamp of their source.
- `internal/surface` - the loopback server and its guards, the sessions it holds, what it writes about itself and about each page, the exchange, and the bounds on waiting.
- `web/` - the client source, the design system and the fonts, built by `make assets`.

What holds them both:

- `internal/convention` - the repository's conventions as executable checks.
- `internal/e2e` - real binaries, a real repository, real git, real processes, each vendor's simulator on `PATH`, a real browser.
- `internal/simulation` - what every simulator shares: the scenario vocabulary, the invocation log, and the run that turns one into the other.
- `internal/claudesim`, `internal/codexsim`, `internal/opencodesim` - each simulator's half of its real CLI's wire contract: the arguments it reads, the stream it writes.
- `internal/testing` - harness, temporary repository, simulator installer, source helpers.
- `cmd/claudesim`, `cmd/codexsim`, `cmd/opencodesim`, `cmd/corpus-outcomes`, `cmd/corpus-record`, `cmd/platforms` - the simulator and the gate binaries.

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

`make release` publishes what a tag names, gated on the verification matrix, and `make release-snapshot` produces the same archives locally without publishing anything.

These two reach the real CLI, and nothing in continuous integration runs them:

```sh
make fixtures    # record what the real claude CLI answers, as shapes
make live-smoke  # one real review against the real binary
```

## Invariants

Each is held by a test that fails when it stops holding.
None is a matter of care, and the reasoning for each is in `docs/adr/`.

**The repository under review**

- Nothing writes in it: no run state, no configuration, no log.
- Every round is snapshotted afresh against the run's original base, and snapshot commits are content-addressed.
- A base a person names resolves to the fork point of that ref and `HEAD`, is recorded once when the run begins and is reused by every round. A ref that does not resolve, and a ref with no merge base against `HEAD`, are refused before a snapshot is taken and before a reviewer is spawned.
- Nothing fetches and nothing rebases: no ref is written, no commit is rewritten, and the reflog is untouched by resolving a base.
- A write in the snapshot worktree is refused by the filesystem, and the worktree is compared against its commit after the round.
- One file is read out of it, `.socrates/review-guide.md`, at the run's base commit and never from the working tree or the snapshot, so a guide the change writes governs the next run rather than this one.
- That guide is bounded, neutralised on the way into the prompt, and additive: it withdraws no rule socrates-review delivers, and nothing socrates-review decides is derived from a word of it.

**The reviewer**

- `internal/agent` is the only package that starts a model process.
- The reviewer CLIs socrates-review drives are stated once, each beside the facts that make it qualify, and support follows from those facts rather than from a list.
- A CLI qualifies on three things: the reviewed repository's own instructions can be kept out of the session, a terminal event says whether the work happened and what it said about a failure, and a model is selectable per invocation.
- Whether the vendor constrains output against the schema is a tier rather than a criterion, because socrates-review validates the document either way.
- A configured CLI socrates-review does not drive is refused by name, naming the guarantee it does not offer, before a snapshot is taken and before a reviewer is spawned.
- Nothing above the seam names a vendor, a flag or an event shape, and nothing outside socrates-review shapes a reviewer's answer into the findings document.
- Each vendor's invocation asks for none of the machine's own customisations; the session's tool list and named sources are recorded rather than trusted, and a vendor that reports no inventory produces runs that record that.
- A round after the first inherits no session: the only cross-round context is the labelled prior findings and recorded responses.
- Model output is neutralised per field on the way into a prompt, and never on the way to disk.
- What a person passes as the intent is neutralised on the way into the prompt whatever the run, and the section carrying it names the party the run records as having stated it.
- A review inside a review is refused.

**Verdicts and the dialogue**

- `clean` requires a findings document that parsed; a missing, dead or unparseable one is `failed`, and `failed` is never merged with `escalated`.
- A review of a working tree holding nothing its base commit does not is refused before a run is created, so no pass is recorded for a review that did not happen.
- A failure carries the account the agent gave of it, bounded, quoted for a person and never matched on.
- `findings.EffectiveAction` is the one reader of a finding's raw action.
- `dialogue.After` is the one rule for what a round leads to, and `dialogue.AwaitsAPerson` the one rule for which rounds wait on a person.
- A finding that asks a person parks the round rather than ending the run.
- A rejection and a decision are claims a fresh reviewer judges, never suppressions: the reason is required, and no prompt asks a reviewer to withhold a finding.
- A resolution that leaves no artifact behind carries a reason. A fix leaves a patch and a forward leaves the finding's own text, so a rejection, a decision and a drop each state why, and `dialogue.Accepts` is the one rule for which resolutions a finding takes.
- Where the run records that nobody here can patch the change, every finding of a round is a person's, and `dialogue.Classified` is the one place that is decided. A finding that holds the change while asking for nothing is forwarded rather than ending the run, so `dialogue.After`, `dialogue.AwaitsAPerson` and `dialogue.Answerable` gain no clause and keep their rules.
- Such a finding is forwarded, dropped or rejected whatever its own action, every round of such a run waits on a person, and `continue` ends the run `escalated` with the findings handed over rather than starting a round.
- A drop's reason is recorded on the run and printed to nobody.
- A recorded response is appended, never rewritten.
- The round cap is two, nothing configures it, and a configuration file naming one is refused.

**Presentation**

- Every round of a run that ended in `clean`, `awaiting-response` or `escalated` is handed over as a page, and a run that failed or was stopped is handed over as nothing whatever its rounds got as far as.
- socrates-review states whether a round waits on a person and does nothing further about the page: no command, flag or run record names one, and it states nothing about whether a page is served.
- `internal/presentation` is the one place an event becomes a line, and `internal/exit` owns the exit-code table.
- `status` prints the base commit the run resolved, as a bare fact and not a warning, and socrates-review says nothing about how far a named ref has moved past it.
- `findings --json` reports each finding's recorded resolution and whether its location is inside the change, and `internal/hunks` is the one place that question is answered.
- A follower stops on the run's record or on the log, waits bounded for a log that does not exist yet, and stops on a parked round.
- `status`, `watch` and `logs` stay three commands.
- Every user-facing string lives in its package's `messages.go`, names what the person can do, and carries no em dash.

**The skill**

- `internal/skill` is the one source of both skills, and each committed copy is generated from it.
- The loop for a change nobody here can patch is a second skill rather than a branch inside the first, and each skill carries a ceiling of its own that the suite reads from the skill rather than from a list.
- The second skill relays what somebody else wrote verbatim and acts on none of it, reads the change in a worktree of its own, records the head commit and refuses to post where it moved, takes a comment's placement from what socrates-review reported, and posts nothing the person did not choose.
- It names nothing the binary does not offer, lists every state a run rests in, teaches polling at a stated cadence, names `socrates-review watch` once and outside the loop, says once that the address `agora serve` prints holds the session `agora poll` is given and that the session is kept for the whole review, and says whose a question is.
- Every branch that reached a result serves a page, `failed` and `cancelled` name none, and the session waits only where socrates-review says the round waits on a person.
- `dialogue.ResponseRules` and `skill.PageBar` are each stated once and rendered where they are needed.
- It says every control carries the command it becomes and that the page shows that command to nobody, and it names each decision from the finding's identifier rather than titling it.
- It says every finding in the round gets a card of its own, that a finding the session resolved states what resolved it and offers no options, and that what foregrounds the rest is the severity order rather than shortening, listing or omitting anything.
- The loop it describes is driven end to end through the real binary, not read.

**The binaries**

- `socrates-review` answers to one name and nothing is installed beside it, and `agora` is a different program.
- What a binary says about itself was recorded when it was linked and is never read where it is run.
- socrates-review reaches no network outside a build tag, and what is measured is `cmd/socrates-review` and the packages it reaches rather than this repository.

**Project identity**

- The module path is the project's repository under the organisation that publishes it, and every import of the project resolves under it.
- The build states the path once and is checked against `go.mod`, so what a binary stamps into itself cannot disagree with what the compiler resolved.
- No source file and no document names a code-hosting account other than the organisation's.

**Verification**

- One entry point runs every gate, a run invokes that entry point rather than listing the gates, and formatting fails before the suite does.
- The matrix is read from the platform statement, so the suite runs on every platform the set names and on no runner it does not.
- Every action a workflow uses is named by a commit with its version beside it, and a check fails while any is named by a tag.
- Every workflow states a read-only token at the top and escalates only in the job that needs more, and a superseded run is cancelled rather than left to report.
- A run happens when somebody asks for one and when a release gates on one, and nothing starts on a push: what that gives up is that nobody learns the default branch is red, or that an advisory was published, until somebody asks.
- The module graph is checked for known vulnerabilities on the same footing, and dependency updates arrive as changes somebody reviews.
- The browser the page half drives is installed by the run at a stated version rather than found on the runner, and an absent browser fails rather than skips.
- Nothing a workflow runs reaches the real CLI, by target or by build tag.

**Platform support**

- `internal/platform` states the supported set once, on the basis that reaping by process group is available, the filesystem enforces the seal and the suite runs there rather than on the basis that socrates-review compiles.
- The refusal, what a release builds and what a person reads all read that set, and no second copy of it exists in the source.
- A review is refused where this build cannot reap by process group, before a snapshot is taken and before a reviewer is spawned, naming the missing guarantee and the same way forward the README names.
- A command that neither spawns nor reaps is not refused there, so a finished run stays readable.
- `platforms.json` is generated from the statement, and a platform on one side and not the other fails the suite in both directions.

**Releases**

- A binary reports the release it was linked from beside the commit stamp, the version is set by the linker from the tag rather than written in the source, and a build that is not at a tag says it is a development build instead of borrowing the nearest one.
- What a release builds is the supported set and no others, and a target on one side and not the other fails the suite.
- One archive carries the tool under both of its names, the second binary, `LICENSE` and `NOTICE`, and the two names come from one build so they cannot be at different commits.
- Every job that publishes needs the job that ran the verification matrix, so a tag on a commit whose verification did not pass produces nothing.
- A release carries a checksum, a keyless signature over the checksums, a bill of materials and build provenance, and nothing in the release path reaches the real CLI.
- Each release is recorded in `CHANGELOG.md`, written by hand, in terms of what it means to somebody who has the previous one.

**Documentation**

- `README.md` states what socrates-review needs before it explains how socrates-review works, and the reviewer CLI, its floor and the account behind it are part of that.
- The toolchain floor a document states is read against `go.mod`, so the two cannot drift.
- Every binary `make install` puts on a path is explained rather than mentioned.
- Where what socrates-review reads goes is stated in both directions: socrates-review sends nothing, and the reviewer CLI it drives sends the change to its vendor. Neither half is stated without the other.
- `SECURITY.md` names a private channel and says what socrates-review does not defend against, including the residual `README.md` already names.
- `CONTRIBUTING.md` says a change starts as a spec before it says anything about the toolchain, and it names one verification command rather than the gates.
- `make help` names every target a document tells a contributor to run.
- The documents a first change needs are present, and `.editorconfig` and `.gitattributes` match the conventions the source is held to.

**Licensing**

- `LICENSE` states MIT in full and names the holder, and no source file opens with a licence header.
- A comment carrying a borrowed work's notice and citing `NOTICE` is not that header, and one that states terms and cites nothing still is.
- `NOTICE` names everything a shipped binary carries that socrates-review did not author, with each licence's full text rather than a link: the client bundle's libraries, the font faces embedded beside it, and the modules the two binaries link in.
- Every file `internal/assets` embeds is accounted for in `NOTICE`, read from what is embedded rather than from `web/`.
- The bundle that ships carries the copyright notice of every library bundled into it, asserted against the shipped bytes, and a bundler set to delete legal comments fails the suite.

**The surface**

- It binds a loopback address; nothing widens it, and an address that is not one is refused before any server starts.
- It takes the first free port of a fixed window counting down where nothing states one, and a whole window held is a refusal naming it.
- A port an operator committed to, by flag or by either field of the `surface` block, is bound or refused and never walked past.
- The surface mints a session, `--session` never brings one into existence, and every page and every route of it is served under a path holding that session.
- It turns away a request naming a host it did not bind, and every state-changing route carries a token: a session's own answers for that session and no other, and the one minted when the server started answers everywhere.
- It composes a command and never runs one, the command it hands back is the option's own string, and no command is rendered on the page.
- It writes only under socrates-review's home, in a directory of its own beside the runs, and it resolves no home and reads no configuration file of its own.
- A caller writes a document rather than a page, and the name, severity, audience and location come from the source a decision cites.
- A document may hold no decision, and one holding no decision, no source entry and no preamble is refused.
- A decision may carry no options, and such a card states something and offers no control a person could mistake for a choice; one carrying no options and stating a mode is refused.
- Every entry of the source is rendered ahead of the decisions, and both are ordered by the severity the source carries, most severe first, informational last, an unstated severity after all of them, and ties in the order the document listed them.
- The order is the surface's, so no field of a document reaches it, and a page with nothing to answer offers neither control and says so.
- Every decision carries a free-text field that no document field can remove, including one that offers no options.
- The languages the page highlights are the bundled library's common set, stated once and generated from `web/`, and a document naming one outside it is refused while one naming none renders as plain text.
- A caller's fragment is reduced to a fixed allowlist before it is embedded, and every name it carries is prefixed by the surface.
- A caller's stylesheet is reduced from a parse that reached the end of its input, every selector it keeps is scoped to the illustration, and the page admits it by a hash of exactly the bytes it embedded.
- What a caller paints is contained by the illustration, so a rule inside one decides nothing about where the page's own controls land.
- Delivery happens once, waiting is bounded, collection is per session, and the server stops once it has held no session for the idle period.
- A stream is told where the exchange stands when it connects rather than only when something next changes, so a client that missed a revision catches up rather than waiting on a change that has already happened.
- The end of an exchange stops the page offering to answer and not the page listening, so a caller who writes again reaches the page the person is already looking at, and only the test about reloading reloads before it asserts.
- Every asset is embedded, no diagram engine is, and the committed bundle is generated from `web/`.
- An element it hides paints nothing, every colour it paints is a token stated once, and every text meets a contrast floor under both preferences.
- Nothing in it names review, finding, verdict, round, or run as a noun, and the module path is not part of its vocabulary.

**The suite**

- Every end-to-end scenario runs against every supported vendor's simulator, and a supported vendor with no simulator or no recordings fails the suite.
- Tests assert on the delivered prompt, never on socrates-review's own source text, and the prompt's reporting section is generated from `internal/findings`.
- A test may not prove behaviour by reading implementation source; one that legitimately reads a Go file is named in `convention.ContractsTestsRead`.
- A gate that cannot check what it was asked to check fails rather than reporting a pass.
- What a real reviewer answered is kept verbatim and replayed; the corpus is classified by `dialogue.After`, and a recording is only made from a run that answered.
- A verification task may require that a live run happened and was recorded, never a particular severity, action or count.
- The page half is driven by a real browser against the real binary; asserting on served markup would test the renderer.
- Nothing in continuous integration reaches the real CLI.
- The source carries no commentary.

## Conventions

- **No comments in source.** Names carry the meaning.
- **One sentence per line in Markdown.**
- **No em dashes anywhere.**
- **Commit messages say what the change is**, with no rationale essay and no AI attribution trailer.
- **Tests name their scenario.** Prefer a test name containing the scenario text; use a `// scenario:` marker where the honest name would not.
- **Test first.** A task is not done until a test named for the scenario failed for the right reason and then passed.
- **A task about a live run says record it, never return this.**
- **Mutate a guard you rely on.** A check nobody has broken on purpose is a check nobody knows works.

## Agent skills

### Issue tracker

Issues live in this repo's GitHub Issues (`socrates-tools/socrates-review`), read and written via the `gh` CLI.
See `docs/agents/issue-tracker.md`.

### Domain docs

Single-context: one `CONTEXT.md` and `docs/adr/` at the repo root.
See `docs/agents/domain.md`.
