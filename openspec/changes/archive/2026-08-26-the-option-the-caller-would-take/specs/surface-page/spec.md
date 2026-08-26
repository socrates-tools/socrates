## MODIFIED Requirements

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
