# 57. The surface binds one port and asks it who stands there

- Status: accepted
- Date: 2026-10-08
- Part of: agora-port

## Decision

Where nothing states an address, the surface binds `127.0.0.1:24672`, one fixed port, and nothing else.
A port held by something that is not this root's surface is a refusal that names the port, says what holds it where the surface can tell, and says how to state another address.
There is no window and no walk: the ten ports counted down from 18018 are gone.

`serve` finds a surface already standing by asking the port, rather than by reading `surface.json`.
The surface answers who it is: that it is agora, the root it writes under, and the build it was linked from.
`surface.json` keeps the control token and nothing that decides where the surface is, so a stale record cannot point `serve` at the wrong place.
A port held by an agora surface of another root, such as another `SOCRATES_HOME` or another user's, is refused, naming that root and how to state another address.

`surface.bind` in agora's file and `--address` on `serve` are how a person moves the surface, as they are today, and nothing is added beside them.
A port either commits to is bound or refused, as before.

The window existed so the page could be forwarded out of a devcontainer through a range the person had opened once.
That reason is going away, and a single port means one thing to forward and one address that never silently moves.
This follows lavish-axi's shape, with a default of agora's own.

24672 spells AGORA on a phone keypad, is unassigned in the IANA registry, and sits below every common ephemeral range.

This replaces the window ADR 24 introduced and ADR 42 placed in `agora/surface`.
ADR 42's other decisions stand.

## Rejected

- One fixed port, then any free port the system hands out.
  It never refuses, but the address moves, and a forwarded port silently stops matching.
- Keeping the window and changing only its numbers.
- 18018, today's first port, which would have kept existing forwarding and README lines working.
- 4387, lavish-axi's own port, which collides with lavish-axi on any machine running both.
- 6543, chosen first, and given up once it turned out to be registered with IANA as lds-distrib and to be the port Pyramid's pserve development server listens on by default.
- Adding an environment variable beside `surface.bind` and `--address`, a third place to state one thing.
- Replacing both with an environment variable and `--port`, which is lavish-axi's shape exactly and would take the setting out of agora's configuration file.
- Reading `surface.json` to find the surface, as before, which leaves where the surface is decided by a file that can outlive it.
- Refusing a surface of another root by saying only that the port is held, the same as for any foreign program.
