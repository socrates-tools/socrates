# socrates

socrates is the one binary the suite publishes, and it answers to that one name.
`socrates review` puts an independent reviewer in front of unreviewed local changes.
Everything socrates review prints names `socrates review`.

`socrates agora` is a second namespace inside the same binary, not a second program.
It serves one document as one page on a loopback address, waits for a person, and prints what they answered.
It knows nothing about code review but the three audiences socrates review writes and the reading of the three severities its hints carry, per ADR 47, and nothing review does reaches it yet.

`socrates design` is a third namespace, and a skill rather than a tool.
The namespace holds one operation, `socrates design skill`, which installs its two skills: `socrates-design`, which runs the interview, and `socrates-tickets`, which cuts the spec it left into tickets.
A design interviews a person over agora and leaves a spec and its decisions behind, and it knows nothing about review or implement, which are the other workflows sharing this binary.

`socrates implement` is a fourth namespace, per ADR 48, and builds a change's tickets through implementers.
`socrates implement skill` installs `socrates-implement`, which turns the session that a person invokes it in into the overseer, and the `socrates-implementer` agent definition the overseer spawns.
It reads one configuration file of its own, `~/.socrates/implement/config.yaml`, `socrates implement config` reports every setting in it, and `socrates implement configure` changes one from a menu.
It knows nothing about review, agora or design.

The dispatcher that routes `socrates` to a namespace holds wiring and nothing else: `internal/dispatcher` mounts review's command tree, agora's command tree, design's command tree and implement's command tree, and is the only thing that imports all four.
Each tool reads one configuration file of its own: review reads `~/.socrates/review/config.yaml`, agora reads `~/.socrates/agora/config.yaml` and implement reads `~/.socrates/implement/config.yaml`, per ADR 42.
review and agora write under `~/.socrates/review`, and implement writes under `~/.socrates/implement`.
design reads and writes none of these: nothing about a design outlives the session that ran it beside the spec, the ADRs, the glossary entries and the tickets it wrote to the repository being designed for.
`~/.socrates` is the brand's root, holding one directory per tool published under it, and `SOCRATES_HOME` moves the root.

`CONTEXT.md` defines the vocabulary and `docs/adr/` holds the decisions, each with what it rejected and why.
This file is the map: the layout, the verification sequence, and the invariants the design rests on.
A change starts as a spec, committed to `docs/changes/<slug>/spec.md` with its tickets as files beside it, and `CONTRIBUTING.md` says how.
An issue is only where somebody using socrates reports a bug or asks for something.

## Scope

socrates review reviews.
It does not land, verify, render, serve, or wait for a person about a review.
`socrates review configure` is the one command that waits for a person, and what it waits about is review's own settings rather than anything a run decides, per ADR 25.
Committing, pushing, opening a pull request, posting a comment and watching continuous integration are not review's and never become review's.
It resolves a base, records what a person decided about a finding, and reports whether a finding's location is inside the change; fetching, rebasing and posting are the session's and the person's.
Running the reviewed project's own tests is not in the core layer either, so `clean` means an independent reviewer found nothing blocking, not that the tests pass.

## Layout

What the suite is:

- `cmd/socrates` - the binary: builds `internal/dispatcher`'s tree, runs it, and turns an error into an exit code.
- `internal/dispatcher` - the root command: routes `review`, `agora`, `design` and `implement` to their own trees, owns the error-to-exit-code mapping and the one failure prefix, and holds the top-level `version`. It is not a tool, and nothing above the trees it mounts names a vendor, a namespace's own flag, or anything any tree does.

What review is:

