# socrates: implement, an overseer landing a change's tickets through implementers

Decisions: ADR 48, ADR 49, ADR 50, ADR 51.

## Problem Statement

A design leaves a spec behind, and `socrates-tickets` cuts that spec into tickets with blockers, waves and a main ticket that holds their status.
Nothing in the suite builds them.

Today a person building a change from its tickets has two ways to do it, and neither fits.

**One session per ticket.**
The person opens a session, hands it a ticket, waits, reviews, clears the context, and opens the next one.
They are the dispatcher: they keep track of which tickets are unblocked, which are running and which came back, and they read every diff themselves.
A change of five tickets is five sessions of attention, and the person is disturbed for every one of them.

**One session for everything.**
The person hands the whole spec to one session, which writes every ticket's code in one context.
Its context fills with code, it reviews its own work, and it runs nothing in parallel, though the tickets' blockers say which could.

The nearest prior art outside the suite is mattpocock's `implement-spec`.
It reads tickets as a task graph, hands each ready ticket to a subagent in a worktree of its own, merges each into one integration branch, and runs one code review at the end.
Four things about that shape do not fit how this suite is used.

**The orchestrating session also writes code.**
Nothing stops it, so the session that holds the whole argument spends its context patching rather than judging, and the only review of a ticket is the one the end of the run gets.

**Its fix loop has no owner.**
A review at the end hands every finding to one fix subagent, and its users report review and fix loops running for hours, because nothing says who judges a fix good enough.

**The person is the one asked when anything is unclear.**
The user of socrates wants a change built with as little of their attention as possible, and only something coding itself cannot settle, such as a security block, is theirs.

**It assumes where a change starts.**
Claude Code's own worktree isolation branches from the repository's default branch, and people branch features off `dev` as often as off `main`.

## Solution

`socrates implement` is a fourth namespace of the one binary, per ADR 48.
`socrates implement skill` installs the `socrates-implement` skill and the `socrates-implementer` agent definition.
A person invokes the skill, and only a person does, per ADR 49.

    socrates implement skill            install the skill and the implementer's agent definition
    socrates implement config           print the settings in force, and the default branch of this repository
    socrates implement configure        change one setting from a menu
    socrates implement default-branch   remember the branch this repository's features start from

The session that invokes the skill is the overseer.
She writes no code.
She reads the main ticket, the spec and the ADRs, creates a branch for the change off the current `HEAD`, and hands each ticket on the frontier to an implementer: a subagent of model and effort the person configured, Sonnet at high effort by default, working in a worktree and on a branch of its own.
She runs at most as many implementers at once as the `concurrency` setting says, 4 by default, and picks the ticket with the longest chain of tickets waiting behind it first.

When an implementer reports done, the overseer reviews its diff against the ticket, the spec and the ADRs, knowing the change accumulates as more tickets land.
Where it falls short she sends feedback to the same implementer, which patches its own work and reports back, for as long as her own best judgement says the loop is converging.
Once her review passes, she fast-forwards the change's branch onto the implementer's branch, runs the repository's own verification there, marks the ticket done, and hands out whatever the landing put on the frontier.

The implementers run in the background, so the overseer is always free to talk to the person, and she is the one channel through which the person talks about the change while it is built.
What the person says mid-run she relays at once to every implementer it concerns.
Only what coding cannot fix reaches the person, and it reaches them in the terminal.

Once the last ticket lands, the overseer reads the whole change against the spec, measured from the repository's default branch, hands what falls short to one implementer, removes every worktree, and tells the person the branch is ready.
Nothing is pushed.

## User Stories

**Somebody whose change is being built**

