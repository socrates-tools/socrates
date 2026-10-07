# 55. A simulator is held to recorded shapes, and only invention is drift

- Status: accepted
- Date: 2026-10-08
- Part of: docs-cleanup

## Decision

Every end-to-end scenario runs against a simulator of each vendor's CLI, and a green suite against a simulator says nothing about the real CLI unless something ties the two together.
The drift check is that tie, and it found a real gap on its first run: the simulator emitted `structured_output` when no schema was supplied, which the real CLI does not, and no test could have seen it because every test supplied a schema.

What is recorded from the real CLI is its shape: the field paths each event type carries, never the stream itself.
A field is recorded before its value is walked, so an empty array still records the field, and a guard pins the field rather than a path beneath it that only a non-empty value would produce.

The check runs in one direction.
Every field the simulator writes must exist in the recording, and a field the recording holds that the simulator does not write is not drift.

## Rejected

- Recording raw streams.
  A raw stream carries the recording machine's hook output, memory paths, socket paths and session identifiers, and a reviewer of the diff cannot read it.
- Checking both directions.
  The real CLI emits event types and fields the simulator has no reason to, so the check would fail permanently and be turned off within a week.
- Trusting the simulator because the suite is green.
  That is the confidence the first run of the drift check showed to be manufactured.
