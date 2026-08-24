## ADDED Requirements

### Requirement: `scrat findings` says whether the round is one a person reads

`scrat findings` SHALL report whether the round it is printing is one a person reads, as one line for a reader and one field in the object `--json` prints.

It SHALL report it from the one rule that decides it, and SHALL NOT restate the condition.

`scrat status` SHALL NOT carry it.
`status` is called in a loop and everything it prints lands in the waiting session's context, which is the scarcest resource in the system.
`findings` is read once, on exactly the branches where the fact is needed, and is already where socrates tells the author the rules.

#### Scenario: A round a person reads says so

- **WHEN** `scrat findings` is run against a run that escalated, or against a parked round holding a finding a person answers
- **THEN** it says the round is one a person reads, and the field is in the object `--json` prints

#### Scenario: A round the author answers alone says nothing about a person reading it

- **WHEN** `scrat findings` is run against a clean run, or against a parked round whose findings the author answers alone
- **THEN** it does not say the round is one a person reads

#### Scenario: `scrat status` does not carry it

- **WHEN** the lines and the JSON object of `scrat status` are examined for a run whose round is one a person reads
- **THEN** neither carries it, and what `status` prints is what it printed before
