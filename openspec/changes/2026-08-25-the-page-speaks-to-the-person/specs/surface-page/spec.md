## MODIFIED Requirements

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

## ADDED Requirements

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

An illustration the caller supplied is the caller's own, and the floor is what the page paints rather than what a fragment carries.

Which colours those are is a judgement and is recorded in this change.
What is held mechanically is that they are stated in one place, so that replacing them is one edit, and that quiet never becomes unreadable, which is the direction a stone palette fails in.

#### Scenario: Every colour is stated once

- **WHEN** the stylesheet the page is built from is read
- **THEN** every colour it paints is a token declared in the block that states the palette, and none is written at the point of use

#### Scenario: Every text meets the contrast floor

- **WHEN** each text the page renders is measured against the colour behind it, under a light preference and under a dark one
- **THEN** each meets the stated ratio, and the measurement saw the whole page rather than a handful of elements
