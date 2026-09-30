# socrates: tickets, cut from a spec into files an implementer works through

Decisions: ADR 43, ADR 44.

## Problem Statement

A design leaves a spec behind, and nothing in the suite turns that spec into work somebody can pick up.

The spec says what a change is for and what was decided about it.
It does not say what to build first, what can be built at the same time as something else, or how far an implementation has got.
Today a session implementing a spec works that out for itself, in its own context, and loses it when the context goes.

The shape that answers this already exists outside the suite.
mattpocock's `to-tickets` skill breaks a spec into tracer-bullet tickets, each a narrow vertical slice that is demoable on its own, each naming the tickets that block it.
It publishes them to a tracker, or to uncommitted files under `.scratch/`, and asks a person to approve the breakdown in the terminal.

Three things about that shape do not fit how this suite is used.

**The tickets live somewhere an implementer may not see.**
A tracker needs a network and an account, and an uncommitted `.scratch/` folder is invisible to an implementer working in a worktree of its own.
The person using socrates does not use a tracker for tickets at all.

**Nothing says where the implementation stands.**
A set of tickets with blocking edges says what may start, not what has started.
An implementing session that orchestrates several tickets at once has no single place to record which it handed out, which came back, and which are still waiting, so a session that loses its context cannot tell a ticket nobody started from one a subagent is halfway through.

**The breakdown is approved in a scrollback.**
A person reading a numbered list of tickets in a terminal and typing back which to merge is the interview `socrates design` was built to replace.

## Solution

`socrates-tickets` is a second skill of the `design` namespace.
`socrates design skill` installs it beside `socrates-design`.

A coding session reads it, reads a spec under `docs/changes/<slug>/`, and cuts it into tickets:

    socrates design skill          install both skills
    socrates agora serve           put the proposed breakdown in front of a person
    socrates agora poll            read what they answered

The session drafts tracer-bullet slices on `to-tickets`' rules, gives each ticket its blockers, and derives each ticket's wave from them.
It puts the breakdown on a page, one card per proposed ticket with *keep it* already chosen, so a person who agrees with the whole breakdown presses once, and asks for a merge, a split or a different blocker in their own words on the card concerned.
Once every card comes back kept, the last page is one card naming the folder and every file it will hold, so a person approves a location rather than discovers one.

A design can hand straight on to it: the last card of `socrates-design` offers writing the spec and cutting its tickets next, in the same session and on the same page.

What it leaves behind, beside the spec, for the person to commit:

    docs/changes/<slug>/tickets/
      README.md             the main ticket
      1-<task>.md
      2.1-<task>.md
      2.2-<task>.md
      3-<task>.md

The main ticket is what an implementing session opens first.
It says in a few sentences what the change is for, points at the spec and the ADRs rather than repeating them, lists every ticket with its blockers, states the order, and holds each ticket's status.
The implementing session is the one writer of that status.

## User Stories

**Somebody approving a breakdown**

1. As somebody approving a breakdown, I want it on a page rather than in a scrollback, so that I read the tickets rather than search a list for them.
2. As somebody approving a breakdown, I want one card per ticket, so that a change I ask for is attached to the ticket it is about.
3. As somebody approving a breakdown, I want every card to arrive with *keep it* chosen, so that agreeing with the whole breakdown costs one press.
4. As somebody approving a breakdown, I want to ask for a merge, a split or a different blocker in my own words, so that a breakdown that is nearly right is not a trap.
5. As somebody approving a breakdown, I want each card to say what the ticket delivers and what blocks it, so that I judge the granularity and the edges together.
6. As somebody approving a breakdown, I want the same page to carry the revised breakdown, so that a second look costs me no second address.
7. As somebody approving a breakdown, I want nothing written until I approve it, so that the files in the repository are the breakdown I agreed to.
8. As somebody approving a breakdown, I want one last card naming the folder and every file it will hold, so that I approve where the tickets go rather than discover it.
9. As somebody who has just finished a design, I want its last card to offer cutting the tickets next, so that the session that holds the argument is the one that slices it.
10. As somebody who already has tickets for a change, I want a second cut refused by name, so that nothing an implementation recorded is overwritten behind my back.
11. As somebody committing the result, I want the tickets left for me to commit, so that committing stays mine, as it is for the spec a design writes.

