# 31. The recommendation arrives chosen

- Status: accepted
- Date: 2026-09-16
- Part of: the page becomes a form

## Decision

A document marking an option as the one the caller would take used to change nothing about how the page rendered it: the mark was said, in words, and the option itself arrived exactly like every other one, nothing on the page selected.
A person who agreed with every recommendation on a document of six decisions still had to choose all six themselves, because a mark ranked an option rather than answering for anybody.

The option a document marks is now selected when the page renders.
A decision the document marks nothing on arrives with nothing selected, because marking nothing is the caller saying it has no view, and picking one for it would be the page answering on the caller's behalf.
At most one option may still be marked, which the document layer already enforced and continues to enforce unchanged.
This holds on every card, whatever its severity or audience: the surface does not read either field before deciding whether to select an option, so it makes no judgement about the document it was handed.

Because a radio cannot be unticked by clicking it again, a single-choice decision now renders one further control below its options, sharing an option's furniture but not counted as one.
It carries no command, and choosing it clears the decision's selection and moves the cursor into the words field below, labelled *None of these*.
A multiple-choice decision renders no such control, because unticking every box already leaves it exactly as blank.

Agreeing with a whole document now costs one press: sending without touching anything carries every recommendation the document marked, and omits every decision the document left to the person.
The consequence accepted alongside that is that a person can submit without reading, and every recommendation runs, including a response that rejects a finding or records a decision carrying a reason the caller drafted rather than the person.
Two brakes stand against that consequence rather than the surface refusing to select anything.
The bar sits at the end of the page, in the normal flow rather than fixed to the viewport, so reaching it means having passed every card first.
The notice still fires for any decision that is blank when a press is made, and a recommendation pre-selected is not blank, but a card the caller left to the person still is, and still has to be confirmed or answered before it goes as nothing.

## Rejected

- Leaving an option unselected wherever its decision carries an error severity, so that only a caller's lower-stakes recommendations arrive chosen.
  This asks the surface to read severity before deciding what to select, which is a judgement about the document CLAUDE.md already forbids the surface from making.
- Leaving an option unselected wherever its decision carries an audience naming a person rather than the author, on the reasoning that a person still has to decide it.
  This asks the surface to read audience the same way, and fails for the same reason: severity and audience decide order and nothing else about what the page does with a card.
