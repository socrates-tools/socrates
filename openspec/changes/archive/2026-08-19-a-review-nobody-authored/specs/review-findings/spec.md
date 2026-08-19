## ADDED Requirements

### Requirement: A finding is a schema, not prose

The reviewer SHALL be forced into a JSON document rather than asked for prose, and that document SHALL be the only contract between a reviewer and everything downstream.

Each finding SHALL carry an identifier, a severity, a file, a one-indexed line, a description and an action.
The document SHALL carry a risk level and a one-sentence rationale for it.

```
{
  "findings": [{
    "id":          "f1",
    "severity":    "error" | "warning" | "info",
    "file":        "internal/snapshot/carve.go",
    "line":        42,
    "description": "...",
    "action":      "no-op" | "auto-fix" | "ask-user"
  }],
  "risk_level":     "low" | "medium" | "high",
  "risk_rationale": "one sentence"
}
```

Identifiers SHALL be stable within a run, because a later layer attaches a written response to a specific finding across rounds and an identifier that moves attaches it to the wrong one.

There SHALL be no field classifying a finding by delivery phase.
`no-mistakes` needs one because push, pull request and continuous integration are later steps there and the reviewer kept reporting "the pull request is not open yet" as a defect.
socrates has no delivery phase, so the field has nothing to classify.

#### Scenario: The reviewer is given the schema

- **WHEN** a reviewer is invoked
- **THEN** the invocation carries this schema, and a response that does not satisfy it fails the round

#### Scenario: Identifiers are unique within a run

- **WHEN** a findings document carries two findings with the same identifier
- **THEN** the round fails naming the duplicate, because a response attached to it would be ambiguous

#### Scenario: A finding names a file and a line

- **WHEN** a finding is recorded
- **THEN** it carries a file path relative to the repository root and a one-indexed line

### Requirement: Severity and action are orthogonal

Severity SHALL decide whether the run blocks.
Action SHALL decide who resolves it.
Neither SHALL be derived from the other.

`error` and `warning` block; `info` does not.
An `info` finding actioned `auto-fix` is therefore something to patch that never blocks, and an `error` actioned `ask-user` is something that blocks and is never handed to a patching agent.

#### Scenario: An informational finding does not block

- **WHEN** a review returns only findings of severity `info`
- **THEN** the run ends `clean`

#### Scenario: A warning blocks

- **WHEN** a review returns a finding of severity `warning`
- **THEN** the run does not end `clean`

#### Scenario: Action does not change whether a run blocks

- **WHEN** two runs return the same `error` finding, one actioned `auto-fix` and one actioned `ask-user`
- **THEN** both block, and they differ only in who the finding is for

#### Scenario: An informational finding can still be actionable

- **WHEN** a review returns an `info` finding actioned `auto-fix`
- **THEN** it is recorded as something to fix and the run still ends `clean`

### Requirement: An unclassified finding belongs to the person

A finding whose action is missing, empty, or not one of the three known values SHALL be treated as `ask-user`.
It SHALL NEVER default to `auto-fix`.

An unclassified finding silently handed to a patching agent is a change to somebody's code that nobody decided to make.

#### Scenario: A missing action defaults to needing a person

- **WHEN** a finding carries no action field
- **THEN** it is treated as `ask-user`

#### Scenario: An empty action defaults to needing a person

- **WHEN** a finding carries an empty action
- **THEN** it is treated as `ask-user`

#### Scenario: An unrecognised action defaults to needing a person

- **WHEN** a finding carries an action that is not one of the three known values
- **THEN** it is treated as `ask-user`, and the unrecognised value is recorded so it is visible

### Requirement: One helper answers every question about a finding's action

The question "does anything here need a person" and the question "which findings may be patched" SHALL be answered by the same defaulting helper.

Two implementations of the same defaulting rule can disagree, and a finding that falls into the gap between them is neither patched nor surfaced.
It is the defect this requirement exists to make unrepresentable, and it is invisible to line coverage.

#### Scenario: Every finding is in exactly one bucket

- **WHEN** any findings document is classified
- **THEN** every finding is either needing a person, patchable, or a no-op, in exactly one of the three, and the three together account for all of them

#### Scenario: The two consumers agree by construction

- **WHEN** the source is searched for code deciding a finding's effective action
- **THEN** exactly one function does it, and both consumers call it

#### Scenario: A document mixing every case

- **WHEN** a document carries findings actioned `auto-fix`, `ask-user`, `no-op`, an unknown value, an empty value and no value at all
- **THEN** the classification places four of them with the person, one as patchable and one as a no-op, and nothing is lost

### Requirement: The verdict follows from the findings by one rule

A run whose review completed SHALL end `clean` when no finding blocks, and `escalated` otherwise.
The rule SHALL be a pure function of the classified findings, so that it can be tested exhaustively against the whole input space rather than sampled.

#### Scenario: No findings at all

- **WHEN** a completed review returns an empty findings array
- **THEN** the verdict is `clean`

#### Scenario: Only no-ops and informational findings

- **WHEN** a completed review returns findings that are all `info`
- **THEN** the verdict is `clean` and the findings are still recorded and readable

#### Scenario: Anything blocking

- **WHEN** a completed review returns at least one finding of severity `error` or `warning`
- **THEN** the verdict is `escalated`

#### Scenario: The rule reaches nothing outward

- **WHEN** the package holding the verdict rule is examined
- **THEN** it imports nothing that reads the filesystem, the network, the clock or the process environment

### Requirement: The risk level is recorded and does not decide the verdict

The risk level and its rationale SHALL be recorded and reported to the person.
They SHALL NOT participate in whether the run blocks, which is decided by severity alone.

A second gate on a model's self-assessed risk is a second place for the verdict to be decided, and two places is one too many.

#### Scenario: A high risk level with no blocking findings

- **WHEN** a review returns `risk_level` of `high` and only `info` findings
- **THEN** the verdict is `clean`, and the risk level and its rationale are reported

#### Scenario: The rationale is one sentence and is kept

- **WHEN** a review returns a risk rationale
- **THEN** it is recorded verbatim and shown wherever the verdict is shown
