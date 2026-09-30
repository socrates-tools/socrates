# The page in a devcontainer: unreachable, and then unread

Observed on 2026-09-30, during a design round served with `socrates agora serve` from a coding session inside a devcontainer, with the browser on the host.
Two separate faults, one after the other, cost the person the round they had answered.

## What happened

1. `socrates agora serve --document round1.json` printed `http://127.0.0.1:18018/sessions/<session>`.
2. The person opened it in a browser on the host and got `ERR_EMPTY_RESPONSE`.
3. A hand-started relay made the page reachable, and the person answered the round.
4. `socrates agora poll --session <session>` had been started as a foreground command and was moved to the background while it waited.
5. The poll collected the answer and exited 0, so the answer went to the harness's task output file rather than back to the session.
6. The session's read of that file was denied by the harness's permission classifier, so the answer never reached the session.
7. Delivery happens once, so polling again waits for a new submit rather than handing back the collected one.

The person had hit step 2 before, in earlier sessions.

## Fault 1: a published port cannot reach a loopback listener

### Cause

The devcontainer publishes a fixed range of ports to the host with a compose `ports:` block, and `18018` is inside it.
Docker delivers a published port to the container's own network interface, and not to its loopback.

The surface binds `127.0.0.1:18018`, because it binds a loopback address and nothing widens it.
Before the relay, the only listener on the port was the surface on loopback:

```
LISTEN 0 4096 127.0.0.1:18018 0.0.0.0:* users:(("socrates",...))
```

So a connection from the host's browser reaches Docker's proxy, which accepts it and then finds nothing listening on the container's interface.
A browser reports that as an empty reply.
Inside the container the page answers, which is why nothing looked wrong from the session:

```
$ curl -s -o /dev/null -w '%{http_code}\n' http://127.0.0.1:18018/sessions/<session>
200
```

The two environments' rules collide.
A devcontainer that publishes ports through compose expects a directly reached service to bind `0.0.0.0` so Docker can forward it.
agora's invariant is that the surface binds loopback only.
`README.md` already names this exact case: publishing a port in a compose file does not work, and a forwarder must run inside the namespace the surface is in.

### Why the port being published is not enough

Publishing the port is what produced the empty reply rather than a refused connection.
Published means Docker forwards it to the container's interface.
It does not mean anything inside the container is listening there.

### The relay that made it work

```sh
socat TCP-LISTEN:18018,bind=<the container's interface address>,fork,reuseaddr TCP:127.0.0.1:18018 &
```

It binds the interface address rather than `0.0.0.0`, because Linux refuses `0.0.0.0:18018` while `127.0.0.1:18018` is held.
The surface's host guard accepts the relayed requests, because the browser still names `127.0.0.1:18018`.
After it, the port has two listeners: the surface on loopback, and the relay on the container's interface.

The relay dies with the shell that started it and with the container, so the fault returns on the next start.

### Where a fix could live

- The environment: a relay started on every container start, from the devcontainer's `postStartCommand` or a per-person hook it runs.
  It needs the port the surface binds to be fixed, so `~/.socrates/agora/config.yaml` states `surface.bind: 127.0.0.1:18018`, which stops the window walking down to a port nothing relays.
- The environment, through a reverse proxy the devcontainer already publishes, with `surface.reached_at` naming the address the proxy serves the page at.
  The proxy must send the upstream host as `localhost:18018` or `127.0.0.1:18018`, because the host guard refuses a host naming no port or another port.
- socrates: the notice `socrates agora serve` prints today is the same everywhere.
  Where the surface can tell it is inside a container, it could say that a published port will not reach it and print the relay for the interface address it can see.
- socrates: a relay of its own on a non-loopback interface would contradict the invariant that nothing widens the surface, and would need an ADR that overturns it.

## Fault 2: a poll that outlasts the harness loses the answer to the session

### Cause

`socrates agora poll` waits up to `surface.PollDeadline`, six hours, and offers no flag to wait less.
The harness bounds a foreground command at ten minutes, so a person who takes longer than that, or a session moved to the background while waiting, turns the poll into a background task.
The poll then still collects the answer, and the answer is written to a file only the harness knows about.
Reading that file is a separate action the harness may refuse, and here it did.

Delivery is once, so the collected answer is no longer held by the surface.
The page offers the way back from its receipt only while it would still reach the caller, and collection has already happened.

### What already exists

The surface's delivery route accepts a `within` query parameter and answers `204` when nothing arrived in time, and an abandoned wait consumes nothing, since the wait is bound to the request's context.
`agora/command/poll.go` always passes `surface.PollDeadline` and exposes neither.

### Where a fix could live

- `socrates agora poll --within <duration>`, reaching the `within` the route already accepts, so a session can wait in slices shorter than its harness's foreground bound and repeat until an answer arrives.
  Nothing is lost between slices, because nothing is consumed until an answer is handed back.
- The skills that name `socrates agora poll`, `socrates-design`, `socrates-tickets` and the two review skills, teaching a stated cadence for it the way `socrates-review` already teaches polling `socrates review status`.

## Status

Not fixed.
The relay in this report was started by hand for one session.
