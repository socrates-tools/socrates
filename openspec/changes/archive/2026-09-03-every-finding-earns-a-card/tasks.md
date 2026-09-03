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

- [x] 1.1 Add a document test that a decision carrying no options validates, and verify it fails against the refusal in `internal/document/messages.go` for the right reason
- [x] 1.2 Remove the no-options refusal, and verify 1.1 passes and no other document test regressed
- [x] 1.3 Correct the test covering `A document that asks something is validated as before` so it exercises a decision whose options carry no command, and verify the refusal of a missing command is untouched
- [x] 1.4 Verify the two scenarios of `A document may hold nothing to decide` that this change does not touch still pass: a document with a preamble and entries and no decision serves, and a document with none of the three is refused
- [x] 1.5 Add a render test that a decision with no options emits its question, its meaning and no option control. Covers `A decision with no options renders as a card that asks nothing`
- [x] 1.6 Implement the no-options card in `internal/render`, and verify 1.5 passes
- [x] 1.7 Add a render test that the free-text field is emitted for a decision with no options and that no document field removed, labelled or resized it. Covers `A decision with no options still carries the field`
- [x] 1.8 Add an end-to-end test driving a real browser that a person writing in the field of an option-less decision and sending comes back through `poll` with no option, no command and the words verbatim. Covers `An annotation on a decision that asked nothing comes back`
- [x] 1.9 Add an end-to-end test that a document whose every decision carries no options renders every one as a card and offers no option control anywhere. Covers `A page of cards that ask nothing offers no controls`
- [x] 1.10 Verify the page-level behaviour of a document with no decisions at all is unchanged, so `A page with nothing to answer says so and offers neither control` still holds for the document shape it was written for
- [x] 1.11 Mutate the guard: make the no-options card emit an option control and verify 1.5 fails

## 2. A decision that asks nothing states no mode

- [x] 2.1 Add a document test that a decision with no options and a stated mode is refused and the decision is named, and verify it fails for the right reason. Covers `A mode on a decision with no options is refused before serving`
- [x] 2.2 Implement the refusal, with its message in `internal/document/messages.go`, and verify 2.1 passes
- [x] 2.3 Add a document test that a decision carrying options and stating no mode is still refused and named. Covers `A decision with options still requires its mode`
- [x] 2.4 Verify the em dash check and the user-facing string check pass over the new messages

## 3. The languages a caller actually passes

- [x] 3.1 Record the built bundle's size before the change, so the cost stated in `design.md` is measured rather than asserted
      `surface.js` 80,705 bytes, and 189,768 bytes for everything `internal/assets` embeds.
- [x] 3.2 Replace the hand-registered thirteen languages in `web/src/surface.js` with the bundled library's common set, run `make assets`, and commit the regenerated bundle
- [x] 3.3 Record the built bundle's size after the change, and state both figures in this task
      `surface.js` 80,705 bytes before and 169,214 after, and 189,768 bytes embedded before and 278,277 after.
      Thirteen languages by hand cost 80,705 and the library's common set of thirty-six costs 169,214, so the whole embedded bundle grows by 47 per cent.
- [x] 3.4 Add an end-to-end test that evidence in `csharp` and evidence in `diff` both render highlighted, driving a real browser and asserting on highlight tokens rather than on served markup. Covers `The languages a caller actually passes are highlighted`
- [x] 3.5 Verify `NOTICE` still accounts for every file `internal/assets` embeds, read from what is embedded rather than from `web/`, and verify the shipped bundle still carries the library's copyright notice
- [x] 3.6 Verify the bundle reaches no network, by loading the page with outbound network unavailable. Covers `Rendering needs no network`

## 4. An unknown language is refused, and an absent one is not

- [x] 4.1 State the registered language set in one place the Go side can read, generated from `web/` rather than written a second time, and add a check that a language on one side and not the other fails the suite in both directions
- [x] 4.2 Add a document test that an unknown evidence language is refused, naming the decision, the language and the set the surface knows, and verify it fails for the right reason. Covers `An unknown evidence language is refused before serving`
- [x] 4.3 Implement the refusal, and verify 4.2 passes
- [x] 4.4 Add a document test that evidence naming no language is served and renders as plain text. Covers `Evidence with no language is served`
- [x] 4.5 Mutate the guard: register one fewer language than the bundle carries and verify 4.1 fails in both directions

## 5. The cards arrive in the order a person would triage them

