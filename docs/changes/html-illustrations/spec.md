# An illustration is whatever HTML can draw

Decisions: ADR 76, ADR 77.

## Problem Statement

A decision's illustration today is a fragment of markup that `agora/sanitise` reduces to a fixed allowlist, painted by one document stylesheet that `agora/stylesheet` parses, reduces and scopes.
The allowlist drops script, `img`, `style` attributes, `href`, `src`, and every `url(...)` but `url(#name)`.
So an agent can draw boxes, lines and imitation screens, and nothing more.

Issue 77 hit this from a review handover.
A UI finding such as "the row jumps 13px when an order completes" is decided by looking, but the page could not show a screenshot.
The agent captured screenshots of both commits, built a second standalone HTML page, opened it separately, and wrote "see section N of that page" on every card.
The person flipped between two pages to decide one finding.

More widely, tools such as lavish-axi let an agent draw anything a browser can show: charts, Mermaid diagrams, interactive comparisons, real screenshots side by side.
An agora page cannot show any of that, so a design or a review that needs a picture is limited to what the allowlist keeps.

## Solution

An illustration becomes a whole HTML document, written by the agent as an ordinary file, and drawn on its card in a sandboxed frame, per ADR 76.

- A decision names it as `"illustration": {"path": "q3.html"}`.
  A relative path resolves against the document's folder, or against the working directory where the document arrives on standard input.
- At serve time agora copies the HTML file and every relative file it references into the session, within 16 MiB and a fixed set of web media types, and serves them from a route of the session.
  The page never reads the caller's tree after that.
- The card holds an iframe sandboxed with `allow-scripts` and without `allow-same-origin`, so the drawing runs its own script in an opaque origin and cannot reach the page, its token or its controls.
- A drawing may load scripts, styles and fonts from cdnjs, jsDelivr and unpkg, and sends nothing anywhere, per ADR 77.
- The frame grows to its content, up to the window's height, and an expand control lays it over the whole window.
- agora serves its own tokens and fonts as one stylesheet a drawing may link, so a drawing can read as part of the page.
- The drawing only shows: it cannot write into the answer.

The fragment, the document stylesheet, `agora/sanitise` and `agora/stylesheet` are removed.
Both skills teach the new illustration, and socrates review's skill also teaches a UI finding drawn as a before/after pair.

## User Stories

1. As a person deciding a UI finding, I want to see the screen before the change and after it side by side on the finding's own card, so that I decide by looking rather than by reading a description.
2. As a person deciding a UI finding, I want the screenshots to be the real captures, not markup imitating them, so that what I judge is what a user would see.
3. As a person in a design interview, I want a card to show a real layout, a chart or an interactive prototype, so that I can choose between shapes rather than between descriptions of shapes.
4. As a person answering a round, I want every picture on the page I already have, so that I never flip between two pages to decide one question.
5. As a person answering a round, I want a drawing to be as tall as it needs to be, up to my window, so that a short drawing leaves no gap and a tall one does not hide the card behind it.
6. As a person looking at a detailed drawing, I want to expand it to fill the window and close it again with Escape, so that a dashboard or a full screen is readable even on a narrow window.
7. As a person answering a round, I want nothing a drawing does to change my selections, my words or the presses, so that what I send is only what I chose.
8. As a person whose drawings hold private code or unreleased screens, I want nothing a drawing holds to leave my machine, so that a library loaded from a CDN cannot send my work anywhere.
9. As a person reading a drawing, I want it to look like the page around it unless it is showing my project's own screen, so that the page reads as one thing.
10. As an agent writing a round, I want to write a drawing as an ordinary HTML file with ordinary images beside it, so that I do not escape a page into JSON or base64-encode a screenshot.
11. As an agent writing a round, I want to name a drawing relative to the document I wrote beside it, so that `round.json` and `q3.html` in one folder just work.
12. As an agent writing a round, I want to use Mermaid, D3, a chart library or Tailwind from a CDN, so that I do not inline a library into every drawing.
13. As an agent writing a round, I want a reference outside the drawing's folder, a file too large, or a file of a type the surface does not serve to be refused before anything is served, naming the file, so that I fix it rather than serving a broken drawing.
14. As an agent revising a round, I want each revision's drawings copied again from what is on disk now, so that a drawing I changed is the one the person sees.
15. As an agent drawing a UI finding in a review, I want the review skill to tell me to capture before and after screenshots where I have a browser, and to fall back to words where I have none, so that issue 77's workflow is the expected one rather than an improvisation.
16. As an agent with a drawing that only states something a sentence carries in full, I want the skills to keep telling me to skip it, so that a person's attention is not spent on decoration.
17. As a maintainer, I want one way to draw and no allowlist to keep current, so that the surface does not grow a second rendering pipeline.
18. As a maintainer, I want the frame's isolation proven in a real browser, so that a drawing reading the page or sending data out fails the suite rather than a person.

