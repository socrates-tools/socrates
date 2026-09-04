## MODIFIED Requirements

### Requirement: A recorded response is not lost

The stored responses of a round SHALL be append-only, and recording one SHALL NOT rewrite what an earlier one wrote.

Two invocations recording at the same time is the ordinary case rather than the exotic one, because the author is a session that is encouraged to run independent commands at once.
A read, an append in memory and a rewrite of the whole document loses whichever claim was written first, silently, and the claim is the entire value of a response.
Appending also bounds what a crash can damage to the line being written.
A rewrite can leave a document that neither `socrates-review respond` nor `socrates-review continue` can read, on a run that cannot be advanced any other way.

#### Scenario: Two answers do not overwrite each other

- **WHEN** several responses to a round are recorded concurrently
- **THEN** every response that was recorded is in the stored document

#### Scenario: An answer is appended rather than rewritten

- **WHEN** a response replaces an earlier response to the same finding
- **THEN** the stored document still holds the earlier one, byte for byte, ahead of the later one

### Requirement: Silence does not advance a round

`socrates-review continue` SHALL refuse while any finding of the current round that holds the change and can be answered has no response, naming which findings are still unanswered and which of them need a person.

Silence is neither a fix, nor a rejection, nor a decision.
A round that advanced on it would present an unpatched, unargued finding to the next reviewer as though it had been dealt with, which is the one outcome worse than not asking.
A round that advanced past an unanswered question would do worse: it would carry a question about what the change is for into a round that judges patches, and nothing would ever answer it.

#### Scenario: Continuing too early

- **WHEN** `socrates-review continue` is run while two blocking findings have no response
- **THEN** it refuses, names both, and starts no round

#### Scenario: Continuing when everything is answered

- **WHEN** every blocking finding has a response
- **THEN** `socrates-review continue` takes a fresh snapshot and starts the next round

#### Scenario: Continuing with a question unanswered

- **WHEN** `socrates-review continue` is run while a finding that asks a person has no decision recorded against it
- **THEN** it refuses, names that finding as one needing a person, and starts no round

#### Scenario: Nothing left to review ends the run

- **WHEN** the author answers every finding by reverting the change, so the working tree holds nothing the run's base does not
- **THEN** the run ends `clean` for the same reason a first round would, rather than `failed`

#### Scenario: Non-blocking findings need no response

- **WHEN** the round produced informational findings alongside blocking ones and only the blocking ones have responses
- **THEN** `socrates-review continue` proceeds

### Requirement: socrates tells the author the rules when it hands over findings

Where a run is `awaiting-response`, `socrates-review findings` SHALL print, alongside the findings, the rules the author is being asked to work under:

- The author did not review this code; an independent reviewer did.
- For each finding, either apply a fix or record a rejection with a concrete reason.
- A rejection is a claim being made to the next reviewer, not a decision.
- A finding that asks a person is the person's to answer. Take it to them, and record what they decided.
- Do not resolve a finding by weakening or deleting behaviour that was added on purpose. Where the finding is right about a symptom and wrong about the cause, say so and fix the cause.

socrates does not run the author's session and cannot constrain it.
Putting the rules in front of it, at the moment it has findings in hand, in the output of the command it reads, is the whole of what is available.

#### Scenario: The rules arrive with the findings

- **WHEN** `socrates-review findings` is run against a run awaiting responses
- **THEN** it prints the rules alongside the findings, including the one saying whose a finding that asks a person is

#### Scenario: The rules do not arrive where they do not apply

- **WHEN** `socrates-review findings` is run against a run that is clean, escalated or failed
- **THEN** it prints the findings without the response rules, because no response is being asked for

### Requirement: A response is a fix, a rejection or a decision, and neither a rejection nor a decision is accepted without a reason

`socrates-review respond` SHALL record, against one finding identifier, either that it was fixed, with an optional note, or that it is rejected, with a reason, or that a person decided it, with a reason saying what was decided.

A rejection's reason SHALL be required and SHALL NOT be accepted empty or blank, and a decision's reason SHALL be required on the same terms.
A decision SHALL be the only resolution accepted for a finding whose action asks a person, and SHALL NOT be accepted for any other finding.
A fix or a rejection offered for a finding whose action asks a person SHALL be refused, naming the resolution that applies.

socrates cannot verify that a person rather than the session answered, and SHALL NOT pretend to.
It records the claim, labels it as a decision, and hands it to a reviewer that owes it nothing, which is exactly what it already does with a rejection.
What it does enforce is that answering such a finding cannot be done by asserting a patch, so a session that answers one on its own behalf has made a claim a reader can see rather than taken a silent shortcut.

A response naming a finding the round did not produce SHALL be refused, naming the identifier.
A second response to the same finding SHALL replace the first, and the replacement SHALL be recorded so that changing one's mind is visible rather than silent.

#### Scenario: Recording a fix

- **WHEN** `socrates-review respond` records a finding as fixed
- **THEN** the response is stored against that identifier and the run stays `awaiting-response` until every blocking finding has one

#### Scenario: Recording a rejection

- **WHEN** `socrates-review respond` records a finding as rejected with a reason
- **THEN** the reason is stored verbatim against that identifier

#### Scenario: A rejection with no reason

- **WHEN** `socrates-review respond` records a rejection with an empty or blank reason
- **THEN** it is refused, saying that a rejection is a claim and needs one, and nothing is stored

#### Scenario: Recording a decision

- **WHEN** `socrates-review respond` records a finding that asks a person as decided, with a reason
- **THEN** the reason is stored verbatim against that identifier, labelled as a decision rather than as a fix

#### Scenario: A decision with no reason

- **WHEN** `socrates-review respond` records a decision with an empty or blank reason
- **THEN** it is refused for the same reason a rejection is, and nothing is stored

#### Scenario: A patch offered where a decision is required

- **WHEN** `socrates-review respond` records a fix or a rejection against a finding whose action asks a person
- **THEN** it is refused, naming that the finding takes a decision, and nothing is stored

#### Scenario: A decision offered where none was asked for

- **WHEN** `socrates-review respond` records a decision against a finding whose action asks for a patch
- **THEN** it is refused, naming the resolutions that apply, and nothing is stored

#### Scenario: A response to a finding that does not exist

- **WHEN** `socrates-review respond` names an identifier the round did not produce
- **THEN** it is refused naming the identifier, and nothing is stored

#### Scenario: Changing a response

- **WHEN** a finding already responded to is responded to again
- **THEN** the later response is what the next reviewer receives, and the run records that it was replaced
