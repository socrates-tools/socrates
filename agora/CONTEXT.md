# Context: agora

`socrates agora` serves one document as one page on a loopback address, waits for a person, and prints what they answered.
These are agora's words.
agora knows nothing about what a caller's document is about, beyond the three audiences ADR 47 names and the reading of the three severities its hints carry.
`docs/adr/` here holds agora's decisions and `AGENTS.md` its invariants.

## The page

- **document**: what a caller writes for the surface to render. Not a page.
- **source**: the entries a caller cites, each carrying a name, a severity, an audience and a location.
- **decision**: one card on the page, citing the source entry it came from.
- **index**: the list beside the decisions naming every source entry, or every decision where there is no source, with whether each is answered and a count of how many are.
- **mark**: the icon a severity or an audience is drawn as, saying what it means on hover, on focus and to a screen reader.
- **settled answer**: what a document carries of a question already answered elsewhere: its question, the labels chosen and the words. Shown ahead of the decisions, read rather than answered, and never sent.
- **option**: one answer to a decision. It carries the command it becomes, or none, and a decision does not mix the two: a person choosing between the options of one decision is choosing between things of the same kind. The surface composes a command and never runs one. The control that takes none of them is not an option: it carries no command, and leaves the decision with nothing selected.
- **illustration**: the one fragment a decision may carry, painted by the document's own stylesheet, contained so nothing in it reaches the page's own controls.
- **session**: what the surface mints so that a page and its routes are served under a path holding it.
- **exchange**: the delivery of a document and the collection of what a person answered. It holds at most one uncollected answer, and a later one replaces it. The last answer collected is kept in the session's directory, and `agora poll --last` prints it again.
- **answer**: what one submit carries: every decision the person answered, anything they wrote, and whether it ends the exchange.
- **confirmation**: the dialog a press raises where any decision is blank, naming each one and taking the person to it. The only thing on the page standing between a press and a send.
