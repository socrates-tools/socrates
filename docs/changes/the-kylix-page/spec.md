# agora: the Kylix page

Decisions: ADR 45, ADR 46, ADR 47.

## Problem Statement

The page reads as a form rather than as a place to decide something.

Every card has the same weight, and nothing on a card separates what was seen from where the person answers it.
The summary list at the top of the page names every source entry, but it scrolls away with the page, says nothing about which decisions are answered, and repeats each severity and audience as a word in capitals beside a name in capitals.
A person working down a long page cannot see how far they have got, or get back to a card they skipped, without scrolling for it.

A later design round is worse.
Every question an earlier round settled comes back as a card that offers no options and still carries an empty words field, so a round opening with nine settled questions asks for nothing nine times before it asks for anything, as `docs/changes/a-settled-card-offers-nothing/report.md` records.

The page has no visual identity of its own either.
It paints a light and a dark palette that are each the default look of a form, in IBM Plex, and the accent colour marks links, focus, recommendations and the primary press alike.

## Solution

The page is redrawn to the Kylix design system: a warm black ground, gold for what the person chooses or presses and, beyond that, only the wordmark, the open badge, links, the recommendation mark, the audience that needs a person and the keywords and types of code, Cinzel for the page title, Manrope for everything read or pressed, and JetBrains Mono for code and locations.
It is dark whatever the reader prefers, per ADR 45.

Beside the cards sits an index that stays in view as the page scrolls.
It names every source entry, or every decision where there is no source, with its severity and location, and a mark that fills as the person answers it.
It says how many decisions are answered, and pressing an entry takes the person to its card and puts the cursor on its first control.

Each card is two sections of one surface.
The top says what was seen: the name, the severity and audience marks, the location, the question, its meaning and the evidence.
The bottom, sunk a step, is where the person answers: the options, taking none of them, and their own words.
Code, a diff and an illustration share one evidence frame with a heading strip naming what it holds.

Severity and audience are drawn as marks, an icon in a tinted circle, whose meaning is shown on hover and on focus and read aloud as its label.
A known severity is marked by what it means, and each audience socrates review writes as the design drew it, per ADR 47.

The message field and the two presses become one card, the composer, which shows the presses while the page is open, the receipt and the way back once something is sent, and a status line throughout.
The header carries a badge saying where the exchange stands.

A document may carry what was settled before, which the page reads ahead of every open question as one quiet block that offers nothing to answer, per ADR 46.
socrates-design writes its settled questions there instead of as cards.

## User Stories

**Somebody answering a page**

1. As somebody answering a page, I want to see how many decisions I have answered, so that I know how much is left without scrolling to find out.
2. As somebody answering a page, I want the index to stay in view while I read, so that getting back to a card I skipped is a press rather than a search.
3. As somebody answering a page, I want pressing an entry to put the cursor on that card's first control, so that I can answer it without reaching for it again.
4. As somebody answering a page, I want what was seen and where I answer to look like two different things, so that I read before I choose.
5. As somebody answering a page, I want the option I chose to stand out from the ones I did not, so that I can check my answer at a glance.
6. As somebody answering a page, I want one press on the page drawn as the one to press, so that I never have to work out which button sends.
7. As somebody answering a page, I want to know what a severity or an audience means without looking it up, so that a small icon does not cost me a guess.
8. As somebody answering a page with a keyboard, I want each mark to be reachable and to say what it means when it has the cursor, so that the meaning shown on hover is not a pointer's alone.
9. As somebody answering a page, I want the header to say whether the exchange is open, sent or over, so that I know whether anything I do still reaches anybody.

**Somebody being designed with**

10. As somebody being designed with, I want what I settled in an earlier round shown as settled rather than asked again, so that a later round opens with what is still open.
11. As somebody being designed with, I want what I settled to state what I chose and what I wrote, so that I can see the argument I have already made.
12. As somebody being designed with, I want to fold what was settled away, so that a long history does not push the open questions below the fold.
13. As somebody being designed with, I want to reopen something I settled by saying so in the message, so that a design is still not a ratchet.

**Somebody maintaining the surface**

14. As somebody maintaining the surface, I want one palette stated once, so that every colour has one home and one contrast to hold.
15. As somebody maintaining the surface, I want the one place the surface knows a word of socrates review stated in one file and exempted by name, so that nothing else of review creeps in unnoticed.
16. As somebody maintaining the surface, I want every font face and icon it carries accounted for in NOTICE, so that what ships is what the licences allow.

## Implementation Decisions

**The document gains `settled`, and nothing else.**
`agora/document` gains `Settled`: an identifier, a name, a location, the question, the labels chosen and the words.
Each entry needs an identifier unique among the settled and not naming an open decision, a question, and no blank label.
A document holding only what was settled holds something, so it is served.
Nothing about settled answers reaches what a person sends: they are not decisions, carry no control and are never counted.

