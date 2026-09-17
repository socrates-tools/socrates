package skill

const (
	frontMatter = `---
name: ` + Name + `
description: Interview a person over a page rather than a scrollback, and leave a spec and its decisions behind. Use when a change needs a spec before it has code.
---`

	purpose = `# socrates design

A design is an interview.
It settles what a change is for and leaves behind a record of the decisions, rather than a change to the code.

What it leaves behind:

- |docs/changes/<slug>/spec.md|, the spec.
- |docs/adr/NNNN-<sentence>.md|, one per decision that is hard to reverse, surprising without context, and the result of a real trade-off.
- |CONTEXT.md|, maintained as terms settle.

The interview is put to a person on a page, the same way socrates review hands a finding to one, rather than as prose in a terminal:

    socrates design skill          install this skill
    socrates agora serve           put a round of questions in front of a person
    socrates agora poll            read what they answered`

	roundsHeading = "## Rounds"

	theRoundsSection = roundsHeading + `

An interview is a tree of questions, walked in rounds rather than one question at a time.
The frontier is every question whose prerequisites are settled.
A round is the whole frontier, composed fresh from what is currently settled and served whole, because none of those questions is waiting on another question of the same round.
A question whose answer depends on one still open belongs to a later round, once that one settles.`

	pageHeading = "## Putting a round on a page"

	pageTemplate = pageHeading + `

Never interview a person in the terminal.
Put the round on a page: |socrates agora serve --document ./round.json| prints the loopback address of its page.
The page is served on this machine and reaches nowhere beyond it.

The document is one decision per question in the round:

    {"title": "...", "preamble": "...",
     "decisions": [{"id": "q2", "question": "...", "meaning": "...", "mode": "single-choice",
                    "options": [{"label": "...", "meaning": "...", "recommended": true},
                                {"label": "...", "meaning": "..."}]}]}

A design answer runs nothing, so no option of it carries a command, and every option of one decision carries the same kind of thing.
State no severity anywhere, on a decision or on a source entry: a design question is not more or less severe than another one, it is earlier or later in the argument, and the surface renders every card in the order you listed it.

One option of a decision may carry |"recommended": true|, and at most one is ever marked.
A document marking two is refused before anything is served.
The marked option is the answer the session would take, and it arrives chosen when the page renders, so a person who agrees with the whole round need only press once.
A decision you mark nothing on arrives with nothing chosen: marking nothing says you have no view, and choosing one for the person would be the page answering on their behalf.

A single-choice decision already renders a further control below its options, labelled *None of these*, which clears the selection and moves the cursor into the words field below it.
Write the real options and nothing else, and never write that option yourself.

A decision may also carry one illustration, painted by the document's own stylesheet, in a field of its own:

    {"fragment": {"markup": "<figure>...</figure>"}}

State the look once for the whole document in |"stylesheet"|.
The surface reduces both the markup and the stylesheet to an allowlist and scopes every selector it keeps to the illustration.

Draw one where the choice is a shape: a layout, an arrangement of panels, a screen the person already recognises by its structure rather than by a description of it.
Skip it where the choice is a name, a threshold or a policy, anything a sentence already carries in full.
An illustration that only repeats what the words already said costs the person attention and costs the skill its ceiling for nothing.

The person answers on the page.
|socrates agora poll --session <session>| hands back what they chose and anything they wrote, in their own words.
Their words are not a record you carry to somebody else to judge, the way a reviewer's finding is: here the person is deciding the design itself, so what they write is the decision, taken into the spec and the ADR it becomes exactly as they wrote it.

Mint the session once, with the first |socrates agora serve|, so the person is given one address.
Keep that session for the whole design: serve every later round to the page they already have with |socrates agora serve --session <session>|, so a further round costs them no second address.
agora collects a session after 45 minutes with no reader, and because no design state lives outside the session, losing one costs an address and no work: serve the accumulated document again and the person is handed a new one to replace it.`

	settledHeading = "## What was settled"

	theSettledSection = settledHeading + `

Round two opens with what round one settled, above what is still open, on the page the person already has.
Nothing carries across a revision, so a settled card is not carried over from the round before it: write it afresh, from what |socrates agora poll| last handed back.

A settled card states the question, the answer and the person's own words, and offers nothing:

    {"id": "q2", "question": "should a design round carry a cap the way a review round does",
     "meaning": "settled: no cap, because nothing here spawns a process for a cap to protect"}

Such a card carries no options and states no mode either: a document giving a mode to a decision that offers nothing is refused before anything is served.
The card still carries the field a person writes in, for whoever wants to reopen it.

List every settled card ahead of the round's open questions.
That is what keeps the block first and contiguous: nothing else about a design document reorders a card, so the position you list it in is the position it renders in.

Carry one sentence per settled question into the preamble too, so the argument reads without scrolling past the block below it.
The preamble is one paragraph on the page, so write it as sentences rather than as lines:

    "preamble": "q1: the surface, matching what agora's own source already calls it. q2: no round cap, because nothing here spawns a process for one to protect."

Identifiers are stable for the life of the design: q2 is q2 from the round it is first asked to the round the spec is written.
A person reopening something they settled two rounds ago names it, and the spec and the ADR it becomes are already written against that name.`

	factsHeading = "## Facts"

	theFactsSection = factsHeading + `

Facts are found rather than asked for.
Read the repository, its history and what it already says of itself before a question is written down: a name already in use, a pattern the codebase already follows, a decision already on record.
No question reaches the page that reading the repository would have answered.
A round asks a person to decide something, never to tell you something you could have looked up yourself.`

	glossaryHeading = "## Keeping the glossary"

	theGlossarySection = glossaryHeading + `

|CONTEXT.md| is maintained during the interview, not after it.
A term worth writing down is written down as soon as it settles, in the round it settled in.

Read |CONTEXT.md| before writing an entry.
A term that conflicts with an entry already there, or that renames something the glossary already calls something else, is never applied silently: put it on a card of its own and let the person decide.
The value of the glossary is that the person agreed to the word.`

	specHeading = "## Writing the spec"

	theSpecSection = specHeading + `

The spec goes to |docs/changes/<slug>/spec.md|, on this repository's own template.
Title it, state the ADRs it produced on the line beneath the title, and write these sections in this order and no other: Problem Statement, Solution, User Stories, Implementation Decisions, Testing Decisions, Out of Scope, Further Notes.
Write the user stories at length rather than in summary, and write the implementation and testing decisions in enough detail to say which modules change, how their interfaces change and which seams the tests use.
Further Notes carries anything that fits none of the six sections above it, including, where a design ends early, what it left open.`

	decisionsHeading = "## Writing the decisions"

	theDecisionsSection = decisionsHeading + `

A decision earns an ADR of its own when all three of these hold together: it is hard to reverse, it is surprising without context, and it is the result of a real trade-off.
Not one or two of the three, all three, joined by and rather than or: with or almost every decision would qualify and the rule would stop meaning anything.
A decision that fails even one of the three stays inside the spec's Implementation Decisions and gets no file of its own.

Number the file after the highest number already in |docs/adr/|, and name it |docs/adr/NNNN-<sentence>.md|.
Its front matter states |Status: accepted|, |Date|, and |Part of: <slug>|.
Its Decision section says what was decided and why.

Its Rejected section is built from the options the person declined, and from nothing else.
Every alternative in it was a card the person saw, and the reason it lost is in the words they wrote beside the choice.
An option nobody was offered never appears there: the section is a record of the page, not a reconstruction of what the session considered.`

	lastRoundHeading = "## The last round"

	theLastRoundSection = lastRoundHeading + `

There is no round cap.
The frontier decides when a design stops, and it stops once the frontier is empty.

When it empties, the last round is one card, and it is the only card on that page: the shape of what was agreed, what will be written and the path it will be written to.
A person approves a location rather than discovers one:

    {"title": "one card before the spec is written",
     "decisions": [{"id": "write", "question": "write the spec now",
                    "meaning": "shape: Problem Statement, Solution, User Stories, Implementation Decisions, Testing Decisions, Out of Scope, Further Notes. path: docs/changes/<slug>/spec.md",
                    "mode": "single-choice",
                    "options": [{"label": "write it, and open it as an issue", "meaning": "this repository states an issue tracker", "recommended": true},
                                {"label": "write it, and leave publishing for later"}]}]}

Offer publishing on this same card where the repository states an issue tracker, such as |docs/agents/issue-tracker.md|.
State nothing about publishing where it does not, and offer only the choice to write.
A person who wants a different path takes *None of these* and writes it in the words field: the card proposes a location, and does not insist on it.`

	endingHeading = "## Ending early"

	theEndingSection = endingHeading + `

The person can end a design at any round, not only the one carrying the last card.
The page's own *Send and finish* ends the exchange on any round, and |socrates agora end --session <session>| ends it from outside, over a page nobody is going to answer again.

Write the spec anyway.
Everything settled up to that point becomes its answer in the spec exactly as it would on a normal end, and every question the frontier still held, whether it ever reached a page or not, is listed under a heading in Further Notes that says the design ended early:

    ### The design ended early

Name each one and what it was asking, so whoever picks the spec up next can see what the interview never reached rather than guessing that nothing was left.`

	installingHeading = "## Installing elsewhere"

	installingElsewhere = installingHeading + `

This document installs under |~/.claude/skills| by default.
Write it somewhere else instead with |socrates design skill --to <directory>|.`
)
