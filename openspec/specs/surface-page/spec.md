# surface-page Specification

## Purpose
TBD - created by archiving change a-page-that-answers-back. Update Purpose after archive.

## Requirements

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

An illustration the document's stylesheet paints SHALL be rendered in a box that holds it, and the page SHALL render such an illustration without the caller writing the machinery for the box either.
An illustration that paints itself and an illustration a stylesheet paints are rendered in the same box, because the difference between them is what the caller supplied and not what the page owes it.

#### Scenario: A diff renders from a patch

- **WHEN** a decision carries a patch as evidence
- **THEN** the page renders it as a diff with both sides and per-line changes visible

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
