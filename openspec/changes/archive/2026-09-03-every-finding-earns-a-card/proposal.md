## Why

A person drove a real review to a verdict this week and then asked for the page to be rebuilt three times.
Every one of the three requests was the page withholding something the person wanted, and none of them was the caller writing a bad document.

The round held six findings.
Three needed the person, two the session had already patched, one was informational.
The skill says the findings a person answers go on the page in full and the rest goes on as a short list, so the caller wrote three cards and put the other three into the preamble as prose.

The person's first request was to see all six.
A round is a unit of work and the person is signing off on the round, not on the subset that survived triage.
Prose in a preamble is where a reader's eye goes last, and three findings written there are three findings that were technically shown and actually skipped.

The second request was that the already-patched ones carry no options.
The caller had given each of them an option in order to render a card at all, because the surface refuses a decision that offers none, and the person read a control that changes nothing as a question they had to answer.
An acknowledgement control is worse than no control: it manufactures a decision that does not exist, and the one thing it reliably communicates is that something is still expected of them.

The third was that the code be highlighted.
`highlight.js` is bundled and the page does highlight, but the bundle registers thirteen languages chosen by hand, and the two the caller passed were `csharp` and `diff`.
Neither is registered, so both rendered as plain text, and nothing anywhere said so.
The caller had no way to learn that the language it named was one the page would ignore.

The fourth is the person's own and did not come from a rebuild: the cards came out in the order the caller listed them, which happened to be finding order.
Finding order is an accident of how the reviewer numbered its output.
The page is the one place in socrates where a person triages, and triage order is severity.

None of these is a defect in the caller.
The first two are the surface and the skill forcing a shape the round did not have, the third is a curated list whose failure mode is silence, and the fourth is the page declining to use a fact it already holds.

`a-page-for-every-round` landed while this was being written and closes part of the gap.
The surface now copies every source entry onto the page with its name, severity, audience and location, and opens the page with them, so the person does see that six things were found without the caller retyping anything.
It also allows a document with no decisions at all, on the reasoning that a caller with something to read and nothing to ask is the ordinary case.

That is the right shape one level up and it does not reach the level this is about.
A source entry carries what a finding is and where.
It cannot carry what was done about it or why, and that is what the person asked to see: not that a finding existed, but that it was a lock deleting itself while in use, and that the deletion is gone.
The same change kept the refusal of a decision with no options, and stated it in the specification for the first time, which makes the shape the person wants explicitly disallowed rather than merely unbuilt.

So this change is the same reasoning one step further in.
A caller with something for a person to read and nothing to ask them about that one finding is also the ordinary case.

## What Changes

- The skill SHALL describe a page that carries every finding in the round as a card of its own, rather than the answerable ones in full and the rest as a list.
- A decision MAY carry no options, and the surface SHALL render it as a card that states something rather than asking anything.
- A decision carrying no options SHALL still carry the free-text field, because a person reading what was already done is exactly a person who may want to say something about it.
- A decision carrying no options SHALL NOT state a mode, and a document that gives one a mode SHALL be refused.
- The skill SHALL say that a finding the session resolved is put on the page as what the finding was and what resolved it, with the diff where the diff is what makes it clear.
- The source list the surface already renders SHALL stay what it is, an index of the whole, and SHALL NOT be what stands in for a finding the person never sees in full.
- The page SHALL order both the source entries and the cards by severity, most severe first and informational last, and SHALL hold equal severities in the order the document listed them.
- The bundle SHALL register the languages a caller actually passes rather than a list somebody maintains by hand.
- A document naming a language the page cannot highlight SHALL be refused before serving, naming the language and the decision, so that a caller learns at the terminal rather than the person learning from an unhighlighted block.
- Evidence naming no language SHALL keep rendering unhighlighted, because that is a caller saying it wants none rather than a caller guessing wrong.

## Capabilities

### New Capabilities

None.
Every part of this change is an existing capability being asked for something it already nearly does.

### Modified Capabilities

- `findings-page`: what the skill tells a session to put on the page, and in what shape.
  The foregrounding requirement is withdrawn and replaced, because its remedy was the short list this change stops relying on.
- `surface-document`: what a caller may write, extended to a decision that asks nothing.
  `A document may hold nothing to decide` is amended, because it is where the refusal of an option-less decision is now stated.
- `surface-page`: the order the entries and the cards come out in, and which languages the page can highlight.

## Impact

- `internal/document`: the refusal of a decision with no options is removed, a refusal of a mode without options is added, and a refusal of an unknown evidence language is added.
- `internal/render`: cards are emitted in severity order, and a card with no options renders its statement and its free-text field and no controls.
- `internal/skill`: the page guidance changes from full-and-a-list to a card per finding, and gains the shape of a resolved finding's card.
- `web/src/surface.js`: the registered language set, and the one place a language name is resolved.
- `NOTICE`: the bundled library set is unchanged, so nothing is added, but the assertion runs against the newly built bundle.
- `internal/e2e`: the page half gains scenarios for a card with no options, for severity order, and for a language the surface refuses.
- Nothing reaches the reviewer, the snapshot, the run, the dialogue or the verdict.
  socrates still states that a round is one a person reads and still does nothing further about the page.

## Dependencies

This change SHALL follow `a-page-for-every-round`, which is archived.
That change is what put the source entries on the page and what stated the refusal this change amends, so both of its halves are load-bearing here.
`findings-page`, `surface-document` and `surface-page` are otherwise promoted, and this change asks each of them for something narrower than what it already does.
