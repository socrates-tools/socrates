## Why

The port the surface binds when nothing states one is fixed at `127.0.0.1:18018`, and a port already held by something agora does not know of is a refusal.
That refusal is correct where an operator named the port and wrong where nobody did: a person who has never thought about ports meets a failure about one, and the only way out is to learn a flag or write a configuration file to get back to where they started.

The fixed port is not the thing worth keeping.
What is worth keeping is that the port does not move under something an operator aimed at it, and that a person is handed an address that opens.
A short, ordered, deterministic window satisfies both: where nobody has committed to a port, the surface takes the first free one counting down from the top of the window, and where somebody has, a held port is still a refusal.

## What Changes

- Where nothing states the address the surface binds and nothing states the address it is reached at, `agora serve` SHALL take the first free port counting down from `18018`, over a window of ten ports ending at `18009`, instead of refusing when `18018` is held.
- The window is ordered and closed, so the port chosen on a given machine is determined by what is held rather than by what the operating system hands out, and the top of the window is still what a free machine binds.
- A whole window held is a refusal that names the window, replacing today's refusal that names one port.
- Where an operator stated the address to bind, by flag or by configuration, nothing walks: a held port is refused exactly as today.
- Where an operator stated the address the surface is reached at, nothing walks either: stating it is how an operator says something is in front of the surface aimed at a port, and a port that moves under that hands a person an address that opens onto nothing.
- The record `poll`, `end` and `stop` read already carries the address bound rather than the address asked for, so those commands follow a moved port with no change.
- `README.md` says that the port moves down where nothing states one, and that the recipe it already gives, which is to forward `18018` and then state the address the surface is reached at, is what pins it.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `review-surface`: the requirement `The surface binds a loopback address and nothing else` changes from binding one fixed port where nothing states an address to taking the first free port over a fixed descending window, and states when the window is not walked at all.

## Impact

- `internal/surface/bind.go`: `DefaultAddress` is joined by the window it is the top of, and by the resolution that walks it.
- `internal/surface/messages.go`: the refusal for a held port splits into the stated-port refusal and the whole-window refusal.
- `cmd/agora/serve.go`: resolving the address to spawn on distinguishes stated from unstated, and walks only in the unstated case.
- `internal/e2e/address_test.go`: the scenarios for the default bind and for a held port.
- `internal/fidelity/livesmoke_test.go`: it serves on `DefaultAddress` and asserts the surface bound exactly that; it moves to the resolution so a held `18018` does not fail the live smoke.
- `README.md`: prose beside the forwarding recipe. The step it names stays `18018`, so the check that the documented port is read from the surface rather than written a second time is untouched.
- No change reaches socrates, the document, the renderer, the page, or the session model.
