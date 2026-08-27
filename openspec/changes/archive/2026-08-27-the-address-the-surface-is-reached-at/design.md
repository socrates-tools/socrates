## Context

See proposal.md for why.

What exists today: `--address` defaults to `surface.DefaultAddress`, which is `127.0.0.1:0`, and `CheckLoopback` refuses anything that is not loopback before a server starts.
`serve` prints `open.Address` at both of its exits, the one that started a server and the one that found a live surface and revised the document on it.
`hostIsTheSurface` admits a request only when its `Host` splits into a host and a port, the port equals the port the surface bound, and the host is `localhost` or a loopback address.

socrates has a configuration file at `$XDG_CONFIG_HOME/socrates/config.yaml`, decoded with unknown fields refused.
agora has none.

## Goals / Non-Goals

**Goals:**

- An operator states, once per machine, the loopback address the surface binds and the address it is reached at.
- A session that has been told nothing beyond the skill hands a person an address that opens.
- The bind and the host guard are reached by exactly what reaches them today.

**Non-Goals:**

- Discovering what is in front of the surface. agora is told or it is not.
- Verifying that the stated address serves anything.
- Anything per-run, per-document or per-repository. What is stated belongs to the machine.

## Decisions

### agora gets its own file rather than a section of socrates'

socrates' loader refuses a field it does not know, so an agora section in socrates' file is a field socrates has to learn, and socrates learning the surface's address is socrates knowing about the surface.
It does not know about it today, and the two being separate programs is an invariant with a test behind it.

So `$XDG_CONFIG_HOME/agora/config.yaml`, resolved the way socrates resolves its own, and absent by default.

Considered and rejected: one file for both, keyed by program.
It saves an operator one path and costs the repository the one boundary that keeps `agora` from becoming socrates' third name.

### The file is read in `internal/surface`, not in a new package

What it holds is two addresses, and `DefaultAddress` and `CheckLoopback` are already in `internal/surface`.
The package already resolves the user's directories and writes the session file, so nothing new about the filesystem arrives with it.

Considered and rejected: a second package beside socrates' `internal/config`.
Two packages whose names both mean configuration, one per program, is a pair a reader has to keep straight for the life of the repository, and it buys separation the two programs already have by not importing each other.

### The stated address is an absolute `http` or `https` URL, and that is checked

It is handed to a person to open, so the one property agora can check is that it is shaped like something openable.
An operator who writes a bare hostname gets a refusal at startup rather than a person getting a search results page, and the check costs one parse.

That it is refused for shape and never for reachability is the whole of what agora knows: a probe issued from where the surface runs proves nothing about the machine the page is served to, which is the case this exists for.

Considered and rejected: print whatever is written, verbatim.
The failure it permits is silent and lands on the person the round was for.

### The host guard is not widened, and what is in front of the surface presents the host it bound

This is the load-bearing decision.
Admitting the stated host would make the stated address a thing the surface serves on rather than a thing it says, and a host allowlist an operator can extend is the guard turned off by whoever edits that file next.
The specification already refuses a flag that widens the bind on the reasoning that it is set once and forgotten, and this would be the same flag under a name that does not say what it does.

What it costs is one line wherever the proxy is configured: rewrite the upstream `Host` to the address the surface bound, host and port both, because the guard compares the port too.
What it buys is that the guard has no configuration surface at all, which is the property worth more.

Considered and rejected: admit the stated host as well as the bound one.
It works with an unmodified proxy and it makes the reduction from "the surface accepts what it bound" to "the surface accepts what it bound and whatever this file says" for a saving of one proxy directive.

### A flag beats the file, and there is no environment variable

An operator states a default; a caller who passes `--address` has decided.

No environment variable, for the reason the specification already gives about a flag set once and forgotten: a variable is that failure with less to read back.
A file is something somebody wrote and can open again.

### Both of `serve`'s exits print the same thing

The exit that revises a live surface prints `open.Address` as the exit that started one does.
A caller cannot tell those two paths apart and must not have to, so what is printed is resolved once and used by both.

### The skill is not touched

It says to hand the person the address `agora serve` prints.
That sentence is already right in both arrangements once `serve` prints the right thing, and the alternative is a generated document carrying one machine's port into every repository a session reads it in.

## Risks / Trade-offs

- **An operator states an address nothing serves** → agora prints a claim it cannot check, and the person gets a page that does not load. Nothing detects it, and deliberately so: the only probe agora could issue comes from the wrong side of whatever is in front of it. The failure is visible on the first round a person is handed, and the fix is one line in one file the operator wrote.
- **A configured bind is a fixed port, and a fixed port can be taken** → where today the operating system would choose another, a second `agora serve` on a busy port fails. It is bounded by behaviour that already exists: `serve` reuses a surface that is live and reachable, revising the document on it rather than starting a second server, so the collision is with something that is not agora rather than with another review.
- **The stated address and what the proxy actually serves can drift** → they are two strings maintained in two places by the same person, and nothing reconciles them. Accepted: reconciling them means agora knowing what the proxy is, which is the thing this change exists to avoid.

## Migration Plan

None.
An absent file is today's behaviour exactly: bind `127.0.0.1:0`, print what was bound.
Removing the file returns to it.
