**The pure halves are decided over their input space**, not sampled: what a document accepts, what neutralising removes, and what the renderer emits.

**The page half is driven by a real browser**, against the real binary on a real loopback address, and what it asserts is read back from `poll` on the other side. Asserting on served markup would test the renderer and claim to have tested the page.

**The surface half is driven by a plain HTTP client**, because the bind guards, the token, the bounds and the delivery-once rule are not about the client and a browser there proves less.

**Nothing here knows what a finding is.** The words review, finding, run, round and verdict appear nowhere in this binary. A test that used them would be describing the next change.

**Each guard is mutated on purpose before it is trusted:** the test that proves the new behaviour is run against the old code first.

**`make lint test` after every task, not at the end of a group.**

## 1. The module makes room for a second binary

- [ ] 1.1 Give `NetworkReach` the scope it is asked about, leaving the rule it applies unchanged
- [ ] 1.2 Name socrates' own binary and packages in the test that calls it, and confirm it still fails when something socrates reaches imports the network. Covers `socrates reaches nothing new`
- [ ] 1.3 Add a package outside socrates' reach that imports the network, and confirm what socrates reaches is unchanged. Covers `A second binary in the module is not socrates reaching further`
- [ ] 1.4 Confirm the `socrates` command surface and its run records are untouched by anything in this change. Covers `The surface does not grow`
- [ ] 1.5 Mutate the narrowed scope to include the whole module again, and confirm 1.2 fails rather than passing quietly

## 2. The document contract

- [ ] 2.1 Define the document and its decisions, options, modes and evidence, in a package that imports nothing outward
- [ ] 2.2 Validate a document over its failure space: a missing field, an unknown mode, an option with no command, a decision with no options and no field. Covers `A document that does not validate is refused before serving`
- [ ] 2.3 Take each decision's identifier, severity, audience and location from the cited source, and prove a document disagreeing with its source loses. Covers `A caller cannot rename an identifier`
- [ ] 2.4 Refuse a decision whose identifier appears in no source entry, naming it. Covers `A decision naming nothing in the source is refused`
- [ ] 2.5 Add a contract test reading a checked-in sample of the source shape, so the two cannot drift unnoticed

## 3. Neutralising what a caller supplied

- [ ] 3.1 Reduce a fragment to a fixed allowlist of elements and attributes
- [ ] 3.2 Strip a script, an event-handler attribute, a remote reference and a form control, and confirm serving still succeeds. Covers `A fragment outside the allowlist is stripped`
- [ ] 3.3 Confirm inline SVG within the allowlist survives unchanged. Covers `An illustration survives neutralising`
- [ ] 3.4 Mutate the allowlist to permit a script, and confirm 3.2 fails

## 4. Rendering a document to a page

- [ ] 4.1 Render a document to one HTML page in a package that opens no socket and reads no clock
- [ ] 4.2 Render a decision as a card from the document alone, with no markup supplied by the caller. Covers `A document renders without the caller writing markup`
- [ ] 4.3 Render the free-text field on every decision in both modes, and confirm no document field removes it. Covers `The free-text field cannot be suppressed`
- [ ] 4.4 Enumerate the rendered controls and confirm each belongs to answering or sending. Covers `The surface offers no third interaction`

## 5. Serving, and only on loopback

- [ ] 5.1 Build `cmd/agora` with `serve`, `poll`, `end`, `stop` and `version`, and confirm the surface holds those and no more. Covers `The surface offers nothing beyond the five`
- [ ] 5.2 Confirm `socrates` and `scrat` are still one program and `agora` is a separate one. Covers `socrates keeps its two names and gains no third`
- [ ] 5.3 Return from `serve` as soon as the page is reachable, printing the address. Covers `serve returns before the person has read anything`
- [ ] 5.4 Refuse a bind address that is not loopback, naming it and starting no server. Covers `A wildcard bind is refused rather than quietly narrowed`
- [ ] 5.5 Reject a request naming a host the surface did not bind, before any handler runs. Covers `A request claiming another host is rejected`
- [ ] 5.6 Mint a per-session token, embed it in the page, and reject a state-changing request without it. Covers `An answer without the token is rejected`
- [ ] 5.7 Write the session file and the served document under the user's state directory, and confirm a repository the surface was started inside is byte-identical afterwards. Covers `The repository is untouched`

## 6. The exchange

