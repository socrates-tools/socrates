## MODIFIED Requirements

### Requirement: The skill teaches the loop, including every terminal branch

The skill SHALL describe the sequence a session follows: implement, start a review carrying the person's stated intent, poll until the run is not working, and act on each terminal state.

It SHALL name what to do for each of `clean`, `failed`, `escalated` and `awaiting-response`, and `failed` SHALL be described before `escalated`.

A session scanning for what to do about findings reads past a state that has none, and a `failed` run has none.

It SHALL also name `cancelled`, and SHALL name every state a run can rest in.
A person who stops a review from another terminal leaves a state that is neither a verdict nor work in progress, and a loop told to wait for one of a fixed list of states polls a cancelled run for as long as anybody lets it.

The `awaiting-response` branch SHALL separate the findings the session answers from the findings the person answers.
A finding that can be patched is the session's to patch and record.
A finding that asks a person SHALL be described as the person's, to be taken to them and answered with what they decided, and the skill SHALL say that the session does not answer one on its own behalf.

Until a question ended the run, nothing needed to be said: a terminal state was what stopped a session answering it.
A parked question removes that, so the sentence is what remains, and it is the reason this requirement states it rather than leaving it to be inferred from the loop.

#### Scenario: Every terminal state is covered

- **WHEN** the skill is rendered
- **THEN** it names `clean`, `failed`, `escalated` and `awaiting-response`, and says what to do for each

#### Scenario: A failed run is stated not to be a pass

- **WHEN** the skill describes `failed`
- **THEN** it says that no review completed, that nothing has been checked, and that this is not a clean result

#### Scenario: A cancelled run is not a verdict

- **WHEN** the skill describes `cancelled`
- **THEN** it says that nothing was decided, that the change is not reviewed, and that somebody stopped it on purpose

#### Scenario: Every state a run rests in is named

- **WHEN** the states the skill names are compared with the states a run can rest in
- **THEN** every one of them is named, so the loop has no state it would poll through

#### Scenario: The loop closes

- **WHEN** the skill describes `awaiting-response`
- **THEN** it says to respond to every blocking finding, then continue, then return to polling

#### Scenario: A question is the person's

- **WHEN** the skill describes a parked round holding a finding that asks a person
- **THEN** it says to present that finding to the person, to record what they decided, and that the session does not answer it itself
