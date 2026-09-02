## ADDED Requirements

### Requirement: A page with nothing to answer says so and offers neither control

Where the document holds no decision, the page SHALL offer no way to answer and no way to send, and SHALL say in words that there is nothing on it to answer and that nobody is waiting on it.

A send control with nothing to carry delivers an empty answer to a caller that may not be listening, and a person who meets one reasonably believes something is expected of them.
Saying so is what stops the absence reading as a page that failed to load its controls.

The page SHALL otherwise be the page: the same frame, the same palette, the same layout, and the source rendered as it is on any other page.
A page that is read rather than answered is not a lesser page, and the caller supplies no more markup for it than for any other.

#### Scenario: A page with nothing to answer offers neither control

- **WHEN** a document holding no decision is served and the rendered page's controls are enumerated
- **THEN** there is nothing to answer and nothing to send

#### Scenario: The page says nobody is waiting

- **WHEN** a person opens a page holding no decision
- **THEN** it says there is nothing to answer and that nobody is waiting on it

#### Scenario: The frame is the same frame

- **WHEN** a page holding no decision is rendered beside one holding decisions
- **THEN** both carry the same frame, palette and layout, and both render the source the same way

### Requirement: The page opens with the source at a glance

The page SHALL render the entries of the document's source as a short list, before the decisions.

Each entry SHALL be shown under its name, with its severity, its audience and its location.
An entry a decision covers SHALL take the person to that decision.

A person meeting a page needs to know how much there is before they start reading any of it.
The list is short by construction because its content is copied rather than written, so it costs the caller nothing and it is the one part of the page that cannot disagree with the caller's terminal.

#### Scenario: The source is read before the decisions

- **WHEN** a page is rendered for a document citing source entries and holding decisions
- **THEN** the entries are listed first, each under its name with its severity, audience and location, and the decisions follow

#### Scenario: An entry leads to the decision that covers it

- **WHEN** a person follows an entry that a decision covers
- **THEN** they arrive at that decision

#### Scenario: A document citing no source lists nothing

- **WHEN** a page is rendered for a document that cites no source
- **THEN** no list is shown, and the page is what it is today

## MODIFIED Requirements

### Requirement: The page holds two interactions and nothing else

The page SHALL offer a way to answer a decision and a way to send, and SHALL offer neither where the document holds no decision.
It SHALL NOT offer annotation, a whiteboard, an export, a publish, a share, or a control whose effect the caller never sees.

Every control that is not one of the two is a thing to build, a thing to explain and a thing a person can spend attention on instead of the decision.
The surface exists because a person had one question to answer, or because they had something to read and no question at all, and in the second case both controls are controls whose effect the caller never sees.

Reading the source is not a third interaction.
An entry that takes the person to the decision covering it moves them within the page they are already on, and carries nothing anywhere.

#### Scenario: The surface offers no third interaction

- **WHEN** the rendered page's controls are enumerated
- **THEN** each is part of answering a decision or part of sending, and none annotates, exports, publishes or shares

#### Scenario: A page with nothing to answer offers neither of the two

- **WHEN** the rendered controls of a page whose document holds no decision are enumerated
- **THEN** there are none
