## MODIFIED Requirements

### Requirement: Session state lives outside the repository being discussed

Everything the surface reads and everything it writes SHALL live under agora's own home directory: the session state, the served document and the configuration file an operator may write.
One environment variable SHALL move it, and where none is stated it SHALL be a directory named for agora under the person's own home.
A stated home SHALL be refused unless it is absolute, rather than ignored in favour of the default.

Nothing SHALL be written inside the caller's working directory or the repository it is discussing.

That home SHALL be agora's and not socrates'.
The two are separate programs, and one keeping its state inside the other's directory is the boundary going quietly, in the one place a reader would never look for it.

socrates already holds that nothing writes in the repository under review, and a surface that dropped a directory there would break it from the outside.

#### Scenario: The repository is untouched

- **WHEN** a surface is served, answered and ended from within a git repository
- **THEN** that repository's `git status` is byte-identical to what it was before

#### Scenario: One variable moves everything agora has

- **WHEN** agora's home is stated and a surface is served
- **THEN** the configuration it read and the session state it wrote are both under the stated home, and nothing is written anywhere else

#### Scenario: A relative home is refused

- **WHEN** agora's home is stated as a relative path
- **THEN** the command refuses, names the variable and the path, and does not fall back to the default

#### Scenario: Neither program writes in the other's home

- **WHEN** a review is run and a surface is served from within it
- **THEN** socrates has written only under its own home and agora only under its own
