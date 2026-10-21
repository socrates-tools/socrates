# Review page UI kit

The one browser surface socrates has: the decisions page served by agora (`socrates agora`, used by `socrates review` and `socrates design`).

This kit shows the decisions page as ADR 53 adopted it and ADR 70 amended it: Aegean and sea on limestone.

Source: `agora/render/page.gohtml`, `agora/web/packages/surface/surface.css`, `agora/web/packages/surface/`, `agora/render/{icons,marks,audiences,messages}.go`, sample content `agora/document/testdata/source.json`.

Layout: an app shell that uses the whole window.
- A full-width top bar (`TopBar fluid`): wordmark, the exchange title as a breadcrumb, status on the right.
- A left navigation (`sidebar-width`, 272px) on the sunk surface: progress, the open decisions with the one in view marked `current`, and the settled-earlier link. It scrolls on its own.
- The main area: a page header (counts, title, preamble), then the settled panel and the decision cards at full width. What was seen (question, evidence) gets the card's whole width so code and diffs have room; the answering band sits below it.
- From 1280px, a right reply panel (`panel-width`, 304px): an answer summary at the top (it scrolls), and the message box and the two send buttons, side by side ("Send and finish", then the primary "Send and wait"), pinned at the bottom, always in reach. Below 1280px the composer closes the decisions column instead.
- Below 960px everything stacks into one scrolling column, navigation first.

Files
- `index.html` - interactive page; the "Kit state" switcher at the bottom left is a demo control, not part of the page: it toggles open / over / nothing to answer, and the settled panel.
- `layout.css` - the app shell grid and its breakpoints.
- `ReviewPage.jsx` - app shell: top bar, navigation, main area, reply panel, confirmation dialog.
- `Rail.jsx` - the left navigation: progress and index entries.
- `SettledPanel.jsx` - read-only answers from earlier rounds.
- `DecisionCard.jsx` - one decision: attribution, question, evidence, answering band.
- `Composer.jsx` - message, send / finish, receipt, amend; as the right reply panel or inline.
- `data.js` - sample content.

Behaviour copied from surface.js: a recommended option arrives checked; a decision counts as answered when a choice is made or words are written; sending with blanks opens the confirmation; "Send and finish" closes the exchange.
