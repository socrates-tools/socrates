## RENAMED Requirements

- FROM: `#### Scenario: No tool server is loaded`
- TO: `#### Scenario: No tool socrates did not ask for reaches the reviewer`

- FROM: `#### Scenario: Nothing asks for the machine's own customisations`
- TO: `#### Scenario: Nothing asks for the machine's own customisations`

- FROM: `#### Scenario: A machine that will not let them be ignored`
- TO: `#### Scenario: A machine that hands the reviewer more than it asked for`

- FROM: `#### Scenario: What loaded is read rather than inferred`
- TO: `#### Scenario: What the reviewer held is read rather than assumed`

- FROM: `#### Scenario: A tool server that writes cannot reach the change under review`
- TO: `#### Scenario: A tool that writes cannot reach the change under review`

## MODIFIED Requirements

### Requirement: The Claude implementation speaks the CLI's wire contract

The implementation shipped SHALL invoke the `claude` binary in print mode with a streaming JSON output format, supply the JSON schema on the command line, and take the structured output from the terminal `result` event of the stream.

It SHALL load user-level settings and memory only, and SHALL NOT load the reviewed repository's own project or local settings, because those are authored by the thing under review.

The reviewer SHALL hold the tools socrates asked for and no others.
It asks on the vendor's own mechanism for tools that do not write, and that mechanism covers the built-in tools the model may call and nothing else.
A tool from any other source is therefore not covered by it: an external tool server contributes tools, so does a hook, so does a plugin, and so does a built-in the vendor adds in a release nobody read.
A reviewer holding such a tool is a reviewer that can write, both to the snapshot worktree and outside it.
socrates' first live run found this when the machine's own configured tool server wrote a directory into the snapshot worktree, and a later one found the reviewer's session holding script execution, file upload and browser control.

It SHALL ask for no other source of tools, with an empty allowlist supplied in the invocation's own settings, which the vendor evaluates before another source loads and which an empty value SHALL express as permitting none.
It SHALL NOT ask the vendor to replace the machine's own configuration of those sources.
That request is refused at startup wherever a managed configuration is deployed, because it asks to replace a managed set, and it is the only reason socrates ever needed a second attempt.
An allowlist that permits nothing makes reading those configuration files pointless, so one mechanism that is never refused replaces two of which one always was.

The mechanism SHALL NOT be trusted to have worked.
socrates SHALL take the session's tool list from the stream's opening event, SHALL treat every tool in it that socrates did not ask for as a tool it did not ask for whatever put it there, and SHALL record those.
An allowlist merges with the ones the machine already has, so a person's own settings can permit a source socrates did not, and a policy that makes managed allowlists the only ones that count can ignore what socrates supplied altogether.
Deriving the answer from a vendor's error message instead is a rule that stops working when the prose is edited, and it stops working by failing every review rather than by reporting the truth.

Where such a tool is in the session anyway, socrates SHALL run the review rather than not at all, and the run SHALL record the tools it did not ask for.
A machine whose policy is set by somebody entitled to set it should still get reviews.
Inside the snapshot worktree the seal holds the reviewer read-only, and outside it nothing socrates holds bounds what such a tool reaches, which is why the tool surface is a mechanism of its own rather than something the seal answers.
Neither this requirement nor the seal removes the comparison after the round, which is the check on anything that defeated both.

#### Scenario: No tool socrates did not ask for reaches the reviewer

- **WHEN** a reviewer is invoked on a machine whose user configuration names external sources of tools
- **THEN** the session holds the tools socrates asked for and no others, and the run records nothing to the contrary

#### Scenario: Nothing asks for the machine's own customisations

- **WHEN** any reviewer is invoked
- **THEN** the recorded invocation carries an empty allowlist for other tool sources in its own settings, and carries no request to replace the machine's own configuration of them

#### Scenario: A machine that hands the reviewer more than it asked for

- **WHEN** the session's opening event names a tool socrates did not ask for, whatever put it there
- **THEN** the review completes, the tool restriction is unchanged, and the run records the tools socrates did not ask for

#### Scenario: What the reviewer held is read rather than assumed

- **WHEN** a review completes and the stream's opening event names the session's tools
- **THEN** what the run records about the reviewer's tools is what that event reported, and no part of the decision reads the text of a vendor error

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

#### Scenario: What the agent cost is recorded

- **WHEN** a call completes
- **THEN** the model that answered and the token counts it reported are recorded in the run's event log
