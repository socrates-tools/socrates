# socrates: design, tickets and implement take what their prior art already learned

Decisions: none.

## Problem Statement

`socrates-design`, `socrates-tickets` and `socrates-implement` each started from one of mattpocock's skills: grilling and domain modeling, `to-tickets`, and `implement-spec`.
Each took that skill's core and replaced what did not fit this suite: the page instead of the scrollback, committed files instead of a tracker, a person starting every skill, and an overseer who writes no code.

Those skills have since collected lessons from the people using them, and some of those lessons are missing here.
None of them is about the page, the files or who starts a skill, so taking them changes nothing socrates was built to be different about.

**A design stalls on what it looks up, and accepts the words it is given.**
A session that reads the repository before a round holds back every question until it has read everything, though most of them do not depend on what it is reading.
The glossary is checked only against itself: a fuzzy word, a boundary nobody tested with an example, or a claim the code contradicts reaches the spec unchallenged.
Nothing asks where the change will be tested, so the spec's Testing Decisions are the session's guess.
A card offering a choice of interface offers whatever the session thought of first, and the ADR's Rejected section can only be as good as the options that card offered.
`socrates-tickets` lets a ticket carry a snippet a prototype produced, and nothing in a design ever produces a prototype.

**A breakdown can be cut finer than the change, and its criteria can grade nothing.**
A change that fits one context is still cut into several tickets, because nothing says that one ticket is a valid breakdown.
A blocker that does not gate its ticket serialises work the overseer could have run side by side, and nothing asks the person whether each one gates.
An acceptance criterion that is already true where the ticket starts, that only another ticket can meet, or that restates the title, is ticked without anybody building anything.
Nothing connects a ticket to the user stories it delivers, so a story no ticket covers is found at the last read, if at all.
A ticket that needs a credential or a manual step is discovered to need one halfway through a run.

**An implementation explores the same code once per implementer, and reviews without saying which line was missed.**
Every implementer reads the code its ticket touches from nothing, and two of them can still pick two names for one new thing on a shared surface.
The overseer's feedback says what falls short, but not which line of the ticket, the spec or the repository's rules it falls short of.
An implementer in a repository that states no way of testing has no discipline to fall back on.
`git stash` is shared by every worktree of one repository, so one implementer can pop what another stashed.

## Solution

Each skill gains the lessons that fit it, in its own words and under its own ceiling.

`socrates-design`:

- Facts are looked up by subagents while the round goes out: only the questions that depend on a running lookup wait for it.
- The glossary is challenged rather than only checked: a fuzzy word, an invented scenario at the boundary between two terms, and a claim the code contradicts each go on a card.
- An entry says what a term is, holds no implementation detail, and names the synonyms it gives up.
- Once a round has settled what the change touches, a card asks where it is tested: existing seams before new ones, the highest one that works, ideally one.
- Where a question is the shape of an interface or a module, subagents each draft one under a different constraint, and the drafts that really differ become the card's options.
- Where a card cannot settle a question in words, a throwaway prototype answers it outside the repository, and what it showed goes on the card.
- The decisions section names what usually earns an ADR, without loosening the three-part test.

`socrates-tickets`:

- A change that fits one fresh context is one ticket.
- A blocker is named only where the ticket cannot start without it, and the breakdown's preamble asks whether the cut is too coarse or too fine and whether every blocker gates.
- Each ticket names the spec's user stories it delivers, on its card and in its file, and a story no ticket delivers is named on the page.
- Each acceptance criterion names what would show it false and is false where the ticket starts.
- A ticket that needs a person says what, under a heading of its own.

`socrates-implement` and `socrates-implementer`:

- Before the first hand-out one subagent writes notes about the code the change touches, outside the repository, and every brief points at them.
- A ticket that needs a person is asked about before it is handed out, while the rest keep working.
- Review keeps two counts apart: what the ticket, spec and ADRs ask, quoting the line missed, and the repository's own standards, citing the rule.
- The last read checks the spec's user stories one by one.
- The implementer never uses `git stash`, and where the repository states no way of testing it writes each test first and watches it fail.

## User Stories

**Somebody being interviewed**

1. As somebody answering a design round, I want the questions that do not depend on what the session is reading to reach me straight away, so that a slow lookup does not cost me a round.
2. As somebody answering a design round, I want a word I used loosely to come back to me as a choice between what it could mean, so that the glossary holds the word I meant.
3. As somebody answering a design round, I want an invented example at the edge of two terms, so that I decide where one stops and the other starts before the code does.
4. As somebody answering a design round, I want to hear when the code already does something other than what I said, so that the spec does not rest on a claim the repository contradicts.
5. As somebody answering a design round, I want to be asked where the change will be tested, so that the spec's testing decisions are mine rather than the session's guess.
6. As somebody choosing an interface, I want options drafted under different constraints, so that I am choosing between real alternatives and the Rejected section records something worth having beaten.
7. As somebody answering a question words cannot settle, I want a prototype to show me what each answer does, so that I decide from behaviour rather than from a description of it.
8. As somebody reading the ADRs later, I want the scope decisions, deliberate deviations and invisible constraints written down, so that nobody fixes what was done on purpose.

**Somebody approving a breakdown**

