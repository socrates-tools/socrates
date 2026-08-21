## Why

`agent-seam` already says it: a simulator nobody compares to the real thing becomes a specification of its author's beliefs, and a green suite against one manufactures confidence.

The comparison it requires covers the shape of the stream.
It says nothing about what the reviewer answers, and that is where the beliefs are.

Four live runs during `a-rejection-is-a-claim` produced 28 real findings, and put numbers on the gap:

- **`error` severity was never used.** Not once in 28 findings. Every simulator scenario in the suite leans on it, and the scenario that drives the parked journey opens with an `error`.
- **`{warning, no-op}` never appeared.** That shape is what the routing rule was rewritten for, twice, in this repository, in one afternoon. No real reviewer has ever emitted it.
- **`{info, ask-user}` appeared twice**, and it means a reviewer can escalate a run over a finding it labelled as holding nothing. No scenario in the suite covers it and nobody has decided whether that is right.
- **Not one round would have parked.** Every live round carried an `ask-user` finding, so the whole dialogue layer has never been exercised by a real reviewer.

The last one is why `a-rejection-is-a-claim` shipped with an open task.
Task 9.4 asked for a live dialogue, and a live dialogue is not something socrates can ask for: it depends on what a model chooses to return.
The task was unsatisfiable as written, and the four attempts to satisfy it cost real money to rediscover the same escalation each time.

So the answer is not "test it live" and not "invent a scenario and call it covered".
It is to keep what real reviewers actually answered, replay it, and be honest in the suite's own output about which paths have only ever been walked by fixtures somebody wrote.

## What Changes

- **Real answers are recorded, verbatim, as a corpus.** Every live round's findings document is kept with what it was answered about: the model, the effort, the round, the base and snapshot commits, and the run it came from. The four already performed seed it.
- **The simulator can replay a corpus entry.** A scenario can name a recorded answer instead of carrying an invented one, so a test drives socrates with output a real reviewer actually produced.
- **The suite says which routing outcomes have been observed.** A report classifies every corpus entry through `dialogue.After` and prints, per outcome, how many recorded answers reach it and how many scenarios only assert it. An outcome no recorded answer has ever reached is named as such rather than left to a retrospective.
- **The recorded distribution is checked in and cannot drift silently.** The summary is a file the corpus regenerates and a test compares, so an entry that disappears, or a rule change that moves what the corpus routes to, shows up in a diff rather than in nobody's memory.
- **A verification task may not require a particular answer from a live run.** It may require that a live run happened, that its answer was recorded, and that whatever it returned was routed correctly. What the reviewer chooses to find is not the author's to promise.
- **The corpus is recorded from socrates' own reviews only.** A findings document quotes the code it reviewed. Recording from an arbitrary repository would check somebody else's source into this one.

## Capabilities

### New Capabilities

- `reviewer-corpus`: the recorded answers of real reviewers, what each entry carries, how one is recorded, how a test replays one, and the report that says which outcomes the corpus has actually reached.

### Modified Capabilities

- `agent-seam`: the comparison against the real tool covers what the reviewer answered and not only the shape it answered in, and the simulator gains replay of a recorded answer as a first-class scenario source.

## Impact

- One new package holding the corpus, its entries and the report. Recording reaches the network and is behind a build tag, as the existing recordings are. Replay does not, so the ordinary suite uses the corpus and still reaches nothing.
- The existing hand-written scenarios stay. They cover shapes no reviewer has produced yet, which is legitimate for a rule that must hold anyway, and the report is what stops that being mistaken for evidence.
- No change to socrates' behaviour. Nothing in `cmd/scrat` gains a flag, no verdict moves, and the routing rule is untouched by this change.
- The open task in the archived `a-rejection-is-a-claim` is answerable afterwards: the dialogue is driven from a recorded answer, and the live obligation becomes recording whatever the reviewer said rather than needing it to say something in particular.
