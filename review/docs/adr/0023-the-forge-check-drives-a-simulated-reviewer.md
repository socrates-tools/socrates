# 23. The forge check drives a simulated reviewer

- Status: accepted
- Date: 2026-09-09
- Part of: reviewing a branch, and reviewing a change that is not mine

## Decision

ADR 19 moved the question of whether a forwarded finding can be a comment on a line into socrates-review, where a test in this repository answers it.
What no test here can answer is whether a forge accepts the payload that answer produces and puts the comment where the payload said.
`make forge-smoke` is that check, behind the `forgesmoke` build tag, run by hand, against `socrates-tools/repository-under-review` and its one long-lived pull request.

Its reviewer is a simulator rather than a real CLI.
What the check reads is the forge, so a real reviewer would add a model's judgement to a test about an HTTP response, make the run cost money, and leave which findings exist up to what a model happened to say.
This repository already refuses to require a particular finding of a live run, and a check that needs one finding inside the change and two outside it would be requiring three.
So the findings are stated by the scenario, and the fixture pull request is shaped to match: one changed line, a caller it does not touch, and a defect in each.

It asserts three things: the forge puts the finding socrates-review reports as inside the change on that line and holds it against a position in the diff; the forge refuses a line comment for a finding socrates-review reports as outside it, which is what makes the placement rule load bearing; and everything outside the change reaches the pull request as one comment rather than one each.
It deletes both comments afterwards, whether it passed or failed.

## Rejected

- a real reviewer, which buys nothing the live smoke check does not already buy and makes the forge assertions depend on what a model found; and opening a pull request per run, which leaves debris behind every time the check fails.
