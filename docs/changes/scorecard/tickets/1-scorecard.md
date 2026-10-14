# 1: Scorecard runs weekly, publishes, and is the one run nobody asks for

## What to build

What the spec describes: the workflow, the convention check on what may start unasked, the amended invariant and the badge.

## Acceptance criteria

- [x] `.github/workflows/scorecard.yml` starts on a weekly schedule and on request, and on nothing else.
- [x] It publishes its results and uploads them to code scanning, with every action pinned by commit.
- [x] The suite fails where another workflow starts unasked, or where Scorecard starts on a push.
- [x] `AGENTS.md` names the exception and ADR 64.
- [x] The README carries the badge.
- [x] A run started after it lands is recorded, and the badge reads a score.

## User stories

- 1, 2, 3, 4, 5

## Blocked by

None
