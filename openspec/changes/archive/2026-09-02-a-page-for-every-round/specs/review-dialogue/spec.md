## RENAMED Requirements

- FROM: `#### Scenario: An escalated round is one a person reads`
- TO: `#### Scenario: An escalated round waits on a person`

- FROM: `#### Scenario: A parked round holding a question is one a person reads`
- TO: `#### Scenario: A parked round holding a question waits on a person`

- FROM: `#### Scenario: A parked round the author answers alone is not one a person reads`
- TO: `#### Scenario: A parked round the author answers alone waits on nobody`

- FROM: `#### Scenario: A clean round is not one a person reads`
- TO: `#### Scenario: A clean round waits on nobody`

## ADDED Requirements

### Requirement: Which rounds wait on a person is one rule beside the routing rule

Whether a completed round waits on a person SHALL be a pure function of the classified findings and the state that round led to, in the same package as the rule that decided that state.

A round waits on a person where the round cannot move until a person answers:

- `escalated`: always
- `awaiting-response` holding at least one finding that a person answers: always
- `awaiting-response` where every finding it waits on is one the author patches: never
- `clean`: never
- any other state a run can hold: never

The rule SHALL decide the wait and SHALL NOT decide whether a page is served.
Every round that completed a review is handed over as a page, so the question the rule answers is which of those pages the session stops at, not which of them exist.
The rule it replaces answered both, and the two answers are not the same answer: a clean round is one nobody is waiting on and still one a person wants to read.
Reading a page is the person's to do whenever they like, and waiting on one is the session's to do only where an answer is owed.

It SHALL take the state rather than the round number.
The round number would make it re-decide the routing that has already been decided and recorded, and a rule that recomputes an answer somebody has stored can disagree with the stored one.
Reading the state is reading that one decision's result, which is the property this requirement is for.

A `clean` round waits on nobody even where it returned findings, because nothing is being asked of anybody and there is no decision to put in front of them.
A state that is not the outcome of a completed round waits on nobody either, so a failed run and a cancelled run need no case of their own anywhere else.

The rule SHALL be decided by this one function and SHALL NOT be re-derived anywhere from the run's state and its counts.
Both facts a re-derivation would read are already reported, so the second reader is the cheap thing to write and the expensive thing to be wrong about: two places computing whether a round waits can disagree, and the disagreement is silent, because a session that waits on a round nobody will answer stops for good, and a session that does not wait on one that is owed an answer patches around the person.

The rule SHALL be testable against its whole input space rather than sampled, and it SHALL reach no filesystem, no clock and no process.

#### Scenario: An escalated round waits on a person

- **WHEN** a round leads to `escalated`, whether by a blocking finding whose own action asks for nothing or by the round cap
- **THEN** it waits on a person

#### Scenario: A parked round holding a question waits on a person

- **WHEN** a round parks and at least one of the findings it waits on is answered by a person
- **THEN** it waits on a person

#### Scenario: A parked round the author answers alone waits on nobody

- **WHEN** a round parks and every finding it waits on is one the author patches or rejects
- **THEN** it waits on nobody

#### Scenario: A clean round waits on nobody

- **WHEN** a round leads to `clean`, with or without findings that block nothing
- **THEN** it waits on nobody

#### Scenario: The rule is decided over its whole input space

- **WHEN** the rule is exercised across every combination of classified findings and every state a run can hold
- **THEN** it agrees with the routing rule for every state that rule can reach, and answers that no round waits on a person in any state that rule cannot reach

## REMOVED Requirements

### Requirement: Which rounds a person reads is one rule beside the routing rule

**Reason**: The rule decided two things and only one of them still needs deciding.
It answered whether a page is served and whether the session stops at it, with one predicate over the same inputs.
Every round that completed a review is now handed over as a page, so the serving half is answered for every round and the rule that remains is about the wait.

**Migration**: The predicate keeps its inputs, its outputs and every case of its table, and is stated as the wait rather than the reading: `Which rounds wait on a person is one rule beside the routing rule` replaces it, case for case, with nothing added and nothing dropped.
A reader who wants which rounds a person reads now reads `findings-page`, where the answer is all of them.
