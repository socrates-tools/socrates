# Review page UI kit

The one browser surface socrates has: the decisions page served by agora (`socrates agora`, used by `socrates review` and `socrates design`).

This kit shows the **proposed light redesign** (Aegean + clay on limestone). The page as shipped today is dark (warm black + gold, ADR 0045); its recreation is `/Agora Current.dc.html`.

Source: `agora/render/page.gohtml`, `agora/web/src/surface.css`, `agora/web/src/surface.js`, `agora/render/{icons,marks,audiences,messages}.go`, sample content `agora/document/testdata/source.json`.

Files
- `index.html` — interactive page; switcher bottom-right toggles open / over / nothing to answer, and the settled panel.
- `ReviewPage.jsx` — page shell: top bar, title, rail + column, confirmation dialog.
- `Rail.jsx` — progress and index entries.
- `SettledPanel.jsx` — read-only answers from earlier rounds.
- `DecisionCard.jsx` — one decision: attribution, question, evidence, answering band.
- `Composer.jsx` — message, send / finish, receipt, amend.
- `data.js` — sample content.

Behaviour copied from surface.js: a recommended option arrives checked; a decision counts as answered when a choice is made or words are written; sending with blanks opens the confirmation; "Send and finish" closes the exchange.
