## Context

See `proposal.md` for why.

What exists today, and what each part of it constrains:

- `dialogue.APersonReads(classified, state)` is one pure predicate over the classified findings and the state the round led to.
  `report.FindingsOf` reads it once, and `socrates findings` prints `read by a person` and carries `read_by_a_person` in its JSON.
  Nothing else in the source reads it, by requirement.
- `internal/skill` is the one source of the skill, and the committed copy is generated from it.
  The skill's `awaiting-response` and `escalated` branches serve a page; the `clean` and `failed` branches do not mention one.
- `document.Validate` refuses a document holding no decision, and refuses a decision holding no option or an option with no command.
- `document.Source` is read only to copy each decision's name, severity, audience and location out of the entry it cites.
  An entry no decision cites is accepted and rendered nowhere, so the skill's short list of the rest of the round is prose the session writes into the preamble.
- `agora serve` starts a detached server if none is running, prints the address of the page and returns.
  The surface collects a session that has been idle for `surface.IdlePeriod`, and stops once it has held no session for that period.
  Nothing about serving requires anybody to poll.
- `internal/render` renders title, preamble and decisions; the client in `web/` owns the controls and the committed bundle is generated from it.
- The surface names no review, finding, verdict, round or run as a noun, and a convention check holds that.

## Goals / Non-Goals

**Goals:**

- The page reaches the person on every round that produced a review result, at no cost to the loop where nobody is being asked anything.
- One rule still decides the one thing that still needs deciding, and it is still not re-derived anywhere.
- The short list on the page stops being prose the session retypes.
- agora does the new work. socrates gains no command, no flag, no field in a run record and no network reach.

**Non-Goals:**

- A page for `failed` or `cancelled`. See `findings-page`.
- Any way for the person to acknowledge a page, or for the session to learn that they read one. A page nobody waits on is a page nobody hears about.
- Configuring the page away. There is no flag to suppress it and no configuration field for it, because a person who does not want to read it does not open it.
- Changing how long the surface holds a page, what it binds, or how a session is minted.

## Decisions

**The predicate is renamed rather than split in two.**
The obvious move is two predicates, one for whether a page is served and one for whether the session waits.
The first would answer `true` for every input the second is ever asked about, so it would be a function that decides nothing, sitting beside one that decides everything, and the pair would eventually be called in the wrong order by somebody who assumed they were symmetrical.
So `dialogue.APersonReads` becomes `dialogue.AwaitsAPerson`, with the same signature, the same table and the same exhaustive test over the whole input space, and which rounds reach a page is stated in `findings-page` as all of them rather than computed anywhere.
Alternative rejected: keeping the name and changing what the skill does with it, which leaves the source saying `a person reads` about the one round out of four that a person is not merely reading.

**`failed` and `cancelled` get no page, and the rule is stated over the state the run ended in.**
A page for a failed run would hold the one line socrates already printed, and the person would open a tab to read it a second time.
A cancelled run is the case where the two possible statements of the rule come apart: its first round may well have completed a review and written a findings document, so a rule phrased as "a round completed a review" would hand over a page for it.
It is excluded on purpose. Somebody stopped the review, so nothing was decided about the change, and a page of what it got as far as answers a question they withdrew.
Stating the rule over the state also keeps it agreeing with `report.reviewCompleted`, which reports a cancelled run as unreviewed.
Alternative considered: a uniform "every ended run gets a page", which is a tidier sentence and a worse page.

**A document may hold no decision, and the floor becomes a title plus one of a decision, a source entry or a preamble.**
The alternative was to keep `serve` refusing a decisionless document and have the session invent a decision with a single option, something like "understood".
That is a control that records nothing, on a page nobody is polling, and the surface's own rule is that it offers no control whose effect the caller never sees.
The floor is stated as "the page would hold nothing" rather than "no decisions", so the refusal still catches an empty document.