- `review/command` - the command tree: `begin`, `status`, `findings`, `respond`, `continue`, `watch`, `logs`, `cancel`, `skill`, `config`, `configure`, and the hidden `supervise`.
- `review/snapshot` - the freeze: a temporary index, a commit, a detached worktree, and the seal that makes it unwritable.
- `review/runstore` - run directories, `run.json`, run selection, rounds, retention.
- `review/runstate` - the states and which are terminal.
- `review/findings` - the findings contract, the reviewer's schema, the classification and the verdict rule.
- `review/dialogue` - what a round leads to, which rounds wait on a person, which resolutions a finding accepts, the round cap, and the response document.
- `review/hunks` - the lines a change touches, read from its own diff, and whether a finding's location is one of them.
- `review/hygiene` - model output neutralised before it re-enters a prompt.
- `review/guidance` - the one file review reads out of the repository under review, the commit it is read at, and the ceiling it is held to.
- `review/prompt` - the delivered prompt, assembled from named sections.
- `review/testquality` - the test-quality rule, shared by every role that needs it.
- `internal/platform` - the one place the supported platforms are stated, and the basis on which one qualifies.
- `review/agent` - the only package that starts a model process, the reviewer CLIs review drives and the basis on which one qualifies, and one file per vendor behind the seam.
- `review/reviewer` - beginning, supervising, responding, advancing, cancelling.
- `review/report` - what `status`, `findings` and `config` print.
- `review/configure` - the menu over review's settings: its labels and its refusals, built on `internal/settingsmenu`.
- `review/presentation` - the one place an event becomes a line a person reads.
- `review/follow` - following a run's event log live.
- `review/terminal` - the pager, the plain fallback, and whether a stream is a terminal at all.
- `review/skill` - the two skills a coding session reads, one for its own change and one for somebody else's, and installing them.
- `internal/exit`, `internal/version`, `review/config`, `review/runlog`, `internal/state`, `internal/refusal`, `internal/process` - the small owned things: exit codes, the linked build stamp, the one optional configuration file review reads and what every setting in force is, the JSONL log, the brand's root and each tool's directory under it, refusals, process groups.
- `review/corpus` - what real reviewers answered, verbatim and embedded, and what the routing rule does with it.
- `review/fidelity` - the recorded shapes of each real CLI, what was measured about each one's criteria, the drift check, the live smoke check and the forge check, all network behind a build tag.
- `review/simulation` - what every simulator shares: the scenario vocabulary, the invocation log, and the run that turns one into the other.
- `review/claudesim`, `review/codexsim`, `review/opencodesim` - each simulator's half of its real CLI's wire contract: the arguments it reads, the stream it writes.
- `review/cmd/claudesim`, `review/cmd/codexsim`, `review/cmd/opencodesim`, `review/cmd/corpus-outcomes`, `review/cmd/corpus-record`, `review/cmd/platforms` - the simulator and the gate binaries.

What agora is:

- `agora/command` - the command tree: `serve`, `poll`, `end`, `stop`, and nothing else.
- `agora/config` - the one optional configuration file agora reads, which holds the `surface` block and nothing of review's.
- `agora/document` - what a caller writes for the surface to render, what it copies rather than accepts, and the languages the page can highlight.
- `agora/sanitise` - a caller's HTML fragment reduced to a fixed allowlist, and every name it carries prefixed.
- `agora/stylesheet` - a caller's stylesheet read as a parse, reduced to a fixed allowlist, and scoped to the illustration.
- `agora/render` - a document to one HTML page, the marks a severity and an audience are drawn as, and a patch to the rows of a split diff.
- `agora/assets` - the built client bundle and fonts, embedded, and the stamp of their source.
- `agora/surface` - the loopback server and its guards, the address bound where nothing states one and the window walked down from it, the sessions it holds, what it writes about itself and about each page, the exchange, and the bounds on waiting.
- `agora/web/` - the client source, the design system and the fonts, built by `make assets`.

What design is:

- `design/command` - the command tree: one operation, `skill`, that installs both skills.
- `design/skill` - the source of both skills. The design's: the rounds, the page, what was settled, the glossary, the spec, the decisions, the last round, ending early, and installing elsewhere. The tickets': the spec it starts from, slicing, blockers and waves, the breakdown on a page, the last card, the main ticket, a ticket, and leaving them to be committed.

What implement is:

