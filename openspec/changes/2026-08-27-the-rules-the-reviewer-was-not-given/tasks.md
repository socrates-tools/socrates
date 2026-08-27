**Every assertion is on the delivered prompt**, read from the invocation the simulator recorded, never from socrates' source.
`internal/e2e/prompt_test.go` already reads it that way and is where these join.

**No test asserts what a model concluded.** Whether these rules find more is answered by the corpus and by `make live-smoke`, and neither is a task here.

**Each guard is mutated on purpose before it is trusted:** the assertion that a rule reached the reviewer is run against the prompt without that rule first, and is expected to fail there.
A rule added and asserted in one step is a rule nobody has seen fail.

**`make lint test` after every task, not at the end of a group.**

## 1. The reviewer is told to report a test that breaks the rule

- [x] 1.1 Assert the delivered prompt carries, beside the test-quality rule, the instruction to report a test the change adds that breaks it, and confirm the assertion fails against the prompt as it stands. Establishes the gap against `The reviewer is told to report what the rule forbids`
- [x] 1.2 Add the reviewer's instruction to `internal/testquality` as a second exported block, leaving `Rule` byte-identical
- [x] 1.3 Deliver it under the reviewer prompt's test-quality section, after the rule, and confirm 1.1 passes
- [x] 1.4 Assert the instruction names a test the change touched as in scope and a repository-wide cleanup as out of it. Covers `The reporting duty is bounded to the change`
- [x] 1.5 Confirm `TestOneBlockManyRoles` still passes and that every role receiving the shared rule receives the same text, with the reviewer alone receiving the second block. Covers `One block, many roles`
- [x] 1.6 Take the second block back out of the reviewer's section and confirm 1.1 and 1.4 both fail while 1.5 still passes

## 2. The delivered test-quality rule and the requirement's copy agree

- [x] 2.1 Assert the requirement's quoted block and `testquality.Rule` carry the same three paragraphs, including that a bug-fix change's test reproduces the reported failure
- [x] 2.2 Confirm the assertion passes without changing `Rule`, since the paragraph has been delivered all along and the spec is what this change corrected

## 3. A stated intent is checked against the change

- [x] 3.1 Assert the intent section states that a change omitting a behaviour the intent names, or carrying one it rules out, is reported naming both, and confirm the assertion fails against the prompt as it stands. Establishes the gap against `A contradiction is reported and left to a person`
- [x] 3.2 Add the obligation to the intent section, so it is built only where an intent was stated
- [x] 3.3 Assert the obligation says the contradiction is reported even where nothing else about the change is wrong. Covers `An otherwise clean change is not an exemption`
- [x] 3.4 Assert the obligation says the reviewer does not resolve it alone, and that what it asks for is the action the classifier reads as needing a person. Covers `A contradiction is reported and left to a person`
- [x] 3.5 Assert the obligation says what the intent does not mention is not something it ruled out. Covers `An unstated thing is not a contradiction`
- [x] 3.6 Run a review with no stated intent and confirm the delivered prompt carries neither the intent, the obligation, nor a heading for either. Covers `No intent, no obligation`
- [x] 3.7 Confirm the intent still reaches the reviewer marked as what the author says they wanted rather than as fact. Covers `Intent reaches the reviewer`
- [x] 3.8 Take the obligation out and confirm 3.1, 3.3, 3.4 and 3.5 all fail while 3.6 and 3.7 still pass

## 4. A bug fix is asked whether the bug is still reachable

- [x] 4.1 Assert the reviewing section asks the reviewer to establish whether a bug-fix change claims a durable fix or a stated containment, and confirm the assertion fails against the prompt as it stands. Establishes the gap against `The reviewer is asked which claim the fix makes`
- [x] 4.2 Add the rule to the reviewing section, placed after the anti-astronaut clamp
- [x] 4.3 Assert the rule asks for the failing sequence, the invariant, the sibling paths and the shared state, and asks whether the same failure is still reachable. Covers `A durable fix is traced past the symptom`
- [x] 4.4 Assert the rule requires a path the source shows and says a containment the author stated as one is not a defect for being one. Covers `The evidence bar is the one the rest of the prompt sets`
- [x] 4.5 Confirm the anti-astronaut clamp is still delivered, still ahead of the new rule, and still carries the same text. Covers `Every rule reaches the reviewer`
- [x] 4.6 Take the rule out and confirm 4.1, 4.3 and 4.4 all fail while 4.5 still passes

## 5. Nothing else moved

- [x] 5.1 Add the new rule to the table in `TestEveryRuleReachesTheReviewer` so one test names every rule the requirement lists. Covers `Every rule reaches the reviewer`
- [x] 5.2 Confirm the sections are still delivered in one fixed order, that two reviews of the same change are still byte-identical, and that a section with nothing to say is still absent. Covers `Order is stable`, `An empty section is absent`
- [x] 5.3 Confirm the action definitions still match the classification the code applies, and that no finding schema, classification, verdict, dialogue or command behaviour changed
- [x] 5.4 Run `make scenario-coverage` and `make scenario-drop` and confirm every scenario added here is named by a test and nothing was dropped from the promoted specifications
- [x] 5.5 Run `make corpus-outcomes` and record the distribution it prints, so the shift these rules cause is measured from a stated starting point rather than remembered

## 6. Verification

- [x] 6.1 Run `make fmt`, `make lint`, `make test`, `make scenario-coverage`, `make scenario-drop` and `make corpus-outcomes`
- [x] 6.2 Run `make live-smoke` once against the real CLI and confirm the round returns a schema-conforming document under the longer prompt
