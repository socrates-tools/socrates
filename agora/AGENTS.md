# socrates agora

`socrates agora` is a namespace inside the one binary, not a second program.
It serves one document as one page on a loopback address, waits for a person, and prints what they answered.
It knows nothing about code review but the three audiences socrates review writes and the reading of the three severities its hints carry, per ADR 47.

This file is agora's layout and invariants.
`/AGENTS.md` holds the suite's, `CONTEXT.md` here holds agora's vocabulary, and `docs/adr/` here holds agora's decisions.

## Layout

- `agora/command` - the command tree: `serve`, `poll` (which prints the last collected answer again with `--last`), `end`, `stop`, and nothing else.
- `agora/config` - the one optional configuration file agora reads, which holds the `surface` block and nothing of review's.
- `agora/document` - what a caller writes for the surface to render, what it copies rather than accepts, and the languages the page can highlight.
- `agora/render` - a document to one HTML page, the marks a severity and an audience are drawn as, a patch to the rows of a split diff, and the partials that write each primitive's markup, one `ui-<name>.gohtml` under `partials/` for each, and the catalog, whose states for every primitive are stated once in `catalog.go`.
- `agora/assets` - the built client bundle and fonts, embedded, and the stamp of their source.
  A file under the `agoradev` build tag serves the bundle from disk instead, for `make dev`.
- `agora/surface` - the loopback server and its guards, the one address bound where nothing states one, the sessions it holds, what it writes about itself and about each page, the copy of each illustration and the route that serves it, the exchange, and the bounds on waiting.
- `agora/web/` - the client source, the design system and the fonts, built by `make assets`.
  It is a pnpm workspace of two packages.
  `packages/ui` is `@agora/ui`, the library: its `tokens/` holds `tokens.css`, `fonts.css`, the fonts and the mark, and the bundle serves `tokens.css` at `/assets/tokens.css` with every custom property it holds, because illustrations load it by that path.
  Its `theme.css` is the one Tailwind entry: it clears Tailwind's own theme, names the tokens as utilities, declares the named variants and screens the markup writes, and states with its `@source` lines every path a class is read from, which are the client's TypeScript, the partials and the landing page.
  `make assets` compiles it once, bundles the result into `surface.css` for the document, and hands it to `stylesheet.ts`, which builds the one constructable sheet every primitive adopts as its `static styles` through `UiElement` in `element.ts`.
  A primitive is one TypeScript file in its own folder, exported from `index.ts`: a `@customElement` class extending `UiElement` that paints its host by returning the utilities it wants from `hostClasses()` and draws its parts into its shadow root, with a module of its own where it needs one, such as the code block's highlighting or the diff block's patch.
  No stylesheet styles a primitive, and the library holds none but `theme.css` and what `tokens/` holds.
  Adding a primitive is a Lit element in its own folder exported from `index.ts`, a partial writing the element with the attributes it declares and what it slots, and its states in `catalogued` in `agora/render/catalog.go`, and nothing else.
  The library imports nothing from the surface, it is where highlight.js and Lit are dependencies, and a stylesheet outside `packages/ui/tokens/` that writes a colour or names a font family fails the suite.
  A class written with an arbitrary value, `p-[13px]` or `bg-[#fff]` or `bg-(--x)`, in any file `theme.css`'s `@source` lines name fails the suite naming the file and the line, and the check reads those lines rather than a list of its own, per ADR 84.
  `packages/surface` holds the page's script as strict TypeScript and the surface's stylesheet, which holds only the rules no utility can say.
  The pieces the glossary names are written in `page.gohtml`, `catalog.gohtml` and the partials, composed of the library's primitives and laid out by utilities, each keeping the one class that names it for the script and the suite, per ADR 85.
  The script is one module for each of sending and the receipt, the confirmation, the live revision, the calls that show and hide the tooltip and the navigation, beside the answers, the composer, the expanding drawing, copying and the drawings' heights.
  Each follows the page through a listener on the document attached once when the page starts, so a revision fetches the page and replaces the surface and binds nothing.

## Invariants

