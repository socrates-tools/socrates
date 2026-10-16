# An illustration is whatever HTML can draw

A decision's illustration becomes an HTML file the agent writes, drawn on its card in a sandboxed frame, in place of the allowlisted fragment and the document stylesheet.
It lets a page show real screenshots, charts and interactive drawings, as issue 77 and the design behind this change ask.

Spec: [../spec.md](../spec.md). Decisions: [ADR 76](../../../../agora/docs/adr/0076-an-illustration-is-an-html-document-drawn-in-a-sandboxed-frame.md), [ADR 77](../../../../agora/docs/adr/0077-an-illustration-may-load-from-three-cdns-and-sends-nothing-out.md).

## Working through the tickets

- Commit this folder before handing out any ticket.
- A ticket may start once every ticket it names as a blocker is done.
- Mark a ticket in progress before handing it out.
- Mark it done once its acceptance criteria hold and its work is committed.
- Mark it abandoned, with a reason, where building it turns out to be wrong, and in the same change every ticket it blocked takes its blockers in its place.

## Tickets

| Ticket | Blocked by | Status |
| --- | --- | --- |
| [1: An illustration is drawn in a sandboxed frame](1-an-illustration-is-drawn-in-a-sandboxed-frame.md) | None | in progress |
| [2.1: An illustration brings the files it references](2.1-an-illustration-brings-the-files-it-references.md) | 1 | pending |
| [2.2: The frame grows to its drawing](2.2-the-frame-grows-to-its-drawing.md) | 1 | pending |
| [2.3: A drawing expands to fill the window](2.3-a-drawing-expands-to-fill-the-window.md) | 1 | pending |
| [2.4: A drawing may link the page's tokens](2.4-a-drawing-may-link-the-page-tokens.md) | 1 | pending |
| [3: Both skills teach the drawn illustration](3-both-skills-teach-the-drawn-illustration.md) | 2.1, 2.4 | pending |
