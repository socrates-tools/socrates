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

The interview is put to a person on a page, the same way socrates review hands a finding to one, rather than as prose in a terminal:

    socrates design skill          install this skill
    socrates agora serve           put a round of questions in front of a person
    socrates agora poll            read what they answered

## Rounds

An interview is a tree of questions, walked in rounds rather than one question at a time.
The frontier is every question whose prerequisites are settled.
A round is the whole frontier, composed fresh from what is currently settled and served whole, because none of those questions is waiting on another question of the same round.
A question whose answer depends on one still open belongs to a later round, once that one settles.

## Putting a round on a page

Never interview a person in the terminal.
Put the round on a page: `socrates agora serve --document ./round.json` prints the loopback address of its page.
The page is served on this machine and reaches nowhere beyond it.

The document is one decision per question in the round:

    {"title": "...", "preamble": "...",
     "decisions": [{"id": "q2", "question": "...", "meaning": "...", "mode": "single-choice",
                    "options": [{"label": "...", "meaning": "...", "recommended": true},
                                {"label": "...", "meaning": "..."}]}]}

A design answer runs nothing, so no option of it carries a command, and every option of one decision carries the same kind of thing.
State no severity anywhere, on a decision or on a source entry: a design question is not more or less severe than another one, it is earlier or later in the argument, and the surface renders every card in the order you listed it.

One option of a decision may carry `"recommended": true`, and at most one is ever marked.
A document marking two is refused before anything is served.
The marked option is the answer the session would take, and it arrives chosen when the page renders, so a person who agrees with the whole round need only press once.
A decision you mark nothing on arrives with nothing chosen: marking nothing says you have no view, and choosing one for the person would be the page answering on their behalf.

A single-choice decision already renders a further control below its options, labelled *None of these*, which clears the selection and moves the cursor into the words field below it.
Write the real options and nothing else, and never write that option yourself.

The person answers on the page.
`socrates agora poll --session <session>` hands back what they chose and anything they wrote, in their own words.
Their words are not a record you carry to somebody else to judge, the way a reviewer's finding is: here the person is deciding the design itself, so what they write is the decision, taken into the spec and the ADR it becomes exactly as they wrote it.

Mint the session once, with the first `socrates agora serve`, so the person is given one address.
Keep that session for the whole design: serve every later round to the page they already have with `socrates agora serve --session <session>`, so a further round costs them no second address.
agora collects a session after 45 minutes with no reader, and because no design state lives outside the session, losing one costs an address and no work: serve the accumulated document again and the person is handed a new one to replace it.

## Facts

Facts are found rather than asked for.
Read the repository, its history and what it already says of itself before a question is written down: a name already in use, a pattern the codebase already follows, a decision already on record.
No question reaches the page that reading the repository would have answered.
A round asks a person to decide something, never to tell you something you could have looked up yourself.

## Installing elsewhere

This document installs under `~/.claude/skills` by default.
Write it somewhere else instead with `socrates design skill --to <directory>`.
