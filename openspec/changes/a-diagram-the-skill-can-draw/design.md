## Why the skill may name a field of a document socrates never writes

socrates reviews.
It does not render, serve, or wait for a person, and `agora` knows nothing about code review.
Naming a field of agora's document inside socrates' own skill looks like the coupling that scope rule exists to prevent.

It is not, and the reason is that the coupling is already there and is the point of the skill.
The skill already tells a session to run `agora serve --document`, to hand over the address, and to wait on `agora poll`.
It already tells the session to write that document.
What it withholds is one field of the document it just asked for, which is not independence but an omission that reads as a limit.

The seam that matters is that **socrates reaches agora not at all**: no command, flag or run record names a page, and nothing socrates spawns knows one exists.
That seam is untouched here.
A session reading the skill is on the other side of it, and it is the session that writes both the document and the diagram.

## Why the bar does not name the field

The bar states what the page owes a person.
Each line of it is a way the page fails somebody who is deciding something, and none of them is a fact about the document format.

A bar that named `fragment` would have to change when the field did, and the bar is the one text held to being stated once and rendered from a single source.
Coupling it to a format is how a bar starts drifting toward a manual.

So the bar keeps its line about when a diagram is warranted, the mechanism sits beside it, and the two are read together without either owning the other's job.

## What was rejected

**Teaching the session to reach for a second tool.**
A general-purpose artifact surface will draw a diagram, and a session that believes agora cannot will find one.
It puts the person in front of two pages for a single round, and the second page's controls do not name the commands they become, which is the property the whole exchange rests on.

**Embedding a diagram engine.**
Already refused by `surface-page`, and this change does not reopen it.
A layout engine surrenders position and emphasis to a machine, and an illustration is written by the caller precisely so that the caller decides what the reader looks at first.
