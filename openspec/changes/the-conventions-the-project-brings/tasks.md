**Every assertion about what the reviewer was told is on the delivered prompt**, read from the invocation the simulator recorded, never from socrates' source.

**No test proves the anchor by reading Go.** Whether the guide came from the base commit is proved by a repository whose base and working tree say different things, and by which of the two reached the prompt.

**No test asserts what a model concluded** from the guidance. That the guidance decides nothing is asserted on the verdict, the classification, what a person is asked and the exit code, with the reviewer's answer held fixed.

**Each guard is mutated on purpose before it is trusted.** The anchor is proved by pointing the read at the working tree and watching the anchor tests fail; the ceiling by raising it; the additive framing by taking the sentence out.
A guard added and asserted in one step is a guard nobody has seen fail.

**`make lint test` after every task, not at the end of a group.**

## 1. Reading a file at a commit

- [x] 1.1 Export reading a file's contents at a commit from `internal/snapshot`, so that git invocation stays in the one package that makes it, and verify it against a temporary repository whose working tree and base disagree about the same path
- [x] 1.2 Confirm the read reports absence distinctly from failure, so that a repository with no such path is not a git error, and verify with a path the commit does not carry
- [x] 1.3 Confirm `git status` in the repository is byte-identical before and after the read, and that nothing was written. Covers `Reading the guide writes nothing`

## 2. `internal/guidance`

- [x] 2.1 Add `internal/guidance` holding the fixed path `.socrates/review-guide.md` and the 16 KiB default ceiling, and verify each is stated in exactly one place
- [x] 2.2 Read the guide at a given base commit and answer one of four outcomes: carried, none, present in the change and not at the base, or over the ceiling. Verify each against a temporary repository built for it
- [x] 2.3 Return whitespace-only content as no guidance rather than as guidance holding nothing, and verify with a guide holding only blank lines. Covers `A guide holding nothing`
- [x] 2.4 Carry no part of a guide over the ceiling, answer its size and the ceiling it was held to, and let the review run. Covers `A guide over the ceiling is not carried`
- [x] 2.5 Carry a guide of exactly the ceiling whole, and verify the boundary from both sides. Covers `A guide at the ceiling is carried whole`
- [x] 2.6 Read the size from the base commit's tree entry before any of the guide is read, so an oversized guide is never loaded, and verify with a guide that crosses the ceiling only before neutralisation would collapse it
- [x] 2.7 Put every string the package emits in its own `messages.go`, and confirm none carries an em dash

## 3. The section the reviewer reads

- [x] 3.1 Assert the delivered prompt carries a guidance section holding what the base commit's guide says, and confirm the assertion fails against the prompt as it stands. Establishes the gap against `A project's guidance reaches the reviewer` and `The guidance is its own section`
- [x] 3.2 Add the section to `internal/prompt` in the fixed order, after the test-quality block and before the reporting contract, with socrates' own framing around the guidance, and confirm 3.1 passes
- [x] 3.3 Assert the framing states that the project wrote the guidance, that it adds to the rules above and withdraws none of them, and that it decides nothing about severity, action or what is reported. Covers `The section says what the guidance cannot do` and `The section is framed as the project's, not as socrates'`
- [x] 3.4 Assert the delivered headings put the guidance section after the reviewing rules and the test-quality block and before the reporting contract. Covers `The guidance comes after the rules it cannot remove` and `The order puts socrates' rules first`
- [x] 3.5 Run a review against a repository with no guide and confirm the delivered prompt carries neither the section nor a heading for it. Covers `A repository with no guide` and `No guidance, no heading`
- [x] 3.6 Confirm two reviews of the same change with the same guidance deliver byte-identical prompts. Covers `The prompt stays byte-identical across runs`
- [x] 3.7 Confirm the delivered prompt names no commit for the guidance, since it is not something the reviewer decides anything from
- [x] 3.8 Take the framing sentence out and confirm 3.3 fails while 3.1, 3.4 and 3.5 still pass

## 4. The base commit is what is read

- [x] 4.1 Build a repository whose base carries no guide and whose working tree adds one, run a review, and assert the delivered prompt carries no guidance section. Covers `A guide the change adds does not govern that change`
- [x] 4.2 Build a repository whose base carries one guide and whose working tree says something else in the same file, run a review, and assert the delivered prompt carries the base commit's words and none of the change's. Covers `A guide the change edits is delivered as the base had it`
- [x] 4.3 Change the guide between round one and round two of a run and assert both delivered prompts carry the same guidance, since every round is snapshotted against the run's original base. Covers `The guidance does not move between rounds`
- [x] 4.4 Point the read at the snapshot worktree instead of the base commit and confirm 4.1, 4.2 and 4.3 all fail, then put it back

## 5. The guidance decides nothing

- [x] 5.1 Run two reviews on the same change with the same recorded reviewer answer, one against a guide demanding a severity, demanding an action and asking that a finding be withheld, and one against no guide, and assert the verdict, every finding's classification, what a person is asked and the exit code are the same in both. Covers `The guidance decides nothing socrates decides`
- [x] 5.2 Confirm no code outside `internal/prompt`'s section reads the guidance, by removing the section and confirming that every test except the prompt and recording ones still passes

## 6. Neutralisation