Each is held by a test that fails when it stops holding.
None is a matter of care, and the reasoning for each is in this directory's `docs/adr/`.

**The surface**

- It binds a loopback address; nothing widens it, and an address that is not one is refused before any server starts.
- It binds one fixed port where nothing states one, and that port held by anything else is a refusal naming the port, `surface.bind` and `--address`, and no surface starts on another port.
  The port is stated once, in `agora/surface` beside the guards that judge an address, and every binary the suite builds is linked with a port of its own in its place, so no test reaches the port a real surface stands on, per ADR 74.
- `serve`, `poll`, `end` and `stop` find the surface by asking the port who stands there, never by `surface.json` alone: the surface answers `/identity` with no token, that it is agora, the root it writes under and the build it was linked from, and a surface is trusted only when it answers as agora under this root.
  Where a committed bind names port 0 there is no port to ask, so the address `surface.json` records is asked instead and trusted on the same terms.
  A port held by an agora surface of another root is a refusal naming that root.
- A surface of this root linked from another build, by any difference in the version, the commit stamp or whether it was built uncommitted, is reused while it holds a session still being answered and replaced once every session it holds is over, per ADR 58.
  Replacing stops it through its control token on a route that answers only when no session is open and refuses new ones from then on, waits a bounded time for the port, and starts this build's surface on it, and a port that is never let go is the refusal for a port held by something else.
  A `serve` naming a session never replaces the surface holding it, and nothing without a token learns whether the surface holds a session.
- A port an operator committed to, by flag or by either field of the `surface` block, is bound or refused and never walked past.
- The surface mints a session, `--session` never brings one into existence, and every page and every route of it is served under a path holding that session.
- It turns away a request naming a host it did not bind, and every state-changing route carries a token: a session's own answers for that session and no other, and the one minted when the server started answers everywhere.
- It composes a command and never runs one, the command it hands back is the option's own string, and no command is rendered on the page.
- It writes only under agora's own directory in the brand's root, beside agora's configuration file, and nothing of it lands in another tool's directory.
- What an operator states about it is read from `~/.socrates/agora/config.yaml`, in agora's own directory under the brand's root, and from no other file, and a file it cannot read or a field the `surface` block does not have is refused, naming the file.
- A caller writes a document rather than a page, and the name, severity, audience and location come from the source a decision cites.
- A document may hold no decision, and one holding no decision, no source entry, nothing settled and no preamble is refused.
- A decision may carry no options, and such a card states something and offers no control a person could mistake for a choice; one carrying no options and stating a mode is refused.
- Every entry of the source is rendered ahead of the decisions, in an index beside them that names every decision where there is no source, and both are ordered by the severity the source carries, most severe first, informational last, an unstated severity after all of them, and ties in the order the document listed them.
- The index counts a decision answered by the rule that decides what is sent, as the person answers it, so a decision whose recommendation arrives chosen arrives counted.
- What was settled is rendered ahead of every open question, offers no control, is counted by nothing and is never sent, and an identifier is either settled or still open, per ADR 46.
- What was settled arrives folded while a question is still open, and unfolded when nothing is, per ADR 72.
- The order is the surface's, so no field of a document reaches it, and a page with nothing to answer offers neither control and says so.
- From 960px the page fills the window and only its navigation and its decisions scroll, the navigation marks the decision in view, and from 1280px a reply panel summarises what will go and keeps the message and both presses in reach, per ADR 75.
  There is one composer, and below 1280px it closes the decisions.
- Every decision carries a free-text field that no document field can remove, including one that offers no options.
- The option a document marks arrives chosen, a decision that marks none arrives with nothing chosen, and a decision offering options also offers taking none of them, which is not an option, carries no command and leaves the decision with nothing chosen.
- What the controls hold at the submit is what is sent, so nothing is committed per decision, a decision holding no selection and no words is omitted, and a decision left blank is named before either press sends.
  Either press over a blank decision raises a modal dialog naming every one of them and takes no armed state, a press inside it is what sends, what is behind it is dimmed, and each of its names closes it and puts the cursor on that card's first control.
  Going back, Escape and a press on the backdrop are one path: nothing is sent, the page is as it was, and the cursor returns to the button that raised it.
