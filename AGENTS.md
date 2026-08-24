# socrates

socrates puts an independent reviewer in front of unreviewed local changes.
The binary is `scrat`.
The project is `socrates`.

`openspec/` is normative.
Where this file and `openspec/` disagree, `openspec/` wins, and this file is what needs fixing.

## The scope rule

socrates reviews.
It does not land, and it does not verify.
Committing, pushing, opening a pull request and watching continuous integration are not socrates' and never become socrates'.
Running the reviewed project's own tests is not in the core layer either, so `clean` means an independent reviewer found nothing blocking, not that the tests pass.

## Layout

- `cmd/scrat` - the one binary a person runs: `review`, `status`, `findings`, `respond`, `continue`, `watch`, `logs`, `cancel`, `skill`, and the hidden `supervise` that a review spawns.
- `cmd/claudesim` - the simulator, which speaks the vendor CLI's wire contract and is symlinked as `claude` onto `PATH` by the test harness.
- `cmd/scenario-coverage` - the coverage gate behind `scripts/scenario-coverage`.
- `cmd/scenario-drop` - the promotion gate behind `scripts/scenario-drop`.
- `cmd/corpus-outcomes` - the report behind `scripts/corpus-outcomes`, and the `-write` that rewrites the recorded distribution.
- `cmd/corpus-record` - records one completed run's answer into the corpus.
- `internal/snapshot` - the freeze: a temporary index from `HEAD`, `add -A`, `write-tree`, `commit-tree`, a detached worktree of the result, and the seal that makes that worktree unwritable.
- `internal/runstore` - run directories, `run.json` written by rename, run selection, rounds, and the retention sweep.
- `internal/runstate` - the run states and which of them are terminal.
- `internal/findings` - the findings contract, the schema handed to the reviewer, the classification and the verdict rule. It imports nothing outward.
- `internal/dialogue` - what a round leads to, as one pure rule over the classified findings and the round number; which rounds a person reads, as one pure rule over the classified findings and the state that round led to; the round cap and the reason for it; the response document and which of its three resolutions each finding accepts. It imports nothing outward.
- `internal/hygiene` - neutralising a field of model output before it is re-embedded in a prompt.
- `internal/prompt` - the delivered prompt, assembled from named sections in one fixed order.
- `internal/testquality` - the test-quality rule, as one block, shared by every role that needs it.
- `internal/agent` - the only package that starts a model process or parses its output.
- `internal/reviewer` - beginning a run, supervising the round, recording a response, advancing to the next round, cancelling.
- `internal/report` - what `status` and `findings` print, for a person and as JSON.
- `internal/presentation` - the one place an event becomes a line a person reads. A function from events to lines, and it imports nothing that could reach anything.
- `internal/follow` - following a run's event log live: the two stopping conditions, the bounded wait for a log that does not exist yet, and reading what the log holds once.
- `internal/terminal` - whether the output is a terminal, what `PAGER` says, and the plain fallback where a pager will not start.
- `internal/skill` - the skill a coding session reads to drive the loop itself: its text as the one source, what the committed copy is generated from, the commands and flags it names, the bar a page a person reads is held to, and installing it at user level.
- `internal/exit` - the exit-code table.
- `internal/config` - the one optional configuration file under the user's configuration directory.
- `internal/runlog` - the JSONL event log, its byte ceiling and its in-place rotation.
- `internal/convention` - the repository's own conventions as executable checks: no commentary, the test-quality rule, the message registry, what socrates reaches and what is compiled out of the ordinary suite.
- `internal/coverage` - scenario coverage: what `openspec/specs` and the changes under way promise against what the tests name. An archived change is skipped, because its specification was promoted.
- `internal/promotion` - what archiving a change would do to the promoted specifications: a `MODIFIED` requirement that carries fewer scenarios than the one it substitutes, one that substitutes nothing, one added beside a requirement that already exists, and a declared scenario rename the promoted specification already holds the new name of.
- `internal/corpus` - what real reviewers answered, stored verbatim beside the circumstances it was answered in, embedded so replay needs no path. `runs` takes an answer out of a completed run; `outcomes` classifies the corpus through the routing rule and says which outcomes only fixtures reach.
- `internal/fidelity` - the recorded shapes of the real CLI, the drift check, and the live smoke check. Everything that reaches the network is behind a build tag.
- `internal/e2e` - the end-to-end suite: the real binary, a real repository, real git, real processes, the simulator on `PATH`.
- `internal/testing` - the harness, the temporary repository, the simulator installer, and source helpers. Test support only.

