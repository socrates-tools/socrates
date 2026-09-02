## RENAMED Requirements

- FROM: `#### Scenario: Both branches say to serve the page`
- TO: `#### Scenario: Every branch that has a result says to serve the page`

- FROM: `#### Scenario: A round the author answers alone gets no page`
- TO: `#### Scenario: A round the author answers alone gets a page nobody waits on`

- FROM: `#### Scenario: A second round reaches the page the person already has`
- TO: `#### Scenario: A later round reaches the page the person already has`

## ADDED Requirements

### Requirement: The skill hands every round over as a page

The skill SHALL say that on every branch where a round completed a review, the session drafts a document, serves it with `agora`, hands the person the address, and says briefly what it holds.
Those branches are `clean`, `awaiting-response` and `escalated`.

It SHALL say that the session waits on `agora poll` only where socrates reports that the round waits on a person, and that on the other branches it hands over the address and carries on with patching or reporting.
A session told to wait on every page would stop on a clean run for an answer nobody was asked for.

It SHALL say that `failed` and `cancelled` get no page: a failed run completed no review, a cancelled one was stopped on purpose, and in both cases the reason is one line socrates already prints.

It SHALL say that the address `agora serve` prints holds the session that page belongs to, and that `agora poll` is given that session.
It SHALL say it once, where the loop is described, rather than on each branch that serves a page.
A session that is not told this polls a surface holding more than one page and is told about somebody else's round, or about none, and it has no way to find out which.

It SHALL say that the session is kept for the whole review, and that a later round is written to the page the person already has rather than served as a new one.
A review is two rounds about one change and one person reading them, so a second address is a second tab, and the first goes stale with nothing on it saying so.
Every round of a run that is still going reaches the page, so the address is handed over once, at the first round that reached a result.

It SHALL take the fact from what socrates reports rather than restating the condition, so that a session never decides for itself whether a round waits on a person.

It SHALL hold the page to the stated bar rather than describing the page's qualities again, and the bar SHALL appear before the branches that refer to it.

It SHALL say that the person answers on the page, that each control names the command it becomes, and that the session is what records the answer.
It SHALL NOT tell the session to ask the person to restate in the terminal a decision they already made on the page.

Nothing about the page SHALL appear in the waiting step.
Everything in that step is run in a loop, and a page is made once, after the run has stopped moving.

#### Scenario: Every branch that has a result says to serve the page

- **WHEN** the skill is rendered
- **THEN** the `clean`, `awaiting-response` and `escalated` branches each say to serve the page and hand over the address, and the `failed` and `cancelled` branches say nothing about a page

#### Scenario: The skill says which branches wait

- **WHEN** the skill describes the branches that serve a page
- **THEN** it says to wait on `agora poll` where socrates reports the round waits on a person, and to carry on without waiting where it does not

#### Scenario: A round the author answers alone gets a page nobody waits on

- **WHEN** the skill describes a parked round whose findings the author answers alone
- **THEN** it says to serve the page and hand over the address, and then to patch, respond and continue without waiting on it

#### Scenario: The skill says to carry the session from serve to poll

- **WHEN** the skill describes serving a page and waiting on it
- **THEN** it says the printed address holds the session, and that `agora poll` is given that session

#### Scenario: The session is stated once rather than on every branch

- **WHEN** the skill is rendered
- **THEN** it says how the session reaches `agora poll` once, and no branch restates it

#### Scenario: A later round reaches the page the person already has

- **WHEN** the skill describes a second round
- **THEN** it says to write it to the session the first round was served under, and does not say to hand over a second address

#### Scenario: The session does not decide the condition itself

- **WHEN** the skill describes when to wait on a page
- **THEN** it names what socrates reports, and states no condition of its own over the state or the counts

#### Scenario: The bar is stated before the branches that refer to it

- **WHEN** the skill is rendered
- **THEN** the bar appears once, before the branches that serve a page, and no branch restates it

#### Scenario: The skill says how an answer reaches the run

- **WHEN** the skill describes making the page
- **THEN** it says the person answers on the page, that each control names the command it becomes, and that the session records it with `socrates respond`

#### Scenario: The skill says the page and the terminal use the same identifiers

- **WHEN** the skill describes what goes on the page
- **THEN** it says every finding keeps the identifier, severity and audience `socrates findings` reported, so the person can name one and the session knows which

#### Scenario: A decision is still the person's to make

- **WHEN** the skill describes a finding that needs a person
- **THEN** it says the session does not answer one on its own behalf, whichever surface the answer arrives through

## REMOVED Requirements

### Requirement: The skill renders a round a person reads as a page

**Reason**: It taught the page as a branch, reached only where socrates reported that a person had something to act on.
Every round that completed a review is now handed over as a page, so the branch is the loop and what is left to teach on each branch is whether the session waits.

**Migration**: `The skill hands every round over as a page` replaces it and keeps everything it required: the session carried from `agora serve` to `agora poll` and stated once, one session for the whole review, the fact taken from what socrates reports rather than restated, the bar stated before the branches, the answer recorded by `socrates respond`, and nothing about the page in the waiting step.
What is added is the branch that serves a page and does not wait on it.
