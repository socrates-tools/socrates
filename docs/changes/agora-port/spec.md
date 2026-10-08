# agora binds one port, and asks it who stands there

Decisions: ADR 57, ADR 58.

## Problem Statement

Where nothing states an address, agora walks down a window of ten ports, 18018 to 18009, and binds the first one free.
The window was chosen for one machine: a devcontainer whose person forwarded that range to the Windows host once and could then open any page agora served.
That person is leaving the devcontainer, and on every other machine the window is a cost with nothing to show for it.

A held 18018 does not fail, it moves.
A person who forwarded one port, or wrote `ssh -L 18018:...`, opens an address nothing answers on, while the page they were handed sits on 18017.
The move is silent, so what they see is an empty reply and not a reason.

`serve` finds the surface already standing by reading `surface.json` and reaching the address it names.
The record decides where the surface is, so a record that outlived its surface, or one written under another root, points `serve` somewhere it should not go.

A surface keeps serving whatever build started it.
After `make install` or an upgrade, the page a person opens next is drawn by the old binary until the old surface has held no session for 45 minutes, so a fix to the page is not seen until then.

lavish-axi solves the same problem with less: one default port, overridable, a server found by asking the port, a held port refused by name, and an older server replaced by a newer one.

## Solution

agora binds `127.0.0.1:24672` where nothing states an address, and nothing else.
24672 spells AGORA on a phone keypad, is unassigned in the IANA registry, and sits below every common ephemeral range.

`serve` asks the port who stands there:

- Nothing: it starts a surface on the port, as today.
- This root's surface, from this build: it reuses it, as today.
- This root's surface, from another build, holding a session: it reuses it, so no open page is taken away.
- This root's surface, from another build, holding no session: it stops it through the control token and starts this build's surface on the port.
- An agora surface of another root: it refuses, naming that root, and saying how to state another address.
- Anything else: it refuses, naming the port, and saying how to state another address.

`surface.bind` in `~/.socrates/agora/config.yaml` and `--address` on `serve` move the surface exactly as they do today, and nothing is added beside them.
A port either commits to is bound or refused, never walked past, as today.

`surface.json` keeps the control token, and no longer decides where the surface is.

## User Stories

1. As a person running a design, I want `socrates agora serve` to bind the same port every time on my machine, so that the address I forwarded or bookmarked is the one my page is on.
2. As a person whose port 24672 is taken by another program, I want `serve` to refuse and name the port, so that I know why no page appeared and what to change, rather than opening an address that is empty.
3. As a person whose port is taken, I want the refusal to say that `surface.bind` in agora's file or `--address` moves the surface, so that I can fix it in one edit without reading the README.
4. As a person who already set `surface.bind`, I want that setting to keep working unchanged, so that upgrading costs me nothing.
5. As a person running two roots on one machine, such as a test suite under its own `SOCRATES_HOME` beside my own, I want a `serve` from one root to refuse a surface of the other by naming its root, so that a page is never opened under somebody else's sessions and token.
6. As a person who just ran `make install`, I want the next page I open to be drawn by the build I installed, once nothing is open on the old one, so that a fix to the page is seen without waiting 45 minutes or running `socrates agora stop`.
7. As a person in the middle of answering a page, I want a `serve` from a newer build elsewhere on the machine to leave my page alone, so that I do not lose an answer to an upgrade.
8. As a person whose `surface.json` was left behind by a surface that crashed, I want `serve` to take no notice of where it says the surface is, so that a stale record never points `serve` at the wrong place.
9. As a person forwarding the page out of a container or over ssh, I want the README to name the one port to forward, so that I forward it once.
10. As a contributor, I want the suite to fail when the port the README tells a person to forward is not the port the surface binds, so that the two cannot drift.

## Implementation Decisions

**`agora/surface`**

- `window.go` goes.
  `Addresses`, `FirstFreeIn` and the walk go with it, and so does the refusal naming a whole window.
- The default address is stated once, in `agora/surface` beside the guards that judge an address, as `127.0.0.1:24672`.
  `DefaultAddress` keeps its name.
- The surface answers who it is on a route that takes no token, because a caller asking whether the port is its surface does not yet know whose token to send.
  The answer carries that it is agora, the root it writes under, and the build it was linked from, as `internal/version` reports it, and nothing a caller could act on: not the token, not a session.
  The answer is subject to the same host guard as every other route.
