# findings-page Specification

## Purpose
TBD - created by archiving change a-round-a-person-reads-becomes-a-page. Update Purpose after archive.
## Requirements
### Requirement: The page is what a round a person reads is handed over as

Where a round is one a person reads, the session driving the loop SHALL hand it over as a page served locally by `agora`, give the person the address, and say briefly what it holds rather than pasting the findings.

The page SHALL NOT be published to any host outside the machine.
A findings document quotes the code it was written about, `internal/corpus` refuses to record one for that reason, and a skill that told a session to upload one contradicted that refusal.
Serving it locally is what makes the two rules agree.

The reviewer SHALL NOT draft it.
It is sealed read-only inside the snapshot worktree, it runs with the machine's own customisations disabled, and its findings are the untrusted thing under judgement.
A reviewer that also composed the person's summary of its own findings would be arguing its case in the document the person decides from.

socrates SHALL NOT draft it either.
What is being asked for is plain English about what a finding means for the change, which a template over the findings document cannot produce, and socrates' agent seam exists to review rather than to write for a person.

The page SHALL NOT replace what the session reads.
`socrates findings` keeps printing the findings, because the session patches from them, and the page is a reading aid rather than the record.

#### Scenario: A round a person reads is handed over as a page

- **WHEN** the loop reaches a round socrates reports as one a person reads
- **THEN** the session serves it as a page, hands over the address, and does not paste the findings into the terminal

#### Scenario: The page is served locally and published nowhere

- **WHEN** the skill describes handing a round over
- **THEN** it names the local surface, and names no host, upload, publish or artifact outside the machine

#### Scenario: Nothing socrates spawns drafts the page

- **WHEN** the prompt a reviewer is given and the roles socrates spawns are examined
- **THEN** none of them is asked for a summary written for a person, and nothing about a page reaches any prompt socrates builds

### Requirement: The page is held to a stated bar, stated once

The bar the page is held to SHALL be stated in exactly one place in the source and rendered into the skill from it.

It SHALL require each of the following, because each is a way the page fails a person who is deciding something:

- Get to the point. Length is not thoroughness, and a person who has to search the page for their decision has been handed the terminal output again.
- Speak plain English about what the finding means for the change, rather than restating the reviewer's technical prose in a nicer typeface.
- Say what each option gets the person, in the words of the thing being built rather than the words of the code. A sentence they could not act on without opening the code is not plain English, and a symbol, a constant or a test named as though they already know it is that sentence.
- Where a finding is purely technical, illustrate it with an example rather than leaving the reader to work out what the term costs them.
- Use a diagram only where a workflow is what makes the finding hard to follow, and not otherwise.
- Spend tokens as though they were the person's, because they are.

Two copies of a bar drift, and the direction is predictable: the copy nobody reads while working stays as it was, and the copy in front of the session gets softened.

Plain English was already required, and a page shipped with an option meaning "the window endpoints clamp with PAGE_LIMIT again, and the smoke assertions on the 25 cap come out".
It passed because a session that has just read the findings believes the reviewer's vocabulary is plain.
The bar therefore states a test the writer can fail rather than a standard they grade themselves against.

#### Scenario: The bar is stated once and rendered from there

- **WHEN** the bar the skill states is compared with its source
- **THEN** they are the same text, and the source is the only place it appears

#### Scenario: The bar names each way the page fails a person

- **WHEN** the bar is rendered
- **THEN** it requires getting to the point, plain English about impact, an option said in the words of the thing being built rather than the code, an example where a finding is purely technical, a diagram only where a workflow needs one, and restraint about tokens

### Requirement: The page foregrounds what the person has to act on

The page SHALL give the findings a person answers in full, and the rest of the round as a short list.

A parked round holding one question about intent and five mechanical defects is one decision and five things the session is already fixing.
A page that treats all six evenly buries the only thing the person has to act on, and a page that drops the five leaves them unable to see what else was found.

#### Scenario: The decisions are in full and the rest is a list

- **WHEN** the skill describes a page for a round holding both findings a person answers and findings the author patches
- **THEN** it says the first are given in full and the second are listed briefly

### Requirement: The page is read only, and the terminal is where the person acts

Every control the page carries SHALL carry the command it becomes, the surface SHALL hand that command back to the session, and the page SHALL NOT show it to the person.
Nothing the person does on the page SHALL be recorded against a run.
The session running `socrates respond` SHALL remain the only thing that records an answer.

The page previously carried no control at all, because there was nothing behind it and an answer given on it would have reached no run.
That reasoning is kept and its conclusion moves: a local surface now hands the answer back, so a control is permitted exactly where the session can say what command it becomes.
A control whose command the session cannot name is refused by the same reasoning that refused all of them before, because the session would be offering a decision it cannot honour.
The command is not printed on the page, because the person deciding is not the one who runs it: a command under an option is the same choice again, in a language they did not come here to read, and one they may believe they have to type.

The page SHALL store nothing that outlives the exchange, and SHALL declare no capability that reaches a run.

The page and the terminal SHALL be matchable in both directions without translation.

