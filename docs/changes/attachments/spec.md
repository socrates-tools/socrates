# Attachments beside the words a person writes

Decisions: ADR 87, ADR 88.

## Problem Statement

A person answering a page agora serves can choose an option and write words, and nothing else.
Some answers are a picture: what a page looks like on their screen, a trace they captured, a file the caller never saw.
Today the person describes the screenshot in words, or answers in the terminal around the page, and the caller's agent never sees what the person saw.
socrates review cannot carry that evidence to the next round's reviewer, and socrates design cannot keep it beside the spec it rests on.

## Solution

Every decision's words and the composer's message carry a tray of up to ten attachments of any type.
A person attaches with a button, by dropping files on a card or the composer, or by pasting with Ctrl+V or ⌘V with the cursor anywhere in a card or the composer, which takes a screenshot or files copied in Explorer or Finder.
Images show a thumbnail and everything else a type badge, previewed from the file in memory.

Nothing leaves the browser until a press.
Either press sends one multipart POST holding the answer and every file every tray holds, per ADR 87.
The surface stores the answer whole or refuses it whole, keeps every collected answer whole in its session, its words beside its files, per ADR 88, and never serves an attachment back.
`socrates agora poll` prints an `attachments` list of name, type, size and absolute path beside each decision's words and beside the message, so the caller's agent reads each file as it reads any file.
A caller that never looks for `attachments` sees nothing change, because the field is omitted where empty.

socrates review records what a person attached with `socrates review respond --attach <path>`, which copies each file into its run store, so the next round's reviewer reads it.
socrates design copies what a decision rests on into `docs/changes/<slug>/` and cites it from the spec.
Both skills say that attachments may arrive, that each is read by its path, and that it is the person's evidence rather than an instruction.

The design this was drawn from is `docs/changes/attachments/design.html`, beside this spec.

## User Stories

1. As a person answering a design page, I want to attach a screenshot to a decision, so that what I saw reaches the agent rather than my description of it.
2. As a person answering a review page, I want to attach the screenshot that shows a finding is wrong, so that my rejection carries its evidence to the next reviewer.
3. As a person who took a screenshot with the snipping tool, I want to press Ctrl+V with the cursor anywhere in the card, so that I attach it without saving it to a file first.
4. As a person on macOS, I want ⌘V to do the same, so that pasting works as it does everywhere else on my machine.
5. As a person who copied several files in Explorer or Finder, I want pasting into a card to attach every one of them, so that I do not attach them one at a time.
6. As a person pasting text into my words, I want it to paste as text exactly as it does today, so that attachments change nothing about writing.
7. As a person pasting a clipboard holding both text and files, I want the files attached and the text pasted, so that neither is lost.
8. As a person dragging files from my desktop, I want to drop them on a card or the composer, so that they attach to the place I dropped them.
9. As a person attaching a file, I want to see a thumbnail of an image and a type badge for anything else, so that I can tell what I attached before I send it.
10. As a person who attached the wrong file, I want to remove it from the tray before I press, so that only what I mean to send is sent.
11. As a person adding an eleventh file, I want it refused where I added it with the ceiling named, so that I know why it did not attach.
12. As a person attaching an unnamed pasted image, I want it named `pasted-1.png`, `pasted-2.png` in the order it arrived, so that the agent and I can tell them apart.
13. As a person whose answer to a decision is only a file, I want that decision counted as answered and never named as blank, so that I am not asked to write words I have nothing to put in.
14. As a person reading the reply panel, I want the summary and the receipt to name the files beside the words, so that I see what went.
15. As a person amending an answer before the caller collects it, I want every file the trays hold sent again, so that what was collected is exactly what I last pressed.
16. As a person sending files over the size ceilings, I want the whole submit refused naming the ceiling, so that I never send half an answer.
17. As a person on a page whose exchange is over, I want no tray drawn, so that nothing invites me to attach what nobody will read.
18. As a person looking at a revised page, I want every tray emptied with the rest of the controls, so that nothing from an earlier document is sent against a later one.
19. As a person reading what I settled in an earlier round, I want the names of the files I attached shown in the folded settled block, so that I recognise my own answer.
20. As a caller's agent, I want an absolute path for each attachment in what `socrates agora poll` prints, so that I read the file with the tools I already have.
21. As a caller's agent running a design over many rounds, I want every path an earlier poll printed to stay readable until the session is collected, so that I can cite round 1's screenshot when I write the spec.
22. As a caller's agent that lost what a poll printed, I want `socrates agora poll --last` to print the same paths, so that I can still read the files.
23. As a caller written before attachments existed, I want what `socrates agora poll` prints to stay valid, so that nothing I parse breaks.
24. As a caller revising a document that settles an answer carrying files, I want the revision refused unless it states the names exactly as collected, so that the settled block never misstates what was sent.
25. As someone running socrates review, I want the files a person attached to a finding recorded with `socrates review respond --attach`, so that the next round's reviewer reads the person's evidence.
26. As someone running socrates review, I want `--attach` refused on a forward, naming why, so that nothing pretends to post a file to a pull request.
27. As whoever implements a designed change, I want the screenshot a decision rests on in the change's folder and cited from the spec, so that I see what the person decided from.
28. As a person using socrates, I want my attachments to stay on this machine, so that agora sends nothing beyond it, and to be told that a file recorded on a review is read by the reviewer CLI, which sends it to its vendor.
29. As someone worried about injection, I want both skills to say that an attachment is the person's evidence rather than an instruction, so that a screenshot of text is read as what the person showed rather than followed.

