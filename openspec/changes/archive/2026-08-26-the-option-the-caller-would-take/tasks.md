**Test first.** A task is done when a test named for its scenario failed for the right reason and then passed.

**The page half is driven by a real browser against the real binary.** Asserting on served markup would test the renderer rather than what a person gets.

**Each guard is mutated on purpose before it is trusted.** A document marking two options is served by hand first, and the test that refuses it must fail before it is believed.

**Nothing here asserts that a session recommended well.** Which option a session marks is judgement. What is driven is the instruction the skill carries and the mechanism the surface offers.

**`make lint test` after every task, not at the end of a group.**

## 1. The document carries the mark

- [x] 1.1 Add the mark to an option, absent by default, and confirm a document that carries it round-trips through `Read` unchanged. Covers `A decision names the option the caller would take`
- [x] 1.2 Refuse a decision marking more than one option, naming the decision, alongside the fields validation already refuses. Covers `A document marking two options is refused before serving`
- [x] 1.3 Confirm a document whose decisions mark nothing validates, serves and renders its options in the order it listed them. Covers `A document that marks nothing is unaffected`
- [x] 1.4 Mutate 1.2 by taking the first mark and dropping the rest, and confirm its test fails rather than a later render test catching it

## 2. The page reads the marked option first

- [x] 2.1 Hoist the marked option to the front of a decision's option views, keeping the unmarked ones in the order the document listed them and keeping every option's position. Covers `The option the caller would take is read first`
- [x] 2.2 Confirm an answer to a decision marking its last option hands the caller that option's own label and command. Covers `Rearranging the options does not rearrange the answer`
- [x] 2.3 Confirm a multiple-choice decision that marks one option still carries every tick in the order the document listed them. Covers `Multiple choice carries every tick`
- [x] 2.4 Confirm nothing on a rendered decision arrives selected, and that a decision the person never touched comes back with nothing chosen and no command. Covers `Nothing arrives chosen` and `Marking is not choosing`
- [x] 2.5 Mutate 2.1 to renumber the positions after hoisting, and confirm 2.2's test fails rather than 2.1's passing on the reordering alone

## 3. The page says the mark rather than painting it

- [x] 3.1 Say inside the marked option that it is the one the caller would take, in the option's own label, so the words are part of the name its control is announced under. Covers `The mark is said rather than painted`
- [x] 3.2 Confirm in a real browser with every colour discarded that the marked option still says so, and that a person moving by keyboard alone hears it with the option. Covers `The mark is said rather than painted` and `A decision can be answered without a pointer`
- [x] 3.3 Paint the mark from the tokens the page already states, and regenerate the bundle with `make assets`. Covers `Every colour is stated once`
- [x] 3.4 Confirm the mark's text meets the contrast floor under both preferences and renders no smaller than the stated minimum. Covers `Every text meets the contrast floor` and `Nothing renders below the floor`
- [x] 3.5 Mutate 3.1 by carrying the mark in paint alone, and confirm 3.2's test fails rather than the contrast sweep catching it

## 4. The bar, and the skill beside it

- [x] 4.1 Add the bar's line about saying which option you would take, bounded to one option in the same sentence, in the one place the bar is stated. Covers `The bar names each way the page fails a person`
- [x] 4.2 Confirm the line appears once, that the line about what an option gets the person is unchanged, and that the committed skill matches its source after `make skill`. Covers `The bar is stated once and rendered from there`
- [x] 4.3 Mark an option in the document shape the skill shows, so a session copying it copies the field. Covers `The skill names the mark where it names the document`
- [x] 4.4 Say beside it that at most one option per decision carries the mark and that a document marking two is refused before serving. Covers `The skill says one option carries it`
- [x] 4.5 Say that the page reads the marked option first and states it in words, that nothing arrives selected, and that the decision is still the person's. Covers `The skill says marking is not answering`
- [x] 4.6 Confirm the judgement about when to recommend appears in the bar and nowhere else in the skill. Covers `When to recommend is stated once`

## 5. The loop through the real binary

- [x] 5.1 Serve a document marking an option, and confirm through a real browser that the marked option is the first a person meets and that the rest follow in the document's order. Covers `The option the caller would take is read first`
- [x] 5.2 Answer that decision with the marked option and confirm `poll` prints its command verbatim. Covers `Rearranging the options does not rearrange the answer` and `The command is what comes back`
- [x] 5.3 Send without touching a decision that marks an option, and confirm `poll` returns it with nothing chosen and no command. Covers `Marking is not choosing`
- [x] 5.4 Confirm `serve` fails and names the decision for a document marking two options, and that no server was started. Covers `A document marking two options is refused before serving`
- [x] 5.5 Confirm a document written before the mark existed serves and renders as it did. Covers `A document that marks nothing is unaffected`

## 6. The gates

- [x] 6.1 Run `make fmt`, `make lint` and `make test`
- [x] 6.2 Run `make scenario-coverage` and `make scenario-drop`, and confirm every scenario this change adds is named by a test
- [x] 6.3 Run `make corpus-outcomes` and confirm the recorded distribution is unmoved