1. As somebody whose change is being built, I want to invoke one skill once the tickets are cut, so that the change is built without my opening a session per ticket.
2. As somebody whose change is being built, I want the session I invoked to stay free to talk to me while the tickets are built, so that I can add an instruction without waiting for a ticket to finish.
3. As somebody whose change is being built, I want what I say mid-run to reach every implementer it concerns at once, so that no implementer keeps building to a brief I have since corrected.
4. As somebody whose change is being built, I want the overseer to be the one channel I talk through, so that I never have to find out which subagent is doing what.
5. As somebody whose change is being built, I want the overseer to settle every coding question by her own judgement, so that I am not asked about anything a capable engineer would decide.
6. As somebody whose change is being built, I want to be told only about what coding cannot fix, such as a security block, so that an interruption always means something I alone can unblock.
7. As somebody whose change is being built, I want that question asked in the terminal I invoked the skill from, so that I answer it where I am already looking.
8. As somebody whose change is being built, I want the work to land on a branch named for the change, so that my own branch is untouched until I decide what to do with it.
9. As somebody whose change is being built, I want nothing pushed, so that publishing the change stays mine.
10. As somebody whose change is being built, I want to be told once the branch is ready, so that I know when to read it.
11. As somebody whose change is being built, I want the whole change read against the spec once the last ticket lands, so that tickets that each passed do not add up to something the spec did not ask for.

**Somebody configuring the implementers**

12. As somebody configuring the implementers, I want at most 4 implementers running at once by default, so that a change with a wide frontier does not exhaust my account's limits.
13. As somebody configuring the implementers, I want to change that number from `socrates implement configure`, so that I change it the way I change a reviewer's model.
14. As somebody configuring the implementers, I want an implementer to run on Sonnet at high effort by default, so that the code is written by a capable model at a cost below the overseer's.
15. As somebody configuring the implementers, I want the implementer's model and effort to be settings too, so that I can move both when a model changes.
16. As somebody configuring the implementers, I want a change I make under `configure` to reach the next implementer spawned, so that I never wonder which model is building my code.
17. As somebody configuring the implementers, I want `socrates implement config` to report every setting with its value in force and whether I set it, so that it reads the way `socrates review config` reads.
18. As somebody configuring the implementers, I want socrates implement's settings in a file of its own, so that review's file and agora's file say nothing about building a change.

**Somebody whose repository has a branch of its own to start from**

19. As somebody whose features start from `dev` rather than `main`, I want socrates never to assume the repository's default branch, so that a change is never measured against the wrong branch.
20. As somebody using a repository for the first time with socrates, I want the overseer to propose the branches that look like a starting point, `main`, `master`, `dev` or `develop`, and never one like `feature/*` or `hotfix/*`, so that I pick rather than type.
21. As somebody who picked a default branch once, I want socrates to remember it for that repository, so that I am never asked again.
22. As somebody who picked a default branch, I want to see it in `socrates implement config` and change it in `socrates implement configure`, so that a wrong pick is mine to correct.
23. As somebody with a worktree of a repository, I want that worktree to read the same default branch as the clone it belongs to, so that I am not asked once per worktree.

**The overseer**

