## MODIFIED Requirements

### Requirement: The Claude implementation speaks the CLI's wire contract

The implementation shipped SHALL invoke the `claude` binary in print mode with a streaming JSON output format, supply the JSON schema on the command line, and take the structured output from the terminal `result` event of the stream.

It SHALL load user-level settings and memory only, and SHALL NOT load the reviewed repository's own project or local settings, because those are authored by the thing under review.

It SHALL restrict the agent's tools to those that do not write, on the vendor's own mechanism.
That restriction covers the tools the model may call and nothing else.

It SHALL load no external tool server, whatever the person has configured for their own sessions.
A tool server is a second source of tools that the vendor's own tool restriction does not cover, and a review that loads one is a reviewer that can write, both to the snapshot worktree and outside it.
The first live run of socrates found this: the person's own configured servers wrote a directory into the snapshot worktree, and the round failed for a write the reviewer never asked to make.

Where the vendor's CLI refuses to ignore the machine's own tool server configuration, socrates SHALL run the review without that restriction rather than not at all, and SHALL record that the servers were loaded.
What holds the reviewer read-only in that case is the seal on the snapshot worktree, which no configuration of the vendor's CLI can lift.
Neither this requirement nor the seal removes the comparison after the round, which is the check on anything that defeated both.

#### Scenario: No tool server is loaded

- **WHEN** a reviewer is invoked on a machine whose user configuration names external tool servers
- **THEN** the recorded invocation loads none of them, and nothing outside the vendor's restricted tools reaches the reviewer

#### Scenario: A machine that will not let them be ignored

- **WHEN** the CLI refuses the request to ignore the machine's own tool server configuration
- **THEN** the review is invoked again without that request, the run records that the servers were loaded, and the tool restriction is unchanged

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
