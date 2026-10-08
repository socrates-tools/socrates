# 1: The surface binds one port, and a held port is refused

## What to build

With nothing stated, `socrates agora serve` binds `127.0.0.1:24672` and nothing else.
Where something that is not this root's surface holds that port, `serve` refuses, naming the port, `surface.bind` in agora's file and `--address`, and starts no surface anywhere else.
The window of ten ports and the walk down it are gone.
`surface.bind` and `--address` move the surface exactly as before, and a port either commits to is bound or refused.

What a person reads names the one port: the README's configuration example and its section on reading the page from another machine, agora's own instructions, and the changelog's Unreleased section, which tells somebody who has the previous release that the page moved to 24672, that a held port is now refused rather than walked past, and that whoever forwarded 18018 forwards 24672 instead.

The live smoke picks a free loopback port rather than one from the window, and no free-port pick the suite makes ever lands on 24672.

## Acceptance criteria

- [ ] With nothing stated and 24672 free, the page `serve` prints is on `127.0.0.1:24672`.
- [ ] With 24672 held by another program and nothing stated, `serve` fails, naming 24672, `surface.bind` and `--address`, and nothing listens on 18018 or any other port afterwards.
- [ ] A held port stated by `--address` or by `surface.bind` is still refused, and a free one is still bound.
- [ ] The README names 24672 as the port to forward, says nothing of a window, and the suite fails if the port it names differs from the one the surface binds.
- [ ] agora's instructions state one fixed port refused when held, and no window.
- [ ] The changelog's Unreleased section tells a person who forwarded 18018 to forward 24672.
- [ ] No address in the source or the documents names the old window, except a `surface.bind` value a test states as a person's own.

## User stories

- 1, 2, 3, 4, 9, 10

## Blocked by

None