- [x] 6.1 Assert a guide whose lines would open a prompt section delivers the same set of sections socrates assembled, with the guide's text inside the guidance section. Covers `A guide cannot open a section`
- [x] 6.2 Assert a guide carrying conflict markers is delivered with them removed, on `internal/hygiene` and with no second mechanism added. Covers `A guide carrying conflict markers`
- [x] 6.3 Confirm the guide on disk is unchanged after the round, since socrates never writes in the repository under review
- [x] 6.4 Take the neutralisation out and confirm 6.1 and 6.2 both fail

## 7. What the run records, and what a person reads

- [x] 7.1 Record the outcome in the run and assert a round that carried guidance says so. Covers `Guidance carried is recorded`
- [x] 7.2 Assert a repository with no guide at the base and none in the change records that there was no guidance, distinctly from a guide that did not reach the reviewer. Covers `No guide is recorded as no guide`
- [x] 7.3 Assert a repository whose change adds a guide the base does not carry records that a guide is present in the change and did not reach this reviewer, and that the round is not failed for it. Covers `A newly written guide is reported to the person`
- [x] 7.4 Give the outcome a line in `internal/presentation`, so the person reading the run is told the first time they write the file, and assert the line names the file and says the guide governs the next review. Covers `A newly written guide is reported to the person`
- [x] 7.5 Confirm every string added here lives in its package's `messages.go`, names what the person can do, and carries no em dash
- [x] 7.6 Collapse the third outcome into the second and confirm 7.2, 7.3 and 7.4 all fail

## 8. Nothing configures it

- [x] 8.1 Assert a configuration file naming a guidance path is refused, naming the field socrates does not have, on the unknown-field refusal `internal/config` already applies. Covers `Nothing configures the path`
- [x] 8.2 Add `guidance_ceiling` to the user's configuration file, assert a stated ceiling is what the guide is held to while a repository stating one reaches nothing, and assert a negative one is refused with the rest of the configuration. Covers `The ceiling is the person's to state, not the repository's`, `An unstated ceiling is the default`
- [x] 8.3 Confirm the diff ceiling still applies to every repository and every run with no per-repository override, and that no repository file reaches any other setting

## 9. The seam, restated rather than reopened

- [x] 9.1 Reword the `agent-seam` scenario's outcome to name the repository's *agent* instructions, and confirm the existing test still passes against the invocation unchanged. Covers `Project settings are not loaded`
- [x] 9.2 Add a review against a repository carrying both its own agent instructions and a guide, and assert the invocation still restricts settings to user level, the prompt carries none of the agent instructions, and the guidance reaches its own section. Covers `The one file socrates reads is not the vendor loading the repository`
- [x] 9.3 Confirm nothing in the invocation changed: the same settings restriction, the same customisation mechanism, the same tool set, and no settings of socrates' own

## 10. The documents that describe it

- [x] 10.1 Remove the deferred roadmap note about repository-scoped review guidance from `openspec/config.yaml`, since it is no longer deferred and it named an answer this change replaced
- [x] 10.2 Add the invariants to `CLAUDE.md` beside the ones about the repository under review: one file is read from it, at the base commit, additively, and nothing is written there
- [x] 10.3 Name `internal/guidance` in `CLAUDE.md`'s layout, and confirm every package the file lists still exists
- [x] 10.4 State in `README.md` what a project writes in `.socrates/review-guide.md`, that it takes effect once committed, and that it adds to socrates' rules rather than replacing them

## 11. socrates' own guide

- [x] 11.1 Write `.socrates/review-guide.md` for this repository, carrying only what no deterministic gate here catches, and confirm it restates nothing `internal/convention`, `make lint` or the coverage gates already enforce
- [ ] 11.2 Run `socrates review` on this change with that guide at the base and record what the reviewer did with it, so the channel is judged on a real round rather than on intent

## 12. Verification

- [x] 12.1 Run `make fmt` and `make verify`
- [x] 12.2 Confirm `make scenario-coverage` names a test for every scenario added here and `make scenario-drop` reports nothing dropped from the promoted specifications
- [x] 12.3 Run `make corpus-outcomes` and record the distribution, so a shift the longer prompt causes is measured from a stated starting point
- [ ] 12.4 Run `make live-smoke` once against the real CLI with a guide in place, and record that the round returned a schema-conforming document

## 13. What the reviewer found

- [x] 13.1 Neutralise once rather than twice, by making `hygiene.Neutralise` idempotent and asserting it, so the text a run records and the text a reviewer is given cannot differ by a second pass
- [x] 13.2 Read the guide's size from the base commit's tree entry, so a guide over the ceiling is refused without being loaded
- [x] 13.3 Answer `only-in-the-change` only where the base carried nothing at the path, so a change that removes or empties a guide is not reported as one that adds it
- [x] 13.4 Name the run's field `ReviewGuide` and serialise it as `review_guide`, so nothing competes with the `guidance` key `findings --json` already publishes
- [x] 13.5 Enumerate `EventGuidanceUnread` and `EventGuidanceOverTheCeiling` in `TestEveryKindTheCoreWritesIsRendered`, so the guard that catches an unrendered kind covers the kinds this change adds
- [x] 13.6 Bring `proposal.md` and this task list back in step with what was built, since `openspec/` is the normative source and a stale one is worse than none
