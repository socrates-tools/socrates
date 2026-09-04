# agent-seam Specification

## Purpose
TBD - created by archiving change a-review-nobody-authored. Update Purpose after archive.

## Requirements

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

### Requirement: A role resolves to its own model and effort

Each role SHALL resolve independently to a model identifier and a reasoning effort, from configuration, with a default.
The reviewer wants the strongest model available and nothing else in socrates does, so a single global model setting is a decision that would have to be unmade.

An unknown role SHALL be a programming error caught at construction, not a silent fallback to a default at the call site.

#### Scenario: The reviewer's model is configurable on its own

- **WHEN** configuration names a model for the reviewer role and nothing else
- **THEN** the reviewer is invoked with that model and every other role uses its own default

#### Scenario: An empty configuration still runs

- **WHEN** no configuration file exists
- **THEN** every role resolves to its default and the run proceeds

### Requirement: An agent that never answers does not hold a run forever

Every call SHALL be bounded by a deadline.
On expiry socrates SHALL terminate the agent's process group and fail the call naming the deadline.

A model process that accepts a prompt and never writes is indistinguishable from one thinking hard, and without a deadline the difference is a run that is `reviewing` until the machine is rebooted.

#### Scenario: A silent agent is killed

- **WHEN** the agent writes nothing and does not exit before the deadline
- **THEN** its process group is terminated, the call fails naming the deadline, and no process from it survives

#### Scenario: A slow but live agent is not killed early

- **WHEN** the agent emits events steadily and completes before the deadline
- **THEN** the call succeeds

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

### Requirement: The claude implementation speaks that CLI's wire contract

The implementation shipped SHALL invoke the `claude` binary in print mode with a streaming JSON output format, supply the JSON schema on the command line, and take the structured output from the terminal `result` event of the stream.

It SHALL load user-level settings only, and SHALL NOT load the reviewed repository's own project or local settings, because those are authored by the thing under review.
What is refused there is what the vendor would load on its own account: the repository's project and local settings, and the agent instructions, memory and skills that arrive with them.
They are refused because nothing bounds them, nothing anchors them to a version the change under review did not write, and they reach the session as instructions rather than as something quoted to it.
A file socrates reads itself is a different thing.
`review-guidance` specifies the one there is: read at the base commit, bounded, neutralised, delivered inside a section as the project's own words, and deciding nothing.
The distinction is which side reads it, on what terms, and from which commit, rather than whether the text came from the repository at all.

A call that failed SHALL carry the account the agent gave of its own failure, where the agent gave one.
That terminal event is where the CLI says why the work did not happen, and the agent's error stream is where a process complains about its arguments, its environment and its own startup.
The account SHALL therefore be preferred to the error stream, and both SHALL be carried where both are there.
socrates read the exit status first and returned on it.
A live review that died on the provider being overloaded therefore reported an exit status and an empty error stream.
The message naming the provider sat parsed and unread in the same stream socrates had already finished reading.
A reason answers one question, which is whether to run the review again or to fix something, and an exit status answers it either way.

The subtype beside that account SHALL NOT be the reason where the account is there.
On the same live stream the subtype was `success`, on the event whose error flag was true and whose message named the provider, so of that event's three fields socrates named the only one that said nothing.
It SHALL remain what is said where the event reports an error and carries no account, because a failure has to be reported as something.

Where the agent said nothing on either channel the failure SHALL say that it said nothing, rather than quoting an empty one.

The reason SHALL be bounded, on the ceiling and the mechanism the error stream already has rather than a second one.
It reaches a run's record and the context of a session polling `socrates status`, and the agent's text on some failures is a whole turn.

socrates SHALL NOT decide anything from that account.
It is not parsed, not matched and not compared, no state, verdict or exit code depends on a word in it, and the call had already failed before it was read.
Deriving a decision from a vendor's prose is refused below for the reason it is refused: quoting prose to a person fails by being less useful, and branching on it fails by failing every review.

The reviewer SHALL hold the tools socrates asked for and no others.
It asks on the vendor's own mechanism for tools that do not write, and that mechanism covers the built-in tools the model may call and nothing else.
A tool from any other source is therefore not covered by it: an external tool server contributes tools, so does a plugin, and so does a built-in the vendor adds in a release nobody read.
A hook contributes no tool and is worse than one for it, because it runs a command of its own on an event rather than waiting to be called.
A reviewer holding such a tool is a reviewer that can write, both to the snapshot worktree and outside it.
socrates' first live run found this when the machine's own configured tool server wrote a directory into the snapshot worktree.
A later run found the reviewer's session holding script execution, file upload and browser control.