- The languages the page highlights are the bundled library's common set and every alias it registers for them, stated once and generated from `web/`, and a document naming one outside them is refused while one naming none renders as plain text.
- Every code block of evidence wraps at the width of its card, whatever its language, and nothing in one scrolls sideways.
  A continuation line starts further in than the source line it continues, measured against that line's own indentation, so a wrapped line still reads as one line.
  The split diff of a patch is a table and keeps its own scroll.
- Every code block of evidence numbers its source lines from 1 in a gutter, a wrapped line once, and the numbers are never part of its text, and every block is painted the white of the card it sits on.
- Every code block of evidence offers a copy control that puts the code on the clipboard exactly as the document stated it, says what it does on hover, on focus and to a screen reader, confirms the copy where it was pressed, and sends nothing and changes nothing about the answer.
  A patch's split diff offers none.
- An illustration is an HTML file a caller names by path, drawn in an iframe sandboxed with `allow-scripts` and without `allow-same-origin`, so its script runs in an opaque origin and reads neither the page, its token nor its storage, and cannot navigate the page or open a popup, per ADR 76.
- The sandbox is stated again as a `Content-Security-Policy` header on every response of the illustration's route, beside `X-Content-Type-Options: nosniff`, so a drawing opened on its own is sandboxed as well.
- That policy admits script, style and font from the route, the surface's own origin and three CDN hosts, admits image and media from the route and `data:`, and names nothing for connecting, framing, forms or the base, so a drawing may load a library and has no way to send what it holds anywhere, per ADR 77.
- The page's own policy admits the frame and embeds nothing a caller wrote, so no caller markup and no caller stylesheet is reduced, scoped or hashed anywhere in agora.
- A relative illustration path resolves against the document's folder, or against the working directory where the document arrives on standard input, and a path that is missing or names no HTML file is refused before anything is written or served, naming the decision and the path.
- The surface copies the file into the session rather than reading the caller's tree after that, serves it from a read-only route by decision and file name, and a revision copies it afresh and removes the copy before it once the new one stands.
- The copy brings every relative `src`, `href` and `srcset` of the HTML and every relative `url(...)` and `@import` of any CSS it reaches, and leaves an absolute reference, one naming a scheme and a fragment alone as they are.
  A reference leaving the HTML file's folder, a file that is missing or of a type outside HTML, CSS, JavaScript, JSON, SVG, PNG, JPEG, WebP, GIF, WOFF and WOFF2, and files totalling over 16 MiB are each refused before anything is written or served, naming the decision and the file.
- The page reads one thing from a frame: a finite non-negative number posted by that card's own frame window, which sets that frame's height, capped at the window's height.
  The page asks each frame for it once with a fixed word, so a post that came before the page was listening is not lost.
  Anything else posted to the page, whether a value of another type, a number from a window that is not that frame, or a message from the page itself, changes nothing, so a drawing has no channel into an answer, per ADR 76.
- Delivery happens once, waiting is bounded, collection is per session, and the server stops once it has held no session for the idle period.
- An exchange holds at most one uncollected answer, and a later submit replaces it rather than queueing behind it, so what a caller collects is a person's last word.
- Every collected answer is kept whole in its session's directory as `answers/<n>/answer.json`, numbered in the order collected, before it is handed back, and one the session could not keep stays uncollected, so `agora poll --last` prints again the last answer collected, without waiting and without the surface, per ADR 88.
- Nothing removes a collected answer but the session's directory going when the surface collects the session, so every path a poll printed stays readable until then, per ADR 88.
- Either press sends one multipart POST holding the answer as JSON in the part `answer` and every file every tray holds, each in a part of its own that the answer names in `attached`, and what the controls hold at the press is exactly what is sent, so an amend sends every file again, per ADR 87.
  The surface refuses the whole submit, keeping nothing, for a body that is not form data, for a name the answer states that no part carries and for a part the answer does not name.
  It refuses it too for a file over 20 MiB, for files totalling over 100 MiB and for an eleventh file in one tray, each ceiling stated once in `agora/surface` beside the 1 MiB `BodyCeiling` that keeps guarding the document and control routes, and the page says which ceiling it was and that nothing was sent.