- Every finding on the page SHALL be shown under a name derived from the identifier socrates printed for it, and SHALL carry the same severity and the same audience `socrates findings` gave it. The surface SHALL take the name, the severity and the audience from the source the session cites rather than from what it wrote on the decision, so the match is mechanical rather than a matter of care.
- The name SHALL say which finding it is and SHALL NOT be a title of the session's invention: `f3` reaches the person as "Finding 3".
- Each decision the page presents SHALL say what choosing each option means for the change, and SHALL NOT say what the session will type.
- The page SHALL name no command and no flag that no binary offers.

An identifier is a key, and a person reading a page is not looking anything up.
A person who has read a page about "the endpoint that the migration reads" and must then work out which of six findings that was has been handed the search problem the page was supposed to remove, and a person handed `f1` as a heading has been handed the reviewer's filing system.
A name derived from the identifier is neither: the person says "finding 3", the session runs `socrates respond --finding f3`, and neither of them translates.

#### Scenario: A control carries its command and shows it to nobody

- **WHEN** the skill describes the page
- **THEN** it says every control carries the command it becomes, that the page shows it to nobody, and that nothing is recorded until the session runs it

#### Scenario: A decision made on the page is recorded in the terminal

- **WHEN** a person answers a decision on the page and sends
- **THEN** nothing about that answer is in the run until the session runs `socrates respond`, and the run then holds it

#### Scenario: The person's own words reach the record

- **WHEN** a person answers with an option or writes their own words and the session records it
- **THEN** the reason recorded is what the person chose or wrote, and not a paraphrase composed by the session

#### Scenario: Every finding keeps the identifier socrates gave it

- **WHEN** a page is served for a round
- **THEN** each finding on it is shown under the name its source gives it and carries the severity and audience `socrates findings --json` reported, taken from there rather than retyped

#### Scenario: The page says what choosing an option means for the change

- **WHEN** the skill describes a decision on the page
- **THEN** it says the decision is named from the finding's identifier and that each option says what choosing it means for the change rather than what to type

### Requirement: Where the page is drafted follows what a token costs at that branch

The skill SHALL say that the document behind the page is drafted inline where the run escalated and by a delegated subagent where the round parked, and SHALL say why.

Where the run escalated, nothing follows, so drafting is paid once, and a fresh subagent's own system prompt and tool schemas would make delegating it more total tokens for no benefit.

Where the round parked, the session is about to patch each finding, record one response per finding, continue, poll and read the next round.
Everything added to that context is re-sent on every one of those turns, and it is added at the point the round cap exists to protect, which is where the session begins compacting mid-review and loses the thing it was defending.
Writing a document is cheaper than writing a whole page was, and the reason to delegate is what the finding prose costs the driving context rather than what the markup costs, so it survives the change.

The delegated draft SHALL read the round for itself rather than being handed it, so that what returns to the driving session is the document and a sentence.

The driving session SHALL serve the page and SHALL wait on it.
A delegate that served it would own a surface it is about to exit and leave the driving session polling something it did not start.

#### Scenario: The escalated branch drafts inline

- **WHEN** the skill describes handing over an escalated run
- **THEN** it says to draft the document in the session, and gives the reason

#### Scenario: The parked branch delegates

- **WHEN** the skill describes handing over a parked round holding a finding a person answers
- **THEN** it says to delegate the drafting, says the delegate reads the round itself and returns the document, and gives the reason

#### Scenario: The driving session owns the surface

- **WHEN** the skill describes the delegated branch
- **THEN** the session that drives the loop is the one that serves the page and waits on it

#### Scenario: The page is not in the polling loop

- **WHEN** the loop the skill teaches is examined
- **THEN** nothing about the page appears in the step that waits on the review

### Requirement: socrates states the fact and does nothing further about the page

socrates SHALL NOT render, publish, store, name or verify a page.
It states that a round is one a person reads, and stops.

No command SHALL be added to `socrates`, no flag SHALL be added to `socrates`, and nothing about a page SHALL be written into a run's record.
Recording a link would be a field, a flag to write it and a thing to be stale, in exchange for a link the session has just printed.

socrates SHALL gain no network reach, no dependency and no build tag from this.
What is measured SHALL be `socrates` itself: the binary, and the packages it reaches.
A second binary in the same module is not socrates gaining reach, and measuring the module rather than the tool would answer a question nobody asked.
The rule is that the reviewing tool reaches no network, not that this repository holds no networking code.

#### Scenario: The surface does not grow

- **WHEN** the commands and flags the `socrates` binary offers are compared with what it offered before
- **THEN** none has been added, and no run record carries anything about a page

#### Scenario: socrates reaches nothing new

- **WHEN** what `socrates` itself reaches is examined
- **THEN** it reaches no network, and the packages behind a build tag are the ones that were behind one before

#### Scenario: A second binary in the module is not socrates reaching further

- **WHEN** another binary in the same module imports the network
- **THEN** what `socrates` reaches is unchanged, and the check measures the tool rather than the repository
