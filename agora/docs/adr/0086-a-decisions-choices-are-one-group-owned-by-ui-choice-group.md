# 86. A decision's choices are one group owned by ui-choice-group

- Status: accepted
- Date: 2026-10-21
- Part of: lit-and-tailwind-primitives

## Decision

A decision's options and its *None of these* are `<ui-choice>` elements slotted into one `<ui-choice-group>`.
The group keeps at most one choice checked in a single-choice decision, roves the tab stop and the arrow keys between its choices, and carries `role="radiogroup"`, with each choice carrying `role="radio"`.
A multiple-choice decision uses the same group, each choice carrying `role="checkbox"`, with no single selection to keep.

Native radios sharing a name are one group only inside one tree, and each `ui-choice` renders into a shadow root of its own, per ADR 83.
Radios in different shadow roots are never one group, so the browser stops keeping one selected and stops moving between them on the arrow keys.
The group is what owns that behaviour once the browser does not, which is how Shoelace and Web Awesome build radios.
The next engineer will see a radio group written by hand and be tempted to put native radios back, and this record is why they are not there.

## Rejected

- Each choice form-associated through `ElementInternals`, unchecking every choice of the same name in the document when it is checked.
  The card said keyboard movement between choices would then be code of each choice, and it was declined.
- `ui-choice` alone keeping the light DOM, so a native radio groups as it does today.
  The card called it the one exception to ADR 83's rule, with Tailwind's sheet also styling it from the document, and it was declined.
