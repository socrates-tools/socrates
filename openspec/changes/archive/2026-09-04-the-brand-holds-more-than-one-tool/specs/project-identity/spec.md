## MODIFIED Requirements

### Requirement: The module path names the organisation

The module path SHALL be the project's repository under the organisation that publishes it, and every import within the project SHALL use it.

The path SHALL be written once as data where a build reads it, so that what a build stamps into itself and what the compiler resolves cannot disagree.

The repository SHALL be named for the tool it holds, and the organisation SHALL be named for the brand that publishes it.
A person holding the command holds the repository's name, and a person holding the repository's name knows which of the organisation's tools it is.
An organisation that publishes more than one tool and a repository named after the organisation leave the second tool with no name to take.

Renaming the module on its own SHALL change where the source lives and nothing about what the tool is called: the binaries, the commands, the home directory and the configuration file are unaffected by it.
Where the tool itself is renamed the repository SHALL follow it, and that is a change to both stated as one rather than a consequence of either.

#### Scenario: Every import resolves under the organisation

- **WHEN** the module is built and its imports are read
- **THEN** the module path names the organisation's repository, and no import names any other module of this project

#### Scenario: The build reads the module path rather than repeating it

- **WHEN** the build stamps what a binary says about itself
- **THEN** it composes that from the module path stated once, and no second copy of the path exists for it to disagree with

#### Scenario: The repository is named for the tool

- **WHEN** the last segment of the module path is compared with the command the tool installs
- **THEN** they are the same name, and the segment before it is the organisation rather than a repeat of that name
