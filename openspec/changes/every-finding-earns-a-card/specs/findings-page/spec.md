## REMOVED Requirements

### Requirement: The page foregrounds what the person has to act on

**Reason**: Its reasoning stands and its remedy does not.
It was right that a round of one question and five mechanical defects must not read as six equal questions.
Its remedy was to give the answerable findings in full and the rest as a short list, and a person who read a page built that way asked for all of it in full, because a person signs off on the round rather than on the part of it that survived triage.

The short list is no longer the skill's to arrange either.
The surface now copies every source entry onto the page with its name, severity, audience and location, and opens the page with them.
So the requirement's own rule that the list be taken from what socrates reported rather than retyped is held by the surface, and restating it here would be a second copy of a rule that already has an owner.

What the list cannot carry is what was done about a finding and why, and that is the part a person asked to see.

**Migration**: The requirement below replaces it.
Foregrounding is kept and moves to two mechanisms that demote nothing: the severest card is the first the person meets, and a finding nobody needs to answer carries no control.
Its rule about the list being copied rather than retyped is kept by `surface-document` and `surface-page`, which own it.
Both of its scenarios are withdrawn with it, and the scenarios below cover the shape that replaces them, including the round that holds nothing to answer.

## ADDED Requirements

### Requirement: The page carries the whole round, one card per finding

The page SHALL carry every finding in the round, one card for each.
No finding SHALL be reduced to a line in a list or a sentence in the preamble.

The source list the surface renders SHALL remain what it is, an index of everything found, and SHALL NOT be what stands in for a finding the person never sees in full.

The page SHALL foreground what the person has to act on by the order the cards arrive in and by the absence of controls on the findings they do not.

A parked round holding one question about intent and five mechanical defects is one decision and five things the session is already fixing.
A page that treats all six evenly buries the only thing the person has to act on, and a page that drops the five leaves them unable to see what else was found.
What separates them is that the severest card is the first the person meets, and that a finding the session resolved carries no option for them to weigh.

The skill SHALL say that a finding the session resolved is put on the page as what the finding was and what resolved it, and that the diff goes on it where the diff is what makes the resolution clear.

The skill SHALL say that such a finding carries no options, and SHALL say why: an option on a settled finding is a control that changes nothing, and a person reads it as a question still owed an answer.

Where the round holds nothing a person answers, the page SHALL still carry every finding as a card, and SHALL present none of them as a decision.
A clean run and a round the author patched alone are rounds with nothing to act on, and a card that asks nothing is not a decision: it is the account of what was found and what was done, which is the whole of what those rounds have to hand over.

#### Scenario: Every finding in the round is a card

- **WHEN** the skill describes a page for a round holding both findings a person answers and findings the author resolved
- **THEN** it says every finding in the round gets a card of its own, and none is reduced to a line in a list or a sentence in the preamble

#### Scenario: A resolved finding is stated rather than asked

- **WHEN** the skill describes a card for a finding the session already resolved
- **THEN** it says the card states what the finding was and what resolved it, carries the diff where that is what makes it clear, and offers no options

#### Scenario: The reason for offering no options is given

- **WHEN** the skill says a resolved finding carries no options
- **THEN** it says that a control which changes nothing reads as a question the person still has to answer

#### Scenario: A round with nothing to answer is still a round in full

- **WHEN** the skill describes a page for a clean run, or for a parked round whose findings the author answers alone
- **THEN** it says every finding is still a card, that none carries an option, and that nothing is presented as a decision

#### Scenario: Foregrounding survives without demoting anything

- **WHEN** the skill describes how the page separates what needs the person from what does not
- **THEN** it names the severity order and the absence of controls, and does not tell the session to shorten, list or omit any finding
