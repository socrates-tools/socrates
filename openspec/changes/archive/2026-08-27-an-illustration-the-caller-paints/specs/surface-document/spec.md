## MODIFIED Requirements

### Requirement: A caller writes a document, not a page

The surface SHALL take one JSON document and render the page itself.
A caller SHALL supply only what needs judgement: the question, what it means in plain English, the options, an optional illustration, and an optional stylesheet that paints the illustrations.

A caller that hand-writes a whole page pays for the frame, the palette, the layout and the diff rendering on every exchange, in the context that is about to do the work the exchange was for.
The frame is the same every time, so the surface owns it.

The frame is not negotiable and the illustration is not part of it.
The surface owns the page, its palette, its layout, its evidence rendering and every control a person operates, and a caller cannot reach any of them.
What a caller may paint is the inside of one box per decision, which shows the person something rather than asking them anything.
This is a narrowing of the sentence above and is written down as one: a caller may now paint, and what stops that becoming a caller writing a page is that the box is bounded, its stylesheet is reduced, and every selector in it is scoped by the surface rather than by the caller.

#### Scenario: A document renders without the caller writing markup

- **WHEN** a document holding one decision and no illustration is served
- **THEN** the page renders that decision as a card, and the caller supplied no HTML

#### Scenario: A caller cannot paint the frame

- **WHEN** a document carries a stylesheet that declares rules for the page's own frame, palette or controls
- **THEN** the page renders its frame, palette and controls exactly as it would with no stylesheet at all

#### Scenario: A document that does not validate is refused before serving

- **WHEN** a document is missing a required field, names a mode that does not exist, carries an option with no command, or carries a stylesheet larger than the stated size
- **THEN** `serve` fails, names the field and the decision it belongs to, and starts no server

### Requirement: A supplied fragment is neutralised before it is embedded

The surface SHALL reduce any HTML fragment a document carries to a fixed allowlist of elements and attributes before embedding it, and SHALL drop anything outside that allowlist.
A document MAY carry one such fragment per decision, for an illustration the prose cannot carry.

The fragment is output from the same machine that composed the document, placed in a page that machine also composed.
This is the shape `internal/hygiene` already exists for on the prompt side, against a different threat.

The fragment SHALL carry its markup as a named field rather than being a bare string, so that the document says what an illustration is made of rather than assuming it is made of one thing.
An illustration is markup the document's stylesheet paints, or it is a drawing that paints itself, and a field that can only be a string can only be the second.

The names a fragment may carry SHALL be reduced along with its elements and attributes, so that a fragment cannot claim a name the page uses for its own chrome.
An unconstrained name is inert while nothing paints from names, and it is a fragment reaching the page's own controls the moment something does.

#### Scenario: A fragment outside the allowlist is stripped

- **WHEN** a fragment carries a script, an event-handler attribute, a remote reference or a form control
- **THEN** the page renders the fragment with those removed, and serves without failing

#### Scenario: A fragment cannot claim the page's own names

- **WHEN** a fragment carries a name the page uses for its own chrome or for one of its decisions
- **THEN** the page renders the fragment without it, and nothing the page paints or addresses is affected

#### Scenario: An illustration survives neutralising

- **WHEN** a fragment carries inline SVG using only allowed elements and attributes
- **THEN** the page renders it unchanged

#### Scenario: An illustration the stylesheet paints survives neutralising

- **WHEN** a fragment carries markup using only allowed elements, attributes and names
- **THEN** the page renders it unchanged and the document's stylesheet paints it
