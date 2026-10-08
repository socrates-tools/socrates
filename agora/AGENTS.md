# socrates agora

`socrates agora` is a namespace inside the one binary, not a second program.
It serves one document as one page on a loopback address, waits for a person, and prints what they answered.
It knows nothing about code review but the three audiences socrates review writes and the reading of the three severities its hints carry, per ADR 47.

This file is agora's layout and invariants.
`/AGENTS.md` holds the suite's, `CONTEXT.md` here holds agora's vocabulary, and `docs/adr/` here holds agora's decisions.

## Layout

- `agora/command` - the command tree: `serve`, `poll`, `end`, `stop`, and nothing else.
- `agora/config` - the one optional configuration file agora reads, which holds the `surface` block and nothing of review's.
- `agora/document` - what a caller writes for the surface to render, what it copies rather than accepts, and the languages the page can highlight.
- `agora/sanitise` - a caller's HTML fragment reduced to a fixed allowlist, and every name it carries prefixed.
- `agora/stylesheet` - a caller's stylesheet read as a parse, reduced to a fixed allowlist, and scoped to the illustration.
- `agora/render` - a document to one HTML page, the marks a severity and an audience are drawn as, and a patch to the rows of a split diff.
- `agora/assets` - the built client bundle and fonts, embedded, and the stamp of their source.
- `agora/surface` - the loopback server and its guards, the one address bound where nothing states one, the sessions it holds, what it writes about itself and about each page, the exchange, and the bounds on waiting.
- `agora/web/` - the client source, the design system and the fonts, built by `make assets`.

## Invariants

Each is held by a test that fails when it stops holding.
None is a matter of care, and the reasoning for each is in this directory's `docs/adr/`.

**The surface**

- It binds a loopback address; nothing widens it, and an address that is not one is refused before any server starts.
- It binds one fixed port where nothing states one, and that port held by anything else is a refusal naming the port, `surface.bind` and `--address`, and no surface starts on another port.
  The port is stated once, in `agora/surface` beside the guards that judge an address.
- `serve`, `poll`, `end` and `stop` find the surface by asking the port who stands there, never by `surface.json` alone: the surface answers `/identity` with no token, that it is agora, the root it writes under and the build it was linked from, and a surface is trusted only when it answers as agora under this root.
  Where a committed bind names port 0 there is no port to ask, so the address `surface.json` records is asked instead and trusted on the same terms.
  A port held by an agora surface of another root is a refusal naming that root.
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
- The order is the surface's, so no field of a document reaches it, and a page with nothing to answer offers neither control and says so.
- Every decision carries a free-text field that no document field can remove, including one that offers no options.
- The option a document marks arrives chosen, a decision that marks none arrives with nothing chosen, and a decision offering options also offers taking none of them, which is not an option, carries no command and leaves the decision with nothing chosen.
- What the controls hold at the submit is what is sent, so nothing is committed per decision, a decision holding no selection and no words is omitted, and a decision left blank is named before either press sends.
  Either press over a blank decision raises a modal dialog naming every one of them and takes no armed state, a press inside it is what sends, what is behind it is dimmed, and each of its names closes it and puts the cursor on that card's first control.
  Going back, Escape and a press on the backdrop are one path: nothing is sent, the page is as it was, and the cursor returns to the button that raised it.
- The languages the page highlights are the bundled library's common set, stated once and generated from `web/`, and a document naming one outside it is refused while one naming none renders as plain text.
- A caller's fragment is reduced to a fixed allowlist before it is embedded, and every name it carries is prefixed by the surface.
- A caller's stylesheet is reduced from a parse that reached the end of its input, every selector it keeps is scoped to the illustration, and the page admits it by a hash of exactly the bytes it embedded.
- What a caller paints is contained by the illustration, so a rule inside one decides nothing about where the page's own controls land.
- Delivery happens once, waiting is bounded, collection is per session, and the server stops once it has held no session for the idle period.
- An exchange holds at most one uncollected answer, and a later submit replaces it rather than queueing behind it, so what a caller collects is a person's last word.
- Collecting an answer is visible to the page, so the way back from a receipt is offered exactly while it would still reach the caller.
- Nothing carries across a revision: not a selection, not a typed word, not the message, not the receipt and not an open confirmation, because a decision's identifier is unique only inside the document that names it.
- A stream is told where the exchange stands when it connects rather than only when something next changes, so a client that missed a revision catches up rather than waiting on a change that has already happened.
- The end of an exchange stops the page offering to answer and not the page listening, so a caller who writes again reaches the page the person is already looking at, and only the test about reloading reloads before it asserts.
  A closed exchange says so under the title and draws no field nothing will read, per ADR 53.
- Every asset is embedded, no diagram engine is, and the committed bundle is generated from `web/`.
- An element it hides paints nothing, every colour it paints is a token stated once, and every text meets a contrast floor under both preferences.
  The page is light whatever the reader prefers, so both preferences are painted alike, per ADR 53.
- A severity the surface orders by is marked by what it means, each audience socrates review writes is marked as the design drew it and any other in the caller's own words, and every mark says what it means on hover, on focus and to a screen reader, per ADR 47.
  A mark on a card also says its name in a word, and a mark in the index is its icon alone, per ADR 53.
- Nothing in it names review, finding, verdict, round, or run as a noun, and the module path is not part of its vocabulary, except the three audiences socrates review writes, stated in `agora/render/audiences.go` and exempted by name, and the severity hints, which carry review's reading of a severity in words the check does not know, per ADR 47.

**The suite**

- The page half is driven by a real browser against the real binary; asserting on served markup would test the renderer.
