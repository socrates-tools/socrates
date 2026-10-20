# 4: The catalog draws every primitive in every state

## What to build

The catalog is a page drawn by the real partials, holding every primitive in every state it has, with no document behind it.
The test suite serves it to the browser, and the development build serves it so a person can look at it.
A released binary has no route to it.
The dialog is drawn in each of its states, and the browser suite proves its behaviour against the catalog rather than through the confirmation alone.
Each primitive added to the library from here on adds itself to the catalog.

## Acceptance criteria

- [ ] The browser suite drives the catalog and proves the dialog opening, closing on Escape, closing on a press on the backdrop, and returning the cursor.
- [ ] The development build serves the catalog at a route under the surface, and a person can open it.
- [ ] A released binary answers that route with not found.
- [ ] A primitive declared in the contract but absent from the catalog fails the suite, naming the primitive.

## User stories

- 13, 14

## Blocked by

- 3.1-the-dialog-as-the-first-primitive.md
- 2.1-development-build.md