**An implementing session**

12. As an implementing session, I want one file to open first, so that I know what the change is for before I read any ticket.
13. As an implementing session, I want that file to point at the spec rather than repeat it, so that I read the decisions where they were made and not a paraphrase of them.
14. As an implementing session, I want every ticket listed with its blockers in one place, so that I can see the frontier without opening every file.
15. As an implementing session, I want a ticket to start the moment its own blockers are done, so that I do not wait on work it does not need.
16. As an implementing session, I want the number in a ticket's name to agree with its blockers when the tickets are cut, so that the listing and the graph start out telling me the same order.
17. As an implementing session, I want tickets that can run at the same time to share a wave, so that I can see at a glance what I may hand out together.
18. As an implementing session, I want to mark a ticket in progress before I hand it out, so that a session that resumes after losing its context does not hand it out twice.
19. As an implementing session, I want one place that says what is done, so that I never reconcile several files to learn where things stand.
20. As an implementing session, I want to record a ticket as abandoned with a reason, so that the table still says what was planned and why part of it was not built.
21. As an implementing session, I want to be told when a ticket counts as done, so that done means the same thing for every ticket.
22. As an implementing session, I want the tickets committed, so that an implementer in a worktree of its own reads the same tickets I do.
23. As an implementer given one ticket, I want the ticket to say what to build from the user's side, so that I build a slice rather than a layer.
24. As an implementer given one ticket, I want acceptance criteria I can check, so that I know when to stop.
25. As an implementer given one ticket, I want no file paths or code snippets in it, so that what I read has not gone stale since the ticket was cut.

**Somebody reading the repository later**

26. As somebody reading a change later, I want its tickets beside its spec, so that the work and the argument for it are one folder.
27. As somebody reading a change later, I want the history to show each ticket marked done, so that the order the work landed in is recorded without a tracker.

**A coding session cutting tickets**

28. As a coding session cutting tickets, I want to start from a spec in a change folder and nothing else, so that the slug, the folder and the ADRs come with it.
29. As a coding session cutting tickets, I want to be told to send anything that is not such a spec through a design first, so that no ticket is cut from an argument nobody settled.
30. As a coding session cutting tickets, I want the skill to name only commands the binary offers, so that an invocation I write is one that runs.
31. As a coding session cutting tickets, I want prefactoring to come first and a wide refactor to be sequenced as expand and contract, so that every ticket can land green.

**Somebody installing the suite**

32. As somebody installing the suite, I want the tickets skill installed by the same command as the design skill, so that the two cannot be at different versions.
33. As somebody installing the suite, I want one binary still, so that a further skill costs me nothing on my path.

## Implementation Decisions

**`socrates-tickets` is a second skill of the `design` namespace.**
`design/skill` renders a second document beside `socrates-design`, and `skill.All()` returns both, the way `review/skill` returns review's two.
Each carries a ceiling of its own, stated beside the document it bounds, per ADR 21.
`socrates design skill` installs both under the same directory and prints one line per skill it wrote, which its existing loop already does.
The namespace still holds one operation, so ADR 34's shape stands, and the dispatcher and the other two namespaces do not change.
The design namespace's help text changes from installing a skill to installing its skills, in `design/command/messages.go`.
Both committed copies are generated by `make skill`, and `.claude/skills/socrates-tickets/SKILL.md` joins the committed set.

**The input is a spec under `docs/changes/<slug>/`, and nothing else.**
The skill reads the spec, the ADRs its Decisions line names and `CONTEXT.md`.
A plan, an issue or a conversation is not an input: the skill tells the session to run a design first.
The slug and the folder come from the spec's own path, so the tickets folder is never a question.

