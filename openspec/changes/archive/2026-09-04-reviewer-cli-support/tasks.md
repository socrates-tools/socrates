**The suite must be green after every task, not only at the end of a group.**
The ordering below exists for that reason: the supported set is stated and read before a second vendor is added to it, and the parsed tier is built and tested against the existing vendor's simulator before a vendor that needs it arrives.

**Group 1 comes first because it decides which vendors this change ships.**
Each recording establishes one criterion for one vendor, and a negative answer takes that vendor out of the supported set and into the deferred issue rather than into the implementation.
The recordings are cheap and they are worth having before the adapters, the simulators, the matrix and the fidelity suite are built on them.

**Every task here records what a live run did. None returns what it should have done.**

## 1. Which vendors clear the criteria

- [x] 1.1 Add a recording behind the existing build tag that invokes the real `codex exec` against a temporary repository whose `AGENTS.md` carries a distinctive instruction, with `-c project_doc_max_bytes=0`, and record whether the answer shows any sign of having read it. Record the same invocation without the override as the control, so the recording says the override did something rather than that the model ignored the file
- [x] 1.2 Add the same recording for the real `opencode run` under `OPENCODE_DISABLE_PROJECT_CONFIG=1`, with the same control. That variable is documented as skipping project discovery and not the global file, so record what a repository's own instructions do under it
- [x] 1.3 Add a recording that invokes the real `opencode run --format json` and record whether the terminal event carrying the account of a failure arrives, against anomalyco/opencode#26855 and #31435. Record it for a failing run as well as a succeeding one
- [x] 1.4 Add a recording that invokes the real `codex exec` with a schema file and record whether structured output came back conforming to it, against openai/codex#15451
- [x] 1.5 Add a recording that invokes the real `codex exec` with `--ignore-user-config` on a machine holding a configured external tool source, and record whether the run authenticates and whether that source reached the session
- [x] 1.6 Add a recording of a failed turn under each of the two CLIs, and record which event and which field carry the account of the failure
- [x] 1.7 Add a recording that establishes whether `opencode run` can state a reasoning effort, and record what the model variant covers
- [x] 1.8 Write the recordings up in `design.md` under a "What was measured" heading, replacing the inferences it states. Move any vendor that failed a criterion into the deferred issue and out of this change, and revise `proposal.md` to match

## 2. The bar, stated once

- [x] 2.1 Add the vendor statement to `internal/agent`: each CLI beside the facts that make it qualify, `Qualifies()` as the conjunction of the three criteria, `Tier` and `ReportsItsToolInventory` stated beside them rather than inside, and `Supported()` as the filter, in the shape `internal/platform` uses. State the vendors group 1 confirmed, and state `cursor-agent`, `copilot` and `pi` as not qualifying with the criterion each fails
- [x] 2.2 Add a test that exactly one statement of the supported set exists in the source and that everything naming a supported CLI reads it. Covers `The supported set is stated once`
- [x] 2.3 Add a test that every unsupported entry names the criterion it fails. Covers `An unsupported CLI is stated with the reason it is unsupported`
- [x] 2.4 Add a test that an entry on the extracting tier still qualifies, and mutate `Qualifies()` to verify the test fails if the tier is folded into the conjunction. Covers `A CLI that constrains no schema still qualifies`
- [x] 2.5 Add a test that flipping an entry's isolation criterion to false takes it out of `Supported()` with nothing else edited, and mutate the conjunction to verify the test fails when the criterion stops being load-bearing. Covers `Support follows from the criteria rather than from a list`

## 3. Naming the vendor, and refusing one socrates does not support

- [x] 3.1 Add the reviewer CLI field to `internal/config`, defaulting to `claude`, and add a test that an absent configuration file and an absent field both resolve to it. Covers `An absent field runs the CLI socrates shipped with`
- [x] 3.2 Refuse a configured CLI that is not in `Supported()`, naming it and the guarantee it does not offer, at the point an unsupported platform is refused. Covers `An unsupported CLI is refused before anything happens`
- [x] 3.3 Add an end-to-end test that such a refusal takes no snapshot and starts no process, asserting against the simulator's invocation log and against `git status` being byte-identical
- [x] 3.4 Add an end-to-end test that a command which neither spawns nor reaps is not refused under an unsupported configured CLI. Covers `A finished run stays readable under any configuration`
- [x] 3.5 Verify `make lint test` passes with `claude` still the only implementation

## 4. The two tiers