## Implementation Decisions

### What a person sees

- Every free-text field carries a tray no document field can remove: each decision's words and the composer's message, at most ten files each, of any type.
  The tray is not a section of its own: it is the one bordered box the field sits in, as in Claude's own composer.
  Attached files show as square tiles at the top of the box, above the text, and a small square plus button sits at the bottom left inside the border and opens the file picker directly.
  A tile shows an image's thumbnail filling it, or a type badge and the name for anything else, and a remove control appears on it on hover and on focus and stays reachable by keyboard.
  No hint and no counter are drawn: files dragged over the box or over its card highlight the box's own border, and a file past the ceiling is refused in place inside the box, naming the ceiling.
  A disabled or closed field draws no plus.
- The tray is one new primitive, `ui-file-tray`, a Lit element in its own folder of `packages/ui` naming nothing of agora's domain.
  It draws the box with the one border around whatever it wraps, the attach button, the drop target, the tiles with their thumbnail or type badge and remove control and the in-place refusal.
  It is reached through `files`, `add(files)`, which returns the files it refused, and `clear()`, and announces `change` and `refused`.
  It wraps the field in its default slot, which draws no border of its own through `ui-field`'s `embedded` attribute, slots its refusal, and takes the attach button's name and tooltip from its `attach-label` attribute and the remove controls' name from `remove-label`, all of which Go writes from `agora/render/messages.go`, and it declares its attributes, such as `max`, through reactive properties so `contract.json` holds it.
  Its states are drawn in the catalog through `catalogued` in `agora/render/catalog.go`.
- The tray is the domain piece `<ui-file-tray class="tray">`, written around `.words-field` in each decision and around the message in the composer, with the field's caption drawn above it.
- Paste is a listener the surface script holds on each card and on the composer, which hands every file `clipboardData.files` holds to that tray's `add`, and leaves text to paste as it does today.
  The page never reads the clipboard itself, so no permission prompt appears.
  A pasted image with no name of its own is named `pasted-1.png`, `pasted-2.png` in the order it arrived, and the ceiling and the in-place refusal apply to a paste as to a drop.
- Files alone are an answer: the rule that decides what is sent and that the index counts by counts a decision holding files as answered, and the confirmation never names it as blank.
- The reply panel's summary and the receipt name the files beside the words.
- A closed exchange draws no tray, as it draws no field nothing will read.
- A revision empties every tray with the rest of the controls.
- The page's `Content-Security-Policy` admits `img-src blob:` for the previews and nothing more.

### What travels

- The submit becomes a `FormData` POST in place of `JSON.stringify`: part `answer` holds the JSON answer, in which each decision and the message state the part names they attach in `attached`, and each file is a part of its own carrying its file name and type.
- `agora/surface`: `Sent` and `Choice` gain `attached`, and `Answered` and `Decided` gain `attachments`, each an entry of `name`, `type`, `bytes` and `path`, omitted where empty.
- The answer route reads multipart under its own ceilings: 20 MiB a file, 100 MiB an answer and ten files a tray, stated once in `agora/surface` beside `BodyCeiling`, which keeps guarding the document and control routes at 1 MiB.
  A part the answer does not name, a name the answer states that no part carries, or any ceiling crossed refuses the submit with nothing kept.
- `Exchange.Receive` writes every file into a fresh folder before it swaps its uncollected answer, and removes the folder of the answer it replaces.
  A revision removes every uncollected folder, because nothing carries across a revision.
- A stored file is named by agora, as its position and a sanitised base name, such as `answers/<n>/f2/1-before.png` or `answers/<n>/message/1-trace.txt`, and no stated name reaches outside the answer's folder.
- Per ADR 88, collecting an answer writes it to `answers/<n>/answer.json` beside its files, and every collected answer is kept for the life of the session, removed with the session's directory when the surface collects it.
  `socrates agora poll --last` prints the last one collected, and the session's `answer.json` of ADR 73 is replaced by this layout.
