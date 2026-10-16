# 1: An illustration is drawn in a sandboxed frame

## What to build

A decision's illustration is `"illustration": {"path": "..."}`, naming an HTML file the agent wrote, in place of the fragment.
`socrates agora serve` copies that file into the session and the card draws it in an iframe sandboxed with `allow-scripts` and without `allow-same-origin`, served from a read-only route of the session.
A relative path resolves against the document's folder, or against the working directory where the document arrives on standard input.
Every response of the route carries the sandbox and the policy of ADR 77 as a header, with `nosniff`, and the page's own policy admits the frame.
A revision copies the file afresh, and the copy of the revision before it goes once the new one is served.

The fragment and the document stylesheet go in the same change: a document naming either is refused as naming a field it does not have.
The sanitiser, the stylesheet parser, the page's painting rules, the name prefix, the stylesheet's hash in the page's policy and the conventions that named them go with them.
agora's invariants describe the frame, its sandbox and its policy in place of the allowlist and the scoped stylesheet.
Both skills stop teaching the fragment and say only that an illustration is an HTML file named by path, drawn in a sandboxed frame; ticket 3 teaches the rest.
This ticket copies the HTML file alone: the files it references come in 2.1.

## Acceptance criteria

- [ ] A document naming an illustration by a path relative to its own folder is served, and the card shows the drawing with its own inline script having run.
- [ ] A document piped on standard input resolves the path against the working directory.
- [ ] A path that does not exist, or names a file that is not HTML, is refused before anything is served, naming the decision and the path.
- [ ] A document carrying `fragment` or `stylesheet` is refused before anything is served, naming the field.
- [ ] A drawing's script cannot read the page's document, the session token or the page's storage.
- [ ] A drawing's script cannot fetch, open a socket, submit a form, or navigate the page or open a popup.
- [ ] A drawing loads a script from an admitted CDN origin, and one from an origin not admitted is refused, with no test reaching the real network.
- [ ] Opening the illustration route on its own still runs the drawing sandboxed, because the sandbox arrives as a header.
- [ ] Editing the HTML file and serving the same document again shows the edited drawing.
- [ ] No package of agora reduces caller markup or caller CSS, and neither skill nor any committed rendering mentions the fragment, markup, the document stylesheet or the allowlist.

## User stories

- 3, 4, 7, 8, 10, 11, 12, 14, 17, 18

## Blocked by

None