9. As somebody approving a breakdown of a small change, I want one ticket proposed where one will do, so that I am not merging twelve cards back into one.
10. As somebody approving a breakdown, I want to be asked whether every blocker gates, so that the overseer is not made to wait on work a ticket does not need.
11. As somebody approving a breakdown, I want each card to name the user stories it delivers and the page to name any story nobody delivers, so that nothing the spec asked for is left out of the cut.
12. As somebody approving a breakdown, I want to see which tickets need me, so that I know before the run starts what I will be asked for.

**Somebody whose change is being built**

13. As the person whose change is built, I want a ticket that needs a credential asked about before it is handed out, so that no implementer stalls halfway through waiting for me.
14. As the person whose change is built, I want every criterion to be something that fails before the work and holds after it, so that a ticked box means something was built.
15. As the person whose change is built, I want the last read to go story by story, so that a story no ticket delivered is caught before the branch is called ready.

**The overseer**

16. As the overseer, I want notes about the code the change touches written once, so that each implementer starts from what the first reader found and two implementers use one name for one new thing.
17. As the overseer, I want my feedback to quote the line it falls short of, so that the implementer patches what I meant rather than what it guessed I meant.

**An implementer**

18. As an implementer, I want a default way of testing where the repository states none, so that a test proves the behaviour rather than repeating the code.
19. As an implementer sharing a repository with others, I want never to stash, so that what I set aside cannot be popped into somebody else's worktree.

## Implementation Decisions

**Every change is to a rendered document.**
`design/skill/messages.go` and `implement/skill/messages.go` hold the sections, and `make skill` regenerates the committed copies under `.claude/skills/` and `.claude/agents/`.
No command, flag or setting is added, so no skill names an invocation the binary does not offer.

**Each ceiling rises by what its document gained, and no further.**
`design/skill.Ceiling`, `design/skill.TicketsCeiling`, `implement/skill.Ceiling` and `implement/skill.AgentCeiling` are each set just above the rendered size, per ADR 21, so the next addition is again a deliberate raise.

**The facts section of `socrates-design` gains lookups that do not block a round.**
A lookup runs in a subagent, a running lookup is an unsettled prerequisite, and only the questions downstream of it wait for it.

**The glossary section of `socrates-design` gains three challenges and the shape of an entry.**
A fuzzy word, an invented scenario at the boundary between two terms, and a claim the code contradicts each go on a card of their own.
An entry says what a term is rather than what it does, holds no implementation detail, and names the synonyms it gives up.

**`socrates-design` gains a section, Seams and alternatives, between the glossary and the spec.**
It holds the seams card, the drafts under different constraints, and the prototype.
The prototype is written outside the repository and removed once its question is settled, because a design leaves nothing behind but the spec, the ADRs and the glossary, and anything else uncommitted stops `socrates-implement` from starting.
What the prototype showed goes on the card, through an illustration where it is a shape, and the snippet that states the decision is what the spec keeps, marked as coming from the prototype, which is the snippet `socrates-tickets` already admits.

**The decisions section of `socrates-design` lists what usually earns an ADR, inside the three-part test.**
A scope decision, a deliberate deviation and a constraint the code does not show are named as where all three often hold, and each still has to pass all three.

**`socrates-tickets` gains its rules in the sections they belong to.**
Slicing gains one ticket for a change that fits one fresh context.
Blockers and waves gains that a blocker is named only where the ticket cannot start without it.
The breakdown page gains the two questions of the preamble, the stories each card delivers, and naming a story no card delivers.
A ticket gains two headings, User stories and Needs a person, and the falsifiable criterion.
Needs a person is written only where building the ticket needs what coding cannot do, and is left out otherwise.

**`socrates-implement` gains notes, an early question, two counts and stories in its last read.**
The notes are written by one subagent before the first hand-out, to `notes/<slug>.md` beside the `worktrees/` directory, which is `~/.socrates/implement/notes/` by default.
They are outside the repository, so they are never committed and never swept into the change's branch, and a resumed run reads the notes already there.
The overseer still edits only the main ticket's status table: she writes no notes herself.
A brief's pointers gain the notes.

**The implementer's agent definition gains two rules.**
Never `git stash`, because the stash is shared by every worktree of the repository.
Where the repository states no way of testing: write each test before the code and watch it fail, test behaviour through the public interface, take expected values from the ticket or a worked example rather than computing them the way the code does, and mock only at the system's edges.
The repository's own instructions still win wherever they say anything.

**`CONTEXT.md` gains prototype, notes and needs a person.**

## Testing Decisions

Every promise is held by a test that asserts the statement in the rendered document, in the section it belongs to, the way `design/skill` and `implement/skill` already test every other promise.
The seam is the rendered document: `skill.Document()`, `skill.TicketsDocument()` and `skill.DefaultDefinition()`.
Each test is written first, fails against the current document, and passes once the section says it.
The existing gates hold the rest: the committed copies match what is rendered, every invocation named exists, no option carries a command or a severity, the agent definition names no socrates command, and each document is under its ceiling.

## Out of Scope

- A way to split an oversized ticket mid-run. The overseer still edits only the status table.
- Anything from mattpocock's skills that is about trackers, triage, handoff, retrospectives or a skill starting another skill.
- An end-of-run review by `socrates review`, which ADR 49 already rejects.

## Further Notes

The prior art is mattpocock's skills repository at v1.3: grilling, domain-modeling, codebase-design, prototype, to-spec, to-tickets, implement-spec, tdd and code-review, and the "common questions" pages its documentation keeps for each.
