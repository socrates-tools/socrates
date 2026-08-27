## Why

The reviewer's rules are the ones that were paid for, and three of them were paid for somewhere else and never collected.

An audit of `no-mistakes`' review step against socrates' delivered prompt found three rules that are in force there and absent here.
None of them is a preference.
Each was arrived at through an incident in the tool socrates' reviewer rules were taken from, and each names a class of defect the current prompt asks nobody to look for:

- socrates carries the test-quality rule and never tells the reviewer to report a test that breaks it.
  The rule defines what a worthless test is, and the reviewer is left to infer that finding one is worth a finding.
- socrates asks the reviewer to notice where the code and the stated intent disagree, and asks nothing of it when they do.
  A change that quietly drops a behaviour the person asked for is the case the intent exists for, and "something to report" is not an instruction that survives a reviewer deciding the change is otherwise clean.
- socrates asks whether a change is wrong and never asks whether the bug it claims to fix is still reachable.
  A fix that patches the symptom it was shown passes a review that traced the new code and never looked at the sibling path.

The test-quality gap is the one with a recorded cost in this repository already: `internal/convention` exists because tests that read source were written here, and the reviewer that read them was carrying the rule that forbids them.

## What Changes

- **The reviewer is told to report a test that breaks the test-quality rule.**
  The rule stays one block shared by every role that judges tests.
  The instruction to act on it is a second block, delivered only to the reviewer, and bounded: a same-pattern test inside the change is the change's to answer for, and the rest of the repository is not.
- **A contradiction between the change and the stated intent is an obligation, not an observation.**
  Where the change omits a behaviour the intent names, or adds one it rules out, the reviewer states the criterion and the contradicting code and asks a person, and does so even where the change is otherwise clean.
  It is never actioned as something the author can patch alone, because what the change is for is not the author's to settle.
  Silence in the intent is not a contradiction: an intent that does not mention something has not ruled it out.
- **A change that claims to fix a bug is asked whether the bug is still reachable.**
  The reviewer establishes whether the change claims a durable fix or a stated containment, and for a durable fix reconstructs the failing sequence and the invariant it needs, reads the sibling paths and the shared state, and reports a concrete reachable path where the source shows one.
  A containment the author stated is not blocked for being one.

The intent obligation is delivered only where an intent was stated, so a review started without one carries none of it.
The other two are unconditional.

## Capabilities

### Modified Capabilities

- `reviewer-prompt`: three rules added to the ones the delivered prompt SHALL carry, and the test-quality requirement gains the reviewer-side instruction that acts on the rule it already states.

## Impact

- `internal/testquality`: a second exported block beside `Rule`, for the reviewer alone.
- `internal/prompt`: the reviewing section gains the durable-fix rule; the test-quality section gains the reviewer's instruction; the intent section states the obligation and stays absent when no intent was stated.
- `internal/e2e` and `internal/prompt`: assertions on the delivered prompt, as the existing rules are asserted.
- **No change to** `internal/findings`, `internal/dialogue`, `internal/agent`, the findings schema, the classification, the verdict rule, the round cap, or any command.
  An intent contradiction reaches a person through `ask-user`, which already parks the round, and a durable-fix finding is a finding like any other.
- **No new input.** Nothing here reads the reviewed repository, and nothing here is configurable.

## Dependencies

None.
The most recent change to `reviewer-prompt` is `2026-08-22-no-customisations-reach-the-reviewer`, which decided what the reviewer may be told.
Nothing here re-opens it: all three rules are socrates' own text, delivered to every repository identically.

The same audit found a fourth mechanism, `no-mistakes`' repository-scoped review guidance, which cannot be borrowed without overturning that decision.
It is out of scope here and recorded as an open question rather than deferred silently.