**The slicing rules are `to-tickets`' rules.**
Tracer-bullet vertical slices, each demoable or verifiable on its own and sized for one fresh context; prefactoring first; a wide refactor sequenced as expand, migrate in batches, contract.
A ticket carries no file path and no code snippet, except a snippet a prototype produced that encodes a decision more precisely than prose.
Ticket titles and bodies use the glossary's words.

**Each ticket names its blockers, and its wave is derived from them.**
Per ADR 44, a ticket's wave is one more than the highest wave among its blockers, 1 where it has none.
A wave with one ticket names it `N-<task>.md`, a wave with several names them `N.1-<task>.md`, `N.2-<task>.md`, in the order the breakdown listed them.
The blocker list is the whole graph, and the name gates nothing.
A name is fixed once the tickets are written: abandoning a ticket passes its blockers to the tickets it blocked and renames nothing.

**The breakdown is approved on an agora page, one card per proposed ticket.**
Each card states the ticket's title, what it delivers and what blocks it, and offers *keep it*, marked recommended so it arrives chosen.
A merge, a split or a changed blocker is asked for in the card's words field.
The page follows every rule a design page follows: no severity, no command on any option, no *None of these* written by the skill, one session minted once and kept.
Where a design handed on to the skill, the breakdown is served to the design's own session, and a new one is minted only where there was none, or where that exchange ended or the session was collected.
A revised breakdown is served to the same session.
Once a round comes back with every card kept and nothing written beside it, the last page is one card, and the only card on it: every file name the folder will hold, the main ticket, and `docs/changes/<slug>/tickets/`, the same shape as a design's last round.
Nothing is written until that card is approved.

**The tickets are files under `docs/changes/<slug>/tickets/`.**
Per ADR 43, beside the spec, and published nowhere.
The skill leaves them uncommitted, as a design leaves its spec, and the main ticket says to commit them before any ticket is handed out, because an implementer in a worktree of its own reads only what is committed.

**An existing tickets folder is refused by name.**
Where `docs/changes/<slug>/tickets/` already exists, the skill writes nothing and names the folder.
A person who wants a fresh cut removes it themselves, knowing what its main ticket recorded, so no recorded status is ever overwritten.

**The main ticket is `README.md` in the tickets folder.**
It sits outside the numbering, so no number has to mean both a wave and something else.
It holds, in this order:

- A short introduction for whoever implements the change: what the change is for, in a few sentences, with a pointer to `spec.md` and to each ADR the spec names. It repeats no decision and argues for none, because that is the spec's and the ADRs' job.
- How to work through the tickets: commit the tickets folder before handing anything out; a ticket may start once every ticket it names as a blocker is `done`; mark it `in progress` before handing it out; mark it `done` once its acceptance criteria hold and its work is committed; mark it `abandoned`, with a reason in its row, where implementation finds it should not be built, and in the same change every ticket it blocked takes its blockers in its place, so nothing starts earlier than it would have.
- One table row per ticket: the file, what blocks it, and its status.

**Status lives only in the main ticket, as `pending`, `in progress`, `done` or `abandoned`.**
A task ticket carries no status line, so there is one place to read and one place to write.
An abandoned ticket keeps its row and carries a reason, so the table still says what was planned.
The word is `abandoned` rather than review's `dropped`, which already means something else in the glossary.
The implementing session is the one writer of the table.

**The main ticket names no review.**
How a ticket is implemented is the implementing session's own, and design knows nothing about review, so the main ticket names neither `socrates review` nor a loop of subagents.

**A task ticket is `to-tickets`' local template, less its status line.**
A title carrying its number, what to build from the user's side, acceptance criteria as checkboxes, and the tickets that block it by file name, or none.
Acceptance criteria are behaviour to check and nothing more: a ticket names no test and no scenario, because test names are code-shaped and belong to the implementer, and a repository that wants scenario-named tests says so in its own instructions.

