## MODIFIED Requirements

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

## ADDED Requirements

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
