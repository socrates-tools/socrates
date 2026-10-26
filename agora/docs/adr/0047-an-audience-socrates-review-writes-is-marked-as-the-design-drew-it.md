# 47. An audience socrates review writes is marked as the design drew it

- Status: accepted, its gold tone and its marks on a card amended by ADR 53, its audience marks drawn outlined by ADR 89
- Date: 2026-09-30
- Part of: the-kylix-page

## Decision

A decision's severity and audience are drawn as marks, an icon in a tinted circle whose meaning is shown on hover and on focus and read aloud as its label.

The three severities the surface already orders by, `error`, `warning` and `info`, are each marked by what they mean: a hue and an icon of their own, and a sentence saying what the severity asks of a person.
A severity outside those three is marked neutrally, under its own name.

The three audiences socrates review writes are marked the same way, as the design drew them.
*needs a person* is gold, a person's icon, and says it is *Your call*, that only a person can decide it and that the reviewer will not answer it.
*can be patched* is a wrench and says the author can patch it without the person.
*informational* is an eye and says there is nothing to decide.
Any other audience is marked neutrally, labelled *Decides:* and the audience exactly as the caller wrote it.

The severity hints carry socrates review's reading of a severity too: the one for `error` says the change should not land as it is.
They are kept as the design wrote them, and they name nothing the check that the surface names nothing of the other tool looks for.

The audiences are the one place the surface knows a word of socrates review.
It is stated in one file, `agora/render/audiences.go`, which the check that the surface names nothing of the other tool exempts by name, and that exemption fails the suite once the file names nothing of review.

What this costs is that agora knows socrates review's audiences, and a tool that writes the same words with another meaning gets review's sentence under them.
What it buys is the distinction the design draws at a glance: a card that is the person's own call stands out from one the author will patch without them.

## Rejected

- One neutral mark for every audience, labelled with the caller's words.
  It keeps the surface knowing nothing of review, and it was the first thing built, but it gives up the distinction the design exists to draw, and the person asked for the design.
- A field on the document carrying an audience's icon and sentence.
  It extends what every caller writes for the sake of how one tool would like to be drawn.
