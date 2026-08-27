## ADDED Requirements

### Requirement: The skill says how an illustration reaches the page

The skill SHALL say that a decision may carry one illustration, and that it travels beside the evidence rather than inside it.

It SHALL say an illustration takes one of two forms, and which each is for.
A drawing the session writes as inline SVG is for a workflow, because the surface embeds no diagram engine and renders what the caller supplied.
Markup the document's stylesheet paints is for a panel a person would recognise, because a person who already knows the screen is being asked to recognise it rather than to read a description of it.

It SHALL say the document carries one stylesheet for all of its illustrations rather than one per decision, so a session writing a round about one product states its look once.

It SHALL say that both the markup and the stylesheet are reduced before they are embedded, and name what does not survive: script, event handlers, remote references and form controls from the markup; imports, font faces, properties outside the allowlist and values naming a resource from the stylesheet.
It SHALL say the surface scopes the stylesheet itself, so a session does not scope its own selectors and does not expect a rule aimed at the page to reach it.

It SHALL say this beside the bar's lines about showing and about diagrams, so that a session reads when and how in one place.

It SHALL NOT restate the bar's judgement about when either is warranted.
The bar owns that, each line is stated once, and a second statement beside the mechanism is the copy that gets softened.

#### Scenario: The skill names the illustration as something a decision may carry

- **WHEN** the skill is rendered
- **THEN** it says a decision may carry one illustration, beside its evidence rather than within it

#### Scenario: The skill says what forms an illustration takes

- **WHEN** the skill describes the illustration
- **THEN** it names both a drawing the session writes and markup a stylesheet paints, says which is for what, and does not name or imply a diagram engine the surface would run

#### Scenario: The skill says the look is stated once per document

- **WHEN** the skill describes the stylesheet
- **THEN** it says one stylesheet paints every illustration in the document, rather than one per decision

#### Scenario: The skill says what is stripped before either is embedded

- **WHEN** the skill describes the illustration
- **THEN** it says both the markup and the stylesheet are reduced to an allowlist, and names what does not survive each, so that a session does not write what will be dropped

#### Scenario: The skill says the surface does the scoping

- **WHEN** the skill describes the stylesheet
- **THEN** it says the surface scopes every selector to the illustration, so a session neither scopes its own nor expects one to reach the page

#### Scenario: How to illustrate is stated beside when to

- **WHEN** the skill is rendered
- **THEN** the mechanism appears beside the bar's lines about showing and about diagrams, and the bar's judgement about when each is warranted appears once and only there
