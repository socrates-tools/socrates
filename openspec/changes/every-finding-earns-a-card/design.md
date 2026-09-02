## Context

See `proposal.md` for why.

The surface already holds almost everything this change needs.

`surface-document` requires that a decision always renders a free-text field beside its options, and that the field is the surface's rather than the document's.
It also already specifies the answer this change depends on: a person who writes in the field, chooses no option and sends comes back through `poll` as that decision with no option chosen and the words carried verbatim.
So an annotation with no choice beside it is a shape the surface returns today.

`surface-document` also requires that a decision's identifier, name, severity, audience and location are copied from the source the decision cites rather than taken from the decision itself.
Severity is therefore already on the page, already trustworthy, and already beyond the caller's reach.
Nothing new has to be carried for cards to be ordered by it.

`internal/document/messages.go` holds the refusal this change removes:

    "%s offers no option, and the field a person writes in is the surface's rather than the document's"

That message states its own reasoning, and the reasoning was sound when it was written.
A caller with no options and a free-text field is a caller that could have been reaching for a notes box it does not own.
What the review session showed is that there is a second reason to write a decision with no options, and it is the honest one: the caller has already resolved the finding and the person still needs to see it.

`a-page-for-every-round` archived on 2 September and moved the ground under this.
Two of its results matter here and they pull in opposite directions.

It put the source on the page: `surface-document`'s `The source is what the page shows of the whole` and `surface-page`'s `The page opens with the source at a glance` mean every entry a caller cites is rendered with its name, severity, audience and location, copied rather than composed, before the decisions.
That is the short list, owned by the surface, and it is strictly better than the skill telling a session to write one.
It is also why `findings-page`'s foregrounding requirement is withdrawn rather than edited: the half of it worth keeping now has a different owner.

It also allowed a document with no decisions at all, and in doing so stated the refusal of an option-less decision in the specification for the first time, as `A document may hold nothing to decide`'s third scenario.
Before that change the refusal lived only in `internal/document`.
So this change now amends a requirement rather than removing a line of code, and the amendment is narrow: the same reasoning that admitted a document with nothing to ask admits a decision with nothing to ask.

`web/src/surface.js` registers thirteen languages from `highlight.js/lib/core` by hand and calls `highlightElement` on every `pre code[class^="language-"]`.
`highlight.js` leaves an unregistered language unhighlighted and warns to the console, which nobody reading the page is watching.

`findings-page` requires that the page gives the findings a person answers in full and the rest of the round as a short list.

## Goals / Non-Goals

**Goals:**

- A person sees the whole round on the page, one card per finding.
- A finding the session already resolved reads as settled, with no control suggesting otherwise.
- A person can still write about a resolved finding, because reading what somebody else decided is when a person most wants to comment.
- A caller learns at its own terminal that a language will not highlight, rather than the person learning it from a grey block.
- Cards arrive in the order a person would triage them.

**Non-Goals:**

- socrates does not gain a command, flag or run record that names a page.
  The invariant that socrates states a round is one a person reads and does nothing further about it is untouched, and this change moves no part of the page into socrates.
- The surface does not learn what a finding is.
  Every requirement here is written in the surface's own nouns, and nothing in it names review, finding, verdict, round or run.
- No diagram or layout engine is embedded.
  This change adds language definitions to an already-bundled library and nothing else.
- The caller does not get to order the cards.
  Ordering moves to the surface rather than becoming a new document field.

## Decisions

### A decision may carry no options, and that is a narrowing rather than a reversal

The rule being removed guarded against a caller reaching for a free-text field it does not own.
The guard is kept by keeping the field the surface's: a decision with no options renders the same field, rendered by the surface, on the same terms, and the caller still cannot label it, size it, require it or read it except through `poll`.

What changes is only that a caller may state something without asking anything.

The alternative was leaving the refusal in place and telling callers to write an acknowledgement option.
That is what the review session actually did, and it is what the person rejected.
A control that changes nothing is not a neutral cost: the person read it as a question, and a page whose whole purpose is to separate what needs them from what does not had just put a control under the three things that did not.

The second alternative was a separate document field, a list of statements beside the list of decisions.
Rejected because it splits one page into two kinds of thing that render as the same card, and because the person's request was precisely that the settled findings look like the others rather than like an appendix.
A decision with no options already renders as a card, and one code path for cards is worth more than a name that distinguishes them.

The third alternative is the one that now exists: leave the settled findings to the source list.
It is the strongest of the three and it still does not answer the request.
An entry says what a finding is, where it is, how severe it is and who it is for, and every word of it is copied.
It cannot say that the lock deleted itself while another request held it, or that the deletion is now gone, because a caller composing that would be composing prose, and prose is the one thing the entry deliberately does not carry.
The list tells a person how much there is. A card tells them what happened. The person asked for the second and already had the first.

