# review-surface Specification

## Purpose
TBD - created by archiving change a-page-that-answers-back. Update Purpose after archive.

## Requirements

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

### Requirement: The surface mints a session and the caller names it thereafter

`serve` given no session SHALL mint one, and SHALL serve that session's page under a path holding it.
A caller SHALL name a session to poll it, to end it, or to write a new document to it.

`--session` naming a session the surface does not hold SHALL be refused, and SHALL NOT bring one into existence.
`serve` given no session SHALL therefore be the one way a session is created.

The surface owns the identifier so that every one is unique and a safe path segment without a caller having to invent either.
Two callers who know nothing of each other then cannot land on the same page by choosing the same name.
One way to create and one way to address is also the smaller rule.
A flag that names an existing session and creates a missing one is a flag with two meanings, and the second only reveals itself when a caller mistypes.

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

### Requirement: The surface binds a loopback address and nothing else

The server SHALL bind a loopback address.
No flag, no environment variable, no document field and no configured address SHALL widen that bind, and there SHALL be no tunnel, no share command and no third-party host.

Where nothing states an address, the surface SHALL bind the first free port of a fixed window of loopback ports, taken from the top of the window down, rather than a port the operating system chooses.
The window SHALL be closed and stated in one place, and the port at its top SHALL be what a machine holding none of them binds.
Where every port in the window is held, the surface SHALL refuse, name the window, and start no server.

An ordered closed window is not a port that moves.
What made a port the operating system chooses unusable is that nothing can be aimed at it and nobody can return to it.
Neither is true here.
What an operator aims something at is a port they stated, which pins the bind rather than being walked past.
What a person returns to is the address the surface printed, which is the one it bound and never the one it was asked for.
What it buys is the case the single fixed port refused on, which is a person who has never named a port meeting a failure about one, whose only way out was to learn a flag to get back to where they started.

An operator MAY state the loopback address the surface binds, and that address SHALL be refused by the same check a stated flag is refused by.
A flag SHALL override it, so what an operator configured is a default rather than a decision taken away from the caller.

Where an operator has stated the address to bind, or the address the surface is reached at, the window SHALL NOT be walked.
The port that follows from what they stated is the port bound, and a port held by something the surface does not know of SHALL be refused rather than stepped over.
Stating either is an operator committing to a port, because something in front of the surface is aimed at one.
A surface that moved out from under it would hand a person the address the operator stated, which by then opens onto nothing.
That is the failure serving a page states its reachable address in order to prevent.
Which port is taken where nobody has named one is the only thing the reached-at address reaches.
It reaches neither the check on what may be bound nor the guard on what host a request may name.

An operator MAY state the address the surface is reached at.
It is handed to a person to open, so it SHALL be refused unless it is shaped like something openable, and for its shape alone.
Whether anything serves it is not something the surface can answer from where it runs.
That address SHALL reach neither the bind nor the guard on what host a request may name: it is what the surface says about itself and never what it admits.
A request naming it SHALL be turned away exactly as a request naming any host the surface did not bind is turned away.

The whole reason this exists is that the document quotes the code it was written about.
A surface that can be reached from another machine has given that away, and a flag that widens the bind is a flag somebody sets once and forgets.
An address the surface merely states gives nothing away, because nothing is served on the strength of it.
Admitting a host on the strength of it would be the same forgotten flag under another name, so anything put in front of the surface presents the host the surface bound.

#### Scenario: The default bind is a fixed loopback port

- **WHEN** the surface is started with no address stated by flag and none configured, and no port of the window is held
- **THEN** it binds the fixed port at the top of the window on loopback, and a surface started again after it stops binds the same one

#### Scenario: A held default is the next port down rather than a refusal

- **WHEN** the port at the top of the window is held by something the surface does not know of, and the surface is started with no address stated by flag and none configured
- **THEN** it binds the next port down, prints that address, and the page is reachable there

#### Scenario: The window is walked from the top down rather than at random

- **WHEN** several ports from the top of the window down are held and the surface is started with nothing stated
- **THEN** it binds the highest free port of the window, and a surface started again after it stops binds that same one

#### Scenario: A whole window held is a refusal that names the window

- **WHEN** every port of the window is held and the surface is started with nothing stated
- **THEN** it refuses, names the window rather than one port, and starts no server

#### Scenario: A stated address that is held is refused rather than walked

- **WHEN** a caller states a loopback address whose port is held by something the surface does not know of
- **THEN** it refuses, names that address, starts no server, and binds no other port

