## Context

socrates already parses the whole stream before it looks at how the process exited.
Everything this change needs is therefore already in memory when the reason is composed, and the change is about which of it gets used.

## Goals / Non-Goals

**Goals:**

- A person reading a failed run can tell a transient provider failure from something they have to fix, without opening the stream.
- The reason names what the agent said rather than which field socrates happened to reach first.
- A failure with no reason behind it is distinguishable from one whose reason was lost.

**Non-Goals:**

- Retrying a review because the provider was overloaded. The vendor's CLI already retried ten times over four and a half minutes before it gave up, and a second retry loop on top of that one would multiply a wait socrates already bounds with a deadline. A person who reads `overloaded` runs `socrates review` again.
- Deciding anything from that text. Rejected below.
- Reading the status code the same event carries. Rejected below.
- Neutralising the text through `internal/hygiene`. That package exists for a field being re-embedded in a prompt socrates composes, and no prompt socrates composes carries a failure reason: a run that failed does not continue, so there is no next round to carry it into.

## Decisions

### The account the agent gave is preferred to what it wrote to its error stream

Both are the agent talking about its own failure, and they are not equally informative.
The error stream is where a process complains about its arguments, its environment and its own startup, and it is empty on every failure the CLI handles itself.
The result event is where the CLI says why the work did not happen, which is the question a review's failure reason is answering.

Where both are there both are carried, in that order.
There is no case where dropping one of two accounts of the same failure buys anything, and the ceiling below bounds the pair.

### The subtype is a fallback, not the reason

The live stream is the argument: `subtype` was `success` on the event whose error flag was true and whose message named the provider.
A field that reads `success` on a failure is not a description of the failure, and quoting it produced `claude reported the result "success", so no review completed`, which reads as a defect in socrates rather than an account of what happened.

It stays as the fallback, because an event that reports an error and carries no message still has to be reported as something, and the subtype is then the only thing said.

### Quoting is not deciding

The requirement this delta modifies already refuses to derive an answer from the vendor's prose, because a rule built on prose stops working when the prose is edited, and it stops working by failing every review.
That refusal is about decisions, and it stands.
Nothing branches on this text: it is not parsed, not matched, not compared, and no state, verdict or exit code depends on a word in it.
It is carried to a person verbatim, and the failure was already a failure before it was read.

### The status code beside it is not read

The same event carries a numeric status for the API error, and reading it would let socrates say `529` from a field rather than from prose.
It would also be a second field to keep, whose only content is already the first four words of the message, and socrates makes no decision either could inform.
One field read verbatim and quoted is the smaller thing to be wrong about.

### The reason is bounded by the ceiling the error stream already has

The reason reaches two places that make length matter: a run record that is written once and read afterwards, and `socrates status`, whose every line lands in the context of a session that is about to need that context for patching.
The result event's message is the agent's text, which on some failures is the whole assistant turn.

It is bounded by the same ceiling and the same mechanism the error stream is bounded by, rather than a second one, so there is one truncation rule and one marker for a person to recognise.

### The simulator reproduces the recorded shape rather than an invented one

The behaviour emits the error flag, keeps the subtype the live stream carried, puts the message where the live stream put it, writes nothing to the error stream, and exits with what the scenario asks for.
Exiting is what the scenario controls, because the two branches this change fixes differ only in that, and a scenario that asks for a zero exit drives the branch a non-zero exit would mask.

## Risks / Trade-offs

- **A vendor release renames the field carrying the account, and the reason silently becomes the subtype again.** The fidelity guard pins that field and the error flag beside it against the checked-in recordings, so its disappearance is a diff somebody reads. That is what the opening event's tool list already has and what this borrows.
- **The message is the agent's own text, and an agent under review composed the change being reviewed.** It reaches a person and a run record, never a prompt socrates composes, and nothing branches on it. Where it reaches a polling session's context it is bounded, which is the reason for the ceiling above.
- **Two accounts of one failure make a longer reason than one.** Deliberate. The pair is bounded, and a person diagnosing a failed review is the reader who wants both.

## Migration Plan

No stored run changes shape, and no promoted scenario moves.
A run that already failed keeps the reason it was recorded with.

## Open Questions

None.
