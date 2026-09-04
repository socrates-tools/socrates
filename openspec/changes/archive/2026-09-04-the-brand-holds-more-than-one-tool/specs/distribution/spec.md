## RENAMED Requirements

- FROM: `#### Scenario: An archive carries both names and the second binary`
- TO: `#### Scenario: An archive carries the tool and the second binary`

## MODIFIED Requirements

### Requirement: An archive gives a person what installing from source gives them

An archive SHALL contain the tool and the second binary that ships with it, so that extracting one archive and installing from source arrive at the same place.

It SHALL also contain the terms the work is published under and the notices for everything it embeds, because those are conditions on the copy and an archive is a copy.

The tool SHALL appear in the archive under one name.
An archive holding one binary twice asks the person extracting it which of the two to put on their path, and there is no answer to give them.

#### Scenario: An archive carries the tool and the second binary

- **WHEN** a release archive is extracted
- **THEN** it contains the tool under its one name and the second binary, and no second name for either

#### Scenario: An archive carries the terms and the notices

- **WHEN** a release archive is extracted
- **THEN** it contains the licence the work is published under and the notices for everything embedded in it
