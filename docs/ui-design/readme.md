# Socrates Design System

The visual language for **socrates**, a command-line tool that puts an independent reviewer in front of changes nobody has read yet.
The tool is mostly terminal, and its one visual surface is the **decisions page** that `socrates agora` serves in the browser.
A person reads findings there and sends answers back to the coding session.

This is the system that page is painted in: Aegean blue and terracotta clay on whitewashed limestone.
ADR 53, *"The page is light whatever the reader prefers"*, adopted it and superseded the dark page of ADR 45.
Where this folder and the shipped page disagree, the shipped page is right: `agora/web/src/surface.css` holds the tokens the page actually paints, and this folder is the reference a change to the page is designed against.

> Naming: **socrates** is the product brand and goes on every surface. **agora** is the internal web server and is never shown to the person.

## Sources

- `agora/render/page.gohtml`: page structure.
- `agora/web/src/surface.css`: the shipped tokens and styles.
- `agora/web/src/surface.js`: answering and sending behaviour.
- `agora/render/icons.go`, `marks.go`, `audiences.go`, `messages.go`: icons, marks, tooltip copy.
- `agora/document/testdata/source.json`: sample content.
- `agora/web/fonts/`: Cinzel, Manrope and JetBrains Mono, which `fonts/` here copies byte for byte.
- `Agora Redesign.dc.html`: the mock this system was extracted from.

---

## Content fundamentals

- **Plain, literal, unhurried.**
  Copy says exactly what happens and nothing more: "Sending waits for a reply on this page. Finishing sends and closes the exchange."
  There is no marketing voice.
- **Sentence case** everywhere, including buttons ("Send and wait for a reply").
  Session-written titles and questions arrive **lowercase** ("two things to settle before this lands") and are shown verbatim.
- **Person:** the page speaks to *you* ("Your call") and names the other side as **"whoever asked"**, never "the AI" or "the agent".
  "The author" is the session that wrote the change.
  "The reviewer" is the independent model.
- **Domain words:** finding, decision, exchange, settled, round, change.
  Use them as the glossary does.
  Don't swap in synonyms ("issue", "ticket", "conversation").
- **Consequences over labels:** every option carries a *meaning* line ("correctness over the extra wait").
  Blank states explain their result: "A decision left blank goes as nothing. Whoever asked will see no answer to it."
- **Overlines** are the only uppercase text (DECISIONS, STILL OPEN).
- **No emoji.**
  No exclamation marks.
  Curly quotes for quoted answers (“…”).
  Middle dot (·) for counts ("Settled earlier · 1").

## Visual foundations

- **Palette:** Greek coast.
  **Aegean** `#1259A8` is the one accent: what you choose, what you press, links, focus.
  **Clay** `#B4532D`/`#A3471F` is used sparingly, only for *Recommended* and *Your call*.
  **Limestone** neutrals (`#F6F2EA` page, `#FFFFFF` cards, `#FBF9F5` sunk bands) carry the ground.
  **Ink** navy `#14283C` is all body text.
  States: olive (resolved), cinnabar (error), ochre (warning), always as ink on a pale wash.
- **Type:** Manrope for everything at standard web-app sizes: 14px body, 13px labels and buttons, 12px micro and code, 18px question, 28px page title.
  Headings are 700 with slight negative tracking.
  JetBrains Mono for code, file locations and evidence headings.
  Cinzel appears only in the wordmark; it is an all-capitals face, so "socrates" renders as small capitals.
  **12px is the floor** for any text.
- **Spacing:** 4px base.
  The rhythm is 4 / 8 / 12 / 16 / 24 / 32.
  Cards are 24px padded and 16px apart; options are 8px apart and 12×16 padded; the rail and column sit 32px apart.
  Controls are 36px tall.
  Everything left-aligns to the 24px gutter shared by the top bar, title and rail.
- **Layout:** 1200px page and a 56px sticky top bar, wordmark left and status right.
  Below it, a full-width title and preamble.
  Then a 240px sticky rail beside a column of up to 800px.
  The two wrap to one column when narrow.
