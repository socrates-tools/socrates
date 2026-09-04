# findings-page Specification

## Purpose
TBD - created by archiving change a-round-a-person-reads-becomes-a-page. Update Purpose after archive.

## Requirements

### Requirement: The page is held to a stated bar, stated once

The bar the page is held to SHALL be stated in exactly one place in the source and rendered into the skill from it.

It SHALL require each of the following, because each is a way the page fails a person who is deciding something:

- Get to the point. Length is not thoroughness, and a person who has to search the page for their decision has been handed the terminal output again.
- Speak plain English about what the finding means for the change, rather than restating the reviewer's technical prose in a nicer typeface.
- Say what each option gets the person, in the words of the thing being built rather than the words of the code. A sentence they could not act on without opening the code is not plain English. A symbol, a constant or a test named as though they already know it is that sentence.
- Say which option you would take, on every decision where you would take one, by marking that option and no other. You have read the finding and they have not. A list where every option looks equally likely makes them work out what you already know.
- Where a finding is purely technical, illustrate it with an example rather than leaving the reader to work out what the term costs them.
- Where what a finding costs is what a person sees, show them it rather than describing it.
- Use a diagram only where a workflow is what makes the finding hard to follow, and not otherwise.
- Spend tokens as though they were the person's, because they are.

Two copies of a bar drift, and the direction is predictable: the copy nobody reads while working stays as it was, and the copy in front of the session gets softened.

Plain English was already required, and a page shipped with an option meaning "the window endpoints clamp with PAGE_LIMIT again, and the smoke assertions on the 25 cap come out".
It passed because a session that has just read the findings believes the reviewer's vocabulary is plain.
The bar therefore states a test the writer can fail rather than a standard they grade themselves against.

Recommending is on the bar rather than left to the mechanism, because a field nothing calls for is a field the sessions that most need it leave empty.
It is bounded to one option for the same reason it is asked for at all: a page that ranks everything has ranked nothing.

Showing is a separate line from the example beside it and from the diagram beneath it, and the reason is that the bar's other lines are about a reader who has to be taught something.
A finding about what a person sees has a reader who already knows the thing and is being asked to recognise it.
A sketch of a panel retyped in a monospace block is the one form that reader cannot recognise it in.
The line about diagrams says "and not otherwise", so without a line of its own a page that showed a person their own screen would be a page in breach of the bar.

#### Scenario: The bar is stated once and rendered from there

- **WHEN** the bar the skill states is compared with its source
- **THEN** they are the same text, and the source is the only place it appears

#### Scenario: The bar names each way the page fails a person

- **WHEN** the bar is rendered
- **THEN** it requires getting to the point, plain English about impact, an option said in the words of the thing being built rather than the code, saying which option you would take by marking one and no other, an example where a finding is purely technical, showing what a finding costs where that is what a person sees, a diagram only where a workflow needs one, and restraint about tokens

#### Scenario: Showing and drawing are separate lines

- **WHEN** the bar is rendered
- **THEN** the line about what a person sees and the line about diagrams each appear once and neither restates the other

### Requirement: The page is read only, and the terminal is where the person acts

Every control the page carries SHALL carry the command it becomes, the surface SHALL hand that command back to the session, and the page SHALL NOT show it to the person.
Nothing the person does on the page SHALL be recorded against a run.
The session running `socrates respond` SHALL remain the only thing that records an answer.

The page previously carried no control at all, because there was nothing behind it and an answer given on it would have reached no run.
That reasoning is kept and its conclusion moves: a local surface now hands the answer back, so a control is permitted exactly where the session can say what command it becomes.
A control whose command the session cannot name is refused by the same reasoning that refused all of them before, because the session would be offering a decision it cannot honour.
The command is not printed on the page, because the person deciding is not the one who runs it.
A command under an option is the same choice again, in a language they did not come here to read, and one they may believe they have to type.

The page SHALL store nothing that outlives the exchange, and SHALL declare no capability that reaches a run.

The page and the terminal SHALL be matchable in both directions without translation.

- Every finding on the page SHALL be shown under a name derived from the identifier socrates printed for it, and SHALL carry the same severity and the same audience `socrates findings` gave it. The surface SHALL take the name, the severity and the audience from the source the session cites rather than from what it wrote on the decision, so the match is mechanical rather than a matter of care.
- The name SHALL say which finding it is and SHALL NOT be a title of the session's invention: `f3` reaches the person as "Finding 3".
- Each decision the page presents SHALL say what choosing each option means for the change, and SHALL NOT say what the session will type.
- The page SHALL name no command and no flag that no binary offers.

An identifier is a key, and a person reading a page is not looking anything up.
A person who has read a page about "the endpoint that the migration reads" and must then work out which of six findings that was has been handed the search problem the page was supposed to remove.
A person handed `f1` as a heading has been handed the reviewer's filing system.
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

