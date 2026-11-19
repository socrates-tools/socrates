# The agora page is drawn in the browser

Decisions: ADR 96, ADR 97.

## Problem Statement

The agora page is drawn in two languages.
Go writes the pieces the glossary names, such as the decision, the composer, the receipt and the index, in `agora/render/page.gohtml` and 18 `ui-*.gohtml` partials.
TypeScript defines the primitives those partials write, as Lit elements in `agora/web/packages/ui`.
The two are held together by `contract.json`, which `make assets` writes from what each element declares, and by `agora/render/contract_test.go`, which checks every partial against it.

Every new primitive costs a Lit file, a partial and an entry in `catalog.go`, and every change to what a primitive takes is made twice.
The page's words live in Go, in `agora/render/messages.go`, while the code that shows and hides them lives in TypeScript.
Reading a patch into diff rows is written twice, in `agora/render/diff.go` and in `ui/diff-block/patch.ts`.
A revision fetches a whole page of HTML and parses it to lift one element out.
`agora/render` has been among the busiest folders since October, and most of that churn is keeping the two halves in step.

## Solution

TypeScript draws the whole page in the browser, per ADR 96.
Go serves a shell holding the address, the token, the bundle and one line saying the page needs JavaScript.
The page fetches the document Go validated, as the caller wrote it, together with what only the server knows: the revision, whether the exchange is over, each illustration's address and what is answered.
It fetches it the same way when it opens and on every revision.

One function, `draw(page)`, returns a new `.surface` for each of those, and each piece is a private template in it, written into the light DOM with the one class that names it, per ADR 97.
The page's script and the browser suite keep finding pieces as they do today.

The TypeScript types for what crosses are generated from Go's structs and committed, and `make test` fails while they disagree.
The page's words move to a TypeScript messages module, held by the same wording check as every `messages.go`.
The catalog becomes a TypeScript page in the bundle, still served at `/catalog` under the `agoradev` build tag.

## User Stories

1. As a contributor adding a primitive, I want to write one Lit element and its catalog states in TypeScript, so that I no longer write a partial and a Go catalog entry and keep three files in step.
2. As a contributor changing what a primitive takes, I want the type checker to tell me every place that draws it, so that a renamed attribute fails the build rather than a contract test after the fact.
3. As a contributor changing how a decision looks, I want its markup, its classes and its behaviour in one language, so that I read one file and not a template, a view struct and a script.
4. As a contributor changing what a document holds, I want the TypeScript types to be generated from the Go structs, so that the page cannot read a field Go does not send.
5. As a contributor reviewing a change to the page, I want the page's words beside the code that draws them, so that a label and its use are reviewed together.
6. As a contributor, I want the wording check to keep holding every word the page shows, wherever it now lives, so that no em dash and no word that names nothing a person can do reaches the page.
7. As a contributor, I want reading a patch into diff rows to be written once, so that the heading's counts and the rows never disagree.
8. As a contributor writing a browser test, I want the classes I find pieces by to stay as they are, so that moving the drawing does not mean rewriting the suite.
9. As a person answering a design round, I want the page to look and behave exactly as it does today, so that the change is invisible to me.
10. As a person answering, I want a revision to replace what I see with the new document and carry nothing across, as it does today.
11. As a person whose browser runs no script, I want the page to say it needs JavaScript, so that I am not shown questions I cannot answer.
12. As a caller of `socrates agora serve`, I want the document I write and every refusal of it to stay exactly as they are, so that no skill or caller changes.
13. As a contributor proving a primitive's behaviour, I want the catalog to draw every primitive in every state with the same code the page uses, so that the browser suite proves what ships.
14. As a contributor who never touches the client, I want `make verify` to keep running without Node, so that the toolchain I need does not grow.
15. As a contributor changing a Go struct that crosses to the page, I want the suite to tell me the generated types and the bundle are stale, so that I know to run `make assets`.

## Implementation Decisions

**What Go serves.**
`agora/surface` serves the shell on the page's route and the page's JSON on a route under the same session, guarded by the session's token as the page is today.
The JSON is the validated document as `agora/document` holds it, beside a server part: the revision, whether the exchange is over, each decision's illustration address, what is answered, and the tray ceiling of ten files.
The stream on `/sessions/<id>/live` is unchanged and still carries no token, so it only says a revision exists, per ADR 93.

