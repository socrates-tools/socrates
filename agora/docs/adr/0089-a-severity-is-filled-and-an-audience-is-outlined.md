# 89. A severity is filled and an audience is outlined

- Status: accepted
- Date: 2026-10-26
- Amends: ADR 47, ADR 70
- Part of: the-aegean-page

## Decision

A severity mark is drawn filled, its ink on its own wash, so its hue says how serious a finding is.
An audience mark is drawn outlined, its ink on the card inside a 1px border of its tone, so it says who decides rather than reading as a second severity.
`ui-badge` carries the difference as its `variant`, `filled` or `outlined`, and the renderer chooses it from the kind of mark rather than from its tone.

*Your call* stays Aegean, as ADR 70 decided, now outlined rather than on the Aegean wash.
The index rail draws the *Your call* mark, icon alone, beside an entry's severity, and draws no other audience there.
An outlined mark drawn as an icon alone is ringed in its tone's faint ring, `accent-ring` for *Your call*, rather than in its ink.
A filled wash fades at its edge and is read at the size of its icon, while a ring in full ink shows its whole edge, so at the same 18px the outlined circle read larger than the filled one beside it.

## Rejected

- Keeping every mark filled and telling the kinds apart by hue alone.
  Side by side on a card the two looked like two marks of one kind, and *Your call* on the Aegean wash read as a chosen option, whose wash is nearly the same.
- A key on the page naming each colour.
  The tooltip on every mark already says what it means, and a shape that answers the question needs no legend to be read.
- Moving the highlighting of code off the hues the states use.
  Code is read inside its own block, where a string is not mistaken for something settled, and the change would have spent four new hues for no reading the page needs.
- Drawing the outlined mark smaller in the rail to match the filled washes by eye.
  The circle crowds the person glyph and sits off the row's rhythm, where a fainter ring at the same size matches without either.
- Drawing every audience in the index rail.
  The rail is scanned for what waits on the person, and a wrench or an eye beside each entry is noise for that.
