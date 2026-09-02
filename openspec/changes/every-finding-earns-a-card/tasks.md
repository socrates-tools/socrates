**The suite must be green after every task, not only at the end of a group.**
The ordering below exists for that reason.
The language set is registered before anything refuses an unknown language, because a refusal written first would reject `csharp` and `diff` while they are still unregistered and fail the suite on its own fixtures.

**Nothing in this change reaches socrates.**
No command, flag or run record may come to name a page, and no task here adds one.
A task that moves any part of the page into `cmd/socrates` or `internal/reviewer` is in the wrong change.

**The surface learns no new nouns.**
`internal/surface`, `internal/document`, `internal/render` and `web/` may not name review, finding, verdict, round or run.
The convention check already holds this, and every task below is written in the surface's own nouns for that reason.

## 1. A decision may state something rather than ask

- [ ] 1.1 Add a document test that a decision carrying no options validates, and verify it fails against the refusal in `internal/document/messages.go` for the right reason
- [ ] 1.2 Remove the no-options refusal, and verify 1.1 passes and no other document test regressed
- [ ] 1.3 Correct the test covering `A document that asks something is validated as before` so it exercises a decision whose options carry no command, and verify the refusal of a missing command is untouched
- [ ] 1.4 Verify the two scenarios of `A document may hold nothing to decide` that this change does not touch still pass: a document with a preamble and entries and no decision serves, and a document with none of the three is refused
- [ ] 1.5 Add a render test that a decision with no options emits its question, its meaning and no option control. Covers `A decision with no options renders as a card that asks nothing`
- [ ] 1.6 Implement the no-options card in `internal/render`, and verify 1.5 passes
- [ ] 1.7 Add a render test that the free-text field is emitted for a decision with no options and that no document field removed, labelled or resized it. Covers `A decision with no options still carries the field`
- [ ] 1.8 Add an end-to-end test driving a real browser that a person writing in the field of an option-less decision and sending comes back through `poll` with no option, no command and the words verbatim. Covers `An annotation on a decision that asked nothing comes back`
- [ ] 1.9 Add an end-to-end test that a document whose every decision carries no options renders every one as a card and offers no option control anywhere. Covers `A page of cards that ask nothing offers no controls`
- [ ] 1.10 Verify the page-level behaviour of a document with no decisions at all is unchanged, so `A page with nothing to answer says so and offers neither control` still holds for the document shape it was written for
- [ ] 1.11 Mutate the guard: make the no-options card emit an option control and verify 1.5 fails

## 2. A decision that asks nothing states no mode

- [ ] 2.1 Add a document test that a decision with no options and a stated mode is refused and the decision is named, and verify it fails for the right reason. Covers `A mode on a decision with no options is refused before serving`
- [ ] 2.2 Implement the refusal, with its message in `internal/document/messages.go`, and verify 2.1 passes
- [ ] 2.3 Add a document test that a decision carrying options and stating no mode is still refused and named. Covers `A decision with options still requires its mode`
- [ ] 2.4 Verify the em dash check and the user-facing string check pass over the new messages

## 3. The languages a caller actually passes

- [ ] 3.1 Record the built bundle's size before the change, so the cost stated in `design.md` is measured rather than asserted
- [ ] 3.2 Replace the hand-registered thirteen languages in `web/src/surface.js` with the bundled library's common set, run `make assets`, and commit the regenerated bundle
- [ ] 3.3 Record the built bundle's size after the change, and state both figures in this task
- [ ] 3.4 Add an end-to-end test that evidence in `csharp` and evidence in `diff` both render highlighted, driving a real browser and asserting on highlight tokens rather than on served markup. Covers `The languages a caller actually passes are highlighted`
- [ ] 3.5 Verify `NOTICE` still accounts for every file `internal/assets` embeds, read from what is embedded rather than from `web/`, and verify the shipped bundle still carries the library's copyright notice
- [ ] 3.6 Verify the bundle reaches no network, by loading the page with outbound network unavailable. Covers `Rendering needs no network`