**What leaves Go.**
`page.gohtml`, `catalog.gohtml`, `partials/`, `partials.go`, `primitives.go`, `laid.go`, `diff.go`, `marks.go`, `catalog.go` and `contract_test.go` are removed, with the tests that read the HTML they wrote.
`agora/render` shrinks to the shell and its policy, or folds into `agora/surface` if nothing else is left in it.
The audiences socrates review writes, stated in `audiences.go`, move to the TypeScript messages module, and the wording check's exemption names them there.
`make assets` stops writing `contract.json`.

**Generated types.**
A Go program run by `make assets` writes the TypeScript types for the document and the server part from the Go structs, to a committed file under `packages/surface`.
A Go test regenerates them in memory and fails while the committed file differs, naming the file, so the check needs no Node.
A changed struct also leaves the bundle stale, and the existing check that the bundle agrees with its source catches that.

**The drawing.**
`packages/surface` gains `drawing.ts`, exporting `draw(page): HTMLElement` and the generated `Page` type.
Each piece is a private function returning a Lit template built from the library's primitives, writing the same classes, ids and `data-` attributes the Go markup writes today.
On opening, `main.ts` fetches the page's JSON, mounts `draw(page)`, then starts every `followThe...` module as it does now.
`revision.ts` fetches the JSON instead of the HTML and replaces the surface with `draw(page)`, without `DOMParser`.
Marks, the ordering of the index by severity and the diff heading's file names and counts are worked out in TypeScript, the counts from `patch.ts`.

**Words.**
The page's words move from `agora/render/messages.go` to one messages module in `packages/surface`.
The wording check learns to read the plain string literals of that module as it reads every `messages.go`.

**Tailwind.**
`theme.css`'s `@source` lines stop naming the partials and keep naming the client's TypeScript and the landing page.
The check on arbitrary values keeps reading those lines.

**The catalog.**
The states `catalog.go` states move to a TypeScript module beside the library, drawn by the same primitives, and the catalog page is part of the bundle.
`agora/surface` keeps routing `/catalog` only under `agoradev`, now serving a shell that loads the catalog.
The check that every primitive is drawn by the catalog reads the library's declarations rather than `contract.json`.

**Documents to update in the change that lands it.**
`agora/AGENTS.md`'s layout and invariants for `agora/render`, `agora/web`, the contract and the catalog, the root `AGENTS.md` where it names them, and `CONTRIBUTING.md` where it says who needs Node.

## Testing Decisions

What the page draws is tested only in the browser suite, `internal/e2e`, against the real binary, as the suite's own invariant already asks.
Its selectors stay as they are, because each piece keeps its class.
The harness gains one shared wait for the drawn surface where a test waited for the document to load.

Go tests stop at the JSON Go hands over: `agora/surface` tests read the page's JSON route and assert the document and the server part, including the illustration addresses and the token guard on the route.
The shell is tested for its policy, the token and the address it carries, and the line for a page without script.
The generated-types check is a Go test, and it is broken on purpose once, by editing the committed file, before it is trusted.

No TypeScript test runner is added.
The 8 Go test files that read rendered HTML go with the code they tested, and each behaviour they held that the browser suite does not already hold moves into it before they are removed.

## Out of Scope

- Rewriting any of socrates outside the page in TypeScript.
- Changing what a caller writes in a document, or any refusal of one.
- Changing the stream, its events or its guards.
- Making pieces Lit elements, or moving the page's state out of the DOM into Lit's re-rendering.
- Adding a TypeScript test runner.
- Changing how the page looks or behaves to a person.

## Further Notes

The direction was chosen on a comparison page before the interview, against keeping the page as it is and against rewriting the whole suite in TypeScript.
That page found that about 90% of the shipped Go has nothing to do with HTML, and that the one place two languages met was the partials and the contract test this change removes.

`draw(page)` replacing the whole surface keeps today's rule that nothing carries across a revision, but it also keeps the DOM as the only place an answer is held.
Moving to Lit's in-place re-rendering later would mean rewriting the surface modules, and ADR 97 records why that was not chosen now.
