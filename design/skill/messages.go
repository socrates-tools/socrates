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

    socrates design skill          install this skill and socrates-tickets
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
Nothing carries across a revision, so what was settled is not carried over from the round before it: write it afresh, from what |socrates agora poll| last handed back.

A settled answer goes in the document's |settled| list rather than among its decisions, and states the question, the labels the person chose and their own words:

    "settled": [{"id": "q2", "question": "should a design round carry a cap",
                 "chosen": ["no cap"], "words": "nothing here spawns a process for a cap to protect"}]

The page reads what was settled ahead of every open question, as a block that offers nothing to answer.
An identifier is either settled or still open: a document naming one in both is refused before anything is served.

A settled answer offers no field of its own.
A person reopening one names it in the message the page carries for the whole round, and the next round asks it again as a question.

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
                                {"label": "write it"},
                                {"label": "write it, and cut its tickets next"},
                                {"label": "write it, open it as an issue, and cut its tickets next"}]}]}

Offer publishing on this same card where the repository states an issue tracker, such as |docs/agents/issue-tracker.md|, and offer cutting its tickets next beside it, so every combination of the two is an option of its own, as above.
State nothing about publishing where it does not.
Where it states none, the card offers two options: write it, and write it and cut its tickets next.
Cutting the tickets next means that once the spec is written, you read |socrates-tickets| and cut it, in this same session and to the page the person already has.
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