**`socrates-design`'s last card offers cutting tickets next.**
The last round of a design gains the choice to cut the spec's tickets straight after writing it, in the same session and to the page the person already has.
Every combination with publishing is an option of its own: where the repository states a tracker the card offers write it, write it and open it as an issue, write it and cut its tickets next, and all three; where it states none, write it, and write it and cut its tickets next.
`socrates-design` therefore names `socrates-tickets`, the skill, and names no command of it, because the skill has none.

**`CONTEXT.md` gains ticket, blocker, wave, main ticket and abandoned, and change now holds its tickets as well as its spec.**
They were written as the terms settled.

**The map and the invariants follow.**
`AGENTS.md` names the second skill under what design is, and its invariants for the design skill gain the tickets skill's: the committed copy matches its source, its ceiling is enforced, every invocation it names exists, and it writes no severity, no command and no *None of these*.
`README.md` explains the second skill where it explains the first.

## Testing Decisions

The seam is the delivered document, which is where `design/skill`'s tests already sit, and the real binary for the page.

**The skill is asserted, never read as source.**
The tickets document is tested the way `socrates-design` is: its ceiling is enforced, the committed copy must match what the package renders, and every invocation it names must exist in the binary's command surface with the flags it names.
Prior art is `design/skill/ceiling_test.go`, `committed_test.go` and `name_test.go`.
`socrates design skill` installing both skills, each under its own name, is asserted where `name_test.go` asserts the first.

**What the skill promises is asserted as statements in the rendered document.**
One test per promise, named for its scenario: the input is a spec under a change folder and nothing else; a ticket names its blockers; the wave is derived from them and a wave of one is plain `N`; the tickets are committed under the change's tickets folder; the main ticket is `README.md`; status lives only there, in four states, the fourth being `abandoned` with a reason; the main ticket points at the spec rather than repeating it; the breakdown is approved on a page with one card per ticket and *keep it* chosen; writing is approved on one last card naming the folder and every file; an existing tickets folder is refused by name; the tickets are left uncommitted and the main ticket says to commit them first; a ticket's acceptance criteria name no test.
On the design side, `design/skill/lastround_test.go` gains the test that the last card offers cutting the tickets next.

**The prohibitions are asserted as prohibitions.**
The tickets document writes no *None of these* option, states no severity, puts no command on any option, and names no `socrates review` invocation.
Each is a test that fails when the document stops holding it, following `design/skill/prohibitions_test.go`, and each is broken on purpose once.

**The approval loop is driven end to end through the real binary.**
`internal/e2e` serves a breakdown page whose cards each offer *keep it* marked recommended, answers it in a real browser with one card kept and one carrying words, polls, and asserts what reached the session.
Prior art is `internal/e2e/designround_test.go`.
What the binary cannot drive, because nothing in it writes a file, is the tickets themselves: those are the session's to write from what `socrates agora poll` handed it, exactly as the spec is for a design.

## Out of Scope

- Publishing a ticket to any tracker.
- Implementing the tickets. The main ticket says how to read and mark them, and the implementing session's own skills do the rest.
- Naming a review, or any other workflow, in a ticket.
- Cutting tickets from anything but a spec in a change folder.
- A `tickets` namespace of its own.
- This repository's own earlier practice. ADR 39, `CONTRIBUTING.md` and `docs/agents/issue-tracker.md` say tickets are issues here; the person left reconciling them outside this change.

## Further Notes

The design was run on its own subject over two rounds and a last card: the questions were put on a page, and the person's answers are what this spec and ADRs 43 and 44 record.
None of round two's answers met all three of the ADR criteria, so they are Implementation Decisions here rather than files of their own.

The first round was nearly lost on the way, twice, and `docs/changes/the-page-in-a-devcontainer/report.md` records why.
The person also asked, in round two's message, that a settled card offer no words field, which the surface's invariants forbid today; `docs/changes/a-settled-card-offers-nothing/report.md` records that for a later change.
