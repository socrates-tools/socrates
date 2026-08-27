**The page half is driven by a real browser** against the real binary, as the rest of it already is.
Asserting on served markup would test the renderer rather than the page.

**Each guard is mutated on purpose before it is trusted:** the test that proves the new behaviour is run against the old code first, and is expected to fail there.
All three defects render as a live page or a silent discard, so a test that passes before the fix is a test that is asserting nothing.

**`make lint test` after every task, not at the end of a group.**

## 1. The closed state is something a page can be served in

- [x] 1.1 Drive a browser to end an exchange, reload the page, and confirm it renders live and offers its controls. Establishes the defect against `A reloaded page renders closed`
- [x] 1.2 Carry whether the exchange is over in the render frame, and render the closed state from it
- [x] 1.3 Have `page` read `Over` and pass it in
- [x] 1.4 Confirm a page loaded for the first time onto an ended exchange renders closed. Covers `A page opened for the first time after the end renders closed`
- [x] 1.5 Take the frame field back out and confirm 1.1 and 1.4 both fail

## 2. A stream onto an ended exchange is told at once

- [x] 2.1 Open a live stream onto an already ended exchange, change nothing, and confirm nothing arrives. Establishes the defect against `A stream opened after the end is told at once`
- [x] 2.2 Check `Over` before waiting on `Changed` rather than after
- [x] 2.3 Confirm a stream onto an open exchange still reports revisions and still reports the close when it comes

## 3. Finishing over an uncommitted answer asks first

- [x] 3.1 Select options in a decision, leave it uncommitted, finish, and confirm the exchange ends with nothing delivered and nothing said. Establishes the defect against `Finishing over an uncommitted answer asks first`
- [x] 3.2 Have the finish path name the decisions holding uncommitted selections or words, say what will not reach the caller, and wait for confirmation
- [x] 3.3 Confirm confirming ends the exchange and delivers no answer for the uncommitted decision. Covers `Confirming carries nothing that was not committed`
- [x] 3.4 Confirm declining leaves the exchange open with the selections intact. Covers `Declining leaves the exchange open`
- [x] 3.5 Confirm a finish with every decision committed, and one with no decision touched at all, both proceed without asking
- [x] 3.6 Confirm an uncommitted selection is still never delivered by either send. Covers `Choosing sends nothing`
- [x] 3.7 Remove the confirmation and confirm 3.1's replacement fails

## 4. A new document reopens an ended exchange

- [x] 4.1 End an exchange, revise it with a new document, answer from the reloaded page, and confirm the answer is refused. Establishes the defect against `A new document reopens an ended page`
- [x] 4.2 Have `Revise` clear `over`
- [x] 4.3 Confirm `serve` onto an ended surface returns an address that accepts an answer, and that `poll` delivers it. Covers `Serving again onto an ended exchange yields an answerable page`
- [x] 4.4 Confirm serving onto an open surface is unchanged. Covers `Serving again while the exchange is open is unchanged`
- [x] 4.5 Confirm a poll over an ended exchange still leaves it ended. Covers `A poll over an ended exchange leaves it ended`
- [x] 4.6 Confirm delivery still happens once, and that ending is still distinguishable. Covers `Delivery happens once`, `Ending is distinguishable`

## 5. The journey the defects were found in

- [x] 5.1 Drive the whole reported sequence end to end: a person finishes without committing, is told what would be discarded, declines, commits, and sends, and the caller receives it
- [x] 5.2 Drive the recovery sequence: a person finishes with nothing, the caller serves the same document again, and the person answers from the page at the address `serve` printed
- [x] 5.3 Confirm no wording added by this change carries an em dash, and that every new string lives in its package's `messages.go`
