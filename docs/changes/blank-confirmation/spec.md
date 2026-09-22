# agora: confirming a send over a decision left blank

Decisions: ADR 40, ADR 41.

## Problem Statement

The page arms rather than asks.

Press *Send and wait for a reply* or *Send and finish* while any decision holds neither a selection nor words, and nothing is sent.
Instead an inline panel opens below the two buttons, headed *Left blank, and about to go as nothing:*, listing each blank decision as a link to its card, and closing with the hint *Press the same button again to send anyway* beside a *Keep answering* button.
The client records which button was armed, and the next press of that same button sends.

Three things are wrong with it, and they compound.

**The same press means two different things, and nothing on the button says which.**
The first press does not send and the second does.
Whether a press sends depends on invisible state the person cannot see and the page never shows: which button was armed, and whether the panel below is describing the press they are about to make or the one they made a moment ago.
Pressing the other button, or answering a card and coming back, silently disarms.

**Discovering the mechanism is reading a hint.**
Nothing about a button labelled *Send and finish* suggests that pressing it twice is how it works.
The person has to read the panel, find the sentence explaining that the button they already pressed now does something else, and act on it.
A confirmation people have to be taught by a line of body text is a confirmation that has already failed.

**The panel is easy to miss and easy to ignore.**
It opens in the normal flow at the end of a page that may be long, below the buttons and near the browser's bottom edge, and it neither takes the cursor nor stops the person doing anything else.
On a page with cards above the fold, a person who pressed send can be looking somewhere other than the one region that is telling them their press did nothing.

ADR 29 removed the per-decision commit press and left this as an open question in so many words: whether a person sending over a decision they left blank should be told before it goes.
The arming panel was the answer given, and ADR 31 then leaned on it as one of the two brakes standing against a person sending a whole page of recommendations without reading it.
The question this change reopens is not whether to tell them, which is settled and stays settled, but how.

## Solution

The press raises a dialog that asks outright.

Pressing either send button while any decision is blank opens a modal dialog over the page.
The dialog names every decision left blank, says what sending will do to them, and offers two presses: *Send anyway*, which sends exactly what the controls hold, and *Go back and answer*, which closes the dialog and changes nothing.
There is no armed state and no second meaning for a button: a press either sends, or raises a dialog whose own press sends.

The dialog names each blank decision as a control that closes the dialog and takes the person to that card, so the jump the panel offers today survives.
Cancelling leaves nothing behind: the page is exactly as it was and the cursor returns to the button that was pressed, whether the person cancelled with *Go back and answer*, with Escape, or by pressing the backdrop.
The dialog opens on *Send anyway*, which carries the cursor and is drawn as the primary, because the person already asked to send and confirming should cost a press rather than aiming.

A page with nothing blank is untouched.
It sends on one press, on either button, exactly as it does today.

The inline panel goes, along with its *Keep answering* button, its styling and the client's record of what was armed.

## User Stories

**Somebody answering a page**

1. As somebody answering a page, I want to be asked a question rather than warned, so that I understand what my press will do without reading a hint about how the button works.
2. As somebody answering a page, I want the confirmation in front of what I was looking at, so that a press that did not send cannot be missed because I was looking at a card rather than at the foot of the page.
3. As somebody answering a page, I want a press to mean one thing, so that pressing send twice is never the mechanism and never the accident.
4. As somebody answering a page, I want confirming to cost one press with nothing to aim at, so that deliberately leaving a card blank is cheap.
5. As somebody answering a page, I want the confirmation to name the decisions it is talking about, so that I can tell at once whether I left them blank on purpose.
6. As somebody answering a page, I want a name in the confirmation to take me to its card, so that changing my mind about one of them is a press rather than a search.
7. As somebody answering a page, I want going back to leave the page exactly as it was, so that cancelling never costs me anything I typed or chose.
8. As somebody answering a page, I want Escape and the backdrop to mean going back, so that the dialog behaves like every other dialog I have used.
9. As somebody answering a page who agrees with every recommendation on it, I want to be shown nothing at all, so that the common path stays one press.

**Somebody reading the page with a screen reader or a keyboard**

