## ADDED Requirements

### Requirement: The skill renders a round a person reads as a page

The skill SHALL say that where socrates reports the round as one a person reads, the session renders it as an HTML page through the `artifact-design` skill, hands the person the link, and summarises rather than pastes.

It SHALL say it on both branches where such a round arrives: `escalated`, where the findings are the person's and socrates is finished, and `awaiting-response` holding a finding a person answers.

It SHALL take the fact from what socrates reports rather than restating the condition, so that a session never decides for itself whether a round is one a person reads.

It SHALL hold the page to the stated bar rather than describing the page's qualities again, and the bar SHALL appear before the branches that refer to it.

Nothing about the page SHALL appear in the waiting step.
Everything in that step is run in a loop, and a page is made once, after the run has stopped moving.

#### Scenario: Both branches say to render the page

- **WHEN** the skill is rendered
- **THEN** the `escalated` branch and the `awaiting-response` branch each say to make the page and hand over the link

#### Scenario: The session does not decide the condition itself

- **WHEN** the skill describes when to make a page
- **THEN** it names what socrates reports, and states no condition of its own over the state or the counts

#### Scenario: A round the author answers alone gets no page

- **WHEN** the skill describes a parked round whose findings the author answers alone
- **THEN** it says to patch, respond and continue, and says nothing about a page

#### Scenario: The bar is stated before the branches that refer to it

- **WHEN** the skill is rendered
- **THEN** the bar appears once, before the `escalated` and `awaiting-response` branches, and neither branch restates it

#### Scenario: The skill says the page is read only

- **WHEN** the skill describes making the page
- **THEN** it says the page collects nothing, declares no capability, and that the person answers in the terminal

#### Scenario: The skill says the page and the terminal use the same identifiers

- **WHEN** the skill describes what goes on the page
- **THEN** it says every finding keeps the identifier, severity and audience `scrat findings` reported, so the person can name one and the session knows which