**The renderer builds an index, marks and a settled block.**
`agora/render` indexes the source entries where there are any and the decisions otherwise, ordered by severity like the cards.
An entry a decision covers is a link to its card, and carries a mark saying whether the decision arrives answered, which is whether it recommends an option.
Marks come from `agora/render/marks.go`, the three audiences socrates review writes from `agora/render/audiences.go`, and their icons from `agora/render/icons.go`, six Lucide icons carried as path data.
A diff states how many lines it adds and removes in its heading, and names the language its lines are written in: the one the evidence states, or else its file's extension.

**The page keeps every hook the client and the suite rely on.**
`.decision`, `.choice`, `.decline`, `.words-field`, `.message-field`, `.send`, `.finish`, `.receipt`, `.amend`, `.state`, `.bar`, `.illustration` and `.painted` keep their names and their meaning.
The glance becomes `nav.index`, whose entries keep `.entry` and `.entry-name`, and the composer card holds the message field and the bar.
Choices are native inputs drawn with `appearance: none`, so the control a person clicks is the control the form reads.

**The client counts as the person answers.**
On every change and every keystroke inside a card, `agora/web/src/surface.js` recomputes the card's mark, its entry's mark and the count, using the rule that decides what is sent: a selection or words, and taking none of them alone is not an answer.
An entry scrolls its card to 24px from the top and focuses its first control.
One tooltip, outside `main.surface` so no card clips it and no revision removes it, shows a mark's meaning on hover and focus and hides on scroll unless its mark has the cursor.
Each line of a diff is highlighted on its own where highlight.js knows the diff's language, and on a changed line the comment and punctuation colours step up to one that clears the contrast floor against the line's background.
The composer's phase, `open`, `sent` or `over`, decides what it shows, and the header badge follows it.

**The count reopens a rejection.**
The form-submit proposal in `docs/form-submit/proposal.html` rejected a counter, because every recommendation arrives chosen and a counter would read full before the person touched anything.
The index counts anyway, and a decision that arrives with its recommendation chosen counts as answered.
The confirmation is unchanged and still names every decision left blank, so the count informs and the confirmation is still the brake.

**The page is dark whatever the reader prefers.**
Per ADR 45, the stylesheet states one palette under `:root` with `color-scheme: dark`, and the text floor it states is 12px.

**The fonts and the icons are accounted for.**
Cinzel, Manrope and JetBrains Mono replace the four IBM Plex faces, each a variable latin subset under the SIL Open Font License.
`NOTICE` names each font's copyright and the licence text once, and names Lucide's ISC licence and the MIT licence of the Feather icons three of the six derive from.

**socrates-design writes what was settled to the list.**
The skill's settled section says a settled answer goes in `settled` rather than among the decisions, that an identifier is either settled or open, and that a person reopens a settled answer through the message.

**Two conventions stop misfiring, and one gains a named exemption.**
The palette check reads a custom property whose name holds a digit as a custom property, and the check that the surface knows nothing of the other tool does not read `stroke-linecap` and `stroke-linejoin` values of `round` as the word.
That check exempts the files `convention.TheOtherToolNamedIn` names, each with its reason, and reports an exempted file that names nothing of the other tool as stale.

## Testing Decisions

The seams are the ones the page already has.

**The document and the renderer are tested in process.**
`agora/document/settled_test.go` reads what was settled as the caller wrote it and refuses each way a settled answer cannot be shown.
`agora/render/settled_test.go` asserts the block is read before the decisions, states the chosen labels, the words and *Left to whoever asked*, offers no control, and is led to from the index without being counted.
`agora/render/marks_test.go` asserts what each severity mark is announced as and that a card's marks can be reached by keyboard, and `agora/render/audiences_test.go` asserts each of review's audiences is announced and drawn as the design drew it.

**The page is driven by a real browser against the real binary.**
The existing page, closing, painting and presentation scenarios keep their shape against the new markup.
The contrast floor is still measured under both preferences, and a scenario asserts the page is painted the same under either and states a dark colour scheme.
An entry in the index is pressed, and the cursor is asserted to land on the card's first control.
`internal/e2e/index_test.go` drives the count through choosing, unticking the only box, words alone and taking none, and asserts the index and the cards agree at every step.
It asserts a mark's meaning shows while the mark has the cursor and goes when the cursor leaves, that the header badge follows the exchange from open to sent to over, and that what was settled folds away and the index opens it again.
The count and the tooltip were each broken on purpose, and their scenarios failed.
A scenario colours a Go diff whose changed lines hold comments and holds it to the contrast floor, and putting the design's comment colour back on the changed lines failed it.
The second design round serves what was settled through the list and asserts it is shown ahead of the one open card and offers no control.

**The conventions are mutated.**
Each convention fix has a fabricated tree that fails before the fix: a custom property named with a digit, and an icon drawn with round line ends beside a stylesheet that says *each round*.

## Out of Scope

- A light palette, per ADR 45.
- A round on the document, a round badge, or anything that counts rounds, per ADR 46.
- Anything about what is sent, the confirmation's behaviour, or the receipt's wording.

## Further Notes

The design system was handed over as a prototype in HTML and React.
Its components are not shipped: the page stays server-rendered HTML, one stylesheet and one script, as it was.
