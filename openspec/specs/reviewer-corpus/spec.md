# reviewer-corpus Specification

## Purpose
TBD - created by archiving change what-a-reviewer-actually-answered. Update Purpose after archive.
## Requirements
### Requirement: What a real reviewer answered is kept

socrates SHALL keep a corpus of findings documents returned by real reviewers, each stored verbatim.

An entry SHALL carry, beside the document: the model, the effort, the round number, the base commit, the snapshot commit, the run's name, the date it was recorded, and one line naming in words what was reviewed.

The document SHALL be stored exactly as the reviewer returned it.
A corpus that has been tidied is a record of somebody's tidying, and the fields worth studying are the ones nobody expected.

#### Scenario: An entry holds the answer and its circumstances

- **WHEN** a recorded answer is read back
- **THEN** the document is byte-identical to what the reviewer returned, and the model, effort, round, base, snapshot, run and subject are all present

#### Scenario: An entry that names no model is refused

- **WHEN** an entry is missing the model, the effort or the round it came from
- **THEN** it is refused, naming what is absent, because an answer without them is not evidence of anything

### Requirement: The corpus is socrates reviewing socrates

Recording SHALL refuse unless the repository the run reviewed is socrates' own.

A findings document quotes the code it was written about.
Recording from any other repository would check somebody else's source into this one, and a redaction pass over model prose is a weaker rule than never holding the prose at all.

#### Scenario: Recording from another repository is refused

- **WHEN** recording is attempted against a run whose repository is not socrates
- **THEN** it is refused, naming the repository, and nothing is written

### Requirement: The suite says which outcomes only fixtures have reached

A report SHALL classify every corpus entry through the routing rule and print, for each outcome the rule can produce, how many recorded answers reach it and how many tests in the suite name it.

The second number SHALL be counted from the tests rather than from a list somebody maintains, and SHALL be labelled as what it is.
A test that names an outcome is not the same as a test that asserts it, and a report that called it coverage would be making the claim this report exists to stop.

It SHALL name in words any outcome that no recorded answer has ever reached.

It SHALL NOT fail for such an outcome.
Whether a real reviewer ever produces one is not this repository's to control, and a check that is permanently red is a check people learn to skip.

The routing rule itself SHALL do the classifying.
A report that restated the rule would be an oracle derived from the sentence the implementation was derived from, which is the defect this repository has recorded twice.

#### Scenario: An outcome no reviewer has reached is named

- **WHEN** the report runs against a corpus holding no answer that would park a round
- **THEN** it says the parked outcome has never been observed from a real reviewer, and does not fail

#### Scenario: The report counts what it classified

- **WHEN** the report runs
- **THEN** every outcome carries the number of recorded answers reaching it and the number of tests naming it, and the totals account for every entry

### Requirement: The recorded distribution cannot drift silently

The distribution the report prints SHALL be written to a file that is kept under version control, and a test SHALL fail when the corpus no longer produces it.

An entry that disappears, and a change to the routing rule that moves what the corpus reaches, are both things that should appear as a diff somebody reads rather than as a number nobody was watching.

#### Scenario: A removed entry is visible

- **WHEN** an entry is removed from the corpus
- **THEN** the test comparing the recorded distribution fails, naming what changed

#### Scenario: A routing change that moves the corpus is visible

- **WHEN** the routing rule changes so that a recorded answer reaches a different outcome
- **THEN** the test comparing the recorded distribution fails, and the file is what records the new truth once somebody accepts it

### Requirement: A verification task may not require a particular answer

A task that calls for a live run SHALL be satisfiable by what a live run can produce.

It MAY require that the run happened, that its answer was recorded into the corpus, and that socrates routed whatever came back as this specification requires.
It SHALL NOT require the answer to hold a particular severity, action, count or shape.

What a reviewer finds is the reviewer's.
A task that asks a model to return something in particular is a task the author cannot satisfy by working, and the four attempts that closed no such task are what this requirement exists to prevent repeating.

#### Scenario: A live task is satisfied by recording whatever came back

- **WHEN** a live run returns an answer that reaches an outcome the task's author did not expect
- **THEN** the task is satisfied by the answer being recorded and correctly routed, and the unexpected outcome is what the report then reports

