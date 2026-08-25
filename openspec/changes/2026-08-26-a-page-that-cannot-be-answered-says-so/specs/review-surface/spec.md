## MODIFIED Requirements

### Requirement: What poll prints is data, and it says whether the exchange is over

`poll` SHALL print, for each decision, the identifier, what was chosen, the words the person added, and the command the chosen option carries.
It SHALL print any message the person wrote to the caller.
It SHALL state whether the person ended the exchange or expects the caller to answer and re-render.

A poll SHALL NOT reopen an exchange the person ended.
Writing a new document is the one act that does, and it is the caller's.

Two sends that a caller cannot tell apart make a caller either stop too early or wait for a person who has finished.
The two rules are stated together because they are the pair that is easy to confuse: reading an ended exchange leaves it ended, and asking a new question does not.

#### Scenario: An answer round-trips

- **WHEN** a person chooses in two decisions, writes a message, and sends without ending
- **THEN** `poll` prints both answers with their commands, the message, and that the exchange is not over

#### Scenario: Ending is distinguishable

- **WHEN** a person sends and ends
- **THEN** `poll` prints that the exchange is over, and a later `poll` does not reopen it

#### Scenario: Delivery happens once

- **WHEN** an answer has been delivered to one `poll`
- **THEN** a second `poll` does not deliver it again

#### Scenario: A poll over an ended exchange leaves it ended

- **WHEN** the person has ended the exchange and `poll` is run again
- **THEN** it reports that nobody answered rather than reopening the exchange, and the page stays closed

## ADDED Requirements

### Requirement: A caller who writes a new document is asking again

Writing a new document to a surface whose exchange has ended SHALL reopen that exchange, and the surface SHALL then accept an answer.
`serve` SHALL NOT return an address whose only possible reply to an answer is a refusal.

A caller reaches this path in one situation, and it is the situation where it matters: the person ended the exchange without answering, so the caller asks again.
Handing back the address of an exchange that can only refuse sends that person to a page that was never going to work, and they find out after deciding a second time.
Reopening is the smaller of the two available rules and it changes nothing while an exchange is open, where the exchange is not over and there is nothing to reopen.
It also needs no new command: recovering is serving a document, which `serve` already does, so the five commands stay five.

#### Scenario: Serving again onto an ended exchange yields an answerable page

- **WHEN** the person ends an exchange without answering and the caller serves a document to that surface again
- **THEN** the address returned holds a page that accepts an answer, and `poll` delivers what the person sends from it

#### Scenario: Serving again while the exchange is open is unchanged

- **WHEN** the caller serves a new document to a surface whose exchange is open
- **THEN** the page takes up the new content and the exchange is as answerable as it was
