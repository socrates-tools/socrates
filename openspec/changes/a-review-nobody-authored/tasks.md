**Nothing here writes in the repository under review.** No task commits, stages, stashes, pushes, creates a ref, or leaves a file behind. The one thing that reaches a repository is the snapshot, and the test that `git status` is byte-identical afterwards is what holds it.

**Test first.** Every task below is written so the test can be written before the code: the specification scenario it covers is named on the task, and the task is not done until a test named for that scenario fails for the right reason and then passes.

**`make lint test` after every task, not at the end of a group.** The ordering exists to make that possible.

## 1. The module and the shape of a run

- [x] 1.1 Create the Go module, `cmd/scrat` with cobra and nothing in it, and the `Makefile` targets `build`, `test`, `lint`, `fmt`. `lint` fails on anything `gofmt` would rewrite, so formatting is never a separate thing to remember
- [x] 1.2 Add the commentary check: a test that reads every source file in the trees this repository writes in and fails on any comment that is not a directive a tool parses or a `// scenario:` marker. It holds no comments of its own, which is the only honest way to state it
- [x] 1.3 Add the run identity: the directory name, the repository key derived from the absolute repository root, and the sort. Covers `Runs sort by when they started`, `Two repositories do not collide`, `The same path always gets the same key`
- [x] 1.4 Add `run.json`, written by temporary file and rename. Covers `A reader never sees half a state`, `A crash mid-write leaves the previous state`
- [x] 1.5 Add the six states and the terminality answer. Covers `A terminal run does not move`, `A failed run reports no findings as its verdict`

## 2. The snapshot

- [x] 2.1 Add the snapshot: temporary index from `HEAD`, `add -A`, `write-tree`, `commit-tree`. Covers `An untracked file is in the snapshot`, `A modified tracked file is in the snapshot`, `A staged change and an unstaged change are both in the snapshot`, `A deleted file is deleted in the snapshot`, `An ignored file is not in the snapshot`
- [x] 2.2 Add the untouched-repository test before anything else in this group is trusted. Covers `Git status is unchanged`, `The index is unchanged`, `No ref names the snapshot`, `The snapshot's parent is the current head`. This is the test that fails if the temporary index leaks, and a leak is silent otherwise
- [x] 2.3 Add the detached worktree, its placement under the run directory, and its removal. Covers `The worktree is the snapshot`, `Editing while the reviewer reads`, `The worktree is cleaned up`
- [x] 2.4 Add the refusals. Covers `Not a repository`, `A repository with no commits`, `A bare repository`
- [x] 2.5 Add the cleanup on failure. Covers `Worktree creation fails`, `The temporary index never joins the repository`, `A snapshot in progress leaves nothing behind when it fails`
- [x] 2.6 Add the nothing-to-review path. Covers `A change with nothing to review`

## 3. The findings contract, as a pure core

- [x] 3.1 Add the findings types and the schema document handed to the agent. Covers `The reviewer is given the schema`, `A finding names a file and a line`, `Identifiers are unique within a run`
- [x] 3.2 Add the single classification helper and test it against the whole input space, not a sample. Covers `A missing action defaults to needing a person`, `An empty action defaults to needing a person`, `An unrecognised action defaults to needing a person`, `Every finding is in exactly one bucket`, `A document mixing every case`
- [x] 3.3 Add the test that only one function decides a finding's effective action and both consumers call it. Covers `The two consumers agree by construction`. It is the requirement that cannot be re-derived from behaviour once the second copy exists
- [x] 3.4 Add the verdict rule as a pure function, and the test that its package imports nothing outward. Covers `No findings at all`, `Only no-ops and informational findings`, `Anything blocking`, `The rule reaches nothing outward`
- [x] 3.5 Add severity and action orthogonality. Covers `An informational finding does not block`, `A warning blocks`, `Action does not change whether a run blocks`, `An informational finding can still be actionable`
- [x] 3.6 Add the risk level as recorded and not deciding. Covers `A high risk level with no blocking findings`, `The rationale is one sentence and is kept`

## 4. The simulator, before the agent that needs it