24. As the overseer, I want to read the main ticket, the spec and the ADRs before handing out anything, so that I review every ticket against the argument rather than against its own text.
25. As the overseer, I want the settings read from `socrates implement config` at the start of every run, so that a changed cap holds from the next run without anything reinstalled.
26. As the overseer, I want the frontier to be every ticket whose blockers are done, so that a ticket starts the moment it can.
27. As the overseer, I want to hand out the ticket with the longest chain of tickets waiting behind it first, so that the last ticket can start as early as the graph allows.
28. As the overseer, I want ties broken by how many tickets a ticket unblocks and then by name order, so that two runs over one graph hand out tickets in the same order.
29. As the overseer, I want to hold back a ready ticket that shares a surface with one already running, such as a message catalogue, a type or a registry, so that two implementers do not add one thing under two names.
30. As the overseer, I want holding one back to write no blocker into the tickets, so that the tickets as cut stay what the person approved.
31. As the overseer, I want each implementer in a worktree I made off the change's branch, on a branch I named, so that I know where every ticket's work is and can remove it when it lands.
32. As the overseer, I want to brief an implementer with pointers to the spec, the ticket, the ADRs and its worktree rather than with a paraphrase, so that its prompt is small and nothing in it has drifted from the source.
33. As the overseer, I want to review each implementer's diff myself, so that the review comes from the session holding the whole argument and not from the one that wrote the code.
34. As the overseer, I want to send feedback to the implementer that wrote the code and have it patch its own work, so that I never write code and the fix comes from the context that knows the code.
35. As the overseer, I want no fixed number of feedback exchanges per ticket, so that I can keep going where the work converges and change course where it does not.
36. As the overseer, I want to be free to hand a stalled ticket to a fresh implementer briefed with what went wrong, so that a context that has gone astray does not hold the ticket.
37. As the overseer, I want the repository's own verification run on the change's branch after every landing, so that two tickets that each passed alone but fail together are caught at the landing that broke them.
38. As the overseer, I want a failure after a landing to go back to the implementer whose ticket landed last, so that the fix comes from the work that broke it.
39. As the overseer, I want an implementer to merge the change's branch into its own before it reports done, so that landing it is a fast-forward and I resolve no conflict myself.
40. As the overseer, I want to commit an uncommitted tickets folder as the change branch's first commit, so that every implementer's worktree holds the tickets.
41. As the overseer, I want to commit in progress before handing a ticket out and done once it has landed, so that a session resuming after I lose my context hands nothing out twice.
42. As the overseer, I want to mark a ticket abandoned with a reason where building it turns out wrong, so that the main ticket still says what was planned and why part of it was not built.
43. As the overseer resuming after a lost context, I want a fresh implementer to continue each in-progress ticket from its own branch, so that no committed work is thrown away.
44. As the overseer, I want the default branch read from socrates rather than guessed, so that the whole-change read at the end measures the change from where it started.
45. As the overseer, I want a `socrates implement` command that saves the default branch the person picked, so that I record it without editing a file socrates owns.

**An implementer**

46. As an implementer, I want one ticket and pointers to everything it rests on, so that I build a slice and not the change.
47. As an implementer, I want to work in the worktree I was briefed with and nowhere else, so that nothing I do reaches another ticket's work.
48. As an implementer, I want to follow the repository's own instructions, including how it tests, so that the code I write is the code that repository asks for.
49. As an implementer, I want to run the repository's own verification before I report, so that I report done only for work that passes there.
50. As an implementer, I want feedback to arrive in the context I wrote the code in, so that I patch what I wrote rather than rediscover it.
51. As an implementer, I want to be told never to invoke a socrates skill, so that no skill starts another workflow underneath the overseer.

**Somebody installing the suite**

52. As somebody installing the suite, I want `socrates implement skill` to install the skill and the agent definition together, so that the two cannot be at different versions.
53. As somebody installing the suite, I want to install either somewhere else with `--to` and `--agents-to`, so that a project-scoped install works as the other namespaces' do.
54. As somebody installing the suite, I want one binary still, so that a further namespace costs me nothing on my path.
55. As somebody installing the suite, I want every socrates skill invoked by me only, so that no skill reaches for another and a workflow I did not start never starts.

**Somebody who has just finished a design**

56. As somebody who has just finished a design, I want its last card to offer writing the spec and publishing it, and nothing that starts another skill, so that I start cutting tickets myself with `/socrates-tickets` when I choose to.

## Implementation Decisions

**`socrates implement` is a fourth namespace, per ADR 48.**
`internal/dispatcher` mounts `implement/command`'s tree beside review's, agora's and design's, and stays wiring.
The tree holds `skill`, `config`, `configure` and `default-branch`, and nothing else.
`internal/state` gains `Implement` as a tool, so its home is `~/.socrates/implement` and `SOCRATES_HOME` moves it with the rest.
Every user-facing string lives in each package's `messages.go`, as every package's does.

The packages:

