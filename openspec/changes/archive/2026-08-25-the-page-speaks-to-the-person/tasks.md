**The page half is driven by a real browser** against the real binary, as the rest of it already is. Asserting on served markup would test the renderer.

**The skill half asserts on the rendered document**, never on the constants it is assembled from.

**Nothing here asserts what a reviewer found.** The journeys require that a round arrived and that the handover happened.

**Each guard is mutated on purpose before it is trusted:** the test that proves the new behaviour is run against the old code first.

**`make lint test` after every task, not at the end of a group.**

## 1. The grey box beside the commit control

- [x] 1.1 Render the page with nothing committed, find every element the page hides that still occupies space, and confirm the check fails against the stylesheet as it stands. Covers `Nothing hidden is painted`
- [x] 1.2 State the rule once, above the palette, so that no later declaration can beat what hides an element
- [x] 1.3 Give a hidden element a layout again and confirm the check fails naming it

## 2. No command reaches the page

- [x] 2.1 Serve a document whose every option carries a command, and confirm no command text is anywhere on the rendered page. Covers `No command reaches the page`
- [x] 2.2 Take the command out of the template, the stylesheet and what the client reads
- [x] 2.3 Show a committed answer back in the labels the page offered. Covers `A committed answer is shown in the words the person read`
- [x] 2.4 Confirm a committed decision still renders differently from one merely selected in. Covers `A committed answer is distinguishable`
- [x] 2.5 Confirm `poll` still hands back each chosen option's command verbatim. Covers `The command is what comes back`
- [x] 2.6 Put a command back into the template and confirm 2.1 fails

## 3. A decision is shown under a name

- [x] 3.1 Serve a source naming one decision and not another, and confirm the first is shown under its name and the second under its identifier. Covers `A decision is shown under the name its source gives it`
- [x] 3.2 Carry a name on a source entry and copy it onto the decision, through the door the severity already uses
- [x] 3.3 Head each decision with the name in the renderer
- [x] 3.4 Write a different name on the decision itself and confirm the source wins. Covers `A caller cannot rename an identifier`
- [x] 3.5 Confirm a decision naming nothing in the source is still refused before serving. Covers `A decision naming nothing in the source is refused`

## 4. The palette

- [x] 4.1 Read the stylesheet and confirm every colour is a token stated in the block that states the palette. Covers `Every colour is stated once`
- [x] 4.2 Measure every text the page renders against the colour behind it, under both preferences, against the stated ratio. Covers `Every text meets the contrast floor`
- [x] 4.3 Repaint the page in the product's own palette, and confirm the floor, both preferences, the narrow case and the keyboard all still hold
- [x] 4.4 Dim one token below the floor and confirm the contrast check fails naming it
- [x] 4.5 Rebuild the bundle with `make assets` and confirm the suite fails while the committed bundle and its source disagree

## 5. The skill

- [x] 5.1 Add to the bar what an option has to say and in whose words. Covers `The bar names each way the page fails a person`
- [x] 5.2 Say that a control carries the command it becomes and that the page shows it to nobody. Covers `A control carries its command and shows it to nobody`
- [x] 5.3 Say to name each decision from the finding's identifier, and that an option says what choosing it means for the change rather than what to type. Covers `The page says what choosing an option means for the change`
- [x] 5.4 Confirm the bar is still stated once and rendered from there. Covers `The bar is stated once and rendered from there`
- [x] 5.5 Confirm the skill still says how an answer reaches the run and that the decision is the person's
- [x] 5.6 Regenerate `.claude/skills/socrates/SKILL.md` with `make skill`, and confirm the suite fails while the two differ

## 6. The loop, driven

- [x] 6.1 Serve a round through the real binary with the names the skill says to write, and confirm each finding is shown under its name with the severity and audience socrates reported. Covers `Every finding keeps the identifier socrates gave it`
- [x] 6.2 Confirm a decision answered on the page still reaches the run only through `socrates respond`. Covers `A decision made on the page is recorded in the terminal`
- [x] 6.3 Confirm the reason recorded is still the person's own. Covers `The person's own words reach the record`
- [x] 6.4 Confirm every journey left the repository under review byte-identical

## 7. Closing out

- [x] 7.1 Correct the invariants in `CLAUDE.md` that this change makes wrong
- [x] 7.2 Run the full sequence: `make fmt`, `make lint`, `make test`, `make scenario-coverage`, `make scenario-drop`, `make corpus-outcomes`
- [ ] 7.3 Serve a page from a real round and read it, because every defect this change fixes was found that way and none of them by a test
