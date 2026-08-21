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

The skill SHALL say that the session continues once every answerable finding is answered, without asking whether the round is worth spending.
A patch is unreviewed new code written by the session that was just told it was wrong, and a round that judges it is the round the findings earned.

Where a run ends, the skill SHALL say how many rounds it took and what the last one found.
Two rounds ending with nothing blocking and one round ending the same way are different results, and a session that reports only the verdict tells the person nothing about which of the two happened.

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

#### Scenario: A round is spent without asking

- **WHEN** the skill describes a parked round whose findings can all be patched
- **THEN** it says to answer them and continue, and it does not ask the person whether to spend the round

#### Scenario: What the person is told at the end

- **WHEN** the skill describes reporting a finished run
- **THEN** it says to report how many rounds ran and what the last one found

#### Scenario: A question is the person's

- **WHEN** the skill describes a parked round holding a finding that asks a person
- **THEN** it says to present that finding to the person, to record what they decided, and that the session does not answer it itself

## ADDED Requirements

### Requirement: A further look is a new run, and only where the person asks

The skill SHALL say that a run has at most two rounds, that a further look is a fresh `scrat review` rather than a third round, and that one is started only where the person asks for it.

It SHALL say that the reviewer of a new run is told nothing about the run before it, so the person knows what they are getting: a cold reading of the code as it now stands, and not a continuation of the argument.

A session that believes a third round exists will look for the flag that gives it one, and there is none.
A session that starts a fresh review on its own initiative spends a full review on a change nobody asked it to re-read, which is the cost this whole design is careful with.

#### Scenario: There is no third round

- **WHEN** the skill describes a run that escalated at the round cap
- **THEN** it says that no further round is possible and names a fresh review as the only way to look again

#### Scenario: A fresh review is the person's to ask for

- **WHEN** the skill describes starting another review after a run ended
- **THEN** it says to offer it and to wait for the person to ask, rather than starting one unprompted

#### Scenario: A new run knows nothing of the old one

- **WHEN** the skill describes a fresh review of patched code
- **THEN** it says the reviewer is told nothing about the earlier run, its findings or its responses
