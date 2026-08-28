## Context

See proposal.md for why.

`The surface binds a loopback address and nothing else` already settles the shape of this.
The bind is loopback and no flag, variable, document field or configured address widens it.
An operator may state the address the surface is reached at, it is refused for its shape alone, and it reaches neither the bind nor the guard on what host a request may name.

That requirement already contains the sentence this change leans on: whether anything serves the stated address "is not something the surface can answer from where it runs".
Where the person's browser is, is the same question one step further out.

## Goals / Non-Goals

**Goals:**

- A person who cannot open the address finds out why from the thing that printed it.
- An install where the page opens first time, with the one step stated once.

**Non-Goals:**

- Widening the bind, under a flag, a configuration field, an environment variable or anything else.
- Shipping a forwarder, or spawning one. A child process binding what the parent refuses to bind is the same decision with an extra hop.
- Detecting where the person is.
- Teaching the skill anything new.
- Making the page reachable with no step at all. That cannot be done without reversing the bind rule, and the bind rule is the reason the document may quote the code.

## Decisions

### It states where the page is served, rather than guessing whether the person can reach it

Detecting a container was considered and rejected: `/.dockerenv`, a cgroup naming docker or containerd, `/run/.containerenv`.
Every one is a heuristic, and a heuristic is wrong in both directions.
Wrong one way it warns somebody whose browser is right there and who now distrusts a working address.
Wrong the other way it says nothing to the person it exists for, which is the whole failure again.

The fact needs no heuristic.
The page is reachable from the machine it is served on and from nowhere else, and that is true on every machine, in a container or not.
A person reading it either knows their browser is elsewhere or learns the address is not the problem, and both are better than an empty reply.

This is the same position the requirement already takes about the reached-at address, for the same reason: the surface answers about where it serves, never about what is true somewhere it cannot see.

### It goes on the stream a person reads

`agora serve` prints the address and the caller reads it.
Anything else on that stream is a caller's parse breaking, so the sentence goes on the error stream and the address stream carries the address and nothing else.

The cost, stated rather than hidden: the caller is usually a coding session rather than the person, so the sentence lands in the session's output. That is the right place anyway. The session is what hands the address over, and a session told the page is machine-local can say so in words fitted to whoever it is talking to. The skill gains nothing for this, because the sentence is already in front of whatever ran the command.

### A stated reached-at silences it

An operator who stated the address the surface is reached at has answered where the page is opened from, and repeating the question at them is noise on every serve forever.
It also gives the sentence an end: follow what it says, and it stops.

### The step is stated at install, and the sentence is what makes a skipped step diagnosable

The forward is what makes it work; the sentence is what makes forgetting it cost a minute rather than an afternoon.
The fixed port is what lets the forward be stated once, which is what that port was for.

The address a person is handed names socrates.
agora is a dependency with one consumer, and somebody who installed socrates has no reason to learn the name of the thing serving the page.
The binary keeps its own name where it is talking about itself, which `agora version` already settles.

## Risks / Trade-offs

- **One more line on every serve where nothing is configured** -> one line, on the error stream, and stating the reached-at address ends it.
- **The session sees it before the person does** -> covered above, and it is the session that was going to hand over the address anyway.
- **The sentence is true and unhelpful for somebody whose browser is right there** -> it is one line and it tells them the address is fine, which is the thing they would otherwise have to establish themselves.
- **It still does not work with no step at all** -> it cannot, without reversing the bind rule. What this change buys is that the step is stated where somebody will read it, twice: once at install and once at the moment it is missing.

## Migration Plan

None. Nothing that works today stops working.
