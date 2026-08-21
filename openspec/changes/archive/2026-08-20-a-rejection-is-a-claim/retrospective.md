# Retrospective: a-rejection-is-a-claim

## The live runs, and the dialogue that never happened

Task 9.4 asks for a real dialogue with a real reviewer, including one rejection the author genuinely believes in, and for this document to record whether the second reviewer accepted it and whether it should have.

There is no such record, because in four live runs no round ever parked.

| run | subject | findings | for the person | outcome |
| --- | --- | --- | --- | --- |
| `20260820T091048Z-6616db11` | the whole change, 40 files | 7 | 5 | escalated at round one |
| `20260820T093222Z-ea4284ba` | the whole change after the first round of fixes, 45 files | 7 | 4 | escalated at round one |
| `20260820T095643Z-43015b70` | one fix commit, 18 files | 8 | 1 | escalated at round one |
| `20260820T100854Z-80725392` | one fix commit, 12 files | 6 | 4 | escalated at round one |

Every run reached the reviewer, every reviewer returned a parseable document, and every one of them found at least one finding it actioned `ask-user`.
`ask-user` dominates the routing rule, so the run escalated to the person at round one and nothing was offered for a response.

That is the design working exactly as specified, and it is the most useful thing four live runs established: **on real work, the parked path is the exception rather than the rule.**
A reviewer reading real code finds things that touch what the author decided on purpose, and the moment it finds one the dialogue does not start.
The dialogue machinery in this change is driven end to end by the simulator suite, in four journeys, and by a real reviewer it has never been driven at all.

Two consequences worth stating rather than discovering later:

- The cost of this change is carried mostly by runs that never use it.
- Whatever tuning eventually makes the dialogue common lives in the prompt's action definitions, not in the routing rule. The reviewer is instructed to prefer `ask-user` when in doubt, and it does.

The rejection in 9.4 therefore remains unwritten, and the task stays open.
Manufacturing a defect to provoke a park, or writing a rejection I do not believe, would have produced a document that says the mechanism works when nothing had tested it.

## What four live reviews found that a green suite did not

Before the first run, `make fmt`, `make lint`, `make test` and `make scenario-coverage` were clean and 201 scenarios were named by tests.
The four reviews returned 28 findings between them.
Twenty-seven were real, and three of those were regressions introduced by the fix for a previous one.

The count is the point of this section.
A green suite, a coverage gate, an exhaustive property test and a convention checker did not catch any of what follows.

### The routing rule was wrong three times

It shipped as: nothing blocking **and nothing patchable** is clean, otherwise park until the cap.

The first review traced `{severity: info, action: auto-fix}`, which `findings.Validate` accepts and the prompt's own reporting section invites, since it says severity and action are independent.
Such a round blocks nothing, so `Unanswered` waits on nothing, so `scrat status` says `unanswered: 0` and "ready to continue" with nothing for the author to do; continue, get the same finding back, and the cap escalates a run whose findings all asked for nothing.

The specification said the same wrong thing in the same words, which is why the code matched it.
So did the exhaustive property test: it asserted `!blocks && !patchable` implies clean across all 512 subsets and passed, because it was written from the sentence the implementation was written from.
**A test derived from the same sentence as the code proves the sentence was implemented and nothing else.**

The fix made it "nothing blocking is clean", which the second review then showed was still incoherent: a `{warning, no-op}` finding holds the change while its action asks for nothing, so the run parked and demanded a response to a finding `scrat findings` labelled informational.
The fix for that was to park only on findings a response could address.

The fourth review then showed that fix had opened a hole: with `[{warning, no-op}, {error, auto-fix}]` the round parked on the patchable one, nothing ever waited on the other, and because `scrat findings` reports the latest round only, a blocking finding disappeared from the run entirely.
The rule now escalates on anything blocking that a response cannot address, whatever else is in the set, for the same reason `ask-user` does.

Three reviews, three corrections, one rule of five clauses.
It is four lines of code and it was wrong in a different way each time.

### A rejection reason could compose the prompt that judged it

`hygiene.Neutralise` removed conflict markers and collapsed whitespace, and left a line beginning `## ` alone.
`prompt.Sections` parses exactly that marker.

Removing the fix and rerunning the test that now guards it prints the attack:

```
got:  [... What earlier rounds found | What to report | Your task | How to review | Test quality | What to report]
want: [... What earlier rounds found | How to review | Test quality | What to report]
```

Two sections composed by the session under review, one of them a second `Your task`, both ahead of socrates' own instructions.
The requirement that nothing about a claim's placement gives it the standing of an instruction did not hold, and the test that was meant to hold it used an instruction-shaped reason with no heading marker in it.

The guard is now the sequence of section headings compared against the same prompt built without the injected fields, which does not depend on guessing which heading an attacker picks.

### A completed review could be made unreachable, twice

A parked run kept the process identifiers of processes that had already exited, and `Cancel` kills what the run records with no liveness check, so a parked run cancelled hours later would signal whatever held those numbers by then.

