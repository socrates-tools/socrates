## ADDED Requirements

### Requirement: The skill says how an illustration reaches the page

The skill SHALL say that a decision may carry one illustration, and that it travels as a fragment of its own beside the evidence rather than inside it.

It SHALL say the fragment is inline SVG the session writes, because the surface embeds no diagram engine and renders what the caller supplied.

It SHALL say the fragment is reduced to a fixed allowlist before it is embedded, and that script, event handlers, remote references and form controls do not survive it.

It SHALL say this beside the bar's line about when to draw a diagram, so that a session reads when and how in one place.

It SHALL NOT restate the bar's judgement about when a diagram is warranted.
The bar owns that, it is stated once, and a second statement of it beside the mechanism is the copy that gets softened.

#### Scenario: The skill names the fragment as the way an illustration reaches the page

- **WHEN** the skill is rendered
- **THEN** it says a decision may carry one fragment for an illustration, beside its evidence rather than within it

#### Scenario: The skill says what form the illustration takes

- **WHEN** the skill describes the fragment
- **THEN** it says the session writes inline SVG itself, and does not name or imply a diagram engine the surface would run

#### Scenario: The skill says what is stripped before the fragment is embedded

- **WHEN** the skill describes the fragment
- **THEN** it says the fragment is reduced to a fixed allowlist, so that a session does not write what will be dropped

#### Scenario: How to draw one is stated beside when to draw one

- **WHEN** the skill is rendered
- **THEN** the mechanism appears beside the bar's line about diagrams, and the bar's judgement about when one is warranted appears once and only there
