**The suite must be green after every task, not only at the end of a group.**
Two orderings below exist for that reason: the tests state the address they bind before the default moves, and `internal/surface` holds sessions before anything outside it knows they exist.

**A session is reached through the real binary.**
No task asserts on served markup, and the page half is driven by a real browser, as the surface's tests already are.

## 1. The suite stops depending on a port nobody chose

- [ ] 1.1 State `--address 127.0.0.1:0` explicitly wherever a test starts a surface and does not already, reading the bound address back from the surface as those tests already do, and verify `make test` passes with the default unchanged
- [ ] 1.2 Add a test that a surface started with no address stated by flag and none configured binds the fixed default port, failing for the right reason while the default is still `127.0.0.1:0`. Covers `The default bind is a fixed loopback port`
- [ ] 1.3 Move `DefaultAddress` to `127.0.0.1:18018` and verify 1.2 passes and `make test` is green

## 2. internal/surface holds sessions

- [ ] 2.1 Mint a session identifier from eight bytes of `crypto/rand` rendered hex, and verify a unit test that two mints differ and each is a legal path segment
- [ ] 2.2 Give the server a registry of exchanges keyed by session, each with its own document, its own page token and its own last-touch, with `Serve` still opening exactly one so the existing tests keep passing unchanged
- [ ] 2.3 Add opening a further session on a running server, and verify a test that two sessions hold different documents and different page tokens. Covers `Two sessions are served at once and neither disturbs the other`
- [ ] 2.4 Refuse a named session the server does not hold rather than creating one, and verify the refusal names the session and opens nothing. Covers `A session the surface never minted is refused`
- [ ] 2.5 Collect a session whose exchange has had no reader for the idle period, leaving the others, and verify with `testing/synctest`. Covers `An idle session is collected and the ones being read are not`
- [ ] 2.6 Stop the server when it holds no sessions and has held none for the idle period, and verify a server that is never given a session stops too. Covers `An idle server stops`

## 3. The routes carry the session

- [ ] 3.1 Serve the page and its routes under `/sessions/{session}`, leaving `/assets/` shared, and verify a request to a session's page returns that session's document and a request to an unknown session returns not found
- [ ] 3.2 Mint the control token when the server starts, and accept it on every route
- [ ] 3.3 Accept a session's own page token only on routes under that session's path, and verify one session's token is rejected by another's. Covers `One session's token does not answer for another`
- [ ] 3.4 Accept the control token alone on opening a session and on stopping the server, and verify a page token is rejected by both
- [ ] 3.5 Confirm `hostIsTheSurface` still turns away a request naming a host the surface did not bind, now that every session shares one port. Covers `A request claiming another host is rejected`

## 4. The page knows what it is served under

- [ ] 4.1 Carry the base a page's requests are made against in `render.Frame`, emit it as a `data-` attribute beside the token, and verify the rendered page carries it
- [ ] 4.2 Build the client's three requests from that base rather than writing them absolute, regenerate the committed bundle with `make assets`, and verify `make test` still fails nothing on a generated file disagreeing with its source
- [ ] 4.3 Drive two pages in a real browser against one server, answer on the second, and verify the first is untouched and still answerable. Covers `Two sessions are served at once and neither disturbs the other`

## 5. cmd/agora

- [ ] 5.1 Write the surface record holding the address, the process and the control token, and read it where `session.json` was read
- [ ] 5.2 Bring the detached child up holding no sessions, and open the first session over HTTP as every later `serve` does, and verify the first `serve` after a reboot returns an address that holds a page
- [ ] 5.3 Print the address of the session's page, composed from the address the surface is reached at where an operator stated one and from the address bound otherwise. Covers `A session is minted and the address holds it` and `serve prints the address stated over the address bound`
- [ ] 5.4 Take `--session` on `serve`, `poll` and `end`, and verify `poll` given one of two open sessions reports that session's answer and nothing from the other. Covers `A session is named to poll it`
- [ ] 5.5 Write a new document to the named session only, and verify the other page is unchanged. Covers `A new document reaches the session it named and no other`
- [ ] 5.6 Keep `stop` naming no session, and verify it takes down every page while `end` takes down one. Covers `Ending one session leaves the others being read` and `Stopping is the server, not a session`
- [ ] 5.7 Refuse a `serve` whose port is held by something that is not a surface, naming the port and `--address`, and verify against a listener the test holds open
- [ ] 5.8 Keep the command surface at five, and verify the test that examines it still passes. Covers `The surface offers nothing beyond the five`

## 6. The state directory

- [ ] 6.1 Write each session's record and served document under a directory of its own, remove that directory when the session is collected, and verify a collected session leaves nothing behind
- [ ] 6.2 Confirm a repository a surface was started inside is byte-identical after two sessions are served, answered and ended. Covers `The repository is untouched`

## 7. The skill

- [ ] 7.1 Say once, where the loop is described, that the address `agora serve` prints holds the session and that `agora poll` is given that session, and verify neither branch restates it. Covers `The skill says to carry the session from serve to poll` and `The session is stated once rather than on every branch`
- [ ] 7.2 Regenerate the committed skill with `make skill`, and verify the check that every command and flag the skill names exists still passes with `--session` named
- [ ] 7.3 Drive the loop the skill describes end to end through the real binaries, serving a page, answering it and recording the answer, and verify it reaches a verdict

## 8. Verification

- [ ] 8.1 Mutate the guard that a session's token is rejected by another session and confirm the suite fails, then restore it
- [ ] 8.2 Mutate the refusal of an unknown session and confirm the suite fails, then restore it
- [ ] 8.3 Run `make fmt`, `make lint`, `make test`, `make scenario-coverage`, `make scenario-drop` and `make corpus-outcomes`, and verify every scenario this change states is named by a test
- [ ] 8.4 Record one live run with the real CLI through `make live-smoke`, reaching a page served on the fixed port, and record that it happened