- `implement/command` - the command tree: `config`, `configure`, `default-branch` and `skill`, which installs the overseer's skill and the implementer's agent definition, and nothing else yet.
- `implement/configure` - the menu over implement's settings: its labels and its refusals, the model aliases it offers, the default branch row inside a repository, and rewriting the installed agent definition after a model or an effort is confirmed, built on `internal/settingsmenu`.
- `implement/skill` - the source of the `socrates-implement` skill and of the implementer's agent definition, installing both, and the record of where the definition was installed.
- `implement/config` - the one optional configuration file implement reads, every setting it holds, the default each falls back to, the default branch each repository's features start from and how a repository is named, by the top level of its main worktree, and the report `socrates implement config` prints.

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
None is a matter of care, and the reasoning for each is in `docs/adr/`.

**The repository under review**

- Nothing writes in it: no run state, no configuration, no log.
- Every round is snapshotted afresh against the run's original base, and snapshot commits are content-addressed.
- A base a person names resolves to the fork point of that ref and `HEAD`, is recorded once when the run begins and is reused by every round. A ref that does not resolve, and a ref with no merge base against `HEAD`, are refused before a snapshot is taken and before a reviewer is spawned.
- Nothing fetches and nothing rebases: no ref is written, no commit is rewritten, and the reflog is untouched by resolving a base.
- A write in the snapshot worktree is refused by the filesystem, and the worktree is compared against its commit after the round.
- One file is read out of it, `.socrates/review-guide.md`, at the run's base commit and never from the working tree or the snapshot, so a guide the change writes governs the next run rather than this one.
- That guide is bounded, neutralised on the way into the prompt, and additive: it withdraws no rule review delivers, and nothing review decides is derived from a word of it.

**The reviewer**

- `review/agent` is the only package that starts a model process.
- The reviewer CLIs review drives are stated once, each beside the facts that make it qualify, and support follows from those facts rather than from a list.
- A CLI qualifies on three things: the reviewed repository's own instructions can be kept out of the session, a terminal event says whether the work happened and what it said about a failure, and a model is selectable per invocation.
- Whether the vendor constrains output against the schema is a tier rather than a criterion, because review validates the document either way.
- A configured CLI review does not drive is refused by name, naming the guarantee it does not offer, before a snapshot is taken and before a reviewer is spawned.
- Nothing above the seam names a vendor, a flag or an event shape, and nothing outside review shapes a reviewer's answer into the findings document.
- Each vendor's invocation asks for none of the machine's own customisations; the session's tool list and named sources are recorded rather than trusted, and a vendor that reports no inventory produces runs that record that.
- A round after the first inherits no session: the only cross-round context is the labelled prior findings and recorded responses.
- Model output is neutralised per field on the way into a prompt, and never on the way to disk.
- What a person passes as the intent is neutralised on the way into the prompt whatever the run, and the section carrying it names the party the run records as having stated it.
- A review inside a review is refused.

**Verdicts and the dialogue**

- `clean` requires a findings document that parsed; a missing, dead or unparseable one is `failed`, and `failed` is never merged with `escalated`.
- A review of a working tree holding nothing its base commit does not is refused before a run is created, and a round of one ends the run `failed`, so `clean` is reported for no round a reviewer did not read.
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
- review states whether a round waits on a person and does nothing further about the page: no command, flag or run record names one, and it states nothing about whether a page is served.
- `review/presentation` is the one place an event becomes a line, and `internal/exit` owns the exit-code table.
- `status` prints the base commit the run resolved, as a bare fact and not a warning, and review says nothing about how far a named ref has moved past it.
- `config` names the file it reads and then every setting by its YAML key, with the value in force and whether a person set it, and a setting nothing states in force prints empty rather than invented.
  Each default it reports is read from the package that owns it, so a default has one home, and a test fails where a reported one and its owner's disagree.
  It reports and judges nothing: a reviewer CLI review will not drive is printed as the file states it, and refused where a review starts rather than here.