- [x] 4.1 Add extraction to `internal/findings`: the last fenced JSON candidate in a final message, reaching the existing schema validation rather than a second one. Add unit tests for no candidate, a candidate that is not JSON, a candidate violating the schema, and several candidates where the last one is the answer. Covers `Reasoning before the answer is not mistaken for the answer`
- [x] 4.2 Add the tier to the vendor statement and route the seam through it: structured output on the enforced tier, extraction on the parsed tier, one validation for both. Covers `The enforced tier delivers the vendor's structured output`, `The parsed tier delivers the document from the final message`
- [x] 4.3 Add a test that a vendor stated to constrain the schema and returning no structured output fails without reading the text. Mutate the tier check to verify the test fails when the fallback becomes reachable. Covers `A vendor that was to constrain the schema and did not fails`
- [x] 4.4 Add the delivery instruction to the reporting section of the prompt for the parsed tier, generated from `internal/findings` as that section already is, and add a test asserting the delivered prompt
- [x] 4.5 Add a test that nothing in the source or in any delivered prompt asks the authoring session to produce, reshape or complete a findings document. Covers `The author never shapes the document`
- [x] 4.6 Add a parsed-tier scenario to the claude simulator and run the existing verdict and dialogue scenarios through it, so the tier is exercised before a vendor needs it. Covers `An unparseable document fails on either tier`

## 5. The failure rules move to the seam

- [x] 5.1 Hoist the preference order, the bounding and the refusal to decide from the account out of the claude implementation into the seam, leaving the claude file contributing only which event and which fields carry the account. No behaviour changes
- [x] 5.2 Add a test that each of those rules is stated once rather than once per vendor. Covers `A vendor contributes only where the account is`
- [x] 5.3 Verify every existing failure scenario in `internal/agent` still passes unchanged, and mutate the hoisted preference order to verify they fail for the right reason

## 6. The codex implementation

- [x] 6.1 Add the codex argument construction: `exec`, the JSON event stream, the read-only sandbox stated, the schema file, `--ignore-user-config`, the project document suppressed, and the role's model and effort. Add tests asserting each against the recorded invocation. Covers `The sandbox is stated`, `Nothing asks for the machine's own customisations`
- [x] 6.2 Write the schema to the round's directory under socrates' home, and add a test that the file holds the exact schema and is under socrates' home. Covers `The invocation carries the schema as a file under socrates' home`
- [x] 6.3 Add an end-to-end test that the schema file is written neither in the repository under review nor in the snapshot worktree, asserting `git status` is byte-identical and the worktree still matches its commit
- [x] 6.4 Add the codex stream parser: the terminal turn event, the error event, the usage and the model. Add tests for a truncated stream, a missing terminal event, a schema violation, a non-zero exit and silence, reusing the seam's failure rules
- [x] 6.5 Record, on every completed codex round, that the session did not say what it held, and add a test that no codex round reads as a session holding only what socrates asked for. Covers `A session that cannot say what it held`
- [x] 6.6 Add an end-to-end test that a repository carrying its own agent instructions does not reach the codex reviewer's prompt, and that the guidance file still does. Covers `The repository does not instruct its own reviewer`, `The one file socrates reads still reaches the prompt`
- [x] 6.7 Add an end-to-end test that a write attempted in the codex reviewer's working directory is refused and the review still completes. Covers `A tool that writes cannot reach the change under review`

## 7. The opencode implementation

- [x] 7.1 Add the opencode argument construction: `run`, the JSON event stream, the role's model, and the environment that keeps the repository's instruction discovery out of the session. Add tests asserting each against the recorded invocation. Covers `The repository does not instruct its own reviewer`
- [x] 7.2 Add the opencode stream parser: the terminal event, the account of a failure, the usage and the model, reusing the seam's failure rules
- [x] 7.3 Place opencode on the parsed tier and add end-to-end tests for a reviewer that ends with a document and one that does not. Covers `The document is extracted and validated`, `A reviewer that answered without a document`
- [x] 7.4 Resolve a role to its model with no effort passed where the CLI states none, and add a test that the model is unchanged by the absent effort. Covers `A role with no effort to state`
- [x] 7.5 Record that the session did not say what it held on every completed opencode round, and add a test. Covers `A session that cannot say what it held`
- [x] 7.6 Select the implementation from the configured vendor at the one construction site in `cmd/socrates/supervise.go`, and add a test that nothing above the seam names a vendor, a flag or an event shape. Covers `Nothing above the seam knows which vendor answered`

## 8. The simulators