10. As somebody using a screen reader, I want the confirmation announced as a dialog and its heading read, so that I know a question was asked rather than hearing a status update from somewhere below.
11. As somebody using a keyboard, I want the cursor inside the dialog when it opens, so that the next thing I type reaches the question I was just asked.
12. As somebody using a keyboard, I want the cursor back on the button I pressed when it closes, so that cancelling puts me where I was.
13. As somebody using a keyboard, I want the cursor kept inside the dialog while it is open, so that tabbing does not wander into a page I cannot act on.

**Somebody maintaining the surface**

14. As somebody maintaining the surface, I want no armed state in the client, so that whether a press sends is a function of the page rather than of what happened before.
15. As somebody maintaining the surface, I want the confirmation driven end to end by a real browser, so that the dialog's behaviour is asserted rather than its markup.

## Implementation Decisions

**The trigger is unchanged, and it is the only trigger.**
A press raises the dialog where at least one decision holds no selection and no words, which is `blankCards` today and stays exactly that.
A page with nothing blank sends on one press, on either button.
*Send and finish* over a fully answered page raises nothing, so finishing is not confirmed on its own account.

**The dialog is a native `<dialog>` opened with `showModal`.**
The browser gives the modal semantics the stories ask for: the cursor moves inside, it is kept inside while the dialog is open, Escape closes it, and the rest of the page is inert.
`agora/render/page.gohtml` renders it once per page, as a sibling of `form.exchange` rather than inside it, so that its presses are not submit buttons and the form's first submit button in document order stays *Send and wait for a reply*, which is what ADR 29 relies on for Enter.
Its presses carry `type="button"` and are driven by the client, as every other control on the page is.
The page's `form-action 'none'` policy is untouched and needs no clause: a dialog is not a form action, and nothing about the dialog navigates.

**The dialog's body is built by the client from the cards, as the panel's was.**
`showNotice` becomes the code that fills and opens the dialog, and `hideNotice` becomes the code that closes it.
Each blank decision renders as a `<button type="button">` carrying the decision's identifier, styled as a link, rather than as the `<a href="#decision-...">` the panel uses: an anchor inside an open modal would move the page behind it without closing it, and the fragment would stay in the address bar afterwards.
Pressing one closes the dialog, moves the card into view and puts the cursor on that card's first control, which is its first option where it has options and its words field where it has none.

**Cancelling is one path with three doors.**
*Go back and answer*, Escape and a press on the backdrop all close the dialog, change nothing, and return the cursor to the button that raised it.
Escape is the browser's own `cancel` event on the dialog, and the backdrop is a press whose target is the dialog element itself rather than anything inside it.
The client holds which button raised the dialog only for the length of the dialog being open, to return the cursor to it and to know which press to make on *Send anyway*, and it is cleared on close.
That is not the armed state this change removes: it does not outlive the dialog, and no press behaves differently because of it.

**The dialog opens on *Send anyway*.**
It carries `autofocus` and the primary styling, per ADR 41.

**The arming state and the panel are removed outright.**
`armedFor`, the `.notice` block in `page.gohtml`, the `.notice`, `.notice-message`, `.notice-decisions` and `.notice-hint` rules in `agora/web/src/surface.css`, and the `.keep-answering` control and its binding all go.
The `.state` live region stays and keeps saying what it says today.
`take`, which replaces the surface on a revision, closes any open dialog rather than clearing an armed flag, because a revision replaces the cards the dialog is naming.

**The words the dialog says.**
The heading names the count: *Send with 2 of 5 left blank?*
The body states the consequence in the page's own register: *A decision left blank goes as nothing. Whoever asked will see no answer to it.*
The presses are *Send anyway* and *Go back and answer*.
These are the spec's proposal and not a settled answer; see Further Notes.

**The glossary gains one word.**
`CONTEXT.md`, under The page, gains **confirmation**: the dialog a press raises where any decision is blank, naming each one, and the only thing on the page standing between a press and a send.
The entry lands with the change rather than ahead of it, because the glossary describes the page that exists.

