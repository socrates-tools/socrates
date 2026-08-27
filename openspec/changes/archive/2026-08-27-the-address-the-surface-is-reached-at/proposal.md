## Why

`agora serve` prints the address it bound, and the skill tells the session to hand the person that address.
Where nothing sits in front of the surface those are the same address, and the loop works.

They are not the same address everywhere the loop runs.
A surface bound inside a container is bound on that container's loopback, and the person reading the page is on the machine outside it.
The address agora printed resolves for the process that printed it and for nobody the page was served to, so the session hands over a URL that opens nothing, and the one round a person was meant to read is the round they cannot.

Reaching the surface through something is already provided for: the bind stays loopback, and an operator who wants the page reachable puts a proxy in front of a port that already works.
What is not provided for is agora knowing that happened.
A proxy is a fixed dial, so a surface behind one cannot take a port the operating system chose, and the address it prints is not the address to hand over.
Both facts are the operator's, both are stable for the life of a machine, and today neither has anywhere to live except a flag typed identically every time by a session that has not been told to type it.

The alternative is teaching the skill the flag, and it is the wrong place twice over.
The port and the hostname belong to one machine's arrangement rather than to every reader of the skill, and the skill is generated, so what it carries is carried into every repository a session reads it in.

## What Changes

- `agora` reads an optional configuration file holding two addresses: the loopback address to bind, and the address the surface is reached at.
- The bind address supplies the default for `--address`, which continues to override it.
  It is refused unless it is loopback, by the check the flag already passes through, so configuration reaches the bind by the same gate and widens nothing.
- `serve` prints the address the surface is reached at: the configured one where one is configured, and the address it bound otherwise.
- **The host guard is unchanged, deliberately.** A request naming the configured address is rejected exactly as it is today.
  The surface accepts what it bound, so a proxy in front of it rewrites the host on the way through, and the configured address stays something agora says and never something it admits.
- `agora` gains a configuration file where it had none.
  socrates' file is not it: the two are separate programs, and socrates' file refuses a field it does not know.
- **The skill is untouched, and that is the result rather than an omission.**
  It already says to hand the person the address `agora serve` prints, and that sentence becomes true in both arrangements without being rewritten.
- Nothing in the repository names a proxy, a container or a port.
  What the surface is reached at is a string an operator writes on their own machine.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `review-surface`: `One binary, five commands, beside socrates` says what address `serve` prints, which is now the address the surface is reached at rather than the address it bound.
  `The surface binds a loopback address and nothing else` gains configuration as a source of the bind, held to the same refusal, and states that the address the surface is reached at reaches neither the bind nor the host guard.

## Impact

- `internal/surface`: a configuration file beside the session file it already writes, holding the two addresses and validating them.
  It is the surface's own address, and `DefaultAddress` and `CheckLoopback` are already here.
- `cmd/agora/serve.go`: the `--address` default comes from configuration, and both exits print what the surface is reached at rather than `open.Address`.
- `internal/e2e`: a surface configured to be reached at an address it did not bind, driven through the real binary, and a request naming that address still turned away.
- No new dependency: the configuration is YAML, which the repository already parses.
