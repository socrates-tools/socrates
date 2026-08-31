## Context

See `proposal.md` for why.

Three things about the current shape decide most of this design.

`surface.Serve(address, idle)` binds exactly the address it is handed and returns an error if it cannot.
Two callers reach it: `cmd/agora/serve.go`'s `hold`, which is the detached child, and tests.

`cmd/agora/serve.go`'s `start` runs in the caller's process.
It reads the record, reuses a standing surface if one answers, and otherwise calls `surface.CheckFree(address)` and spawns the child with `--address address`, then waits for the record the child writes.
The child is detached with no standard output and no standard error, so nothing the child refuses reaches the person who ran the command: a child that cannot bind exits silently and the parent times out after twenty seconds with "the surface never answered".

`surface.Record` carries `Address` taken from `listener.Addr()`, which is what the surface actually bound.
`poll`, `end` and `stop` read the record, so a moved port needs nothing from them.

## Goals / Non-Goals

**Goals:**

- One place states the window, and everything else reads it from there: the top of the window is `DefaultAddress`.
- A person who states nothing and whose `18018` is held gets a page, and the address they are handed is the one bound.
- A person who states something gets today's refusal, with today's message, on today's port.
- The refusal for a whole window held reaches the person who ran the command, not a silent child exit and a twenty-second timeout.

**Non-Goals:**

- Choosing a port for a *stated* address. Stating one is the way to pin one, and a stated port that is held stays a refusal.
- Making the window configurable. A configured window is a third way to say what a stated address already says, and its only new case is an operator who wants a moving port aimed at nothing.
- Any change to the record, the session model, the host guard, the page, or socrates.

## Decisions

### The parent chooses the port, not the child

The walk happens in `start`, before the spawn, and the chosen address is passed to the child on `--address` exactly as today.

The alternative is for the child to walk at bind time, which is race-free: the child's own `net.Listen` is the claim, so nothing can take the port between the check and the bind.
It is rejected because the child cannot speak.
It is detached with no standard error by construction, so a whole window held would surface to a person as the parent's twenty-second "the surface never answered" timeout rather than as a refusal naming the window.
Making the child speak means a channel back from a detached process, which is a daemon's problem this tool does not otherwise have.

The race the parent's walk keeps is the one `CheckFree` already has today: something can take the chosen port between the check and the child's bind.
Walking does not widen that gap, because the last probe of the walk is the check the gap starts at, and the failure at the end of it is the same as today's.

### The window is ten ports, `18018` down to `18009`

Ten is enough for the case that motivates this, which is several checkouts or worktrees on one machine and a stray process holding the top port, and small enough that an operator can forward the whole range in one line where they want to.

Descending rather than ascending because the top is the documented port.
Counting up would make the documented port the floor of the range rather than its landmark, and `18019` upward is territory nothing has reserved.

The window lives in `internal/surface/bind.go` beside `DefaultAddress`, as the top port and a count, with `DefaultAddress` composed from the top rather than written twice.

### Stating the address the surface is reached at pins the port

This is the one coupling the design adds, and it is deliberate.

`review-surface` already says the address the surface is reached at reaches neither the bind nor the host guard, and that stays true: what it reaches here is only whether the *default* port is allowed to move, and it reaches nothing about what may be bound or what host may be named.

Without it the README's own recipe breaks silently.
That recipe is: forward `18018`, then state `reached_at` so the address printed is one that opens.
An operator who followed it has aimed a forward at `18018` and told the surface to print an address that resolves through it.
A surface that walked past a held `18018` would print exactly the address the operator stated, which by then opens onto nothing, and that is the failure `a-page-nobody-can-reach-says-so` exists to prevent.
A refusal naming the held port is the right answer for that operator, because they are the one person on the machine who knows what is holding it.

The alternative considered was to key the walk on the stated bind alone and change the README to tell operators to state `bind` as well.
It is rejected because it breaks configurations already written against the current README, and it breaks them by handing a person a dead address rather than by refusing.

### `CheckFree` keeps its shape; the walk is a new function beside it

`CheckFree(address) error` stays what it is and keeps its message, because that message is exactly right for a stated address.

Two new functions sit beside it.
`surface.WindowAddresses()` composes the window from its top port and its count, and `surface.FirstFreeIn(addresses)` walks any list of addresses with the same listen-and-close probe, returning either one to bind or a refusal naming the ends of the list.
Splitting the window from the probe keeps `internal/surface`'s own tests off the real window ports, which the end-to-end tests hold, and those two suites run as separate packages and so run at the same time.

The held-port message splits in two: today's `portIsHeldByAnother`, unchanged, and a new one for a window with nothing free in it.

`serve.go` picks between them from whether a port was committed to, which is read from the flag and the configuration rather than from the address:

```
command.Flags().Changed(addressFlag) || configured.Bind != "" || configured.ReachedAt != ""
```

Reading it from the flag rather than from the address keeps `--address ""` a refusal.
The flag's default stays `DefaultAddress`, so an operator who stated only the address the surface is reached at commits to that port without having to name it, and `CheckLoopback` runs on it exactly as on any stated address.
Where nothing was committed to, the address is dropped and the window is walked instead.
A walked address needs no `CheckLoopback` of its own: the window is composed from the loopback host, and `surface.Serve` checks whatever it is handed before it binds.

### `livesmoke` moves off the fixed port

`internal/fidelity/livesmoke_test.go` serves on `surface.DefaultAddress` and asserts the surface bound exactly that, which is a test of the fixed default sitting in the live smoke suite.
It moves to the resolution, so a developer whose `18018` is held can still run the live smoke.
What the live smoke is for is the real CLI, and the default bind is covered by its own scenario in the end-to-end suite.

## Risks / Trade-offs

**Two cold `agora serve` calls racing now stand up two surfaces rather than one refusing.**
Today the second loses `CheckFree` and refuses; with the walk it lands on the next port down and overwrites the record, orphaning the first until the idle period collects it.
The window for it is between the first parent's spawn and the child's record write, which the reuse path closes as soon as the record exists.
Mitigation: none beyond the idle collection that already exists. An orphaned surface holds a port and a page nobody was handed, and stops itself.
Widening this into a lock is a daemon's mechanism, which is out of scope by the tool's own architecture.

**A person who forwarded `18018` and stated nothing now gets a page on `18017` behind a forward that does not carry it.**
Mitigation: the page already says where it can be opened from, and it names the address bound, so what they read names `18017` rather than `18018`.
An operator who stated `reached_at` is pinned and refuses instead.

**The window is a range an operator has to know about to forward.**
Mitigation: the README says so, and stating `bind` or `reached_at` reduces it to one port, which is what an operator putting something in front of the surface does anyway.

## Migration Plan

None. Nothing persisted changes shape, and no configuration written against the current behaviour changes meaning: every configuration that states `bind` or `reached_at` binds exactly the port it binds today, and every one that states neither only stops refusing in a case it used to refuse in.