- [x] 5.1 Add a render test that decisions whose sources give error, warning and info render in that order whatever order the document listed them in, and verify it fails against document order. Covers `The most severe card is the first a person meets`
- [x] 5.2 Implement the severity order in `internal/render`, reading the severity the source gives rather than anything the decision states, and verify 5.1 passes
- [x] 5.3 Add a render test that the source entries are listed in the same severity order and still before the decisions. Covers `The source list is ordered the same way`
- [x] 5.4 Add a render test that equal severities hold the order the document listed them. Covers `Equal severities keep the order the document listed them`
- [x] 5.5 Add a render test that a document listing an informational decision first and an error last still renders the error first. Covers `The caller cannot order the page`
- [x] 5.6 Add a render test that an entry whose source states no severity is listed after every entry that carries one. Covers `An entry whose source states no severity sorts last`
- [x] 5.7 Review every existing test that asserts on the order entries or decisions render in, and correct the ones that assumed document order rather than deleting them
      `TestTheSourceIsReadBeforeTheDecisions` and `TestAnEntryTakesThePersonToTheDecision` are the two, and both cite a source already listed error, warning, info, so neither asserted an order the surface no longer renders and neither needed correcting.
- [x] 5.8 Verify an entry a decision covers still takes the person to that decision once both are reordered, so `An entry leads to the decision that covers it` still holds
- [x] 5.9 Verify option order within a decision is untouched: a marked option still renders first and the rest hold document order

## 6. What the skill tells a session to put on the page

- [x] 6.1 Change the page guidance in `internal/skill` from findings-in-full-and-a-list to one card per finding in the round, and add a skill prose test. Covers `Every finding in the round is a card`
- [x] 6.2 Add to the guidance the shape of a resolved finding's card: what the finding was, what resolved it, and the diff where that is what makes it clear. Covers `A resolved finding is stated rather than asked`
- [x] 6.3 State in the guidance that a resolved finding carries no options, and why. Covers `The reason for offering no options is given`
- [x] 6.4 Change the guidance for a round with nothing to answer so that it is every finding as a card carrying no options rather than a short list, and verify the skill no longer tells a session to write a list of its own. Covers `A round with nothing to answer is still a round in full`
- [x] 6.5 Add a skill prose test that the guidance names the severity order and the absence of controls as how the page separates what needs the person, and tells the session to shorten, list or omit nothing. Covers `Foregrounding survives without demoting anything`
- [x] 6.6 Verify the guidance does not restate the rule that the short list is copied rather than retyped, since `surface-document` and `surface-page` now own it
- [x] 6.7 Run `make skill`, commit the regenerated copy, and verify the suite fails while the generated copy and `internal/skill` disagree
- [x] 6.8 Verify the skill still names nothing the binary does not offer, still lists every state a run rests in, and still carries no em dash

## 7. The loop still runs

- [x] 7.1 Drive the whole loop end to end through the real binaries against a real repository, with a round holding an error, a warning, an informational finding and a finding the session resolved, and record what the page rendered and what `poll` returned
      `TestTheWholeRoundIsHandedOverAsCardsInSeverityOrder` drives it: `socrates review`, `socrates status`, `socrates findings --json`, `agora serve --document`, a real browser, `agora poll --session`.
      The reviewer returned f1 warning, f2 info, f3 error asking a person and f4 error the session resolved, listed in that order, and the document was written in that order too.

      What the page rendered:

          f3 severity=error   choices=2 diff=0 words=1
          f4 severity=error   choices=0 diff=1 words=1
          f1 severity=warning choices=0 diff=0 words=1
          f2 severity=info    choices=0 diff=0 words=1

      The glance listed Finding 3, Finding 4, Finding 1, Finding 2, in the same order.

      What `agora poll` returned, after the person chose on f3 and wrote on f4:

          {"decisions": [
             {"id": "f3", "chosen": ["the function stays"],
              "commands": ["socrates respond --finding f3 --decided --reason \"the endpoint stays, because the migration reads it\""],
              "words": ""},
             {"id": "f4", "chosen": null, "commands": null,
              "words": "the deletion was the bug I was chasing all morning"}],
           "message": "", "ended": true}
- [x] 7.2 Run `make fmt` and `make verify`, and record that every gate passed
      `make fmt` rewrote nothing. `make lint` reports 0 issues. `make test` passes every package including `internal/e2e` at 176s, `make scenario-drop`, `make corpus-outcomes` and `make fixtures-diff` all pass.
      `make scenario-coverage` names two scenarios and no others: `The decisions are in full and the rest is a list` and `A round with nothing to answer is a page with nothing to answer`, both of them scenarios of the requirement this change withdraws, still standing in `openspec/specs/findings-page/spec.md` until the change is archived.
      `make scenario-drop` reads the same delta as 2 requirements substituted and 12 added, with the promoted specifications losing nothing, so archiving is what closes this and no test is missing.
- [x] 7.3 Update `AGENTS.md` where the surface invariants it lists no longer match, in particular the line stating that every decision carries a free-text field, so the map and `openspec/` agree