- `implement/command` - the command tree.
- `implement/config` - the one optional configuration file socrates implement reads, every setting it holds, the default each falls back to, and the repository an entry under `default_branches` belongs to.
- `implement/configure` - the menu over those settings, built on the machinery lifted out of `review/configure`.
- `implement/skill` - the source of the `socrates-implement` skill and of the `socrates-implementer` agent definition, and installing both.

**The settings are `concurrency`, `roles.implementer.model` and `roles.implementer.effort`, and the default branches.**
They live in `~/.socrates/implement/config.yaml`, per ADR 42, and that file holds only what a person set, per ADR 28.
`concurrency` is a positive whole number, default 4.
`roles.implementer.model` is a Claude Code model alias or a full model identifier, default `sonnet`, and `configure` offers the aliases.
`roles.implementer.effort` is one of `low`, `medium`, `high`, `xhigh` and `max`, default `high`.
`default_branches` maps a repository to the branch its features start from, per ADR 51.
A value the loader would refuse is refused before the file is touched, and a field the file does not have is refused naming the file, as review's is.
Each default is stated once, in `implement/config`, and what `config` reports is read from there.

**`socrates implement config` reports, in the shape `socrates review config` reports.**
It names the file it reads, then every setting by its YAML key with the value in force and whether a person set it.
Run inside a repository, it also reports that repository's default branch as `default_branch`, and an unremembered one prints empty rather than guessed.
It judges nothing.

**`socrates implement configure` is review's menu over implement's settings.**
The menu, the YAML document that keeps a person's comments, key order and blank lines through a write, the atomic write and the reading of an unreadable line leave `review/configure` for a package both menus build on, so the two cannot drift.
That lift is a prefactor of its own, and `socrates review configure` behaves exactly as before it.
Run inside a repository, the menu also offers that repository's default branch, where *use the default* removes the entry.
It refuses where nothing is a terminal, naming `socrates implement config`.

**An implementer's model and effort are written into an agent definition, per ADR 50.**
`socrates implement skill` writes `socrates-implementer.md` into Claude Code's agents directory, `~/.claude/agents` by default, with `model` and `effort` in its front matter taken from the settings in force, `background: true`, and the implementer's standing instructions as its body.
It records where it wrote the definition in `~/.socrates/implement/installed.yaml`, a file of its own because the configuration file holds only what a person set.
`socrates implement configure` writes the definition again wherever that record says, as soon as a person confirms a new model or effort, so the next implementer spawned runs on it.
Where nothing has been installed, `configure` writes the setting and says that the definition is written by `socrates implement skill`.
`--to` names the skills directory, as it does for the other namespaces, and `--agents-to` names the agents directory.

**The implementer's standing instructions live in the definition rather than in each brief.**
Work in the worktree you were briefed with and nowhere else; follow the repository's own instructions, including how it tests; commit on your own branch; run the repository's own verification before reporting; merge the change's branch into your own before reporting done; report what you did and what you did not; patch your own work when feedback arrives; never invoke a socrates skill.
The definition names no `socrates` command.

**Every socrates skill is invoked by a person only, per ADR 49.**
`socrates-review`, `socrates-review-handover`, `socrates-design`, `socrates-tickets` and `socrates-implement` each carry `disable-model-invocation: true` in their front matter.
The rule is stated once in `internal/skill`, and every namespace's skills are rendered through it.
`socrates-design`'s last card stops offering to cut the tickets next: it offers write it, and write it and open it as an issue where the repository states a tracker, and write it alone where it states none.
`socrates-tickets` stops naming a design that handed on to it, and mints its own session.
`socrates-implement` names no other socrates skill and no `socrates review` invocation, and nothing in it tells the overseer to invoke a skill.

**The overseer writes no code.**
She reads, reviews, briefs, runs the repository's verification, runs git to make worktrees, branches and fast-forwards, and commits the main ticket's status.
A conflict is never hers to resolve, because every landing is a fast-forward of a branch the implementer already merged the change's branch into.

