## Why

A live review of another repository failed, and this is what socrates said about it:

```
state: failed
verdict: claude exited: exit status 1:
```

The trailing colon is the defect.

The reviewer died because the model provider was overloaded.
The vendor's CLI retried ten times over four and a half minutes, gave up, wrote `API Error: 529 Overloaded` into the terminal result event of the stream socrates was parsing, and exited non-zero.
socrates read the exit status, took the agent's error stream as the reason, found it empty because the CLI had said everything on the other channel, and printed an empty quotation.
The account it needed was in memory, parsed, one field away, and dropped unread.

A failure reason answers one question: run it again, or fix something.
A provider that was overloaded is the first, a prompt the vendor refused is the second, and `exit status 1` is both.
The person who hit this read `r1/stream.jsonl` out of the run directory to find out which, which is the work the reason exists to save.

The same event is misread one branch along.
Where the agent reports an error and exits zero, socrates names that event's subtype.
On the live stream the subtype was `success`, beside an error flag of true and a message naming the provider.
Of the three fields on that event, socrates named the only one that said nothing.

This is the same shape as the change beside it: not a wrong answer, but the absence of one printed as though it were an answer.
It is a separate change because it is a change to what a failure reports rather than to a check, and because the live run that would have closed that change is the run this defect made illegible.

## What Changes

- **A failed call carries the reason the agent gave.** The result event's own account of why it stopped is preferred to what the agent wrote to its error stream, both are carried where both are there, and the subtype beside the account is a fallback rather than the answer.
- **A failure with nothing behind it says so.** Where the agent said nothing on either channel, the reason says that, rather than trailing a colon onto an empty string.
- **The reason is bounded.** It reaches a run's record and a polling session's context, on the same ceiling the error stream already has, and a vendor is free to make it arbitrarily long.
- **The simulator can report an error in its result event**, so both branches are driven through the real binary against a stream shaped like the one the live run recorded.

## Capabilities

### Modified Capabilities

- `agent-seam`: what a failed call carries, and which of the result event's fields is the reason.

## Impact

- `internal/agent`: three more messages, one reading of the account the result event carries, and the ordering that prefers it to the error stream.
- `internal/claudesim`: a behaviour that reports an error in the result event and exits with whatever the scenario asks for.
- `internal/fidelity`: the guard pins the two fields on the result event that this reading depends on.
- No stored run changes shape, no verdict changes, and no routing changes.
