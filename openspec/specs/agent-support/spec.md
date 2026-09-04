# agent-support Specification

## Purpose
TBD - created by archiving change reviewer-cli-support. Update Purpose after archive.

## Requirements

### Requirement: socrates states which reviewer CLIs it drives, and the basis on which one qualifies

socrates SHALL state the reviewer CLIs it supports in exactly one place in the source, each stated beside the facts that make it qualify rather than as a bare name on a list.

A CLI SHALL qualify on all of:

- The reviewed repository's own instruction files can be kept out of the reviewer's session.
- It reports, on a terminal event, whether the work happened, and where it did not, what it says about why.
- Its model is selectable per invocation, so that a role resolves to its own model.

The criteria are what socrates cannot repair after the fact.
A review conducted by a session that read the reviewed repository's own instructions is defective in its judgement rather than in its format, and nothing downstream recovers it.
A failure the CLI does not report is a run a person cannot act on.
A model that cannot be selected defeats per-role routing, which exists so that the reviewer gets the strongest model available.

Whether the CLI constrains its output to a supplied schema SHALL NOT be a criterion.
socrates validates the findings document against the schema itself, and does so on both tiers, so a CLI without that capability delivers a document that is checked by the same rule.
It SHALL be stated beside the criteria as the tier that CLI is on.

Whether a CLI reports the inventory of tools its session holds SHALL be stated beside the criteria and SHALL NOT be one of them.
socrates records what it observes about the tool surface and never fails a round for it, so a CLI that reports no inventory is one whose runs record that.

A CLI that does not qualify SHALL be stated in the same place with the criterion it fails, rather than omitted.
An omission states that socrates does not support a CLI without stating why, and the refusal message, what a person reads and a reader of the source all need the reason.
A CLI whose vendor ships the missing capability is then supported by doing the work that capability demands, rather than by repeating the survey.

No second copy of the supported set SHALL exist in the source.
What refuses an unsupported CLI, what constructs an implementation, and what a person reads SHALL each read the one statement.

#### Scenario: The supported set is stated once

- **WHEN** the source is searched for the set of reviewer CLIs socrates supports
- **THEN** exactly one statement of it exists, each entry carries the facts that make it qualify, and everything that names a supported CLI reads that statement

#### Scenario: An unsupported CLI is stated with the reason it is unsupported

- **WHEN** the statement is read for a CLI socrates does not support
- **THEN** it is present, and it names the criterion that CLI does not meet

#### Scenario: A CLI that constrains no schema still qualifies

- **WHEN** an entry states that its vendor constrains no schema and meets every criterion
- **THEN** it is among the supported CLIs, on the tier that says its document is extracted and validated

#### Scenario: Support follows from the criteria rather than from a list

- **WHEN** an entry states that the reviewed repository's instructions cannot be kept out of its session
- **THEN** it is not among the supported CLIs, and nothing has to be edited in a second place for that to be so

### Requirement: A findings document arrives on one of two tiers and reaches one validation

socrates SHALL obtain a reviewer's findings document in one of two ways, stated per vendor.

Where the vendor constrains output against a supplied schema, the document SHALL be that structured output.

Where it does not, the reviewer SHALL be asked to end its final message with the document, and socrates SHALL extract it from that message.

Both SHALL then be validated against the same schema by the same check.
A vendor's schema flag constrains what a model may emit and does not perform socrates' validation, and socrates SHALL NOT skip its validation because that flag was passed.

A document that is absent, that does not parse, or that does not satisfy the schema SHALL fail the call on either tier.
The verdict rule is unchanged: `clean` requires a findings document that parsed, and a round whose call failed is `failed`.

Extraction SHALL take the last candidate in the final message.
A reviewer that reasons before answering quotes the code it is discussing, and an earlier block in the same message is that quotation rather than the answer.

Nothing outside socrates SHALL shape a reviewer's answer into a findings document.
The main coding session SHALL NOT be asked to, because that session is the author.
A document the author transcribed is one in which the author chose the severity, chose the action and chose what to omit, and the default that sends an unclassified finding to a person rather than to a patch is unreachable in a document the author wrote.

#### Scenario: The enforced tier delivers the vendor's structured output

- **WHEN** a review runs under a vendor that constrains output against the supplied schema
- **THEN** the findings document is that structured output, and it is validated against the schema before it is used

