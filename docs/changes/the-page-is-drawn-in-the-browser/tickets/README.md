# The agora page is drawn in the browser

This change moves the drawing of agora's page from Go templates to the client's TypeScript.
Go keeps validating the document and serves it as JSON with what only the server knows, and the client draws the whole page from it.
To a person answering, the page looks and behaves as it did.

Spec: [../spec.md](../spec.md). Decisions: [ADR 96](../../../../agora/docs/adr/0096-the-browser-draws-the-whole-page-from-the-document-go-serves.md), [ADR 97](../../../../agora/docs/adr/0097-a-piece-is-a-template-one-draw-function-writes-into-the-light-dom.md).

## Working through the tickets

- Commit this folder before handing out any ticket.
- A ticket may start once every ticket it names as a blocker is done.
- Mark a ticket in progress before handing it out.
- Mark it done once its acceptance criteria hold and its work is committed.
- Mark it abandoned, with a reason, where building it turns out to be wrong, and in the same change every ticket it blocked takes its blockers in its place.

## Tickets

| Ticket | Blocked by | Status |
| --- | --- | --- |
| [1.1-the-document-the-page-fetches](1.1-the-document-the-page-fetches.md) | None | done |
| [1.2-the-wording-check-reads-the-clients-words](1.2-the-wording-check-reads-the-clients-words.md) | None | done |
| [1.3-the-catalog-drawn-by-the-client](1.3-the-catalog-drawn-by-the-client.md) | None | done |
| [1.4-the-browser-suite-holds-the-page](1.4-the-browser-suite-holds-the-page.md) | None | done |
| [2-the-page-drawn-in-the-browser](2-the-page-drawn-in-the-browser.md) | [1.1-the-document-the-page-fetches](1.1-the-document-the-page-fetches.md), [1.2-the-wording-check-reads-the-clients-words](1.2-the-wording-check-reads-the-clients-words.md), [1.4-the-browser-suite-holds-the-page](1.4-the-browser-suite-holds-the-page.md) | done |
| [3-go-writes-no-page-markup](3-go-writes-no-page-markup.md) | [1.3-the-catalog-drawn-by-the-client](1.3-the-catalog-drawn-by-the-client.md), [2-the-page-drawn-in-the-browser](2-the-page-drawn-in-the-browser.md) | pending |
