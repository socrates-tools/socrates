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

