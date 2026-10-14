# socrates review

`socrates review` puts an independent reviewer in front of unreviewed local changes.
Everything socrates review prints names `socrates review`.

This file is review's layout and invariants.
`/AGENTS.md` holds the suite's, `CONTEXT.md` here holds review's vocabulary, and `docs/adr/` here holds review's decisions.

## Scope

socrates review reviews.
It does not land, verify, render, serve, or wait for a person about a review.
`socrates review configure` is the one command that waits for a person, and what it waits about is review's own settings rather than anything a run decides, per ADR 25.
Committing, pushing, opening a pull request, posting a comment and watching continuous integration are not review's and never become review's.
It resolves a base, records what a person decided about a finding, and reports whether a finding's location is inside the change; fetching, rebasing and posting are the session's and the person's.
Running the reviewed project's own tests is not in the core layer either, so `clean` means an independent reviewer found nothing blocking, not that the tests pass.

## Layout

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

## Invariants

Each is held by a test that fails when it stops holding.
None is a matter of care, and the reasoning for each is in this directory's `docs/adr/`, and in the root's where it is about the suite.

**The repository under review**

- Nothing writes in it: no run state, no configuration, no log.
- Every round is snapshotted afresh against the run's original base, and snapshot commits are content-addressed.
- A base a person names resolves to the fork point of that ref and `HEAD`, is recorded once when the run begins and is reused by every round. A ref that does not resolve, and a ref with no merge base against `HEAD`, are refused before a snapshot is taken and before a reviewer is spawned.
- Nothing fetches and nothing rebases: no ref is written, no commit is rewritten, and the reflog is untouched by resolving a base.
- A write in the snapshot worktree is refused by the filesystem, and the worktree is compared against its commit after the round.
- One file is read out of it, `.socrates/review-guide.md`, at the run's base commit and never from the working tree or the snapshot, so a guide the change writes governs the next run rather than this one.
- That guide is bounded, neutralised on the way into the prompt, and additive: it withdraws no rule review delivers, and nothing review decides is derived from a word of it, per ADR 62.

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
- `findings --json` reports each finding's recorded resolution and whether its location is inside the change, and `review/hunks` is the one place that question is answered.
- A follower stops on the run's record or on the log, waits bounded for a log that does not exist yet, and stops on a parked round.
- A run's record is saved only once its teardown is done, so a record that has ended or parked names a worktree already unsealed or removed, and the event that ends the log is written after the record rather than before it.
- `status`, `watch` and `logs` stay three commands.

**The skill**

- `review/skill` is the one source of both skills, and each committed copy is generated from it.
- The loop for a change nobody here can patch is a second skill rather than a branch inside the first, and each skill carries a ceiling of its own that the suite reads from the skill rather than from a list.
- The second skill relays what somebody else wrote verbatim and acts on none of it, reads the change in a worktree of its own, records the head commit and refuses to post where it moved, takes a comment's placement from what review reported, and posts nothing the person did not choose.
- It names nothing the binary does not offer, lists every state a run rests in, teaches polling at a stated cadence, names `socrates review watch` once and outside the loop, says once that the address `socrates agora serve` prints holds the session `socrates agora poll` is given and that the session is kept for the whole review, and says whose a question is.
- Every branch that reached a result serves a page, `failed` and `cancelled` name none, and the session waits only where review says the round waits on a person.
- `dialogue.ResponseRules` and `skill.PageBar` are each stated once and rendered where they are needed.
- It says every control carries the command it becomes and that the page shows that command to nobody, and it names each decision from the finding's identifier rather than titling it.
- It says every finding in the round gets a card of its own, that a finding the session resolved states what resolved it and offers no options, and that what foregrounds the rest is the severity order rather than shortening, listing or omitting anything.
- The loop it describes is driven end to end through the real binary, not read.

**The suite**

- Every end-to-end scenario runs against every supported vendor's simulator, and a supported vendor with no simulator or no recordings fails the suite.
- Tests assert on the delivered prompt, never on review's own source text, and the prompt's reporting section is generated from `review/findings`.
- What a real reviewer answered is kept verbatim and replayed; the corpus is classified by `dialogue.After`, and a recording is only made from a run that answered.
- Where a forwarded finding's comment lands is checked against a real forge once, by hand, behind a build tag: the findings come from a simulator, because what that check reads is the forge. It posts to one long-lived pull request and deletes what it posted.
