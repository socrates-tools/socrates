---
name: socrates-tickets
description: Cut a spec into tickets an implementer works through, written as files beside the spec, with one main ticket that holds the order and the status. Use when a spec under docs/changes/<slug>/ is ready to be built.
disable-model-invocation: true
---

# socrates tickets

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
    socrates agora poll            read what they answered

## The spec it starts from

Start from a spec at `docs/changes/<slug>/spec.md`, and nothing else.
Read it whole, the ADRs its Decisions line names, and `CONTEXT.md`, or the glossaries a `CONTEXT-MAP.md` names, whose words every ticket uses.
A plan, an issue or a conversation is not a spec: tell the person to run a design first, and cut nothing.
The slug comes from the spec's own path, and the tickets go to `docs/changes/<slug>/tickets/`.

Where `docs/changes/<slug>/tickets/` already exists, write nothing and name the folder.
What it holds may already record work in progress or done, and a person who wants a fresh cut removes it themselves, knowing what its main ticket recorded.

## Slicing

Read the code the spec touches before slicing it.
Cut the change into tracer-bullet tickets: each one a vertical slice, a narrow path through every layer the behaviour needs, demoable or verifiable on its own, and sized to fit one fresh context.
Never a horizontal slice of one layer.
A change that fits one fresh context whole is one ticket, and one ticket is a breakdown.

Prefactoring comes first: where the code can be made easier to change, that is a ticket of its own, and every slice it helps names it as a blocker.

A wide refactor is the exception to slicing vertically: one mechanical change whose reach breaks every call site at once, so no slice of it lands green.
Sequence it as expand and contract instead: expand by adding the new form beside the old, migrate the call sites in batches each blocked by the expand, and contract by removing the old form in a ticket blocked by every batch.

A ticket carries no file path and no code snippet, because both go stale before the ticket is built.
The exception is a snippet a prototype produced that states a decision more precisely than prose can, such as a state machine or the shape of a type, trimmed to the decision and marked as coming from the prototype.

## Blockers and waves

Every ticket names the tickets that block it: the ones that must be done before it can start, and no others.
That list is the whole dependency graph.
Name a blocker only where the ticket cannot start without its work.

A ticket's wave is derived from its blockers rather than chosen: one more than the highest wave among its blockers, and 1 for a ticket with none.
A wave holding one ticket names it `N-<task>.md`, and a wave holding several names them `N.1-<task>.md`, `N.2-<task>.md` and on, in the order the breakdown listed them.
`<task>` is a few words in the glossary's vocabulary, lower case and joined by hyphens.

The name orders the folder and gates nothing.
A ticket starts once its own blockers are done, so `3.1` blocked only by `2.1` starts while `2.2` is still being built.

A name is fixed once the tickets are written.
Abandoning a ticket later passes its blockers to the tickets it blocked and renames nothing, so no link to a ticket breaks.

## Putting the breakdown on a page

Never put a breakdown to a person in the terminal.
Put it on a page: `socrates agora serve --document ./breakdown.json` prints the loopback address of its page, and `socrates agora poll --session <session>` hands back what the person answered.

The document is one card per proposed ticket, in the order of their waves, each identified by the file name it would be written to:

    {"title": "...", "preamble": "...",
     "decisions": [{"id": "2.1-parse", "question": "2.1: ...", "meaning": "delivers: ... stories: 3, 4. blocked by: 1-prefactor",
                    "mode": "single-choice",
                    "options": [{"label": "keep it", "meaning": "cut as proposed", "recommended": true}]}]}

Each card states the ticket's title, what it delivers from the user's side, the spec's user stories it delivers, what blocks it and anything it needs a person for, and offers *keep it*, marked recommended so it arrives chosen: a person who agrees with the whole breakdown presses once.
A merge, a split or a different blocker is asked for in the words field of the card concerned, after taking *None of these*, which the page already renders, so never write that option yourself.
A breakdown answer runs nothing, so no option carries a command, and no card states a severity.

The preamble asks two things of the whole breakdown: whether it is cut too coarse or too fine, and whether every blocker gates the ticket it blocks.
It also names every user story no card delivers.

Mint the session once, with the first `socrates agora serve`, so the person is given one address.
Keep that session: serve every revision to the page the person already has with `socrates agora serve --session <session>`.
Revise the breakdown from what they wrote and serve it again, until a round comes back with every card kept and nothing written beside it.

## The last card

Once every card comes back kept, the last page is one card, and it is the only card on that page: every file the folder will hold, the main ticket among them, and the folder they go to.
A person approves a location rather than discovers one:

    {"title": "one card before the tickets are written",
     "decisions": [{"id": "write", "question": "write the tickets now",
                    "meaning": "path: docs/changes/<slug>/tickets/. files: README.md, 1-prefactor.md, 2.1-parse.md, 2.2-render.md, 3-cli.md",
                    "mode": "single-choice",
                    "options": [{"label": "write them", "recommended": true}]}]}

Nothing is written until that card is approved.
A person who wants something else takes *None of these* and says what in the words field.

## The main ticket

`README.md` in the tickets folder is the main ticket: the first file whoever implements the change opens, and outside the numbering, so no number has to mean two things.
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
The implementing session is the one writer of that table.

## A ticket

Each task ticket is a file of its own, and holds:

    # 2.1: <title>

    ## What to build

    <the behaviour this ticket makes work, from the user's side, not a list of layers>

    ## Acceptance criteria

    - [ ] <a behaviour somebody can check>

    ## User stories

    - 3, 4

    ## Needs a person

    <what building it needs that coding cannot do: a credential, an account, a permission or a manual step>

    ## Blocked by

    - 1-prefactor.md

A ticket that needs nothing coding cannot do leaves out Needs a person.
A ticket with no blockers says None under Blocked by, and can start at once.
A ticket carries no status: its status is the main ticket's.
Its acceptance criteria are behaviour to check and nothing more: a ticket names no test and no scenario, because a test's name is code and belongs to whoever builds it, and a repository that wants its tests named for scenarios says so in its own instructions.
Each criterion names what somebody would observe were it false, and is false at the commit the ticket starts from: one already true there, one only another ticket's work can meet, or one that restates the title grades nothing.

## Leaving them to be committed

Leave the tickets uncommitted, as a design leaves its spec: committing is the person's.
Say what was written and where, and say that the folder must be committed before any ticket is handed out, because an implementer working in a worktree of its own reads only what is committed.

## Installing elsewhere

This document installs under `~/.claude/skills` beside `socrates-design`.
Write both somewhere else instead with `socrates design skill --to <directory>`.