## Verification

Run this sequence before reporting work done:

```sh
make fmt
make lint
make test
make scenario-coverage
make scenario-drop
make corpus-outcomes
```

`make lint` fails on anything `gofmt` would rewrite, so formatting is never a separate thing to remember.
`make skill` rewrites the committed copy of the skill from its source, and `make test` fails while the two differ.
`make scenario-coverage` fails when a scenario belonging to a change somebody has started is named by no test, and prints where each gap is.
`make scenario-drop` fails when archiving a change would take a scenario out of the promoted specifications, which a `MODIFIED` delta does by substituting a requirement and carrying fewer scenarios than it holds.
It also fails where a delta declares a scenario rename whose new name the promoted specification already holds, because there is nothing left for the rename to be checked against and the comparison would report a clean promotion.
With no change under way it says so and passes, because what the rule does is held by tests over a fabricated tree rather than by whatever happens to be proposed today.
`make corpus-outcomes` reports nothing but the truth about coverage: for each outcome the routing rule can produce, how many answers real reviewers gave reach it and how many tests in the suite name it.
Naming is not asserting, and the report says naming, because an abstract-syntax walk cannot tell a test that requires an outcome from one that requires its absence.
It is a statement rather than a gate and never fails for an outcome no reviewer has reached, because whether one ever appears is not this repository's to control.
`make test` holds both gates and the comparison of the corpus against `internal/corpus/outcomes/distribution.txt`, so a scenario nobody tests, a delta that would drop one, and a routing change that moves what the corpus reaches all fail the ordinary suite.
Rewrite that distribution with `scripts/corpus-outcomes -write` once the move is deliberate.
Neither gate reads prose: a delta that rewrites a requirement's sentences from a stale copy is still something only a reader catches.

These reach the network and are not part of the ordinary sequence:

```sh
make fixtures       # record what the real claude CLI answers, as shapes
make fixtures-diff  # compare the simulator against those recordings
make live-smoke     # one real review against the real binary
```

The scheduled workflow in `.github/workflows/fidelity.yml` runs them weekly rather than per commit, because a failure there means the vendor's CLI changed rather than that a change here broke something.

## Load-bearing invariants

Each of these is held by a test that fails when the invariant stops holding.
None of them is a matter of care.

