## Context

socrates hands findings over in a terminal, in the shape the reviewer returned them: identifier, severity, path, line, one sentence.
That shape is right for the session, which reads a finding in order to patch the file it names, and wrong for a person, who reads a finding in order to decide something.

Two branches of the loop put the person in the reader's seat.
A run that escalated is finished, and `dialogue.After` reached that state because a blocking finding's own action asks for nothing or because the last permitted round still found something blocking.
A parked round holding a finding that needs a person cannot move until they answer, and the skill already says the session does not answer one on its own behalf.

Three constraints bound how this can be built.

The reviewer cannot draft the page.
It is sealed read-only inside the snapshot worktree, it runs under `--safe-mode`, and its output is the untrusted thing under judgement.
A reviewer that also composed the person's summary of its own findings would be arguing its case in the document the person decides from.

socrates cannot draft the page either.
A Go template over the findings document can restate what the reviewer wrote in a nicer typeface, and the thing being asked for is the opposite: plain English about what a finding means for the change, with an example where the finding is purely technical.
That requires a model, and socrates' agent seam exists to review, not to write for a person.

The session driving the loop already has what is needed: it holds the whole implementation context, it has the `artifact-design` skill, and it is the one thing in the system that is talking to the person.

## Goals / Non-Goals

**Goals:**

- One rule, in `internal/dialogue`, deciding which rounds a person reads, over the same input as the rule it sits beside.
- `scrat findings` stating it, so nothing re-derives it.
- The skill instructing the page, holding it to a stated bar, and putting the drafting where it costs least.
- The end-to-end journey asserting the signal appears on exactly those branches.

**Non-Goals:**

- socrates rendering, publishing, storing or verifying the page. It states a fact about a round and stops.
- Recording the page's URL in the run. That is a new field, a new flag to write it, and a new thing to be stale, in exchange for a link the session has just printed.
- Replacing what the session reads. `scrat findings` keeps printing findings, because the session patches from them.
- A page for a `clean` run, or for a parked round the session can answer by itself.
- Anything about `scrat status`.
- An interactive page. Nothing is behind it, so a control on it would collect an answer that reaches no run, and building one spends the tokens and the time the page exists to save.

## Decisions

### The rule sits beside `dialogue.After`, and takes the state that rule reached

A predicate on the classified findings and the state, in the package that reaches no filesystem, no clock and no process, tested against its whole input space rather than sampled.

It takes the state rather than the round number, which was the first shape and is the wrong one.
A round number makes the predicate re-decide the routing, and everything that reads it holds a state that was already decided and recorded, so the two can disagree: a stored `escalated` whose findings would now route somewhere else answers one way through the record and another through the rule.
Taking the state reads that one decision's result instead, and it makes the predicate total, so `failed`, `cancelled` and a run still reviewing need no case anywhere else.

The alternative is the skill naming the condition in prose: render when the state is `escalated`, or when it is `awaiting-response` and the count of findings needing a person is above zero.
Both facts are already in `scrat status`, so the skill could say it with no Go change at all.
That is the second reader this repository fails the suite over.
`findings.EffectiveAction` has an abstract-syntax check for exactly this shape, and the reason generalises: two places computing who a round is for can disagree, and the disagreement is silent, because a round the skill thinks is the session's own is a round nobody renders and nobody notices was not rendered.

### `scrat findings` carries it, and `scrat status` does not

`status` is polled in a loop and every line it prints lands in the session's context, which `cli-foundation` calls the scarcest resource in the system.
`findings` is read once, on exactly the two branches where the fact matters, and is already where the response rules print.

### The bar is stated once, in the skill

The same pattern as `dialogue.ResponseRules`: one source, and the branches refer to it rather than restating it.
Two copies of the bar drift in the predictable direction, and the copy that gets softened is the one in front of the session.

It does not print from `scrat findings`.
The response rules print there because they govern a command the session is about to run and get wrong; the bar governs a document the skill has already told it to make.

### Inline at `escalated`, delegated at a parked round

The asymmetry is in what a token costs at each branch, not in the work being different.

At `escalated` the run is over. Loading `artifact-design` and drafting the page is paid once, and a fresh subagent's own system prompt and tool schemas make delegating it strictly more total tokens for nothing.

At a parked round the session is about to patch each finding, run one `scrat respond` per finding, run `scrat continue`, poll, read round two's findings and report.
Everything added to that context is re-sent on every one of those turns, and it is added at the exact point `CapIsForContextRatherThanCost` names as where the session starts compacting mid-review and quietly loses the thing it was defending.
Delegating means the driving session pays a prompt and a link, and the subagent reads `scrat findings --json` itself.

The alternative of one uniform rule was considered in both directions.
Delegating everywhere costs the fresh-agent overhead on the branch that cannot benefit from it.
Inlining everywhere spends the context the cap exists to protect, to save one sentence of skill text.

### What the page holds

The findings a person must decide, in full, and the rest of the round as a short list.
A parked round holding one question and five mechanical defects is one decision and five things the session is already fixing, and a page that treats all six evenly buries the only thing the person has to act on.
The five are listed rather than dropped, because the person cannot see what else was found without them.

### The page is read only, and the identifiers are the join

The person's path is: read the page, decide, go back to the terminal and say what they decided.
That makes the page one half of a pair, and the pair only works if the halves are matchable without translation.

So every finding on the page keeps the identifier, severity and audience `scrat findings` printed, and each decision says what answering it looks like in the terminal, naming the finding.
A page that renamed a finding to something more readable would hand the person a search problem in exchange for a nicer sentence, which is the problem the page was made to remove.

Nothing on the page collects anything: no form, no button, no checkbox, no stored state, no declared capability.
There is no backend, so a control would take an answer that reaches no run, and the person would reasonably believe they had answered.

### The journey asserts the signal, not the page

Publishing reaches the network, and `internal/fidelity` is where anything that reaches the network lives, behind a build tag.
So the end-to-end journey drives the loop as it does now and asserts that `scrat findings` states the fact on the escalated branch and on the parked-with-a-question branch, and does not state it on a clean run or on a parked round whose findings are all patchable.
That the session then made a page is not something a test can drive, and claiming otherwise would be a test that reads as passing while checking nothing.

## Risks / Trade-offs

**The page is cooperative, and socrates cannot tell whether one was made.** → Stated rather than papered over. socrates' half is a fact in a command's output and a sentence in a generated skill, both of which fail the suite when they go missing. A session that ignores the skill is the same exposure the loop already has, and the terminal output it would fall back to is what exists today.

**A findings document quotes the code it was written about, and publishing puts it on a hosted page.** → It is the session's reach, in the harness the person is already driving, and the page is private until they share it. socrates gains no network reach, no dependency and no build tag. `internal/corpus` refuses to record from any repository but socrates for this same reason, and that refusal is untouched: the corpus and the page are different documents going to different places. The asymmetry between the two treatments is recorded as https://github.com/rui-ktei/socrates/issues/1 rather than settled here.

**Plain English about impact is a model summarising a reviewer.** → A summary can soften a finding. The page is a reading aid and not the record: `scrat findings` still prints what socrates concluded, `findings.jsonl` still holds what the reviewer said, and the resolution the person records travels to the next reviewer as a claim it judges. Nothing about the page enters the prompt.

**One more thing in the skill, on the two branches that already carry the most text.** → The bar is stated once and the branches gain a sentence each. The surface check still holds, because the page names no scrat command that does not exist, and the polling invariant still holds, because nothing here goes in the loop.