#### Scenario: The parsed tier delivers the document from the final message

- **WHEN** a review runs under a vendor that constrains no schema and the reviewer ends its final message with a findings document
- **THEN** the document is extracted from that message and validated against the same schema, and the round proceeds as it would on the other tier

#### Scenario: An unparseable document fails on either tier

- **WHEN** a reviewer on either tier returns no document, a document that does not parse, or one that violates the schema
- **THEN** the call fails, the round is `failed`, and no verdict is read from it

#### Scenario: Reasoning before the answer is not mistaken for the answer

- **WHEN** a reviewer's final message quotes a candidate document earlier in the message and ends with the real one
- **THEN** the one at the end is what is extracted

#### Scenario: The author never shapes the document

- **WHEN** the source and the delivered prompts are searched for what turns a reviewer's answer into a findings document
- **THEN** it is socrates, and nothing asks the session that authored the change to produce, reshape or complete one

### Requirement: A person names the reviewer CLI once, and an unsupported one is refused

Configuration SHALL carry the reviewer CLI in one field, and it SHALL default to the CLI socrates shipped with, so that an existing configuration and an absent one both keep working.

socrates SHALL NOT choose the CLI by looking at what is installed.
Which reviewer reads a person's unreviewed change is theirs to state, and a tool that picks one silently changes who reviewed the change when a person installs something unrelated.

A configured CLI socrates does not support SHALL be refused, naming it and naming the guarantee it does not offer, before a snapshot is taken and before a reviewer is spawned.
That is where an unsupported platform is refused, and for the same reason.
A refusal taken later has already frozen a commit or started a process against a guarantee socrates cannot hold.

A command that neither spawns nor reaps SHALL NOT be refused there, so a run recorded under any CLI stays readable.

#### Scenario: An absent field runs the CLI socrates shipped with

- **WHEN** no configuration file exists, or one exists naming no reviewer CLI
- **THEN** the run proceeds under the default CLI and nothing is refused

#### Scenario: A supported CLI is invoked

- **WHEN** configuration names a supported reviewer CLI
- **THEN** the reviewer is that CLI, and the run records which one answered

#### Scenario: An unsupported CLI is refused before anything happens

- **WHEN** configuration names a reviewer CLI socrates does not support and a review is asked for
- **THEN** the review is refused naming that CLI and the guarantee it does not offer, no snapshot is taken, and no reviewer process is started

#### Scenario: A finished run stays readable under any configuration

- **WHEN** configuration names an unsupported reviewer CLI and a command that neither spawns nor reaps is run
- **THEN** it is not refused, and it reports what the run holds

### Requirement: The run records which reviewer answered and how the document arrived

A run SHALL record the reviewer CLI it was carried out under and the tier that CLI is on, beside the model that answered.

A findings document says what a reviewer concluded and not which reviewer concluded it, and a person comparing two runs of the same change needs to know whether they are comparing rounds or vendors.
The tier is recorded because a round that failed for want of a document failed differently on each, and a person reading a run needs to know which.

A run recorded before socrates named a vendor SHALL be read as the CLI socrates shipped with, because that is what answered.

#### Scenario: The vendor and tier reach the run's record

- **WHEN** a round completes under any supported reviewer CLI
- **THEN** the run records that CLI and its tier beside the model that answered

#### Scenario: A run from before the field is read as the original CLI

- **WHEN** a run recorded without a reviewer CLI is read
- **THEN** it reads as the CLI socrates shipped with, and nothing fails for the absent field

### Requirement: What a person reads names the supported set by reading it

What a person reads before they decide whether socrates is usable SHALL name every supported reviewer CLI, and SHALL take that from the one statement rather than restating it.

A document naming one CLI where socrates supports three turns away a person who could have used it, which is the same failure `README.md` already corrected by stating its requirements at the top instead of in passing.

It SHALL state which tier each supported CLI is on, so that a person choosing between them knows which asks more of the reviewer.

It SHALL say which capability an unsupported CLI lacks, rather than only that it is unsupported.

#### Scenario: A document that names fewer CLIs than socrates supports fails

- **WHEN** the supported set and what a person reads disagree
- **THEN** the check fails, naming both

#### Scenario: A reader learns why their CLI is not supported

- **WHEN** what a person reads is read for a CLI socrates does not support
- **THEN** it names the capability that CLI does not offer
