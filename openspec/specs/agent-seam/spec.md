# agent-seam Specification

## Purpose
TBD - created by archiving change a-review-nobody-authored. Update Purpose after archive.
## Requirements
### Requirement: Every model call passes through one interface

socrates SHALL reach a coding agent through exactly one internal interface, taking a role, a prompt, a working directory and a JSON schema, and answering structured output, text, the model that answered and what it cost.

No other code SHALL start a model process, construct a vendor's arguments or parse a vendor's output.
The seam exists so that the test suite can substitute the agent and so that a per-role model choice has one place to live, not so that many vendors can plug in.

#### Scenario: One place starts the process

- **WHEN** the source is searched for the construction of a process that runs a coding agent
- **THEN** exactly one package does it, and every caller reaches it through the interface

#### Scenario: A schema is honoured or the call fails

- **WHEN** a call supplies a JSON schema and the agent returns no structured output
- **THEN** the call fails naming that structured output was required and absent, and does not fall back to parsing the text

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
Substituting at the Go interface SHALL NOT be the primary mechanism, because argument construction, stream parsing, process handling and reaping are where the defects are and a Go-level fake exercises none of them.

The simulator SHALL be driven by a scenario file, SHALL be able to produce every failure named in this specification on demand, and SHALL append every invocation it received, including arguments, working directory, environment and prompt, to a log a test can assert on.

A scenario SHALL be able to name a recorded answer instead of carrying one, and the simulator SHALL then emit that entry's findings document as its structured result, byte for byte.
It SHALL NOT adapt a recorded document to the prompt it received, renumber its findings, or supply a field the document lacks: a simulator that repairs a recorded answer is answering from its author's beliefs again.
A scenario naming an answer the corpus does not hold SHALL be refused, naming it, rather than answered with something else.

#### Scenario: A test asserts what was asked

- **WHEN** a review runs against the simulator
- **THEN** the invocation log holds the exact prompt, arguments and working directory socrates used

#### Scenario: The simulator produces a failure on demand

- **WHEN** a scenario asks for a truncated stream, a missing result event, a schema violation, a non-zero exit, silence, or a child process that outlives the agent
- **THEN** the simulator produces exactly that, and socrates behaves as this specification requires for it

#### Scenario: A scenario names a recorded answer

- **WHEN** a scenario names a corpus entry and a review runs
- **THEN** the reviewer's structured result is that entry's document, byte for byte

#### Scenario: A named answer that does not exist is refused

- **WHEN** a scenario names an entry the corpus does not hold
- **THEN** the simulator refuses, naming the entry, and no review is answered with a substitute

### Requirement: The simulator is checked against the real tool

Responses from the real `claude` CLI SHALL be recorded, and a drift check SHALL compare the simulator's output shape against those recordings.
A separate check SHALL run one real review against the real binary and assert that the arguments are still accepted and the schema still honoured.

What a real reviewer answered SHALL be recorded too, and not only the shape it answered in.
The shapes say the stream still parses; the answers are the only evidence about what the reviewer actually reports, and the scenarios the suite runs on were written before anybody had looked at one.

Both SHALL be behind build tags, so that the ordinary suite compiles neither and cannot reach the network by accident.
Replaying a recorded answer SHALL NOT be, because it reaches nothing.

A simulator nobody compares to the real thing becomes a specification of its author's beliefs, and a green suite against one is worse than no suite because it manufactures confidence.

#### Scenario: Drift is visible

- **WHEN** the recorded responses no longer match the shapes the simulator produces
- **THEN** the drift check fails, naming the fields that differ

#### Scenario: A live review's answer is kept

- **WHEN** a live review against the real binary returns a findings document
- **THEN** the document is recorded into the corpus with the model, effort and round it came from, whatever it found

#### Scenario: The ordinary suite reaches no network

- **WHEN** the default test command runs
- **THEN** neither the recording nor the live check is compiled, no test invokes the real `claude` binary, and replaying a recorded answer still works