- Every file is written into a fresh folder of the session before the exchange swaps its uncollected answer, the swap removes the folder of the answer it replaces, and a revision removes the folder of the uncollected answer, so no uncollected answer's files outlive the answer.
- A stored file is named by agora as its position and a sanitised base name, such as `answers/<n>/<decision>/1-before.png` or `answers/<n>/message/1-trace.txt`, and no stated name reaches outside the answer's folder, whatever separators or dots it holds.
- `agora poll` and `agora poll --last` print an `attachments` list of name, type, size and absolute path beside a decision's words and beside the message, omit it where it is empty so an answer with no attachment prints as it always did, and print the same paths wherever the answer is read from.
- The surface serves no attachment and no route reads back what a person sent, so the page previews a file from memory, its policy admits `blob:` for images and for nothing more, and a person's bytes are never rendered from the surface's origin, per ADR 87.
- Every decision's words and the composer's message carry a tray, a revision empties every tray with the rest of the controls, a tray is read-only while the answer is sent and a closed exchange draws no tray.
- Collecting an answer is visible to the page, so the way back from a receipt is offered exactly while it would still reach the caller.
- Nothing carries across a revision: not a selection, not a typed word, not the message, not the receipt and not an open confirmation, because a decision's identifier is unique only inside the document that names it.
- A revision settling an identifier the exchange collected an answer for states the question, the labels, the words and the attachment names that were collected, exactly, or is refused before it is written or served, per ADR 59.
  Only a collected answer binds, the last one collected for an identifier is the one that does, and what the exchange remembers is never sent to the page.
- A stream is told where the exchange stands when it connects rather than only when something next changes, so a client that missed a revision catches up rather than waiting on a change that has already happened.
- The end of an exchange stops the page offering to answer and not the page listening, so a caller who writes again reaches the page the person is already looking at, and only the test about reloading reloads before it asserts.
  A closed exchange says so under the title and draws no field nothing will read, per ADR 53.
- Every asset is embedded, no diagram engine is, and the committed bundle is generated from `web/`.
- A released binary serves its assets from its embedded bundle and from nowhere else, and only the `agoradev` build serves them from the bundle on disk, so a change to the client is seen with a reload in `make dev` and nothing on a person's machine can change what the page they answer looks like or does.
- Every primitive is a Lit element rendering into an open shadow root that adopts the one shared Tailwind sheet, per ADR 82 and ADR 83, and a `ui-` element drawn without a shadow root fails the browser suite on the catalog and on the page, naming it.
  Nothing outside `packages/ui` names `shadowRoot`, and `make assets` refuses a file that does, naming the file and the line, so the surface's script and the browser suite reach a primitive through its properties, methods and events and never inside it.
  A primitive's parts are inside its shadow root and nothing names them from outside, so a partial that writes `data-part` fails the Go contract test, and `data-interaction` on the primitive's element says which of answering, sending and viewing a control belongs to.
  A primitive paints its host through classes it sets on itself, so the browser suite reads the host's box and colour, and the sheet that styles the host is the one loaded on the document.
  A `ui-` element the bundle has not yet defined paints nothing, through one `:not(:defined)` rule in the shared sheet.
  The page's policy is unchanged by the adopted sheet, and the browser suite records no violation of it on the page or the catalog.