- [x] 4.1 Build `cmd/claudesim`: argument acceptance, the streaming JSON event format, the terminal result event, and the exit code, matching the real CLI's contract
- [x] 4.2 Add the scenario file it is driven by, and the invocation log it appends to. Covers `A test asserts what was asked`
- [x] 4.3 Add every failure the specification names as a scenario the simulator can produce on demand. Covers `The simulator produces a failure on demand`
- [x] 4.4 Add the PATH substitution the end-to-end harness uses: a temporary directory holding the simulator symlinked under the name socrates invokes, prepended to `PATH`

## 5. The agent seam

- [x] 5.1 Add the interface, the role type, and the per-role model and effort resolution. Covers `The reviewer's model is configurable on its own`, `An empty configuration still runs`
- [x] 5.2 Add the boundary test that exactly one package constructs a model process. Covers `One place starts the process`
- [x] 5.3 Add the Claude implementation: print mode, streaming JSON, the schema on the command line, user-level settings only, and non-writing tools. Covers `The invocation carries the schema`, `Project settings are not loaded`
- [x] 5.4 Add the stream parser and its failures. Covers `A stream that ends mid-event`, `A stream with no result event`, `Structured output that violates the schema`, `A non-zero exit after a complete stream`, `A schema is honoured or the call fails`
- [x] 5.5 Add the deadline and the process-group kill on expiry, tested with `testing/synctest`. Covers `A silent agent is killed`, `A slow but live agent is not killed early`
- [x] 5.6 Record the model and the token counts in the event log. Covers `What the agent cost is recorded`

## 6. The reviewer run

- [ ] 6.1 Add the detached spawn, the recorded process identity, and the released standard streams. Covers `The command returns while the review runs`, `The reviewer survives its caller`
- [ ] 6.2 Add the process group and reap it on every exit path. Covers `A child that outlives the agent is killed`, `Reaping happens on error too`, `Reaping happens on cancellation`, `The worktree can be removed after the round`. Write the clean-exit test first: reaping only on cancellation is the version that looks correct and is not
- [ ] 6.3 Add the re-entrancy marker and the refusal. Covers `The reviewer's environment is marked`, `A review inside a review is refused`
- [ ] 6.4 Add the worktree comparison after the round and the failure it causes. Covers `A reviewer that wrote fails the round`, `A reviewer that created a file fails the round`, `A clean reviewer passes the check`, `A dirtied worktree is kept`
- [ ] 6.5 Add the positive-evidence rule. Covers `A missing findings document is a failure`, `A dead reviewer is a failure`, `An empty findings array with evidence is clean`, `An unparseable findings document is a failure`
- [ ] 6.6 Stream the reviewer's output to disk as it arrives. Covers `Output is visible while the review is running`, `A failed round keeps what it got`
- [ ] 6.7 Add the dead-reviewer answer on read. Covers `A killed reviewer stops the run`, `A live reviewer is still working`

## 7. The prompt

- [ ] 7.1 Add the named sections and the fixed assembly order, with an absent section rather than an empty one. Covers `Order is stable`, `An empty section is absent`
- [ ] 7.2 Record the delivered prompt in the round directory and make it what tests read. Covers `The prompt is on disk`, `Assertions are made on the delivered prompt`
- [ ] 7.3 Add the context section, the intent section, and the absence of any channel for the author's rationale. Covers `Intent reaches the reviewer`, `No intent is stated`, `The author's rationale has no channel`
- [ ] 7.4 Add the diff section and the ceiling. Covers `A small change carries its diff`, `A large change is read`, `The ceiling is one setting`
- [ ] 7.5 Add the reviewer rules and the action definitions. Covers `Every rule reaches the reviewer`, `The action definitions are the ones the classifier implements`
- [ ] 7.6 Add the test-quality rule as one named block, and the test that socrates' own suite obeys it. Covers `One block, many roles`, `socrates holds itself to it`

## 8. The command surface

- [ ] 8.1 Add `scrat review`, `--wait`, and the stated intent. Covers `The default does not block`, `Waiting returns the verdict`, `Starting is not a verdict`
- [ ] 8.2 Add `scrat status` and `--json`. Covers `A running review is a few lines`, `A terminal run reports its verdict`, `Machine-readable on request`
- [ ] 8.3 Add `scrat findings` and `--json`, printing the effective action rather than the raw field. Covers `Findings are readable`, `The effective action is resolved for the reader`, `A clean run has findings to read too`, `A failed run has no findings to give`
- [ ] 8.4 Add the run selection default. Covers `The bare command finds this repository's run`, `A repository with no runs`, `An identifier wins over the default`
- [ ] 8.5 Add the exit codes as a table one function owns. Covers `Each verdict has its own code`, `A usage error is not a verdict`, `Failure is not escalation`
- [ ] 8.6 Add `scrat cancel`. Covers `Cancelling stops everything`, `Cancelling a terminal run changes nothing`
- [ ] 8.7 Add the configuration file and its refusal. Covers `No configuration file`, `The repository is never written to`, `An unreadable configuration file is a refusal`
- [ ] 8.8 Add the command surface test. Covers `The command surface is what it says`
- [ ] 8.9 Move every user-facing string into one registry per package and add the scan. Covers `A refusal names the next action`, `No em dash reaches a person`

