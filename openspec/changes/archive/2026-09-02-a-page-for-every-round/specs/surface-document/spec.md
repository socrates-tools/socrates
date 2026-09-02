## ADDED Requirements

### Requirement: A document may hold nothing to decide

A document SHALL be accepted with no decision at all, and the surface SHALL serve it as a page that is read rather than answered.

A document SHALL carry a title, and SHALL carry at least one of a decision, a source entry or a preamble.
A document carrying none of the three is refused, because the page would be a heading over nothing.

`serve` previously refused a document holding no decision, on the reasoning that there would be nothing for a person to answer.
The reasoning was sound about answering and wrong about serving: a caller with something for a person to read and nothing to ask them is the ordinary case, not a mistake, and refusing it forces the caller to invent a decision or to fall back to the terminal.

Every other rule about a decision SHALL be unchanged.
A document that holds decisions is validated exactly as it is today, so nothing about a page that asks something is different from what it is now.

#### Scenario: A document with nothing to decide is served

- **WHEN** a document carrying a title, a preamble and source entries but no decision is served
- **THEN** `serve` starts, prints the address, and the page renders what the document holds

#### Scenario: A document with nothing at all is refused

- **WHEN** a document carries a title and no decision, no source entry and no preamble
- **THEN** `serve` fails, says the page would hold nothing, and starts no server

#### Scenario: A document that asks something is validated as before

- **WHEN** a document holding a decision with no option, or with an option carrying no command, is served
- **THEN** `serve` fails and names the decision, exactly as it does today

### Requirement: The source is what the page shows of the whole

The surface SHALL show every entry of the source a caller cites, whether or not a decision cites that entry.

An entry SHALL be shown under the name its source gives it, with the severity, audience and location that source carries.
An entry a decision covers SHALL lead the person to that decision rather than being shown twice as though it were a second thing.

The source is the only part of a document that is copied rather than composed, and it is the part that lets a person see everything that was found without reading it all.
A caller that had to compose that list in prose would be retyping what it was already handed, and a retyped list is one that can disagree with what the caller's own terminal says.

#### Scenario: An entry no decision cites is still shown

- **WHEN** a document cites six source entries and holds one decision
- **THEN** the page shows all six under their names, severities and audiences, and the one the decision covers leads to it

#### Scenario: The list is not the caller's prose

- **WHEN** a document holds source entries
- **THEN** what the page lists comes from those entries rather than from the preamble the caller wrote
