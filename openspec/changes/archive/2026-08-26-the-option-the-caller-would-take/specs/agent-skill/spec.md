## ADDED Requirements

### Requirement: The skill says how a recommendation reaches the page

The skill SHALL say that a decision may mark one option as the one the session would take, and SHALL show the field where it shows the document's shape.

It SHALL say that at most one option per decision carries the mark, and that a document marking two is refused before anything is served.

It SHALL say what the mark does and what it does not do: the page reads the marked option first and says in words that it is the one the session would take, and it selects nothing.
The person's answer is still the person's, and a session that has marked an option has recommended rather than decided.

It SHALL NOT restate the bar's judgement about when to recommend.
The bar owns that, each line is stated once, and a second statement beside the mechanism is the copy that gets softened.

#### Scenario: The skill names the mark where it names the document

- **WHEN** the skill shows the shape of a decision
- **THEN** the option the session would take is marked there, in the same shape a session copies

#### Scenario: The skill says one option carries it

- **WHEN** the skill describes the mark
- **THEN** it says at most one option per decision carries it and that a document marking two is refused before serving

#### Scenario: The skill says marking is not answering

- **WHEN** the skill describes the mark
- **THEN** it says the page reads the marked option first and states it in words, that nothing arrives selected, and that the decision remains the person's

#### Scenario: When to recommend is stated once

- **WHEN** the skill is rendered
- **THEN** the judgement about when to recommend appears in the bar and nowhere else
