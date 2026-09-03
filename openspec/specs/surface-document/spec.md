# surface-document Specification

## Purpose
TBD - created by archiving change a-page-that-answers-back. Update Purpose after archive.

## Requirements

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

### Requirement: The surface copies the identifiers rather than accepting them

Where a caller supplies a source document of its own, the surface SHALL take each decision's identifier, name, severity, audience and location from that source and SHALL NOT take them from the decision the caller wrote.

A source entry MAY carry a name, and the surface SHALL show a decision under the name its source gives it.
A decision whose source names none SHALL be shown under its identifier.

A page that renamed an identifier to something more readable hands the person back the search problem the page existed to remove.
Making the copy mechanical is what stops it being a matter of the caller's care, and a name is copied through the same door for the same reason: a caller free to title a decision is a caller free to make it unfindable.

An identifier is a key, and a person reading a page is not looking anything up.
The name is what the identifier is called in words, so that the person and the caller are still naming one thing.

#### Scenario: A caller cannot rename an identifier

- **WHEN** a document names a decision whose identifier, name, severity or audience disagrees with the source it cites
- **THEN** the rendered page shows what the source said, not what the decision said

#### Scenario: A decision is shown under the name its source gives it

- **WHEN** a source entry carries a name for one decision and none for another
- **THEN** the first is shown under that name and the second under its identifier

#### Scenario: A decision naming nothing in the source is refused

- **WHEN** a document holds a decision whose identifier appears in no source entry
- **THEN** `serve` fails and names the identifier

### Requirement: Every option carries the command it becomes

Each option SHALL carry the exact command that answering with it produces in the caller's terminal.
The surface SHALL return that command verbatim when the person chooses the option, SHALL NOT run it, and SHALL NOT render it.

Two copies of an answer drift, and the copy the person read is the one that is right.
Handing back a command rather than a paraphrase means the caller re-derives nothing, and the surface stays ignorant of what the command does.
The command is required because a caller that cannot say where an answer goes is offering a decision it cannot honour, and it is not rendered because the person deciding is not the one who runs it.

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

### Requirement: An option may be marked as the one the caller would take

A decision MAY mark one of its options as the one the caller would take.
At most one option per decision SHALL carry the mark, in either mode, and `serve` SHALL refuse a document marking two and SHALL name the decision.

The mark SHALL rank the option and SHALL NOT answer for the person.
A marked option SHALL arrive unselected, and a decision the person never touches SHALL come back with nothing chosen and no command, exactly as an unmarked one does.

A document that marks no option SHALL be served unchanged, and nothing about a decision without a mark SHALL be different from what it is today.

The caller has just read the finding and usually knows which option it would take.
A page that cannot carry that judgement makes the person re-derive it, and the caller that wanted to state it puts it in the prose or in the option order instead, which is a convention the surface never agreed to honour.

Two marks are a caller that has not decided, and a person who has to combine two recommendations has been handed the decision back with extra steps.

Marking is separated from selecting because a pre-selected option is an answer nobody gave, and by the time it reaches the terminal it is indistinguishable from one the person made.

#### Scenario: A decision names the option the caller would take

- **WHEN** a document marks one option of a decision
- **THEN** the page shows that option as the one the caller would take, and the other options carry no such mark

#### Scenario: A document marking two options is refused before serving

- **WHEN** a document marks two options of one decision
- **THEN** `serve` fails, names the decision, and starts no server

#### Scenario: Marking is not choosing

- **WHEN** a document marking an option is served and the person sends without touching that decision
- **THEN** `poll` returns that decision with nothing chosen and no command

#### Scenario: A document that marks nothing is unaffected

- **WHEN** a document whose decisions mark no option is served and answered
- **THEN** the page renders its options in the order the document listed them, and what `poll` returns is what it would have returned before the mark existed

### Requirement: A document may hold nothing to decide

A document SHALL be accepted with no decision at all, and the surface SHALL serve it as a page that is read rather than answered.

A document SHALL carry a title, and SHALL carry at least one of a decision, a source entry or a preamble.
A document carrying none of the three is refused, because the page would be a heading over nothing.

`serve` previously refused a document holding no decision, on the reasoning that there would be nothing for a person to answer.
The reasoning was sound about answering and wrong about serving: a caller with something for a person to read and nothing to ask them is the ordinary case, not a mistake, and refusing it forces the caller to invent a decision or to fall back to the terminal.

