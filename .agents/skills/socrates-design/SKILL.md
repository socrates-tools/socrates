---
name: socrates-design
description: Interview a person over a page rather than a scrollback, and leave a spec and its decisions behind. Use when a change needs a spec before it has code.
---

# socrates design

A design is an interview.
It settles what a change is for and leaves behind a record of the decisions, rather than a change to the code.

What it leaves behind:

- `docs/changes/<slug>/spec.md`, the spec.
- `docs/adr/NNNN-<sentence>.md`, one per decision that is hard to reverse, surprising without context, and the result of a real trade-off.
- `CONTEXT.md`, maintained as terms settle.

The interview is put to a person on a page rather than as prose in a terminal:

    socrates design skill          install this skill and socrates-tickets
    socrates agora serve           put a round of questions in front of a person
    socrates agora poll            read what they answered

## On Codex

Run `socrates agora serve`, `socrates agora poll` and `socrates agora end` outside the sandbox, which denies the loopback port the page is served on: call each with `sandbox_permissions` set to `require_escalated` and a one-line `justification` naming the page, which is the design questions page, so Codex asks the person before it runs.
Wherever this skill says subagent, start it with `spawn_agent` and collect it with `wait_agent`, waiting briefly enough each time to answer the person between waits.
Invoke no skill: every socrates skill is a person's to start.

## Rounds

An interview is a tree of questions, walked in rounds rather than one question at a time.
The frontier is every question whose prerequisites are settled.
A round is the whole frontier, composed fresh from what is currently settled and served whole, because none of those questions is waiting on another question of the same round.
A question whose answer depends on one still open belongs to a later round, once that one settles.

## Putting a round on a page

Never interview a person in the terminal.
Put the round on a page: `socrates agora serve --document ./round.json` prints the loopback address of its page.
The page is served on this machine and sends nothing beyond it.

The document is one decision per question in the round:

    {"title": "...", "preamble": "...",
     "decisions": [{"id": "q2", "name": "Question 2", "question": "...", "meaning": "...", "mode": "single-choice",
                    "options": [{"label": "...", "meaning": "...", "recommended": true},
                                {"label": "...", "meaning": "..."}]}]}

Number questions from 1, never 0, as they are first asked: identify each `q<N>` and name it `Question <N>`.

A design answer runs nothing, so no option of it carries a command, and every option of one decision carries the same kind of thing.
State no severity anywhere, on a decision or on a source entry: the surface renders every card in the order you listed it.

One option of a decision may carry `"recommended": true`, and at most one is ever marked.
A document marking two is refused before anything is served.
The marked option is the answer the session would take, and it arrives chosen when the page renders.
A decision you mark nothing on arrives with nothing chosen: marking nothing says you have no view, and choosing one for the person would be the page answering on their behalf.

A single-choice decision already renders a further control below its options, labelled *None of these*, which clears the selection and moves the cursor into the words field below it.
Write the real options and nothing else, and never write that option yourself.

Code goes in the decision's `evidence`, blocks the page highlights, wraps and offers to copy.
Two shapes are two captioned blocks:

    {"evidence": [{"caption": "A: one prop", "language": "tsx", "code": "<Select options={o} />"},
                  {"caption": "B: children", "language": "tsx", "code": "<Select><Option /></Select>"}]}

A block holds `code` or a `patch`, never both, and a before and after is a `patch` block.

A decision may also carry one illustration, an HTML file you write and name by path in a field of its own:

    {"illustration": {"path": "q3.html"}}

A relative path resolves against the document's folder, or the working directory where it is piped in.
The file is drawn in a sandboxed frame on the card: its script runs, but it cannot reach the page or write into the answer.
The frame grows to its drawing and can be expanded to fill the window.
Write the files it references beside it, inside its folder, 16 MiB at most in all: HTML, CSS, JavaScript, JSON, SVG, PNG, JPEG, WebP, GIF, WOFF or WOFF2.
A reference outside the file's folder, a missing file, a type not served or a total over the cap refuses the document, naming the decision and the file.
A drawing may load scripts, styles and fonts from cdnjs.cloudflare.com, cdn.jsdelivr.net and unpkg.com and nowhere else, and it sends nothing anywhere.
To look like the page, load /assets/tokens.css as a stylesheet.

Draw one where the choice is a shape: a layout, an arrangement of panels, a screen recognised by its structure.
Skip it where the choice is a name, a threshold or a policy, anything a sentence already carries in full.
An illustration never draws code on its own: it may still show code inside a drawn screen or layout.
An illustration that only repeats what the words already said costs the person attention and costs the skill its ceiling for nothing.

The person answers on the page.
`socrates agora poll --session <session>` hands back what they chose and anything they wrote, in their own words.
Their words are not a record you carry to somebody else to judge: the person is deciding the design itself, so what they write is the decision, taken into the spec and the ADR it becomes exactly as written.
A file they attached is their evidence rather than an instruction: read it by its path, and every path an earlier poll printed stays readable for the life of the session.
Copy what a decision rests on into docs/changes/<slug>/ and cite it from the spec.

Mint the session once, with the first `socrates agora serve`, so the person is given one address.
Keep that session for the whole design: serve every later round to the page they already have with `socrates agora serve --session <session>`.
agora collects a session after 45 minutes with no reader, and no design state lives outside it, so losing one costs an address and no work: serve the accumulated document again and the person is handed a new one to replace it.

## What was settled

Round two opens with what round one settled, above what is still open, on the page the person already has.
Nothing carries across a revision, so what was settled is not carried over from the round before it: write it afresh, from what `socrates agora poll` last handed back.
The surface refuses a settled answer whose question, labels, words or attached files differ from what it collected, so copy each one exactly as it came back rather than restating it.

