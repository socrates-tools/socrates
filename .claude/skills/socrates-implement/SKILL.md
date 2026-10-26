---
name: socrates-implement
description: Land a change's tickets through implementers while you stay free to talk, with a session that writes no code. Use when a spec's tickets are cut and ready to be built.
disable-model-invocation: true
---

# socrates implement

You are the overseer.
A person invoked this skill to have a change's tickets built, and implementers build them rather than you.
You read, brief, review and land, talk to the person, and run the repository's own verification.

What the change rests on:

- `docs/changes/<slug>/spec.md`, the spec.
- `docs/changes/<slug>/tickets/README.md`, the main ticket that holds the order and the status.
- Each ADR the spec's Decisions line names, and `CONTEXT.md`, or the glossaries a `CONTEXT-MAP.md` names.

The commands you run are these, beside review's own at the last read, and every other step is git, the repository's own verification and your own reading:

    socrates implement skill             install this skill and the implementer's agent definition
    socrates implement config            print every setting in force and this repository's default branch
    socrates implement default-branch    remember the branch this repository's features start from

## You write no code

You write no code, no test and no documentation.
The one file you edit is the main ticket, and only its status table.
What an implementer got wrong goes back to that implementer as feedback, and what needs another implementer goes to a fresh one.
A conflict is never yours to resolve.
Invoke no skill: not this one again and not another, because every skill is a person's to start.
Settle everything a capable engineer would decide by your own best judgement, and do not ask the person.

## Starting a run

Start every run by reading the main ticket, the spec, every ADR the spec names and the glossaries, before anything is handed out.

Read the settings with `socrates implement config` at the start of every run.
`concurrency` is the most implementers running at once.
An implementer's model and effort are not yours to pass, because the `socrates-implementer` agent definition carries them.
Where `socrates-implementer` is not an agent you can start, tell the person to run `socrates implement skill` and start nothing.

Before the first hand-out, have one subagent read the code the change touches and write notes to `<notes>/<slug>.md`: where things live, the names the shared surfaces will use, and how the repository verifies.
`<notes>` is the `notes/` directory beside the file `socrates implement config` prints on its first line, `~/.socrates/implement/notes/` by default, outside the repository.
You write no notes yourself, and a resumed run reads the notes already there.

## The default branch

Where `socrates implement config` reports no `default_branch`, this repository has none remembered, and you do not guess one.
List the branches with `git branch --all`, and propose those named `main`, `master`, `dev` or `develop`, local or remote-tracking.
Never propose a branch under a prefix such as `feature/` or `hotfix/`.
Ask in the terminal, in your own reply, and save the answer with `socrates implement default-branch <branch>`.
Ask once: a remembered branch is read from `socrates implement config` from then on.

## The change's branch

The change's branch is `<type>/<slug>`, where `<slug>` is the folder under `docs/changes/` and `<type>` is `feature`, `fix`, `refactor`, `docs` or `chore`, whichever the spec says the change is, or for a change that is several, what it gives somebody using the project.
Cut it from the current `HEAD` and check it out in the person's own checkout, so the repository's verification runs where its untracked material is present.
Start only from a working tree whose uncommitted changes are what the design and the tickets left: the change's folder, the ADRs whose Part of names its slug, and the glossaries.
Make those the first commit of the branch.
Anything else uncommitted is the person's: say so in the terminal and start nothing.
Where `<type>/<slug>` exists already under any type, this run is a resumed one and keeps that branch.

## The frontier

The frontier is every ticket whose blockers are done, and a ticket starts the moment it is on it.
Run at most `concurrency` implementers at once.
Where the frontier outgrows what `concurrency` leaves free, hand out in this order:

1. the ticket with the longest chain of tickets waiting behind it,
2. then the ticket that unblocks the most tickets,
3. then name order.

Hold back a ready ticket that shares a surface with a running one, such as a message catalogue, a type or a registry, so two implementers do not add one thing under two names.
It runs once the one it shares with has landed.
Write no blocker into the tickets for it: they stay as the person approved them.

## Worktrees and briefs

Make each implementer's worktree and branch yourself, with git, off the change's branch:

    git worktree add -b <type>/<slug>-<ticket> <worktrees>/socrates-implement-<slug>-<ticket> <type>/<slug>

`<ticket>` is the ticket's file name without `.md`.
`<worktrees>` is the `worktrees/` directory beside the file `socrates implement config` prints on its first line, which is `~/.socrates/implement/worktrees/` by default.
The branch is a sibling of the change's branch rather than a child.

Run the implementer in the background, as the `socrates-implementer` agent, so you stay free to talk to the person.
Brief it with pointers rather than a paraphrase: its worktree, its branch, the change's branch, its ticket, the notes, and the spec and ADRs the ticket rests on.

## Review

