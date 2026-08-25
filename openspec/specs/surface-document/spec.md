# surface-document Specification

## Purpose
TBD - created by archiving change a-page-that-answers-back. Update Purpose after archive.
## Requirements
### Requirement: A caller writes a document, not a page

The surface SHALL take one JSON document and render the page itself.
A caller SHALL supply only what needs judgement: the question, what it means in plain English, the options, and an optional fragment.

A caller that hand-writes a whole page pays for the frame, the palette, the layout and the diff rendering on every exchange, in the context that is about to do the work the exchange was for.
The frame is the same every time, so the surface owns it.

#### Scenario: A document renders without the caller writing markup

- **WHEN** a document holding one decision and no fragment is served
- **THEN** the page renders that decision as a card, and the caller supplied no HTML

#### Scenario: A document that does not validate is refused before serving

- **WHEN** a document is missing a required field, names a mode that does not exist, or carries an option with no command
- **THEN** `serve` fails, names the field and the decision it belongs to, and starts no server

### Requirement: The surface copies the identifiers rather than accepting them

Where a caller supplies a source document of its own, the surface SHALL take each decision's identifier, severity, audience and location from that source and SHALL NOT take them from the decision the caller wrote.

A page that renamed an identifier to something more readable hands the person back the search problem the page existed to remove.
Making the copy mechanical is what stops it being a matter of the caller's care.

#### Scenario: A caller cannot rename an identifier

- **WHEN** a document names a decision whose identifier, severity or audience disagrees with the source it cites
- **THEN** the rendered page shows what the source said, not what the decision said

#### Scenario: A decision naming nothing in the source is refused

- **WHEN** a document holds a decision whose identifier appears in no source entry
- **THEN** `serve` fails and names the identifier

### Requirement: Every option carries the command it becomes

Each option SHALL carry the exact command that answering with it produces in the caller's terminal.
The surface SHALL return that command verbatim when the person chooses the option, and SHALL NOT run it.

Two copies of an answer drift, and the copy the person read is the one that is right.
Handing back a command rather than a paraphrase means the caller re-derives nothing, and the surface stays ignorant of what the command does.

#### Scenario: The command is what comes back

- **WHEN** a person chooses an option and sends
- **THEN** `poll` prints that option's command exactly as the document stated it

#### Scenario: The surface runs nothing

- **WHEN** a person chooses any option and sends
- **THEN** no process is started by the surface, and nothing outside its own session file is written

### Requirement: Every decision carries a way to answer in the person's own words

A decision SHALL always render a free-text field beside its options, whether the mode is single-choice or multiple-choice.
A document SHALL NOT be able to remove it.

The options are the caller's reading of the question and not the question itself.
Three buttons make a question feel closed when it is not, and the field is what keeps the person able to say the thing nobody offered.

#### Scenario: The free-text field cannot be suppressed

- **WHEN** a decision is rendered in either mode
- **THEN** a free-text field is present, and no field of the document removed it

#### Scenario: Words alone are a complete answer

- **WHEN** a person writes in the field, chooses no option, and sends
- **THEN** `poll` returns that decision with no option chosen and the words carried verbatim

### Requirement: A supplied fragment is neutralised before it is embedded

The surface SHALL reduce any HTML fragment a document carries to a fixed allowlist of elements and attributes before embedding it, and SHALL drop anything outside that allowlist.
A document MAY carry one such fragment per decision, for an illustration the prose cannot carry.

The fragment is output from the same machine that composed the document, placed in a page that machine also composed.
This is the shape `internal/hygiene` already exists for on the prompt side, against a different threat.

#### Scenario: A fragment outside the allowlist is stripped

- **WHEN** a fragment carries a script, an event-handler attribute, a remote reference or a form control
- **THEN** the page renders the fragment with those removed, and serves without failing

#### Scenario: An illustration survives neutralising

- **WHEN** a fragment carries inline SVG using only allowed elements and attributes
- **THEN** the page renders it unchanged
