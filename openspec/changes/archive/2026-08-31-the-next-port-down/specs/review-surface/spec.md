## MODIFIED Requirements

### Requirement: The surface binds a loopback address and nothing else

The server SHALL bind a loopback address.
No flag, no environment variable, no document field and no configured address SHALL widen that bind, and there SHALL be no tunnel, no share command and no third-party host.

Where nothing states an address, the surface SHALL bind the first free port of a fixed window of loopback ports, taken from the top of the window down, rather than a port the operating system chooses.
The window SHALL be closed and stated in one place, and the port at its top SHALL be what a machine holding none of them binds.
Where every port in the window is held, the surface SHALL refuse, name the window, and start no server.

An ordered closed window is not a port that moves.
What made a port the operating system chooses unusable is that nothing can be aimed at it and nobody can return to it, and neither is true here: what an operator aims something at is a port they stated, which pins the bind rather than being walked past, and what a person returns to is the address the surface printed, which is the one it bound and never the one it was asked for.
What it buys is the case the single fixed port refused on, which is a person who has never named a port meeting a failure about one, whose only way out was to learn a flag to get back to where they started.

An operator MAY state the loopback address the surface binds, and that address SHALL be refused by the same check a stated flag is refused by.
A flag SHALL override it, so what an operator configured is a default rather than a decision taken away from the caller.

Where an operator has stated the address to bind, or the address the surface is reached at, the window SHALL NOT be walked: the port that follows from what they stated is the port bound, and a port held by something the surface does not know of SHALL be refused rather than stepped over.
Stating either is an operator committing to a port, because something in front of the surface is aimed at one.
A surface that moved out from under it would hand a person the address the operator stated, which by then opens onto nothing, and an address that opens onto nothing is the failure serving a page says where it can be opened from in order to prevent.
Which port is taken where nobody has named one is the only thing the address the surface is reached at reaches: it still reaches neither the check on what may be bound nor the guard on what host a request may name.

An operator MAY state the address the surface is reached at.
It is handed to a person to open, so it SHALL be refused unless it is shaped like something openable, and it SHALL be refused for its shape alone: whether anything serves it is not something the surface can answer from where it runs.
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