### A decision with no options states no mode, and one that states a mode is refused

`mode` names how a person chooses among options.
With no options there is nothing to choose among, so a stated mode is a contradiction rather than a harmless extra field.

Refusing is chosen over ignoring because the caller here is a coding session assembling JSON, and a field silently ignored is a field it will keep emitting and eventually rely on.
The refusal names the decision, which is what makes it a one-line fix at the terminal.

### The surface orders the cards by severity, and the caller cannot

Severity is copied from the source and the caller cannot reach it, so ordering by severity is ordering by something already beyond the caller's influence.
Putting the order in the surface rather than in the document follows the frame rule: the surface owns the page, its layout and every control, and the sequence of cards is layout.

The order is most severe first and informational last.
Ties hold the order the document listed them, so a caller retains a say over the sequence within one severity without having a say over the sequence between severities.
A decision whose source carries no severity sorts after every named severity, on the same stable basis, because an unclassified finding is not evidence of unimportance and putting it last is the reading that does not overstate it.

The alternative was a document field for order.
Rejected for the reason the identifier is copied rather than accepted: a caller free to order the page is a caller free to bury the error, and the page exists to stop that.

There is a real cost.
A caller with a narrative across several findings loses the ability to sequence it, and this change does not give that back.
The narrative belongs in each card's own prose, where a reader who jumped straight to the error will still find it.

### The registered language set is the library's common set, not a curated list

Thirteen hand-picked languages is a list with an owner, and its failure mode is a caller naming a language nobody thought of and getting silence.
`csharp` and `diff` were both named by the first real caller outside the project.

Registering the library's common set replaces the list with the library's own judgement about what is common, which is maintained by somebody else and updated when the library is.
The bundle grows, and that is the price.
It is paid once, into a binary that already embeds its own fonts, and it buys the removal of an entire class of silent quality loss.

Bundling every language the library ships was rejected: the common set already covers what a code reviewer's evidence is written in, and the tail is large.

### An unknown language is refused, and an absent one is not

These are different statements by the caller and they get different treatment.

No language means the caller wants no highlighting, which is a legitimate thing to want and renders as plain text.
A language the page cannot highlight means the caller believes it will be highlighted and is wrong, and the person is the one who pays for that belief.

So `serve` refuses the second, names the decision and the language, and names the set it does know.
Refusal is consistent with how the surface already treats a mode that does not exist and a document marking two options: the caller is at a terminal, the error is one edit away, and the person never sees the degraded page.

The cost is that a caller naming something outside the common set cannot serve at all until it either picks another name or drops the field.
That is the intended pressure.
The refusal names the alternative, so the way forward is in the error.

### The skill describes a card per finding, and says what a resolved card holds

The requirement being changed is `findings-page`'s foregrounding rule.
Its reasoning stands and its remedy does not: a round of one question and five mechanical defects should still not read as six equal questions, and the way to achieve that is now the severity order and the absence of controls, rather than demoting five findings to a list.

The skill gains the shape of a resolved finding's card: what the finding was, what resolved it, and the diff where the diff is what makes it clear.
That is prose guidance rather than a schema, because the diff already has a document field and a renderer, and what the skill is teaching is when to reach for it.

## Risks / Trade-offs

**A longer page.**
Six cards where there were three, and every card carrying prose.
Mitigated by the severity order putting what matters first, and by the skill continuing to require brevity.
The person asked for this explicitly, having read both shapes.

**A larger bundle.**
The common language set is bigger than thirteen languages.
Mitigated by it being a one-time cost in an already-embedded asset, and measured rather than assumed: the task list records the before and after.

**The refusal of an unknown language will annoy somebody.**
A caller mid-flow, with a page half written, will hit a refusal over a field it thought was decorative.
Accepted, because the alternative is the person reading unhighlighted code and nobody finding out.
The refusal naming the known set is what keeps it a one-line fix.

**A caller that wanted a narrative order loses it.**
Stated in the decision above rather than mitigated.
No caller has asked for it, and the guard against burying an error is worth more than a sequence no caller has yet wanted.

## Migration

The document format gains no required field, so every document that serves today still serves.

Two documents stop serving: one that gives a mode to a decision with no options, and one that names a language outside the registered set.
Neither can exist today, because the first is refused for having no options at all and the second serves while rendering plain.
The second is the only real break, and it breaks loudly at the caller's terminal with the fix in the message.

`A document may hold nothing to decide` keeps both of its first two scenarios unchanged, and its third narrows from any decision with no option to a decision whose options carry no command.
A caller relying on the old refusal was relying on being stopped, and nothing was built on it.

The skill is generated from `internal/skill`, so the committed copy is regenerated and the suite fails while the two disagree.
