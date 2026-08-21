## Why

The core layer ships `scrat status`, and it ships it because a waiting session needs three lines and nothing more.
That leaves the person with nothing.

A review takes minutes.
While it runs, the person wants to watch it happen in a second tab, and when it has finished they want to read what it actually did, and neither of those is `status`.
Pointing them at `status` in a loop is the answer that looks like it works: it prints a fixed few lines, so watching a review through it means running it repeatedly and learning nothing about what the reviewer is doing between polls.
Pointing them at `cat` on the round's output file is the other answer that looks like it works, until the file is JSONL and the reviewer is chatty.

Three needs look like one command and are not.
The consumer differs, the follow behaviour differs, and the cost of being verbose differs by an order of magnitude, because one of the three prints into a language model's context and two print into a person's terminal.

This change is where the JSONL log written by the core stops being an implementation detail and becomes something a person reads.

## What Changes

- **`scrat watch` follows a run live and exits with its verdict.** It renders the structured log readably as events arrive, stops at the terminal state, prints the verdict, and exits with the verdict's code. That last part is what later lets it be the body of a hook with no rewriting.
- **`scrat logs` dumps or pages a finished run.** No follow, no waiting, no exit-code theatre beyond the verdict.
- **A follower that races the spawn waits rather than erroring.** The log may not exist for the first moments of a run, and a command that fails in that window fails exactly when a person is most likely to run it.
- **A follower never sits forever after a run ends.** Stopping at the terminal state is a requirement, not a convenience: a naive follower on a finished run is indistinguishable from a stalled review.
- **Rendering is separate from the log.** The log stays structured, and the renderer is what turns it into something to read. Nothing writes pre-rendered text.

## Capabilities

### New Capabilities

- `run-presentation`: following a run live, rendering the structured log for a person, dumping a finished run, and the follow semantics that keep either from lying about what is happening.

### Modified Capabilities

- `cli-foundation`: the command surface gains `watch` and `logs`, and the boundary between what each of the three reading commands is for becomes a requirement rather than a convention.

## Impact

- One new package for rendering, which reads the event log and writes to a terminal, and is the only place that knows how an event looks to a person.
- No change to what is written. The log format, the run states, the verdicts and the exit codes are the core's and are untouched.
- No terminal library and no multiplexer integration. Detecting whether the output is a terminal is the whole of the terminal awareness, so the suite runs with no terminal at all.
