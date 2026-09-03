# surface-page Specification

## Purpose
TBD - created by archiving change a-page-that-answers-back. Update Purpose after archive.

## Requirements

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

### Requirement: A decision is answered by choosing, and the mode says how many

A decision SHALL render as single-choice or multiple-choice as its document states.
Single-choice SHALL accept one option, and multiple-choice SHALL accept any number including none.

Choosing SHALL change only what the page holds, and SHALL send nothing.
A separate control per decision SHALL commit that decision's answer.

Where a decision marks the option the caller would take, the page SHALL render that option first and SHALL leave the remaining options in the order the document listed them.
No option SHALL arrive selected, marked or not.
What an answer carries SHALL be keyed to the order the document listed the options in, so where the page shows an option SHALL make no difference to what the caller is handed back.

A control that sent on every click sends an answer the person was still deciding, and takes back nothing when they change their mind.

A recommendation the person reaches after the others is one they read after they have decided, so placement is where the mark does its work.
Order on the page and order in the answer are held apart because the page is arranged for the person reading it and the answer is keyed for the caller receiving it, and one is free to change without disturbing the other.

#### Scenario: Choosing sends nothing

- **WHEN** a person selects, deselects and reselects options in a decision
- **THEN** nothing is delivered, and `poll` is still waiting

#### Scenario: Multiple choice carries every tick

- **WHEN** a person ticks three options in a multiple-choice decision and commits it
- **THEN** all three are carried, in the order the document listed them

#### Scenario: Committing again replaces the earlier answer

- **WHEN** a person commits a decision, changes it, and commits it again before sending
- **THEN** one answer for that decision is delivered, and it is the later one

#### Scenario: The option the caller would take is read first

- **WHEN** a decision marking its third option is rendered
- **THEN** that option is the first the person meets, and the rest follow in the order the document listed them

#### Scenario: Rearranging the options does not rearrange the answer

- **WHEN** a person chooses the marked option of a decision that marked one listed last, and sends
- **THEN** `poll` returns that option's own label and command, unchanged by where the page showed it

#### Scenario: Nothing arrives chosen

- **WHEN** a decision marking an option is rendered
- **THEN** no option is selected, and the person's first selection is the first the page holds

### Requirement: Selected is shown separately from committed

A decision whose answer has been committed SHALL render differently from one the person has merely clicked in, and SHALL show back what was committed in the words the page offered for it.

A person who cannot tell what the caller will receive from what they have merely touched does not know what they have said.
What they are shown back is the option they chose, because that is the sentence they read when they chose it.
The command that option carries is the caller's and is not shown, for the reason the page shows no command anywhere.

#### Scenario: A committed answer is distinguishable

- **WHEN** one decision is committed and another only has options selected
- **THEN** the two render differently, and only the committed one shows what was committed

#### Scenario: A committed answer is shown in the words the person read

- **WHEN** a person commits an answer to a decision
- **THEN** what is shown back names the options they chose, in the labels the page offered, and carries no command

#### Scenario: A committed answer can be taken back

- **WHEN** a person reopens a committed decision before sending
- **THEN** it returns to the uncommitted state and its previous answer is not delivered

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

### Requirement: The page renders its own evidence

The page SHALL render highlighted source, a diff in a split or unified layout, and an inline illustration the caller supplied, without the caller writing the machinery for any of them.

Every asset needed SHALL be embedded in the binary.
No diagram engine SHALL be embedded, because an illustration is written by the caller and a layout engine surrenders position and emphasis to a machine.

The languages the page can highlight SHALL be the bundled library's own common set rather than a set chosen and maintained in this repository.
A set curated here is a list with an owner, and its failure mode is a caller naming a language nobody thought of and getting silence: the first caller outside this project named two, `csharp` and `diff`, and both rendered as plain text with nothing anywhere saying so.
The library's common set is maintained by the library and moves when it moves, which is the whole reason to prefer it.
Bundling every language the library ships is not required, because the common set covers what evidence in a code review is written in and the tail is large.

A language the page cannot highlight SHALL be refused when the document is served rather than rendered unhighlighted, and evidence naming no language SHALL render unhighlighted.

