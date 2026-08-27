# review-surface Specification

## Purpose
TBD - created by archiving change a-page-that-answers-back. Update Purpose after archive.
## Requirements
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

### Requirement: A surface answers only the session it was started for

Every route that changes state SHALL require a token minted when the surface started and embedded in the page it served.

Another process on the same machine can reach a loopback port.
Without the token it could answer on the person's behalf, and the caller would record a decision nobody made.

#### Scenario: An answer without the token is rejected

- **WHEN** an answer is submitted without the token the page carries
- **THEN** it is rejected, and `poll` continues to wait

### Requirement: Session state lives outside the repository being discussed

The session file, the served document and anything else the surface writes SHALL live under the user's own state directory.
Nothing SHALL be written inside the caller's working directory or the repository it is discussing.

socrates already holds that nothing writes in the repository under review, and a surface that dropped a directory there would break it from the outside.

#### Scenario: The repository is untouched

- **WHEN** a surface is served, answered and ended from within a git repository
- **THEN** that repository's `git status` is byte-identical to what it was before

### Requirement: Waiting for a person is bounded and says so when nobody comes

`poll` SHALL wait on a person rather than on a machine, and SHALL therefore wait far longer than a poll on a process would.
It SHALL still be bounded, and on reaching that bound SHALL return saying nobody answered rather than continuing to block.

The server SHALL shut itself down after an idle period with no session open.

A caller that blocks forever is a caller that has to be killed, and an unattended session that never returns reports a surface nobody opened as a surface still being read.

#### Scenario: Nobody answers

- **WHEN** the bound passes with no answer sent
- **THEN** `poll` returns saying nobody answered, and does not report an answer

#### Scenario: A poll that started before the page was open still waits

- **WHEN** `poll` is started in the window between `serve` returning and the page being opened
- **THEN** it waits rather than failing, and returns the answer when it comes

#### Scenario: An idle server stops

- **WHEN** no session has been open for the idle period
- **THEN** the server exits, and the port is free

### Requirement: What poll prints is data, and it says whether the exchange is over

`poll` SHALL print, for each decision, the identifier, what was chosen, the words the person added, and the command the chosen option carries.
It SHALL print any message the person wrote to the caller.
It SHALL state whether the person ended the exchange or expects the caller to answer and re-render.

A poll SHALL NOT reopen an exchange the person ended.
Writing a new document is the one act that does, and it is the caller's.

Two sends that a caller cannot tell apart make a caller either stop too early or wait for a person who has finished.
The two rules are stated together because they are the pair that is easy to confuse: reading an ended exchange leaves it ended, and asking a new question does not.

#### Scenario: An answer round-trips

- **WHEN** a person chooses in two decisions, writes a message, and sends without ending
- **THEN** `poll` prints both answers with their commands, the message, and that the exchange is not over

#### Scenario: Ending is distinguishable

- **WHEN** a person sends and ends
- **THEN** `poll` prints that the exchange is over, and a later `poll` does not reopen it

#### Scenario: Delivery happens once

- **WHEN** an answer has been delivered to one `poll`
- **THEN** a second `poll` does not deliver it again

#### Scenario: A poll over an ended exchange leaves it ended

- **WHEN** the person has ended the exchange and `poll` is run again
- **THEN** it reports that nobody answered rather than reopening the exchange, and the page stays closed

### Requirement: A caller who writes a new document is asking again

Writing a new document to a surface whose exchange has ended SHALL reopen that exchange, and the surface SHALL then accept an answer.
`serve` SHALL NOT return an address whose only possible reply to an answer is a refusal.

A caller reaches this path in one situation, and it is the situation where it matters: the person ended the exchange without answering, so the caller asks again.
Handing back the address of an exchange that can only refuse sends that person to a page that was never going to work, and they find out after deciding a second time.
Reopening is the smaller of the two available rules and it changes nothing while an exchange is open, where the exchange is not over and there is nothing to reopen.
It also needs no new command: recovering is serving a document, which `serve` already does, so the five commands stay five.

#### Scenario: Serving again onto an ended exchange yields an answerable page

- **WHEN** the person ends an exchange without answering and the caller serves a document to that surface again
- **THEN** the address returned holds a page that accepts an answer, and `poll` delivers what the person sends from it

#### Scenario: Serving again while the exchange is open is unchanged

- **WHEN** the caller serves a new document to a surface whose exchange is open
- **THEN** the page takes up the new content and the exchange is as answerable as it was
