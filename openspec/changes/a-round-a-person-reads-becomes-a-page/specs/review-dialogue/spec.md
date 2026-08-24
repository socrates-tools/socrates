## ADDED Requirements

### Requirement: Which rounds a person reads is one rule beside the routing rule

Whether a completed round is one a person reads SHALL be a pure function of the classified findings and the state that round led to, in the same package as the rule that decided that state.

A round is one a person reads where socrates is finished and the findings are the person's, and where the round cannot move until a person answers:

- `escalated`: always
- `awaiting-response` holding at least one finding that a person answers: always
- `awaiting-response` where every finding it waits on is one the author patches: never
- `clean`: never
- any other state a run can hold: never

It SHALL take the state rather than the round number.
The round number would make it re-decide the routing that has already been decided and recorded, and a rule that recomputes an answer somebody has stored can disagree with the stored one.
Reading the state is reading that one decision's result, which is the property this requirement is for.

A `clean` round is never one a person reads even where it returned findings, because nothing is being asked of anybody and there is no decision to put in front of them.
A state that is not the outcome of a completed round is never one either, so a failed run and a cancelled run need no case of their own anywhere else.

The rule SHALL be decided by this one function and SHALL NOT be re-derived anywhere from the run's state and its counts.
Both facts a re-derivation would read are already reported, so the second reader is the cheap thing to write and the expensive thing to be wrong about: two places computing who a round is for can disagree, and the disagreement is silent, because a round nobody renders is not a round anybody notices was not rendered.

The rule SHALL be testable against its whole input space rather than sampled, and it SHALL reach no filesystem, no clock and no process.

#### Scenario: An escalated round is one a person reads

- **WHEN** a round leads to `escalated`, whether by a blocking finding whose own action asks for nothing or by the round cap
- **THEN** it is one a person reads

#### Scenario: A parked round holding a question is one a person reads

- **WHEN** a round parks and at least one of the findings it waits on is answered by a person
- **THEN** it is one a person reads

#### Scenario: A parked round the author answers alone is not one a person reads

- **WHEN** a round parks and every finding it waits on is one the author patches or rejects
- **THEN** it is not one a person reads

#### Scenario: A clean round is not one a person reads

- **WHEN** a round leads to `clean`, with or without findings that block nothing
- **THEN** it is not one a person reads

#### Scenario: The rule is decided over its whole input space

- **WHEN** the rule is exercised across every combination of classified findings and every state a run can hold
- **THEN** it agrees with the routing rule for every state that rule can reach, and answers that no person reads a round in any state that rule cannot reach
