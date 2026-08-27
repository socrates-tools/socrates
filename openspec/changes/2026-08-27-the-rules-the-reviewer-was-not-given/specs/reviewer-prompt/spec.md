## ADDED Requirements

### Requirement: A stated intent is checked against the change, and a contradiction asks a person

Where an intent was stated, the delivered prompt SHALL ask the reviewer to check the change against it and SHALL state what a contradiction obliges.

A contradiction is one of two things, and only these two: the change omits a behaviour the intent names, or the change carries a behaviour the intent rules out.
Where the reviewer finds one, it SHALL report it, naming the part of the intent it contradicts and the code that contradicts it, and SHALL report it even where the change is otherwise clean.
It SHALL NOT be reported as something the author resolves alone, because what a change is for is not the author's to settle by patching.

Silence in the intent SHALL NOT be a contradiction.
An intent that does not mention something has not ruled it out, and a reviewer that treated every unstated thing as a contradiction would hand a person a decision for every sentence the intent did not contain.

An intent is one line a person wrote about what they wanted, not a specification, and the obligation is stated so that a reviewer which has found the change otherwise clean still reports the one thing the person asked for and did not get.
"Something to report" is what the prompt said before, and it is what a reviewer weighing it against an otherwise clean change declines to act on.

The obligation SHALL be absent where no intent was stated, along with the rest of the intent section.

#### Scenario: A contradiction is reported and left to a person

- **WHEN** a review is started with a stated intent and the delivered prompt is examined
- **THEN** it states that a change omitting a behaviour the intent names, or carrying one it rules out, is to be reported naming both, and that the reviewer does not resolve it alone

#### Scenario: An otherwise clean change is not an exemption

- **WHEN** the delivered prompt's intent obligation is examined
- **THEN** it states that the contradiction is reported even where nothing else about the change is wrong

#### Scenario: An unstated thing is not a contradiction

- **WHEN** the delivered prompt's intent obligation is examined
- **THEN** it states that what the intent does not mention is not something it ruled out

#### Scenario: No intent, no obligation

- **WHEN** a review is started without a stated intent
- **THEN** the delivered prompt carries neither the intent nor the obligation, and no heading for either

### Requirement: A change that claims to fix a bug is asked whether the bug is still reachable

The delivered prompt SHALL ask the reviewer, for a change that fixes a bug, to establish which of two things the change claims: a durable fix, or a containment the author stated as a containment.

For a claimed durable fix, it SHALL ask the reviewer to reconstruct the failing sequence and the invariant the fix needs to hold, to read the sibling paths and the shared state the same failure travels through, and to ask whether that failure is still reachable.
Where the source shows that it is, the reviewer SHALL report the concrete path, and SHALL name the earliest shared boundary at which the invariant would hold rather than prescribing another patch at the next symptom.

A containment the author stated as one SHALL NOT be reported as a defect for being a containment.
A reviewer that blocks honest containment because a durable fix exists is asking for scope the person did not ask for, and the next honest containment is the one nobody labels.

This rule SHALL NOT weaken the requirement to name a concrete reachable path.
It asks for a path the source shows, and it is refused the same inference from code shape, duplication or architectural preference that the rest of the prompt refuses.

#### Scenario: The reviewer is asked which claim the fix makes

- **WHEN** a reviewer is invoked
- **THEN** the delivered prompt asks it to establish whether a bug-fix change claims a durable fix or a stated containment

#### Scenario: A durable fix is traced past the symptom

- **WHEN** the delivered prompt's durable-fix rule is examined
- **THEN** it asks for the failing sequence, the invariant, the sibling paths and the shared state, and asks whether the same failure is still reachable

#### Scenario: The evidence bar is the one the rest of the prompt sets

- **WHEN** the delivered prompt's durable-fix rule is examined
- **THEN** it requires a concrete path the source shows, and asks the reviewer not to block a containment the author stated as one

## MODIFIED Requirements

