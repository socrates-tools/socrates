## ADDED Requirements

### Requirement: An option may be marked as the one the caller would take

A decision MAY mark one of its options as the one the caller would take.
At most one option per decision SHALL carry the mark, in either mode, and `serve` SHALL refuse a document marking two and SHALL name the decision.

The mark SHALL rank the option and SHALL NOT answer for the person.
A marked option SHALL arrive unselected, and a decision the person never touches SHALL come back with nothing chosen and no command, exactly as an unmarked one does.

A document that marks no option SHALL be served unchanged, and nothing about a decision without a mark SHALL be different from what it is today.

The caller has just read the finding and usually knows which option it would take.
A page that cannot carry that judgement makes the person re-derive it, and the caller that wanted to state it puts it in the prose or in the option order instead, which is a convention the surface never agreed to honour.

Two marks are a caller that has not decided, and a person who has to combine two recommendations has been handed the decision back with extra steps.

Marking is separated from selecting because a pre-selected option is an answer nobody gave, and by the time it reaches the terminal it is indistinguishable from one the person made.

#### Scenario: A decision names the option the caller would take

- **WHEN** a document marks one option of a decision
- **THEN** the page shows that option as the one the caller would take, and the other options carry no such mark

#### Scenario: A document marking two options is refused before serving

- **WHEN** a document marks two options of one decision
- **THEN** `serve` fails, names the decision, and starts no server

#### Scenario: Marking is not choosing

- **WHEN** a document marking an option is served and the person sends without touching that decision
- **THEN** `poll` returns that decision with nothing chosen and no command

#### Scenario: A document that marks nothing is unaffected

- **WHEN** a document whose decisions mark no option is served and answered
- **THEN** the page renders its options in the order the document listed them, and what `poll` returns is what it would have returned before the mark existed
