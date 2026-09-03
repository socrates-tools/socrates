## ADDED Requirements

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

## MODIFIED Requirements

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