- [x] 8.1 Add `internal/codexsim` and `cmd/codexsim` speaking codex's command-line and event-stream contract, reading the same scenario files and appending to the same invocation log
- [x] 8.2 Add `internal/opencodesim` and `cmd/opencodesim` on the same terms, delivering a document in its final message rather than in a structured field
- [x] 8.3 Make the scenario vocabulary vendor-neutral and tier-neutral, and verify every scenario is producible by every simulator. Covers `The simulator produces a failure on demand`
- [x] 8.4 Teach `internal/testing/simulator` to install the simulator for the vendor under test, under that vendor's name
- [x] 8.5 Add corpus replay to both new simulators, emitting a recorded entry's document byte for byte in the way that vendor delivers one, and add a test that a scenario naming an absent entry is refused. Covers `A scenario names a recorded answer`, `A named answer that does not exist is refused`
- [x] 8.6 Add a test that a vendor stated as supported with no simulator fails the suite, and mutate the vendor statement to verify it fails. Covers `A supported vendor with no simulator fails the suite`

## 9. Every scenario, every vendor

- [x] 9.1 Parameterise the end-to-end suite over `agent.Supported()`, running the same scenarios and the same assertions against each vendor's simulator. Covers `One scenario, every vendor`
- [x] 9.2 Verify `make scenario-coverage` and `make scenario-drop` still hold, and that coverage is now claimed against every supported vendor
- [x] 9.3 Add an end-to-end test per vendor that a process left behind by a reviewer is reaped by process group, since reaping is per-vendor process handling rather than shared code
- [x] 9.4 Verify `make verify` passes and record what the matrix costs in wall time. `go test -race -count=1 -timeout 30m ./...` takes 213s on this machine with three vendors in the matrix. What it took before the matrix was not measured, so the multiple the design accepted is not something this number settles; what it settles is that the suite still finishes well inside its timeout

## 10. Fidelity, per vendor

- [x] 10.1 Move the recorded shapes under a per-vendor directory, keeping the claude recordings byte-identical, and make the drift check read `Supported()`
- [x] 10.2 Promote the group 1 recordings into each vendor's recordings, and add a test that a supported vendor with no recordings fails the drift check rather than passing silently. Covers `Every supported vendor is checked against its own recordings`, `Drift is visible`
- [x] 10.3 Add a live smoke check per vendor behind the build tag, running one real review against the real binary
- [x] 10.4 Record the number of attempts a live review took on the parsed tier. Covers `What the extracting tier costs is recorded`
- [x] 10.5 Add a test that the ordinary suite compiles none of them, reaches no network, invokes no real reviewer binary, and still replays a recorded answer. Covers `The ordinary suite reaches no network`
- [x] 10.6 Verify nothing in continuous integration reaches any real CLI, by target and by build tag
- [x] 10.7 Record the vendor beside the model, effort and round on any live answer kept in the corpus. Covers `A live review's answer is kept`

## 11. What the run says

- [x] 11.1 Record the reviewer CLI and its tier on the run beside the model that answered, and add a test. Covers `The vendor and tier reach the run's record`
- [x] 11.2 Read a run recorded without the field as `claude`, and add a test against a fixture written before the field existed. Covers `A run from before the field is read as the original CLI`
- [x] 11.3 Add the vendor to what `status` and `findings` print, and verify the presentation and user-facing string checks still pass

## 12. What a person reads

- [x] 12.1 Name every supported CLI in `README.md`, taking the set from the statement, and add a convention check that fails while the two disagree. Covers `A document that names fewer CLIs than socrates supports fails`
- [x] 12.2 State which tier each supported CLI is on, and what that means for a person choosing between them
- [x] 12.3 Say in `README.md` which capability each unsupported CLI lacks, and link the deferred issue. Covers `A reader learns why their CLI is not supported`
- [x] 12.4 Update the `agora` and telemetry statements so that "the reviewer CLI it drives sends the change to its vendor" reads correctly with more than one vendor, and verify both halves of the reach are still stated
- [x] 12.5 Name the configuration field in `README.md`, and verify the em dash check and the user-facing string check pass
- [x] 12.6 Update `internal/skill` where it names the reviewer CLI, regenerate the committed copy with `make skill`, and verify the generated-file check passes
- [x] 12.7 Update `AGENTS.md`'s layout, invariants and verification sections, and `openspec/config.yaml` where it states one shipped implementation and one simulator
- [x] 12.8 Add the release note to `CHANGELOG.md`, by hand, in terms of what it means to somebody who has the previous release
- [x] 12.9 Run `make fmt` and `make verify`
