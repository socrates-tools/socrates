## RENAMED Requirements

- FROM: `#### Scenario: Both branches say to render the page`
- TO: `#### Scenario: Both branches say to serve the page`
- FROM: `#### Scenario: The skill says the page is read only`
- TO: `#### Scenario: The skill says how an answer reaches the run`

## MODIFIED Requirements

### Requirement: The skill drives only the surface that exists

Every command and flag the skill names SHALL exist, and a test SHALL check each against the binary that offers it.

The skill now names two binaries.
An extraction that recognised only one of them would leave every command of the other unchecked, so a skill naming an `agora` flag that does not exist would pass the check that exists to catch exactly that.
The check SHALL therefore attribute each named command to its binary and SHALL fail on a name no binary offers.

A skill naming a flag that does not exist is a session doing something plausible and wrong, and it fails at the point where the person has stopped watching.

#### Scenario: Every command named is real

- **WHEN** the commands and flags the skill names are extracted and checked against the binaries that offer them
- **THEN** every one of them exists

#### Scenario: A removed flag fails the suite

- **WHEN** a flag the skill names is removed from the binary that offered it
- **THEN** the check fails, naming the flag and the skill

#### Scenario: A flag of the second binary is checked too

- **WHEN** the skill names a command or flag of the surface binary that the surface binary does not offer
- **THEN** the check fails, and it does so for the same reason and in the same way as for the first binary

### Requirement: The skill renders a round a person reads as a page

The skill SHALL say that where socrates reports the round as one a person reads, the session drafts a document, serves it with `agora`, hands the person the address, waits on `agora poll`, and records what came back with `socrates respond`.

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
