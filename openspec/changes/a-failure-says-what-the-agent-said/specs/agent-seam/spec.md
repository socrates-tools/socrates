## MODIFIED Requirements

### Requirement: The Claude implementation speaks the CLI's wire contract

The implementation shipped SHALL invoke the `claude` binary in print mode with a streaming JSON output format, supply the JSON schema on the command line, and take the structured output from the terminal `result` event of the stream.

It SHALL load user-level settings only, and SHALL NOT load the reviewed repository's own project or local settings, because those are authored by the thing under review.

A call that failed SHALL carry the account the agent gave of its own failure, where the agent gave one.
That terminal event is where the CLI says why the work did not happen, and the agent's error stream is where a process complains about its arguments, its environment and its own startup.
The account SHALL therefore be preferred to the error stream, and both SHALL be carried where both are there.
socrates read the exit status first and returned on it, so a live review that died on the provider being overloaded reported an exit status and an empty error stream, while the message naming the provider sat parsed and unread in the same stream socrates had already finished reading.
A reason answers one question, which is whether to run the review again or to fix something, and an exit status answers it either way.

The subtype beside that account SHALL NOT be the reason where the account is there.
On the same live stream the subtype was `success`, on the event whose error flag was true and whose message named the provider, so of that event's three fields socrates named the only one that said nothing.
It SHALL remain what is said where the event reports an error and carries no account, because a failure has to be reported as something.

Where the agent said nothing on either channel the failure SHALL say that it said nothing, rather than quoting an empty one.

The reason SHALL be bounded, on the ceiling and the mechanism the error stream already has rather than a second one, because it reaches a run's record and it reaches the context of a session that polls `socrates status` while waiting to patch, and the agent's text on some failures is a whole turn.

socrates SHALL NOT decide anything from that account.
It is not parsed, not matched and not compared, no state, verdict or exit code depends on a word in it, and the call had already failed before it was read.
Deriving a decision from a vendor's prose is refused below for the reason it is refused: quoting prose to a person fails by being less useful, and branching on it fails by failing every review.

The reviewer SHALL hold the tools socrates asked for and no others.
It asks on the vendor's own mechanism for tools that do not write, and that mechanism covers the built-in tools the model may call and nothing else.
A tool from any other source is therefore not covered by it: an external tool server contributes tools, so does a plugin, and so does a built-in the vendor adds in a release nobody read.
A hook contributes no tool and is worse than one for it, because it runs a command of its own on an event rather than waiting to be called.
A reviewer holding such a tool is a reviewer that can write, both to the snapshot worktree and outside it.
socrates' first live run found this when the machine's own configured tool server wrote a directory into the snapshot worktree, and a later one found the reviewer's session holding script execution, file upload and browser control.

The tools socrates asks for SHALL be those that read the change and those that check a claim against the world, and SHALL NOT be those that write.
A reviewer that cannot look up what a library actually does writes a finding about its own recollection, and patching is the author's half of the asymmetry this design rests on.
The tools that write are also the ones that reach outside the snapshot worktree, where nothing socrates holds bounds anything.

It SHALL ask the vendor to disable the customisations the machine brings on its own account, on the one mechanism that disables all of those at once and leaves auth, model selection and the built-in tools working.
It SHALL supply no settings of its own, because that mechanism makes an allowlist for one source of tools pointless and covers the sources such an allowlist said nothing about.
It SHALL NOT ask the vendor to replace the machine's own configuration of any source.
That request is refused at startup wherever a managed configuration is deployed, because it asks to replace a managed set, and it is the only reason socrates ever needed a second attempt.
One mechanism that is never refused therefore replaces two of which one always was.

What that mechanism does not reach is a hook an administrator's policy configures, which the vendor keeps in force by design, along with a policy-configured status line and policy-configured file-suggestion commands.
Only the policy-configured ones survive it, and a policy-configured tool source is not among them: those do not load, which socrates measured on such a machine, where a plugin's own tool left the session's list at the same time.
The vendor documents managed plugins, managed skills and managed memory going the same way, and socrates holds no measurement that distinguishes a managed one from a person's own, because the session's tool list does not say what contributed a tool.
A policy hook contributes no tool, so the session's tool list is not what would catch it.
What does catch it inside the snapshot worktree is the seal and the comparison, because its command is a child of the session and the session runs there.
What nothing socrates holds bounds is where that command reaches outside the worktree, and that one residual is accepted and stated rather than covered, for the same reason the review runs at all on such a machine: a policy set by somebody entitled to set it should still get reviews.

The mechanism SHALL NOT be trusted to have worked.
socrates SHALL take the session's tool list from the stream's opening event, SHALL treat every tool in it that socrates did not ask for as a tool it did not ask for whatever put it there, and SHALL record those.
The set socrates asked for SHALL be derived from what the invocation carries rather than stated a second time.
A vendor release that adds a built-in outside that set therefore records an unexpected tool on every run until somebody looks, which is the direction to be wrong in.
Deriving the answer from a vendor's error message instead is a rule that stops working when the prose is edited, and it stops working by failing every review rather than by reporting the truth.
The list is what the reviewer could use at the point socrates read it rather than a proof that nothing arrived later, because a source that connects after the opening event can contribute a tool that event did not name.

socrates SHALL therefore also take the external sources of tools the same opening event names, and SHALL record every one it names, because socrates asked for none.
That event names a source the session connected whether or not that source's tools had reached the tool list beside it, which is what makes the two readings incomplete in different places: the tool list misses a source whose tools are still arriving, and the source list misses a built-in the vendor added.
socrates measured the gap rather than assuming it, on an opening event that named four sources while the tool list held none of their tools.
That field SHALL be read where the same event's inventory of what is installed on the machine SHALL NOT, and what separates them is also measured: the sources named went from four to none under the one mechanism, while the inventory beside them did not move.
Recording a source SHALL NOT fail the round, for the reason recording an absent tool list does not: the observation bounds nothing, and a machine whose policy is set by somebody entitled to set it should still get reviews.
Unlike an absence, a source SHALL be recorded whether or not the round completed, because it is something the session said rather than something it failed to say.
What covers a source that connects after socrates has read the event is those customisations not loading at all, and the comparison after the round.

An opening event that reports no tools at all is not a session that held only what socrates asked for, and socrates SHALL NOT record it as one.
Where a round completes and the stream reported no tool list, the run SHALL record that it did not, distinctly from a session holding a tool socrates did not ask for, because a person reading the run needs to know which of the two happened.
That SHALL NOT fail the round: a release that stops emitting the field should report the truth rather than fail every review, which is the same direction the built-in above is wrong in.
A round that failed before the reviewer answered has no tool list either, and the failure is what there is to read, so the absence is recorded only where the round otherwise completed.

Where such a tool is in the session anyway, socrates SHALL run the review rather than not at all, and the run SHALL record the tools it did not ask for.
A machine whose policy is set by somebody entitled to set it should still get reviews.
Inside the snapshot worktree the seal holds the reviewer read-only, and outside it nothing socrates holds bounds what such a tool reaches, which is why the tool surface is a mechanism of its own rather than something the seal answers.
Neither this requirement nor the seal removes the comparison after the round, which is the check on anything that defeated both.

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
- **THEN** the recorded invocation restricts settings to user level, and the reviewer's prompt contains none of the repository's instructions

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
