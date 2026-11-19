# 2: The page is drawn in the browser

## What to build

The page a person opens is drawn by the client.
Go serves a shell holding the address, the token, the bundle and one line saying the page needs JavaScript, and nothing the document says.
The page fetches its document when it opens and on every revision, the same way, and draws the whole surface from it, each piece keeping the one class it is found by today.
Marks, the index's order and the diff heading are worked out in the client, the heading's counts from the same reading of the patch the rows come from.
Every word the page shows comes from the client's messages module.
To a person, the page looks and behaves as it did.

## Acceptance criteria

- [ ] The browser suite passes with its selectors unchanged, beside one shared wait for the drawn page.
- [ ] The HTML the page's route serves holds none of the document's questions, options or words.
- [ ] With script turned off, the page shows one line saying it needs JavaScript and nothing of the document.
- [ ] A revision draws the new document from what it fetched and carries nothing across: no selection, no word, no receipt and no open confirmation.
- [ ] Renaming a property a primitive declares makes the asset build fail, naming where the page draws it.
- [ ] Every word the page shows is a string of the client's messages module, so an em dash added to any of them fails the suite.

## User stories

- 2, 3, 5, 7, 8, 9, 10, 11, 12

## Left to other tickets

- 3-go-writes-no-page-markup.md: removing the Go page, its partials, its words and the contract it no longer uses

## Blocked by

- 1.1-the-document-the-page-fetches.md
- 1.2-the-wording-check-reads-the-clients-words.md
- 1.4-the-browser-suite-holds-the-page.md