The skill SHALL say that the document behind the page is drafted inline where nothing follows the round and by a delegated subagent where another round does, and SHALL say why.

Nothing follows a `clean` run or an `escalated` one, so drafting is paid once, and a fresh subagent's own system prompt and tool schemas would make delegating it more total tokens for no benefit.

Where the round parked, the session is about to patch each finding, record one response per finding, continue, poll and read the next round.
Everything added to that context is re-sent on every one of those turns.
It is added at the point the round cap exists to protect, where the session begins compacting mid-review and loses the thing it was defending.
That is true of a parked round whose findings the author answers alone as much as of one holding a question, because what costs the driving context is the finding prose rather than who answers it.

The delegated draft SHALL read the round for itself rather than being handed it, so that what returns to the driving session is the document and a sentence.

The driving session SHALL serve the page, and SHALL wait on it where the round waits on a person.
A delegate that served it would own a surface it is about to exit and leave the driving session polling something it did not start.

#### Scenario: The escalated branch drafts inline

- **WHEN** the skill describes handing over an escalated run
- **THEN** it says to draft the document in the session, and gives the reason

#### Scenario: The clean branch drafts inline

- **WHEN** the skill describes handing over a clean run
- **THEN** it says to draft the document in the session, for the same reason the escalated branch gives

#### Scenario: The parked branch delegates

- **WHEN** the skill describes handing over a parked round, whether or not it holds a finding a person answers
- **THEN** it says to delegate the drafting, says the delegate reads the round itself and returns the document, and gives the reason

#### Scenario: The driving session owns the surface

- **WHEN** the skill describes the delegated branch
- **THEN** the session that drives the loop is the one that serves the page and the one that waits on it

#### Scenario: The page is not in the polling loop

- **WHEN** the loop the skill teaches is examined
- **THEN** nothing about the page appears in the step that waits on the review

### Requirement: socrates states the wait and does nothing further about the page

socrates SHALL NOT render, publish, store, name or verify a page.
It states whether a round waits on a person, and stops.

No command SHALL be added to `socrates`, no flag SHALL be added to `socrates`, and nothing about a page SHALL be written into a run's record.
Recording a link would be a field, a flag to write it and a thing to be stale, in exchange for a link the session has just printed.

socrates SHALL NOT state whether a page is served, because every round that completed a review is handed over as one and a fact true of all of them tells the session nothing it did not have.

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

### Requirement: Every round that completed a review is handed over as a page

Where a run holds a round that completed a review, the session driving the loop SHALL hand that round over as a page served locally by `agora`, give the person the address, and say briefly what it holds rather than pasting the findings.

That is every round reached in `clean`, in `awaiting-response` and in `escalated`, whatever the round holds and whoever answers it.

A run that reached `failed` or `cancelled` SHALL NOT be handed over as a page, and the reason is not the same reason for the two.

A `failed` run completed no review, so there is nothing on such a page but the one-line reason socrates already prints.
A tab to open for a sentence already in front of the person is worse than the sentence.

A `cancelled` run may hold a round that did complete a review, and it is still handed over as nothing.
Somebody stopped the review on purpose, so nothing was decided about the change, and handing them a page of what the run got as far as would be answering a question they withdrew.
The rule is therefore stated over the state the run ended in rather than over whether any round of it completed a review, and `socrates findings` reports a cancelled run as unreviewed for the same reason.

The page was once only for a round that could not move until a person answered, because one condition decided both whether a page was served and whether the session waited on it.
The waiting needed deciding and the serving did not.
The page exists because a review result reads better on a page than in a scrollback.
That is as true of the two commonest outcomes, a clean run and a round the author patches alone, as of the ones that stop for the person.

The page SHALL NOT be published to any host outside the machine.
A findings document quotes the code it was written about, `internal/corpus` refuses to record one for that reason, and a skill that told a session to upload one contradicted that refusal.
Serving it locally is what makes the two rules agree.

The reviewer SHALL NOT draft it.
It is sealed read-only inside the snapshot worktree, it runs with the machine's own customisations disabled, and its findings are the untrusted thing under judgement.
A reviewer that also composed the person's summary of its own findings would be arguing its case in the document the person decides from.

socrates SHALL NOT draft it either.
What is being asked for is plain English about what a finding means for the change, which a template over the findings document cannot produce.
socrates' agent seam exists to review rather than to write for a person.

The page SHALL NOT replace what the session reads.
`socrates findings` keeps printing the findings, because the session patches from them, and the page is a reading aid rather than the record.

#### Scenario: A round that completed a review is handed over as a page