- **Backgrounds:** flat limestone.
  No gradients, textures, photos or illustrations.
  Cards sit on the ground with a hairline border.
- **Cards:** white, 1px `#E5DED1` border, 12px radius, soft two-layer navy-tinted shadow.
  The answering band inside a decision is *sunk* (`#FBF9F5`) and divided by a hairline.
  The settled panel is sunk with no shadow.
- **Radii:** 6 (checkbox), 10 (wordmark tile, rail rows), 12 (cards, options, fields, evidence, notices), 16 (dialog), pill (buttons, chips, status).
- **Selected state:** Aegean border, a 1px Aegean ring drawn as box-shadow so nothing shifts, and an `#EEF4FB` wash.
  Radio fills to a 5px ring; checkbox fills solid with a white check.
- **Hover:** surfaces get a limestone wash (rail rows) or a darker border (fields, options).
  Primary darkens to `#0E4A8C`.
  Secondary turns its border and text Aegean.
- **Press:** darkest Aegean `#0B3D73` plus `scale(0.98)`.
  No bounces.
- **Focus:** 3px `#D6E5F5` ring, with an Aegean border on fields.
- **Motion:** 160ms colour/border transitions and 220ms for the progress bar, on `cubic-bezier(0.2,0,0,1)`.
  Disabled under reduced motion.
- **Transparency & blur:** only the top bar (85% limestone, 12px blur) and the dialog scrim (36% navy, 4px blur).
- **Shadows:** card `0 1px 2px / 0 12px 32px -12px` at 4–10% navy.
  Dialog is deeper.
  The primary button has a soft Aegean glow.
- **Closed states hide, not disable.**
  When the exchange is over, the text fields and the composer are removed, and a single neutral notice under the preamble says so.

## Iconography

- The set is the six Lucide stroke icons that socrates ships in `agora/render/icons.go`: circle-alert, triangle-alert, info, user, wrench and eye.
  There is also a square-capped check from `page.gohtml`.
  They are inline SVG, 2px stroke, round caps and single colour (currentColor).
  The `Icon` component reproduces them from those exact paths, so components don't repeat SVG markup.
- Icons appear only inside **marks** (severity/audience) and **checks**.
  There are no decorative icons, no icon font, no emoji and no unicode-as-icon (except ·, →, − in counts and diffs).
- If more icons are needed, take them from Lucide at the same 2px stroke.
- **Mark:** `assets/brand/socrates-mark.svg` (and `favicon.svg`): a simplified front-facing Socrates bust (bald dome, furrowed brow, snub nose, curled beard) in limestone on an Aegean rounded tile, with a clay chiton at the shoulders.
  Flat shapes, no gradients.
  Use it as the favicon and app icon; never recolour or redraw it.
  Link it with `<link rel="icon" type="image/svg+xml" href="assets/brand/favicon.svg">`.
- **Wordmark:** set type: Cinzel 700, 0.12em tracking, lowercase "socrates".
  In the top bar it is led by the mark at 28px with a 10px gap.

## Index

- `styles.css`: entry point (imports only), which pulls in `tokens/`: fonts, colors, typography, spacing, effects and base.
- `fonts/`: Cinzel, Manrope and JetBrains Mono variable woff2, each under the SIL Open Font License 1.1, whose text and copyright notices are in `NOTICE` at the repository root.
- `guidelines/`: 15 foundation specimen cards (Colors, Type, Brand, Spacing, Effects).
- `components/`:
  - `brand/`: Wordmark, Icon.
  - `actions/`: Button.
  - `marks/`: StatusPill, Mark, Check, Tag.
  - `forms/`: Option, OptionGroup, TextArea, Progress.
  - `evidence/`: CodeBlock, Tok, DiffView, Illustration.
  - `surfaces/`: Card, CardSection, Notice.
  - `navigation/`: TopBar, IndexEntry.
  - `feedback/`: Dialog, Tooltip.
- `ui_kits/review-page/`: the interactive decisions page (open / sent / confirm / over / nothing to answer).
- `assets/dev/preview-fallback.js`: lets specimen cards and the kit render from the component sources before the bundle is compiled.