- `configure` draws a menu of the same settings `config` reports, read from the same place, so the two cannot disagree about what is in force.
  It writes only the key a person confirmed, as soon as they confirm it, and a person's comments, key order and blank lines survive the write.
  `use the default` removes a setting rather than writing today's default in, an emptied block goes with its last setting, and a value the loader would refuse is refused before the file is touched.
  It refuses where nothing is a terminal, naming `socrates review config` and no library, and it is the one command that opens a file the others refuse, naming the line it could not read and offering nothing but resetting that setting.
  A block that is wrong as a whole is named and reset as the block rather than by its first field, and what the menu says is wrapped at the terminal's width rather than cut off.
- A configuration file carrying a `surface` block is refused by every review command but `configure`, naming `~/.socrates/agora/config.yaml` as where the block belongs, before a snapshot is taken and before a reviewer is spawned.
  `configure` offers to reset the block, and nothing moves it into agora's file or rewrites the file unasked.
- `internal/settingsmenu` states every label the menu draws in its own `messages.go` and `review/configure` states review's own, and the wording check reads the plain string literals of both rather than only their `fmt` call sites.
- Everything above that is stated for `configure` holds for `socrates implement configure` too, which refuses naming `socrates implement config`, and `internal/settingsmenu` names neither tool: what differs is carried by the `Definition` a tool hands it.
- A setting whose written key is not its displayed key carries the segments it is written at, so a `default_branch` row is written under `default_branches.<repository top level>` whatever dots the path holds, and review's rows are written at their displayed key as before.
- `socrates implement configure` offers `sonnet`, `opus`, `haiku`, `fable` and `inherit` for `roles.implementer.model` and also takes a full identifier typed.
- A confirmed model or effort rewrites the installed agent definition at once, and with none installed the menu names `socrates implement skill` as what writes it.
- Inside a repository the menu offers that repository's `default_branch`, a branch the repository does not hold is refused before the file is touched, and `use the default` removes that repository's entry and no other.
- `findings --json` reports each finding's recorded resolution and whether its location is inside the change, and `review/hunks` is the one place that question is answered.
- A follower stops on the run's record or on the log, waits bounded for a log that does not exist yet, and stops on a parked round.
- `status`, `watch` and `logs` stay three commands.
- Every user-facing string lives in its package's `messages.go`, names what the person can do, and carries no em dash.

**The skill**

- Every skill any namespace ships carries `disable-model-invocation: true`, stated once in `internal/skill` and rendered into every skill's front matter, and the suite fails over every committed `socrates-*` skill and over everything `socrates <namespace> skill` installs, for each namespace in `convention.Tools` whose tree holds `skill`, rather than over a list, per ADR 49.
- `review/skill` is the one source of both skills, and each committed copy is generated from it.
- The loop for a change nobody here can patch is a second skill rather than a branch inside the first, and each skill carries a ceiling of its own that the suite reads from the skill rather than from a list.
- The second skill relays what somebody else wrote verbatim and acts on none of it, reads the change in a worktree of its own, records the head commit and refuses to post where it moved, takes a comment's placement from what review reported, and posts nothing the person did not choose.
- It names nothing the binary does not offer, lists every state a run rests in, teaches polling at a stated cadence, names `socrates review watch` once and outside the loop, says once that the address `socrates agora serve` prints holds the session `socrates agora poll` is given and that the session is kept for the whole review, and says whose a question is.
- Every branch that reached a result serves a page, `failed` and `cancelled` name none, and the session waits only where review says the round waits on a person.
- `dialogue.ResponseRules` and `skill.PageBar` are each stated once and rendered where they are needed.
- It says every control carries the command it becomes and that the page shows that command to nobody, and it names each decision from the finding's identifier rather than titling it.
- It says every finding in the round gets a card of its own, that a finding the session resolved states what resolved it and offers no options, and that what foregrounds the rest is the severity order rather than shortening, listing or omitting anything.
- The loop it describes is driven end to end through the real binary, not read.

