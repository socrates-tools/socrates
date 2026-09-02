## RENAMED Requirements

- FROM: `### Requirement: socrates states the fact and does nothing further about the page`
- TO: `### Requirement: socrates states the wait and does nothing further about the page`

- FROM: `#### Scenario: A round a person reads is handed over as a page`
- TO: `#### Scenario: A round that completed a review is handed over as a page`

- FROM: `#### Scenario: The page is served locally and published nowhere`
- TO: `#### Scenario: A page is served locally and published nowhere`

- FROM: `#### Scenario: Nothing socrates spawns drafts the page`
- TO: `#### Scenario: Nothing socrates spawns drafts a page`

## ADDED Requirements

### Requirement: Every round that completed a review is handed over as a page

Where a run holds a round that completed a review, the session driving the loop SHALL hand that round over as a page served locally by `agora`, give the person the address, and say briefly what it holds rather than pasting the findings.

That is every round reached in `clean`, in `awaiting-response` and in `escalated`, whatever the round holds and whoever answers it.

A run that reached `failed` or `cancelled` SHALL NOT be handed over as a page, and the reason is not the same reason for the two.

A `failed` run completed no review, so there is nothing on such a page but the one-line reason socrates already prints, and a tab to open for a sentence already in front of the person is worse than the sentence.

A `cancelled` run may hold a round that did complete a review, and it is still handed over as nothing.
Somebody stopped the review on purpose, so nothing was decided about the change, and handing them a page of what the run got as far as would be answering a question they withdrew.
The rule is therefore stated over the state the run ended in rather than over whether any round of it completed a review, and `socrates findings` reports a cancelled run as unreviewed for the same reason.

The page was once only for a round that could not move until a person answered, because one condition decided both whether a page was served and whether the session waited on it.
The waiting needed deciding and the serving did not.
The page exists because a review result reads better on a page than in a scrollback, which is as true of the two commonest outcomes, a clean run and a round the author patches alone, as of the ones that stop for the person.

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

## MODIFIED Requirements

### Requirement: The page foregrounds what the person has to act on

The page SHALL give the findings a person answers in full, and the rest of the round as a short list.

A parked round holding one question about intent and five mechanical defects is one decision and five things the session is already fixing.
A page that treats all six evenly buries the only thing the person has to act on, and a page that drops the five leaves them unable to see what else was found.

Where the round holds nothing a person answers, the page SHALL be that short list and the account of what was reviewed, and nothing more.
A clean run and a round the author patches alone are rounds with nothing to act on, so a page for one of them that presented anything as a decision would be inventing one.

The short list SHALL be taken from what socrates reported rather than retyped by the session, on the same footing as the name, severity and audience a decision carries.
A list the session retypes is a list that can disagree with the terminal, and it is the part of the page that exists precisely so the person can see what else was found.

#### Scenario: The decisions are in full and the rest is a list

- **WHEN** the skill describes a page for a round holding both findings a person answers and findings the author patches
- **THEN** it says the first are given in full and the second are listed briefly

#### Scenario: A round with nothing to answer is a page with nothing to answer

- **WHEN** a page is served for a clean run, or for a parked round whose findings the author answers alone
- **THEN** it holds the round as a short list and presents nothing as a decision

### Requirement: Where the page is drafted follows what a token costs at that branch

The skill SHALL say that the document behind the page is drafted inline where nothing follows the round and by a delegated subagent where another round does, and SHALL say why.

Nothing follows a `clean` run or an `escalated` one, so drafting is paid once, and a fresh subagent's own system prompt and tool schemas would make delegating it more total tokens for no benefit.

Where the round parked, the session is about to patch each finding, record one response per finding, continue, poll and read the next round.
Everything added to that context is re-sent on every one of those turns, and it is added at the point the round cap exists to protect, which is where the session begins compacting mid-review and loses the thing it was defending.
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

## REMOVED Requirements

### Requirement: The page is what a round a person reads is handed over as

**Reason**: It made the page conditional on the round being one that could not move until a person answered, which left a clean run and a round the author patches alone with no page at all.
Those are the two commonest outcomes of a review, and a result reads better on a page than in a scrollback whether or not anybody is being asked to answer it.

**Migration**: `Every round that completed a review is handed over as a page` replaces it.
Everything it required of a page that is served is kept word for word: served locally, published nowhere, drafted by neither the reviewer nor socrates, and not a replacement for what `socrates findings` prints.
What changes is which rounds reach one, and the condition it carried is now `Serving the page and waiting on it are separate`, where it decides the wait alone.
