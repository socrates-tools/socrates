# 83. A primitive renders into a shadow root and is reached only through its own properties and events

- Status: accepted
- Date: 2026-10-21
- Part of: lit-and-tailwind-primitives

## Decision

Every primitive renders into a shadow root and composes what Go writes inside it by slots.
A card, a dialog and the app shell wrap content Go writes, and a slot is how a Lit element wraps content without code of its own to move it.
The tokens reach inside as inherited custom properties, and one constructable Tailwind sheet, built once, is adopted by every primitive.
The same stylesheet is also loaded on the document, so Tailwind's `@property` registrations hold and the pieces Go writes are styled.

Nothing outside a primitive opens its shadow root.
The surface's script and the browser suite reach a primitive through what it exposes, such as `value`, `checked`, `disabled`, `focus()` and `click()`, and through composed events it dispatches.
A class the script or the suite finds a control by sits on the primitive's element, so `.send` is the `<ui-button>` rather than a `<button>` inside it.
A delegated listener on the document sees an event from inside a shadow root retargeted to the primitive, so the surface's listeners stay attached once, on the document.

This supersedes ADR 81, which drew every primitive without a shadow root so that the page's stylesheet and the suite's selectors reached inside it.
The test asserting that no primitive has a shadow root becomes one asserting that every primitive has one.

## Rejected

- The light DOM, as ADR 81 said.
  The card said every wrapping primitive would then carry its own code to move the content Go wrote into its frame, which is the imperative DOM code this change is for removing, and it was declined.
- The light DOM for leaves and a shadow root for wrappers.
  The card called it two rules for one library, and it was declined.
- Open shadow roots that the script and the suite query into through one piercing helper.
  The card said callers would then depend on each primitive's internal markup, and it was declined.
