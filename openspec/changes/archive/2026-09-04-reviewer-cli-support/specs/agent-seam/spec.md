## RENAMED Requirements

- FROM: `### Requirement: The Claude implementation speaks the CLI's wire contract`
- TO: `### Requirement: The claude implementation speaks that CLI's wire contract`


## MODIFIED Requirements

### Requirement: Every model call passes through one interface

socrates SHALL reach a coding agent through exactly one internal interface, taking a role, a prompt, a working directory and a JSON schema, and answering structured output, text, the model that answered and what it cost.

No other code SHALL start a model process, construct a vendor's arguments or parse a vendor's output.

The seam SHALL admit more than one vendor, and a vendor SHALL be one file behind that interface: an argument construction, a stream parser, and nothing above the seam that knows which one answered.
It exists so that the test suite can substitute the agent, so that a per-role model choice has one place to live, and so that a person who runs a supported CLI other than the one socrates shipped with can be reviewed by it.

It SHALL NOT become a plugin interface, a registry, or configuration describing an arbitrary CLI's flags.
A vendor is a CLI somebody read the documentation for, measured against `agent-support`'s criteria, wrote a simulator for and recorded the shapes of.
Configuration that lets a person describe a CLI socrates has never seen would state support socrates cannot hold, because every guarantee in this specification is a claim about a particular wire contract.

Where a vendor is stated to constrain its output against the supplied schema, a call that supplies a schema and receives no structured output SHALL fail, and SHALL NOT read the agent's text instead.
`agent-support` states the two tiers on which a findings document arrives, and which tier a vendor is on is stated in advance rather than discovered from what came back.
A vendor that was to constrain the output and did not has failed, and reading its prose would turn a vendor defect into a review of unknown provenance.

Which vendor a run was carried out under is `agent-support`'s, and what each vendor's wire contract is is stated here, once per vendor.

#### Scenario: One place starts the process

- **WHEN** the source is searched for the construction of a process that runs a coding agent
- **THEN** exactly one package does it, and every caller reaches it through the interface

#### Scenario: A schema is honoured or the call fails

- **WHEN** a call supplies a schema to a vendor stated to constrain output against it, and the agent returns no structured output
- **THEN** the call fails naming that structured output was required and absent, and the agent's text is not read instead

#### Scenario: Nothing above the seam knows which vendor answered

- **WHEN** the source above the seam is searched for a vendor's name, flag or event shape
- **THEN** none is found, and what varies between vendors is confined to that vendor's file

### Requirement: The agent is simulated for tests, at the wire rather than in Go

The test suite SHALL substitute the agent with a simulator binary that speaks the same command-line and output-stream contract, placed on `PATH` under the name socrates invokes.
Substituting at the Go interface SHALL NOT be the primary mechanism.
Argument construction, stream parsing, process handling and reaping are where the defects are, and a Go-level fake exercises none of them.

Every supported vendor SHALL have a simulator of its own, speaking that vendor's contract and installed under that vendor's name.
One simulator carrying a switch over vendors SHALL NOT be that mechanism, for the reason a Go-level fake is not.
Argument construction and stream parsing are what differ between vendors, and they are what a simulator exists to exercise.
A shared parser exercises one of them however many vendors are configured.

A scenario SHALL be stated in terms of what happens rather than in terms of a vendor or a tier, and every vendor's simulator SHALL be able to produce every scenario.
A truncated stream, a missing terminal event, a schema violation, a non-zero exit, silence and a surviving child are things any coding agent can do.
On a vendor that constrains no schema, a schema violation is a final message whose document violates the schema and a missing document is a final message carrying none, so one vocabulary covers both tiers.
A scenario that only one simulator can produce is a scenario the other vendors are untested against.

The simulator SHALL be driven by a scenario file, SHALL be able to produce every failure named in this specification on demand, and SHALL append every invocation it received, including arguments, working directory, environment and prompt, to a log a test can assert on.

A scenario SHALL be able to name a recorded answer instead of carrying one, and the simulator SHALL then emit that entry's findings document as its answer, byte for byte, in whatever way its vendor delivers one.
It SHALL NOT adapt a recorded document to the prompt it received, renumber its findings, or supply a field the document lacks.
A simulator that repairs a recorded answer is answering from its author's beliefs again.
A scenario naming an answer the corpus does not hold SHALL be refused, naming it, rather than answered with something else.

The end-to-end suite SHALL run its scenarios against every supported vendor.
Scenario coverage is the release metric, and a coverage claim measured against one vendor's adapter is a claim about that adapter rather than about socrates.

A vendor stated as supported and having no simulator SHALL fail the suite, in the way a release target on one side of the platform statement and not the other already does.

#### Scenario: A test asserts what was asked

- **WHEN** a review runs against the simulator
- **THEN** the invocation log holds the exact prompt, arguments and working directory socrates used

#### Scenario: The simulator produces a failure on demand

- **WHEN** a scenario asks for a truncated stream, a missing result event, a schema violation, a non-zero exit, silence, or a child process that outlives the agent
- **THEN** the simulator produces exactly that, and socrates behaves as this specification requires for it