## Implementation Decisions

### The document

- `agora/document`: `Decision.Fragment` and `Document.Stylesheet` are removed, with `StylesheetCeiling`.
  `Decision.Illustration` is added, as `"illustration": {"path": "..."}`, per Question 10.
  A document still naming `fragment` or `stylesheet` is refused by the decoder's unknown-field rule, as any other field it does not have.
- A document's reader is told the folder relative paths resolve against: the document file's folder for `--document <file>`, the working directory for `--document -`, per Question 11.

### Copying an illustration

- At `serve`, for each decision carrying an illustration, the surface copies the HTML file into the session's directory, under a revision of its own, and follows its references: `src`, `href` and `srcset` in the HTML, and `url(...)` and `@import` in any CSS it reaches, inline or copied.
  A reference that is absolute, names a scheme, or is a fragment alone is not copied.
  A relative reference that leaves the HTML file's folder is refused.
- An illustration's files total at most 16 MiB, and each is one of HTML, CSS, JavaScript, JSON, SVG, PNG, JPEG, WebP, GIF, WOFF or WOFF2, per Question 12.
  Anything else, a missing file, or a total over the cap is refused before the document is written or served, naming the decision and the file.
- A revision copies afresh, and the copy of the revision before it is removed once the new one is served.
- What was settled carries no illustration, as it carries no options.

### Serving it

- A new read-only route under the session's path serves an illustration's copied files by decision and relative path, with each file's stated media type and `X-Content-Type-Options: nosniff`.
- Every response of that route carries a `Content-Security-Policy` header: `sandbox allow-scripts`, and the policy of ADR 77.
  That policy admits script, style and font from the illustration's own route, from the surface's own origin named explicitly, and from `https://cdnjs.cloudflare.com`, `https://cdn.jsdelivr.net` and `https://unpkg.com`; image and media from the illustration's route and `data:`; and nothing for `connect-src`, `frame-src`, `form-action` and `base-uri`.
- The page's own policy, in `agora/render`, gains `frame-src` for the illustration route and drops the stylesheet hash, since no caller stylesheet is embedded.

### Sizing and expanding

- The surface adds one script of its own, served from the embedded assets, to the end of each copied HTML document's body.
  It posts the document's height to the parent whenever it changes, per Question 8.
- The page listens for a message whose source is that card's frame and whose data is a finite non-negative number, sets the frame's height to it, capped at the window's height, and ignores anything else.
  This is the only thing the page reads from a frame, so a drawing still cannot write into an answer, per Question 7.
- Each frame has an expand control that lays the frame over the whole window, with a close button, Escape and the backdrop all closing it, and the cursor returning to the control, per Question 9.
  The frame is laid over the page by style rather than moved in the document, since moving an iframe reloads it and loses its state.
- `.illustration`'s `max-height: 32rem`, `.painted` and the containment rule in `agora/web/src/surface.css` go, with `render.NamePrefix`.

### The look

