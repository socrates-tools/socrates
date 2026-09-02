# agent-skill Specification

## Purpose
TBD - created by archiving change the-loop-an-agent-can-drive. Update Purpose after archive.

## Requirements

### Requirement: The skill has one source and the committed copy is generated from it

The skill's text SHALL live in exactly one place in the source, and the copy committed to the repository SHALL be generated from it.
A test SHALL regenerate and compare, and SHALL fail while the two differ.

A committed file that is read rather than generated has one failure mode and it is the one that matters: somebody edits the file, the generator is never run, and the installed skill and the repository's copy differ with nothing failing.

#### Scenario: Drift fails the suite

- **WHEN** the committed skill file differs from what the source generates
- **THEN** the test fails, naming the file and the command that regenerates it

#### Scenario: Regenerating is one command

- **WHEN** the regeneration target is run
- **THEN** the committed copy matches the source and the test passes

### Requirement: The skill teaches the loop, including every terminal branch

The skill SHALL describe the sequence a session follows: implement, start a review carrying the person's stated intent, poll until the run is not working, and act on each terminal state.

It SHALL name what to do for each of `clean`, `failed`, `escalated` and `awaiting-response`, and `failed` SHALL be described before `escalated`.

A session scanning for what to do about findings reads past a state that has none, and a `failed` run has none.

It SHALL also name `cancelled`, and SHALL name every state a run can rest in.
A person who stops a review from another terminal leaves a state that is neither a verdict nor work in progress, and a loop told to wait for one of a fixed list of states polls a cancelled run for as long as anybody lets it.

The `awaiting-response` branch SHALL separate the findings the session answers from the findings the person answers.
A finding that can be patched is the session's to patch and record.
A finding that asks a person SHALL be described as the person's, to be taken to them and answered with what they decided, and the skill SHALL say that the session does not answer one on its own behalf.

Until a question ended the run, nothing needed to be said: a terminal state was what stopped a session answering it.
A parked question removes that, so the sentence is what remains, and it is the reason this requirement states it rather than leaving it to be inferred from the loop.

The skill SHALL say that the session continues once every answerable finding is answered, without asking whether the round is worth spending.
A patch is unreviewed new code written by the session that was just told it was wrong, and a round that judges it is the round the findings earned.

Where a run ends, the skill SHALL say how many rounds it took and what the last one found.
Two rounds ending with nothing blocking and one round ending the same way are different results, and a session that reports only the verdict tells the person nothing about which of the two happened.

#### Scenario: Every terminal state is covered

- **WHEN** the skill is rendered
- **THEN** it names `clean`, `failed`, `escalated` and `awaiting-response`, and says what to do for each

#### Scenario: A failed run is stated not to be a pass

- **WHEN** the skill describes `failed`
- **THEN** it says that no review completed, that nothing has been checked, and that this is not a clean result

#### Scenario: A cancelled run is not a verdict

- **WHEN** the skill describes `cancelled`
- **THEN** it says that nothing was decided, that the change is not reviewed, and that somebody stopped it on purpose

#### Scenario: Every state a run rests in is named

- **WHEN** the states the skill names are compared with the states a run can rest in
- **THEN** every one of them is named, so the loop has no state it would poll through

#### Scenario: The loop closes

- **WHEN** the skill describes `awaiting-response`
- **THEN** it says to respond to every blocking finding, then continue, then return to polling

#### Scenario: A round is spent without asking

- **WHEN** the skill describes a parked round whose findings can all be patched
- **THEN** it says to answer them and continue, and it does not ask the person whether to spend the round

#### Scenario: What the person is told at the end

- **WHEN** the skill describes reporting a finished run
- **THEN** it says to report how many rounds ran and what the last one found

#### Scenario: A question is the person's

- **WHEN** the skill describes a parked round holding a finding that asks a person
- **THEN** it says to present that finding to the person, to record what they decided, and that the session does not answer it itself

### Requirement: The skill teaches polling rather than following

The skill SHALL name `scrat status` as the command to wait with, and SHALL say that everything the session reads lands in its own context.

It SHALL name `scrat watch` only as something to suggest to the person, in another terminal, and SHALL NOT put it in the session's own loop.

It SHALL name how often to poll, and SHALL say that the wait is bounded by socrates rather than by the session.

A session given no cadence and no bound invents both, and each way of inventing them is wrong in its own direction: polling several times a second spends the review re-reading a document nobody has written yet, and a timeout of the session's own reports a live review as stuck.
Neither is needed, because a reviewer that answers nothing before its deadline ends the run as `failed`, so every run reaches a state to act on.

#### Scenario: The waiting command is status

- **WHEN** the skill describes waiting for a review
- **THEN** it names `scrat status` and gives the reason

