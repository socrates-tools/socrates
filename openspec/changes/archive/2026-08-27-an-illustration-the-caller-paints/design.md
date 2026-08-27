## Context

See proposal.md for why.

What already exists is most of the frame.
`document.Decision` carries a fragment, `sanitise.Fragment` reduces it to a fixed allowlist of elements and attributes, and the page renders the survivor in a box of its own.
`surface-page` already treats what a caller supplied as the caller's own, and the presentation sweeps already step around that box.

Three things stop that box holding a panel rather than a drawing.
`style` is not an allowed attribute and `<style>` is dropped with its content.
The page states a policy that admits no inline style, so an allowed attribute would be dropped by the browser anyway.
And the box styles a drawing rather than markup, so markup put in it renders as body copy.

The third is a stylesheet edit.
The first two are the same question asked twice, which is whether a caller may paint, and that question is what this change answers.

## Goals / Non-Goals

**Goals:**

- A caller may paint the inside of the illustration well enough that a person recognises a panel from the product they are being asked to decide about.
- Nothing a caller writes can reach, move, recolour or cover the page's own frame, palette or controls, and that holds mechanically rather than by the caller being careful.
- A reduction that cannot do its job fails rather than passing what it managed to read.

**Non-Goals:**

- Fidelity to a real product's brand. What a caller can paint is what a caller can describe in a reduced stylesheet, and a mock that is recognisable is the bar rather than one that is exact.
- Any reach from socrates into agora. The session writes the document, as it already does, and no command, flag or run record learns that a page exists.
- Interaction inside an illustration. It shows a person something. Everything a person operates is the page's, and the reduction that drops form controls from the markup is not softened.
- A general styling channel. One box per decision, painted by one stylesheet per document, and no second use.

## Decisions

### The stylesheet belongs to the document, the markup to the decision

A round is several findings about one product, and a product has one look.
Per-decision stylesheets are the same bytes repeated in the context that is about to do the patching, and they are two panels of one mock free to disagree about what a panel is.

The alternative considered was a stylesheet per decision, on the argument that it keeps everything about one decision in one place.
It was rejected because the thing being described is not the decision, it is the product, and the document is the only object in the contract that is about the product.

### The fragment becomes an object

An illustration is either markup a stylesheet paints or a drawing that paints itself, and a field that can only be a string can only be the second.

This is a breaking change to the document contract and it is free right now.
The only change that would have written against the string is `a-diagram-the-skill-can-draw`, which this supersedes and which was never implemented.
Doing it later means a migration for a contract that has no other reason to move.

### The reduction reads a parse, not a pattern

A CSS reduction built on pattern matching is a gate that reports a pass it cannot justify, which is the failure the repository already refuses in `internal/coverage` and `internal/promotion`.
Escapes, comments inside values, nested at-rules and a `url` written as `\75 rl` all defeat matching and none of them defeat a parse.

This brings a dependency, and the repository carries only `golang.org/x/net` today.
`github.com/tdewolff/parse/v2/css` is the choice: a tokeniser and parser, no transitive dependencies, and it hands back structure rather than an opinion about it.

Alternatives considered.
Hand-writing a tokeniser was rejected because the correctness of this gate is the whole change and a CSS tokeniser is a thing to maintain rather than a thing to write once.
`gorilla/css` was rejected because it is a scanner rather than a parser, which leaves the structure to be assembled here, which is the part that has to be right.
Regular expressions were rejected on the paragraph above.

Where the parse does not reach the end of the input, nothing from that stylesheet is embedded.
A partially reduced stylesheet is the case where the reduction saw less than the browser will.

### The surface scopes, and the caller does not

Every selector that survives is prefixed with the illustration's own scope before it is emitted, so the stylesheet the page carries cannot name anything outside the box, whatever the caller wrote.

The alternative was to require callers to scope their own selectors and reject the ones that are not scoped.
It was rejected for the reason the identifier-copying requirement already gives: a rule held by the caller's care is a rule one typo away from failing, and the page cannot tell a typo from an attempt.

A selector that cannot be scoped is dropped rather than emitted unscoped.
That covers `:has()`, `:host`, `::part()` and `::slotted()`, each of which can select a subject the prefix does not constrain.

### Names are rewritten rather than restricted

Every `class` and `id` a fragment carries is prefixed by the surface, and the matching tokens in the stylesheet's selectors and in its `url(#...)` references are prefixed with it.
The caller writes `panel` and the page carries the prefixed form in both places.

This solves two different problems with one pass.
An `id` collision is a functional break: a fragment carrying `id="decision-f3"` breaks what the client script addresses and what an anchor reaches.
A `class` collision is a cosmetic one in the other direction: the page's own chrome rules would match an element inside the mock and paint it as page furniture.