- **Nothing writes in the repository under review.** The only thing that reaches it is the snapshot, and the test that `git status` is byte-identical afterwards is what holds it. No run state, no configuration and no log lives inside a reviewed repository.
- **`clean` requires positive evidence.** A verdict of `clean` needs a findings document that parsed. A missing document, a dead reviewer, an unparseable document and a stream with no result event are all `failed`, and `failed` is never merged with `escalated`.
- **One function decides a finding's effective action.** `findings.EffectiveAction` is it, and an abstract-syntax check fails the suite when a second reader of the raw field appears.
- **One package starts a model process.** `internal/agent` is it, and a boundary test fails when another package constructs one or names the vendor's binary.
- **The reviewer is read-only structurally, not by trust.** Three mechanisms hold it, in the order they act: the vendor's tool restriction, which covers the tools the model may call and nothing else; the seal on the snapshot worktree, which refuses a write from anything in that session; and the comparison of the snapshot worktree against the snapshot commit after the round, which catches whatever defeated the seal. The machine can take the first two away and did, which is why the third stays.
- **A write in the snapshot worktree is refused by the filesystem, not by good behaviour.** The worktree is sealed after checkout by clearing every write bit in it, the seal is proved by one write before a reviewer is spawned, and a filesystem that allows that write is a refusal rather than a review socrates cannot stand behind. The seal is not a sandbox: the session owns those files, so something determined to write can clear the bit back, and that is what the comparison after the round is for. The seal is lifted when the round ends, so a worktree kept for inspection is an ordinary directory.
- **The reviewer loads user settings only.** The reviewed repository's own agent instructions are input authored by the thing under review, and a reviewer that reads them can be seated by them.
- **No customisation the machine brings on its own account reaches the reviewer.** A plugin and an external tool source put a tool socrates never asked for in the reviewer's session, and one of them wrote a directory into the snapshot worktree on socrates' first live run. A memory file and a skill contribute no tool and seat the reviewer instead, and a hook contributes no tool and runs a command of its own. `--safe-mode` is the one mechanism: it disables the ones the machine brings on its own account, it leaves auth, model selection and the built-in tools working, and unlike the flags that ask to replace a managed set it is not refused where one is deployed. It arrived in the vendor's CLI in 2.1.169 and socrates keeps no fallback, so the floor is stated where the CLI is named as a prerequisite and an older CLI fails every round on an unknown option. A policy-configured tool source does not load, which socrates measured on such a machine, and the vendor documents managed plugins, managed skills and managed memory going the same way. The mechanism is not trusted: socrates reads the session's tool list from the stream's opening event, subtracts the tools it asked for, and records the remainder, so one check covers every source of a tool rather than the one it was written for, and it reads the sources that event names beside it. Where such a tool is in the session anyway the review still runs, because a locked-down machine should still get reviews, and the run records it. The seal holds the reviewer read-only inside the snapshot worktree and bounds nothing outside it, which is why the tool surface is a mechanism rather than something the seal answers. What the flag does not reach is a hook an administrator's policy configures, which the vendor keeps in force by design: it contributes no tool, so the tool list is not what would catch it; its command is a child of the session and runs in the sealed worktree, where the seal and the comparison do; and where that command reaches outside the worktree nothing socrates holds bounds it. That one residual is stated rather than covered.
- **A tool list socrates could not read is recorded, not passed.** An opening event that reports no tools at all leaves the subtraction empty, so a run that says nothing would be indistinguishable from a session holding exactly what socrates asked for. A round that completed on a stream with no tool list records that it did, under its own kind, and the round does not fail for it: a release that stops emitting the field should report the truth rather than fail every review, which is the direction the built-in case is already wrong in. A round that failed before answering records the failure and not the absence, because the failure is what there is to read.
- **A source the session named is recorded, whatever its tools did.** The tool list and the source list come from the same opening event and are blind in different places: a source is named before its tools arrive, which socrates measured on an event naming four while the list held none of their tools, and a built-in the vendor added belongs to no source. Both are read, every source is recorded because socrates asked for none, and neither reading fails the round. That field is worth reading where the same event's inventory of what is installed on the machine is not, and the difference is measured rather than assumed: the sources named went from four to none under the one mechanism, while the inventory beside them did not move. Unlike an absence it is recorded whether or not the round completed, because it is something the session said rather than something it failed to say.
- **A gate that cannot check what it was asked to check says so.** The drop gate follows a scenario rename a delta declares, and there is nothing to follow where the promoted specification already holds the new name: the comparison sees everything carried and reports a clean promotion. That state fails the gate now, because it is what editing the promoted specifications while the change is still open does, and a check nobody can tell apart from a passing one gets run and ticked.
- **One rule decides what a round leads to, and nothing about a finding's action escapes it.** `dialogue.After` is it, tested against its whole input space rather than a sample. A round parks where every finding that holds the change can be answered, and `dialogue.Answerable` is the one owner of that set, so what the run waits on is what parked it. `escalated` is reached by a finding that holds the change while its own action asks for nothing, and by the round cap, and by nothing else: those are the two cases where socrates is finished and nothing more it can do.
- **One rule decides which rounds a person reads, and nothing re-derives it.** `dialogue.APersonReads` is it: `escalated` always, a parked round holding a finding a person answers always, and nothing else. It takes the state that round led to rather than the round number, because everything that reads it holds a state already decided and recorded, and a rule that recomputes an answer somebody has stored can disagree with the stored one. `scrat findings` reports it as a line and a field and `scrat status` does not, because status is polled and every line it prints lands in a waiting session's context. Both facts a re-derivation would read are already reported, which is what makes the second reader the cheap thing to write and the expensive thing to be wrong about: a round nobody renders is not a round anybody notices was not rendered.
- **socrates states that a round is one a person reads, and does nothing further about the page.** It renders nothing, publishes nothing, stores nothing and verifies nothing, no command or flag names a page, and no run record carries one: a link socrates neither made nor can refresh is a field that goes stale. The publishing is the session's reach in the harness the person is already driving, socrates gains no network reach from it, and a check over the whole module holds that. What that leaves outside socrates is a findings document, which quotes the code it was written about, on a page hosted elsewhere, and `internal/corpus` refuses to record for that same reason: the asymmetry is recorded as https://github.com/rui-ktei/socrates/issues/1 rather than settled.
- **The bar a page is held to has one source, and where it is drafted follows what a token costs at that branch.** `skill.PageBar` is the source, the skill renders it once, and neither branch restates it, because two copies of a bar drift and the copy in front of the session is the one that gets softened. The page is drafted in the session where the run escalated, since nothing follows that would re-read the drafting and a fresh subagent would cost its own start-up for no benefit, and delegated where the round parked, since the session is about to patch, record one answer per finding, continue and read another round, and everything added now is read again on every one of those turns. Nothing about the page is in the polling step.
- **The page collects nothing, and the identifiers are what join it to the terminal.** There is nothing behind it, so a form, a button, a checkbox or a stored answer would take a decision that reaches no run while the person believed they had made it. Every finding keeps the identifier, severity and audience `scrat findings` gave it, and each decision says what answering it looks like in the terminal, because a page that renamed a finding to something more readable hands the person back the search problem the page existed to remove.
- **A finding that asks a person parks the round rather than ending the run.** A machine arguing with intent is the thing this design refuses to build, and what carries that refusal is which resolution the finding accepts rather than a terminal state. Ending the run is what left the answer nowhere to go: the author decided, patched what they decided, and no round remained to judge either. A round holding one question and five mechanical defects waits on all six.
- **A rejection is a claim, never a suppression, and so is a decision.** The author may reject a finding, the reason is required, and both the finding and the reason travel to a fresh reviewer that judges them. No prompt socrates builds ever tells a reviewer to refrain from reporting a finding because it was found, rejected or answered before, and a test fails when one does.
- **A decision is a claim rather than a shortcut.** `scrat respond --decided` is the only resolution a finding that asks a person accepts, and it is accepted for no other finding, so answering an intent question by asserting a patch requires stating something false rather than choosing the loosest available word. socrates cannot verify that a person rather than the session answered and does not pretend to: it records the claim, labels it as a decision in the next round's prompt, and hands it to a reviewer that owes it nothing. What replaced the terminal state that used to make this impossible is a required reason in a durable record, and a sentence in the skill that a test asserts is there.
- **A round after the first inherits no session.** Each round is a new process, and the only cross-round context is the labelled block of prior findings and recorded responses. A reviewer resuming the round whose findings prescribed the patches would be seating the prescriber as the certifier.
- **The round cap is two, and it is for context rather than for cost.** The reason is recorded beside the constant, nothing configures it, and a configuration file that names one is refused.
- **Every round is snapshotted afresh against the run's original base.** A later round must see the patches as commits it can diff rather than as a description of them, and diffing against the previous round's snapshot would read a patch out of the context it was applied to. A base that has become unreachable fails the round naming it.
- **Model output is neutralised per field on the way into a prompt, and never on the way to disk.** `internal/hygiene` is it, and what it removes includes the heading marker socrates' own sections are built from, because a field that opens a section is a session under review composing the prompt that reviews it. The stored findings and responses hold what was actually said.
- **A recorded response is appended, never rewritten.** The author is a session that runs commands at once, and a read-append-rewrite of the whole document loses whichever claim was written first, silently. `responses.jsonl` is append-only and the last response to a finding is the one that counts.
- **A review inside a review is refused.** The reviewer's environment carries `SOCRATES_RUN`, and `scrat review` refuses when it is already set.
- **What a real reviewer answered is kept, verbatim, and replayed rather than invented.** `internal/corpus` holds every live answer beside the model, effort, round, base, snapshot, run and subject it was an answer to, and a scenario naming one is answered with those bytes and no others: the simulator does not adapt a recorded document to the prompt, renumber its findings or supply a field it lacks. Recording refuses any repository that is not socrates, because a findings document quotes the code it was written about.
- **The suite says which outcomes only fixtures have reached.** The corpus is classified by `dialogue.After` rather than by a second reading of the rule, the report names any outcome no recorded answer reaches, and it never fails for one. The distribution is checked in, so an entry that disappears and a routing change that moves the corpus are both a diff somebody reads.
- **A verification task may not require a particular answer.** A task about a live run may require that the run happened, that its answer was recorded and that socrates routed whatever came back. It may not require a severity, an action, a count or a shape: what a reviewer finds is the reviewer's, and four live reviews that closed no such task are what this rule exists to prevent repeating.
- **One place turns an event into a line a person reads.** `internal/presentation` is it, and a check fails when any other package names an event kind, unless that package is named in `convention.EventKindsNamed` with why it reads the log as data rather than rendering it. The renderer imports nothing beyond the event type and text formatting, so it cannot reach the filesystem, the clock or a terminal even by accident.
- **A follower stops on the run's record and on the log, and either is enough.** Waiting only for a terminal event is the version that hangs, because a run that failed before its reviewer wrote anything has a terminal state and an empty log. A parked round stops a follower too: nothing of that run is running, and only the author ends it.
- **A follower started too early waits, bounded, rather than failing.** `scrat review` returns as soon as the run exists, so the window between that and the reviewer's first byte is exactly when a person is most likely to run `scrat watch`. The wait is bounded, and a log that never appears is reported as a failed run.
- **`status`, `watch` and `logs` stay three commands.** `status` prints into a language model's context and the other two print into a person's terminal, which is an order of magnitude difference in what verbosity costs. A test reads the flags of each, so `status` cannot grow a way to follow and `watch` cannot grow a way to return after one read.
- **The skill has one source, and the committed copy is generated from it.** `internal/skill` holds the text, `make skill` writes `.claude/skills/scrat/SKILL.md`, and a test fails while the two differ. A committed file that is read rather than generated fails the way that matters: somebody edits it, nobody runs the generator, and the installed skill and the repository's copy disagree with nothing failing.
- **The response rules are stated once and rendered twice.** `dialogue.ResponseRules` is the source, `scrat findings` prints it where a round parks, the skill states it before the session has read a finding, and a test compares the two. Two copies of five sentences drift, and the copy in front of the session is the one that gets softened.
- **The skill names nothing the binary does not offer.** Every command and flag its text names is extracted and checked against the built binary, and removing one from the surface fails the check. A skill naming a flag that does not exist is a session doing something plausible and wrong at the point where the person has stopped watching.
- **The skill teaches polling, and names `scrat watch` once.** `scrat status` is the command in the session's loop, because everything a session reads lands in the context it is about to need for patching. The one mention of `scrat watch` is a suggestion for the person in another terminal, it sits outside the loop, and a test fails on a second mention or on one inside the loop.
- **The skill names a polling cadence and says the wait is socrates' to bound.** A session given neither invents both, and both inventions are wrong: a fast poll re-reads a document nobody has written, and a session's own timeout reports a live review as stuck. The reviewer's deadline ends the run as `failed`, so every run reaches a state to act on.
- **The skill says whose a question is, that the round is spent without asking, and that there is no third one.** A parked round used to hold nothing but the session's own work, and now it can hold a finding only a person closes, so the skill splits the branch and says the session does not answer one on its own behalf. It also says that a round every answerable finding has an answer for is spent without asking the person, because the information that decision needs does not exist until the patches have been read, and that a further look after the cap is a fresh `scrat review` the person asks for rather than a third round.
- **The skill names every state a run rests in.** The list is compared against `internal/runstate`, `cancelled` included. A loop told to wait for a fixed list of states polls through the one state a person creates, and the state a person creates is the one nothing else will resolve.
- **The loop the skill describes is driven end to end, not read.** Every branch it names is walked through the real binary against the simulator, and each journey asserts that nothing socrates ran wrote in the repository and that no command outside what the skill names was used. Text fails quietly, and a skill whose sequence nobody executes is prose.
- **The prompt is a generated interface.** Tests assert on the prompt that was delivered, read back from the round directory or from the simulator's invocation log, and never on socrates' own source text.
- **The reporting section is derived from the classifier.** The severities, actions and risk levels the prompt defines are generated from what `internal/findings` implements, so the two cannot drift.
- **Snapshot commits are content-addressed.** Author and committer dates are pinned to the epoch, so the same worktree and base always produce the same commit, and therefore a byte-identical prompt.
- **Every user-facing string lives in its package's `messages.go`.** A refusal names the condition and what the person can do about it, and no message carries an em dash.
- **The exit codes are a table one function owns.** `internal/exit` is it: `0` clean, `1` escalated, `2` failed, `3` usage, `4` still working.
- **The source carries no commentary.** A directive a tool parses and a `// scenario:` marker are the only comments the check allows.
- **Every scenario under way is named by a test.** A test function name, a subtest name or a `// scenario:` marker naming the scenario is what counts.
- **socrates holds itself to the test-quality rule it hands the reviewer.** A test may not prove behaviour by reading implementation source. Where a test reads a Go file legitimately, the contract it reads is named in `convention.ContractsTestsRead`.

## Conventions

- **No comments in source.** Names carry the meaning. `internal/convention` fails the suite on anything else.
- **One sentence per line in Markdown.**
- **No em dashes anywhere.**
- **Commit messages say what the change is**, with no rationale essay and no AI attribution trailer.
- **Tests name their scenario.** Prefer a test name that contains the scenario text; use a `// scenario:` marker where the honest test name would not.
- **Test first.** A task is not done until a test named for the scenario failed for the right reason and then passed.
- **A task about a live run says record it, never return this.** Require the run, the recording and the routing; never a severity, an action or a count.
- **Mutate a guard you rely on.** A check nobody has broken on purpose is a check nobody knows works.