**The design skill**

- `design/skill` is tested the way `review/skill` is: the rendered document is the subject, its ceiling is enforced, the committed copy must match what the package renders, and every invocation it names must exist in the binary's command surface with the flags it names, across both namespaces it covers.
- The last card offers write it and nothing that starts another skill, per ADR 49: every socrates skill is a person's to start.
- The skill writes no severity, no `None of these` option and no command on any option it composes: a design answer runs nothing, and a design question is earlier or later in the argument rather than more or less severe.
- A settled question is written to the document's `settled` list rather than as a decision, stating the question, the labels chosen and the person's own words, and a person reopens one through the page's message, per ADR 46.
- The ADR criterion the skill states is a conjunction, tested as one: hard to reverse, surprising without context, and the result of a real trade-off, never any one or two of the three.
- There is no round cap: a person can end a design at any round, on the page or with `socrates agora end`, and either way the exchange ends carrying only what was actually answered, leaving everything else for the session to record as unsettled.
- The loop is driven end to end through the real binary: a round renders with its recommendation chosen, and what the person answers, including approving the last round's card, reaches the session through `socrates agora poll`. What the binary cannot drive, because nothing in it writes a file, is the spec and the ADRs themselves: those are the session's to write from what `socrates agora poll` handed it.

**The tickets skill**

- `socrates-tickets` is a second skill of `design/skill`, installed by `socrates design skill` beside the first, held to a ceiling of its own and to the same gates: the committed copy matches its source, and every invocation it names exists.
- It says nothing about a design handing on to it and mints its own session, per ADR 49.
- It writes no severity, no `None of these` option and no command on any option, and names no `socrates review` invocation: how a ticket is implemented is the implementing session's own.
- It cuts from a spec under `docs/changes/<slug>/` and nothing else, refuses a tickets folder that already exists by name, and leaves what it writes uncommitted.
- A ticket names its blockers and its wave is derived from them, per ADR 44, and status is kept in the main ticket and nowhere else.
- The breakdown is driven end to end through the real binary: every card arrives with *keep it* chosen, and a card answered in the person's own words reaches the session through `socrates agora poll`. What the binary cannot drive is the tickets themselves, which are the session's to write.

**The implement skill**

- `socrates-implement` is rendered by `implement/skill`, installed by `socrates implement skill` under the skills directory or `--to` beside the agent definition, carries `disable-model-invocation: true` through `internal/skill`'s front matter, and is held to a ceiling of its own.
- The committed copy at `.claude/skills/socrates-implement/SKILL.md` must match what the package renders, and `make skill` writes it.
- Every invocation it names exists in the binary with the flags it names, and it names only `socrates implement` ones.
- It names no `socrates review` invocation and no other skill, and tells the overseer to invoke no skill, per ADR 49.
- The overseer writes no code, edits only the main ticket's status table, and settles every coding question by her own judgement.
- She reads the main ticket, the spec and the ADRs, and the settings from `socrates implement config`, at the start of every run.
- A default branch that is not remembered is proposed from the branches named `main`, `master`, `dev` or `develop`, never one under a prefix such as `feature/` or `hotfix/`, asked in the terminal and saved with `socrates implement default-branch`, per ADR 51.
- The change's branch is cut from the current `HEAD` and starts with what the design and the tickets left, and anything else uncommitted stops the run.
- The frontier is every ticket whose blockers are done, at most `concurrency` implementers run, the longest chain goes first, then the most unblocked, then name order, and a ready ticket sharing a surface with a running one is held back without a blocker being written.
- Each implementer is briefed with pointers, runs in the background in a worktree and on a branch the overseer made with git, and gets feedback from the overseer with no cap on the exchanges.
- Landing is a fast-forward followed by the repository's own verification, a failure goes to the implementer whose ticket landed last, and where the only commits an implementer's branch lacks are the overseer's own status commits she merges them into its branch herself, while any other new commit goes back to the implementer to merge.
- In progress is committed before a hand-out and done once a ticket has landed and passed, and a resumed overseer continues each in-progress ticket from its own branch with a fresh implementer.
- What the person says mid-run is relayed at once to every implementer it concerns, and the person is reached in the terminal only about what coding cannot fix.
- The whole change is read against the spec from the default branch's fork point, every worktree is removed, and nothing is pushed.
- The orchestration itself is the session's to perform from what the skill says: the suite asserts each promise as a statement in the rendered document and every command the skill names against the real binary.

