**Nothing here changes what is written.** The log format, the run states, the verdicts and the exit codes belong to the core, and no task in this change edits one. A task that finds itself wanting to is a task in the wrong change.

**`make lint test` after every task, not at the end of a group.**

## 1. The renderer, as a pure function

- [x] 1.1 Add the renderer package taking a sequence of events and producing lines, and the test that its package imports nothing outward. Covers `Rendering needs no terminal`, `The renderer decides nothing about the run`
- [x] 1.2 Render each event kind the core writes, and render an unrecognised kind legibly rather than dropping it. Covers `An unknown event kind is still shown`
- [x] 1.3 Render the rotation record. Covers `A rotated log is marked`
- [x] 1.4 Assert that nothing outside the renderer produces a line a person reads, so a second rendering cannot appear beside it

## 2. Following

- [x] 2.1 Add the tail over the event log, and the two stopping conditions. Write the empty-log test first: waiting only for a terminal event is the version that hangs, and it passes every other test. Covers `A run that failed before writing anything`
- [x] 2.2 Add `scrat watch` and its exit codes. Covers `Following a review to a clean verdict`, `Following a review to an escalation`, `A follower does not outlive the run`, `A follower stops where the round stops`
- [x] 2.3 Add the cancellation path. Covers `Cancelling while following`
- [x] 2.4 Add the bounded wait for a log that does not exist yet, tested with `testing/synctest`. Covers `Watching immediately after starting`, `Waiting is bounded`
- [x] 2.5 Add the immediate refusal for a run identifier that does not exist. Covers `A run identifier that does not exist`. It is the case the bounded wait would otherwise swallow for the length of its patience

## 3. Reading

- [x] 3.1 Add `scrat logs`. Covers `Reading a finished run`, `Reading a run still in progress`
- [x] 3.2 Add the terminal detection and the pager, with the plain fallback. Covers `Piping works`, `A pager that cannot be started`

## 4. The boundary between the three

- [x] 4.1 Update the command surface and its test. Covers `The command surface is what it says`
- [x] 4.2 Add the tests that hold the boundary. Covers `Status stays small when the reviewer is loud`, `Neither command grows the other's flag`, `All three agree about the run`

## 5. End to end

- [x] 5.1 Drive a run to each verdict with a follower attached, and assert the follower exits with the same code the run's own verdict gives
- [x] 5.2 Drive a follower attached before the reviewer's first event, and assert it follows normally afterwards
- [x] 5.3 Add the `testscript` scripts for `watch` and `logs`, including piped output and a missing run

## 6. What a person reads

- [x] 6.1 Any new sentence goes in its package's own registry, and the plain-language scan passes with nothing added to a permit list
- [x] 6.2 Update `README.md` with the three commands and which is for whom

## 7. Verification

- [x] 7.1 `make fmt && make lint test` clean
- [x] 7.2 `make scenario-coverage` naming no scenario this change adds
- [x] 7.3 `make scenario-drop` before archiving, because the command surface requirement is replaced by a `MODIFIED` delta and a scenario the delta does not carry leaves the specification silently