- [ ] 6.1 Block in `poll` until the person sends, and print each decision's identifier, choices, words and command
- [ ] 6.2 Round-trip two decisions and a message, and confirm the exchange is reported as not over. Covers `An answer round-trips`
- [ ] 6.3 Confirm the command `poll` prints is the option's own string, verbatim. Covers `The command is what comes back`
- [ ] 6.4 Confirm no process is started and nothing outside the session file is written when an answer arrives. Covers `The surface runs nothing`
- [ ] 6.5 Carry an answer that is words alone, with no option chosen. Covers `Words alone are a complete answer`
- [ ] 6.6 Report ending distinguishably, and refuse to reopen a closed exchange. Covers `Ending is distinguishable`
- [ ] 6.7 Deliver an answer to one `poll` only. Covers `Delivery happens once`
- [ ] 6.8 Wait rather than fail where `poll` starts before the page is open. Covers `A poll that started before the page was open still waits`
- [ ] 6.9 Bound the wait and return saying nobody answered, tested with `testing/synctest`. Covers `Nobody answers`
- [ ] 6.10 Shut the server down after the idle period with no session open, freeing the port. Covers `An idle server stops`
- [ ] 6.11 Mutate the delivery-once rule to deliver twice, and confirm 6.7 fails

## 7. The page, driven by a real browser

- [ ] 7.1 Stand up the end-to-end harness: the real binary, a real loopback address, a real headless browser, and `poll` read on the other side
- [ ] 7.2 Select, deselect and reselect within a decision, and confirm nothing was delivered. Covers `Choosing sends nothing`
- [ ] 7.3 Tick three options in a multiple-choice decision, commit, and confirm all three arrive in document order. Covers `Multiple choice carries every tick`
- [ ] 7.4 Commit, change, commit again, and confirm one answer arrives and it is the later one. Covers `Committing again replaces the earlier answer`
- [ ] 7.5 Confirm a committed decision renders differently from one merely selected in, and shows its command. Covers `A committed answer is distinguishable`
- [ ] 7.6 Reopen a committed decision before sending, and confirm its previous answer is not delivered. Covers `A committed answer can be taken back`
- [ ] 7.7 Send without ending, and confirm the page stays open and accepts a re-render. Covers `Sending without ending keeps the page open`
- [ ] 7.8 Send and end, and confirm the page says so and accepts no further answer. Covers `Ending closes the exchange`
- [ ] 7.9 Write a new document for an open surface, and confirm the page takes it up without a reload. Covers `A re-render reaches an open page`

## 8. Evidence, embedded

- [ ] 8.1 Add `web/` holding the client source and the fonts, and `make assets` building it into `internal/assets`
- [ ] 8.2 Commit the built bundle and fail the suite while it differs from what its source builds, as `make skill` already does
- [ ] 8.3 Render a patch as a diff with both sides and per-line changes visible. Covers `A diff renders from a patch`
- [ ] 8.4 Confirm the page renders with the machine's outbound network unavailable. Covers `Rendering needs no network`
- [ ] 8.5 Confirm no diagram engine is embedded, and measure the binary so its weight is a number somebody has seen

## 9. The floor the page is held to

- [ ] 9.1 Measure the rendered page's text and confirm none is below the stated minimum. Covers `Nothing renders below the floor`
- [ ] 9.2 Render under a light preference and a dark one, and confirm each paints its own background and text with no token left undefined. Covers `Both preferences are painted`
- [ ] 9.3 Reach, name and operate every option, field and send by keyboard alone. Covers `A decision can be answered without a pointer`
- [ ] 9.4 Confirm at a narrow width that the diff and any illustration stay legible within their own bounds and the page does not scroll sideways. Covers `The narrow case does not shrink the evidence to nothing`
- [ ] 9.5 Record every request the loaded page issues and confirm each resolves to the surface itself. Covers `The page reaches nothing`

## 10. Closing out

- [ ] 10.1 Run the full sequence: `make fmt`, `make lint`, `make test`, `make scenario-coverage`, `make scenario-drop`, `make corpus-outcomes`
- [ ] 10.2 Confirm no package socrates reaches gained an import, and that `internal/reviewer`, `internal/runstore`, `internal/dialogue`, `internal/findings` and `internal/skill` are untouched
- [ ] 10.3 Confirm the words review, finding, run, round and verdict appear nowhere in the new binary or its packages