An illustration the document's stylesheet paints SHALL be rendered in a box that holds it, and the page SHALL render such an illustration without the caller writing the machinery for the box either.
An illustration that paints itself and an illustration a stylesheet paints are rendered in the same box, because the difference between them is what the caller supplied and not what the page owes it.

#### Scenario: A diff renders from a patch

- **WHEN** a decision carries a patch as evidence
- **THEN** the page renders it as a diff with both sides and per-line changes visible

#### Scenario: The languages a caller actually passes are highlighted

- **WHEN** a decision carries evidence in a language of the bundled library's common set
- **THEN** the page renders it highlighted, and no language in that set renders as plain text

#### Scenario: The narrow case does not shrink the evidence to nothing

- **WHEN** the page is viewed at a narrow width
- **THEN** the diff and any illustration remain legible, scrolling within their own bounds rather than the page scrolling sideways

#### Scenario: Rendering needs no network

- **WHEN** the page is loaded with the machine's outbound network unavailable
- **THEN** the highlighting, the diff, the fonts and the illustration all render

#### Scenario: A painted illustration renders in the same box as one that paints itself

- **WHEN** one decision carries a drawing and another carries markup the document's stylesheet paints
- **THEN** both render in the illustration's box, and the caller wrote the machinery for neither

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

### Requirement: The presentation is held to a stated floor

No interface text SHALL render below a stated minimum size, and the page SHALL paint its own background and text colours rather than inheriting them.
The page SHALL follow the reader's light or dark preference.
Every control SHALL be operable by keyboard and SHALL carry an accessible name.

Where an option is marked as the one the caller would take, the page SHALL say so in words within that option, and SHALL NOT convey it by colour, weight or position alone.
The words SHALL be part of the option's accessible name, so that a person hearing the page and a person seeing it learn the same thing at the same moment.

A floor that is written down is a floor that survives somebody being in a hurry, and a page that inherits its colours from whatever opened it is a page nobody has actually seen.

A mark carried only by paint is a mark that a screen reader, a monochrome screen and a person who does not know the convention all miss, and it is placement alone that is left telling them, which is the unstated convention this replaced.

#### Scenario: Nothing renders below the floor

- **WHEN** the rendered page's text is measured
- **THEN** no interface text is below the stated minimum

#### Scenario: Both preferences are painted

- **WHEN** the page is rendered under a light preference and under a dark one
- **THEN** each paints its own background and text, and neither leaves a colour undefined

#### Scenario: A decision can be answered without a pointer

- **WHEN** a person moves through the page by keyboard alone
- **THEN** every option, field and send can be reached, named and operated

#### Scenario: The mark is said rather than painted

- **WHEN** a decision marking an option is rendered and every colour is discarded
- **THEN** the marked option still says that it is the one the caller would take, and it says so in the name the option is announced under

### Requirement: Nothing on the page is written for a terminal

The page SHALL NOT render any command a document carries, and SHALL NOT present a flag, an argument or an identifier to the person as something to run.
What answering a decision does SHALL be said in the words of the decision.

A person reading this page is not at a prompt.
A command under every option is a second and harder copy of a choice already written above it, in the language the page exists to spare them, and a person who reads it as an instruction believes there is something left for them to type.
The command still travels: the caller gets each chosen option's command back verbatim, which is what carrying it was for.

The same holds for what the page renders around the evidence: a file in a diff SHALL be shown under the path the person knows it by, rather than under the one the patch format prefixes for a tool that is about to apply it.

#### Scenario: No command reaches the page

- **WHEN** a document whose every option carries a command is served and answered
- **THEN** no command text appears anywhere on the rendered page, and `poll` still returns each chosen option's command verbatim

### Requirement: An element the page hides paints nothing

An element carrying the attribute that hides it SHALL occupy no space and paint nothing, whatever else the stylesheet declares about it.

The page hides what it will show later: the line saying what was committed, and the control that takes an answer back.
A layout declared for one of those outranks the attribute that hides it, and what the person then sees is a small empty box beside the commit control, which they click and nothing happens.
The two halves of that failure live in different files and neither is wrong alone, so this is held by a rule no later declaration can beat rather than by whoever writes the next style.