#### Scenario: The waiting step says how often and for how long

- **WHEN** the skill describes waiting for a review
- **THEN** it names a polling interval and says that the run ends on its own, so the session invents neither

#### Scenario: Watch is for the person

- **WHEN** the skill mentions `scrat watch`
- **THEN** it is as a suggestion for the person in another terminal, and it appears in no loop the session is told to run

### Requirement: The response rules come from the same source the tool prints

The four response rules the skill states SHALL be read from the same source `scrat findings` prints them from, and a test SHALL compare the two renderings.

Two copies of four sentences drift, and the direction is predictable: the copy nobody reads while implementing stays as it was, and the copy in front of the session gets softened.

#### Scenario: One source, two surfaces

- **WHEN** the rules the skill states are compared with the rules `scrat findings` prints
- **THEN** they are the same text

#### Scenario: The rules are stated before the findings are read

- **WHEN** the skill is rendered
- **THEN** the response rules appear before the step that reads the findings

### Requirement: The skill states what socrates will not do

The skill SHALL state each of the following, because each prevents a specific attempt a session makes:

- socrates never patches anything. The session patches, in its own working tree.
- There is no way to explain the change to the reviewer. The reviewer receives the diff, the repository and the stated intent, and that is the design rather than a limitation to work around.
- A rejection is a claim that goes to the next reviewer, not a way to dismiss a finding, and a rejection with no reason is refused.

#### Scenario: The three prohibitions are present

- **WHEN** the skill is rendered
- **THEN** it states each of the three

#### Scenario: No channel is implied for the author's account

- **WHEN** the skill is rendered
- **THEN** it names no command, flag or file through which a session could send the reviewer its own account of the change

### Requirement: The skill drives only the surface that exists

Every command and flag the skill names SHALL exist, and a test SHALL check each against the binary that offers it.

The skill now names two binaries.
An extraction that recognised only one of them would leave every command of the other unchecked, so a skill naming an `agora` flag that does not exist would pass the check that exists to catch exactly that.
The check SHALL therefore attribute each named command to its binary and SHALL fail on a name no binary offers.

A skill naming a flag that does not exist is a session doing something plausible and wrong, and it fails at the point where the person has stopped watching.

#### Scenario: Every command named is real

- **WHEN** the commands and flags the skill names are extracted and checked against the binaries that offer them
- **THEN** every one of them exists

#### Scenario: A removed flag fails the suite

- **WHEN** a flag the skill names is removed from the binary that offered it
- **THEN** the check fails, naming the flag and the skill

#### Scenario: A flag of the second binary is checked too

- **WHEN** the skill names a command or flag of the surface binary that the surface binary does not offer
- **THEN** the check fails, and it does so for the same reason and in the same way as for the first binary

### Requirement: `scrat skill` installs it and says where

`scrat skill` SHALL write the skill at user level and print the path it wrote.
`--to` SHALL write it to a named directory instead, which is what the repository's own generation target uses.

Running it again with the same target SHALL succeed and produce the same result, so that whether it has been run is not a question anybody has to answer.

#### Scenario: Installing names the path

- **WHEN** `scrat skill` is run
- **THEN** the skill exists at user level and the command printed the path

#### Scenario: Installing twice

- **WHEN** `scrat skill` is run twice
- **THEN** both succeed, the second says the same thing, and the file is the same

#### Scenario: Installing elsewhere

- **WHEN** `scrat skill --to` names a directory
- **THEN** the skill is written there and nothing is written at user level

#### Scenario: A directory that cannot be written

- **WHEN** the target directory cannot be written
- **THEN** the command refuses naming the path and the reason, and exits with the usage code

### Requirement: A further look is a new run, and only where the person asks

The skill SHALL say that a run has at most two rounds, that a further look is a fresh `scrat review` rather than a third round, and that one is started only where the person asks for it.

It SHALL say that the reviewer of a new run is told nothing about the run before it, so the person knows what they are getting: a cold reading of the code as it now stands, and not a continuation of the argument.

A session that believes a third round exists will look for the flag that gives it one, and there is none.
A session that starts a fresh review on its own initiative spends a full review on a change nobody asked it to re-read, which is the cost this whole design is careful with.

#### Scenario: There is no third round

- **WHEN** the skill describes a run that escalated at the round cap
- **THEN** it says that no further round is possible and names a fresh review as the only way to look again

#### Scenario: A fresh review is the person's to ask for

- **WHEN** the skill describes starting another review after a run ended
- **THEN** it says to offer it and to wait for the person to ask, rather than starting one unprompted

#### Scenario: A new run knows nothing of the old one

