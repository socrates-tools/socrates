## ADDED Requirements

### Requirement: What happens after a round is one pure rule

The state a run reaches after a completed round SHALL be a pure function of the classified findings and the round number:

- any finding whose effective action is `ask-user`: `escalated`
- nothing blocking: `clean`
- nothing blocking that a response could address: `escalated`
- the round number equals the cap: `escalated`
- otherwise: `awaiting-response`

Whether a document blocks SHALL have one owner, which is the classification's verdict.

A round SHALL park only where it holds a finding that both holds the change and can be patched, and that SHALL be the same set the run waits on responses for.
A round that parks with nothing for the author to answer asks a question nobody can close and then escalates at the cap for findings whose own action asked for nothing.
A finding that holds the change while its action asks for nothing is a contradiction in what the reviewer returned, and it goes to the person rather than being demanded an answer socrates presents as unnecessary.

The reason a round escalated SHALL name the cap only where the cap is why it escalated.
A run that concluded `clean` naming the cap tells the author that something is left for them when nothing is.

The rule SHALL be testable against its whole input space rather than sampled, and its package SHALL reach no filesystem, no clock and no process.

The `ask-user` case SHALL dominate the others.
A set holding one question about intent and five mechanical defects goes to the person with all six, because the person's answer to the sixth may change what the other five should have been.

#### Scenario: A clean round finishes the run

- **WHEN** a round completes with no blocking and no patchable findings
- **THEN** the run is `clean` and no further round is started

#### Scenario: A question about intent stops the loop at once

- **WHEN** a round completes with one `ask-user` finding and several `auto-fix` findings
- **THEN** the run is `escalated`, all the findings are handed to the person, and none of them is offered for a response

#### Scenario: An informational round finishes the run

- **WHEN** a round completes with patchable findings that none of them block
- **THEN** the run is `clean` rather than parked, because nothing is being asked of the author

#### Scenario: A blocking finding nobody can patch goes to the person

- **WHEN** a round completes with a finding that holds the change and carries the action `no-op`
- **THEN** the run is `escalated`, and the run waits on no response to that finding

#### Scenario: A clean second round does not name the cap

- **WHEN** the final permitted round completes with nothing blocking
- **THEN** the run is `clean` and its recorded reason does not name the round cap

#### Scenario: Mechanical findings park for the author

- **WHEN** a round completes with blocking findings that are all `auto-fix`
- **THEN** the run is `awaiting-response`

#### Scenario: The cap escalates

- **WHEN** the final permitted round completes with blocking findings
- **THEN** the run is `escalated`, naming that the round cap was reached

#### Scenario: The rule reaches nothing outward

- **WHEN** the package holding the rule is examined
- **THEN** it imports nothing that reads the filesystem, the network, the clock or the process environment

### Requirement: The round cap is two, for context rather than for cost

socrates SHALL permit at most two rounds in a run.

The reason SHALL be recorded where the cap is: the session being handed findings is already large from implementing the change, and each round adds findings, patches and rebuttals to it.
Round three is where that session starts compacting mid-review and quietly loses the thing it was defending.

The cap SHALL NOT be configurable.
A number chosen for a reason that is not cost is a number somebody raises the moment they meet a run that nearly finished.

#### Scenario: A second round is permitted

- **WHEN** a run's first round parks and responses are recorded
- **THEN** a second round starts

#### Scenario: A third round is not

- **WHEN** a run's second round completes with blocking findings
- **THEN** the run escalates rather than starting a third round, and says the cap was reached

#### Scenario: No setting changes it

- **WHEN** the configuration file is examined
- **THEN** it holds no round cap

### Requirement: A response is a fix or a rejection, and a rejection needs a reason

`scrat respond` SHALL record, against one finding identifier, either that it was fixed, with an optional note, or that it is rejected, with a reason.

A rejection's reason SHALL be required and SHALL NOT be accepted empty or blank.
A response naming a finding the round did not produce SHALL be refused, naming the identifier.
A second response to the same finding SHALL replace the first, and the replacement SHALL be recorded so that changing one's mind is visible rather than silent.

#### Scenario: Recording a fix

- **WHEN** `scrat respond` records a finding as fixed
- **THEN** the response is stored against that identifier and the run stays `awaiting-response` until every blocking finding has one

#### Scenario: Recording a rejection

- **WHEN** `scrat respond` records a finding as rejected with a reason
- **THEN** the reason is stored verbatim against that identifier

#### Scenario: A rejection with no reason

- **WHEN** `scrat respond` records a rejection with an empty or blank reason
- **THEN** it is refused, saying that a rejection is a claim and needs one, and nothing is stored

#### Scenario: A response to a finding that does not exist

- **WHEN** `scrat respond` names an identifier the round did not produce
- **THEN** it is refused naming the identifier, and nothing is stored

