## MODIFIED Requirements

### Requirement: Sending is two acts, and the page says which is which

The composer SHALL offer one send that expects the caller to answer and re-render, and one that ends the exchange.
The page SHALL state, in words, what each does.

Where a decision holds selected options or written words that have not been committed, the send that ends SHALL state which decisions those are and that what they hold will not reach the caller, and SHALL end only once the person confirms.
An uncommitted answer SHALL NOT be delivered by either send.

A person who cannot tell "I want to discuss this" from "I am finished" either ends too early or waits for a caller that has stopped listening.
Ending is the one act on this page that cannot be taken back, and it is reached by a control sitting beside an answer the person believes they have given.
Delivering that answer instead would be the control that sends on every click, arriving late: somebody who selected an option to read it against the others would have their reading delivered as their decision.
So the answer stays undelivered and the silence goes, because a person who is told what they are about to discard can still choose to discard it.

#### Scenario: Sending without ending keeps the page open

- **WHEN** a person sends without ending
- **THEN** the page stays open, and the caller may re-render it

#### Scenario: Ending closes the exchange

- **WHEN** a person sends and ends
- **THEN** the page says the exchange is over and accepts no further answer

#### Scenario: Finishing over an uncommitted answer asks first

- **WHEN** a person selects options in a decision, does not commit it, and sends with the control that ends
- **THEN** the page names that decision, says what it holds will not reach the caller, and nothing is delivered until the person confirms

#### Scenario: Confirming carries nothing that was not committed

- **WHEN** a person confirms that finish
- **THEN** the exchange ends, and `poll` reports no answer for the decision that was never committed

#### Scenario: Declining leaves the exchange open

- **WHEN** a person declines that finish
- **THEN** the exchange is not over, nothing is delivered, and the selections the decision held are still there to commit

### Requirement: The page updates when the caller re-renders it

Where the caller writes a new document for a surface, the page SHALL take up the new content without the person reloading it.
An answer the person had committed but not sent SHALL NOT be silently carried into content that no longer holds that decision.

Where the exchange had ended, a new document SHALL make the page answerable again, because a caller writing one is asking again.

The exchange that this exists for is one where the caller answers and re-renders.
A person who has to reload to see the reply does not know the reply arrived.
A person sent back to a page to answer a question the caller has just asked again is owed a page that can take the answer.

#### Scenario: A re-render reaches an open page

- **WHEN** the caller writes a new document for the surface being read
- **THEN** the open page shows the new content without being reloaded

#### Scenario: A new document reopens an ended page

- **WHEN** the caller writes a new document for a surface whose exchange the person ended
- **THEN** the page becomes answerable, and an answer sent from it is accepted and delivered

## ADDED Requirements

### Requirement: A page that cannot be answered says so, however it was loaded

Where the exchange is over, the page SHALL render its closed state and SHALL NOT offer to take an answer, whether it was open when the exchange ended, loaded afresh afterwards, or reloaded.
A live stream opened onto an exchange that is already over SHALL be told so on connection rather than on the next change.

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