## 9. The event log

- [ ] 9.1 Add the JSONL event log and its writer. Covers `Events are one JSON object per line`, `A started run is on disk immediately`, `Round output is under the round`
- [ ] 9.2 Add the byte ceiling and in-place truncating rotation. Covers `A chatty reviewer does not fill the disk`, `The writer keeps writing across a rotation`
- [ ] 9.3 Add the retention sweep on run start. Covers `An old run is swept`, `A live run is never swept`

## 10. End to end

- [ ] 10.1 Build the harness: real binary, real temporary repository, real git, real spawning, the simulator on `PATH`, `synctest` for time. Nothing else substituted
- [ ] 10.2 Drive the clean journey: implement, review, verdict `clean`, exit `0`, repository byte-identical
- [ ] 10.3 Drive the escalated journey: findings returned, verdict `escalated`, exit `1`, findings readable and classified
- [ ] 10.4 Drive each failure journey the simulator can produce, and assert that none of them reports `clean`
- [ ] 10.5 Drive the writing-reviewer journey and assert the run fails and the worktree is kept
- [ ] 10.6 Drive the recursive journey: a review started from inside a review is refused
- [ ] 10.7 Drive the concurrent-editing journey: the person edits while the reviewer reads, and the findings still anchor to the snapshot
- [ ] 10.8 Add the CLI behaviour tests as `testscript` txtar scripts, covering the exit codes, the run selection default and the refusals

## 11. Fidelity to the real tool

- [ ] 11.1 Add the recording, behind a build tag, of responses from the real `claude` CLI, and `make fixtures` to run it
- [ ] 11.2 Add the drift check comparing the simulator's shapes against the recordings, and `make fixtures-diff`. Covers `Drift is visible`
- [ ] 11.3 Add the live smoke check that runs one real review against the real binary and asserts the arguments are accepted and the schema honoured, behind its own tag, and `make live-smoke`
- [ ] 11.4 Add the test that the ordinary suite compiles none of it. Covers `The ordinary suite reaches no network`
- [ ] 11.5 Add the scheduled workflow that runs the drift and live checks weekly rather than per commit, because a failure means the CLI changed rather than that this change broke something

## 12. Coverage and documentation

- [ ] 12.1 Add `scripts/scenario-coverage`: list every scenario in `openspec/specs/` and `openspec/changes/` that no test names, exiting non-zero on a gap
- [ ] 12.2 Run it and close every gap this change introduced. A scenario with no test is a promise socrates has not made
- [ ] 12.3 Write `AGENTS.md` with the layout, the verification sequence and the load-bearing invariants, and symlink `CLAUDE.md` to it
- [ ] 12.4 Write `README.md`: what socrates is, what it deliberately is not, and the first run
- [ ] 12.5 Start `docs/gotchas.md` with what this change had to re-derive: the `git stash create` untracked-file trap, `GIT_INDEX_FILE` needing a path that does not exist, and the `synctest` bubble rules

## 13. Verification

- [ ] 13.1 `make fmt && make lint test` clean
- [ ] 13.2 `make scenario-coverage` naming no scenario this change adds
- [ ] 13.3 A live run: review a real change in a real repository with the real `claude` CLI, and record what it found in this change's `retrospective.md`. A live run against a real agent routinely finds defects a green suite against a simulator cannot see, and the ratio is the signal for how novel the surface was
