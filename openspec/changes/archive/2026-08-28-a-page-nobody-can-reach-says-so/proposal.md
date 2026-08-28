## Why

The surface binds a loopback address and nothing widens it, which is what keeps a document quoting reviewed code on one machine.
A person whose browser is on that machine opens the address and reads the page.

A person whose browser is somewhere else does not.
A devcontainer, a remote host, anything with a network namespace of its own: the address is forwarded to the machine's own interface, the surface is not listening there, and the connection opens and carries nothing.
The browser reports an empty reply, which names no cause and looks exactly like a port that was never forwarded.

Nothing said anything.
The address printed was correct, the page was served, and the one fact that would have explained it, that the page can be opened from the machine serving it and from nowhere else, was never stated.
So the first person to meet this spends an hour on a port mapping that was working the whole time.

The fixed port exists so an operator can put something in front of the surface once.
Nothing says so at the moment somebody needs to know it.

## What Changes

- Serving a page says that the page can be opened from the machine serving it and from nowhere else, and names the field that changes the address printed.
  It says it on the stream a person reads, never on the one a caller parses.
- It says it only where nothing states the address the surface is reached at, because an operator who stated one has already answered the question.
- The README carries the step that makes it work: forward the fixed port once, state the address the surface is reached at, and name that address for socrates rather than for the binary serving it.

Nothing about the bind changes.
The surface still binds loopback, still refuses every address that is not, still admits only the host it bound, and still ships no forwarder and spawns none.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `review-surface`: gains `A page says where it can be opened from`, next to the requirement that binds loopback and refuses to widen.

## Impact

- `cmd/agora`: says it when a page is served and nothing states the address the surface is reached at.
- `internal/surface`: the sentence, in its `messages.go`.
- `README.md`: the install-time step, and the port in it read from the source rather than written twice.
- `internal/skill`: untouched. The sentence is already in front of whoever runs the command.
- `cmd/socrates` and every package it reaches: untouched.
- No new dependency.