**The implementer's agent definition**

- `implement/skill` renders `socrates-implementer.md` once, and the committed copy at `.claude/agents/socrates-implementer.md` must match what it renders at the defaults.
- Its `model` and `effort` are the settings in force, and `background` is true.
- It names no `socrates` command, read against every namespace in `convention.Tools`.
- It is held to a ceiling of its own.
- Where it was written is recorded in `~/.socrates/implement/installed.yaml`, a second install elsewhere replaces the record, and `config.yaml` never holds it.

**The default branch**

- A repository is keyed by the top level of its main worktree, so every worktree of one clone reads one entry and a second clone at another path reads none.
- A branch the repository does not hold, local or remote-tracking, is refused by `socrates implement default-branch` and the file is untouched.
- Outside a repository and in a bare one, `default-branch` is refused naming what it needs, and `config` reports no `default_branch`.
- A save keeps a person's comments, key order and blank lines, through the same machinery as the menu, and a repository path holding a dot is one key.

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
- The matrix is read from the platform statement, so the suite runs on every platform the set names and on no runner it does not.
- Every action a workflow uses is named by a commit with its version beside it, and a check fails while any is named by a tag.
- Every workflow states a read-only token at the top and escalates only in the job that needs more, and a superseded run is cancelled rather than left to report.
- A run happens when somebody asks for one and when a release gates on one, and nothing starts on a push: what that gives up is that nobody learns the default branch is red, or that an advisory was published, until somebody asks.
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

**The surface**

- It binds a loopback address; nothing widens it, and an address that is not one is refused before any server starts.
- It takes the first free port of a fixed window counting down where nothing states one, and a whole window held is a refusal naming it.
  The window is stated once, in `agora/surface` beside the guards that judge an address.
- A port an operator committed to, by flag or by either field of the `surface` block, is bound or refused and never walked past.
- The surface mints a session, `--session` never brings one into existence, and every page and every route of it is served under a path holding that session.
- It turns away a request naming a host it did not bind, and every state-changing route carries a token: a session's own answers for that session and no other, and the one minted when the server started answers everywhere.
- It composes a command and never runs one, the command it hands back is the option's own string, and no command is rendered on the page.
- It writes only under review's home, in a directory of its own beside the runs.
- What an operator states about it is read from `~/.socrates/agora/config.yaml`, in agora's own directory under the brand's root, and from no other file, and a file it cannot read or a field the `surface` block does not have is refused, naming the file.
- A caller writes a document rather than a page, and the name, severity, audience and location come from the source a decision cites.
- A document may hold no decision, and one holding no decision, no source entry, nothing settled and no preamble is refused.
- A decision may carry no options, and such a card states something and offers no control a person could mistake for a choice; one carrying no options and stating a mode is refused.
- Every entry of the source is rendered ahead of the decisions, in an index beside them that names every decision where there is no source, and both are ordered by the severity the source carries, most severe first, informational last, an unstated severity after all of them, and ties in the order the document listed them.
- The index counts a decision answered by the rule that decides what is sent, as the person answers it, so a decision whose recommendation arrives chosen arrives counted.
- What was settled is rendered ahead of every open question, offers no control, is counted by nothing and is never sent, and an identifier is either settled or still open, per ADR 46.
- The order is the surface's, so no field of a document reaches it, and a page with nothing to answer offers neither control and says so.
- Every decision carries a free-text field that no document field can remove, including one that offers no options.
- The option a document marks arrives chosen, a decision that marks none arrives with nothing chosen, and a decision offering options also offers taking none of them, which is not an option, carries no command and leaves the decision with nothing chosen.
- What the controls hold at the submit is what is sent, so nothing is committed per decision, a decision holding no selection and no words is omitted, and a decision left blank is named before either press sends.
  Either press over a blank decision raises a modal dialog naming every one of them and takes no armed state, a press inside it is what sends, what is behind it is dimmed, and each of its names closes it and puts the cursor on that card's first control.
  Going back, Escape and a press on the backdrop are one path: nothing is sent, the page is as it was, and the cursor returns to the button that raised it.
