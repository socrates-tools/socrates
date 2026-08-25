## MODIFIED Requirements

### Requirement: socrates states the fact and does nothing further about the page

socrates SHALL NOT render, publish, store, name or verify a page.
It states that a round is one a person reads, and stops.

No command SHALL be added to `socrates`, no flag SHALL be added to `socrates`, and nothing about a page SHALL be written into a run's record.
Recording a link would be a field, a flag to write it and a thing to be stale, in exchange for a link the session has just printed.

socrates SHALL gain no network reach, no dependency and no build tag from this.
What is measured SHALL be `socrates` itself: the binary, and the packages it reaches.
A second binary in the same module is not socrates gaining reach, and measuring the module rather than the tool would answer a question nobody asked.
The rule is that the reviewing tool reaches no network, not that this repository holds no networking code.

#### Scenario: The surface does not grow

- **WHEN** the commands and flags the `socrates` binary offers are compared with what it offered before
- **THEN** none has been added, and no run record carries anything about a page

#### Scenario: socrates reaches nothing new

- **WHEN** what `socrates` itself reaches is examined
- **THEN** it reaches no network, and the packages behind a build tag are the ones that were behind one before

#### Scenario: A second binary in the module is not socrates reaching further

- **WHEN** another binary in the same module imports the network
- **THEN** what `socrates` reaches is unchanged, and the check measures the tool rather than the repository