- **WHEN** the skill describes a fresh review of patched code
- **THEN** it says the reviewer is told nothing about the earlier run, its findings or its responses

### Requirement: The skill is named for the product

The skill SHALL be named `socrates`, SHALL be installed at `.claude/skills/socrates/`, and SHALL name `socrates` in every invocation its steps give.

Exactly one constant in the source SHALL carry that name, and it SHALL be the one the installer writes the directory from and the one the surface check reads invocations with.
Two constants would let the skill teach a command name that the check never looks for, and the check that proves every named command exists would go on passing while naming nothing.

#### Scenario: The skill is installed under the product's name

- **WHEN** `socrates skill` is run
- **THEN** the skill is written at `.claude/skills/socrates/SKILL.md`

#### Scenario: One constant names the skill and reads its invocations

- **WHEN** the constant that names the skill is changed
- **THEN** the installer's directory and the invocations the surface check reads both follow it, and no second copy of the name remains

#### Scenario: The skill names the full command

- **WHEN** the skill's steps are read
- **THEN** every invocation in them names `socrates`, and none names the short name

### Requirement: The skill writes an invocation as code

The skill SHALL write every invocation it names as code, inline or as an indented block, and the check that reads its invocations SHALL read code alone.

The command and the product are now one word, so a sentence whose subject is the product is indistinguishable from an invocation by any rule that reads prose.
Read as prose, `socrates patches nothing` names a command called `patches`, and the check that exists to prove every named command is real would report five commands the binary has never had.

#### Scenario: Prose naming the product is not an invocation

- **WHEN** the skill states a sentence whose subject is the product rather than a command
- **THEN** the check reads no command from it, and every invocation written as code is still read

### Requirement: The skill says how a recommendation reaches the page

The skill SHALL say that a decision may mark one option as the one the session would take, and SHALL show the field where it shows the document's shape.

It SHALL say that at most one option per decision carries the mark, and that a document marking two is refused before anything is served.

It SHALL say what the mark does and what it does not do: the page reads the marked option first and says in words that it is the one the session would take, and it selects nothing.
The person's answer is still the person's, and a session that has marked an option has recommended rather than decided.

It SHALL NOT restate the bar's judgement about when to recommend.
The bar owns that, each line is stated once, and a second statement beside the mechanism is the copy that gets softened.

#### Scenario: The skill names the mark where it names the document

- **WHEN** the skill shows the shape of a decision
- **THEN** the option the session would take is marked there, in the same shape a session copies

#### Scenario: The skill says one option carries it

- **WHEN** the skill describes the mark
- **THEN** it says at most one option per decision carries it and that a document marking two is refused before serving

#### Scenario: The skill says marking is not answering

- **WHEN** the skill describes the mark
- **THEN** it says the page reads the marked option first and states it in words, that nothing arrives selected, and that the decision remains the person's

#### Scenario: When to recommend is stated once

- **WHEN** the skill is rendered
- **THEN** the judgement about when to recommend appears in the bar and nowhere else

### Requirement: The skill says how an illustration reaches the page

The skill SHALL say that a decision may carry one illustration, and that it travels beside the evidence rather than inside it.

It SHALL say an illustration takes one of two forms, and which each is for.
A drawing the session writes as inline SVG is for a workflow, because the surface embeds no diagram engine and renders what the caller supplied.
Markup the document's stylesheet paints is for a panel a person would recognise, because a person who already knows the screen is being asked to recognise it rather than to read a description of it.

It SHALL say the document carries one stylesheet for all of its illustrations rather than one per decision, so a session writing a round about one product states its look once.

It SHALL say that both the markup and the stylesheet are reduced before they are embedded, and name what does not survive: script, event handlers, remote references and form controls from the markup; imports, font faces, properties outside the allowlist and values naming a resource from the stylesheet.
It SHALL say the surface scopes the stylesheet itself, so a session does not scope its own selectors and does not expect a rule aimed at the page to reach it.

It SHALL say this beside the bar's lines about showing and about diagrams, so that a session reads when and how in one place.

It SHALL NOT restate the bar's judgement about when either is warranted.
The bar owns that, each line is stated once, and a second statement beside the mechanism is the copy that gets softened.

#### Scenario: The skill names the illustration as something a decision may carry

- **WHEN** the skill is rendered
- **THEN** it says a decision may carry one illustration, beside its evidence rather than within it

#### Scenario: The skill says what forms an illustration takes

- **WHEN** the skill describes the illustration
- **THEN** it names both a drawing the session writes and markup a stylesheet paints, says which is for what, and does not name or imply a diagram engine the surface would run

#### Scenario: The skill says the look is stated once per document

- **WHEN** the skill describes the stylesheet
- **THEN** it says one stylesheet paints every illustration in the document, rather than one per decision