**The source is rendered, and it is what the short list is made of.**
Today the source is a lookup table that only decisions read, and the list of everything else in the round is prose.
Rendering the source makes the list mechanical for the same reason the name, severity and audience are already copied from it: a caller free to retype the list is a caller free to write a list that disagrees with its own terminal.
It also makes a decisionless page worth serving at all, because otherwise a clean round's page is a heading and a paragraph.
An entry a decision covers links to that decision rather than repeating it, so the list is a table of contents rather than a second rendering.
Alternative rejected: a new top-level array of read-only items beside `decisions`, which is a second way to say what a source entry already says.

**A page with nothing to answer offers neither control, and says so.**
An empty send would deliver an empty answer to a caller that is not listening, which is the one interaction on this surface that cannot be taken back being used for nothing.
The absence has to be stated in words, because an absent control and a control that failed to render look the same to the person.
This is also why the page keeps the frame, the palette and the layout: a page that is read rather than answered is not a degraded page, and the caller supplies no more markup for one.

**Where the document is drafted follows what follows the round, not which branch it is.**
The existing rule is already this rule stated over two branches: inline where the run is over, delegated where another round follows.
Adding two branches makes the general form the shorter one to write, and it decides the new branches without a new argument.
`clean` drafts inline for exactly the reason `escalated` does.
`awaiting-response` delegates whether or not a person is being asked anything, because what costs the driving context is the finding prose, and that context is about to patch, record, continue and read another round.

**The address is handed over once, at the first round that produced a result.**
This is the existing "one session for the whole review" rule, now reached on every branch rather than on two of them.
A second round writes to the same session with `agora serve --session`, which already exists and needs no change.

**socrates changes only what it says.**
The predicate's name, the line, and the JSON field.
No command, no flag, no run-record field, no dependency and no build tag, so the `findings-page` requirement that socrates does nothing further about the page holds unchanged.

## Risks / Trade-offs

**Every clean run now costs a drafted document where it used to cost a sentence.**
Mitigation: the source list is copied rather than composed, so a clean round's document is a title, a preamble of a sentence or two and the entries socrates already reported.
The page bar's line about spending tokens as though they were the person's is what governs the preamble, and it is already stated once and rendered into the skill.

**A person who never opens a page leaves a surface process running for the idle period.**
Mitigation: none needed, and none added.
This is what already happens when a person ignores a page they were asked to answer, the collection rule already handles it, and shortening the period for the new case would shorten it for the case where somebody is being waited on.

**A session could hand over an address and then exit before the person opens the page.**
Mitigation: it already can, and the server is detached from the session that started it.
The page's lifetime is the surface's idle period rather than the session's, which is the property that makes serving a page nobody waits on safe.

**Renaming the reported field breaks anything reading `read_by_a_person`.**
Mitigation: the skill is the only reader, it is generated from `internal/skill`, and it is rewritten by this change.
A test asserts the field's name, so a stale reader fails the suite rather than reading `false` forever.

**A delegated subagent asked for a document about a round with nothing for a person could invent a decision to fill the page.**
Mitigation: the skill says what a page with nothing to answer holds, and `serve` accepts a decisionless document, so there is nothing to work around.
An invented decision would carry a command the session did not offer, which the surface already refuses.

**The convention check that no surface package names review, finding, verdict, round or run could be tripped by the new rendering.**
Mitigation: the new work is stated in the surface's own vocabulary, which is entries of a source, and the check is what enforces it.

## Migration Plan

The order is chosen so `make lint test` passes after every task rather than only at the end.

1. agora accepts a decisionless document and renders the source. Nothing in socrates reads either yet, so the suite stays green on the existing skill.
2. socrates renames the predicate, the line and the field, and its tests move with them.
3. The skill is rewritten over the four branches and regenerated, and the end-to-end tests that drive the loop follow.
4. `AGENTS.md` is brought back into step with the invariants as they then are.

There is no rollback story and none is needed: nothing persists across the change, no run record gains a field, and a run started before it is read by the same commands afterwards.