Every other rule about a decision SHALL be unchanged, except that a decision may carry no options.
That exception is stated as a requirement of its own below, and it is the same reasoning one step further in: a caller with something for a person to read and nothing to ask them about that one thing is also the ordinary case, and refusing it forces the caller to invent an option.
A document that holds decisions is validated in every other respect exactly as it is today, so a decision that does ask something is no different from what it is now.

#### Scenario: A document with nothing to decide is served

- **WHEN** a document carrying a title, a preamble and source entries but no decision is served
- **THEN** `serve` starts, prints the address, and the page renders what the document holds

#### Scenario: A document with nothing at all is refused

- **WHEN** a document carries a title and no decision, no source entry and no preamble
- **THEN** `serve` fails, says the page would hold nothing, and starts no server

#### Scenario: A document that asks something is validated as before

- **WHEN** a document holding a decision that carries options, one of them with no command, is served
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

### Requirement: A decision may state something rather than ask

A decision MAY carry no options.
The surface SHALL render such a decision as a card that states something and asks nothing, and SHALL render no control on it that a person could mistake for a choice.

Such a decision SHALL still carry the free-text field, on the same terms as every other decision: rendered by the surface, and beyond the caller's reach to label, size, require or remove.

This is a narrowing of the rule that a caller supplies the options, and is written down as one.
A caller may now write a decision with none, and what stops that becoming a caller reaching for a field it does not own is that the field is still the surface's and still the only one.

The source already shows a person everything that was found, in a list the surface copies rather than the caller writes.
What the list cannot carry is what was done about an entry and why, and a caller that has already resolved something has exactly that to say.
Its choices were to write it into the preamble, where a reader's eye goes last, or to give it an option in order to render a card at all.
An option on something already settled is a control that changes nothing, and a person reads a control as a question still owed an answer.

#### Scenario: A decision with no options renders as a card that asks nothing

- **WHEN** a document holding a decision with no options is served
- **THEN** the page renders it as a card, showing its question and its meaning and no option control of any kind

#### Scenario: A decision with no options still carries the field

- **WHEN** a decision with no options is rendered
- **THEN** the free-text field is there, and no field of the document removed, labelled or resized it

#### Scenario: An annotation on a decision that asked nothing comes back

- **WHEN** a person writes in the field of a decision that carries no options and sends
- **THEN** `poll` returns that decision with no option chosen, no command, and the words carried verbatim

#### Scenario: A page of cards that ask nothing offers no controls

- **WHEN** a document whose every decision carries no options is served
- **THEN** the page renders every one as a card and offers no option control anywhere on it

### Requirement: A decision that asks nothing states no mode

A decision carrying no options SHALL NOT state a mode, and `serve` SHALL refuse a document that gives one a mode and SHALL name the decision.

A mode names how a person chooses among options.
Where there are none there is nothing to choose among, so a stated mode is a contradiction rather than a harmless extra field.

Refusing is chosen over ignoring because the caller is a session assembling a document, and a field silently ignored is a field it keeps emitting and eventually relies on.

#### Scenario: A mode on a decision with no options is refused before serving

- **WHEN** a document holds a decision that carries no options and states a mode
- **THEN** `serve` fails, names the decision, and starts no server

#### Scenario: A decision with options still requires its mode

- **WHEN** a document holds a decision that carries options and states no mode
- **THEN** `serve` fails, names the decision, and starts no server

### Requirement: A language the page cannot highlight is refused, and no language is not

Where a decision's evidence names a language, `serve` SHALL refuse a document naming one the page cannot highlight, SHALL name the decision and the language, and SHALL name the languages it can.

Evidence naming no language SHALL be served and SHALL render unhighlighted.

These are two different statements by the caller and they earn different treatment.
Naming no language is a caller that wants no highlighting, which is a legitimate thing to want.
Naming a language the page ignores is a caller that believes its evidence will be highlighted and is wrong, and the person reading a grey block is the one who pays for the belief.

The caller is at a terminal where the error is one edit away, and the person never sees the degraded page.
The refusal names the set it knows so that the way forward is in the message.

#### Scenario: An unknown evidence language is refused before serving

- **WHEN** a document names an evidence language the page cannot highlight
- **THEN** `serve` fails, names the decision and the language, names the languages it can highlight, and starts no server

#### Scenario: Evidence with no language is served

- **WHEN** a document carries evidence that names no language
- **THEN** the document is served and that evidence renders as plain text
