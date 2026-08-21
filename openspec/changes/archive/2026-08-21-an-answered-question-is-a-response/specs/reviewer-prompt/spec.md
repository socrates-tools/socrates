## MODIFIED Requirements

### Requirement: A re-review carries the provenance clause

A round after the first SHALL carry a clause stating:

- The code under re-review was written by the change's author in response to findings. It is unreviewed new code, not a settled resolution of the findings that prompted it, and it is reviewed to exactly the same adversarial standard as the original change.
- Prior findings, claimed fixes and rejections are claims, not evidence. Verify each against the current code, and judge independently whether what the patches introduced is correct, not merely whether it does what was prescribed.
- A decision recorded against a finding that asked a person is a claim that a person answered it, and what was decided is a claim about what the change is for. Judge whether the code now matches what was decided, and say so where the decision leaves the finding unresolved.
- A test written or changed in the same round as the code it exercises is part of that round's claim, not independent proof. Judge whether its asserted outcome is the right outcome and whether it could still pass with the code wrong.

The clause SHALL be absent from the first round, where there is no prior claim for it to be about.

A recorded decision SHALL reach the next round labelled as a decision rather than as a fix, in the same labelled block that carries prior findings and responses.
A reviewer told that a mechanical change closed a question about purpose is being told something nobody claimed, and it is exactly the substitution the resolutions are kept separate to prevent.

#### Scenario: The clause reaches a re-review

- **WHEN** a second round is invoked
- **THEN** the delivered prompt carries each of the statements in this requirement

#### Scenario: The clause is absent from a first review

- **WHEN** the first round is invoked
- **THEN** the delivered prompt carries none of them

#### Scenario: A claimed fix is not taken on trust

- **WHEN** the author records a finding as fixed and the next round is invoked
- **THEN** the delivered prompt presents that as a claim to verify against the code, not as a statement that the finding was resolved

#### Scenario: A decision arrives as a decision

- **WHEN** a decision is recorded against a finding that asked a person and the next round is invoked
- **THEN** the delivered prompt carries what was decided, attributes it to a person rather than to a patch, and asks the reviewer to judge whether the code matches it
