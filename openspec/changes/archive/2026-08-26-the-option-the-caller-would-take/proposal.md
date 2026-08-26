## Why

A decision arrives on the page as a flat list of options, and every option looks equally likely.
The session that wrote the document has just read the findings, and in almost every case it knows which option it would take.
It has nowhere to say so.

The person deciding is the person who built the thing, and what they want first is not a menu but a proposal they can disagree with.
A flat list makes them re-derive a judgement that was already made one step upstream, which is the same failure the bar already names: handing them the terminal output in a nicer typeface.

The absence is not neutral either.
A session that wants to recommend something and has no field for it puts the recommendation in the prose, or leans on option order and hopes it reads as ranking.
Both are unstated conventions the page does not honour, so what the person sees is a hint the surface never agreed to carry.

## What Changes

- An option may be marked as the one the caller would take.
  At most one per decision, in either mode, and a document marking two is refused before any server starts, because a recommendation the person has to assemble is not one.
- The page says so in words, inside the option, so the mark survives a screen reader, a monochrome screen and a person who does not know what the colour means.
  Nothing about it is carried by colour alone.
- The page renders the marked option first and leaves the rest in the order the document gave them.
  A recommendation below the fold is a recommendation the person reads after they have already decided.
- **Marking is not choosing.** Nothing on the page is selected until the person selects it, and a decision they never touched still comes back with nothing chosen and no command.
  A pre-ticked option is an answer nobody gave, and it would be indistinguishable from theirs by the time it reached the terminal.
- The value an option posts back stays the position the document gave it, so what `poll` hands back is unchanged by the reordering, and the caller's own ordering is what it always was.
- The bar gains a line: say which option you would take.
  Without it the field is a mechanism arriving with nothing that calls for its use, and the sessions that most need to be told are the ones that would leave it out.
- The skill teaches the field where it teaches the document's shape, and says that the mark ranks the option rather than answering for the person.
- **socrates reaches nothing new.** No command, no flag, no run record. The session writes the document, as it already does.

## Capabilities

### New Capabilities

None.
The recommendation is a field on something a caller already writes and a mark on something the page already renders.

### Modified Capabilities

- `surface-document`: a requirement for the mark an option may carry, the one-per-decision bound and the rule that marking never selects, beside the existing requirement that every option carries the command it becomes.
- `surface-page`: `A decision is answered by choosing, and the mode says how many` gains the page's placement of the marked option and the requirement that nothing arrives selected; `The presentation is held to a stated floor` gains that the mark is stated in words rather than in colour.
- `findings-page`: `The page is held to a stated bar, stated once` gains the line about saying which option you would take, and the scenario that enumerates the bar names it.
- `agent-skill`: the skill teaches the field, what marking does and does not do, and that at most one option per decision carries it.

## Impact

- `internal/document`: `Option` gains the field, and validation refuses a decision marking more than one.
- `internal/render`: the marked option is hoisted to the front of the option views while keeping the position the document gave it, and the page states the mark in the option's own text.
- `internal/surface`: nothing changes in the exchange; a test holds that what comes back is still keyed by document position and that an untouched decision carries nothing.
- `web/src/surface.css`: the mark is painted from existing tokens; `make assets` regenerates the committed bundle.
- `internal/skill`: the passage beside the document shape, and one line of `PageBar`; `make skill` regenerates the committed copy.
- `internal/e2e`: a real browser reads the mark, finds the recommended option first, and finds nothing selected until it selects something.