**The change's branch is `socrates-implement/<slug>`, cut from the current `HEAD`.**
The overseer creates it and checks it out in the person's own checkout, so the repository's verification runs where untracked material the tests may read is present.
She starts only from a working tree whose uncommitted changes are what the design and the tickets left: the change's folder, the ADRs whose Part of names its slug, and `CONTEXT.md`.
Those are the branch's first commit.
Anything else uncommitted is the person's, and she says so in the terminal and starts nothing.

**Each implementer works in a worktree the overseer makes.**
She adds it with git off the change's branch, on `socrates-implement/<slug>-<ticket>`, where `<ticket>` is the ticket's file name without `.md`, under `~/.socrates/implement/worktrees/`.
A sibling rather than a child of the change's branch, because git cannot hold a branch and a branch beneath it.
She removes the worktree and its branch once the ticket lands, and every remaining one at the end.

**The frontier is every ticket whose blockers are done, and the cap bounds what runs.**
Where the frontier holds more tickets than `concurrency` leaves free, the ticket with the longest chain of tickets waiting behind it goes first, then the one that unblocks the most, then name order.
Two ready tickets the overseer judges to share a surface run one after the other, and she writes no blocker into the tickets.

**The review is the overseer's, and has no cap.**
She reviews each implementer's diff against the ticket, the spec and the ADRs, and sends feedback to the same implementer until her own judgement says it holds.
She may hand a stalled ticket to a fresh implementer briefed with what the first got wrong.
socrates review is not named, per ADR 49.

**Status follows the main ticket's rules, with the overseer as its one writer.**
In progress is committed on the change's branch before a ticket is handed out, and done once her review passed, the ticket landed and the verification passed there.
Abandoned carries a reason, and every ticket it blocked takes its blockers in its place, as `socrates-tickets` states.
A resumed overseer reads the table: each in-progress ticket is handed to a fresh implementer briefed to continue from its own branch.

**The person is reached in the terminal, and only about what coding cannot fix.**
The overseer escalates nothing a capable engineer would decide.
A security block, a credential, an account or a permission she cannot grant herself is the person's, asked in her own reply while every implementer it does not concern keeps working.

**The default branch is asked for once per repository, per ADR 51.**
Where `socrates implement config` reports none, the overseer proposes the repository's local and remote branches named `main`, `master`, `dev` or `develop`, never one under a prefix such as `feature/` or `hotfix/`, asks in the terminal, and saves the answer with `socrates implement default-branch <branch>`.
That command refuses a branch the repository does not hold, and writes the entry for the repository it runs in.
A repository is keyed by the top level of its main worktree, so every worktree of one clone reads one entry.
The whole-change read at the end diffs from `merge-base(<default branch>, socrates-implement/<slug>)`.

**`CONTEXT.md` gains overseer, implementer and default branch, and frontier is widened to tickets.**
They were written as the terms settled.

**The map and the invariants follow.**
`AGENTS.md` names the fourth namespace in what the suite is, its packages under what implement is, and its own file beside review's and agora's.
Its invariants gain a section for the implement skill and lose the design hand-on, and `README.md` explains the namespace where it explains design.

## Testing Decisions

The seams are the delivered documents, the configuration file as the loader reads it, and the real binary.

**The skill and the definition are asserted as documents, never read as source.**
`implement/skill` is tested the way `design/skill` is: each document's ceiling is enforced, the committed copies at `.claude/skills/socrates-implement/SKILL.md` and `.claude/agents/socrates-implementer.md` must match what the package renders at the defaults, and every invocation the skill names must exist in the binary's command surface with the flags it names.
Prior art is `design/skill/ceiling_test.go`, `committed_test.go` and `name_test.go`.

