## MODIFIED Requirements

### Requirement: The reviewer is told what it is looking at and nothing about why

The prompt SHALL carry the base commit, the snapshot commit, the changed-file list, and the person's stated intent where one was given.

It SHALL NOT carry the authoring session's account of what it did, its reasoning, its summary of its own change, or its confidence.
Piping the author's rationale to the reviewer contaminates the independence the whole design pays for, and an author's account of a change is exactly the thing a reviewer is there to not take on trust.

There is exactly one exception, and it is explicit rather than a softening of the rule.
A round after the first SHALL carry the prior rounds' findings and the author's recorded responses to them, in a block labelled as claims to be judged.
The exception is stated as an exception so that a second one has to be argued for rather than assumed.

#### Scenario: Intent reaches the reviewer

- **WHEN** a review is started with a stated intent
- **THEN** the delivered prompt carries that intent, marked as the author's stated intent rather than as fact

#### Scenario: No intent is stated

- **WHEN** a review is started without a stated intent
- **THEN** the delivered prompt carries no intent section at all, rather than an empty one

#### Scenario: The author's rationale has no channel

- **WHEN** the command surface is examined for a way to pass the authoring session's own account of the change
- **THEN** there is none, and the only free text a caller may supply is the stated intent and a response attached to a finding

#### Scenario: The first round carries no claims

- **WHEN** the first round is invoked
- **THEN** the delivered prompt carries no prior-round block, because there is nothing to judge yet

#### Scenario: A later round carries the claims and nothing more

- **WHEN** a second round is invoked
- **THEN** the delivered prompt carries the prior findings and the recorded responses, and carries no other account of the change from the author's session

## ADDED Requirements

### Requirement: A re-review carries the provenance clause

A round after the first SHALL carry a clause stating:

- The code under re-review was written by the change's author in response to findings. It is unreviewed new code, not a settled resolution of the findings that prompted it, and it is reviewed to exactly the same adversarial standard as the original change.
- Prior findings, claimed fixes and rejections are claims, not evidence. Verify each against the current code, and judge independently whether what the patches introduced is correct, not merely whether it does what was prescribed.
- A test written or changed in the same round as the code it exercises is part of that round's claim, not independent proof. Judge whether its asserted outcome is the right outcome and whether it could still pass with the code wrong.

The clause SHALL be absent from the first round, where there is no prior claim for it to be about.

#### Scenario: The clause reaches a re-review

- **WHEN** a second round is invoked
- **THEN** the delivered prompt carries each of the three statements in this requirement

#### Scenario: The clause is absent from a first review

- **WHEN** the first round is invoked
- **THEN** the delivered prompt carries none of them

#### Scenario: A claimed fix is not taken on trust

- **WHEN** the author records a finding as fixed and the next round is invoked
- **THEN** the delivered prompt presents that as a claim to verify against the code, not as a statement that the finding was resolved