#### Scenario: The skill says what is stripped before either is embedded

- **WHEN** the skill describes the illustration
- **THEN** it says both the markup and the stylesheet are reduced to an allowlist, and names what does not survive each, so that a session does not write what will be dropped

#### Scenario: The skill says the surface does the scoping

- **WHEN** the skill describes the stylesheet
- **THEN** it says the surface scopes every selector to the illustration, so a session neither scopes its own nor expects one to reach the page

#### Scenario: How to illustrate is stated beside when to

- **WHEN** the skill is rendered
- **THEN** the mechanism appears beside the bar's lines about showing and about diagrams, and the bar's judgement about when each is warranted appears once and only there

### Requirement: The skill hands every round over as a page

The skill SHALL say that on every branch where a round completed a review, the session drafts a document, serves it with `agora`, hands the person the address, and says briefly what it holds.
Those branches are `clean`, `awaiting-response` and `escalated`.

It SHALL say that the session waits on `agora poll` only where socrates reports that the round waits on a person, and that on the other branches it hands over the address and carries on with patching or reporting.
A session told to wait on every page would stop on a clean run for an answer nobody was asked for.

It SHALL say that `failed` and `cancelled` get no page: a failed run completed no review, a cancelled one was stopped on purpose, and in both cases the reason is one line socrates already prints.

It SHALL say that the address `agora serve` prints holds the session that page belongs to, and that `agora poll` is given that session.
It SHALL say it once, where the loop is described, rather than on each branch that serves a page.
A session that is not told this polls a surface holding more than one page and is told about somebody else's round, or about none, and it has no way to find out which.

It SHALL say that the session is kept for the whole review, and that a later round is written to the page the person already has rather than served as a new one.
A review is two rounds about one change and one person reading them, so a second address is a second tab, and the first goes stale with nothing on it saying so.
Every round of a run that is still going reaches the page, so the address is handed over once, at the first round that reached a result.

It SHALL take the fact from what socrates reports rather than restating the condition, so that a session never decides for itself whether a round waits on a person.

It SHALL hold the page to the stated bar rather than describing the page's qualities again, and the bar SHALL appear before the branches that refer to it.

It SHALL say that the person answers on the page, that each control names the command it becomes, and that the session is what records the answer.
It SHALL NOT tell the session to ask the person to restate in the terminal a decision they already made on the page.

Nothing about the page SHALL appear in the waiting step.
Everything in that step is run in a loop, and a page is made once, after the run has stopped moving.

#### Scenario: Every branch that has a result says to serve the page

- **WHEN** the skill is rendered
- **THEN** the `clean`, `awaiting-response` and `escalated` branches each say to serve the page and hand over the address, and the `failed` and `cancelled` branches say nothing about a page

#### Scenario: The skill says which branches wait

- **WHEN** the skill describes the branches that serve a page
- **THEN** it says to wait on `agora poll` where socrates reports the round waits on a person, and to carry on without waiting where it does not

#### Scenario: A round the author answers alone gets a page nobody waits on

- **WHEN** the skill describes a parked round whose findings the author answers alone
- **THEN** it says to serve the page and hand over the address, and then to patch, respond and continue without waiting on it

#### Scenario: The skill says to carry the session from serve to poll

- **WHEN** the skill describes serving a page and waiting on it
- **THEN** it says the printed address holds the session, and that `agora poll` is given that session

#### Scenario: The session is stated once rather than on every branch

- **WHEN** the skill is rendered
- **THEN** it says how the session reaches `agora poll` once, and no branch restates it

#### Scenario: A later round reaches the page the person already has

- **WHEN** the skill describes a second round
- **THEN** it says to write it to the session the first round was served under, and does not say to hand over a second address

#### Scenario: The session does not decide the condition itself

- **WHEN** the skill describes when to wait on a page
- **THEN** it names what socrates reports, and states no condition of its own over the state or the counts

#### Scenario: The bar is stated before the branches that refer to it

- **WHEN** the skill is rendered
- **THEN** the bar appears once, before the branches that serve a page, and no branch restates it

#### Scenario: The skill says how an answer reaches the run

- **WHEN** the skill describes making the page
- **THEN** it says the person answers on the page, that each control names the command it becomes, and that the session records it with `socrates respond`

#### Scenario: The skill says the page and the terminal use the same identifiers

- **WHEN** the skill describes what goes on the page
- **THEN** it says every finding keeps the identifier, severity and audience `socrates findings` reported, so the person can name one and the session knows which

#### Scenario: A decision is still the person's to make

- **WHEN** the skill describes a finding that needs a person
- **THEN** it says the session does not answer one on its own behalf, whichever surface the answer arrives through
