**The page half is driven by a real browser**, against the real binary on a real loopback address. Asserting on served markup would test the renderer and claim to have tested the page.

**A test that reloads before asserting proves the weaker claim.** Nothing added here may reach the behaviour through a reload except the scenario that is about reloading.

**Each guard is mutated on purpose before it is trusted:** the test that proves the new behaviour is run against the old code first, and the reproduction is confirmed to fail before anything is fixed.

**Nothing in agora knows what a finding is.** Review, finding, verdict, round and run appear nowhere in that binary.

**`make lint test` after every task, not at the end of a group.**

## 1. The reproduction, before anything is fixed

- [x] 1.1 Drive a real browser: open a page, commit an answer, finish, then write a second document for the same session, and confirm the open tab does not take it up. Confirm it fails for that reason and not for a slow wait
- [x] 1.2 Confirm the same tab takes the second document up after a reload, so the reproduction names the page rather than the server
- [x] 1.3 Confirm dropping the client's close of the stream does not on its own fix 1.1, so the server half is not skipped on a guess

## 2. A stream says where the exchange stands when it connects

- [x] 2.1 State the revision the exchange is at when a stream connects onto one that is not over, without waiting for a change. Covers `A stream is told where the exchange stands when it connects`
- [x] 2.2 Write a document while nothing is streaming, then open a stream, and confirm it is told the current revision rather than waiting for a change that has already happened. Covers `A stream that missed a revision catches up on connecting again`
- [x] 2.3 Confirm a stream opened onto an exchange that is already over is still told so on connection, and is told that rather than a revision. Covers `A stream opened after the end is told at once`
- [x] 2.4 Confirm a stream still reports each later change, and that a page reading one is not told the same revision twice for one change
- [x] 2.5 Mutate the stream to say nothing on connection, and confirm 2.2 fails

## 3. The page keeps listening after the exchange is over

- [x] 3.1 Stop the page closing its stream when it is told the exchange is over, and confirm the closed state still stops every control and still says the exchange is over. Covers `A reloaded page renders closed` and `A page opened for the first time after the end renders closed`
- [x] 3.2 Take up a later document on the page the person is looking at, by the path a re-render already takes rather than a second one
- [x] 3.3 Confirm the reopened page no longer says the exchange is over and that every control the new content carries is operable. Covers `The page stops saying the exchange is over`
- [x] 3.4 Confirm an answer the person committed but did not send is still not carried into content that no longer holds that decision, over an ended exchange as over an open one
- [x] 3.5 Regenerate the committed bundle from `web/` and confirm the suite fails while the generated file and its source disagree
- [x] 3.6 Mutate the page to close its stream again, and confirm 1.1 fails

## 4. The scenario is proved without a reload

- [x] 4.1 Rewrite the test named for the scenario so it asserts on the open page, with no reload, and confirm an answer sent from it is accepted and delivered. Covers `A new document reopens an ended page`
- [x] 4.2 Keep the reload as a scenario of its own, so removing the reload does not remove a guarantee. Covers `A reloaded page still reopens`
- [x] 4.3 Confirm a re-render onto an exchange that never ended still reaches the open page unchanged. Covers `A re-render reaches an open page`
- [x] 4.4 Confirm no test of this behaviour reloads except the one that is about reloading

## 5. The map is brought back into step

- [x] 5.1 Update `AGENTS.md` where it states what the surface does about delivery and waiting, if what it says has stopped being true
- [x] 5.2 Run `make fmt` and `make verify`, and confirm scenario coverage names every scenario this change adds
