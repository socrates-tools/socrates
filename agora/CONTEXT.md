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
- **mark**: the icon a severity or an audience is drawn as on a badge, saying what it means on hover, on focus and to a screen reader.
- **settled answer**: what a document carries of a question already answered elsewhere: its question, the labels chosen, the words and the names of the files attached. Shown ahead of the decisions, read rather than answered, and never sent.
- **option**: one answer to a decision. It carries the command it becomes, or none, and a decision does not mix the two: a person choosing between the options of one decision is choosing between things of the same kind. The surface composes a command and never runs one. The control that takes none of them is not an option: it carries no command, and leaves the decision with nothing selected.
- **evidence**: the code a decision shows, as a list of blocks, each a piece of code in a language or a patch, with an optional caption. Drawn and highlighted by the page itself, unlike an illustration. Not code, and not a snippet.
- **illustration**: the one HTML document a decision may carry, drawn in a frame of its own that reaches neither the page nor the network beyond the CDNs it names. Not a fragment, and not painted by the document.
- **session**: what the surface mints so that a page and its routes are served under a path holding it.
- **exchange**: the delivery of a document and the collection of what a person answered. It holds at most one uncollected answer, and a later one replaces it. Every answer collected is kept whole in the session's directory until the surface collects the session, and `agora poll --last` prints the last one again.
- **answer**: what one submit carries: every decision the person answered, anything they wrote, and whether it ends the exchange.
- **attachment**: a file a person sends with an answer, beside one decision's words or beside the message. Always the person's and never the caller's, so not evidence and not an illustration. Not an upload.
- **tray**: where a decision or the composer holds its attachments until a press sends them. Not a drop zone.
- **confirmation**: the dialog a press raises where any decision is blank, naming each one and taking the person to it. The only thing on the page standing between a press and a send.
- **composer**: the one place a person sends from: the field for anything to say to whoever asked, the two presses, and the receipt. Not the answers themselves, which stay on their decisions.
- **receipt**: what the composer shows once a press has sent, naming every decision that went and the message beside them. Not a record: it is drawn from the page and kept nowhere.

## The client

- **primitive**: a generic piece of the page that names nothing in agora's domain, such as a button, a field, a dialog or a code block. Every piece of the page this glossary does not name is a primitive, and every piece it does name is built from primitives. Not a component, and not a widget.
- **catalog**: every primitive drawn in every state it has, by the same markup the page uses, for the browser suite to drive. Not a page: no document is behind it and nobody answers it.
