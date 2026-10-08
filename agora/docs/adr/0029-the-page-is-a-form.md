# 29. The page is a form

- Status: accepted
- Date: 2026-09-16
- Part of: agora's page

## Decision

A person answering a document used to press twice per decision.
Once to choose an option or write a word, and once on *Commit this answer* to tell the page that the choice it could already see was real.
Only a committed answer reached the send, and an answer left in a control and never committed reached nobody.

That second press bought a state worth naming, a card that was settled rather than merely touched, and a way to take a settled card back before sending.
It cost every person who agreed with the page twice the presses, and it hid a silent failure mode: choosing without committing looked identical to answering, right up to the point where a send carried nothing.

The page is now a real `<form>`.
Choosing is answering: what a control holds at the moment of the one press at the end of the page is what is sent.
A decision goes if it has an option selected or words typed, and is omitted otherwise, which is the rule the exchange already applied to a committed answer and now applies directly to a control.

The per-decision *Commit this answer* and *Take this answer back* controls are gone, along with the settled state a committed card took, the client's record of what was committed, and the restore-after-revision behaviour that record existed only to drive.
The discard warning is gone with them: it asked about an answer that was chosen but not yet committed, and once choosing is answering there is nothing left uncommitted to warn about.
What replaces the warning is a separate decision, because it is a separate question: whether a person sending over a decision they left blank should be told before it goes.

The two presses move out of the composer to a bar at the end of the page.
The bar is not fixed to the viewport: it sits in the normal flow, after the last decision and after the composer's message field, so reaching it means having passed every card.
It carries one primary press that sends and waits for a reply, and one quiet secondary that sends and finishes, and it renders identically whether a document holds six decisions or a page whose every decision offers nothing to choose among.
The primary is the form's first submit button in document order, so pressing Enter anywhere that does not swallow it, such as an option, reaches the primary rather than the secondary.

The page's Content-Security-Policy already states `form-action 'none'`, adopted when the page was still a collection of buttons rather than a form, because a caller's illustration is confined to a box of its own and nothing on the page should be able to navigate it away.
A `<form>` submitting through that policy would be refused by the browser and logged as a security violation, so the submit handler calls `preventDefault()` and posts over `fetch`, exactly as the buttons did before.

Nothing carries across a revision: no selection, no typed words, and not the message field.
A decision's identifier is the reviewer's and unique only inside the one document that names it, so an identifier is not a safe key to restore an answer by across a revision: a later document can reuse an identifier for an unrelated question, and restoring by that identifier would land the earlier answer on a card that never asked for it.
The client's restore-after-revision behaviour is removed rather than extended to guard against that, because a revision already replaces the whole surface with a freshly rendered one, and a fresh render already carries nothing to restore.

## Rejected

- keeping the per-decision commit and reopen controls and adding the form and the bar around them, which halves nothing: a person still presses twice per decision and once more at the end.
- keeping the committed state and the discard warning as the page's own defence against sending a blank decision, and building the notice this ticket's parent describes as an addition to them rather than a replacement, which leaves two mechanisms saying almost the same thing in different words.
- letting the `<form>` navigate natively and dropping `form-action 'none'` to allow it, which trades a real defence in depth for one fewer line of client script.
- restoring a committed answer after a revision by matching it against a card carrying the same identifier, which is the behaviour this decision removes: the identifier's scope is one document, and treating it as durable across documents can answer a question on a person's behalf that they never saw.
- fixing the identifier collision by namespacing identifiers to a revision instead, which asks every caller to invent a scheme for a problem that goes away entirely once nothing is restored.