#### Scenario: A scenario names a recorded answer

- **WHEN** a scenario names a corpus entry and a review runs
- **THEN** the reviewer's answer carries that entry's document, byte for byte, delivered the way that vendor delivers one

#### Scenario: A named answer that does not exist is refused

- **WHEN** a scenario names an entry the corpus does not hold
- **THEN** the simulator refuses, naming the entry, and no review is answered with a substitute

#### Scenario: One scenario, every vendor

- **WHEN** the end-to-end suite runs
- **THEN** each scenario runs against every supported vendor's simulator, and the assertions are the same ones

#### Scenario: A supported vendor with no simulator fails the suite

- **WHEN** a vendor is stated as supported and no simulator speaks its contract
- **THEN** the suite fails, naming that vendor

### Requirement: The simulator is checked against the real tool

Responses from each supported vendor's real CLI SHALL be recorded, and a drift check SHALL compare that vendor's simulator against those recordings.
A separate check SHALL run one real review against each real binary and assert that the arguments are still accepted and that a findings document arrives and satisfies the schema.

Where a vendor is on the tier that extracts the document from the final message, that check SHALL record how many attempts the review took.
The tier's cost is how often a reviewer answers without a usable document, and that is a number to measure rather than a claim to make.

What a real reviewer answered SHALL be recorded too, and not only the shape it answered in.
The shapes say the stream still parses; the answers are the only evidence about what the reviewer actually reports, and the scenarios the suite runs on were written before anybody had looked at one.

Recorded answers SHALL NOT be required from every vendor.
The corpus exists so that the scenarios stop resting on their author's beliefs about what a reviewer says, and that holds once for the tool rather than once per vendor.

Both SHALL be behind build tags, so that the ordinary suite compiles neither and cannot reach the network by accident.
Replaying a recorded answer SHALL NOT be, because it reaches nothing.

A simulator nobody compares to the real thing becomes a specification of its author's beliefs, and a green suite against one is worse than no suite because it manufactures confidence.

#### Scenario: Drift is visible

- **WHEN** a vendor's recorded responses no longer match the shapes that vendor's simulator produces
- **THEN** the drift check fails, naming the vendor and the fields that differ

#### Scenario: Every supported vendor is checked against its own recordings

- **WHEN** the drift check runs
- **THEN** it covers every supported vendor, and a supported vendor with no recordings fails rather than passes silently

#### Scenario: What the extracting tier costs is recorded

- **WHEN** a live check runs against a vendor whose document is extracted from the final message
- **THEN** the number of attempts the review took is recorded

#### Scenario: A live review's answer is kept

- **WHEN** a live review against a real binary returns a findings document
- **THEN** the document is recorded into the corpus with the vendor, model, effort and round it came from, whatever it found

#### Scenario: The ordinary suite reaches no network

- **WHEN** the default test command runs
- **THEN** neither the recording nor the live check is compiled, no test invokes any real reviewer binary, and replaying a recorded answer still works

## ADDED Requirements

### Requirement: How a failure is reported is the seam's rule rather than a vendor's

Where a call failed, what socrates says about it SHALL be decided once at the seam, and a vendor SHALL contribute only which event and which fields carry the agent's account.

A call that failed SHALL carry the account the agent gave of its own failure, where the agent gave one.
The vendor's terminal event is where the CLI says why the work did not happen, and the agent's error stream is where a process complains about its arguments, its environment and its own startup.
The account SHALL therefore be preferred to the error stream, and both SHALL be carried where both are there.

Where the agent said nothing on either channel the failure SHALL say that it said nothing, rather than quoting an empty one.

The reason SHALL be bounded, on the ceiling and the mechanism the error stream already has rather than a second one.

socrates SHALL NOT decide anything from that account.
It is not parsed, not matched and not compared, no state, verdict or exit code depends on a word in it, and the call had already failed before it was read.

These are rules about failures rather than rules about a vendor, and a second vendor restating them is a second place for them to drift.
The one thing that is per-vendor is which event reports the failure and which of its fields is the account, and that is stated in that vendor's requirement.

#### Scenario: Every vendor reports a failure the same way

- **WHEN** any supported vendor's agent fails with an account of its own failure
- **THEN** the failure carries that account, bounded on the shared ceiling, in preference to the error stream, and nothing socrates decides reads a word of it

#### Scenario: A vendor contributes only where the account is

- **WHEN** the source is searched for the preference order, the bounding and the refusal to decide from the account
- **THEN** each is stated once at the seam rather than once per vendor

### Requirement: The codex implementation speaks that CLI's wire contract

The `codex` implementation SHALL invoke `codex exec` with a newline-delimited JSON event stream, a read-only sandbox, and the JSON schema supplied as a file.
It is on the tier where the vendor constrains the output against that schema.

The sandbox SHALL be stated rather than relied on as a default.
`codex exec` is read-only unless told otherwise, and a default can change in a release.
That sandbox is enforced by the operating system rather than by a model choosing not to call a tool, so it is a stronger mechanism than a tool allowlist for the same guarantee.
It does not replace the seal or the comparison after the round.