The tools socrates asks for SHALL be those that read the change and those that check a claim against the world, and SHALL NOT be those that write.
A reviewer that cannot look up what a library actually does writes a finding about its own recollection, and patching is the author's half of the asymmetry this design rests on.
The tools that write are also the ones that reach outside the snapshot worktree, where nothing socrates holds bounds anything.

It SHALL ask the vendor to disable the customisations the machine brings on its own account, on the one mechanism that disables all of those at once and leaves auth, model selection and the built-in tools working.
It SHALL supply no settings of its own, because that mechanism makes an allowlist for one source of tools pointless and covers the sources such an allowlist said nothing about.
It SHALL NOT ask the vendor to replace the machine's own configuration of any source.
That request is refused at startup wherever a managed configuration is deployed, because it asks to replace a managed set, and it is the only reason socrates ever needed a second attempt.
One mechanism that is never refused therefore replaces two of which one always was.

What that mechanism does not reach is a hook an administrator's policy configures, which the vendor keeps in force by design, along with a policy-configured status line and policy-configured file-suggestion commands.
Only the policy-configured ones survive it, and a policy-configured tool source is not among them.
Those do not load, which socrates measured on such a machine, where a plugin's own tool left the session's list at the same time.
The vendor documents managed plugins, managed skills and managed memory going the same way.
socrates holds no measurement distinguishing a managed one from a person's own, because the session's tool list does not say what contributed a tool.
A policy hook contributes no tool, so the session's tool list is not what would catch it.
What does catch it inside the snapshot worktree is the seal and the comparison, because its command is a child of the session and the session runs there.
Nothing socrates holds bounds where that command reaches outside the worktree.
That residual is accepted and stated rather than covered, for the same reason the review runs at all on such a machine: a policy set by somebody entitled to set it should still get reviews.

The mechanism SHALL NOT be trusted to have worked.
socrates SHALL take the session's tool list from the stream's opening event, SHALL treat every tool in it that socrates did not ask for as a tool it did not ask for whatever put it there, and SHALL record those.
The set socrates asked for SHALL be derived from what the invocation carries rather than stated a second time.
A vendor release that adds a built-in outside that set therefore records an unexpected tool on every run until somebody looks, which is the direction to be wrong in.
Deriving the answer from a vendor's error message instead is a rule that stops working when the prose is edited, and it stops working by failing every review rather than by reporting the truth.
The list is what the reviewer could use at the point socrates read it, not a proof that nothing arrived later.
A source connecting after the opening event can contribute a tool that event did not name.

socrates SHALL therefore also take the external sources of tools the same opening event names, and SHALL record every one it names, because socrates asked for none.
That event names a source the session connected whether or not that source's tools had reached the tool list beside it.
The two readings are therefore incomplete in different places: the tool list misses a source whose tools are still arriving, and the source list misses a built-in the vendor added.
socrates measured the gap rather than assuming it, on an opening event that named four sources while the tool list held none of their tools.
That field SHALL be read where the same event's inventory of what is installed on the machine SHALL NOT.
What separates them is measured: the sources named went from four to none under the one mechanism, while the inventory beside them did not move.
Recording a source SHALL NOT fail the round, for the reason recording an absent tool list does not.
The observation bounds nothing, and a machine whose policy is set by somebody entitled to set it should still get reviews.
Unlike an absence, a source SHALL be recorded whether or not the round completed, because it is something the session said rather than something it failed to say.
What covers a source that connects after socrates has read the event is those customisations not loading at all, and the comparison after the round.

An opening event that reports no tools at all is not a session that held only what socrates asked for, and socrates SHALL NOT record it as one.
Where a round completes and the stream reported no tool list, the run SHALL record that it did not.
That is distinct from a session holding a tool socrates did not ask for, and a person reading the run needs to know which of the two happened.
That SHALL NOT fail the round: a release that stops emitting the field should report the truth rather than fail every review, which is the same direction the built-in above is wrong in.
A round that failed before the reviewer answered has no tool list either, and the failure is what there is to read, so the absence is recorded only where the round otherwise completed.

Where such a tool is in the session anyway, socrates SHALL run the review rather than not at all, and the run SHALL record the tools it did not ask for.
A machine whose policy is set by somebody entitled to set it should still get reviews.
Inside the snapshot worktree the seal holds the reviewer read-only.
Outside it nothing socrates holds bounds what such a tool reaches, so the tool surface is a mechanism of its own rather than something the seal answers.
Neither the tool surface nor the seal removes the comparison after the round, which is the check on anything that defeated both.

#### Scenario: No tool socrates did not ask for reaches the reviewer

