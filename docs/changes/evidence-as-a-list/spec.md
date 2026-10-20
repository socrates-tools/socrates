# Evidence is a list of captioned blocks, and design shows code through it

Decisions: ADR 78.

## Problem Statement

A design round often asks a person to choose between shapes of code, such as two interfaces for a component.
The design skill only teaches illustrations, so a session draws code by hand: plain `<pre>` blocks in an HTML file, unhighlighted, in a font and frame each session invents, with long lines cut off at the right edge.
Issue 78 reports this from a round comparing two shapes for a React `Select`.

agora already has a field for code: a decision's `evidence`, which review and socrates-tickets use and the page highlights with highlight.js.
It does not fit the case.
It holds one piece of code and one patch, with no caption, so two named shapes cannot sit side by side.
Its long lines scroll sideways, and on macOS the hidden scrollbar makes them look cut off.
It offers no way to copy.
It accepts only canonical highlight.js names, so `tsx` is refused even though highlight.js highlights it.

## Solution

A decision's `evidence` becomes a list of blocks, per ADR 78.
Each block holds either code in a language or a patch, never both, and may carry a caption.
The list is the only shape agora reads: a document whose evidence is an object is refused before anything is served, and the refusal says to write a list and to install the skill again.

The page draws every block the same way, whoever wrote it.
A block's heading shows its caption where there is one, with the language or the patch's file beside it, and shows the language or file alone where there is none, as today.
Every line of a code block wraps at the width of the card, with continuation lines indented, so nothing is hidden.
Every code block offers a copy control, which writes the block's code to the person's clipboard and sends nothing anywhere.
A document may name a language by any name highlight.js answers to, aliases such as `tsx`, `jsx`, `ts`, `sh` and `yml` included, and a name highlight.js does not know is still refused.

The design skill teaches evidence.
Code goes in evidence, and an illustration never draws code on its own: it may still show code as part of a screen or a layout, but a card that only shows code uses evidence.
review and socrates-tickets write a list of one where they wrote an object.

## User Stories

1. As a person answering a design round, I want two shapes of an interface shown as two named, highlighted blocks, so that I can compare them at a glance.
2. As a person answering a design round, I want each block to carry a caption that matches the option it belongs to, so that I know which code is which option.
3. As a person answering a design round, I want long lines to wrap rather than disappear past the edge of the card, so that I read every line whole without scrolling sideways.
4. As a person answering a design round, I want wrapped lines indented, so that I can still tell where one source line ends and the next begins.
5. As a person answering any round, I want to copy a block of code with one press, so that I can try it without selecting it by hand.
6. As a person answering a review round, I want evidence to look as it does today where the reviewer named no caption, so that nothing I already know how to read changes beyond wrapping and copy.
7. As a person answering a ticket round, I want each ticket's file shown as it is shown today, so that approving the breakdown is unchanged.
8. As a session running a design, I want to write code and its language as data rather than as markup, so that there is less to get wrong and fewer tokens to spend.
9. As a session running a design, I want to name a language as `tsx` or `jsx`, so that I write the name I would use anywhere else and the code is highlighted.
10. As a session running a design, I want the skill to tell me that code goes in evidence, so that I never hand-draw code in an illustration.
11. As a session running a design, I want an illustration still to be allowed to show code inside a drawn screen or layout, so that a mockup that contains code stays possible.
12. As a session writing a document in the old object shape, I want a refusal that says to write a list and to install the skill again, so that I can fix the document and the person can fix the stale skill.
13. As a session writing a block holding both code and a patch, I want a refusal naming the decision, so that I split it into two blocks.
14. As a session writing a language highlight.js does not know, I want a refusal naming the decision and the accepted names, as today.
15. As a reviewer whose finding shows code and a patch, I want to pass them as two blocks, so that each is drawn and captioned on its own.
16. As a maintainer, I want one field that shows code for every caller, so that highlighting, wrapping and copy are built and tested once.

## Implementation Decisions

