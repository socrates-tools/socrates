## Context

socrates is finished as a tool and unfinished as a thing anybody uses.
The gap is that the sequence which drives it lives in a person's head and is re-explained to a session on every change.

This layer is small and its risk is unusual: everything in it is text, none of it executes, and its failure mode is a session that behaves plausibly and wrongly.
The mitigations that work here are the ones that turn text into something mechanical: one source for a statement made twice, a committed copy that a test compares against its generator, and end-to-end journeys that drive the sequence the skill describes rather than reading the skill and believing it.

## Goals / Non-Goals

**Goals:**

- The loop runs correctly without anybody explaining it.
- The response rules reach the session before it has read a finding, in the same words `scrat findings` uses.
- The skill and the committed copy cannot disagree.
- Installation is boring: idempotent, outside the repository, and it says where it wrote.

**Non-Goals:**

- A second surface for a second agent tool. The skill format is one vendor's, and generating three from one source is a problem to have when a second is asked for.
- Any new command, flag or output. If the skill needs one, the surface was incomplete and the fix belongs in the layer that owns it.
- Teaching the session how to review. It is not reviewing; it is patching and arguing.
- Asserting that a model followed the skill. That is an evaluation, not a test.

## Decisions

### The skill's text is one constant and the committed file is generated from it

The skill lives in Go as its source, the committed file is generated, and a test regenerates and compares.

The alternative, a committed Markdown file read at build time, has one failure mode and it is the one that matters: somebody edits the file, the generator is never run, and the installed skill and the repository's copy differ with nothing failing.

### The response rules have one source, used twice

`scrat findings` prints them where a run is parked, and the skill states them up front.
Both read the same constant, and a test compares the two renderings.

Two copies of the same four sentences will drift, and the direction they drift is predictable: the one nobody reads while implementing stays as it was, and the one in front of the session gets softened.

### The skill teaches polling, not following

`scrat status` is the command the skill names in the wait loop, and the skill says why: everything the session reads lands in its context, and it is about to need that context for patching.

`scrat watch` is named in the skill exactly once, as the thing to tell the person to run in another tab if they want to see it happen.

### The loop the skill teaches

1. Implement, as asked.
2. `scrat review --intent <what the person asked for>`.
3. Poll `scrat status` until it is not working.
4. On `clean`: say so, say what was reviewed, and stop.
5. On `failed`: say the review did not run and why. This is not a pass, and nothing has been checked.
6. On `escalated`: read `scrat findings`, present them to the person, and stop. These are theirs.
7. On `awaiting-response`: read `scrat findings`, and for each finding either patch it or `scrat respond --reject` with a concrete reason. Then `scrat continue`, and go to 3.

The branch that matters is 5, and it is placed before 6 deliberately: a session scanning a list for "what do I do about findings" reads past a state that has none.

### The skill says what socrates will not do for you

Three sentences that prevent three specific attempts:

- socrates never patches anything. The session patches, in its own working tree.
- There is no way to explain the change to the reviewer. The reviewer gets the diff, the repository and the stated intent, and that is the design rather than a limitation to work around.
- A rejection is a claim that goes to the next reviewer. It is not a way to dismiss a finding, and a finding rejected without a reason is refused by the tool.

### Installing

`scrat skill` writes the skill at user level and prints the path.
`--to` writes it somewhere else, which is what the repository's own generator target uses.

Running it twice is a no-op with the same message, because an installation command that a person has to remember whether they ran is a command they run again anyway.

## Risks / Trade-offs

- **A skill is advisory.** The session may ignore it, and nothing in socrates can stop that. What the design buys is that the correct sequence is stated once, in front of the session, at the moment it needs it, and that every branch of it has been driven end to end so the sequence itself is known to work.
- **Testing text is testing text.** The skill's assertions are about generation, drift and the shared source, not about whether the wording persuades a model. The wording is evaluated by using it, and what a live run found belongs in the retrospective rather than in an assertion.
- **One vendor's skill format.** A second agent tool means a second renderer over the same source. Accepted now, and cheap later precisely because the source is one constant rather than a file in that vendor's shape.
