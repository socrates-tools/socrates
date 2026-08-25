## MODIFIED Requirements

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