- What counts as another build is any difference in the version, the commit stamp or whether the build was uncommitted.
- The surface reports whether it holds a session, so `serve` can tell a surface it may replace from one it must reuse.
  This may be the identity answer itself or the existing control-tokened route; the implementer chooses, and no session identifier leaves through a route without a token.
- `Record` keeps the token and the process, and its address stops deciding where `serve` reaches.
- New refusals, each in `messages.go`, each naming what the person can do:
  - the port is held by something that is not an agora surface, naming the port, `surface.bind` in agora's file and `--address`;
  - the port is held by an agora surface of another root, naming that root, the port, `surface.bind` and `--address`.

**`agora/command`**

- `serve` with nothing committed binds `DefaultAddress` or refuses, rather than walking.
- `serve` decides from the port's answer rather than from `ReadRecord`, by the six cases under Solution, and reads the token from `surface.json` only once the port has said it is this root's surface.
- Replacing a surface of another build stops it through the existing `POST /stop` with the token, waits until the port is free, and starts this build's surface on it.
  The wait is bounded, and a port that never frees is the refusal for a port held by something else.
- `stop` reaches the surface on the port in the same way.
- `--address` and `surface.bind` keep their meaning, their names and their checks.

**Documents**

- `README.md`'s configuration example and "Reading the page from another machine" name 24672, and the sentence about counting down a window goes.
- `agora/AGENTS.md`: the layout line for `agora/surface` loses "the window walked down from it", the window invariant becomes one fixed port refused when held and found by asking it, and two invariants are added, one for a surface of another root and one for a surface of another build.
- `CHANGELOG.md`'s Unreleased section says, for somebody who has the previous release, that the page moved to 24672, that a held port is now refused rather than walked past, and that whoever forwarded 18018 forwards 24672 instead.
- `review/config/config_test.go` names 18018 only as a value inside a refused block, and is updated only so the suite no longer names the old window.

**Elsewhere**

- `review/fidelity/livesmoke_test.go` picks its address from the window today, and picks a free loopback port instead.
- `internal/e2e`'s free-port range, 20000 to 29999, contains 24672, so a test could take the default from under a test that holds the machine-wide lock.
  The range moves off it, or skips it, so no free-port pick ever lands on the default.

## Testing Decisions

The behaviour is tested by the existing end-to-end address tests in `internal/e2e/address_test.go`, rewritten for one port.
They drive the real binary against real sockets, and the tests that bind the default hold `internal/testing/exclusive`'s machine-wide lock, so they never race each other or another suite run.
`agora/surface`'s own tests keep only the address guards.

Each scenario below is a test named for it, and each fails before the change for the right reason:

- With nothing stated, the surface binds 24672.
- A foreign listener on the default is a refusal naming the port, `surface.bind` and `--address`, and no surface is started on another port.
- An agora surface of another root on the default is a refusal naming that root.
- A surface of this root from another build that holds no session is replaced, and the page it hands back is served by this build.
- A surface of this root from another build that holds a session is reused, and that session still answers.
- A `surface.json` naming an address nothing answers on is not where `serve` reaches.
- A stated or configured address that is held is refused, as today.
- The port the README tells a person to forward is the one the surface binds, as today.

The window tests go: `TestAHeldDefaultIsTheNextPortDownRatherThanARefusal`, `TestTheWindowIsWalkedFromTheTopDownRatherThanAtRandom`, `TestAWholeWindowHeldIsARefusalThatNamesTheWindow` and `TestWhatAPersonIsToldToForwardIsThePortTheWindowWalkedDownTo`, with `agora/surface/window_test.go`.

A surface of another build is simulated by starting a surface whose identity answer differs, by whatever seam the real binary already allows, such as a build of the binary stamped with another commit through the linker flags `make` already passes.
No test reaches into the server in process to fake an identity.

## Out of Scope

- An environment variable for the port.
- A `--port` flag.
- Binding anything other than loopback, and any equivalent of lavish-axi's Tailscale or LAN listeners, link host or allowed hosts.
- Changing the idle period, or stopping the surface the moment its last session ends.
- Migrating anything a person stated: a `surface.bind` naming 18018 keeps binding 18018.

## Further Notes

The window was always a devcontainer accommodation, deferred as its own change while the monorepo change landed.

24672 was not the person's first choice: 6543 was, and it was given up once it turned out to be registered with IANA as lds-distrib and to be the default port of Pyramid's pserve development server.
