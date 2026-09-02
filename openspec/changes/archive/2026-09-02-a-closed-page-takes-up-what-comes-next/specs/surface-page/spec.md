## MODIFIED Requirements

### Requirement: The page updates when the caller re-renders it

Where the caller writes a new document for a surface, the page SHALL take up the new content without the person reloading it.
An answer the person had committed but not sent SHALL NOT be silently carried into content that no longer holds that decision.

Where the exchange had ended, a new document SHALL make the page answerable again, because a caller writing one is asking again.
It SHALL do so on the page the person is already looking at, without them reloading it, on the same terms as a re-render onto an exchange that never ended.

The end of an exchange SHALL NOT be the end of the page's attention.
A page told the exchange is over stops offering to answer and keeps listening, because the thing it is waiting for is the caller asking again, and that is the one event it would otherwise be unable to hear.

The exchange that this exists for is one where the caller answers and re-renders.
A person who has to reload to see the reply does not know the reply arrived.
That is as true after they have sent and finished as before: a page reading "the exchange is over" is a page nobody reloads, so a reply written to it reaches nobody, and the person is left waiting on a surface that already holds what they are waiting for.
A person sent back to a page to answer a question the caller has just asked again is owed a page that can take the answer.

Reloading SHALL remain sufficient, and SHALL NOT be what the requirement rests on.
A test that reloads before asserting proves the weaker of the two claims and reports it as the stronger.

#### Scenario: A re-render reaches an open page

- **WHEN** the caller writes a new document for the surface being read
- **THEN** the open page shows the new content without being reloaded

#### Scenario: A new document reopens an ended page

- **WHEN** the caller writes a new document for a surface whose exchange the person ended, and the person does not reload
- **THEN** the page they are looking at shows the new content, becomes answerable, and an answer sent from it is accepted and delivered

#### Scenario: A reloaded page still reopens

- **WHEN** the caller writes a new document for a surface whose exchange the person ended, and the person reloads
- **THEN** the page shows the new content and is answerable, as it is without the reload

#### Scenario: The page stops saying the exchange is over

- **WHEN** a new document reaches a page that was showing the closed state
- **THEN** it no longer says the exchange is over, and every control the new content carries is operable

### Requirement: A page that cannot be answered says so, however it was loaded

Where the exchange is over, the page SHALL render its closed state and SHALL NOT offer to take an answer, whether it was open when the exchange ended, loaded afresh afterwards, or reloaded.

A live stream SHALL be told the state of the exchange it is opened onto when it connects, rather than only when something next changes.
Where that exchange is already over it SHALL be told so on connection.
Where it is not, the stream SHALL state the revision the exchange stands at on connection.

A stream that says nothing until the next change cannot be caught up.
A client that connects, or reconnects, after a revision it did not see waits for a change that has already happened, and the page stays on content the caller has replaced with no way to discover it has.
Stating the revision on connection is what makes a stream something a client can rejoin rather than something it must not lose.

Whether an answer can be accepted is decided by one piece of state, and until now that state was consulted only where an answer was already being refused.
A page that renders every control live over an ended exchange is a page that invites an answer it will reject, and the person finds out by composing one and losing it.
Refusing the answer is not the same as saying it cannot be given, and only the second is any use before the fact.

#### Scenario: A reloaded page renders closed

- **WHEN** a person ends the exchange and then reloads the page
- **THEN** the page says the exchange is over and offers no control that would send an answer

#### Scenario: A page opened for the first time after the end renders closed

- **WHEN** a page is loaded for the first time onto an exchange that has already ended
- **THEN** it says the exchange is over rather than rendering as though it could be answered

#### Scenario: A stream opened after the end is told at once

- **WHEN** a live stream is opened onto an exchange that is already over, and nothing further changes
- **THEN** it reports the exchange closed without waiting for a change

#### Scenario: A stream is told where the exchange stands when it connects

- **WHEN** a live stream is opened onto an exchange that is not over, and nothing further changes
- **THEN** it states the revision the exchange stands at without waiting for a change

#### Scenario: A stream that missed a revision catches up on connecting again

- **WHEN** a document is written for a surface while nothing is streaming, and a stream is then opened
- **THEN** it is told the current revision rather than waiting for a change that has already happened
