## Context

socrates has two things that stand in for a real reviewer.

`cmd/claudesim` speaks the vendor CLI's wire contract and answers from a scenario file, and it is what every end-to-end test runs against.
`internal/fidelity` records what the real CLI's stream looks like and fails when the simulator's shapes drift from it.

Between them they cover the transport completely and the content not at all.
The scenario files hold findings documents somebody wrote by hand, and nothing has ever compared one to a document a reviewer returned.

The four live runs during `a-rejection-is-a-claim` are the first real answers this project has ever looked at, and they disagreed with the fixtures in ways nobody predicted.
They are also the reason a task shipped unsatisfied, because the task asked a model to return a particular kind of answer and it declined four times.

## Goals / Non-Goals

**Goals:**

- Edge-case scenarios that are grounded in output a real reviewer produced, and are replayed byte for byte.
- An honest statement, in the suite's own output, of which outcomes have only ever been reached by invented input.
- A way to write a verification task about a live run that a live run can actually satisfy.
- Keeping every live answer, because they cost money to obtain and were thrown away with a temporary directory four times already.

**Non-Goals:**

- Making the dialogue happen live. What the reviewer finds is the reviewer's; this change makes the outcome recorded rather than requested.
- Prompt tuning to shift the distribution. The numbers this change surfaces are what a later change would tune against, and tuning without them is guessing.
- A corpus of other repositories' reviews. The corpus is socrates reviewing socrates.
- Replacing the hand-written scenarios. A rule must hold for shapes no reviewer has produced yet, and those scenarios are how that is tested.
- Any inference about model quality. The corpus records what was answered, not whether it was right.

## Decisions

### An entry is the answer plus what it was an answer to

A findings document alone is not evidence of anything: the same reviewer answers differently at a different effort, on a different subject, in a later round.

So an entry carries the document verbatim, and beside it: the model, the effort, the round number, the base and snapshot commits, the run's name, the day it was recorded, and one line naming the subject in words.
The document is stored exactly as the reviewer returned it, because a corpus that has been tidied is a corpus of somebody's tidying.

### Replay is by name, and the simulator does not interpret

A scenario names an entry, and the simulator emits that entry's document as its structured result.
It does not adapt it to the prompt it was given, does not renumber findings, and does not fill in a field the recorded document lacks.

A simulator that repaired a recorded answer would be back to answering from its author's beliefs, which is the thing being fixed.

### The report is a statement, not a gate

For every outcome the routing rule can produce, the report prints how many recorded answers reach it and how many scenarios assert it.

It does not fail when an outcome has no recorded answer, because whether one ever appears is not in this repository's control, and a permanently red check is one people learn to ignore.
What it does is name the gap in the output of a command somebody runs anyway, next to the number that closes it.

The distribution is also written to a file that is checked in and compared by a test.
That is the part that cannot drift: an entry removed, or a routing change that moves what the corpus reaches, appears as a diff in review rather than as a number nobody was watching.

### The routing rule is what classifies an entry, not a copy of it

The report calls `dialogue.After` on the classified document.

Restating the rule in the report would be the same defect this repository has now written down twice: an oracle derived from the sentence the implementation was derived from.
It also makes the report answer the question that matters, which is what socrates would have done with that answer, not what a second implementation thinks it would have done.

### Recording is guarded, because a findings document quotes source

An entry holds descriptions written about the code under review, and those descriptions quote it.

Recording therefore refuses unless the repository under review is socrates itself, and the refusal names the condition.
This is a smaller rule than a redaction pass and a stronger one: nothing has to be scrubbed if nothing foreign is ever recorded.

### A task about a live run says "record it", never "return this"

The lesson from task 9.4 is small and specific, and it is worth a requirement rather than a note in a retrospective.

A verification task may require that a live run happened, that its answer was recorded into the corpus, and that socrates routed whatever came back correctly.
It may not require the answer to hold a particular severity, action or shape, because the author cannot promise what the model will say and a task that asks for one is a task that is failed by the weather.

## Risks / Trade-offs

- **A corpus of four entries invites conclusions it cannot support.** Twenty-eight findings from one repository on one day at one effort is an anecdote. The report prints the count beside every number for that reason, and the entries carry their model and effort so a later reader can see how narrow the sample was.
- **The corpus grows with every live run and is checked in.** Findings documents are small, but a document is a few kilobytes and the seed is four. A retention rule is deliberately not part of this change: the first time the corpus is inconvenient is the first time anybody knows what shape it should have.
- **Replay freezes a moment.** A recorded answer is what one model version returned; the vendor's model moves. That is what the existing drift check is for on the shape side, and on the content side the entry's recorded model is what tells a reader the answer is old.
- **The report can be read as coverage.** It is the opposite: it exists to say which coverage is invented. The wording of what it prints is part of this change rather than left to whoever writes the code, because a line that says "0 recorded" reads as a failure and a line that says the parked outcome has never been observed live reads as the truth.
