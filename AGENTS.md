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

- `cmd/scrat` - the one binary a person runs: `review`, `status`, `findings`, `cancel`, and the hidden `supervise` that a review spawns.
- `cmd/claudesim` - the simulator, which speaks the vendor CLI's wire contract and is symlinked as `claude` onto `PATH` by the test harness.
- `cmd/scenario-coverage` - the coverage gate behind `scripts/scenario-coverage`.
- `internal/snapshot` - the freeze: a temporary index from `HEAD`, `add -A`, `write-tree`, `commit-tree`, a detached worktree of the result, and the seal that makes that worktree unwritable.
- `internal/runstore` - run directories, `run.json` written by rename, run selection, rounds, and the retention sweep.
- `internal/runstate` - the six states and which of them are terminal.
- `internal/findings` - the findings contract, the schema handed to the reviewer, the classification and the verdict rule. It imports nothing outward.
- `internal/prompt` - the delivered prompt, assembled from named sections in one fixed order.
- `internal/testquality` - the test-quality rule, as one block, shared by every role that needs it.
- `internal/agent` - the only package that starts a model process or parses its output.
- `internal/reviewer` - beginning a run, supervising the round, cancelling.
- `internal/report` - what `status` and `findings` print, for a person and as JSON.
- `internal/exit` - the exit-code table.
- `internal/config` - the one optional configuration file under the user's configuration directory.
- `internal/runlog` - the JSONL event log, its byte ceiling and its in-place rotation.
- `internal/convention` - the repository's own conventions as executable checks: no commentary, the test-quality rule, the message registry.
- `internal/coverage` - scenario coverage: what `openspec/specs` and the changes under way promise against what the tests name. An archived change is skipped, because its specification was promoted.
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
```

`make lint` fails on anything `gofmt` would rewrite, so formatting is never a separate thing to remember.
`make scenario-coverage` fails when a scenario belonging to a change somebody has started is named by no test, and prints where each gap is.
`make test` holds the same gate, so a scenario nobody tests fails the ordinary suite as well.

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
- **Nothing can write in the snapshot worktree, because the filesystem refuses.** The worktree is sealed after checkout by clearing every write bit in it, the seal is proved by one write before a reviewer is spawned, and a filesystem that allows that write is a refusal rather than a review socrates cannot stand behind. The seal is lifted when the round ends, so a worktree kept for inspection is an ordinary directory.
- **The reviewer loads user settings only.** The reviewed repository's own agent instructions are input authored by the thing under review, and a reviewer that reads them can be seated by them.
- **The reviewer loads no external tool server.** A tool server is a second source of tools that the vendor's tool restriction does not cover, and one wrote a directory into the snapshot worktree on socrates' first live run. Where the CLI refuses to ignore the machine's own configuration, the review runs anyway and the run records that the servers were loaded, because the seal is what holds the reviewer read-only and no configuration of the vendor's CLI can lift it.
- **A review inside a review is refused.** The reviewer's environment carries `SOCRATES_RUN`, and `scrat review` refuses when it is already set.
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
- **Mutate a guard you rely on.** A check nobody has broken on purpose is a check nobody knows works.
