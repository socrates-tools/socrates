## MODIFIED Requirements

### Requirement: The page is held to a stated bar, stated once

The bar the page is held to SHALL be stated in exactly one place in the source and rendered into the skill from it.

It SHALL require each of the following, because each is a way the page fails a person who is deciding something:

- Get to the point. Length is not thoroughness, and a person who has to search the page for their decision has been handed the terminal output again.
- Speak plain English about what the finding means for the change, rather than restating the reviewer's technical prose in a nicer typeface.
- Say what each option gets the person, in the words of the thing being built rather than the words of the code. A sentence they could not act on without opening the code is not plain English, and a symbol, a constant or a test named as though they already know it is that sentence.
- Where a finding is purely technical, illustrate it with an example rather than leaving the reader to work out what the term costs them.
- Where what a finding costs is what a person sees, show them it rather than describing it.
- Use a diagram only where a workflow is what makes the finding hard to follow, and not otherwise.
- Spend tokens as though they were the person's, because they are.

Two copies of a bar drift, and the direction is predictable: the copy nobody reads while working stays as it was, and the copy in front of the session gets softened.

Plain English was already required, and a page shipped with an option meaning "the window endpoints clamp with PAGE_LIMIT again, and the smoke assertions on the 25 cap come out".
It passed because a session that has just read the findings believes the reviewer's vocabulary is plain.
The bar therefore states a test the writer can fail rather than a standard they grade themselves against.

Showing is a separate line from the example beside it and from the diagram beneath it, and the reason is that the bar's other lines are about a reader who has to be taught something.
A finding about what a person sees has a reader who already knows the thing and is being asked to recognise it, and a sketch of a panel retyped in a monospace block is the one form that reader cannot recognise it in.
The line about diagrams says "and not otherwise", so without a line of its own a page that showed a person their own screen would be a page in breach of the bar.

#### Scenario: The bar is stated once and rendered from there

- **WHEN** the bar the skill states is compared with its source
- **THEN** they are the same text, and the source is the only place it appears

#### Scenario: The bar names each way the page fails a person

- **WHEN** the bar is rendered
- **THEN** it requires getting to the point, plain English about impact, an option said in the words of the thing being built rather than the code, an example where a finding is purely technical, showing what a finding costs where that is what a person sees, a diagram only where a workflow needs one, and restraint about tokens

#### Scenario: Showing and drawing are separate lines

- **WHEN** the bar is rendered
- **THEN** the line about what a person sees and the line about diagrams each appear once and neither restates the other
