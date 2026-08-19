## Why

The core hands findings back and stops.
That is a linter with a very expensive model behind it, and it has the linter's failure mode: the author reads the findings, decides four of the six are wrong, fixes two, and nobody ever checks either half of that judgement.

The interesting property of socrates is not that a reviewer produced findings.
It is that the author, who has all the context, gets to disagree, and that the disagreement is judged rather than assumed.

There is one rule to get right, and getting it backwards makes this layer worse than not having it.

`no-mistakes` opens its fixer prompt with *"always start with double checking whether the findings are legitimate."*
That is safe advice for a stranger and dangerous advice for an author.
The session being handed findings here wrote the code; told it is wrong, it will defend.
The same sentence that stops a neutral fixer thrashing licenses an author to rationalise, and the result is a review the author always wins.

So invert it.
The author may reject a finding, and a rejection is never silent dismissal: it is a written response, attached to that finding's identifier, that goes back to a fresh reviewer, which accepts it or escalates to the person.
That single rule is the difference between a dialogue and a linter with extra steps.

## What Changes

- **A run can be awaiting the author.** Where a completed review found something blocking and nothing that challenges intent, the run parks awaiting responses instead of escalating.
- **`scrat respond` records a fix or a rejection against a finding.** A rejection carries a reason and the reason is required. A response naming a finding that does not exist is refused.
- **`scrat continue` starts the next round.** It re-snapshots, so the reviewer sees the patches as commits it can diff rather than as a description of what changed.
- **The next reviewer is a fresh process that receives the prior findings and the written responses, labelled as claims to be judged.** This is the single, explicit exception to the rule that the reviewer never receives the author's account of the change.
- **The re-review carries the provenance clause.** The code under re-review was written by the change's author in response to findings. It is unreviewed new code, not a settled resolution. A test written in the same round as the code it exercises is part of the claim.
- **Anything that challenges intent escalates immediately and is never handed to a machine to argue with.** An `ask-user` finding stops the loop at the round it appeared in.
- **The round cap is two, and it is for context rather than for cost.** The session being handed findings is already large from implementing the feature, and each round adds findings, patches and rebuttals to it. Round three is where it starts compacting mid-review and quietly loses the thing it was defending.
- **Everything re-embedded in a prompt is neutralised first.** socrates now has two channels of model output feeding the next prompt where the core had one, so this becomes its own capability rather than a helper.

## Capabilities

### New Capabilities

- `review-dialogue`: rounds, the awaiting-response state, recording a fix or a rejection, advancing to the next round, the round cap, and what escalates immediately.
- `feedback-hygiene`: neutralising model output before it is re-embedded in a prompt, per field, on both channels.

### Modified Capabilities

- `run-store`: a run gains the `awaiting-response` state and holds one round directory per round rather than one.
- `snapshot`: every round is snapshotted afresh, and the base a round is reviewed against is the original base rather than the previous round's snapshot.
- `reviewer-prompt`: the prohibition on the author's rationale gains its one explicit exception, and the re-review carries the provenance clause.
- `cli-foundation`: the command surface gains `respond` and `continue`, and `findings` gains the response guidance.

## Impact

- One new package for the dialogue: the round decision, which is a pure function of the classified findings, the responses recorded, and the round number.
- The responses document joins the findings document as a per-round artifact.
- No change to the agent seam, to the verdict vocabulary, or to the exit codes. A run that used to escalate at round one now may escalate at round two, and `escalated` still means the same thing.
- The reviewer is still spawned per round and is still a fresh process, so the cold-reviewer property costs nothing extra here.
