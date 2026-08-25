## RENAMED Requirements

- FROM: `#### Scenario: Every control names the command it becomes`
- TO: `#### Scenario: A control carries its command and shows it to nobody`

- FROM: `#### Scenario: The page says what the person does back in the terminal`
- TO: `#### Scenario: The page says what choosing an option means for the change`

## MODIFIED Requirements

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