- `agora/document`: `Decision.Evidence` becomes a list of blocks.
  Each block carries `caption`, `language`, `code` and `patch`.
  Reading a document whose `evidence` is an object fails with a refusal, stated in `messages.go`, that names the decision, says evidence is a list of blocks, and says to install the skill again.
  A block holding both `code` and `patch` is refused, naming the decision.
  A block holding neither is refused, naming the decision.
  `caption` is optional.
  A block's language is checked against every name highlight.js answers to, so the existing unknown-language refusal stays and its list grows.
- `agora/document/languages.json`: generated by `agora/web/build.mjs` from highlight.js's common set together with every alias each of those languages registers, rather than from `listLanguages()` alone.
  It stays generated by `make assets`, and `make test` still fails while it and its source disagree.
- `agora/render`: what a card shows becomes a list, one entry per block, each with its caption, its language and either its code or its read patch.
  The written language of a patch is worked out per block, as `writtenIn` does today.
- `agora/render/page.gohtml`: the evidences region ranges over the blocks.
  A block's heading carries the caption where there is one, with the language or the patch's file beside it, and the language or file alone where there is none.
  Every code block carries a copy control, labelled in `agora/render/messages.go`.
- `agora/web/src/surface.css`: every code block wraps at the width of its card with continuation lines indented, replacing the special case for markdown.
  The split diff of a patch keeps its table and is not part of this change.
- `agora/web/src/surface.js`: the copy control writes the block's code, as the document stated it, to the clipboard.
  Highlighting resolves a block's language through highlight.js, which answers to aliases already.
- `agora/AGENTS.md`: the invariant that only markdown wraps becomes one that every code block wraps with continuation lines indented, and the languages invariant says the accepted names are the common set and its aliases.
- `agora/CONTEXT.md`: evidence is defined as the code a decision shows, as a list of blocks, each a piece of code in a language or a patch, with an optional caption, drawn and highlighted by the page itself, unlike an illustration.
  The entry gives up `code` and `snippet`.
  This entry was written during the design.
- `design/skill/messages.go`: the design skill teaches evidence for code, with a captioned two-block example, and states that code goes in evidence and an illustration never draws code on its own.
  Its tickets section writes each ticket's file as a list of one markdown block.
- `review/skill/messages.go`: the review and handover skills write evidence as a list, and a finding showing code and a patch passes two blocks.
- The committed skills under `.claude/skills` and `.agents/skills` are regenerated with `make skill`.

## Testing Decisions

The existing seams, each for its part, with no new seam.

- `agora/document` tests: a list of blocks is read; an object is refused with the refusal naming the list and the skill install; a block holding both code and a patch is refused; a block holding neither is refused; a caption is optional; `tsx` and `jsx` are accepted; a name highlight.js does not know is still refused.
- `internal/e2e`, driving the real binary and a real browser: a caption shows in the block's heading beside its language; a block without a caption shows its language alone; `tsx` code is highlighted; a line wider than the card wraps with its continuation indented and nothing scrolls sideways; the copy control puts the block's code on the clipboard.
  `highlighting_test.go` and `presentation_test.go` are the natural homes, and the tests that write evidence today, `page_test.go`, `index_test.go`, `wholeround_test.go` and `ticketsround_test.go`, move to the list.
- `design/skill` tests: the skill teaches evidence for code, states that an illustration never draws code on its own, and its tickets section writes a list of one.
  `review/skill` tests move to the list in the same way.

## Out of Scope

- The split diff of a patch: its layout, its wrapping and any copy control on it.
- Per-option code, where an option carries the code it stands for.
- A second field for code beside evidence.
- Reading evidence written as an object.
- Any change to illustrations beyond what the design skill says about code in them.

## Further Notes

- Issue 78 proposed a `diff` language or a before and after pair for code changes.
  A block holding a patch already covers that, and the design skill's example shows it.
- `agora/document/testdata/source.json` writes evidence as an object today and moves to the list.