**One invariant is reworded and none is dropped.**
`CLAUDE.md` says today that a decision left blank is named before either press sends.
That stays true and gains the mechanism: it is named in a dialog the press raises, and a second press inside that dialog is what sends.
The neighbouring clauses, that what the controls hold at the submit is what is sent, that nothing is committed per decision, and that a decision holding no selection and no words is omitted, are untouched.
ADR 31's statement that the notice still fires for any blank decision is superseded by ADR 40 and the brake it describes is unchanged in effect.

## Testing Decisions

The seam is the real binary driven by a real browser, which is where the panel's tests already sit.

**The existing tests move rather than multiply.**
`internal/e2e/closing_test.go` holds the panel's scenarios: that a press over a blank decision delivers nothing, that the notice names the blank decision, that its links point at the right cards, that pressing the same button again sends, and that `.keep-answering` closes it.
Each is rewritten against the dialog, and the one that asserts a second press on the same button sends becomes the one that asserts a press on *Send anyway* sends.
`internal/e2e/page_test.go` holds the two scenarios where a decision is left blank by declining every option and by unticking every box, and both keep their shape with the dialog as the subject.

**The dialog's modality is asserted, not assumed.**
A scenario asserts that the dialog is open and modal after a press over a blank decision, that the cursor is on *Send anyway* when it opens, and that the cursor is back on the button that raised it after each of the three ways of cancelling.
These are read from the live document in the browser, never from served markup, because asserting on markup would test the renderer.

**The three doors are three scenarios.**
Cancelling with *Go back and answer*, with Escape and with a press on the backdrop each get a scenario, and each asserts the same two things: nothing was delivered, and what the controls held is still held.

**Jumping to a card is a scenario.**
Pressing a blank decision's name in the dialog closes it, and the cursor lands on that card's first control.
It is asserted for a decision with options and for a decision with none, because the first control differs.

**The absence of arming is asserted.**
A scenario presses send over a blank decision, cancels, and presses the same button again, and asserts that the dialog opens a second time rather than sending.
This is the fault the change exists to remove, so it is named for it.

**The untouched path keeps its scenario.**
A press with nothing blank delivers on one press and raises no dialog, which `closing_test.go` already asserts and which keeps asserting it against the dialog.

**A guard is broken on purpose.**
Per the repository's conventions, the check that a press over a blank decision delivers nothing is mutated once to confirm it fails.

**Contrast and paint are checked as the page's own rules require.**
The dialog's colours are tokens already stated, and the closed dialog paints nothing, which is the existing rule that an element the page hides paints nothing.

## Out of Scope

- Confirming anything other than a blank decision. A fully answered page sends on one press, on either button, and *Send and finish* is not confirmed on its own account.
- Refusing to send a page where every decision is blank. A document may hold only cards that offer nothing to choose among, and a person who wants to send nothing but a message is entitled to.
- Marking a blank card on the page. The confirmation names the blanks and takes the person to them, and the page carries no new state after it closes.
- Anything about the receipt, the amend path, or what happens after a send.
- Anything about review, which neither raises this dialog nor knows it exists.

## Further Notes

The dialog is a straight replacement for the panel and reaches nothing the panel did not.
It is a client change, a template change, a stylesheet change and the e2e scenarios that cover them, and no Go package outside `agora/render` is touched.

### The design ended early

The interview ended on round one, so the frontier below reached no page.

- **Whether the dialog says something different over *Send and finish*.** Finishing closes the exchange, and nothing further can be sent from the page afterwards, so there is no amending a send that went out with three cards blank. The trigger is settled and covers both presses, but whether the dialog names which press it is confirming, and whether it says more about the one that cannot be taken back, was never asked.
- **The exact words.** The heading, the body sentence and the two presses are proposed in Implementation Decisions and were not put to anybody.
- **Whether the dialog names what will go as well as what will not.** It could list the answered decisions beside the blank ones, which makes it a summary of the whole send rather than a question about the gaps.
- **Whether a page with every decision blank reads differently.** *Send with 5 of 5 left blank?* is a different sentence from a page with two gaps in it, and it may deserve different words rather than the same ones with different numbers.
- **The last card.** The design never reached the round that proposes where the spec lands and whether it is published to the tracker, so `docs/changes/blank-confirmation/spec.md` is the session's choice under the repository's own convention rather than a path anybody approved, and nothing has been published.