- The surface serves the page's colour tokens, type scale and fonts as one stylesheet under its embedded assets, per Question 13, with `Access-Control-Allow-Origin: *` on it and its fonts, since a sandboxed frame fetches them from an opaque origin.
  The skills name its address.

### Removals

- `agora/sanitise` and `agora/stylesheet` are removed, with their entries in `internal/convention/ignorance.go` and `internal/convention/evidence.go`.
- `docs/ui-design/components/evidence/Illustration.*` is redrawn for the frame and its expand control.

### Skills and documents

- `design/skill/messages.go`: the illustration section describes the file, the frame, the CDNs, the token stylesheet and the refusals, keeps the rule to draw only where the choice is a shape, and a prototype's result can now be the prototype itself.
  "reaches nowhere beyond it" becomes "sends nothing beyond it", per ADR 77.
- `review/skill/messages.go`: `theIllustration` is rewritten the same way, and adds that a finding decided by looking carries a before/after drawing, with screenshots the agent captures by whatever browser it has, falling back to words where it has none, per Question 15.
  socrates review itself captures nothing.
  Its network promise changes as design's does.
- `agora/AGENTS.md`: the invariants about the fragment's allowlist, the stylesheet's parse and hash, and painting are replaced by the frame's sandbox, its policy, the copy's bounds and the one number a page reads from a frame.
  "Every asset is embedded, no diagram engine is" stays true of the page's own bundle.
- `SECURITY.md` names what a drawing can reach: three CDN hosts for loading, nothing for sending.
- `make skill` regenerates the committed skills.

## Testing Decisions

The change is tested at two seams, per Question 14.

`internal/e2e`, with the real binary and a real browser, replaces `painting_test.go` with tests named for these scenarios:

- a drawing renders on its card with a copied image and its own script having run;
- a drawing's script cannot read the page's document, its token or its storage;
- a drawing's script cannot fetch, open a socket, submit a form, or navigate the page or a popup;
- a drawing loads a script from an admitted CDN origin and is refused one from an origin not admitted, with a test server standing in on each, so no test reaches the real network;
- a frame grows to a short drawing's height and stops at the window's height for a tall one;
- a message from a frame that is not a number changes nothing on the page;
- the expand control lays the frame over the window, Escape closes it, the cursor returns to the control, and the drawing's state survives both;
- a drawing linking the token stylesheet is painted with the page's own fonts;
- a revision serves the drawing as it now is on disk.

`agora/surface`, with no browser, tests resolving and copying:

- a relative path resolves against the document's folder, and against the working directory on standard input;
- a reference outside the HTML file's folder, a missing file, a file of a type not served, and a total over 16 MiB are each refused before anything is served, naming the decision and the file;
- references in `src`, `srcset`, `href`, CSS `url(...)` and `@import` are all copied, and absolute and scheme references are not;
- every response of the illustration route carries the sandbox, the policy and `nosniff`.

A guard each of these relies on is mutated once on purpose, such as removing the sandbox header or adding `allow-same-origin`, and the browser test is seen to fail.

## Out of Scope

- A built-in before/after image field or layout, as issue 77 proposed, per Question 6: a pair is one more thing a drawing lays out.
- Any channel from a drawing into the answer, such as clicking an element to write a note, per Question 7.
- socrates review capturing screenshots itself: its skill tells the agent to capture where it can.
- A second loopback port or any separate listener for drawings.
- Embedding any diagram engine or library in the binary.

## Further Notes

- lavish-axi was the reference for the technique: it serves an agent's whole HTML file unchanged in an iframe sandboxed without `allow-same-origin`, repeats the sandbox as a header on every artifact response, and serves relative assets from the file's own folder.
  It sets no network rule, which ADR 77 does.
- A CSP `'self'` inside a sandboxed document is matched against the URL's origin in current browsers, but the policy names the surface's origin explicitly so the token stylesheet and the height script do not depend on that.
- Issue 77 is closed by this change.