#### Scenario: Changing a response

- **WHEN** a finding already responded to is responded to again
- **THEN** the later response is what the next reviewer receives, and the run records that it was replaced

### Requirement: A recorded response is not lost

The stored responses of a round SHALL be append-only, and recording one SHALL NOT rewrite what an earlier one wrote.

Two invocations recording at the same time is the ordinary case rather than the exotic one, because the author is a session that is encouraged to run independent commands at once.
A read, an append in memory and a rewrite of the whole document loses whichever claim was written first, silently, and the claim is the entire value of a response.
Appending also bounds what a crash can damage to the line being written, where a rewrite can leave a document that neither `scrat respond` nor `scrat continue` can read, on a run that cannot be advanced any other way.

#### Scenario: Two answers do not overwrite each other

- **WHEN** several responses to a round are recorded concurrently
- **THEN** every response that was recorded is in the stored document

#### Scenario: An answer is appended rather than rewritten

- **WHEN** a response replaces an earlier response to the same finding
- **THEN** the stored document still holds the earlier one, byte for byte, ahead of the later one

### Requirement: Silence does not advance a round

`scrat continue` SHALL refuse while any finding of the current round that holds the change and can be patched has no response, naming which findings are still unanswered.

Silence is neither a fix nor a rejection.
A round that advanced on it would present an unpatched, unargued finding to the next reviewer as though it had been dealt with, which is the one outcome worse than not asking.

#### Scenario: Continuing too early

- **WHEN** `scrat continue` is run while two blocking findings have no response
- **THEN** it refuses, names both, and starts no round

#### Scenario: Continuing when everything is answered

- **WHEN** every blocking finding has a response
- **THEN** `scrat continue` takes a fresh snapshot and starts the next round

#### Scenario: Nothing left to review ends the run

- **WHEN** the author answers every finding by reverting the change, so the working tree holds nothing the run's base does not
- **THEN** the run ends `clean` for the same reason a first round would, rather than `failed`

#### Scenario: Non-blocking findings need no response

- **WHEN** the round produced informational findings alongside blocking ones and only the blocking ones have responses
- **THEN** `scrat continue` proceeds

### Requirement: The next round is a fresh process that inherits no session

Each round's reviewer SHALL be a new process, and SHALL NOT resume, continue or otherwise inherit any prior round's session.

Round `n`'s patches implement round `n-1`'s findings, so a reviewer resuming that round's session would be seating the prescriber of those patches as their certifier.
It is the mechanism behind a real shipped defect elsewhere, where one round wrote both wrong code and the test blessing it and the resumed reviewer passed them.

The cross-round context a re-review legitimately needs SHALL travel in the explicit, labelled block of prior findings and responses, and nowhere else.

#### Scenario: No session is carried

- **WHEN** a second round is invoked
- **THEN** the invocation carries no session identifier, no resume instruction and no reference to the first round's conversation

#### Scenario: The claims travel explicitly

- **WHEN** a second round is invoked
- **THEN** the delivered prompt carries the prior findings and the recorded responses, and they are the only cross-round context in it

### Requirement: Every round's artifacts are kept

Each round SHALL keep its own prompt, findings, responses and reviewer output, under its own directory, and an earlier round's artifacts SHALL NOT be overwritten or removed by a later one.

What the first reviewer said, what the author claimed about it, and what the second reviewer made of that is the whole record of the disagreement, and it is the only place it exists.

#### Scenario: Two rounds, two directories

- **WHEN** a run completes two rounds
- **THEN** both round directories exist, each with its own prompt, findings and reviewer output, and the first round's responses are under the first round

#### Scenario: An escalated run keeps everything

- **WHEN** a run escalates at its second round
- **THEN** every artifact of both rounds is readable

### Requirement: socrates tells the author the rules when it hands over findings

Where a run is `awaiting-response`, `scrat findings` SHALL print, alongside the findings, the rules the author is being asked to work under:

- The author did not review this code; an independent reviewer did.
- For each finding, either apply a fix or record a rejection with a concrete reason.
- A rejection is a claim being made to the next reviewer, not a decision.
- Do not resolve a finding by weakening or deleting behaviour that was added on purpose. Where the finding is right about a symptom and wrong about the cause, say so and fix the cause.

socrates does not run the author's session and cannot constrain it.
Putting the rules in front of it, at the moment it has findings in hand, in the output of the command it reads, is the whole of what is available.

#### Scenario: The rules arrive with the findings

- **WHEN** `scrat findings` is run against a run awaiting responses
- **THEN** it prints the four rules alongside the findings

#### Scenario: The rules do not arrive where they do not apply

- **WHEN** `scrat findings` is run against a run that is clean, escalated or failed
- **THEN** it prints the findings without the response rules, because no response is being asked for
