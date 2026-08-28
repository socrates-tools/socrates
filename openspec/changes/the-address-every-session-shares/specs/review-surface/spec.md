## ADDED Requirements

### Requirement: The surface mints a session and the caller names it thereafter

`serve` given no session SHALL mint one, and SHALL serve that session's page under a path holding it.
A caller SHALL name a session to poll it, to end it, or to write a new document to it.

`--session` naming a session the surface does not hold SHALL be refused, and SHALL NOT bring one into existence.
`serve` given no session SHALL therefore be the one way a session is created.

The surface owns the identifier so that every one is unique and a safe path segment without a caller having to invent either, and so that two callers who know nothing of each other cannot land on the same page by choosing the same name.
One way to create and one way to address is also the smaller rule: a flag that names an existing session and creates a missing one is a flag with two meanings, and the second only reveals itself when a caller mistypes.

A caller who loses the identifier has no way to ask the surface for it, and that is accepted rather than solved.
The recovery is to serve again and hand over a new address, which is correct behaviour, and the session that was lost is collected when it goes idle.

#### Scenario: A session is minted and the address holds it

- **WHEN** `agora serve` is given a valid document and no session
- **THEN** it mints a session, prints an address holding that session, and the page is reachable there

#### Scenario: Two sessions are served at once and neither disturbs the other

- **WHEN** two documents are served without naming a session, and a person answers on the second
- **THEN** each has its own address and its own page, `poll` on the second returns that answer, and `poll` on the first still waits

#### Scenario: A session the surface never minted is refused

- **WHEN** a caller names a session the surface does not hold
- **THEN** it is refused, the session is named, and no session is created

#### Scenario: A session is named to poll it

- **WHEN** two sessions are open and `poll` is given one of them
- **THEN** it returns what was answered on that session, and reports nothing from the other

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
Where nothing is in front of the surface those are one address and the distinction costs nothing; where something is, the address the surface bound opens for the process that printed it and for nobody it was served to.

Two binaries in one module is ordinary in Go, and a second repository for a tool with one consumer is overhead with no payer.
This does not disturb the rule that `socrates` and `scrat` are one program under two names: `agora` is a different program, not a third name.

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

#### Scenario: socrates keeps its two names and gains no third

- **WHEN** the binaries the repository builds are examined
- **THEN** `socrates` and `scrat` are one program, and `agora` is a separate one

### Requirement: The surface binds a loopback address and nothing else

The server SHALL bind a loopback address.
No flag, no environment variable, no document field and no configured address SHALL widen that bind, and there SHALL be no tunnel, no share command and no third-party host.

The address it binds when nothing states one SHALL be a fixed loopback port rather than a port the operating system chooses.
A port that moves is a port an operator cannot put anything in front of and a person cannot return to, and it bought nothing once the server holds many sessions: what one server on one port used to be short of was a way to tell its pages apart, which is now the path rather than the port.

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

#### Scenario: The default bind is a fixed loopback port

- **WHEN** the surface is started with no address stated by flag and none configured
- **THEN** it binds the fixed default port on loopback, and a surface started again after it stops binds the same one

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

### Requirement: A surface answers only the session it was started for

Every route that changes state SHALL require a token, and that token SHALL be minted per session rather than per server.
A session's token SHALL be embedded in the page that session served, and SHALL answer for that session and no other.

Another process on the same machine can reach a loopback port.
Without the token it could answer on the person's behalf, and the caller would record a decision nobody made.

A server holding many sessions makes the second half of that necessary as well as the first.
One token for every page would mean that whoever holds the page for one round can answer every other round open on the machine, and those rounds are about different repositories and different people's work.

#### Scenario: An answer without the token is rejected

- **WHEN** an answer is submitted without the token the page carries
- **THEN** it is rejected, and `poll` continues to wait

#### Scenario: One session's token does not answer for another

- **WHEN** two sessions are open and an answer is submitted to the first carrying the second's token
- **THEN** it is rejected, and `poll` on the first continues to wait

### Requirement: Waiting for a person is bounded and says so when nobody comes

`poll` SHALL wait on a person rather than on a machine, and SHALL therefore wait far longer than a poll on a process would.
It SHALL still be bounded, and on reaching that bound SHALL return saying nobody answered rather than continuing to block.

A session with no reader for the idle period SHALL be collected, and the server SHALL shut itself down when it holds no sessions.
Collection SHALL be per session, so that a page nobody opened costs the pages people are reading nothing.

A caller that blocks forever is a caller that has to be killed, and an unattended session that never returns reports a surface nobody opened as a surface still being read.

#### Scenario: Nobody answers

- **WHEN** the bound passes with no answer sent
- **THEN** `poll` returns saying nobody answered, and does not report an answer

#### Scenario: A poll that started before the page was open still waits

- **WHEN** `poll` is started in the window between `serve` returning and the page being opened
- **THEN** it waits rather than failing, and returns the answer when it comes

#### Scenario: An idle session is collected and the ones being read are not

- **WHEN** two sessions are open and one has had no reader for the idle period
- **THEN** that session is collected, and the other is still reachable and still answerable

#### Scenario: An idle server stops

- **WHEN** every session has been collected and none is open
- **THEN** the server exits, and the port is free

### Requirement: A caller who writes a new document is asking again

Writing a new document to a session whose exchange has ended SHALL reopen that exchange, and the surface SHALL then accept an answer.
`serve` SHALL NOT return an address whose only possible reply to an answer is a refusal.

A document written to a named session SHALL reach that session and no other.

A caller reaches this path in one situation, and it is the situation where it matters: the person ended the exchange without answering, so the caller asks again.
Handing back the address of an exchange that can only refuse sends that person to a page that was never going to work, and they find out after deciding a second time.
Reopening is the smaller of the two available rules and it changes nothing while an exchange is open, where the exchange is not over and there is nothing to reopen.
It also needs no new command: recovering is serving a document, which `serve` already does, so the five commands stay five.

#### Scenario: Serving again onto an ended exchange yields an answerable page

- **WHEN** the person ends an exchange without answering and the caller serves a document to that session again
- **THEN** the address returned holds a page that accepts an answer, and `poll` delivers what the person sends from it

#### Scenario: Serving again while the exchange is open is unchanged

- **WHEN** the caller serves a new document to a session whose exchange is open
- **THEN** the page takes up the new content and the exchange is as answerable as it was

#### Scenario: A new document reaches the session it named and no other

- **WHEN** two sessions are open and a new document is served to one of them
- **THEN** that page takes up the new content, and the other page is unchanged
