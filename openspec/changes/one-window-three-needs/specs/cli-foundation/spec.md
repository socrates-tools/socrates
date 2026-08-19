## MODIFIED Requirements

### Requirement: The binary is `scrat` and it does one thing

socrates SHALL ship exactly one binary, named `scrat`.
It SHALL provide `review`, `status`, `findings`, `watch`, `logs` and `cancel`, and it SHALL NOT provide a command that commits, pushes, opens a pull request, watches continuous integration, or runs the repository's tests.

The scope rule is mechanical rather than a judgement call about what is close enough to review, because judgement calls erode.

#### Scenario: The command surface is what it says

- **WHEN** `scrat --help` is run
- **THEN** it lists `review`, `status`, `findings`, `watch`, `logs` and `cancel`, and nothing that lands or verifies a change

## ADDED Requirements

### Requirement: The three reading commands are for three consumers

`status`, `watch` and `logs` SHALL stay distinct, and none SHALL grow into another.

- `status` answers one shot, in a few lines, for the session waiting on the review. It never streams and never follows.
- `watch` follows live, for the person in a second tab. It always follows and always stops at a terminal state.
- `logs` reads what happened, for the person afterwards. It never follows.

`status` SHALL NOT gain a follow option and `watch` SHALL NOT gain a one-shot option.
Either would be the same command reached two ways, and the reason they are separate is that one of them prints into a language model's context and two print into a person's terminal, which is a difference of an order of magnitude in what verbosity costs.

#### Scenario: Status stays small when the reviewer is loud

- **WHEN** a reviewer has emitted a large volume of output and `scrat status` is run
- **THEN** its output is the same handful of lines it prints for a quiet run

#### Scenario: Neither command grows the other's flag

- **WHEN** the flags of `status` and `watch` are examined
- **THEN** `status` offers no way to follow and `watch` offers no way to return after one read

#### Scenario: All three agree about the run

- **WHEN** the three commands are run against one terminal run
- **THEN** all three report the same verdict and exit with the same code
