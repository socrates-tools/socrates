## Context

See proposal.md for why.

What exists today: `surface.Serve(document, address, token, idle)` binds, holds one `Exchange` over one document, and mints one token that guards every state-changing route.
`cmd/agora serve` reads `session.json` from agora's state home; if a surface is live it revises the document on it, and otherwise it spawns itself detached with `AGORA_SERVING=yes` and the document on standard input, and the child writes `session.json` and `document.json`.
`poll`, `end` and `stop` read that file and act on the one surface it names.
The page reads its token from `document.body.dataset.token` and issues three fixed requests: `POST /answer`, `GET /` and `EventSource('/live')`.

`hostIsTheSurface` compares the request's `Host` port against the bound port, and `CheckLoopback` refuses a bind that is not loopback.
An operator may state a bind and the address the surface is reached at in `config.yaml` at the top of agora's home.

## Goals / Non-Goals

**Goals:**

- Two rounds a person reads are two pages, open at once, neither able to disturb the other.
- One address that does not move between runs, so an operator configures a port once and a person can return to a page they left.
- The identifier of a page is minted by the surface and is a safe path segment by construction.
- Nothing the surface admits gets wider.

**Non-Goals:**

- Listing what the surface holds. There is no sixth command, and a caller that lost an identifier serves again.
- Reviving a session the surface no longer holds. A collected session is gone.
- Sharing anything between sessions beyond the assets and the server itself.
- Anything about which caller a session belongs to. The surface knows sessions, not owners.

## Decisions

### The surface mints the identifier, and `--session` never creates one

The caller could name its own session, and socrates' skill would pass the run's name, which has the property that a caller who lost it can recompute it from `socrates status`.
That is a real advantage and it is outweighed twice.

A caller-invented name goes into a URL path, so the surface would have to refuse anything that is not a bounded, safe segment: a gate to write, to test and to get wrong, which minting removes entirely.
And `--session` would then mean "that session, or a new one if you have not seen this name", which is creation and addressing in one flag, where the second meaning only shows itself when somebody mistypes.
Minting leaves one way to create and one way to address, and makes two callers landing on the same page impossible rather than making it their problem.

The cost is that the identifier lives only in the caller's memory, and the surface offers no way to ask for it back.
This is accepted rather than solved: the recovery is to serve again and hand over a new address, which is correct behaviour, and per-session collection makes the abandoned page cost nothing.
A recovery mechanism would mean a command that lists what the surface holds, and the five commands are five.

Considered and rejected: keying a session by the canonical path of the document, as lavish keys an artifact by its file.
agora's documents routinely arrive on standard input, so there is no stable file to key on, and the property that makes it work there is absent here.

The identifier is eight bytes from `crypto/rand`, rendered hex.
Sixteen characters is short enough to read back off a URL and wide enough that a local process cannot enumerate open pages.

### Two kinds of token, and the rule is one sentence

Today one token guards every state-changing route and is both what the page carries and what the CLI sends.
With many sessions that single token would let whoever holds one page answer every other round open on the machine, and those rounds are about different repositories and different people's work.

So: a control token, minted when the server starts and written into the surface record that only the user can read, and a page token per session, embedded in the page that session served.

The rule: **every route under a session's path accepts that session's token, and the control token is accepted everywhere.**
Routes that are not under a session's path, which are opening a session and stopping the server, accept the control token only.

This is strictly tighter than what exists.
A page can answer, close and stream its own session and can reach nothing else; the CLI, which can already read the state directory, keeps doing everything.
It also settles `close`, which both the page and `agora end` call, without a special case: both carry a token the session accepts.

### The detached child becomes a plain server, and sessions arrive over HTTP

Today the spawned child is given the document on standard input and comes up already serving it.
That does not survive a server holding many sessions: the second `serve` finds a live surface and has to open a session on it, so opening a session has to be something that can happen over the wire.

