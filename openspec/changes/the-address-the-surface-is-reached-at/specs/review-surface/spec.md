## MODIFIED Requirements

### Requirement: One binary, five commands, beside socrates

The repository SHALL build a second binary, `agora`, from `cmd/agora`.
It SHALL offer `serve`, `poll`, `end`, `stop` and `version`, and no others.

`serve` SHALL return as soon as the page is reachable, printing the address the surface is reached at, so the caller is not blocked while a person reads.
Where an operator has stated that address, `serve` SHALL print what they stated; where none is stated, `serve` SHALL print the address it bound.
`poll` SHALL block until the person sends and SHALL print what they answered.
`end` SHALL close a surface the caller no longer needs, and `stop` SHALL shut the server down.

A caller hands what `serve` printed to a person, so what `serve` prints has to be what that person can open.
Where nothing is in front of the surface those are one address and the distinction costs nothing; where something is, the address the surface bound opens for the process that printed it and for nobody it was served to.

Two binaries in one module is ordinary in Go, and a second repository for a tool with one consumer is overhead with no payer.
This does not disturb the rule that `socrates` and `scrat` are one program under two names: `agora` is a different program, not a third name.

#### Scenario: serve returns before the person has read anything

- **WHEN** `agora serve` is given a valid document
- **THEN** it prints the address and exits, and the page is reachable at that address

#### Scenario: serve prints the address stated over the address bound

- **WHEN** an operator has stated the address the surface is reached at, and `agora serve` is given a valid document
- **THEN** it prints the stated address rather than the address it bound

#### Scenario: The surface offers nothing beyond the five

- **WHEN** the binary's command surface is examined
- **THEN** it holds exactly `serve`, `poll`, `end`, `stop` and `version`, and nothing that publishes, exports, shares or annotates

#### Scenario: socrates keeps its two names and gains no third

- **WHEN** the binaries the repository builds are examined
- **THEN** `socrates` and `scrat` are one program, and `agora` is a separate one

### Requirement: The surface binds a loopback address and nothing else

The server SHALL bind a loopback address.
No flag, no environment variable, no document field and no configured address SHALL widen that bind, and there SHALL be no tunnel, no share command and no third-party host.

An operator MAY state the loopback address the surface binds, and that address SHALL be refused by the same check a stated flag is refused by.
A flag SHALL override it, so what an operator configured is a default rather than a decision taken away from the caller.

An operator MAY state the address the surface is reached at.
It is handed to a person to open, so it SHALL be refused unless it is shaped like something openable, and it SHALL be refused for its shape alone: whether anything serves it is not something the surface can answer from where it runs.
That address SHALL reach neither the bind nor the guard on what host a request may name: it is what the surface says about itself and never what it admits.
A request naming it SHALL be turned away exactly as a request naming any host the surface did not bind is turned away.

The whole reason this exists is that the document quotes the code it was written about.
A surface that can be reached from another machine has given that away, and a flag that widens the bind is a flag somebody sets once and forgets.
An address the surface merely states gives nothing away, because nothing is served on the strength of it.
Admitting a host on the strength of it would be the same forgotten flag under another name, so anything put in front of the surface presents the host the surface bound.

#### Scenario: A wildcard bind is refused rather than quietly narrowed

- **WHEN** the surface is asked to bind an address that is not loopback
- **THEN** it refuses, names the address, and starts no server

#### Scenario: A configured bind that is not loopback is refused too

- **WHEN** the address an operator stated for the bind is not loopback
- **THEN** it refuses, names the address, and starts no server

#### Scenario: A stated flag beats a configured bind

- **WHEN** an operator has stated the address to bind and a caller passes a different loopback address
- **THEN** the surface binds the address the caller passed

#### Scenario: An address a person could not open is refused

- **WHEN** the address an operator stated for the surface to be reached at is not shaped like something openable
- **THEN** it refuses, names the address, and starts no server

#### Scenario: An address nothing answers on is still served

- **WHEN** the address an operator stated is shaped like something openable and nothing serves it
- **THEN** the surface starts and prints it, because what answers there is not something it can see from where it runs

#### Scenario: The page reaches nothing

- **WHEN** the page is loaded and every request it issues is recorded
- **THEN** every one resolves to the surface itself, and none leaves the machine

#### Scenario: A request claiming another host is rejected

- **WHEN** a request arrives naming a host the surface did not bind
- **THEN** it is rejected before any handler runs

#### Scenario: The address the surface is reached at is not a host it admits

- **WHEN** an operator has stated the address the surface is reached at, and a request arrives naming that address
- **THEN** it is rejected before any handler runs, exactly as any other host the surface did not bind is