The schema file SHALL be written under socrates' home, in the round's own directory, and SHALL NOT be written in the repository under review or in the snapshot worktree.
Nothing writes in the repository under review, and the snapshot worktree is sealed.

It SHALL ask the vendor to ignore the machine's own configuration file, which is where that CLI's external tool sources, plugins and hooks are configured, and SHALL supply no configuration of its own beyond the invocation.

It SHALL suppress the repository's own project document, so that the change under review does not instruct the reviewer judging it.
That is the same criterion the `claude` implementation meets by restricting settings to user level, reached by that CLI's own mechanism.
The one file socrates reads out of the repository is unaffected: `review-guidance` specifies it, socrates reads it at the base commit, and it reaches the prompt in its own section as the project's own words.

The role's model and reasoning effort SHALL be supplied on the invocation.

The account of a failure SHALL be taken from the event by which that CLI reports a failed turn or an error, and what socrates then does with it is the seam's rule rather than this one's.

That CLI reports no inventory of the tools its session holds.
socrates SHALL therefore record, on every round that completes under it, that the session did not say what it held, on the rule this specification already states for a session that never said.
It SHALL NOT record such a round as a session that held only what socrates asked for, and SHALL NOT fail the round for the absence.
An exemption stating that the sandbox makes the observation unnecessary SHALL NOT be added.
It would record a guarantee nobody measured, which is what reading the tool list rather than trusting the mechanism exists to avoid.

#### Scenario: The invocation carries the schema as a file under socrates' home

- **WHEN** a request supplying a schema is made
- **THEN** the recorded invocation names a schema file, that file holds the exact schema, and it is under socrates' home rather than in the repository under review or the snapshot worktree

#### Scenario: The sandbox is stated

- **WHEN** any reviewer is invoked
- **THEN** the recorded invocation states the read-only sandbox rather than relying on it being the default

#### Scenario: Nothing asks for the machine's own customisations

- **WHEN** any reviewer is invoked
- **THEN** the recorded invocation asks the vendor to ignore the machine's own configuration file, and supplies no configuration of its own beyond the invocation

#### Scenario: The repository does not instruct its own reviewer

- **WHEN** a reviewer is invoked against a repository that carries its own agent instructions
- **THEN** the recorded invocation suppresses that repository's project document, and the reviewer's prompt carries none of those instructions

#### Scenario: The one file socrates reads still reaches the prompt

- **WHEN** a reviewer is invoked against a repository carrying both its own agent instructions and the guidance file socrates reads
- **THEN** the prompt carries none of the repository's agent instructions, and the guidance reaches the prompt in its own section

#### Scenario: A session that cannot say what it held

- **WHEN** a round completes under this vendor
- **THEN** the run records that the session did not say what it held, the round is not failed for it, and nothing records the session as holding only what socrates asked for

#### Scenario: A tool that writes cannot reach the change under review

- **WHEN** a reviewer runs under this vendor and something attempts to write into its working directory
- **THEN** the write is refused, the review completes, and the run's verdict is whatever the findings say

### Requirement: The opencode implementation speaks that CLI's wire contract

The `opencode` implementation SHALL invoke `opencode run` with that CLI's JSON event stream and the role's model.
It is on the tier where the document is extracted from the reviewer's final message, because that CLI constrains no schema.

It SHALL keep the reviewed repository's own instruction discovery out of the session, on that CLI's own mechanism for doing so.

The role's model SHALL be supplied on the invocation.
Where that CLI offers no way to state a reasoning effort, the role SHALL resolve to its model with no effort passed, which the routing already permits.
A vendor that cannot carry an effort SHALL NOT cause a role to resolve to a different model.

The account of a failure SHALL be taken from that CLI's terminal event, and what socrates then does with it is the seam's rule rather than this one's.

That CLI reports no inventory of the tools its session holds, and the rule stated for the vendor above applies here unchanged.

#### Scenario: The repository does not instruct its own reviewer

- **WHEN** a reviewer is invoked against a repository that carries its own agent instructions
- **THEN** the recorded invocation keeps that repository's instruction discovery out of the session, and the reviewer's prompt carries none of those instructions

#### Scenario: The document is extracted and validated

- **WHEN** a review completes under this vendor and the reviewer's final message ends with a findings document
- **THEN** the document is extracted, validated against the schema, and the round proceeds as it would under a vendor that constrained it

#### Scenario: A reviewer that answered without a document

- **WHEN** a review completes under this vendor and the reviewer's final message carries no findings document
- **THEN** the call fails naming the absent document, and the round is `failed`

#### Scenario: A role with no effort to state

- **WHEN** a role resolves to a model and an effort and this vendor offers no way to state one
- **THEN** the recorded invocation carries that model, carries no effort, and the role's model is unchanged

#### Scenario: A session that cannot say what it held

- **WHEN** a round completes under this vendor
- **THEN** the run records that the session did not say what it held, the round is not failed for it, and nothing records the session as holding only what socrates asked for
