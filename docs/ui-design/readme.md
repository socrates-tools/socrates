# Socrates Design System

The visual language for **socrates**, a command-line tool that puts an independent reviewer in front of changes nobody has read yet.
The tool is mostly terminal; its one visual surface is the **decisions page** that `socrates agora` serves in the browser.
A person reads findings there and sends answers back to the coding session.

This system is the **light redesign** of that page: Aegean blue and terracotta clay on whitewashed limestone.
It replaced the dark palette (warm black + gold) that ADR 45 had chosen, and ADR 53, *"The page is light whatever the reader prefers"*, records the change.

> Naming: **socrates** is the product brand and goes on every surface.
> **agora** is the internal web server and is never shown to the person.

## Sources
- Codebase: `socrates/` (local folder, Go + a small `agora/web` bundle).
  Key files:
  - `agora/render/page.gohtml`: page structure.
  - `agora/web/src/surface.css`: shipped tokens and styles (dark).
  - `agora/web/src/surface.js`: answering and sending behaviour.
  - `agora/render/icons.go`, `marks.go`, `audiences.go`, `messages.go`: icons, marks, tooltip copy.
  - `agora/document/testdata/source.json`: sample content.
  - `agora/web/fonts/`: Cinzel, Manrope, JetBrains Mono (variable woff2), copied to `fonts/`.
  - `docs/adr/0053-…`: the decision that adopted this system, superseding ADR 45.
- `Agora Current.dc.html`: recreation of the dark page this system replaced.
  `Agora Redesign.dc.html`: the redesign mock this system was extracted from.

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
- **Type:** Manrope for everything, with headings at 700 and tight negative tracking (40px title, 22px question).
  JetBrains Mono for code, file locations and evidence headings.
  Cinzel appears only in the wordmark; it is an all-capitals face, so "socrates" renders as small capitals.
  **12px is the floor** for any text.
- **Spacing:** 4px base.
  The rhythm is 8 / 12 / 16 / 24 / 32 / 48.
  Cards are 32px padded and 24px apart; options are 8px apart; the rail and column sit 48px apart.
  Everything left-aligns to the 32px gutter shared by the top bar, title and rail.
- **Layout:** 1200px page and a 64px sticky top bar, wordmark left and status right.
  Below it, a full-width title and preamble.
  Then a 264px sticky rail beside a column of up to 840px.
  The two wrap to one column when narrow.
- **Backgrounds:** flat limestone.
  No gradients, textures, photos or illustrations.
  Cards sit on the ground with a hairline border.
- **Cards:** white, 1px `#E5DED1` border, 16px radius, soft two-layer navy-tinted shadow.
  The answering band inside a decision is *sunk* (`#FBF9F5`) and divided by a hairline.
  The settled panel is sunk with no shadow.
- **Radii:** 6 (checkbox), 10 (wordmark tile), 12 (options, fields, evidence, notices, rail rows), 16 (cards), 20 (dialog), pill (buttons, chips, status).
- **Selected state:** Aegean border, a 1px Aegean ring drawn as box-shadow so nothing shifts, and an `#EEF4FB` wash.
  Radio fills to a 6px ring; checkbox fills solid with a white check.
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
  The `Icon` component reproduces them from those exact paths.
- Icons appear only inside **marks** (severity/audience) and **checks**.
  There are no decorative icons, no icon font, no emoji and no unicode-as-icon (except ·, →, − in counts and diffs).
- If more icons are needed, take them from Lucide at the same 2px stroke.
- **Logo:** none is supplied.
  The wordmark is set type: Cinzel 700, 0.12em tracking, lowercase "socrates", with an optional Aegean tile holding the initial.

## Index
- `styles.css`: entry point (imports only), which pulls in `tokens/`: fonts, colors, typography, spacing, effects and base.
- `fonts/`: Cinzel, Manrope and JetBrains Mono variable woff2.
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
- `SKILL.md`: Agent Skill entry.

### Intentional additions
- **Icon**: wraps the icon paths that `icons.go` emits, so components don't repeat SVG markup.
- **TopBar**: the redesign moves the wordmark and status into a sticky bar; the shipped page has them in a masthead row.
- **Tag**: one primitive for the *Recommended* mark and the *Settled* overline, which the shipped page styles separately.