- **WHEN** the loop reaches a run in `clean`, in `awaiting-response` or in `escalated`
- **THEN** the session serves it as a page, hands over the address, and does not paste the findings into the terminal

#### Scenario: A run that ended without a result is handed over as nothing

- **WHEN** the loop reaches a run in `failed`, or a run in `cancelled` whose first round completed a review
- **THEN** the session says what happened and why, and serves no page

#### Scenario: A page is served locally and published nowhere

- **WHEN** the skill describes handing a round over
- **THEN** it names the local surface, and names no host, upload, publish or artifact outside the machine

#### Scenario: Nothing socrates spawns drafts a page

- **WHEN** the prompt a reviewer is given and the roles socrates spawns are examined
- **THEN** none of them is asked for a summary written for a person, and nothing about a page reaches any prompt socrates builds

### Requirement: Serving the page and waiting on it are separate

The session SHALL serve a page for every round it hands over, and SHALL wait on that page only where socrates reports that the round waits on a person.

Where the round waits on nobody, the session SHALL hand over the address, say briefly what the page holds, and carry on: patching, recording and continuing where the round parked, and reporting where the run is over.
It SHALL NOT wait, and it SHALL NOT ask the person to say they have read it.

These were one condition, and separating them is the whole of this change.
A session that waited on a clean run would stop for an answer nobody was asked for, and a person who never opened that page would hang the loop for good.
A session that did not wait on a round whose findings are the person's would patch around them and record decisions nobody made.
So the fact socrates reports decides the wait, and the page is served either way.

A page nobody waits on SHALL outlive the session's attention.
The surface holds it until it has been idle, so a person who reads a clean review an hour later is reading a page that is still there and that nothing is blocked on.

#### Scenario: A round that waits on a person is waited on

- **WHEN** the loop reaches a round socrates reports as one that waits on a person
- **THEN** the session serves the page, hands over the address, and waits for what the person sends back before recording anything

#### Scenario: A round that waits on nobody is handed over and not waited on

- **WHEN** the loop reaches a clean run, or a parked round whose findings the author answers alone
- **THEN** the session serves the page, hands over the address, and carries on without waiting on it

#### Scenario: A page nobody waits on does not stop the loop

- **WHEN** a page is served for a round that waits on nobody and the person never opens it
- **THEN** the session patches, records and continues, or reports and stops, exactly as it would have with no page at all

### Requirement: The page carries the whole round, one card per finding

The page SHALL carry every finding in the round, one card for each.
No finding SHALL be reduced to a line in a list or a sentence in the preamble.

The source list the surface renders SHALL remain what it is, an index of everything found, and SHALL NOT be what stands in for a finding the person never sees in full.

The page SHALL foreground what the person has to act on by the order the cards arrive in and by the absence of controls on the findings they do not.

A parked round holding one question about intent and five mechanical defects is one decision and five things the session is already fixing.
A page that treats all six evenly buries the only thing the person has to act on, and a page that drops the five leaves them unable to see what else was found.
What separates them is that the severest card is the first the person meets, and that a finding the session resolved carries no option for them to weigh.

The skill SHALL say that a finding the session resolved is put on the page as what the finding was and what resolved it, and that the diff goes on it where the diff is what makes the resolution clear.

The skill SHALL say that such a finding carries no options, and SHALL say why: an option on a settled finding is a control that changes nothing, and a person reads it as a question still owed an answer.

Where the round holds nothing a person answers, the page SHALL still carry every finding as a card, and SHALL present none of them as a decision.
A clean run and a round the author patched alone are rounds with nothing to act on.
A card that asks nothing is not a decision. It is the account of what was found and what was done, which is the whole of what those rounds have to hand over.

#### Scenario: Every finding in the round is a card

- **WHEN** the skill describes a page for a round holding both findings a person answers and findings the author resolved
- **THEN** it says every finding in the round gets a card of its own, and none is reduced to a line in a list or a sentence in the preamble

#### Scenario: A resolved finding is stated rather than asked

- **WHEN** the skill describes a card for a finding the session already resolved
- **THEN** it says the card states what the finding was and what resolved it, carries the diff where that is what makes it clear, and offers no options

#### Scenario: The reason for offering no options is given

- **WHEN** the skill says a resolved finding carries no options
- **THEN** it says that a control which changes nothing reads as a question the person still has to answer

#### Scenario: A round with nothing to answer is still a round in full

- **WHEN** the skill describes a page for a clean run, or for a parked round whose findings the author answers alone
- **THEN** it says every finding is still a card, that none carries an option, and that nothing is presented as a decision

#### Scenario: Foregrounding survives without demoting anything

- **WHEN** the skill describes how the page separates what needs the person from what does not
- **THEN** it names the severity order and the absence of controls, and does not tell the session to shorten, list or omit any finding