**What the skill promises is asserted as statements in the rendered document.**
One test per promise, named for its scenario: the overseer writes no code; she reads the cap from `socrates implement config` at the start of every run; the frontier is every ticket whose blockers are done; the longest chain goes first, then the most unblocked, then name order; a ready ticket sharing a surface is held back and no blocker is written; each implementer works in a worktree she made on a branch she named; landing is a fast-forward; verification runs after every landing; in progress is committed before a hand-out; a resumed overseer continues each in-progress ticket from its branch; the person is reached in the terminal and only about what coding cannot fix; what the person says mid-run is relayed to every implementer it concerns; the default branch is proposed from the named candidates and saved with `socrates implement default-branch`; the whole change is read against the spec at the end and nothing is pushed.

**The prohibitions are asserted as prohibitions, and each is broken on purpose once.**
The implement skill names no `socrates review` invocation and tells the overseer to invoke no skill.
The definition names no `socrates` command.
Every skill every namespace ships carries `disable-model-invocation: true`, asserted in `internal/skill` over the skills of all three namespaces that ship one.
`socrates-design`'s last card offers no cutting of tickets, in `design/skill/lastround_test.go`, and `socrates-tickets` names no design handing on to it.

**The settings are asserted through the loader and the real binary.**
`implement/config` is tested the way `review/config` is: each default is read from its owner, a reported default and its owner's agree, an out-of-range value is refused naming the file, and a field the file does not have is refused.
`socrates implement config` is driven in `internal/e2e`, inside and outside a repository, and from a worktree of a clone that reads the clone's entry.
`socrates implement default-branch` is driven against a real repository: it writes the entry for that repository, refuses a branch the repository does not hold, and leaves every other line of the file alone.

**The menu is driven by the real binary in a real pseudo-terminal.**
`socrates implement configure` is typed at and read back from what it drew, following `internal/e2e/configure_test.go`.
A confirmed model or effort lands in the file and in the installed definition at once; *use the default* removes the line; a menu run inside a repository offers its default branch.
`socrates review configure`'s own tests stay as they are and pass across the lift, which is what proves the lift changed nothing.

**Installing is asserted through the real binary.**
`socrates implement skill` writes the skill and the definition under their default directories and under `--to` and `--agents-to`, writes the definition's model and effort from the settings in force, and records where it wrote the definition.

**What the binary cannot drive is the overseer.**
Nothing in the binary spawns a subagent, so the orchestration itself, handing out, reviewing, relaying and landing, is the session's to perform from what the skill says, exactly as a design's spec is.
The suite asserts every command the skill names, and the document that says how to use them.

## Out of Scope

- Pushing the change's branch, opening a pull request, or posting anything to a forge.
- socrates review in any part of an implementation.
- Any harness but Claude Code: the skill and the definition install where Claude Code reads them, as every skill the suite ships does.
- A cap on the overseer's feedback to one implementer.
- Putting an escalation on an agora page.
- Detecting the default branch inside the binary: the candidates are the overseer's to propose.
- Renaming a ticket, or writing a blocker, during implementation.

## Further Notes

The design ran over four rounds and a last card on one page, and the person's answers are what this spec and ADRs 48 to 51 record.
It borrowed from `implement-spec` the task graph, the frontier, a worktree and branch per implementer, the implementer merging the branch tip before reporting so landing is a fast-forward, pointers rather than paraphrase in a brief, and the warning that worktrees postpone collisions rather than remove them.
It borrowed nothing of its review loop, its tracker, or its pull request.

The person asked on round one, in the page's message, that the main agent stay available while the subagents work, because it is the one channel by which they talk about the feature being built.
That is why implementers run in the background and why the overseer relays.

On q14 the person chose no option and wrote instead that socrates should not assume the repository's default branch, which is what ADR 51 records.
The session read that as the overseer making each worktree herself, because Claude Code's own isolation branches from the repository's default branch, and said so on the next page; nothing reopened it.

On q22 the person kept the word default branch and added that it may be renamed if it is confused with git's: "This is for socrates and if this default confuses with the "default" in git's world, we can rename this concept."
