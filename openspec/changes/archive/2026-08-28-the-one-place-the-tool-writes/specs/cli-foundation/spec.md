## MODIFIED Requirements

### Requirement: Configuration is one optional user file

socrates SHALL read one configuration file under its own home directory, and SHALL run with defaults where it is absent.
It SHALL hold the per-role model and effort, the inlined-diff ceiling, the agent deadline, the log ceiling, the run retention age, and one block holding what an operator states about the surface, and nothing else.

That block SHALL hold the address the surface binds and the address it is reached at.
The surface SHALL read it from the one file rather than parsing a file of its own, so that an operator configuring this tool opens one file.

Everything socrates reads and everything it writes SHALL live under that one home directory, and so SHALL everything the surface reads and writes.
One environment variable SHALL move all of it, and where none is stated it SHALL be a directory named for socrates under the person's own home.
A stated home SHALL be refused unless it is absolute, rather than ignored in favour of the default.

One directory rather than two is a question with one answer: where socrates keeps things.
The answer is given to a stranger reading an issue that will not reproduce, and a second directory belonging to a second binary is that answer growing a second half for a program the stranger has never heard of.
Refusing a relative home rather than ignoring it is the same argument: a stated home that silently does nothing is a person who believes their state moved and is looking in the wrong place.

socrates SHALL NOT read configuration from the repository under review, SHALL NOT create a configuration file on its own, and SHALL NOT write anything at all inside a repository.

#### Scenario: No configuration file

- **WHEN** `scrat review` is run with no configuration file present
- **THEN** it runs with defaults and creates no configuration file

#### Scenario: One variable moves everything

- **WHEN** socrates' home is stated and a review is run
- **THEN** both the configuration file it reads and the runs it writes are under the stated home, and nothing is written anywhere else

#### Scenario: A relative home is refused

- **WHEN** socrates' home is stated as a relative path
- **THEN** the command refuses, names the variable and the path, and does not fall back to the default

#### Scenario: The repository is never written to

- **WHEN** a review runs to completion
- **THEN** nothing inside the repository has changed except the git objects the snapshot needed, and no socrates file exists anywhere under it

#### Scenario: An unreadable configuration file is a refusal

- **WHEN** the configuration file exists and does not parse
- **THEN** the command refuses naming the file and the parse error, and exits with the usage code, rather than silently running on defaults

#### Scenario: What an operator states about the surface is read from the one file

- **WHEN** the configuration file states the address the surface binds and the address it is reached at, and a surface is served
- **THEN** the surface binds what was stated and prints what was stated, having read no file of its own

#### Scenario: A surface block that names a field socrates does not know is refused

- **WHEN** the surface block of the configuration file names a field that is not one of its own
- **THEN** the command refuses naming the field, exactly as it does for a field of any other block