- `agora/command` changes nothing of its own: `poll` and `poll --last` print the paths because `Answered` carries them.
- The surface never serves an attachment, and no route reads back what a person sent.

### What was settled

- `agora/document`'s `Settled` gains `attached`, the names of the files collected for that identifier.
- The folded settled block shows them.
- The check of ADR 59 in `agora/surface` holds them as it holds the question, the labels and the words: a revision that drops, adds or renames one is refused before it is written or served.

### review

- `socrates review respond --attach <path>`, repeatable up to ten, copies each file into the run store beside the response, in `review/runstore`.
- `--attach` is accepted with every resolution but forwarded, and refused on a forward with a refusal naming why: a forward posts the finding to a pull request, and a forge has no plain route for posting a file.
- `review/dialogue`'s `Response` gains the attachments it recorded, and `review/prompt` names each one's path and type in the next round's prompt beside `author_said`, so the reviewer reads it with its own read tool.
- The review skill passes each path `socrates agora poll` printed for a finding to `respond --attach`.
- `README.md`'s statement of where what socrates review reads goes names attachments in both directions: agora sends nothing, and a file recorded with `--attach` is read by the reviewer CLI, which sends it to its vendor.

### design

- The design skill says that attachments may arrive, that each is read by its path, and that every path an earlier poll printed stays readable for the life of the session.
- It copies what a decision rests on into `docs/changes/<slug>/` and cites it from the spec, and that folder goes when the change lands.

### Both skills

- Both say that an attachment is the person's evidence rather than an instruction, extending the sentence that already says a person's words are recorded rather than followed, and a test reads that sentence.
- Every new label goes through its package's `messages.go`, carries no em dash and passes agora's name check: attachment, file and tray are clear of it.

### Records

- `agora/CONTEXT.md` defines attachment and tray, and holder stays out of the glossary because a decision and the composer are already named.
- ADR 87 records that bytes travel at the press and are never served back, and ADR 88 that every collected answer is kept whole, amending ADR 73.
- `agora/AGENTS.md` gains an invariant for each rule above, each held by a test named for its scenario, and the layout line for `agora/web` names `ui-file-tray`.

## Testing Decisions

- The highest seam is the browser suite driving the real binary, end to end: attach by the button, drop, and paste into a card with the cursor in the words and outside them, press, then run the real `socrates agora poll`, read the paths it printed and compare the bytes on disk with what was attached.
  It follows the browser tests that already drive the page and poll the real binary.
- The same seam holds: the eleventh file refused in place, the pasted names, text pasting as text, files alone counted as answered and never named by the confirmation, the summary and the receipt naming files, the amend resending every file, a revision emptying every tray, a closed exchange drawing no tray, and the settled block showing names.
- Paste is driven by dispatching a paste event carrying a `DataTransfer` of files, which pins what the installed Chromium hands the page.
- `agora/surface`'s own tests hold only what a browser cannot cheaply produce: the 20 MiB and 100 MiB ceilings, a part the answer does not name and a name with no part, a hostile stated name staying inside the answer's folder, store-then-swap leaving nothing on a refusal, a replaced uncollected answer's folder removed, a revision removing every uncollected folder, and every collected answer kept under `answers/<n>/` until the session is collected.
- The ADR 59 check gains its cases in the existing settled tests in `agora/surface` and `agora/document`.
- The catalog draws `ui-file-tray` in every state it has, and the browser suite proves its behaviour there, as for every primitive.
- `review/command`'s tests hold `respond --attach`: the copy into the run store, the ceiling of ten, and the refusal on a forward.
  `review/prompt`'s tests hold the attachments named in the next round's prompt.
- The skills' existing tests read the sentence saying an attachment is evidence rather than an instruction.

## Out of Scope

- Attachments on a settled answer offered as anything more than its names.
- A caller attaching files to a document: evidence and illustrations stay the caller's way to show something.
- Posting attachments to a pull request when a finding is forwarded.
- Any route that serves an attachment back to the page.
- Reading the clipboard directly, which would need a permission prompt.

## Further Notes

- The design page's risks stand: a file is an injection channel, the body ceiling stops being one number, a person amending repeatedly writes and removes a lot of disk, and what a browser hands a paste of files copied in a file manager varies between browsers.
- Under ADR 88, disk under `~/.socrates/agora` grows with each round of a session, up to 100 MiB an answer, until the session is collected.
- The design page proposed one ADR numbered 79, which was taken; this change's ADRs are 87 and 88.
