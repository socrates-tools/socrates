## Why

By the end of the dialogue layer socrates works, and working it is a five-command sequence with three branch points and one rule that is easy to get backwards.

The consumer of that sequence is a coding session, and a coding session is told what to do once, in prose, by whoever is at the keyboard.
Told it once per change, by hand, it will get it right the first few times and then drift: it will run `watch` instead of `status` and fill its own context with the reviewer's output, it will read the findings and start arguing instead of patching, it will treat a `failed` run as a pass because nothing blocking came back, or it will simply stop after round one.

Every one of those is the tool failing, and none of them is a bug in the tool.

A generated, installed skill is how the sequence stops depending on somebody remembering it.
It is also the one place the response rules can be stated to the session that is about to defend its own work, before it has read a single finding.

## What Changes

- **A `/scrat` skill is generated from one source, committed, and installed at user level.** A test fails while the committed copy and its source disagree, so the two cannot drift apart quietly.
- **`scrat skill` installs it**, defaulting to user level and naming where it wrote. Idempotent, so running it twice is not a question anybody has to think about.
- **The skill teaches the loop, including the branches.** Start a review, poll status, read findings, patch or respond, continue, and stop on each terminal state for the right reason.
- **The skill teaches `status` and not `watch`.** Everything a session reads lands in its context, and the reviewer's output stream is the largest thing socrates produces.
- **The skill carries the response rules verbatim from the same source `scrat findings` prints them from.** One statement, two surfaces, and a test that they are the same text.
- **The skill says that a failed run is not a pass**, because a session that got no blocking findings back and a session whose review never ran look identical from the outside if nobody says otherwise.
- **The skill offers no way to explain the change to the reviewer**, because there is none, and a skill that implies otherwise invites a session to try.

## Capabilities

### New Capabilities

- `agent-skill`: the generated skill, what it teaches, where it installs, and the drift check that keeps the committed copy honest.

### Modified Capabilities

- `cli-foundation`: the command surface gains `skill`.

## Impact

- One new package holding the skill's text as the single source, and a generator target that regenerates the committed copy.
- Writes at user level when asked. It creates nothing in the repository under review, so the boundary the core drew is unmoved.
- No change to any command's behaviour, exit code or output. The skill drives the surface that already exists, and a skill needing a new flag is a signal the flag was missing, not that the skill is special.