- The theme is `theme.css`, and Tailwind paints only the tokens: a utility exists for a token or a named step and for nothing else, per ADR 84.
  A state the page carries in a `data-` attribute or in `aria-current` is a utility's variant declared once by name there, a width is a screen named there, and a pseudo-element or a slotted node is a variant, `slotted` for `::slotted(*)`, `host-focused` for `:host(:focus-visible)` and `gutter` for the `::before` that draws a line number.
  A Go view's class is laid out from the template through `laid`, which adds utilities to a class the view names, so Tailwind reads every utility from the `.gohtml` files.
  The surface's stylesheet holds no rule that names a class, an id or an element, no stylesheet under `packages/` does, and the suite fails naming the rule where one does.
  A utility written with `!` is important, and it is written so for two cases: a rule on a slotted node, because the document's Preflight outranks a shadow root's rule on a node it slots, and a rule that overrides what the document says of a host, such as the page's `:focus-visible` outline, which outranks a shadow root's.
  A card styles its slotted body and footer with `slotted:` utilities written important, so a section the page hides paints nothing and a decision whose exchange is over loses its footer and the strip around it together.
- A primitive declares its attributes through its reactive properties and its slots through its template, and what Go writes that a person reads is slotted rather than passed as an attribute.
  The badge takes a tone and an icon and slots its label, the status dot a tone, the overline a count and its words, the progress bar a label, a value and a maximum, the brand a mark and its name, and the tooltip nothing, and a mark is a piece of the surface drawn on `ui-badge`.
  The index is built on `ui-progress-bar`, `ui-status-dot` and `ui-overline`, a decision, the composer and the empty page are drawn on `ui-card`, which slots a body and a footer, the brand in the topbar is `ui-brand`, and the page itself is drawn in `ui-app-shell`, which slots the topbar, the navigation, the content and the panel Go writes, owns where each sits, the scrolling of the navigation and of the content from 960px and the panel's place from 1280px, and marks in the navigation the section in view.
  The exchange's form wraps `ui-app-shell` in the page, so the controls in its slots stay part of the form.
  A domain piece carries its own class on the primitive's element, so the confirmation is `<ui-dialog class="confirmation">`, which draws its dialog in its shadow root and slots the confirmation's heading, body and actions.
  Where the browser suite and the script reach a control by its class, the class sits on the primitive's element, so `.send` is the `<ui-button>`, `.words-field` is the `<ui-field>` and `.choice` is the `<ui-choice>`, and they read and set it through `value`, `checked`, `disabled`, `focus()` and `click()`.
  A control sets `delegatesFocus`, so a click or `focus()` on its element reaches the control inside it, and a submit button asks its form to submit through `requestSubmit` in `packages/ui/form.ts` and is the press `pressOf` names, since a shadow root keeps the control out of the form it stands in.
  A decision's options and its *None of these* are `ui-choice` elements slotted into one `ui-choice-group`, per ADR 86: native radios in different shadow roots are never one group, so the group keeps at most one choice checked where it is single-choice, keeps one tab stop that follows the checked choice, moves the checked choice and the cursor on the arrow keys, and carries `role="radiogroup"`, while a `multiple` group carries `role="group"` and checks several, and each choice carries `role="radio"` or `role="checkbox"` with `aria-checked`.
  A choice announces `change` when a person presses it or presses Space on it, as a native one did, and Enter on it submits the form it is in.
  A focused choice draws the page's focus ring around its indicator, through the `host-focused` variant, and keeps its own chosen ring and no outline, written important.
  The copy control offers what it holds in its `text` attribute, or asks the listeners above it with a `copyrequested` event whose `detail.text` they set, and the expand control announces `expanded` and `collapsed` to its owner, which lays out what it expands.
  The file tray takes files through `add(files)`, which returns the files it refused, and `clear()`, reads them through `files`, holds at most `max`, takes none through a drop or through `add` while `disabled`, returning every file as refused without a refusal drawn or `refused` announced, and announces `change` when what it holds changes and `refused` when a file does not fit, drawing a refusal that names the ceiling where the file was added, a thumbnail previewed from a `blob:` URL for an image and a type badge for anything else, a control removing each file and a counter against the ceiling, and says what it drew through `drawnChips`, `drawnCounter` and `drawnRefusal`.
  A dialog opens modal from `open()`, closes from `close()` and from Escape or a press outside its panel, announces `dismissed` on Escape and on that press, and says whether it stands open and modal in `isOpen` and `isModal`, and the page's script names its heading with `headed()`.
  The tooltip shows on hover and on focus of anything carrying `data-title`, hides when the cursor or the focus leaves and on scroll, and follows a focused anchor when the page scrolls, as methods and listeners of `ui-tooltip` itself.
  The code block takes its code as the text of the element and a `language`, draws the highlighted, numbered and hung lines in its shadow root, and says what it drew through `statedText()`, `highlighted` and `drawnLines`, each line's number, gutter and the column each of its rows starts at, so nothing outside it reads its shadow root to find a line.
  The diff block takes the patch as the text of the element and a `language`, reads it into hunks and paired numbered sides in its own module, draws one scrolling table for each hunk in its shadow root, and says what it drew through `highlighted` and `drawnHunks`, each hunk's heading, rows with their kinds, numbers, text and coloured tokens, and the scroll of its table; Go still reads the patch for the file names and the counts of the heading and to leave out a patch with no hunk.
  Both paint a token by utilities named for the token's colour, which `highlighting.ts` maps from the class highlight.js gives it, so a line's number is never part of the text a person selects.
  The client holds no code that binds listeners to a surface or releases them across a revision, and the suite fails where the page holds a listener of an earlier surface after several revisions in a row.
