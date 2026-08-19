## ADDED Requirements

### Requirement: The prompt delivered to a reviewer is a generated interface

The prompt socrates sends SHALL be treated as an artifact with a contract rather than as an implementation detail: it SHALL be recorded verbatim in the round directory, and the delivered prompt SHALL be what tests assert on.

A test SHALL NOT assert that a prompt rule exists by reading socrates' own source for a sentence.
Matching text can be dead or commented out, and a behaviour-preserving refactor changes it.
A test SHALL NOT assert what a model concluded from a prompt: whether a given wording finds a given bug is a development-time evaluation, never a check in the ordinary suite.

#### Scenario: The prompt is on disk

- **WHEN** a round completes
- **THEN** the round directory holds the exact prompt that was delivered, byte for byte

#### Scenario: Assertions are made on the delivered prompt

- **WHEN** a test asserts that a rule reached the reviewer
- **THEN** it reads the prompt the simulator received, not socrates' source

### Requirement: The reviewer is told what it is looking at and nothing about why

The prompt SHALL carry the base commit, the snapshot commit, the changed-file list, and the person's stated intent where one was given.

It SHALL NOT carry the authoring session's account of what it did, its reasoning, its summary of its own change, or its confidence.
Piping the author's rationale to the reviewer contaminates the independence the whole design pays for, and an author's account of a change is exactly the thing a reviewer is there to not take on trust.

#### Scenario: Intent reaches the reviewer

- **WHEN** a review is started with a stated intent
- **THEN** the delivered prompt carries that intent, marked as the author's stated intent rather than as fact

#### Scenario: No intent is stated

- **WHEN** a review is started without a stated intent
- **THEN** the delivered prompt carries no intent section at all, rather than an empty one

#### Scenario: The author's rationale has no channel

- **WHEN** the command surface is examined for a way to pass the authoring session's own account of the change
- **THEN** there is none, and the only free text a caller may supply is the stated intent

### Requirement: The diff is inlined below a ceiling and read above it

The reviewer's working directory SHALL always be the snapshot worktree, so that it can inspect a call site the diff does not show.

Where the diff between the base commit and the snapshot is smaller than a configured byte ceiling, the prompt SHALL carry the diff itself.
Where it is larger, the prompt SHALL carry the changed-file list and the two commits, and SHALL say that the diff is to be read rather than supplied.

Which of the two happened SHALL be recorded in the run, because it is the largest single factor in what a round costs and how long it takes.

#### Scenario: A small change carries its diff

- **WHEN** the diff is under the ceiling
- **THEN** the delivered prompt contains the diff, and the run records that the diff was inlined

#### Scenario: A large change is read

- **WHEN** the diff is over the ceiling
- **THEN** the delivered prompt contains the changed-file list and the two commits rather than the diff, and the run records that the diff was not inlined

#### Scenario: The ceiling is one setting

- **WHEN** the ceiling is configured
- **THEN** it applies to every repository and every run, and there is no per-repository override

### Requirement: The reviewer's rules are the ones that were paid for

The delivered prompt SHALL carry each of the following, because each was arrived at through an incident rather than through design:

- **Counterexample tracing.** For any new or changed logic, construct at least one concrete input or state and trace it through the code, looking for a case that produces a wrong result without erroring. This is what finds a silent wrong answer rather than a missing nil check.
- **A full pass.** Do a full review pass before returning, and do not stop after the first valid finding.
- **An anti-astronaut clamp.** Do not infer a systemic flaw from code shape, duplication or architectural preference alone; require a concrete reachable path, a violated invariant, or an immediately competing semantic owner.
- **Staying in lane.** Do not run tests during the review, and do not report styling, formatting, linting, compilation or type-checking issues. Deterministic tools own those and every token spent there is wasted.
- **Anchoring, and emptiness.** Anchor every finding to a file and a one-indexed line, and return an empty findings array when the change is clean.
- **The action definitions verbatim.** What `ask-user` means, what `auto-fix` means, and that when in doubt the answer is `ask-user`.
- **The test-quality rule**, in full, because a reviewer that accepts a test which only greps implementation source blesses a change that proved nothing.

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

#### Scenario: One block, many roles

- **WHEN** a role that can write, repair or review tests is invoked
- **THEN** the delivered prompt carries the test-quality rule, and it is the same text in each

#### Scenario: socrates holds itself to it

- **WHEN** socrates' own test suite is examined
- **THEN** no test asserts a behaviour by matching text in socrates' own source, and the tests that read a file read a generated interface named as such

### Requirement: The prompt is assembled from named sections in a fixed order

The prompt SHALL be composed of named sections, assembled in a fixed order, so that a section can be added, changed or made conditional without any other section moving.

A section that has nothing to say SHALL be absent rather than empty, so that an unstated intent does not deliver a heading with nothing under it.

#### Scenario: Order is stable

- **WHEN** two reviews are run against the same change with the same configuration
- **THEN** the delivered prompts are byte-identical

#### Scenario: An empty section is absent

- **WHEN** a section has nothing to contribute
- **THEN** neither its heading nor a placeholder appears in the delivered prompt
