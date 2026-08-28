## ADDED Requirements

### Requirement: A page says where it can be opened from

Where nothing states the address the surface is reached at, serving a page SHALL say that the page can be opened from the machine serving it and from nowhere else, and SHALL name the field that changes the address it prints.

It SHALL say it on the stream a person reads, and the stream a caller parses SHALL carry the address and nothing else.
A caller reads that stream to learn where the page is, and a second line on it is a caller's parse breaking.

Where an operator has stated the address the surface is reached at, nothing further SHALL be said about where the page can be opened from.
That statement is the answer, and asking again on every serve is noise with no end.

What is stated above is the only condition: nothing else, and no mark of the machine it happens to be running on, decides whether it is said.
It states where the page is served rather than judging whether the person can reach it, because whether the browser that will open the address is on this machine is a harder form of a question the surface already declines about its own reached-at address, and a guess at it is wrong in both directions: it warns a person whose browser is right there, or says nothing to the person it exists for.

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
