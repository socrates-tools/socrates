# agora's primitives are Lit elements styled with Tailwind

Every primitive of agora's page becomes a Lit element rendering into a shadow root, styled with Tailwind utilities that paint only the tokens.
Go keeps rendering the page, writing each primitive's element, its attributes and what it slots, and the glossary's pieces and the landing page are laid out with the same utilities.
A person answering a page sees and does what they do today.

Spec: [../spec.md](../spec.md). Decisions: [ADR 82](../../../../agora/docs/adr/0082-a-primitive-is-a-lit-element-that-renders-its-own-parts.md), [ADR 83](../../../../agora/docs/adr/0083-a-primitive-renders-into-a-shadow-root-and-is-reached-through-its-own-properties.md), [ADR 84](../../../../agora/docs/adr/0084-tailwind-paints-only-the-tokens.md), [ADR 85](../../../../agora/docs/adr/0085-the-glossarys-pieces-stay-go-markup-laid-out-by-utilities.md), [ADR 86](../../../../agora/docs/adr/0086-a-decisions-choices-are-one-group-owned-by-ui-choice-group.md).

## Working through the tickets

- Commit this folder before handing out any ticket.
- A ticket may start once every ticket it names as a blocker is done.
- Mark a ticket in progress before handing it out.
- Mark it done once its acceptance criteria hold and its work is committed.
- Mark it abandoned, with a reason, where building it turns out to be wrong, and in the same change every ticket it blocked takes its blockers in its place.

## Tickets

| Ticket | Blocked by | Status |
| --- | --- | --- |
| [1.1: The foundation, proved on the badge](1.1-the-foundation-and-the-badge.md) | None | done |
| [1.2: No class string carries an arbitrary value](1.2-no-arbitrary-value.md) | None | done |
| [2.1: The indicators, the brand and the tooltip as Lit elements](2.1-the-indicators-and-the-tooltip.md) | 1.1 | done |
| [2.2: The button, the field, the copy control and the expand control as Lit elements](2.2-the-controls.md) | 1.1 | done |
| [2.3: A decision's choices as one ui-choice-group](2.3-the-choice-group.md) | 1.1 | done |
| [2.4: The card and the dialog as Lit elements that slot what they wrap](2.4-the-card-and-the-dialog.md) | 1.1 | done |
| [2.5: The app shell as a Lit element](2.5-the-app-shell.md) | 1.1 | done |
| [2.6: The code block and the diff block as Lit elements](2.6-the-code-and-diff-blocks.md) | 1.1 | done |
| [2.7: The glossary's pieces and the catalog laid out by utilities](2.7-the-pieces-and-the-catalog-on-utilities.md) | 1.1 | in progress |
| [2.8: The landing page styled with utilities](2.8-the-landing-page-on-utilities.md) | 1.1 | in progress |
| [3: Every primitive a Lit element, and the light-DOM primitive gone](3-every-primitive-a-lit-element.md) | 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7 | pending |
