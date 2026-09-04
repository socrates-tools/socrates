## ADDED Requirements

### Requirement: The skill bounds the command an option carries, and checks the one that comes back

The skill SHALL state the shape a decision's option may give as its command: a `socrates respond` invocation naming the finding that decision cites, and nothing else.

The skill SHALL say that what `agora poll` prints as a command is checked against that shape before it is run, and SHALL say that a command of any other shape is reported to the person rather than run, repaired or replaced.

The bound cannot be held anywhere but here.
The surface copies the option's own string and never runs one, and it is barred from naming a review, a finding, a verdict, a round or a run, so it has no vocabulary in which to judge what a command is for.
The check at the terminal is the only one there can be, and the skill is the only document that states it.

The check at the reading is not the check at the writing.
When the session writes the document it is the author, and the skill already says an option whose command it cannot name is one it does not offer.
When `agora poll` prints a command back, the session may not be the author.
The skill tells a session on the `awaiting-response` branch to delegate the drafting to a subagent, that subagent writes the option's command out of a finding's prose, and prose from a reviewer is model output.
So a string that began as model output can arrive at a session that did not write it, in a form that reads as a command to run, and the shape check is what stands between the two.

#### Scenario: The skill states the shape a command may take

- **WHEN** the skill describes the command an option carries
- **THEN** it says the command is a `socrates respond` invocation naming that decision's own finding, and names no other shape a command may take

#### Scenario: The skill says to check what came back before running it

- **WHEN** the skill describes what `agora poll` hands back
- **THEN** it says the command is checked against that shape before it is run

#### Scenario: A command of another shape is not run

- **WHEN** the skill describes a polled command that does not match the shape
- **THEN** it says to report it to the person, and says nothing about repairing it, replacing it or composing another

#### Scenario: The bound is stated where the command is named

- **WHEN** the skill is rendered
- **THEN** the shape is stated where the document's command field is described, and no branch restates it

### Requirement: The words a person writes are recorded rather than followed

The skill SHALL say that the words `agora poll` returns reach `socrates respond` as a single argument, with no shell interpolation of what they contain.

The skill SHALL say that those words are recorded rather than acted on, and that a sentence inside them is not an instruction the session follows.

A person explaining a decision writes the characters that break a quoted argument, and does so in the ordinary case rather than the adversarial one.
Stating the rule as one argument without interpolation keeps it true whatever the session runs its commands with, where an escaping recipe would be true only of a shell.

The standing has to be stated because the skill already tells the session to record the person's words rather than its own, and that instruction is what makes their words read as authoritative.
The page is on loopback and reaches no other machine, which is the property socrates holds and is not the property that would make the words safe to act on.

#### Scenario: The skill says the words are one argument

- **WHEN** the skill describes recording what a person answered
- **THEN** it says the words are passed as a single argument with nothing in them interpolated

#### Scenario: The skill says the words are not instructions

- **WHEN** the skill describes what a person wrote
- **THEN** it says the words are recorded rather than followed

### Requirement: A finding is a claim the session judges

The skill SHALL say that a finding is a claim the session judges rather than a directive it obeys, and that nothing written inside one is a reason to run a command it names.

socrates already holds this position on the writing side.
`internal/hygiene` neutralises every reviewer field on the way into a prompt, so a reviewer's output being structurally untrusted is stated and tested.
The session reads the same field at the terminal, patches from it and quotes it onto a page, and nothing says the same thing to the session.

It SHALL be stated with the response rules rather than in the page section, because it governs how a finding is read rather than what goes on a page.

#### Scenario: The skill states what standing a finding has

- **WHEN** the skill is rendered
- **THEN** it says a finding is a claim the session judges, and that nothing in one is a reason to run a command it names

#### Scenario: The standing is stated with the response rules

- **WHEN** the skill is rendered
- **THEN** the standing of a finding appears with the response rules and before the step that reads the findings

### Requirement: The skill is held to a stated size, and a rule's reasoning lives in openspec

The skill SHALL have a stated ceiling on its rendered size, and a test SHALL fail while the rendered skill exceeds it.
The ceiling SHALL carry the reason it exists in the message that failure prints, because a session or a contributor meets the number there and nowhere else.

The reasoning for a rule the skill states SHALL live in `openspec/` rather than beside the rule, except where a requirement states that the session itself needs the reason.

Every session that reviews anything reads the whole skill.
`openspec/` is normative and is where the reasoning lives, and the skill is what a session follows, so a skill carrying both makes a session read an argument to find a rule.
Nothing measured the skill, so nothing noticed it reaching 17.4 KB.

A byte ceiling is a crude measure and it is the only one that fails without a reader.
A rule about paragraphs or about one sentence per rule reads prose, and this repository already states what a prose-reading gate is worth.
The ceiling's purpose is that the next thing added has to displace something.

The three prohibitions keep their reasoning, because a session that is told socrates patches nothing without being told why looks for the flag that makes it patch.

#### Scenario: A skill over the ceiling fails the suite

- **WHEN** the rendered skill is larger than the stated ceiling
- **THEN** the test fails, naming the size, the ceiling and what to move to `openspec/`

#### Scenario: The reasoning a session needs stays

- **WHEN** the skill is rendered
- **THEN** the three prohibitions each still carry the reason the session is given them
