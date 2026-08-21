## Why

A finding actioned `ask-user` ends the run, so the person's decision has nowhere to go.
The author decides, patches what was decided, and there is no round left to judge the patch: `respond` and `continue` both refuse a terminal run, and a fresh `scrat review` is a new run whose reviewer is told nothing about the question that was asked or the answer that was given.

The rule's own stated reason is narrower than the behaviour it produces.
It escalates because "parking on the rest would advance the run past a finding nothing in the dialogue can close", and that is true of a machine answering an intent question.
It stops being true the moment a person has answered one, which is exactly the case the escalation was built to reach.

## What Changes

- **A round holding an `ask-user` finding parks instead of ending.** The run reaches `awaiting-response` rather than `escalated`, and only the person's answer moves it. **BREAKING**: runs that exited `1` for an intent question now exit `4` until they are answered.
- **A third resolution records a decision.** `scrat respond --decided --reason` answers a finding that asked for one. `--fixed` and `--rejected` keep their meanings and remain refused for a finding that needs a decision, so the author cannot answer an intent question by asserting a patch.
- **A decision travels as a claim, like a rejection.** The recorded answer says a person decided it, the reason is required, and both reach the next reviewer, which judges whether the decision and whatever was patched actually resolved the finding.
- **What a round waits on distinguishes a patch from a decision.** `scrat status` and `scrat findings` say how many of the unanswered findings need a person, so a caller that must stop and ask can tell without parsing prose. The exit-code table does not change.
- **`escalated` narrows to what it always meant: socrates is finished and nothing more it can do.** The round cap reached with something still blocking, and a question the person left unanswered when the cap ran out.
- **The skill teaches the branch.** A parked round holding a question is where the session presents to the person and waits, rather than answering it and continuing.
- **The skill teaches that a round is spent without asking, and that there is no third one.** Once every answerable finding is answered the session continues, because a patch is unreviewed new code and judging it is the point. Where the cap escalates, a further look is a fresh `scrat review` rather than a third round, and only where the person asks for one.

## Capabilities

### Modified Capabilities

- `review-dialogue`: what a round leads to when it holds an `ask-user` finding, what the set of answerable findings contains, the third resolution and the reason it requires, and what the next round is told about a decision.
- `cli-foundation`: `scrat respond` gains `--decided`, and `scrat status` says how many of the unanswered findings need a person rather than a patch.
- `run-store`: what `awaiting-response` means, which is no longer only that what the round found is mechanical.
- `agent-skill`: the `awaiting-response` branch splits into what the session may answer and what it must hand to the person, the prohibition against answering an intent question moves from being enforced by a terminal state to being stated, and the skill states who decides that a round is spent.
- `reviewer-prompt`: the labelled prior-round block carries a decision as a decision, so the next reviewer knows a person answered rather than that the author claimed a fix.

## Impact

- `internal/dialogue`: `After` stops treating `ask-user` as dominant over parking, `Answerable` admits findings that need a decision, and `Resolution` gains a third value with a required reason.
- `internal/report`: the counts a parked run prints, and the guidance it hands over.
- `internal/reviewer`, `internal/prompt`, `internal/skill`, `cmd/scrat`: recording a decision, carrying it into the next round's prompt, and teaching the branch.
- The corpus is unaffected as stored bytes and affected as routing: every recorded answer holding an `ask-user` finding reclassifies from `escalated` to `awaiting-response`, which is a deliberate diff in `internal/corpus/outcomes/distribution.txt` and the first time that file will show a parked outcome any real reviewer reached.
- Nothing about what socrates writes in a reviewed repository changes, and no new writer appears: the author still applies every patch.