#### Scenario: Nothing hidden is painted

- **WHEN** the page is rendered with nothing yet committed
- **THEN** every element the page hides has no box and paints nothing, and the page shows no control that does nothing when clicked

### Requirement: The page is painted in one stated palette, held to a contrast floor

Every colour the page paints SHALL be stated once, as a token, in the one block that states the palette, and no colour SHALL be written where it is used.
Every text the page renders SHALL meet a contrast ratio of at least 4.5 against what it is painted on, under a light preference and under a dark one.

The palette is the product's own: cut marble for the page, the grey of the block for its rules, slate for the ink, and the blue of the sea between the islands for the one accent a decision carries.
A page whose colours are the colours of the tool that composed the document tells the person they are reading that tool's output rather than their own decision, and the first person to see one said so within a minute.

An illustration the caller supplied is the caller's own, and the floor is what the page paints rather than what a fragment carries or what a caller's stylesheet paints.
The floor cannot reach a colour that arrives when the document does.
It is measured against a stylesheet that is read, and a caller's colours are not read by anything before a person is looking at them, so a floor claimed over them would be a floor nothing holds.
What is held instead is that the caller's colours land only inside the illustration, so a stylesheet that fails a person fails them in the box that was showing them something rather than in the decision they came to make.

#### Scenario: Every colour is stated once

- **WHEN** the stylesheet the page is built from is read
- **THEN** every colour it paints is a token declared in the block that states the palette, and none is written at the point of use

#### Scenario: Every text meets the contrast floor

- **WHEN** each text the page renders is measured against the colour behind it, under a light preference and under a dark one
- **THEN** each meets the stated ratio, and the measurement saw the whole page rather than a handful of elements

#### Scenario: A caller's colours do not reach the decision

- **WHEN** a document carries a stylesheet that paints its illustrations in colours of its own
- **THEN** every text the page itself renders is painted in the stated palette and still meets the floor

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

### Requirement: What the page lists is ordered by severity

The page SHALL order the entries of the source, and the decisions that follow them, by the severity the source carries: most severe first and informational last.

Entries of equal severity SHALL hold the order the document listed them, and decisions of equal severity SHALL do the same.
An entry or decision whose source carries no severity SHALL sort after every one that carries a severity, on the same stable basis.

The order SHALL be the surface's and SHALL NOT be reachable by the caller.
No document field SHALL set it, and the order a document lists its entries or decisions in SHALL decide nothing between severities.

The entries come before the decisions and this changes nothing about that.
What it decides is the order within each of the two, which the page has not stated until now and which has therefore been whatever order the caller happened to write.

Severity is copied from the source rather than accepted from the decision, so ordering by it is ordering by something already beyond the caller's influence.
The sequence of cards is layout, and the surface owns the layout.

A caller free to order the page is a caller free to put the error last, and a page a person opens to triage is the wrong place to leave that to chance.
The cost is that a caller with a narrative across several decisions cannot sequence it, and that narrative belongs in each card's own prose, where a reader who went straight to the error still finds it.

An unclassified severity sorting last is not a claim that it matters least.
It is the reading that overstates nothing, which is the same basis on which an unclassified action is never silently patched.

#### Scenario: The most severe card is the first a person meets

- **WHEN** a document holds decisions whose sources give severities of error, warning and info, listed in any order
- **THEN** the page renders the error first, the warning next and the info last

#### Scenario: The source list is ordered the same way

- **WHEN** a document cites source entries carrying severities of error, warning and info, listed in any order
- **THEN** the entries are listed most severe first and informational last, and still before the decisions

#### Scenario: Equal severities keep the order the document listed them

- **WHEN** a document holds three decisions whose sources all give the same severity
- **THEN** the page renders them in the order the document listed them

#### Scenario: The caller cannot order the page

- **WHEN** a document lists an informational decision first and an error last
- **THEN** the page renders the error first, and nothing the document stated changed that

#### Scenario: An entry whose source states no severity sorts last

- **WHEN** a document cites an entry carrying no severity alongside entries that carry one
- **THEN** the page lists it after every entry that carries a severity