#### Scenario: A configured bind that is held is refused rather than walked

- **WHEN** an operator has stated the address to bind and that port is held by something the surface does not know of
- **THEN** it refuses, names that address, starts no server, and binds no other port

#### Scenario: A stated reached-at address pins the port too

- **WHEN** an operator has stated the address the surface is reached at, has stated no address to bind, and the port at the top of the window is held
- **THEN** it refuses rather than binding the next port down, and starts no server

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

### Requirement: A page says where it can be opened from

Where nothing states the address the surface is reached at, serving a page SHALL say that the page can be opened from the machine serving it and from nowhere else, and SHALL name the field that changes the address it prints.

It SHALL say it on the stream a person reads, and the stream a caller parses SHALL carry the address and nothing else.
A caller reads that stream to learn where the page is, and a second line on it is a caller's parse breaking.

Where an operator has stated the address the surface is reached at, nothing further SHALL be said about where the page can be opened from.
That statement is the answer, and asking again on every serve is noise with no end.

What is stated above is the only condition: nothing else, and no mark of the machine it happens to be running on, decides whether it is said.
It states where the page is served rather than judging whether the person can reach it.
Whether the browser that will open the address is on this machine is a harder form of a question the surface already declines about its own reached-at address.
A guess at it is wrong in both directions: it warns a person whose browser is right there, or says nothing to the person it exists for.

The surface binds a loopback address and nothing widens it, so a person whose browser is elsewhere gets a connection that opens and carries nothing.
A browser reports that as an empty reply, which names no cause and looks exactly like a port nobody forwarded.

#### Scenario: A page served with nothing stated says where it can be opened from

- **WHEN** a page is served and nothing states the address the surface is reached at
- **THEN** what a person reads says the page can be opened from the machine serving it, names the field that changes the address, and the page is served

#### Scenario: What a caller reads is the address and nothing else

- **WHEN** a page is served and nothing states the address the surface is reached at
- **THEN** the stream a caller parses holds the address and nothing else

#### Scenario: A stated address is an operator who has already answered

- **WHEN** a page is served and an operator has stated the address the surface is reached at
- **THEN** nothing is said about where the page can be opened from

#### Scenario: The port an operator is told to forward is the one the surface binds

- **WHEN** the documented install step names a port to forward
- **THEN** it is the fixed default the surface binds, read from the surface rather than written a second time

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

### Requirement: Session state lives outside the repository being discussed

Everything the surface reads and everything it writes SHALL live under socrates' home directory: the session state, the served document and what an operator states about the surface.
The surface SHALL NOT resolve a home of its own, SHALL NOT read an environment variable of its own, and SHALL NOT read a configuration file of its own.
It SHALL ask for that home through the one package that resolves it, so the rule is stated once for both binaries.

The surface's own state SHALL sit in a directory of its own under that home, beside what socrates writes rather than mixed into it.

Nothing SHALL be written inside the caller's working directory or the repository it is discussing.

agora serves exactly one consumer, and a dependency with one consumer is that consumer's.
Two homes, two variables and two configuration files made "where is your state?" a question with two halves, which is the cost the one-home rule exists to remove and which a second home put straight back.

The variable that once moved the surface's home SHALL be refused rather than ignored where it is still set, naming itself and naming the variable that replaced it.
A stated home that silently does nothing is a person who believes their state moved and is looking in the wrong place.

socrates already holds that nothing writes in the repository under review, and a surface that dropped a directory there would break it from the outside.

#### Scenario: The repository is untouched

- **WHEN** a surface is served, answered and ended from within a git repository
- **THEN** that repository's `git status` is byte-identical to what it was before

#### Scenario: One variable moves everything both binaries have

- **WHEN** socrates' home is stated, a review is run and a surface is served
- **THEN** the runs socrates wrote, the state the surface wrote and the one configuration file they read are all under the stated home, and nothing is written anywhere else

#### Scenario: A relative home is refused

- **WHEN** socrates' home is stated as a relative path
- **THEN** the command refuses, names the variable and the path, and does not fall back to the default

#### Scenario: The surface keeps a directory of its own inside that home

- **WHEN** a review is run and a surface is served from within it
- **THEN** the surface's state is under a directory of its own beside the runs, and no run directory holds any of it

#### Scenario: The variable that moved the surface's own home is refused

- **WHEN** the variable that once moved the surface's home is set and a surface is served
- **THEN** it refuses, names that variable and the one that replaced it, and starts no server

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