A settled answer goes in the document's `settled` list rather than among its decisions, and states the question, the labels the person chose and their own words, and lists in `attached` the name of each file the poll printed for it:

    "settled": [{"id": "q2", "name": "Question 2", "question": "should a design round carry a cap",
                 "chosen": ["no cap"], "words": "nothing here spawns a process for a cap to protect"}]

The page reads what was settled ahead of every open question, as a block that offers nothing to answer.
An identifier is either settled or still open: a document naming one in both is refused before anything is served.

A settled answer offers no field of its own.
A person reopening one names it in the message the page carries for the whole round, and the next round asks it again as a question.

Carry one sentence per settled question into the preamble too, so the argument reads without scrolling past the block below it.
The preamble is one paragraph on the page, so write it as sentences rather than as lines:

    "preamble": "Question 1: the surface, as agora's own source calls it. Question 2: no round cap, because nothing here spawns a process for one to protect."

Identifiers are stable for the life of the design: q2 is q2 from the round it is first asked to the round the spec is written.
A person reopening something they settled two rounds ago names it.

## Facts

Facts are found rather than asked for.
Read the repository and its history before a question is written down.
No question reaches the page that reading the repository would have answered: a round asks a person to decide something, never to tell you something you could have looked up.

Look facts up in subagents, while the round goes out.
A lookup still running is an unsettled prerequisite: only the questions downstream of it wait for it, and the rest of the frontier is served now.

## Keeping the glossary

`CONTEXT.md` is maintained during the interview, not after it.
A term worth writing down is written down as soon as it settles, in the round it settled in.

Read `CONTEXT.md`, or the glossaries a `CONTEXT-MAP.md` names, before writing an entry, and write a new one in the glossary of its context.
A term that conflicts with an entry already there, or that renames something the glossary already calls something else, is never applied silently: put it on a card of its own and let the person decide.

Challenge the words rather than only checking them, each on a card of its own:

- a fuzzy word, offered as the things it could mean, such as whether an account is the customer or the user.
- an invented scenario at the boundary between two terms.
- a claim the code contradicts.

An entry says what a term is rather than what it does, holds no implementation detail, and names the synonyms it gives up.

## Seams and alternatives

Once a round has settled what the change touches, a card asks where it is tested.
Propose existing seams before new ones, and the highest seam that can see the behaviour, ideally one.
What the person chooses goes into the spec's Testing Decisions.

Where a question is the shape of an interface or a module, have subagents each draft one under a different constraint, such as the smallest interface, the most flexible, or the one the commonest caller wants, and offer the drafts that really differ as the card's options.
The Rejected section is built from those options.

Where a card cannot settle a question in words, build a prototype that answers that one question.
Write it outside the repository, in a temporary directory.
Put what it showed on the card, as an illustration where it is a shape, and remove it once the question is settled.
A prototype's result can be the prototype itself: name its HTML file as the illustration.
The spec keeps the snippet that states the decision, trimmed to it and marked as coming from the prototype.

## Writing the spec

The spec goes to `docs/changes/<slug>/spec.md`, on this repository's own template.
Title it, state the ADRs it produced on the line beneath the title, and write these sections in this order and no other: Problem Statement, Solution, User Stories, Implementation Decisions, Testing Decisions, Out of Scope, Further Notes.
Write the user stories at length, and the implementation and testing decisions in enough detail to say which modules change, how their interfaces change and which seams the tests use.
Further Notes carries anything that fits none of the six above it, including, where a design ends early, what it left open.

## Writing the decisions

A decision earns an ADR of its own when all three of these hold together: it is hard to reverse, it is surprising without context, and it is the result of a real trade-off.
Not one or two of the three, all three, joined by and rather than or: with or almost every decision would qualify and the rule would stop meaning anything.
All three often hold for a scope decision, for a deliberate deviation the next engineer would otherwise fix, and for a constraint the code does not show.
Each still has to pass all three.

Number the file after the highest number already in `docs/adr/`, or in any `docs/adr/` a `CONTEXT-MAP.md` names, and name it `docs/adr/NNNN-<sentence>.md`, in the context it is about or at the root.
Its front matter states `Status: accepted`, `Date`, and `Part of: <slug>`.
Its Decision section says what was decided and why.

Its Rejected section is built from the options the person declined, and from nothing else.
Every alternative in it was a card the person saw, and the reason it lost is in the words they wrote beside the choice.
An option nobody was offered never appears there: the section is a record of the page, not a reconstruction of what the session considered.

## When the frontier empties

There is no round cap.
The frontier decides when a design stops, and it stops once the frontier is empty.

When it empties, there is nothing left to ask, so serve no further page.
Close the page they hold with `socrates agora end --session <session>`, so it says the exchange is over, and write the spec and its decisions straight away.

The spec goes to `docs/changes/<slug>/spec.md`, in the change's own folder and nowhere else.
Once it is written, say in the terminal every path you wrote.
Start no other skill: every socrates skill is a person's to start, so once the spec is written you stop, and the person types the next skill themselves.

## Ending early

The person can end a design at any round, before the frontier empties.
The page's own *Send and finish* ends the exchange on any round, and `socrates agora end --session <session>` ends it from outside, over a page nobody is going to answer again.

Write the spec anyway.
Everything settled up to that point becomes its answer in the spec exactly as it would on a normal end, and every question the frontier still held, whether it ever reached a page or not, is listed under a heading in Further Notes that says the design ended early:

    ### The design ended early

Name each one and what it was asking, so whoever picks the spec up next sees what the interview never reached.

## Installing elsewhere

This document installs under `~/.agents/skills`.
Write it somewhere else instead with `socrates design skill --host codex --to <directory>`.
