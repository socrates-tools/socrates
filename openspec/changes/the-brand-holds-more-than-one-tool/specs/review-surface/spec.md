## RENAMED Requirements

- FROM: `#### Scenario: socrates keeps its two names and gains no third`
- TO: `#### Scenario: socrates-review has one name and agora is a different program`

## MODIFIED Requirements

### Requirement: One binary, five commands, beside socrates

The repository SHALL build a second binary, `agora`, from `cmd/agora`.
It SHALL offer `serve`, `poll`, `end`, `stop` and `version`, and no others.

`serve` SHALL return as soon as the page is reachable, printing the address that page is reached at, so the caller is not blocked while a person reads.
That address SHALL hold the session the page belongs to, composed of the address the surface is reached at, the path sessions are served under, and the session itself.
Where an operator has stated the address the surface is reached at, `serve` SHALL compose from what they stated; where none is stated, it SHALL compose from the address it bound.

`poll` SHALL block until the person sends and SHALL print what they answered.
`end` SHALL close a surface the caller no longer needs.
Each SHALL be told which session it is about, because the server holds more than one.

`stop` SHALL shut the server down and SHALL name no session.
The two are deliberately not the same act: `end` finishes one exchange and leaves every other page being read, and `stop` takes the server down and every session with it.
A per-session `stop` would be `end` under a second name, and the surface holds five commands rather than a synonym.

A caller hands what `serve` printed to a person, so what `serve` prints has to be what that person can open.
Where nothing is in front of the surface those are one address and the distinction costs nothing.
Where something is, the address the surface bound opens for the process that printed it and for nobody it was served to.

Two binaries in one module is ordinary in Go, and a second repository for a tool with one consumer is overhead with no payer.
This does not disturb the rule that the tool answers to one name: `agora` is a different program, not a second name for `socrates-review`.

#### Scenario: serve returns before the person has read anything

- **WHEN** `agora serve` is given a valid document
- **THEN** it prints the address and exits, and the page is reachable at that address

#### Scenario: serve prints the address stated over the address bound

- **WHEN** an operator has stated the address the surface is reached at, and `agora serve` is given a valid document
- **THEN** the address it prints is composed from the stated address rather than from the address it bound, and holds the session it minted

#### Scenario: The surface offers nothing beyond the five

- **WHEN** the binary's command surface is examined
- **THEN** it holds exactly `serve`, `poll`, `end`, `stop` and `version`, and nothing that publishes, exports, shares or annotates

#### Scenario: Ending one session leaves the others being read

- **WHEN** two sessions are open and one is ended
- **THEN** that page reports the exchange over, and the other is still reachable and still answerable

#### Scenario: Stopping is the server, not a session

- **WHEN** two sessions are open and the server is stopped
- **THEN** neither page is reachable, and the port is free

#### Scenario: socrates-review has one name and agora is a different program

- **WHEN** the binaries the repository builds are examined
- **THEN** `socrates-review` answers to one name, and `agora` is a separate program rather than a second name for it
