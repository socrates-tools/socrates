**The reduction is the whole change and it is tested as a pure core.** Its input space is large and a wrong answer is expensive, which is the shape this repository writes exhaustive unit tests for. Everywhere else, prefer the end-to-end test.

**Each guard is mutated on purpose before it is trusted.** A construct the allowlist refuses is admitted by hand first, and the test that refuses it must fail before it is believed.

**Nothing here asserts that a session illustrated well.** Which form a session picks, and whether the panel it drew resembles anything, is judgement. What is driven is the instruction the skill carries and the mechanism the surface offers.

**The page half is driven by a real browser against the real binary.** Asserting on served markup would test the renderer rather than what a person gets.

**`make lint test` after every task, not at the end of a group.**

## 1. The reduction reads a parse

- [x] 1.1 Add the CSS parser dependency, behind an interface owned by the new package rather than its types leaking out of it. Verify `go mod tidy` leaves the module graph with one new direct dependency and no transitive ones
- [x] 1.2 Reduce a stylesheet by walking its parse, emitting only what the allowlist admits. Covers `A stylesheet outside the allowlist is reduced`
- [x] 1.3 Embed nothing from a stylesheet whose parse does not reach the end of the input. Covers `A stylesheet that does not parse embeds nothing`
- [x] 1.4 Mutate the parse failure into a partial emit, and confirm 1.3's test fails rather than passing on what was read

## 2. What the reduction admits

- [x] 2.1 Admit only the at-rules the allowlist names, and drop imports, font faces and the rest with their contents. Covers `A stylesheet outside the allowlist is reduced`
- [x] 2.2 Admit only the properties the allowlist names, and drop a declaration whose value names a resource, including one written with escapes or split by a comment. Covers `A stylesheet outside the allowlist is reduced`
- [x] 2.3 Admit custom properties and the references that read them. Covers `A caller may state its palette once`
- [x] 2.4 Refuse a stylesheet above the stated size at the point a document is validated, before any server starts. Covers `An oversized stylesheet is refused before serving`
- [x] 2.5 Mutate each of 2.1, 2.2 and 2.4 by admitting the refused construct, and confirm the test for each fails on its own rather than another catching it

## 3. The surface scopes, and renames

- [x] 3.1 Prefix every surviving selector with the illustration's scope, and drop a selector that cannot be scoped. Covers `A selector naming the page is confined to the illustration`
- [x] 3.2 Confirm a rule the caller wrote for its own names still paints after scoping. Covers `What the caller painted is still painted`
- [x] 3.3 Prefix every `class` and `id` a fragment carries, and the matching tokens in the stylesheet's selectors and in its `url(#...)` references, so that neither can collide with a name the page uses. Covers `A fragment cannot claim the page's own names`
- [x] 3.4 Confirm a drawing whose gradients, clips and masks are referenced by identifier still renders after renaming. Covers `An illustration survives neutralising`
- [x] 3.5 Mutate the scoping to emit one selector unscoped, and confirm 3.1's test fails rather than the containment in group 5 catching it

## 4. The document carries an illustration and a look

- [x] 4.1 Make the fragment an object carrying its markup, and confirm a document written against the old string is refused by name rather than silently ignored. Covers `An illustration the stylesheet paints survives neutralising`
- [x] 4.2 Add the document's stylesheet, applying to every illustration in it. Covers `One stylesheet paints every illustration in a document`
- [x] 4.3 Confirm a document carrying an illustration and no stylesheet serves and renders it unpainted. Covers `A document that carries no stylesheet still serves`
- [x] 4.4 Confirm the validation refuses a stylesheet before a server starts, alongside the fields it already refuses. Covers `A document that does not validate is refused before serving`

## 5. The page embeds what was reduced

- [x] 5.1 Emit the scoped stylesheet into the page once, however many illustrations it paints. Covers `One stylesheet paints every illustration in a document`
- [x] 5.2 State the policy's hash from exactly the bytes emitted, and confirm in a real browser that the illustration paints and no policy violation is reported. Covers `The embedded stylesheet paints`
- [x] 5.3 Confirm the policy admits no inline style beyond that hash, and that style carried on an element is not admitted. Covers `Inline style is not admitted as a class`
- [x] 5.4 Make the illustration a box that contains and clips what is painted in it, and regenerate the bundle with `make assets`. Covers `A painted illustration cannot cover the page`
- [x] 5.5 Drive a stylesheet that positions and oversizes an element over the page's controls, and confirm in a real browser that every control remains visible and operable. Covers `A painted illustration cannot cover the page`
- [x] 5.6 Render a drawing and a painted panel in one document, and confirm both land in the illustration's box. Covers `A painted illustration renders in the same box as one that paints itself`
- [x] 5.7 Confirm a painted illustration at a narrow width stays legible and scrolls within its own bounds. Covers `A painted illustration is legible at a narrow width`
- [x] 5.8 Confirm the page's frame, palette and controls are painted identically with and without a caller's stylesheet that names them. Covers `A caller cannot paint the frame`
- [x] 5.9 Confirm the contrast sweep still sees the whole page, still steps around the illustration, and still holds every text the page itself paints. Covers `A caller's colours do not reach the decision`
- [x] 5.10 Mutate the containment away, and confirm 5.5's test fails rather than the scoping in group 3 catching it

## 6. The bar, and the skill beside it

- [x] 6.1 Add the bar's line about showing what a finding costs where that is what a person sees, in the one place the bar is stated. Covers `The bar names each way the page fails a person`
- [x] 6.2 Confirm the line about showing and the line about diagrams each appear once and neither restates the other. Covers `Showing and drawing are separate lines`
- [x] 6.3 Say beside the bar that a decision may carry one illustration, travelling beside the evidence rather than inside it. Covers `The skill names the illustration as something a decision may carry`
- [x] 6.4 Say an illustration takes two forms and which is for what, and confirm the rendered skill names no diagram engine and implies none the surface would run. Covers `The skill says what forms an illustration takes`
- [x] 6.5 Say the document carries one stylesheet for all of its illustrations. Covers `The skill says the look is stated once per document`
- [x] 6.6 Say what does not survive the reduction of each of the markup and the stylesheet. Covers `The skill says what is stripped before either is embedded`
- [x] 6.7 Say the surface does the scoping, so a session neither scopes its own selectors nor expects one to reach the page. Covers `The skill says the surface does the scoping`
- [x] 6.8 Confirm the mechanism sits beside both bar lines and that the bar's judgement about when each is warranted appears once and only there. Covers `How to illustrate is stated beside when to`
- [x] 6.9 Mutate the bar's new line, and confirm the test holding the bar to one statement fails rather than the mechanism's alone
- [x] 6.10 Regenerate `.claude/skills/socrates/SKILL.md` with `make skill`, and confirm `make test` fails while the committed copy and `internal/skill` disagree
- [x] 6.11 Confirm the skill still names nothing the binaries do not offer, and that the prohibitions it is held to are unchanged

## 7. The superseded change

- [x] 7.1 Remove `openspec/changes/a-diagram-the-skill-can-draw`, and confirm its requirements are carried by this change's `agent-skill` delta rather than dropped

## 8. The whole gate

- [x] 8.1 Confirm nothing in the surface names review, finding, verdict, round or run as a noun, and that no user-facing string this change adds carries an em dash
- [x] 8.2 `make fmt`, `make lint`, `make test`, `make scenario-coverage`, `make scenario-drop`, `make corpus-outcomes`
