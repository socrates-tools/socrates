## MODIFIED Requirements

### Requirement: The skill renders a round a person reads as a page

The skill SHALL say that where socrates reports the round as one a person reads, the session drafts a document, serves it with `agora`, hands the person the address, waits on `agora poll`, and records what came back with `socrates respond`.

It SHALL say that the address `agora serve` prints holds the session that page belongs to, and that `agora poll` is given that session.
It SHALL say it once, where the loop is described, rather than on each branch that serves a page.
A session that is not told this polls a surface holding more than one page and is told about somebody else's round, or about none, and it has no way to find out which.

It SHALL say that the session is kept for the whole review, and that a second round a person reads is written to the page they already have rather than served as a new one.
A review is two rounds about one change and one person reading them, so a second address is a second tab, and the first goes stale with nothing on it saying so.

It SHALL say it on both branches where such a round arrives: `escalated`, where the findings are the person's and socrates is finished, and `awaiting-response` holding a finding a person answers.

It SHALL take the fact from what socrates reports rather than restating the condition, so that a session never decides for itself whether a round is one a person reads.

It SHALL hold the page to the stated bar rather than describing the page's qualities again, and the bar SHALL appear before the branches that refer to it.

It SHALL say that the person answers on the page, that each control names the command it becomes, and that the session is what records the answer.
It SHALL NOT tell the session to ask the person to restate in the terminal a decision they already made on the page.

Nothing about the page SHALL appear in the waiting step.
Everything in that step is run in a loop, and a page is made once, after the run has stopped moving.

#### Scenario: Both branches say to serve the page

- **WHEN** the skill is rendered
- **THEN** the `escalated` branch and the `awaiting-response` branch each say to serve the page and hand over the address

#### Scenario: The skill says to carry the session from serve to poll

- **WHEN** the skill describes serving a page and waiting on it
- **THEN** it says the printed address holds the session, and that `agora poll` is given that session

#### Scenario: The session is stated once rather than on every branch

- **WHEN** the skill is rendered
- **THEN** it says how the session reaches `agora poll` once, and neither branch restates it

#### Scenario: A second round reaches the page the person already has

- **WHEN** the skill describes a second round that is one a person reads
- **THEN** it says to write it to the session the first round was served under, and does not say to hand over a second address

#### Scenario: The session does not decide the condition itself

- **WHEN** the skill describes when to make a page
- **THEN** it names what socrates reports, and states no condition of its own over the state or the counts

#### Scenario: A round the author answers alone gets no page

- **WHEN** the skill describes a parked round whose findings the author answers alone
- **THEN** it says to patch, respond and continue, and says nothing about a page

#### Scenario: The bar is stated before the branches that refer to it

- **WHEN** the skill is rendered
- **THEN** the bar appears once, before the `escalated` and `awaiting-response` branches, and neither branch restates it

#### Scenario: The skill says how an answer reaches the run

- **WHEN** the skill describes making the page
- **THEN** it says the person answers on the page, that each control names the command it becomes, and that the session records it with `socrates respond`

#### Scenario: The skill says the page and the terminal use the same identifiers

- **WHEN** the skill describes what goes on the page
- **THEN** it says every finding keeps the identifier, severity and audience `socrates findings` reported, so the person can name one and the session knows which

#### Scenario: A decision is still the person's to make

- **WHEN** the skill describes a finding that needs a person
- **THEN** it says the session does not answer one on its own behalf, whichever surface the answer arrives through
