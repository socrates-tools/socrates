# review-surface Specification

## Purpose
TBD - created by archiving change a-page-that-answers-back. Update Purpose after archive.
## Requirements
### Requirement: One binary, five commands, beside socrates

The repository SHALL build a second binary, `agora`, from `cmd/agora`.
It SHALL offer `serve`, `poll`, `end`, `stop` and `version`, and no others.

`serve` SHALL return as soon as the page is reachable, printing the address, so the caller is not blocked while a person reads.
`poll` SHALL block until the person sends and SHALL print what they answered.
`end` SHALL close a surface the caller no longer needs, and `stop` SHALL shut the server down.

Two binaries in one module is ordinary in Go, and a second repository for a tool with one consumer is overhead with no payer.
This does not disturb the rule that `socrates` and `scrat` are one program under two names: `agora` is a different program, not a third name.

#### Scenario: serve returns before the person has read anything

- **WHEN** `agora serve` is given a valid document
- **THEN** it prints the address and exits, and the page is reachable at that address

#### Scenario: The surface offers nothing beyond the five

- **WHEN** the binary's command surface is examined
- **THEN** it holds exactly `serve`, `poll`, `end`, `stop` and `version`, and nothing that publishes, exports, shares or annotates

#### Scenario: socrates keeps its two names and gains no third

- **WHEN** the binaries the repository builds are examined
- **THEN** `socrates` and `scrat` are one program, and `agora` is a separate one

### Requirement: The surface binds a loopback address and nothing else

The server SHALL bind a loopback address.
No flag, no environment variable and no document field SHALL widen that bind, and there SHALL be no tunnel, no share command and no third-party host.

The whole reason this exists is that the document quotes the code it was written about.
A surface that can be reached from another machine has given that away, and a flag that widens the bind is a flag somebody sets once and forgets.

#### Scenario: A wildcard bind is refused rather than quietly narrowed

- **WHEN** the surface is asked to bind an address that is not loopback
- **THEN** it refuses, names the address, and starts no server

#### Scenario: The page reaches nothing

- **WHEN** the page is loaded and every request it issues is recorded
- **THEN** every one resolves to the surface itself, and none leaves the machine

#### Scenario: A request claiming another host is rejected

- **WHEN** a request arrives naming a host the surface did not bind
- **THEN** it is rejected before any handler runs

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

Two sends that a caller cannot tell apart make a caller either stop too early or wait for a person who has finished.

#### Scenario: An answer round-trips

- **WHEN** a person chooses in two decisions, writes a message, and sends without ending
- **THEN** `poll` prints both answers with their commands, the message, and that the exchange is not over

#### Scenario: Ending is distinguishable

- **WHEN** a person sends and ends
- **THEN** `poll` prints that the exchange is over, and a later `poll` does not reopen it

#### Scenario: Delivery happens once

- **WHEN** an answer has been delivered to one `poll`
- **THEN** a second `poll` does not deliver it again