- The languages the page highlights are the bundled library's common set, stated once and generated from `web/`, and a document naming one outside it is refused while one naming none renders as plain text.
- A caller's fragment is reduced to a fixed allowlist before it is embedded, and every name it carries is prefixed by the surface.
- A caller's stylesheet is reduced from a parse that reached the end of its input, every selector it keeps is scoped to the illustration, and the page admits it by a hash of exactly the bytes it embedded.
- What a caller paints is contained by the illustration, so a rule inside one decides nothing about where the page's own controls land.
- Delivery happens once, waiting is bounded, collection is per session, and the server stops once it has held no session for the idle period.
- An exchange holds at most one uncollected answer, and a later submit replaces it rather than queueing behind it, so what a caller collects is a person's last word.
- Collecting an answer is visible to the page, so the way back from a receipt is offered exactly while it would still reach the caller.
- Nothing carries across a revision: not a selection, not a typed word, not the message, not the receipt and not an open confirmation, because a decision's identifier is unique only inside the document that names it.
- A stream is told where the exchange stands when it connects rather than only when something next changes, so a client that missed a revision catches up rather than waiting on a change that has already happened.
- The end of an exchange stops the page offering to answer and not the page listening, so a caller who writes again reaches the page the person is already looking at, and only the test about reloading reloads before it asserts.
- Every asset is embedded, no diagram engine is, and the committed bundle is generated from `web/`.
- An element it hides paints nothing, every colour it paints is a token stated once, and every text meets a contrast floor under both preferences.
  The page is dark whatever the reader prefers, so both preferences are painted alike, per ADR 45.
- A severity the surface orders by is marked by what it means, each audience socrates review writes is marked as the design drew it and any other in the caller's own words, and every mark says what it means on hover, on focus and to a screen reader, per ADR 47.
- Nothing in it names review, finding, verdict, round, or run as a noun, and the module path is not part of its vocabulary, except the three audiences socrates review writes, stated in `agora/render/audiences.go` and exempted by name, and the severity hints, which carry review's reading of a severity in words the check does not know, per ADR 47.

**The suite**

- Every end-to-end scenario runs against every supported vendor's simulator, and a supported vendor with no simulator or no recordings fails the suite.
- Tests assert on the delivered prompt, never on review's own source text, and the prompt's reporting section is generated from `review/findings`.
- A test may not prove behaviour by reading implementation source; one that legitimately reads a Go file is named in `convention.ContractsTestsRead`.
- A gate that cannot check what it was asked to check fails rather than reporting a pass.
- What a real reviewer answered is kept verbatim and replayed; the corpus is classified by `dialogue.After`, and a recording is only made from a run that answered.
- A verification task may require that a live run happened and was recorded, never a particular severity, action or count.
- The page half is driven by a real browser against the real binary; asserting on served markup would test the renderer.
- The menu half is driven by the real binary in a real pseudo-terminal, typed at and read back from what was drawn; exercising the model in process would pass on a build where `socrates review configure` cannot start.
- Nothing in continuous integration reaches the real CLI or a forge.
- Where a forwarded finding's comment lands is checked against a real forge once, by hand, behind a build tag: the findings come from a simulator, because what that check reads is the forge. It posts to one long-lived pull request and deletes what it posted.
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

### Domain docs

Single-context: one `CONTEXT.md` and `docs/adr/` at the repo root.
See `docs/agents/domain.md`.