- **WHEN** a reviewer is invoked on a machine whose user configuration names external sources of tools
- **THEN** the session holds the tools socrates asked for and no others, and the run records nothing to the contrary

#### Scenario: Nothing asks for the machine's own customisations

- **WHEN** any reviewer is invoked
- **THEN** the recorded invocation asks the vendor to disable the customisations the machine brings on its own account, supplies no settings of its own, and carries no request to replace the machine's own configuration of any source

#### Scenario: The reviewer is given tools that do not write

- **WHEN** any reviewer is invoked
- **THEN** the recorded invocation asks for the tools that read the change and check a claim, and for none of the tools that write

#### Scenario: A machine that hands the reviewer more than it asked for

- **WHEN** the session's opening event names a tool socrates did not ask for, whatever put it there
- **THEN** the review completes, the tool restriction is unchanged, and the run records the tools socrates did not ask for

#### Scenario: What the reviewer held is read rather than assumed

- **WHEN** a review completes and the stream's opening event names the session's tools
- **THEN** what the run records about the reviewer's tools is what that event reported, and no part of the decision reads the text of a vendor error

#### Scenario: A session that never said what it held

- **WHEN** a round completes and the stream's opening event reports no tools at all
- **THEN** the run records that the session did not say what it held, the round is not failed for it, and nothing records the session as holding only what socrates asked for

#### Scenario: A source named before its tools arrive

- **WHEN** the stream's opening event names an external source of tools and the session's tool list holds none of that source's tools
- **THEN** the run records the source socrates asked for none of, the round is not failed for it, and a tool list holding nothing unexpected does not leave the run reading as a session that held only what socrates asked for

#### Scenario: A tool that writes cannot reach the change under review

- **WHEN** a reviewer runs with a tool socrates did not ask for and it writes into its working directory as it starts
- **THEN** the write is refused, the review completes, and the run's verdict is whatever the findings say

#### Scenario: The invocation carries the schema

- **WHEN** a request supplying a schema is made
- **THEN** the recorded invocation carries print mode, a streaming JSON output format, and that exact schema

#### Scenario: Project settings are not loaded

- **WHEN** a reviewer is invoked against a repository that carries its own agent instructions
- **THEN** the recorded invocation restricts settings to user level, and the reviewer's prompt carries none of the repository's agent instructions

#### Scenario: The one file socrates reads is not the vendor loading the repository

- **WHEN** a reviewer is invoked against a repository carrying both its own agent instructions and the guidance file socrates reads
- **THEN** the recorded invocation restricts settings to user level, the prompt carries none of the repository's agent instructions, and the guidance reaches the prompt in its own section

#### Scenario: A stream that ends mid-event

- **WHEN** the agent's output stream is truncated part-way through a line
- **THEN** the call fails naming a truncated stream, and no partial finding set is returned

#### Scenario: A stream with no result event

- **WHEN** the agent's output stream completes with no terminal result event
- **THEN** the call fails naming the missing result, rather than returning empty output

#### Scenario: Structured output that violates the schema

- **WHEN** the agent returns structured output that parses as JSON but does not satisfy the schema
- **THEN** the call fails naming the violation, and the output is recorded for inspection

#### Scenario: A non-zero exit after a complete stream

- **WHEN** the agent emits a complete stream carrying structured output and then exits non-zero
- **THEN** the call fails, naming the exit status and carrying whatever the agent wrote to its error stream

#### Scenario: A non-zero exit after the agent said why

- **WHEN** the agent's result event reports an error carrying its own account of the failure and the agent then exits non-zero
- **THEN** the call fails naming the exit status and carrying that account, rather than the exit status beside an empty error stream

#### Scenario: A result that reports an error carries what it said

- **WHEN** the agent's result event reports an error carrying its own account of the failure and the agent exits zero
- **THEN** the call fails carrying that account, and the subtype beside it is not what is reported

#### Scenario: A result that reports an error and says nothing

- **WHEN** the agent's result event reports an error carrying no account of it
- **THEN** the call fails naming what that event did say about itself, rather than an empty quotation

#### Scenario: A failure the agent said nothing about

- **WHEN** the agent exits non-zero having written nothing to its error stream and reported nothing in its result event
- **THEN** the failure names the exit status and states that the agent said nothing about why, and carries no empty quotation

#### Scenario: An account too long to carry is shortened

- **WHEN** the account a failed call carries is longer than a failure reason will hold
- **THEN** the reason carries the beginning of it, marks that it was shortened, and does so on the same ceiling the agent's error stream is held to

#### Scenario: What the agent cost is recorded

- **WHEN** a call completes
- **THEN** the model that answered and the token counts it reported are recorded in the run's event log

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
