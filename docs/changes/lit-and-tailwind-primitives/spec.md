# agora's primitives are Lit elements styled with Tailwind

Decisions: ADR 82, ADR 83, ADR 84, ADR 85, ADR 86.

## Problem Statement

A contributor changing how agora's page looks or behaves works in a house style nobody else uses.
A primitive has two halves in two languages: Go writes its markup as a partial in `agora/render/partials`, and `agora/web/packages/ui` gives it a stylesheet and, where it has behaviour, a light-DOM custom element, held together by `contract.json`.
Changing one primitive means editing a `.gohtml` file, a `.css` file and a `.ts` file, and keeping the three in step.
The styling is about 2,260 lines of hand-written CSS across the primitives, the catalog and the surface, most of it long selectors such as `ui-button[data-variant='primary'] > [data-part='control']`, with raw px spacing and no spacing scale.
The behaviour is imperative DOM code: elements query their parts by `data-part` and mutate them by hand.
The landing page in `site/` is styled the same way, in a stylesheet of its own.
Lit and Tailwind are mainstream, known to contributors and to agents, and would let the client read like other web code.

## Solution

Every primitive becomes a Lit element that renders its own parts into a shadow root, styled with Tailwind utilities, with its template, its styling and its behaviour in one TypeScript file, per ADR 82 and ADR 83.
Go keeps rendering the page: a partial writes a primitive's element, the attributes it declares and the content it slots, and nothing inside it.
`make assets` writes the attributes and slots every element declares to `contract.json`, and a Go test holds every partial to it.

Tailwind styles everything the client draws: the primitives, the pieces the glossary names, the catalog and the landing page.
Its own theme is cleared, and every utility is a token named in `tokens.css`, which stays the one place a token is stated, per ADR 84.
No class string carries an arbitrary value.
Preflight comes in wherever Tailwind is, in every shadow root and on the document.

The pieces the glossary names stay Go markup, composing primitives and laid out by utilities in the `.gohtml` files, and the surface's hand-written stylesheet empties, per ADR 85.

The surface's script and the browser suite reach a primitive only through the properties and events it exposes, never inside its shadow root.
A decision's choices are one `<ui-choice-group>`, which owns selection, the arrow keys and the radio-group role once native radios can no longer group across shadow roots, per ADR 86.

A person answering a page sees and does what they see and do today.

## User Stories

1. As a contributor changing how a primitive looks, I want its styling as utilities in its own template, so that I read what it draws in one place rather than tracing selectors across a stylesheet.
2. As a contributor changing a primitive, I want its template, its styling and its behaviour in one TypeScript file, so that I do not keep a Go partial, a stylesheet and an element in step by hand.
3. As a contributor adding a primitive, I want to write one Lit element, a one-line partial and its states in the catalog, so that adding a primitive is no longer a partial, a contract entry, a folder, two imports and a catalog entry.
4. As a contributor giving a primitive behaviour, I want reactive properties and a declarative template, so that the element re-renders from its state rather than my code querying and mutating its parts.
5. As a contributor who knows Lit and Tailwind from elsewhere, I want agora's client to use them the way their own docs do, so that I can work on it without first learning a house style.
6. As an agent working on the client, I want the client written in mainstream Lit and Tailwind, so that what I already know about both applies without translation.
7. As a contributor changing only a partial, for example to pass a new attribute, I want to need no Node, so that a Go-only change stays Go-only.
8. As a contributor who writes an attribute or a slot on a partial that no element declares, I want a Go test to fail naming the element and the attribute or slot, so that a misspelling is caught without a browser.
9. As a contributor who declares an attribute or a slot that no partial writes, I want the same test to fail naming it, so that the contract cannot drift from what the page writes.
10. As a contributor reaching for a colour, a font, a radius or a shadow, I want only the tokens to exist as utilities, so that `bg-red-500` or `rounded-lg` cannot paint the page with a value the design did not state.
11. As a contributor writing `p-[13px]` or `bg-[#fff]`, I want the suite to fail naming the file and the line, so that I add a token or use a step of the scale instead.
12. As a contributor spacing something, I want one spacing scale whose steps are the 2, 8, 12 and 16px the page already uses, so that spacing is a step rather than a number I choose.
13. As a contributor laying out a piece of the page, such as a decision or the composer, I want to write utilities in the Go markup that draws it, so that a piece's layout sits with its markup as a primitive's does.
14. As a contributor working on the surface's script, I want to read and set a primitive through its own properties, such as a field's `value`, so that my code does not depend on the markup inside it.
15. As a contributor writing a browser test, I want to click, type at and read a primitive by its element and its class, so that a test does not break when a primitive's internal markup changes.
16. As a person answering a single-choice decision, I want at most one option chosen and the arrow keys to move between the options and *None of these*, so that the decision behaves as a radio group does today.
17. As a person using a screen reader, I want a decision's options announced as a radio group or a set of checkboxes with each one's state, so that I know what I am choosing and what is chosen.
18. As a person answering a page, I want it to look and behave as it does today, so that the change to how it is built costs me nothing.
19. As a person answering a page, I want the page's policy to stay as strict as it is, so that building the page differently does not widen what it admits.
20. As a person whose illustration loads `/assets/tokens.css`, I want that file to keep stating every token it states today, so that my drawing still looks like the page.
21. As a contributor changing the landing page, I want to style it with the same utilities as the page, so that the site and the page are built one way.
22. As a maintainer publishing a release, I want the landing page's stylesheet already built and committed, so that `make site` and the release job need no Node.
23. As a maintainer, I want `NOTICE` to name Lit and Tailwind and the shipped bundle to carry their notices, so that what the binary carries that socrates did not author is still stated in full.
24. As a maintainer, I want `main` to move from today's build to the new one in one merge, so that `main` never holds two ways of building a primitive.

