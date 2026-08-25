## Why

A person read the first page a real run served and reported five things about it, and every one of them is the same failure seen from a different side: the page is written for the session that composed it rather than for the person deciding.

- Under every option sat the `socrates respond` command that option becomes, argument for argument, quoted reason and all.
  The person is not at a prompt. They read a line of terminal syntax, took it for something they had to type, and were reading the change's plumbing rather than the decision.
- Each decision was titled `f1`.
  That is what socrates calls the finding, and it is not what a person calls anything.
- What an option meant was written in the reviewer's vocabulary: "the window endpoints clamp with PAGE_LIMIT again, and the smoke assertions on the 25 cap come out".
  A person who does not already know the code cannot tell what they are choosing between.
- Beside the commit control sat a small grey box that did nothing when clicked.
  The line that says what was committed is hidden until there is something to say, and the stylesheet gives it a layout, which defeats the attribute that hides it.
- The palette is the palette of the tool that composed the document.
  A page whose colours belong to the machine tells the person they are reading that machine's output rather than their own decision.

The first four are the surface and the skill doing exactly what they were specified to do, so this change moves the specifications rather than working around them.
"Every control names the command it becomes" was written to stop a person deciding something whose destination they could not see.
The destination is now visible where it belongs, which is in the terminal the caller runs it in, and printing it on the page buys the person nothing but a harder second copy of the same choice.

## What Changes

- **No command is rendered.**
  An option still carries the exact command it becomes, `poll` still hands it back verbatim, and the page shows it to nobody.
  What a committed answer is shown back as is the option the person read when they chose it.
- **A decision is shown under a name, not an identifier.**
  The source a document cites may give each decision a name, copied onto the decision by the same mechanism that copies its severity, audience and location.
  The skill says to derive that name from the identifier socrates printed and never to invent a title, so `f1` reaches the person as "Finding 1" and the return to the terminal is still one exchange.
- **The bar the page is held to says what plain English is.**
  It already required plain English about what a finding means for the change. It now also requires that each option say what choosing it gets the person, in the words of the thing being built rather than the words of the code, and that no sentence be carried over from the reviewer untranslated.
- **An element the page hides paints nothing**, whatever else the stylesheet says about it.
- **The page is painted in the product's own palette**, stated once as tokens, and every text on it meets a stated contrast ratio against what it is painted on, under both preferences.

## Capabilities

### Modified Capabilities
- `surface-page`: the command stops being rendered and a committed answer is shown back in the words the person read; a hidden element paints nothing; the palette is stated once and held to a contrast floor.
- `surface-document`: a source entry may name a decision, and the name is copied rather than accepted like the identifier beside it.
- `findings-page`: the control rule keeps its reasoning and loses its display; a finding is shown under a name derived from its identifier; the bar says what plain English means for an option.

## Impact

- `internal/document`: a name on a source entry and on a decision, copied from the source.
- `internal/render`: the option's command and the committed command leave the template; the name is what a decision is headed with.
- `web/src`: the palette, the hidden rule, the committed line, and the styles the command carried.
- `internal/skill`: the page section, `PageBar`, and the regenerated `.claude/skills/socrates/SKILL.md`.
- `internal/convention`: the check that no colour is written outside the block that states the palette.
- `internal/e2e`: the page, presentation and handover journeys, and the browser the suite drives is told to keep its own credential store, because on one platform it stops for a dialog nobody is there to answer and the page then fails to load.
- **No change to** `internal/surface`, `internal/sanitise`, or anything socrates does. The exchange already carries the command back, and which rounds a person reads is unchanged.

## Dependencies

None. `a-control-names-the-command-it-becomes` and `a-page-that-answers-back` are archived, and both `surface-page` and `findings-page` are promoted.
