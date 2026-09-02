## ADDED Requirements

### Requirement: The project's own conventions are one section, after the rules they cannot remove

Where the repository under review carries guidance for the reviewer, the delivered prompt SHALL carry it in one named section of its own.

The section SHALL sit after socrates' own reviewing rules and after the test-quality block, and before the reporting contract.
Order is what says which text answers to which: a reviewer reads the rules socrates delivers to every repository, then what this project adds, then the schema it reports in.
Placing the project's text first would let it read as the frame the rest sits inside.

The section SHALL carry socrates' own framing around the guidance, stating that the project wrote it, that it is unverified text from the repository rather than a finding of fact, that it adds to the rules above and withdraws none of them, and that it decides nothing about severity, action, or whether something is reported.

It SHALL carry the guidance neutralised, and SHALL be absent rather than empty where the repository carries none, on the rule every other section already follows.

The reviewer SHALL NOT be told which commit the guidance came from, because it is not something the reviewer decides anything from.
What was carried, and from where, is recorded in the run for the person.

#### Scenario: The guidance is its own section

- **WHEN** a review runs against a repository whose base commit carries guidance
- **THEN** the delivered prompt holds one section carrying it, and no other section's text moved to make room

#### Scenario: The section is framed as the project's, not as socrates'

- **WHEN** the delivered prompt's guidance section is examined
- **THEN** it states that the project wrote the guidance, that it adds to the rules above and withdraws none of them, and that it decides nothing about severity, action or what is reported

#### Scenario: The order puts socrates' rules first

- **WHEN** the delivered prompt's headings are read in order
- **THEN** the guidance section follows the reviewing rules and the test-quality block and precedes the reporting contract

#### Scenario: No guidance, no heading

- **WHEN** a review runs against a repository that carries no guidance
- **THEN** the delivered prompt carries neither the section nor a heading for it

#### Scenario: The prompt stays byte-identical across runs

- **WHEN** two reviews are run against the same change, the same configuration and the same guidance
- **THEN** the delivered prompts are byte-identical