const (
	ticketsFrontMatter = `---
name: ` + TicketsName + `
description: Cut a spec into tickets an implementer works through, written as files beside the spec, with one main ticket that holds the order and the status. Use when a spec under docs/changes/<slug>/ is ready to be built.
---`

	ticketsPurpose = `# socrates tickets

A spec says what a change is for and what was decided about it.
Tickets say what to build first, what can be built beside something else, and how far the building has got.

What it leaves behind, beside the spec:

    docs/changes/<slug>/tickets/
      README.md             the main ticket
      1-<task>.md
      2.1-<task>.md
      2.2-<task>.md
      3-<task>.md

The breakdown is put to a person on a page, the way a design's questions are, rather than as a list in a terminal:

    socrates design skill          install this skill beside socrates-design
    socrates agora serve           put the proposed breakdown in front of a person
    socrates agora poll            read what they answered`

	startingHeading = "## The spec it starts from"

	theStartingSection = startingHeading + `

Start from a spec at |docs/changes/<slug>/spec.md|, and nothing else.
Read it whole, the ADRs its Decisions line names, and |CONTEXT.md|, whose words every ticket uses.
A plan, an issue or a conversation is not a spec: tell the person to run a design first, and cut nothing.
The slug comes from the spec's own path, and the tickets go to |docs/changes/<slug>/tickets/|.

Where |docs/changes/<slug>/tickets/| already exists, write nothing and name the folder.
What it holds may already record work in progress or done, and a person who wants a fresh cut removes it themselves, knowing what its main ticket recorded.`

	slicingHeading = "## Slicing"

	theSlicingSection = slicingHeading + `

Read the code the spec touches before slicing it.
Cut the change into tracer-bullet tickets: each one a vertical slice, a narrow path through every layer the behaviour needs, demoable or verifiable on its own, and sized to fit one fresh context.
Never a horizontal slice of one layer.

Prefactoring comes first: where the code can be made easier to change, that is a ticket of its own, and every slice it helps names it as a blocker.

A wide refactor is the exception to slicing vertically: one mechanical change whose reach breaks every call site at once, so no slice of it lands green.
Sequence it as expand and contract instead: expand by adding the new form beside the old, migrate the call sites in batches each blocked by the expand, and contract by removing the old form in a ticket blocked by every batch.

A ticket carries no file path and no code snippet, because both go stale before the ticket is built.
The exception is a snippet a prototype produced that states a decision more precisely than prose can, such as a state machine or the shape of a type, trimmed to the decision and marked as coming from the prototype.`

	wavesHeading = "## Blockers and waves"

	theWavesSection = wavesHeading + `

Every ticket names the tickets that block it: the ones that must be done before it can start, and no others.
That list is the whole dependency graph.

A ticket's wave is derived from its blockers rather than chosen: one more than the highest wave among its blockers, and 1 for a ticket with none.
A wave holding one ticket names it |N-<task>.md|, and a wave holding several names them |N.1-<task>.md|, |N.2-<task>.md| and on, in the order the breakdown listed them.
|<task>| is a few words in the glossary's vocabulary, lower case and joined by hyphens.

The name orders the folder and gates nothing.
A ticket starts once its own blockers are done, so |3.1| blocked only by |2.1| starts while |2.2| is still being built.

A name is fixed once the tickets are written.
Abandoning a ticket later passes its blockers to the tickets it blocked and renames nothing, so no link to a ticket breaks.`

	breakdownHeading = "## Putting the breakdown on a page"

	theBreakdownSection = breakdownHeading + `

Never put a breakdown to a person in the terminal.
Put it on a page: |socrates agora serve --document ./breakdown.json| prints the loopback address of its page, and |socrates agora poll --session <session>| hands back what the person answered.

The document is one card per proposed ticket, in the order of their waves, each identified by the file name it would be written to:

    {"title": "...", "preamble": "...",
     "decisions": [{"id": "2.1-parse", "question": "2.1: ...", "meaning": "delivers: ... blocked by: 1-prefactor",
                    "mode": "single-choice",
                    "options": [{"label": "keep it", "meaning": "cut as proposed", "recommended": true}]}]}

Each card states the ticket's title, what it delivers from the user's side and what blocks it, and offers *keep it*, marked recommended so it arrives chosen: a person who agrees with the whole breakdown presses once.
A merge, a split or a different blocker is asked for in the words field of the card concerned, after taking *None of these*, which the page already renders, so never write that option yourself.
A breakdown answer runs nothing, so no option carries a command, and no card states a severity.

Where a design handed on to this skill, serve the breakdown to the design's session with |socrates agora serve --session <session>|, so the person keeps the page they have, and mint a new one only where there was none, or where that exchange ended or the session was collected.
Keep that session: serve every revision to the page the person already has.
Revise the breakdown from what they wrote and serve it again, until a round comes back with every card kept and nothing written beside it.`

	lastCardHeading = "## The last card"

	theLastCardSection = lastCardHeading + `

Once every card comes back kept, the last page is one card, and it is the only card on that page: every file the folder will hold, the main ticket among them, and the folder they go to.
A person approves a location rather than discovers one:

    {"title": "one card before the tickets are written",
     "decisions": [{"id": "write", "question": "write the tickets now",
                    "meaning": "path: docs/changes/<slug>/tickets/. files: README.md, 1-prefactor.md, 2.1-parse.md, 2.2-render.md, 3-cli.md",
                    "mode": "single-choice",
                    "options": [{"label": "write them", "recommended": true}]}]}

Nothing is written until that card is approved.
A person who wants something else takes *None of these* and says what in the words field.`

	mainTicketHeading = "## The main ticket"

	theMainTicketSection = mainTicketHeading + `

|README.md| in the tickets folder is the main ticket: the first file whoever implements the change opens, and outside the numbering, so no number has to mean two things.
It points at the spec and each ADR it names, and repeats no decision and argues for none, because the spec and the ADRs already say what was decided and why.

It holds, in this order:

    # <the change's title>

    <a few sentences: what the change is for, for whoever builds it>

    Spec: ../spec.md. Decisions: one link per ADR the spec names.

    ## Working through the tickets

    - Commit this folder before handing out any ticket.
    - A ticket may start once every ticket it names as a blocker is done.
    - Mark a ticket in progress before handing it out.
    - Mark it done once its acceptance criteria hold and its work is committed.
    - Mark it abandoned, with a reason, where building it turns out to be wrong, and in the same change every ticket it blocked takes its blockers in its place.

    ## Tickets

Under Tickets goes a Markdown table, one row per ticket and three columns: the ticket, linked to its file; the tickets that block it; and its status.

The main ticket is the one place a ticket's status is kept, and the status is one of four: pending, in progress, done, or abandoned with a reason beside it.
Every ticket starts pending, and an abandoned ticket keeps its row, so the table still says what was planned.
The implementing session is the one writer of that table.`

	aTicketHeading = "## A ticket"

	theTicketSection = aTicketHeading + `

Each task ticket is a file of its own, and holds:

    # 2.1: <title>

    ## What to build

    <the behaviour this ticket makes work, from the user's side, not a list of layers>

    ## Acceptance criteria

    - [ ] <a behaviour somebody can check>

    ## Blocked by

    - 1-prefactor.md

A ticket with no blockers says None under Blocked by, and can start at once.
A ticket carries no status: its status is the main ticket's.
Its acceptance criteria are behaviour to check and nothing more: a ticket names no test and no scenario, because a test's name is code and belongs to whoever builds it, and a repository that wants its tests named for scenarios says so in its own instructions.`

	committingHeading = "## Leaving them to be committed"

	theCommittingSection = committingHeading + `

Leave the tickets uncommitted, as a design leaves its spec: committing is the person's.
Say what was written and where, and say that the folder must be committed before any ticket is handed out, because an implementer working in a worktree of its own reads only what is committed.`

	ticketsInstallingElsewhere = installingHeading + `

This document installs under |~/.claude/skills| beside |socrates-design|.
Write both somewhere else instead with |socrates design skill --to <directory>|.`
)
