**What a document accepts is decided over its input space**, not sampled: a document with decisions, one with none, and one with nothing at all.

**The page half is driven by a real browser**, against the real binary on a real loopback address. Asserting on served markup would test the renderer and claim to have tested the page.

**The loop half is driven end to end through the real binaries**, so what the skill teaches is proved by running it rather than by reading it.

**Nothing in agora knows what a finding is.** Review, finding, verdict, round and run appear nowhere in that binary, and the entries the page lists are entries of a source.

**Each guard is mutated on purpose before it is trusted:** the test that proves the new behaviour is run against the old code first.

**`make lint test` after every task, not at the end of a group.**

## 1. A document may hold nothing to decide

- [x] 1.1 Accept a document holding no decision, and serve it. Covers `A document with nothing to decide is served`
- [x] 1.2 Refuse a document carrying no decision, no source entry and no preamble, saying the page would hold nothing and starting no server. Covers `A document with nothing at all is refused`
- [x] 1.3 Confirm a document that does hold decisions is validated exactly as before, over the same failure space: no option, an option with no command, two options marked. Covers `A document that asks something is validated as before`
- [x] 1.4 Mutate the floor to accept a document carrying only a title, and confirm 1.2 fails rather than passing quietly

## 2. The source is what the page shows of the whole

- [x] 2.1 Render every entry of the source as a short list ahead of the decisions, each under its name with its severity, audience and location. Covers `The source is read before the decisions`
- [x] 2.2 Show an entry no decision cites, and confirm one a decision covers leads to that decision rather than being rendered twice. Covers `An entry no decision cites is still shown` and `An entry leads to the decision that covers it`
- [x] 2.3 Confirm what the list holds comes from the entries rather than from the preamble the caller wrote. Covers `The list is not the caller's prose`
- [x] 2.4 Render a document citing no source and confirm no list is shown and the page is what it is today. Covers `A document citing no source lists nothing`
- [x] 2.5 Confirm the convention check still passes: nothing added here names review, finding, verdict, round or run as a noun

## 3. A page with nothing to answer

- [x] 3.1 Offer no way to answer and no way to send where the document holds no decision, and enumerate the rendered controls to prove it. Covers `A page with nothing to answer offers neither control` and `A page with nothing to answer offers neither of the two`
- [x] 3.2 Say in words that there is nothing to answer and that nobody is waiting, and read it back from the page a browser loaded. Covers `The page says nobody is waiting`
- [x] 3.3 Confirm such a page carries the same frame, palette and layout as one holding decisions, and renders the source the same way. Covers `The frame is the same frame`
- [x] 3.4 Confirm a page that does hold decisions still offers both interactions and nothing else. Covers `The surface offers no third interaction`
- [x] 3.5 Regenerate the committed bundle from `web/` and confirm the suite fails while the generated file and its source disagree
- [x] 3.6 Mutate the renderer to emit the send control on a decisionless page, and confirm 3.1 fails

## 4. socrates says which rounds wait

- [x] 4.1 Rename the rule for what it decides, keeping its signature, its table and the property that nothing re-derives it. Covers `An escalated round waits on a person`, `A parked round holding a question waits on a person`, `A parked round the author answers alone waits on nobody` and `A clean round waits on nobody`
- [x] 4.2 Exercise the rule across every combination of classified findings and every state a run can hold, and confirm it agrees with the routing rule and answers that nothing waits in a state that rule cannot reach. Covers `The rule is decided over its whole input space`
- [x] 4.3 Print a line naming the wait and carry `awaits_a_person` in the object `--json` prints, from the one rule that decides it. Covers `A round that waits on a person says so`
- [x] 4.4 Confirm a clean run and a parked round the author answers alone say nothing about waiting, and say nothing that would stop the round being handed over as a page. Covers `A round the author answers alone says nothing about waiting`
- [x] 4.5 Confirm `socrates status` carries neither the line nor the field, and prints what it printed before. Covers `scrat status does not carry the wait`
- [x] 4.6 Confirm no command, no flag and no run-record field was added, and that socrates states nothing about whether a page is served. Covers `The surface does not grow`
- [x] 4.7 Confirm what socrates itself reaches is unchanged and no package moved behind a build tag. Covers `socrates reaches nothing new` and `A second binary in the module is not socrates reaching further`
- [x] 4.8 Mutate the rule to answer that a clean round waits, and confirm 4.2 and 4.4 fail