The teardown of the round's worktree ran after the parked state was published, and the next round checks out that exact path, so `scrat continue` raced a `chmod` walk and an `os.RemoveAll` over the whole tree and failed the run on a git error.
Reordering it fixed the race and introduced a worse failure: a teardown error now returned before the run was ever published, leaving `run.json` reading `reviewing` with the findings on disk, `--wait` polling forever, and `scrat cancel` as the only way out of a review that had succeeded.
The round now parks regardless, records the teardown failure against itself, and the next round clears whatever stands where its worktree goes, so the failure is recoverable rather than deferred.

That last path is the one thing in this change with no test.
Forcing a teardown failure needs a filesystem the harness cannot produce, and the recovery half is tested by leaving a directory where the next round's worktree goes, which is what the failure would leave behind.

### The author's claim could be dropped, silently, three ways

`scrat respond --fixed --reason "..."` built the response from `--note`, never looked at `--reason`, and recorded a fix with no note, after which the next reviewer was told the author claims to have fixed this and says nothing about how.
`--note` on a rejection went the same way.
Both are now refused, naming the flag that applies.

A rejection with no reason read standard input whenever it was not a character device.
In the test suite standard input is `/dev/null`, a character device, so the tests saw the intended refusal; from a script or an agent harness it blocks until the writer closes, or stores whatever unrelated bytes were in the pipe as the author's claim.

Recording a response was a read, an append in memory and a rewrite of the whole file, with no exclusion.
Two `scrat respond` invocations at once, which is the ordinary case for a session encouraged to run independent commands in parallel, lost whichever claim was written first.

The first attempt at that last one was a lock file with a liveness check on the holder, and the test written for it failed immediately: the stale-takeover path removes a lock file it does not own, so two processes could both hold it.
The document was already append-only with last-wins semantics, so the file is now append-only too, one JSON response per line, and there is no lock to get wrong.

### The rest

The verdict of a clean second round named the round cap, telling the author something was left for them when nothing was.
`Continue` treated a tree holding nothing the base does not as `failed` where `Begin` treats the identical condition as `clean`, and the reason it recorded said "identical to HEAD" on the one path where HEAD is not what was compared.
`Continue` spawned a reviewer without the review-inside-a-review guard, so a reviewer session could seat a reviewer of its own; the invariant is stated against the reviewer's environment rather than against one command, and now holds that way.
This change's own `cli-foundation` delta declared `watch` and `logs`, which belong to a change nobody has started, and would have entered the promoted specification as satisfied.

### The test-quality findings, which are the ones worth reading twice

The reviewer found the vacuous tests I wrote while fixing the reviewer's earlier findings.

The concurrency test recorded eight responses in parallel and asserted only that responses which succeeded were present, skipping the ones that failed and never asserting a count, so a run where all eight failed was a pass.

The replacement assertions built a `Responses` literal holding two answers to one finding and then asserted that it held two answers to one finding.

And the new oracle in the exhaustive routing test restated `dialogue.Answerable` line for line, which is the exact flaw this document had already named two sections above.
The oracle is now written from literal severities and actions rather than from the helpers under test.

Socrates found, in socrates' own tests, the failure mode socrates' test-quality rule exists to name, in code written by the author who had just finished writing that the failure mode was what let the first defect through.

## What the design got right

`ask-user` dominating is the rule that fired, and it fired on every real change it saw.
Fourteen of the twenty-eight findings across four runs were actioned `ask-user`.
A tool that routed those through a machine dialogue would have argued with intent fourteen times before a person saw any of it.

Keeping every round's artifacts and recording each round's snapshot in `run.json` cost nothing and made all four runs reconstructable afterwards, exactly, from `reviews`, `events.jsonl` and the round directories.

The seal held on every run, with the machine's enterprise tool servers loaded and `agent-tool-servers-loaded` recorded each time, and each reviewed clone was byte-identical afterwards.

## What this says about reviewing with socrates

The loop converged, but not quickly, and never by agreement.
Round one of run three reviewed the fixes from run two and found eight new things.
Run four reviewed the fixes from run three and found six, four of them about decisions the previous fixes had made.

Three of those were regressions: the reorder that fixed a race created a dead end, the park condition that fixed a contradiction dropped a finding, and the append that fixed a lost update turned a stored claim into a reported failure.
Every one of them was found by the next reviewer, and none of them by the suite.

The honest summary is that a fresh reviewer per round on real work is worth its cost, and that the thing it is best at is finding what the previous fix broke.
That is an argument for the round cap being about context rather than confidence: two rounds is not enough for a change of this size to converge, and the point of stopping is that a person decides, not that socrates has finished.

## How task 9.4 closed

It closed on being reworded, not on a dialogue having happened.

The task asked for a live dialogue including one genuine rejection.
Four live runs later there was still no dialogue, because every round returned an `ask-user` finding and escalated, and no amount of working could change what the model chose to answer.

`what-a-reviewer-actually-answered` is where that became a rule: a task about a live run may require the run, the recording and the routing, and may not require a severity, an action, a count or a shape.
Task 9.4 now asks for what a live run can produce, and the four answers that could not satisfy the old wording are the seed of the corpus.

The rejection is still owed, and it is owed by a scenario replaying a recorded answer rather than by another live run.
