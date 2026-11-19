# 96. The browser draws the whole page from the document Go serves

- Status: accepted
- Date: 2026-11-19
- Part of: the-page-is-drawn-in-the-browser

## Decision

TypeScript draws the whole page in the browser, and Go writes none of it but a shell.
The shell holds the address, the token, the bundle and one line saying the page needs JavaScript, and nothing a document says.
The page fetches the document Go validated, as the caller wrote it, together with what only the server knows: the revision, whether the exchange is over, each illustration's address and what is answered.
It fetches it the same way on opening and on every revision, so there is one way in.
Marks, diff rows and layout are worked out in TypeScript, and `agora/render`'s partials, `page.gohtml`, `catalog.gohtml`, `diff.go`, `marks.go` and the contract test that held Go's markup to what each element declares are removed.

Go keeps validating a document, and its structs stay the source of the shape that crosses.
The TypeScript types are generated from them and committed, and `make test` fails while the two disagree.

This supersedes ADR 85, whose pieces stayed Go markup, and the half of ADR 82 that kept Go writing a primitive's element through a partial checked against `contract.json`.
The next engineer will find a Go server that writes no HTML for its one page, and expect the server to render it as it once did.
It does not, because a page drawn in two languages was held together by a contract test, and every primitive cost a Lit file, a partial and a catalog entry.

Changing what crosses now needs Node, because the bundle that reads it is rebuilt from TypeScript, where ADR 79 asked Node only of somebody who touches `agora/web`.

## Rejected

- Keeping the page as it is, Go templates and Lit joined by a contract test, and rewriting the whole suite in TypeScript.
  Both were cards on the page that started this change, and the person chose to have TypeScript draw the page instead.
- Go sending a prepared view it has already worked out, with TypeScript only drawing it.
  The card's cost was that the shape TypeScript reads would be the view rather than the document a caller writes, and patch splitting would stay written twice.
- A read-only list of the questions Go writes for a page whose script does not run.
  The card's cost was a second way to draw a decision, kept in Go, for a page on which nothing can be sent without the script.
- Writing the first document into the shell as a JSON block.
  The card's cost was two ways in, since a revision must fetch the document anyway: the stream carries no token, per ADR 93.