## 5. The skill hands every round over

- [x] 5.1 Rewrite the loop's branches so `clean`, `awaiting-response` and `escalated` each draft a document, serve it, hand over the address and say briefly what it holds, and so `failed` and `cancelled` say nothing about a page. Covers `Every branch that has a result says to serve the page`
- [x] 5.2 Say to wait on `agora poll` only where socrates reports the round waits on a person, and to carry on without waiting where it does not. Covers `The skill says which branches wait` and `A round the author answers alone gets a page nobody waits on`
- [x] 5.3 Say once, where the loop is described, that the printed address holds the session and that `agora poll` is given it, and confirm no branch restates it. Covers `The skill says to carry the session from serve to poll` and `The session is stated once rather than on every branch`
- [x] 5.4 Say the session is kept for the whole review and that a later round is written to the page the person already has. Covers `A later round reaches the page the person already has`
- [x] 5.5 Say what a page holds where the round holds nothing a person answers, and that the short list comes from what socrates reported rather than from the session. Covers `A round with nothing to answer is a page with nothing to answer` and `The decisions are in full and the rest is a list`
- [x] 5.6 State where the document is drafted as what follows the round: inline where nothing does, delegated where another round does. Covers `The escalated branch drafts inline`, `The clean branch drafts inline` and `The parked branch delegates`
- [x] 5.7 Confirm the driving session is the one that serves and the one that waits, and that nothing about the page appears in the waiting step. Covers `The driving session owns the surface` and `The page is not in the polling loop`
- [x] 5.8 Confirm the skill takes the fact from what socrates reports and states no condition of its own over the state or the counts. Covers `The session does not decide the condition itself`
- [x] 5.9 Confirm the bar appears once, before the branches that serve a page, and that no branch restates it. Covers `The bar is stated before the branches that refer to it`
- [x] 5.10 Confirm the skill still says how an answer reaches the run, that the page and the terminal use the same identifiers, and that a decision is still the person's. Covers `The skill says how an answer reaches the run`, `The skill says the page and the terminal use the same identifiers` and `A decision is still the person's to make`
- [x] 5.11 Confirm the skill names the local surface and no host, upload, publish or artifact outside the machine. Covers `A page is served locally and published nowhere`
- [x] 5.12 Confirm nothing socrates spawns is asked for a summary written for a person, and that nothing about a page reaches any prompt socrates builds. Covers `Nothing socrates spawns drafts a page`
- [x] 5.13 Regenerate the committed skill from `internal/skill` and confirm the suite fails while the generated file and its source disagree

## 6. The loop, driven end to end

- [x] 6.1 Drive a clean run, an escalated run and both flavours of parked round through the real binaries, and confirm each is served as a page. Covers `A round that completed a review is handed over as a page`
- [x] 6.2 Drive a failed run, and a run stopped after a round that completed a review, and confirm no page is served and the reason is what the session reports. Covers `A run that ended without a result is handed over as nothing`
- [x] 6.2a Prove the place that check reads is the place a page lands: serve one into it and confirm it arrives, so the check can fail rather than passing over an empty directory
- [x] 6.3 Drive a round that waits on a person and confirm the session waits for what the person sends before anything is recorded. Covers `A round that waits on a person is waited on`
- [x] 6.4 Drive a clean run and a parked round the author answers alone, and confirm the session hands over the address and carries on. Covers `A round that waits on nobody is handed over and not waited on`
- [x] 6.5 Leave the page of a round that waits on nobody unopened, and confirm the loop patches, records and continues, or reports and stops, as it would with no page. Covers `A page nobody waits on does not stop the loop`
- [x] 6.6 Mutate the loop to wait on every page, and confirm 6.5 fails rather than hanging the suite
- [x] 6.7 Mutate the loop to serve a page for a run that ended without a result, and confirm 6.2 fails

## 7. The map is brought back into step

- [x] 7.1 Update `AGENTS.md` where it names the rule and what socrates states about a page, and confirm the invariants it lists are the ones tests hold
- [x] 7.2 Run `make fmt` and `make verify`, and confirm scenario coverage names every scenario this change adds
