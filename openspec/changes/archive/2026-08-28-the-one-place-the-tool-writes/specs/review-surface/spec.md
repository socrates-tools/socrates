## RENAMED Requirements

- FROM: `#### Scenario: One variable moves everything agora has`
- TO: `#### Scenario: One variable moves everything both binaries have`

- FROM: `#### Scenario: Neither program writes in the other's home`
- TO: `#### Scenario: The surface keeps a directory of its own inside that home`

## MODIFIED Requirements

### Requirement: Session state lives outside the repository being discussed

Everything the surface reads and everything it writes SHALL live under socrates' home directory: the session state, the served document and what an operator states about the surface.
The surface SHALL NOT resolve a home of its own, SHALL NOT read an environment variable of its own, and SHALL NOT read a configuration file of its own.
It SHALL ask for that home through the one package that resolves it, so the rule is stated once for both binaries.

The surface's own state SHALL sit in a directory of its own under that home, beside what socrates writes rather than mixed into it.

Nothing SHALL be written inside the caller's working directory or the repository it is discussing.

agora serves exactly one consumer, and a dependency with one consumer is that consumer's.
Two homes, two variables and two configuration files made "where is your state?" a question with two halves, which is the cost the one-home rule exists to remove and which a second home put straight back.

The variable that once moved the surface's home SHALL be refused rather than ignored where it is still set, naming itself and naming the variable that replaced it.
A stated home that silently does nothing is a person who believes their state moved and is looking in the wrong place.

socrates already holds that nothing writes in the repository under review, and a surface that dropped a directory there would break it from the outside.

#### Scenario: The repository is untouched

- **WHEN** a surface is served, answered and ended from within a git repository
- **THEN** that repository's `git status` is byte-identical to what it was before

#### Scenario: One variable moves everything both binaries have

- **WHEN** socrates' home is stated, a review is run and a surface is served
- **THEN** the runs socrates wrote, the state the surface wrote and the one configuration file they read are all under the stated home, and nothing is written anywhere else

#### Scenario: A relative home is refused

- **WHEN** socrates' home is stated as a relative path
- **THEN** the command refuses, names the variable and the path, and does not fall back to the default

#### Scenario: The surface keeps a directory of its own inside that home

- **WHEN** a review is run and a surface is served from within it
- **THEN** the surface's state is under a directory of its own beside the runs, and no run directory holds any of it

#### Scenario: The variable that moved the surface's own home is refused

- **WHEN** the variable that once moved the surface's home is set and a surface is served
- **THEN** it refuses, names that variable and the one that replaced it, and starts no server