- The catalog has no document behind it and is drawn by the real partials with every primitive in every state it has, and a primitive the contract declares that the catalog does not draw fails the suite naming it.
  Only the `agoradev` build routes `/catalog`, so a released binary answers it with not found, and the browser suite proves each primitive's behaviour against it.
  It loads the page's own `surface.js` and `surface.css`, so the primitives are proved as the page ships them, and it draws the shell in three iframes of 480px, 960px and 1280px, so each width the shell answers to is a window of its own, and the script starts the page's features only where there is a surface, so the catalog follows no exchange.
  Its opening control names a `ui-dialog` by its id with `data-opens`, and the dialog's own `open()` answers it.
- `make assets` writes `contract.json` from the program it type-checks, and it holds elements and nothing else: for each Lit element, the attributes its reactive properties read and the slots its template renders.
  A Go test in `agora/render` fails on an attribute or a slot a partial writes that a Lit element does not declare, on a declared one no partial writes, and on a `data-part`, naming the element and the attribute, slot or part, so changing a partial's markup needs no Node, per ADR 82.
  The global attributes `class`, `id`, `hidden`, `title`, `role`, `slot`, `tabindex`, `autofocus`, `aria-*` and `data-*` are not part of the contract.
- The client holds no JavaScript but its build script, and `make assets` type-checks it under `strict` before bundling anything, so a type error stops the build naming the file and line and writes no bundle, per ADR 79.
  `make lint` does not run the type checker, so `make verify` needs no Node.
- The client's TypeScript holds no comment, and `make assets` refuses one, naming the file and line, except the `/*!` notice that carries a bundled library's licence.
  `scripts/prove-client-build` runs the build over a scratch copy holding each violation, in the step continuous integration runs after `make verify`.
- An element it hides paints nothing, every colour it paints is a token stated once, and every text meets a contrast floor under both preferences.
  The page is light whatever the reader prefers, so both preferences are painted alike, per ADR 53.
- A severity the surface orders by is marked by what it means, each audience socrates review writes is marked as the design drew it and any other in the caller's own words, and every mark says what it means on hover, on focus and to a screen reader, per ADR 47.
  A mark on a card also says its name in a word, and a mark in the index is its icon alone, per ADR 53.
- Nothing in it names review, finding, verdict, round, or run as a noun, and the module path is not part of its vocabulary, except the three audiences socrates review writes, stated in `agora/render/audiences.go` and exempted by name, and the severity hints, which carry review's reading of a severity in words the check does not know, per ADR 47.

**The suite**

- The page half is driven by a real browser against the real binary; asserting on served markup would test the renderer.
- A revision settling an identifier the exchange collected an answer for states the question, the labels, the words and the attachment names that were collected, exactly, or is refused before it is written or served, per ADR 59.
  Only a collected answer binds, the last one collected for an identifier is the one that does, and what the exchange remembers is never sent to the page.