Alternatives considered.
Dropping `id` outright was rejected because SVG needs it: `clipPath`, `mask` and gradients are referenced by identifier, and dropping ids would break the drawing form this change also carries.
Requiring a reserved prefix from the caller was rejected as the same care-based rule as caller-side scoping, and it forces the prefix into the session's context on every name it writes.
A closed vocabulary of names the surface publishes was rejected because the caller is mocking an arbitrary product and cannot be given a fixed set of nouns for it.

### Containment is the second layer, and it is what makes the property allowlist tractable

Scoping decides what a rule may name. Containment decides where its result may land.

The illustration's box takes `contain: layout paint style`, which makes it a containing block for absolutely and fixed-positioned descendants and clips them to it.
A stylesheet that positions an element over the commit control therefore positions it over part of the mock instead, without the allowlist having to reason about positioning at all.

This is why `position` does not need to be restricted to `static` and `relative`, which was the first shape considered.
Restricting it would have made a common mock harder to write for a risk containment already removes.

### The stylesheet is admitted by hash

The page states a policy that admits nothing it did not bring.
Admitting inline style as a class so that one stylesheet can be embedded admits every inline style anything ever puts on this page, and the reduction stops being the only door.

A hash of exactly the bytes emitted keeps the door one-wide, needs no state carried from a render to a request, and fails closed when the two disagree.
A nonce was considered and is equivalent in strength here, since the page is rendered per request and never stored.
The hash was preferred because it is derived rather than generated, so there is nothing to get wrong about where it came from.

Style carried on an element is not admitted at all.
It cannot be scoped, and admitting it would need the policy to permit hashed attributes, which is a second door for a form that has no advantage over the first.

### The bar gains a line, and this is a departure

`a-diagram-the-skill-can-draw` deliberately left the bar alone, on the ground that the bar states what the page owes a person and not which field carries it.
That was right for a change that added a mechanism for a bar line that already existed.

This change adds a reason to illustrate that the bar does not have.
The line about diagrams is about workflows and ends "and not otherwise", so a page that showed a person their own screen would be a page in breach of the bar as it stands.
Adding the mechanism without the line ships a capability the bar forbids.

The line is separate from the example line beside it because the readers differ.
The example line is for a reader who has to be taught what a term costs them.
This one is for a reader who already knows the screen and is being asked to recognise it, and a panel retyped as brackets is the one form that reader cannot recognise it in.

### The contrast floor does not reach a caller's colours, and the specification says so

The floor is measured against a stylesheet that is read.
A caller's colours arrive when the document does, and nothing reads them before a person is looking at them.

Holding them to the floor would mean measuring at serve time and refusing a document whose palette failed, which turns a document into something that can be refused for taste.
The floor stays what the page paints, the existing carve-out is widened to name a caller's stylesheet, and the compensating property is that a caller's colours land only inside the box that was showing the person something and never in the decision they came to make.

## Risks / Trade-offs

- **The reduction is the only thing between a caller and the page, and a gap in it is not visible from the page.** → It is a pure core with a large input space and an expensive wrong answer, which is the shape this repository already writes exhaustive unit tests for. Every construct the allowlist refuses is mutated on purpose before the guard that refuses it is trusted.
- **A new dependency in a repository that has one.** → It is a parser with no transitive dependencies, it is used in exactly one package, and the seam is the reduction's own interface rather than the parser's types, so replacing it is one file.
- **Two illustration forms means two ways to do one thing, and a session may pick wrong.** → The bar carries a line for each and the skill says which is for what. Which one a session picks is judgement, and nothing here asserts that it picked well, because a test claiming a page carried the right kind of picture would read as passing while checking nothing.
- **A caller can paint something illegible.** → It can, inside its own box, and the floor already does not cover what a caller supplied. This is a widening of an existing carve-out rather than a new one.
- **Containment could clip a mock a caller meant to overflow.** → The box scrolls within its own bounds, which is already required of the illustration at a narrow width.

## Migration Plan

Nothing is deployed and nothing is written against the contract that moves.
`fragment` goes from string to object in one commit, and the only document that carries one is the one written by the session reading the skill, which is regenerated in the same change.

`openspec/changes/a-diagram-the-skill-can-draw` is removed as part of this change rather than archived, because it was never implemented and its requirements are carried here.

Rollback is the absence of a stylesheet: a document that carries none serves and renders its illustrations unpainted, which is the behaviour that exists today.

## Open Questions

- Whether `serve` should report what the reduction dropped, so a caller learns its stylesheet was partly ignored rather than discovering it on the page. It changes no requirement here and can be answered once there is a real stylesheet to be surprised by.
