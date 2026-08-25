## ADDED Requirements

### Requirement: The page holds two interactions and nothing else

The page SHALL offer a way to answer a decision and a way to send.
It SHALL NOT offer annotation, a whiteboard, an export, a publish, a share, or a control whose effect the caller never sees.

Every control that is not one of the two is a thing to build, a thing to explain and a thing a person can spend attention on instead of the decision.
The surface exists because a person had one question to answer.

#### Scenario: The surface offers no third interaction

- **WHEN** the rendered page's controls are enumerated
- **THEN** each is part of answering a decision or part of sending, and none annotates, exports, publishes or shares

### Requirement: A decision is answered by choosing, and the mode says how many

A decision SHALL render as single-choice or multiple-choice as its document states.
Single-choice SHALL accept one option, and multiple-choice SHALL accept any number including none.

Choosing SHALL change only what the page holds, and SHALL send nothing.
A separate control per decision SHALL commit that decision's answer.

A control that sent on every click sends an answer the person was still deciding, and takes back nothing when they change their mind.

#### Scenario: Choosing sends nothing

- **WHEN** a person selects, deselects and reselects options in a decision
- **THEN** nothing is delivered, and `poll` is still waiting

#### Scenario: Multiple choice carries every tick

- **WHEN** a person ticks three options in a multiple-choice decision and commits it
- **THEN** all three are carried, in the order the document listed them

#### Scenario: Committing again replaces the earlier answer

- **WHEN** a person commits a decision, changes it, and commits it again before sending
- **THEN** one answer for that decision is delivered, and it is the later one

### Requirement: Selected is shown separately from committed

A decision whose answer has been committed SHALL render differently from one the person has merely clicked in, and SHALL show the command that answer carries.

A person who cannot tell what the caller will receive from what they have merely touched does not know what they have said.
Showing the command is also what makes the page and the terminal matchable in the direction the person needs.

#### Scenario: A committed answer is distinguishable

- **WHEN** one decision is committed and another only has options selected
- **THEN** the two render differently, and only the committed one shows its command

#### Scenario: A committed answer can be taken back

- **WHEN** a person reopens a committed decision before sending
- **THEN** it returns to the uncommitted state and its previous answer is not delivered

### Requirement: Sending is two acts, and the page says which is which

The composer SHALL offer one send that expects the caller to answer and re-render, and one that ends the exchange.
The page SHALL state, in words, what each does.

A person who cannot tell "I want to discuss this" from "I am finished" either ends too early or waits for a caller that has stopped listening.

#### Scenario: Sending without ending keeps the page open

- **WHEN** a person sends without ending
- **THEN** the page stays open, and the caller may re-render it

#### Scenario: Ending closes the exchange

- **WHEN** a person sends and ends
- **THEN** the page says the exchange is over and accepts no further answer

### Requirement: The page renders its own evidence

The page SHALL render highlighted source, a diff in a split or unified layout, and an inline illustration the caller supplied, without the caller writing the machinery for any of them.

Every asset needed SHALL be embedded in the binary.
No diagram engine SHALL be embedded, because an illustration is written by the caller and a layout engine surrenders position and emphasis to a machine.

#### Scenario: A diff renders from a patch

- **WHEN** a decision carries a patch as evidence
- **THEN** the page renders it as a diff with both sides and per-line changes visible

#### Scenario: The narrow case does not shrink the evidence to nothing

- **WHEN** the page is viewed at a narrow width
- **THEN** the diff and any illustration remain legible, scrolling within their own bounds rather than the page scrolling sideways

#### Scenario: Rendering needs no network

- **WHEN** the page is loaded with the machine's outbound network unavailable
- **THEN** the highlighting, the diff, the fonts and the illustration all render

### Requirement: The page updates when the caller re-renders it

Where the caller writes a new document for an open surface, the page SHALL take up the new content without the person reloading it.
An answer the person had committed but not sent SHALL NOT be silently carried into content that no longer holds that decision.

The exchange that this exists for is one where the caller answers and re-renders.
A person who has to reload to see the reply does not know the reply arrived.

#### Scenario: A re-render reaches an open page

- **WHEN** the caller writes a new document for the surface being read
- **THEN** the open page shows the new content without being reloaded

### Requirement: The presentation is held to a stated floor

No interface text SHALL render below a stated minimum size, and the page SHALL paint its own background and text colours rather than inheriting them.
The page SHALL follow the reader's light or dark preference.
Every control SHALL be operable by keyboard and SHALL carry an accessible name.

A floor that is written down is a floor that survives somebody being in a hurry, and a page that inherits its colours from whatever opened it is a page nobody has actually seen.

#### Scenario: Nothing renders below the floor

- **WHEN** the rendered page's text is measured
- **THEN** no interface text is below the stated minimum

#### Scenario: Both preferences are painted

- **WHEN** the page is rendered under a light preference and under a dark one
- **THEN** each paints its own background and text, and neither leaves a colour undefined

#### Scenario: A decision can be answered without a pointer

- **WHEN** a person moves through the page by keyboard alone
- **THEN** every option, field and send can be reached, named and operated
