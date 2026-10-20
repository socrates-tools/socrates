# 78. Evidence is a list of blocks, and the list is the only shape read

- Status: accepted
- Date: 2026-10-20
- Part of: evidence-as-a-list

## Decision

A decision's evidence was one object: one piece of code in a language, one patch, or both drawn one under the other.
A design card comparing two shapes of an interface needs two pieces of code, each named, so sessions drew code by hand in an illustration, unhighlighted and clipped, as issue 78 reports.

Evidence is now a list of blocks.
Each block holds either code in a language or a patch, never both, and may carry a caption.
review, socrates-tickets and socrates-design all show code through this one field, and review and tickets pass a list of one.

The list is the only shape agora reads.
A document whose evidence is an object is refused before anything is served, and the refusal says to write a list and to install the skill again, because a skill installed from an earlier release still writes the object.

## Rejected

- Each option carries its own code.
  The code would sit under the option it belongs to, but shared code that belongs to no option would still need evidence, so a decision would show code in two places.
- A new code list beside evidence, as the issue proposed.
  Two fields on one decision would show code, one for review and tickets and one for design.
- Teaching the design skill evidence as it was, and changing nothing in agora.
  One block per card puts two shapes in one block separated by comments, with no caption naming either.
- Reading both shapes, an object as a list of one.
  Nothing written today would break, but agora would read two shapes for good.
