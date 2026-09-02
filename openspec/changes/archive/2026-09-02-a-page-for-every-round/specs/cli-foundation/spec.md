## RENAMED Requirements

- FROM: `#### Scenario: A round a person reads says so`
- TO: `#### Scenario: A round that waits on a person says so`

- FROM: `#### Scenario: A round the author answers alone says nothing about a person reading it`
- TO: `#### Scenario: A round the author answers alone says nothing about waiting`

- FROM: `#### Scenario: `scrat status` does not carry it`
- TO: `#### Scenario: `scrat status` does not carry the wait`

## ADDED Requirements

### Requirement: `scrat findings` says whether the round waits on a person

`scrat findings` SHALL report whether the round it is printing waits on a person, as one line for a reader and one field in the object `--json` prints.

It SHALL report it from the one rule that decides it, and SHALL NOT restate the condition.

The line and the field SHALL name the wait rather than the reading.
Every round that completed a review is handed over as a page, so a fact reported as whether a person reads the round would be true of every round it is printed for and would decide nothing.
What still needs deciding is whether the session stops at the page, and that is what is reported.

`scrat status` SHALL NOT carry it.
`status` is called in a loop and everything it prints lands in the waiting session's context, which is the scarcest resource in the system.
`findings` is read once, on exactly the branches where the fact is needed, and is already where socrates tells the author the rules.

#### Scenario: A round that waits on a person says so

- **WHEN** `scrat findings` is run against a run that escalated, or against a parked round holding a finding a person answers
- **THEN** it says the round waits on a person, and the field is in the object `--json` prints

#### Scenario: A round the author answers alone says nothing about waiting

- **WHEN** `scrat findings` is run against a clean run, or against a parked round whose findings the author answers alone
- **THEN** it does not say the round waits on a person, and it says nothing that would stop the session handing the round over as a page

#### Scenario: `scrat status` does not carry the wait

- **WHEN** the lines and the JSON object of `scrat status` are examined for a run whose round waits on a person
- **THEN** neither carries it, and what `status` prints is what it printed before

## REMOVED Requirements

### Requirement: `scrat findings` says whether the round is one a person reads

**Reason**: What it reported is now true of every round it can be run against, so it distinguishes nothing.
A person reads every round that completed a review, because every one of them is handed over as a page.

**Migration**: The line and the field are kept and renamed for what they now decide, which is whether the session waits on the page rather than whether one is made.
`read by a person` becomes a line naming the wait, and the `--json` field `read_by_a_person` becomes `awaits_a_person`.
The only reader of either is the skill, which this change rewrites.
