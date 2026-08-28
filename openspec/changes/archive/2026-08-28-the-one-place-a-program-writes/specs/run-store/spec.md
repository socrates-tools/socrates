## MODIFIED Requirements

### Requirement: A run is a directory whose name answers the common question

Every run SHALL own one directory under socrates' home, beside the configuration file rather than under a directory of its own, named `<UTC instant>-<repository key>`, where the instant is second-resolution and lexically sortable and the repository key is derived from the absolute path of the repository root.

Naming SHALL carry both, so that the most recent run for a given repository is answerable by listing directory names and reading no files.
The command that asks that question is called in a loop by a session waiting on a review, and a design that opens every run to answer it pays that cost on every poll.

#### Scenario: Runs sort by when they started

- **WHEN** three runs are started in the same repository in order
- **THEN** their directory names sort lexically into the order they were started

#### Scenario: Two repositories do not collide

- **WHEN** runs are started in two different repositories within the same second
- **THEN** their directory names differ, and each names its own repository's key

#### Scenario: The same path always gets the same key

- **WHEN** a repository key is derived twice for one repository root
- **THEN** the two are equal, and a different repository root produces a different key

#### Scenario: A run is found beside the configuration

- **WHEN** a review is run and socrates' home is examined
- **THEN** the run's directory and the configuration file are under that one home, and no run exists outside it
