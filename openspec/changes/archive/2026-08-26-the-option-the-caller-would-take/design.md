## Context

See proposal.md for why.

The frame is already there.
A document carries decisions, each carrying options of label, meaning and command.
The page renders each option as a label wrapping a radio or a checkbox whose value is the option's position in the document, and the exchange reads those positions back, sorts them, and hands the caller each chosen option's own label and command.

So the document already has a place for one more field per option, and the page already renders options from a view assembled per request rather than from the document directly.
The one thing that does not exist is any notion that one option is different from the others.

Two properties constrain what can be done about it.
The position an option posts back is the key the exchange resolves against the document, so anything that renumbers options changes what `poll` returns.
And the page is held to a floor that requires every control to carry an accessible name, so a mark that only paints is a mark that half its readers do not receive.

## Goals / Non-Goals

**Goals:**

- A person meets the option the caller would take first, and learns that it is that option without having to interpret anything.
- The caller's contract is unchanged for a document that marks nothing, and unchanged in what it receives for one that marks something.
- Marking and choosing cannot be confused with one another anywhere: not in the document, not on the page, not in what comes back.

**Non-Goals:**

- Ranking the rest. One option is marked or none is, and the remaining options stay in the order the caller listed them.
- A reason attached to the mark. `meaning` already says what an option gets the person, and a second field for why this one is recommended is the same sentence in two places.
- Any reach from socrates into agora. The session writes the document, as it already does.
- Telling the caller whether the person took the recommendation. The caller wrote the document and gets the chosen option back, so it can already tell.

## Decisions

### The mark is a boolean on the option

An option carries `"recommended": true`, and no option carrying it is the ordinary case.

The alternative considered was a field on the decision naming which option is recommended, either by position or by label.
It was rejected because it is a second copy of the ordering: it can name a position that does not exist, it goes stale the moment an option is inserted above it, and validating it is work that a boolean makes impossible to need.
A mark that travels inside the thing it marks cannot come detached from it.

A rank per option was also considered, and rejected for the reason the bar states: a page that ranks everything has ranked nothing, and the person then reads six sentences to learn what one mark says.

### The order on the page is the renderer's, and the order in the document is the contract's

The hoist happens where the option views are assembled, and the document keeps the caller's order untouched.

Reordering in the document would renumber the positions, and the positions are what the page posts and what the exchange resolves.
The result is a contract where what an answer means depends on a field the caller set for presentation, which is exactly the coupling that makes a later change to placement a breaking change.

The view already carries `Position` separately from where the option sits in the slice, so this is one stable partition of the slice and nothing else.
The stability matters: the unmarked options keep their relative order, so a caller that ordered them deliberately still gets that order below the mark.

The single scenario worth stating in the specification is the one that catches the mistake: a decision that marks its last option, answered with that option, still hands the caller that option's own label and command.

### The mark is said, and paint is only emphasis

The word goes inside the option's label element, which is what the accessible name of its control is computed from, so the mark reaches a person hearing the page in the same breath as the option itself.
It carries a token background for the same reason the chosen option does, and removing every colour from the page leaves the mark fully intact.

The alternative considered was a coloured border or a heavier weight on the marked option.
It was rejected on the floor the page already states: a mark carried only by paint is missed by a screen reader, by a monochrome screen and by anyone who does not already know the convention, and what is left telling them is placement alone, which is the unstated convention this change exists to replace.

The word itself sits in the page template beside `Choose one` and the other words the page says for itself.
`messages.go` holds what a package's Go code emits, and the page's own copy has always been in the template it renders from.

### Marking never selects, and it is enforced by there being no mechanism

Nothing renders a `checked` attribute, so there is nothing to set and nothing to get wrong.
The scenario exists because the omission is what is load bearing: a pre-ticked option would come back through the exchange as an ordinary answer and be recorded in the terminal as the person's own.

A default that answers for the person is the same failure as an unclassified finding defaulting to `auto-fix`, which the reviewer half of this repository already refuses.

### Two marks are refused at validation rather than resolved at render

`serve` fails and names the decision, alongside the existing refusals for a missing field, an unknown mode and an option with no command.

The alternative was to take the first mark and drop the rest, which would serve a page that silently disagrees with the document that produced it.
A caller that marked two has not decided, and the failure is cheap and immediate: it happens before any server starts, in the session that wrote the document and can fix it.

The bound is one per decision in both modes.
Multiple-choice was considered as a case for several marks, since a person may tick several options.
It was rejected because the mark answers "what would you do", and a caller that would do three things says so in one option that is the three things, which is also what it has to hand back a single command for.

### The bar gains a line, beside the line about what an option gets the person

The two lines are adjacent and different.
The existing one is about the words an option is written in.
The new one is about the caller declining to state a judgement it has already made, which is a way a page fails a person that no existing line names.

Without it the field is a mechanism arriving with nothing calling for its use, and the sessions least likely to volunteer a recommendation are the ones whose pages most need one.
The line is bounded in the same sentence that asks for it, because "mark the one you would take" and "mark one and no other" are one instruction and drift apart if stated twice.

## Risks / Trade-offs

- **A session marks an option it should not have, and the mark carries authority the page cannot check.** → The page states it as the caller's own judgement rather than as a fact, every option keeps its meaning, and the free-text field that cannot be suppressed is untouched. The alternative is a person deriving the same judgement from nothing, which is the state today.
- **Placement and the mark say the same thing twice, and a caller could come to rely on placement alone.** → The specification requires the words and forbids conveying it by position alone, and the scenario that discards colour is what holds it.
- **The bar delta and `an-illustration-the-caller-paints` both modify the requirement that states the bar.** → They add adjacent lines to one list and neither edits the other's. Whichever archives second rebases its delta onto the list as it then stands, which is a one-line edit, and the enumerating scenario names both lines once both are in.
- **A document marking nothing loses nothing, so the field can be quietly ignored forever.** → That is the intended failure mode of a field nobody has to fill in, and the bar is the answer to it. Nothing asserts that a given page recommended well, because a test claiming a page recommended the right option would read as passing while checking nothing.

## Migration Plan

The field is additive and defaults to absent.
A document written before this change serves and renders exactly as it did, which is also the rollback: stop marking options.

The committed skill and the committed asset bundle are generated, so `make skill` and `make assets` run in the same change and the check that a generated file agrees with its source stays green.
