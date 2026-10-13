# socrates design

`socrates design` is a skill rather than a tool.
The namespace holds one operation, `socrates design skill`, which installs its two skills: `socrates-design`, which runs the interview, and `socrates-tickets`, which cuts the spec it left into tickets.
A design interviews a person over agora and leaves a spec and its decisions behind, and it knows nothing about review or implement.

This file is design's layout and invariants.
`/AGENTS.md` holds the suite's, `CONTEXT.md` here holds design's vocabulary, and `docs/adr/` here holds design's decisions.

## Layout

- `design/command` - the command tree: one operation, `skill`, that installs both skills.
- `design/skill` - the source of both skills. The design's: the rounds, the page, what was settled, the glossary, the spec, the decisions, when the frontier empties, ending early, and installing elsewhere. The tickets': the spec it starts from, slicing, blockers and waves, the breakdown on a page, writing them, the main ticket, a ticket, and leaving them to be committed.

## Invariants

Each is held by a test that fails when it stops holding.
None is a matter of care, and the reasoning for each is in this directory's `docs/adr/`.

**The design skill**

- `design/skill` is tested through its rendered documents: their ceiling is enforced, the committed copies must match what the package renders, and every invocation it names must exist in the binary's command surface with the flags it names, across both namespaces it covers.
- No page asks a person to approve writing: once the frontier is empty, or every ticket card comes back kept, the session closes the page with `socrates agora end` and writes, and starts no other skill, per ADR 49: every socrates skill is a person's to start.
- The skill writes no severity, no `None of these` option and no command on any option it composes: a design answer runs nothing, and a design question is earlier or later in the argument rather than more or less severe.
- A settled question is written to the document's `settled` list rather than as a decision, stating the question, the labels chosen and the person's own words, and a person reopens one through the page's message, per ADR 46.
- The ADR criterion the skill states is a conjunction, tested as one: hard to reverse, surprising without context, and the result of a real trade-off, never any one or two of the three.
- Facts are looked up while the round goes out, so only the questions downstream of a running lookup wait for it.
- The glossary is challenged on cards: a fuzzy word, a scenario at the boundary between two terms, and a claim the code contradicts.
- A card asks where the change is tested, an interface question offers drafts made under different constraints, and a prototype is written outside the repository and removed once its question settles.
- There is no round cap: a person can end a design at any round, on the page or with `socrates agora end`, and either way the exchange ends carrying only what was actually answered, leaving everything else for the session to record as unsettled.
- The loop is driven end to end through the real binary: a round renders with its recommendation chosen, and what the person answers reaches the session through `socrates agora poll`, and a page the session ends once nothing is left to ask renders closed. What the binary cannot drive, because nothing in it writes a file, is the spec and the ADRs themselves: those are the session's to write from what `socrates agora poll` handed it.

**The tickets skill**

- `socrates-tickets` is a second skill of `design/skill`, installed by `socrates design skill` beside the first, held to a ceiling of its own and to the same gates: the committed copy matches its source, and every invocation it names exists.
- It says nothing about a design handing on to it and mints its own session, per ADR 49.
- It writes no severity, no `None of these` option and no command on any option, and names no `socrates review` invocation: how a ticket is implemented is the implementing session's own.
- It cuts from a spec under `docs/changes/<slug>/` and nothing else, refuses a tickets folder that already exists by name, and leaves what it writes uncommitted.
- A ticket names its blockers and its wave is derived from them, per ADR 44, and status is kept in the main ticket and nowhere else.
- A change that fits one fresh context is one ticket, and a blocker is named only where its ticket cannot start without it.
- Every ticket names the spec's user stories it delivers, the page names any story no card delivers, and a ticket that needs a person says what for.
- Every acceptance criterion names what would show it false and is false where its ticket starts.
- The breakdown is driven end to end through the real binary: every card arrives with *keep it* chosen, and a card answered in the person's own words reaches the session through `socrates agora poll`. What the binary cannot drive is the tickets themselves, which are the session's to write.
  The skill says to copy each one exactly as `socrates agora poll` handed it back, because agora refuses one restated otherwise, per ADR 59.