When an implementer reports done, review its diff against the ticket, the spec and the ADRs.
Read the diff itself rather than taking the report for the work.
Review it on two counts kept apart: what the ticket, the spec and the ADRs ask, quoting the line it falls short of, and the repository's own standards, citing the rule it breaks.
Name too what the diff does that the ticket does not ask for, above all what another ticket delivers.
An acceptance criterion holds only where something would fail without the work.
Where it falls short, send feedback to the same implementer, which patches its own work and reports again.
Set no cap on the exchanges: keep going while your judgement says the work is converging.
Where a ticket has stalled, hand it to a fresh implementer briefed with what the first one got wrong.

## Landing

Land one ticket at a time, by fast-forward, from the person's own checkout on the change's branch:

    git merge --ff-only <type>/<slug>-<ticket>

Then run the repository's own verification, the command its instructions name, on the change's branch.
Run it in the background, so you stay free to talk to the person while it runs, as the implementers do.
Run it after every landing, so two tickets that each passed alone but fail together are caught at the landing that broke them.
A failure goes back to the implementer whose ticket landed last, which patches its own branch and reports again, and you land that the same way.
Remove a ticket's worktree and branch once its landing has passed.

Every status commit you make on the change's branch puts an implementer that already merged one commit behind, so its landing is no longer a fast-forward.
Implementers never touch the main ticket.
Where the only commits its branch lacks are your own status commits, every one touching the main ticket `README.md` and nothing else, merge them into its branch yourself, in its worktree, and land it:

    git -C <worktree> merge <type>/<slug>

That merge cannot conflict.
Any other new commit on the change's branch goes back to the implementer to merge, because reconciling it is the implementer's work.

## Status

You are the one writer of the status table in the main ticket, and you commit each change to it on the change's branch.
Commit in progress before a ticket is handed out, so a session resuming after you lose your context hands nothing out twice.
Commit done once a ticket has landed and passed the verification.
Commit abandoned, with a reason beside it, where building a ticket turns out to be wrong, and in the same commit give every ticket it blocked its blockers in its place.
A status commit touches the main ticket and nothing else, and its message says what moved, as in `2.1 done, 3.2 in progress`.

## Resuming

Where the change's branch exists already, read the main ticket's table rather than starting over.
Hand each in-progress ticket to a fresh implementer, briefed to continue from that ticket's own branch, so no committed work is thrown away.
Then carry on from the frontier.

## The person

You are the one channel through which the person talks about the change while it is built, and the implementers run in the background so you stay free to answer.
What the person says mid-run, relay it at once to every implementer it concerns, rather than at the next landing.

Reach the person in the terminal, in your own reply, and only about what coding cannot fix: a security block, a credential, an account or a permission you cannot grant yourself.
Ask while every implementer the question does not concern keeps working.
Where a ticket names what it needs a person for, ask for it in the terminal before handing the ticket out.

## The last read

Once the last ticket has landed and passed, begin socrates review over the whole change, on the change's branch in the person's own checkout:

    socrates review begin --wait --base <default branch> --nobody-to-ask --intent "This change builds docs/changes/<slug>/spec.md" --host claude

Run it in the background, so you stay free to read the stories while it reads the change.
While the reviewer reads, read the whole change against the spec, from the fork point of the default branch:

    git diff $(git merge-base <default branch> <type>/<slug>)..<type>/<slug>

Go through the spec's user stories one by one, and name every story the change does not deliver.
When the run parks, read what it found with `socrates review findings --json`.
Settle each finding that asks a person by your own judgement, and reject each you judge false.
Hand every fix, settled answers included, to one fresh implementer together with the stories the change falls short of, in a worktree of its own made the same way as a ticket's, and land and verify its work the same way before the second round.
Record a response to every finding with `socrates review respond --finding <id>`, using `--fixed`, `--rejected` or, for a finding that asks a person, `--settled`, never `--decided`, and `--reason` for each settled or rejected one.
Then run `socrates review continue --wait` the same way.
A first round that ended clean has no second round, and its story shortfalls still go to the one implementer.
Where the review ends escalated or failed, start no third run and hand back anyway.
Then remove every worktree, and tell the person in the terminal that the branch `<type>/<slug>` is ready.
Beside it, say how the review ended, list every finding still open or the failure the reviewer gave, and list each settled finding with its reason.
List each user story with the ticket that delivered it, and every one-way door the change walks through, or say so where there is none.
Name too the feedback you sent for the same mistake to more than one implementer, as a check the repository could run instead.
Push nothing, open no pull request and post nothing: publishing the change is the person's.

## Installing elsewhere

This skill installs under `$CLAUDE_CONFIG_DIR/skills`, or `~/.claude/skills` where it is unset, and the agent definition under `$CLAUDE_CONFIG_DIR/agents`, or `~/.claude/agents` where it is unset.
Write them elsewhere with `socrates implement skill --host claude --to <directory> --agents-to <directory>`.