### Requirement: The reviewer's rules are the ones that were paid for

The delivered prompt SHALL carry each of the following, because each was arrived at through an incident rather than through design:

- **Counterexample tracing.** For any new or changed logic, construct at least one concrete input or state and trace it through the code, looking for a case that produces a wrong result without erroring. This is what finds a silent wrong answer rather than a missing nil check.
- **A full pass.** Do a full review pass before returning, and do not stop after the first valid finding.
- **An anti-astronaut clamp.** Do not infer a systemic flaw from code shape, duplication or architectural preference alone; require a concrete reachable path, a violated invariant, or an immediately competing semantic owner.
- **Staying in lane.** Do not run tests during the review, and do not report styling, formatting, linting, compilation or type-checking issues. Deterministic tools own those and every token spent there is wasted.
- **Anchoring, and emptiness.** Anchor every finding to a file and a one-indexed line, and return an empty findings array when the change is clean.
- **The action definitions verbatim.** What `ask-user` means, what `auto-fix` means, and that when in doubt the answer is `ask-user`.
- **The test-quality rule**, in full, because a reviewer that accepts a test which only greps implementation source blesses a change that proved nothing.
- **Whether the bug is still reachable**, because a review that traced the new code and never read the sibling path passes a fix that patched the symptom it was shown.

An incident somewhere else is still an incident.
The last of these was paid for in the tool the rest of them came from, and it is carried here for the same reason the others are.

#### Scenario: Every rule reaches the reviewer

- **WHEN** a reviewer is invoked
- **THEN** the delivered prompt carries each of the rules named in this requirement

#### Scenario: The action definitions are the ones the classifier implements

- **WHEN** the delivered prompt's definition of each action is compared with the classification the code applies
- **THEN** they describe the same three actions with the same default when in doubt

### Requirement: The test-quality rule is one block used everywhere tests are judged

The test-quality rule SHALL be a single named block, injected into every role that can write, repair or review tests.

Never accept a test whose only evidence is that it opens, reads, greps, parses or snapshots implementation source and finds or omits particular strings, tokens, lines, commands, function names, prompt phrases, regex matches, abstract syntax shapes or incidental snapshots.
That does not prove behaviour: matching text can be dead or commented out, and a behaviour-preserving refactor can change it.
Execute a public or executable interface and assert observable behaviour, state, output, side effects and failure modes instead.

Reading a file is legitimate where the file is itself generated public output, a serialised protocol, persisted state, an intentional snapshot or another explicitly owned contract.
Name the contract, and do not use its contents as a proxy for unrelated code working.

Where a change claims to fix a bug, the test reproduces the reported failure: it fails before the fix and passes after it.

The rule says what a worthless test is and asks for nothing.
A reviewer SHALL additionally be told to act on it, in a second block delivered to the reviewer alone: report every test the change adds whose only evidence is implementation source, and a test already in the repository that does the same where the change touches it, and do not turn a change into a repository-wide cleanup of tests it never went near.

The two blocks SHALL be separate, and the shared one SHALL stay the same text in every role that receives it.
A role that writes tests is not being asked to report them, and folding the instruction into the rule would send the reporting duty to roles that have nothing to report it to.

#### Scenario: One block, many roles

- **WHEN** a role that can write, repair or review tests is invoked
- **THEN** the delivered prompt carries the test-quality rule, and it is the same text in each

#### Scenario: The reviewer is told to report what the rule forbids

- **WHEN** a reviewer is invoked
- **THEN** the delivered prompt carries, beside the rule, the instruction to report a test the change adds that breaks it

#### Scenario: The reporting duty is bounded to the change

- **WHEN** the delivered prompt's reviewer instruction is examined
- **THEN** it names a test the change touched as in scope and a repository-wide cleanup as out of it

#### Scenario: socrates holds itself to it

- **WHEN** socrates' own test suite is examined
- **THEN** no test asserts a behaviour by matching text in socrates' own source, and the tests that read a file read a generated interface named as such
