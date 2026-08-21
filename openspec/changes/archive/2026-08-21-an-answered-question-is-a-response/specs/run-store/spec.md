## MODIFIED Requirements

### Requirement: A run holds one of seven states and four of them are terminal

A run SHALL be in exactly one of `starting`, `reviewing`, `awaiting-response`, `clean`, `escalated`, `failed`, `cancelled`.
`clean`, `escalated`, `failed` and `cancelled` SHALL be terminal.

- `awaiting-response`: a review completed, what it found can be answered, and socrates is waiting for the author to fix, to argue, or to bring a question to the person and record what they decided.
- `clean`: a review completed and reported nothing blocking.
- `escalated`: a review completed and socrates has nothing further it can do, because what it found can be answered by nobody or because the round cap ran out.
- `failed`: no review completed.
- `cancelled`: the person stopped it.

`failed` SHALL NOT be merged with `escalated`.
They call for different actions: escalated means read the findings, failed means nothing has been checked.

`awaiting-response` SHALL NOT be terminal and SHALL NOT be reported as a verdict.
A run parked on the author is a run that has not concluded, and reporting it as one would let a session stop at the point where the dialogue was about to happen.
A run parked on a question the person must answer is parked on the same terms: the dialogue has not happened yet, and what is waited on is a person rather than a patch.

#### Scenario: A terminal run does not move

- **WHEN** a run has reached a terminal state
- **THEN** no later write changes its state or its verdict

#### Scenario: A failed run reports no findings as its verdict

- **WHEN** a run ends as `failed`
- **THEN** its recorded verdict is the failure and its reason, and it does not present an empty findings set as a result

#### Scenario: A parked run is not finished

- **WHEN** a run is `awaiting-response`
- **THEN** commands reporting on it exit with the code meaning still working, and no verdict is printed

#### Scenario: A run parked on a question is not a verdict either

- **WHEN** a run is `awaiting-response` because a finding asks the person a question
- **THEN** it is reported as still working rather than as escalated, and no verdict is printed

#### Scenario: A parked run can still be cancelled

- **WHEN** a run is `awaiting-response` and is cancelled
- **THEN** it becomes `cancelled`, its snapshot worktree is removed, and its rounds' artifacts are kept
