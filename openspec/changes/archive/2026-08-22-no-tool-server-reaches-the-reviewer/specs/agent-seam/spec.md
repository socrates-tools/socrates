## MODIFIED Requirements

### Requirement: The Claude implementation speaks the CLI's wire contract

The implementation shipped SHALL invoke the `claude` binary in print mode with a streaming JSON output format, supply the JSON schema on the command line, and take the structured output from the terminal `result` event of the stream.

It SHALL load user-level settings and memory only, and SHALL NOT load the reviewed repository's own project or local settings, because those are authored by the thing under review.

It SHALL restrict the agent's tools to those that do not write, on the vendor's own mechanism.
That restriction covers the built-in tools the model may call and nothing else.
It does not cover a tool that a tool server contributes, which is why stopping the servers loading is a separate mechanism rather than the same one stated twice.

It SHALL load no external tool server, whatever the person has configured for their own sessions.
A tool server is a second source of tools that the vendor's own tool restriction does not cover, and a review that loads one is a reviewer that can write, both to the snapshot worktree and outside it.
The first live run of socrates found this: the person's own configured servers wrote a directory into the snapshot worktree, and the round failed for a write the reviewer never asked to make.
A later live run found the other half of it: on a machine whose managed configuration the CLI would not let socrates ignore, the reviewer's session held script execution, file upload and browser control, and nothing socrates holds covers what those reach.

It SHALL stop them with an empty tool-server allowlist, supplied with the invocation's own settings, which the vendor evaluates before a server loads and which an empty value SHALL express as permitting none.
It SHALL NOT ask the vendor to replace the machine's own tool server configuration.
That request is refused at startup wherever a managed configuration is deployed, because it asks to replace a managed set, and it is the only reason socrates ever needed a second attempt.
An allowlist that permits nothing makes reading those configuration files pointless, so one mechanism that is never refused replaces two of which one always was.

The mechanism SHALL NOT be trusted to have worked.
socrates SHALL take the tool servers and the session's tool list from the stream's opening event and SHALL record what that event reported.
An allowlist merges with the ones the machine already has, so a person's own settings can permit a server socrates did not, and a policy that makes managed allowlists the only ones that count can ignore what socrates supplied altogether.
Deriving the answer from a vendor's error message instead is a rule that stops working when the prose is edited, and it stops working by failing every review rather than by reporting the truth.

Where the machine loads them anyway, socrates SHALL run the review rather than not at all, and the run SHALL record which servers the opening event named and how many tools they contributed.
A machine whose policy is set by somebody entitled to set it should still get reviews.
Inside the snapshot worktree the seal holds the reviewer read-only, and outside it nothing socrates holds bounds what a loaded tool reaches, which is why loading none of them is a mechanism of its own rather than something the seal answers.
Neither this requirement nor the seal removes the comparison after the round, which is the check on anything that defeated both.

#### Scenario: No tool server is loaded

- **WHEN** a reviewer is invoked on a machine whose user configuration names external tool servers
- **THEN** the recorded invocation loads none of them, and nothing outside the vendor's restricted tools reaches the reviewer

#### Scenario: Nothing asks to replace a managed set

- **WHEN** any reviewer is invoked
- **THEN** the recorded invocation carries an empty tool-server allowlist in its own settings, and carries no request to replace the machine's own tool server configuration

#### Scenario: A machine that will not let them be ignored

- **WHEN** the machine's policy overrides the allowlist socrates supplied, so the opening event names tool servers
- **THEN** the review completes, the tool restriction is unchanged, and the run records the servers that event named and how many tools they contributed

#### Scenario: What loaded is read rather than inferred

- **WHEN** a review completes and the stream's opening event names the tool servers of the session
- **THEN** what the run records about tool servers is what that event reported, and no part of the decision reads the text of a vendor error

#### Scenario: A tool server that writes cannot reach the change under review

- **WHEN** a reviewer runs with the machine's tool servers loaded and one of them writes into its working directory as it starts
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