## Implementation Decisions

### Dependencies and the build

- `lit` becomes a dependency of `@agora/ui`, at the latest stable 3.x.
- `tailwindcss` becomes a dependency of `@agora/ui`, because its theme layer, Preflight and utilities ship in the bundle's stylesheet, and `@tailwindcss/node` becomes a development dependency of `agora/web`, driven from `build.mjs`.
  Tailwind has no official esbuild plugin.
- `build.mjs` compiles the library's Tailwind stylesheet once, with `@source` naming `agora/web/packages/**/*.ts`, `agora/render/**/*.gohtml` and `site/index.html`, and hands the compiled CSS to esbuild.
  The client's TypeScript imports the compiled CSS as text through a module the build resolves, from which `packages/ui` builds one constructable `CSSStyleSheet` that every primitive adopts as its `static styles`.
  The same compiled CSS is bundled into `surface.css`, which the page and the catalog load on the document.
- `build.mjs` also writes the landing page's stylesheet, compiled from its own Tailwind entry in `site/` against `site/index.html`, to a committed generated file in `site/`, and `site/` joins the sources `bundle.sum` stamps.
  `make site` copies the generated file and builds nothing.
- `make lint` still runs no type checker, and `make verify` still needs no Node, per ADR 79.

### The theme

- The library's Tailwind stylesheet opens with `@theme { --*: initial; }` and names every token in `tokens.css` that a utility paints with in an `@theme inline` block: colours, font families, the type scale, leadings, trackings, radii, shadows, durations and easings.
- `tokens.css` gains `--spacing`, set so the scale's steps are the 2, 8, 12 and 16px the stylesheets use today, and stays a single `:root` block that illustrations load at `/assets/tokens.css`.
- The palette check keeps reading the client's stylesheets and the landing page's source stylesheet, and does not read a generated stylesheet.
  The library's Tailwind entry states no colour and no font family outside `var(--...)`, so the check holds over it unchanged.
- A new convention check reads every class string in `agora/web/packages/**/*.ts`, `agora/render/**/*.gohtml` and `site/index.html` and fails on a utility written with an arbitrary value, naming the file and the line.
  It reads the same paths the `@source` lines name, stated once in `internal/convention`.

### The primitives

- Each primitive in `packages/ui` is one Lit element declared with standard decorators, `@customElement` and `@property`, without `experimentalDecorators`.
  Its folder holds its element and, where a primitive needs one, a module of its own, such as the code block's highlighting.
  The per-primitive `.css` files, `primitives.css`, `parts.ts` and `contract.ts` are removed.
- Every primitive renders into an open shadow root, Lit's default, and controls set `delegatesFocus`, so `focus()` and a click on the element reach the control inside.
- What Go writes that a person reads is slotted rather than passed as an attribute, so a button's label, a choice's label and meaning, and a card's body are children of the element.
- A primitive exposes what its callers need as properties, methods and composed events: a field's `value` and `focus()`, a choice's `checked` and `value`, a button's `disabled` and `click()`, the dialog's `open()` and `close()`, the copy control's `copyrequested` event, and the expand control's `expanded` and `collapsed` events.
  The surface's TypeScript never names `shadowRoot`, and the build refuses a file outside `packages/ui` that does, naming the file and the line.
- `ui-choice-group` is a new primitive.
  It slots `ui-choice` children, keeps one checked where it is single-choice, roves `tabindex` and moves on the arrow keys, and carries `role="radiogroup"`, with each choice carrying `role="radio"` or `role="checkbox"` and `aria-checked`.
  *None of these* is a `ui-choice` of the dashed variant in the same group.
- `data-part` is retired: a primitive's parts are inside its shadow root, and nothing outside it names them.
  `data-interaction` moves to the primitive's element and still says which of answering, sending and viewing a control belongs to.
- The tooltip still shows for anything carrying `data-title`, which now sits on the primitive's element where the anchor is a primitive.
- Until the bundle defines them, primitives the page writes paint nothing, through a `:not(:defined)` rule in the document's stylesheet, so no unstyled slotted text flashes on load.
- The page's Content-Security-Policy is unchanged.
  Constructable stylesheets are not governed by `style-src`, and no inline `style` attribute or `<style>` element is written.

### The contract

