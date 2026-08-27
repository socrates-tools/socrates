## Why

A finding whose cost is what a person sees is put on the page as a sketch in a monospace block.
A panel is drawn with spaces, its controls are drawn with square brackets, and the person is asked to read that back into the screen they know.

The bar already requires an example where a finding is purely technical, and the surface already renders one fragment per decision.
What the fragment cannot do is look like anything.
`class` survives neutralising but `style` does not, `<style>` is dropped with its content, and the page's own content security policy admits no inline style at all, so the only illustration the surface can carry today is inline SVG with presentation attributes on it.
That is a drawing.
A session writing one computes its own text metrics, nothing wraps, and what arrives is recognisably a diagram of a panel rather than a panel.

The person deciding is the person who built the thing.
What makes them recognise a panel is its structure, its controls and the words actually on it, and none of that survives being retyped as brackets.
A finding that says two sentences contradict each other is answered in seconds by seeing them sit under one another, and in a paragraph otherwise.

`a-diagram-the-skill-can-draw` proposed teaching the skill that the fragment exists, aimed at inline SVG.
It was never implemented, so `fragment` is still a bare string that nothing has been written against, and the two changes edit one paragraph of the skill beside one line of the bar.
Shipping them separately means writing that paragraph twice and softening it once.
This supersedes it and carries its requirements.

## What Changes

- A document carries one stylesheet for its illustrations, and each decision's fragment carries the markup that stylesheet paints.
  One product has one look and a round has several findings, so the look is stated once per document rather than repeated per decision.
- The surface reduces that stylesheet before embedding it, the way it already reduces the fragment: an allowlist of at-rules, selectors, properties and values, and a parse it cannot complete is refused rather than passed through.
- Every selector that survives is scoped by the surface to the illustration, so the stylesheet the page embeds is structurally incapable of naming the page's own chrome.
- The illustration becomes a contained box, so what a caller paints cannot paint outside it.
- The stylesheet reaches the page under a hash the page states, rather than by the policy admitting inline style generally.
- **BREAKING** `fragment` becomes an object carrying markup rather than a bare string.
  Nothing is written against the string, because the change that would have written against it was not implemented.
- `class` and `id` on a fragment are constrained.
  Both pass through unconstrained today, so a fragment can already claim the page's own names, and that becomes load bearing the moment a caller can paint.
- The bar gains one line: where what a finding costs is what a person sees, show them.
  The bar's existing line about diagrams is a line about workflows and says "and not otherwise", so without this the mechanism arrives with nothing authorising its use.
- The skill says how an illustration reaches the page in both of its forms, beside the bar lines that call for one.
- **socrates reaches nothing new.** No command, no flag, no run record. The session writes the document, as it already does.
- `a-diagram-the-skill-can-draw` is superseded and its directory removed.

## Capabilities

### New Capabilities
- `surface-chrome`: a caller-supplied stylesheet, reduced to an allowlist and scoped to the illustration before the page embeds it, and admitted by a stated hash rather than by a policy that permits inline style.

### Modified Capabilities
- `surface-document`: `A caller writes a document, not a page` is narrowed and the narrowing is stated, because a caller may now paint inside one box; `A supplied fragment is neutralised before it is embedded` gains the stylesheet beside the markup.
- `surface-page`: `The page renders its own evidence` gains the containment the illustration is held in; `The page is painted in one stated palette, held to a contrast floor` extends its existing carve-out for what a fragment carries to cover what a caller paints, and says why the floor cannot reach it.
- `findings-page`: `The page is held to a stated bar, stated once` gains the line about a finding whose cost is what a person sees.
- `agent-skill`: the skill says how an illustration reaches the page, what forms it takes, and what does not survive the reduction. Carried from `a-diagram-the-skill-can-draw`.

## Impact

- `internal/document`: the fragment becomes an object, the document gains a stylesheet, and both are validated before a server starts.
- A new package for the stylesheet reduction, which is the bulk of the work and the only part with a real adversary.
- `internal/sanitise`: `class` and `id` constrained.
- `internal/render`: the scoped stylesheet emitted into the page, and the policy's hash computed from exactly the bytes emitted.
- `web/src/surface.css`: the illustration becomes a contained, sized box rather than a centred one; `make assets` regenerates the committed bundle.
- `internal/skill`: the passage beside the bar, and the bar itself; `make skill` regenerates the committed copy.
- **A new dependency**: a CSS tokeniser. The repository carries only `golang.org/x/net` today, and a reduction built on pattern matching rather than a parse is a reduction that reports a pass it cannot justify.
- `internal/e2e`: the illustration driven by a real browser, including that the page's own chrome is unmoved by what a caller painted.
- `openspec/changes/a-diagram-the-skill-can-draw`: removed.
