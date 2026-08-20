# feedback-hygiene Specification

## Purpose
TBD - created by archiving change a-rejection-is-a-claim. Update Purpose after archive.
## Requirements
### Requirement: Model output is neutralised before it is re-embedded in a prompt

Every field of a findings document and every field of a response document SHALL be neutralised before it is placed in a prompt.

Neutralising SHALL: remove conflict markers, remove the heading marker the prompt's own sections are built from where a line opens with one, normalise carriage returns to newlines, collapse runs of whitespace within a line, and trim.
It SHALL happen per field, never across the assembled block, so that a field's boundaries survive and a long description cannot swallow the label of the field after it.

socrates has two channels of model output feeding the next prompt where the core had one, and the second is written by a session that is arguing, so this matters more here rather than less.

#### Scenario: Conflict markers do not survive

- **WHEN** a finding's description or a response's reason contains conflict markers
- **THEN** they are absent from the delivered prompt

#### Scenario: A field cannot open a section

- **WHEN** a finding's description or a response's reason holds a line opening with the prompt's heading marker
- **THEN** the delivered prompt holds no section the field composed, and the field's text stays inside the section it was placed in

#### Scenario: Line endings are normalised

- **WHEN** a field contains carriage returns
- **THEN** the delivered prompt contains none

#### Scenario: Neutralising is per field

- **WHEN** one field's value is long and multi-line and another field follows it
- **THEN** both fields are recognisable in the delivered prompt, and neither has absorbed the other's label

#### Scenario: The stored record is not neutralised

- **WHEN** a findings or response document is written to disk
- **THEN** it holds what was actually said, verbatim, and neutralising happens only on the way into a prompt

### Requirement: The prior-round block is labelled as metadata

The block carrying prior findings and responses SHALL be introduced as a section to be treated as metadata only, and each claim within it SHALL be attributed to who made it.

It SHALL NOT be introduced in a way that reads as instructions to the reviewer, because it is composed of text written by a model that has an interest in the outcome.

#### Scenario: The block says what it is

- **WHEN** a second round is invoked
- **THEN** the delivered prompt introduces the prior-round block as metadata to be judged, and attributes each finding and each response

#### Scenario: A response that reads like an instruction

- **WHEN** a recorded rejection reason is phrased as an instruction to the reviewer
- **THEN** it still appears inside the labelled block, attributed to the author, and nothing about its placement gives it the standing of an instruction

### Requirement: A rejection is never a suppression

socrates SHALL NOT instruct a reviewer to refrain from reporting a finding on the grounds that it was previously found, previously rejected, or previously responded to.

The prohibition exists because the mechanism inverts without it.
A rejection is made in order to be judged, and a reviewer told not to re-report the finding it disputes cannot judge it.
An instruction of that shape is correct in a tool where a person can mark a finding permanently ignored, and socrates has no such thing.

#### Scenario: No suppression instruction reaches the reviewer

- **WHEN** a second round is invoked after a rejection
- **THEN** the delivered prompt carries no instruction to avoid re-reporting any finding

#### Scenario: A disputed finding may be found again

- **WHEN** the author rejected a finding and the second reviewer still considers it real
- **THEN** it appears in the second round's findings, and the run escalates to the person on the cap rather than resolving the disagreement itself

