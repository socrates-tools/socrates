# 2: serve finds the surface by asking the port who stands there

## What to build

The surface answers, without a token, who it is: that it is agora, the root it writes under, and the build it was linked from, and nothing a caller could act on.
That answer is held to the same host guard as every other route.

`serve` and `stop` decide where the surface is from what the port answers rather than from `surface.json`, which keeps the control token and no longer says where the surface is.
`serve` reads the token only once the port has said it is this root's surface.

Where the port holds an agora surface of another root, such as another `SOCRATES_HOME`, `serve` refuses, naming that root, the port, `surface.bind` and `--address`.

## Acceptance criteria

- [ ] A request to the surface with no token is told it is agora, its root and its build, and is told no token and no session.
- [ ] That answer is refused to a request naming a host the surface did not bind.
- [ ] A `surface.json` naming an address where nothing answers, beside a surface standing on the port, does not stop `serve` from opening the page on the surface standing there.
- [ ] A `serve` under one `SOCRATES_HOME` while a surface of another `SOCRATES_HOME` holds the port fails, naming the other root, the port, `surface.bind` and `--address`, and the other root's sessions are untouched.
- [ ] `socrates agora stop` stops the surface standing on the port whatever address `surface.json` names.

## User stories

- 5, 8

## Blocked by

- 1-one-default-port.md