## 4. An unknown language is refused, and an absent one is not

- [ ] 4.1 State the registered language set in one place the Go side can read, generated from `web/` rather than written a second time, and add a check that a language on one side and not the other fails the suite in both directions
- [ ] 4.2 Add a document test that an unknown evidence language is refused, naming the decision, the language and the set the surface knows, and verify it fails for the right reason. Covers `An unknown evidence language is refused before serving`
- [ ] 4.3 Implement the refusal, and verify 4.2 passes
- [ ] 4.4 Add a document test that evidence naming no language is served and renders as plain text. Covers `Evidence with no language is served`
- [ ] 4.5 Mutate the guard: register one fewer language than the bundle carries and verify 4.1 fails in both directions

## 5. The cards arrive in the order a person would triage them

- [ ] 5.1 Add a render test that decisions whose sources give error, warning and info render in that order whatever order the document listed them in, and verify it fails against document order. Covers `The most severe card is the first a person meets`
- [ ] 5.2 Implement the severity order in `internal/render`, reading the severity the source gives rather than anything the decision states, and verify 5.1 passes
- [ ] 5.3 Add a render test that the source entries are listed in the same severity order and still before the decisions. Covers `The source list is ordered the same way`
- [ ] 5.4 Add a render test that equal severities hold the order the document listed them. Covers `Equal severities keep the order the document listed them`
- [ ] 5.5 Add a render test that a document listing an informational decision first and an error last still renders the error first. Covers `The caller cannot order the page`
- [ ] 5.6 Add a render test that an entry whose source states no severity is listed after every entry that carries one. Covers `An entry whose source states no severity sorts last`
- [ ] 5.7 Review every existing test that asserts on the order entries or decisions render in, and correct the ones that assumed document order rather than deleting them
- [ ] 5.8 Verify an entry a decision covers still takes the person to that decision once both are reordered, so `An entry leads to the decision that covers it` still holds
- [ ] 5.9 Verify option order within a decision is untouched: a marked option still renders first and the rest hold document order

## 6. What the skill tells a session to put on the page

- [ ] 6.1 Change the page guidance in `internal/skill` from findings-in-full-and-a-list to one card per finding in the round, and add a skill prose test. Covers `Every finding in the round is a card`
- [ ] 6.2 Add to the guidance the shape of a resolved finding's card: what the finding was, what resolved it, and the diff where that is what makes it clear. Covers `A resolved finding is stated rather than asked`
- [ ] 6.3 State in the guidance that a resolved finding carries no options, and why. Covers `The reason for offering no options is given`
- [ ] 6.4 Change the guidance for a round with nothing to answer so that it is every finding as a card carrying no options rather than a short list, and verify the skill no longer tells a session to write a list of its own. Covers `A round with nothing to answer is still a round in full`
- [ ] 6.5 Add a skill prose test that the guidance names the severity order and the absence of controls as how the page separates what needs the person, and tells the session to shorten, list or omit nothing. Covers `Foregrounding survives without demoting anything`
- [ ] 6.6 Verify the guidance does not restate the rule that the short list is copied rather than retyped, since `surface-document` and `surface-page` now own it
- [ ] 6.7 Run `make skill`, commit the regenerated copy, and verify the suite fails while the generated copy and `internal/skill` disagree
- [ ] 6.8 Verify the skill still names nothing the binary does not offer, still lists every state a run rests in, and still carries no em dash

## 7. The loop still runs

- [ ] 7.1 Drive the whole loop end to end through the real binaries against a real repository, with a round holding an error, a warning, an informational finding and a finding the session resolved, and record what the page rendered and what `poll` returned
- [ ] 7.2 Run `make fmt` and `make verify`, and record that every gate passed
- [ ] 7.3 Update `AGENTS.md` where the surface invariants it lists no longer match, in particular the line stating that every decision carries a free-text field, so the map and `openspec/` agree