Once it can, the first `serve` should use the same path rather than a second one.
The child comes up holding nothing and writes the surface record; the parent then opens its session over HTTP exactly as every later `serve` does.
One code path, and the case that used to be special (the very first document) stops being special.

Considered and rejected: keeping standard input for the first session.
It saves one round trip on one invocation and leaves two ways for a document to reach the surface, one of which is exercised only by the first caller after a reboot.

### The fixed port is 18018, and a port held by something else is a refusal

`serve` reaches the surface record, and if a surface answers there it opens a session on it.
Where nothing answers, it forgets the record and binds.
Where the bind fails because something holds the port and that something is not a surface, `serve` refuses, names the port, and names `--address`.

It does not stop what is there.
lavish resolves this by asking the occupant to shut down and killing it by PID if it will not, which it needs because its client and server versions must match.
agora has no such coupling, and a tool that kills a process it did not start because it wanted that process's port is a tool that will one day kill the wrong one.

18018 is stated once, in `internal/surface`, as the default the flag and the configured bind both override.

### The session is a path segment, and the assets are not under it

The page is served at `/sessions/{session}`, and its routes are under it.
The assets stay at `/assets/`, shared by every session, so a second page opened in the same browser paints from cache and the embedded bundle is served once.

The page is rendered with the base its requests are made against, carried in `Frame` beside the token and revision it already carries and emitted as a `data-` attribute on the body, as the token already is.
The client's three fixed paths are then built from that base rather than written absolute.
Making them relative to the document URL instead would work only while the page's own path has a trailing slash, which is a property of the URL a person typed rather than one the page can hold.

`hostIsTheSurface` is untouched.
It compares the port, and the port is now the same for every session, which is the case it was already written for.

### Collection is per session, and the server outlives any one of them

`Exchange` already knows when it was last touched and how many readers it has, so per-session collection is the existing rule applied to each entry rather than to the server.

The server stops when it holds no sessions and has held none for the idle period.
Stating it that way rather than "when the last session is collected" covers the server that comes up and is never given a session, which is otherwise a process that lives forever because nothing ever ends.

### The state directory gains a directory per session

`<home>/surface.json` holds the address, the process and the control token.
`<home>/sessions/<session>/` holds that session's record and the document it was served, and is removed when the session is collected.

This mirrors socrates' run directories, which is deliberate: a person opening either state directory finds the same shape, and a session is legible with `ls` and `cat` at the moment the surface itself is what is broken.

## Risks / Trade-offs

- **A caller loses the identifier and cannot ask for it back** → serve again and hand over a new address. The abandoned session is collected after the idle period and costs one directory and one map entry until then. Solving it properly means a command that lists sessions, which is a sixth command; the cheap escape hatch, if this ever bites, is to let `--session` recreate a session under an identifier the surface minted earlier, which buys stable identifiers at the price of a ledger of minted identifiers that outlives the sessions and needs its own retention rule. Not now.
- **A fixed port can be held by something that is not a surface** → `serve` refuses and names both the port and the flag that moves it. Where today the operating system would silently choose another port, a caller now gets a message. That is the trade the fixed port is for, and the refusal is the visible half of it.
- **One server holds every page, so `stop` ends every round being read** → `stop` is not in the loop the skill teaches, which names only `serve` and `poll`, so no review runs it. It is the operator's blunt instrument and it says so. Ordinary cleanup is `end`, which is one session, and collection, which is automatic.
- **A page token and a control token are two things where there was one** → the rule is one sentence and the tighter half is the one a browser holds. The failure it prevents is real and would be invisible: answering the wrong round from the wrong page.
- **Two sessions on one server share a process** → a panic serving one page takes down every page. This is already true of every request the server handles, and the answer is the same one: the caller's `poll` returns, the surface record is stale, and the next `serve` starts a server again.

## Migration Plan

None.
An existing `session.json` and `document.json` at agora's state home name a surface the new server does not recognise; `serve` finds nothing answering, forgets them and binds.
Anything a person had open is lost, which is what stopping the old binary would have done anyway.