- `build.mjs` reads, from the TypeScript program it already builds, every class decorated with `@customElement`: its element name, the attribute each `@property` reflects or reads, and the name of every `<slot>` its `render()` template holds, the default slot included.
- `contract.json` becomes, for each element, its attributes and its slots.
- The Go test in `agora/render` renders every partial and the catalog and fails on an attribute a `ui-*` element carries that its element does not declare, on a `slot` attribute naming a slot it does not render, and on a declared attribute or slot no partial writes, naming the element and the attribute or slot.
  Global attributes, `class`, `id`, `hidden`, `title`, `role`, `aria-*` and `data-*`, are not part of the contract.

### The partials and the pieces

- Each `ui-<name>.gohtml` shrinks to the element, its attributes and its slotted children.
- `page.gohtml`, its partials and the catalog lay out the glossary's pieces with utilities, and each piece keeps the one class that names it for the script and the suite.
- `surface.css` keeps only what utilities cannot say, such as the `:not(:defined)` rule, and `catalog/catalog.css` goes the same way.
- `convention.PiecesOfTheSurface` and the check reading it are removed.
- The landing page in `site/` is restyled with utilities against the same theme, and its hand-written `site.css` becomes its Tailwind entry.

### Licensing

- `NOTICE` gains Lit, BSD-3-Clause, bundled into `bundle/surface.js`, and Tailwind CSS, MIT, whose theme, Preflight and utilities are bundled into `bundle/surface.css` and the landing page's stylesheet, each with its full licence text.
- The notice test reads every library in a workspace package's `dependencies`, as it does today, and asserts each one's notice against the shipped file its code is bundled into, so Lit's `@license` comments are asserted against `surface.js` and Tailwind's banner against `surface.css`.

### Documents

- `agora/AGENTS.md` is rewritten where it describes `agora/web`, `agora/render` and the primitives: how a primitive is added, the shadow root, the contract of attributes and slots, the theme, and the invariants this change replaces.
- `CONTRIBUTING.md` still asks for Node and `pnpm` only of somebody who touches `web/`, and now names `site/` beside it.
- `agora/CONTEXT.md` keeps its definition of a primitive and gains no word: a choice group is a primitive, which names nothing in agora's domain.

### Landing

- The change lands as tickets on one change branch, merged to `main` once, as the previous client change did.
  The foundation comes first: Lit, the Tailwind build, the theme, the contract of attributes and slots, and the arbitrary-value check.
  Then the primitives in groups, then the pieces and the catalog, then the landing page, each ticket passing `make verify` on the branch.

## Testing Decisions

The change is tested at the seams that already exist, and no new kind of test is added.

- **The catalog in a real browser.**
  Each primitive is proved against `/catalog` in chromedp, in every state the catalog draws, as today.
  The existing catalog tests are the regression net, and change only where a selector moves from a control inside a primitive to the primitive's element.
  `ui-choice-group` gets scenarios of its own: one choice checked at a time in a single-choice group, the arrow keys moving the checked choice and the focus, Tab leaving the group from the checked choice, and a multiple-choice group checking several.
- **The real page in a real browser.**
  The pieces are proved against the real binary as today, with the suite clicking, typing at and reading primitives by their element and class.
  The test asserting that no primitive has a shadow root becomes one asserting that every primitive has one.
  The suite's collected policy violations stay empty on the page and the catalog, which proves the adopted stylesheets admitted under the unchanged policy.
- **The Go contract test.**
  It fails on an undeclared attribute or slot and on a declared one nobody writes, each mutated on purpose once to prove it fails.
- **The convention checks.**
  The arbitrary-value check is proved over a scratch tree holding one in a `.ts` file, one in a `.gohtml` file and one in `site/index.html`, and the palette check is proved still to fail on a colour written in the library's Tailwind entry outside `tokens.css`.
- **The bundle and the notice.**
  `bundle.sum` covers `site/`, so a stale landing stylesheet fails the suite, and the notice test fails where Lit's or Tailwind's notice is missing from the file it ships in.
  `scripts/prove-client-build` gains the case of a surface file naming `shadowRoot`.

No component test runner in Node is added, so `make verify` still needs no Node.

## Out of Scope

- Lit rendering the whole page, or Go handing the client the document as data.
- The glossary's pieces as Lit elements.
- Server-side rendering a primitive, or declarative shadow DOM written by Go.
- A component test runner such as Web Test Runner.
- Redesigning the page: what a person sees and does stays as it is.
- Changing the page's Content-Security-Policy or the illustration's.

## Further Notes

- The facts the interview rested on: chromedp's selectors reach nothing inside a shadow root as the suite writes them today; the page's policy is a `<meta>` tag with `style-src 'self'` and `style-src-attr 'none'`; Tailwind 4.3 runs no script in the browser, emits its theme on `:root, :host`, and registers `@property` only at document level; Lit 3.3 is BSD-3-Clause and its shipped files carry `@license` comments esbuild keeps; Tailwind's output opens with an MIT banner esbuild keeps.
- That constructable stylesheets fall outside `style-src` was not confirmed in this repository during the design; the first ticket proves it with the suite's violation collection before anything is built on it.
- ADR 80 and ADR 81 are marked superseded by ADR 82 and ADR 83.
  ADR 79 stands.
