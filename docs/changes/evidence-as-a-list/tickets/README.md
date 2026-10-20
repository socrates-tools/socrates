# Evidence is a list of captioned blocks, and design shows code through it

A decision's evidence becomes a list of blocks, each code in a language or a patch, with an optional caption.
Every code block wraps its long lines and offers a copy control, and a language may be named by any alias highlight.js knows.
The design skill then teaches evidence for code, so a design session stops drawing code by hand in an illustration.

Spec: [../spec.md](../spec.md). Decisions: [ADR 78](../../../../agora/docs/adr/0078-evidence-is-a-list-of-blocks-and-the-only-shape-read.md).

## Working through the tickets

- Commit this folder before handing out any ticket.
- A ticket may start once every ticket it names as a blocker is done.
- Mark a ticket in progress before handing it out.
- Mark it done once its acceptance criteria hold and its work is committed.
- Mark it abandoned, with a reason, where building it turns out to be wrong, and in the same change every ticket it blocked takes its blockers in its place.

## Tickets

| Ticket | Blocked by | Status |
| --- | --- | --- |
| [1.1: Evidence is a list of captioned blocks](1.1-evidence-as-a-list.md) | None | pending |
| [1.2: Every code block wraps its long lines](1.2-every-code-block-wraps.md) | None | pending |
| [1.3: Copy a code block](1.3-copy-a-code-block.md) | None | pending |
| [1.4: A language may be named by any name highlight.js answers to](1.4-language-aliases.md) | None | pending |
| [2: The design skill shows code through evidence](2-design-skill-teaches-evidence.md) | 1.1-evidence-as-a-list, 1.4-language-aliases | pending |
