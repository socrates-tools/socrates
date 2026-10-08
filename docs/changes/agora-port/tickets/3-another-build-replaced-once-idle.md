# 3: A surface of another build is replaced once it holds no session

## What to build

Where the port holds this root's surface but it was linked from another build of socrates, `serve` replaces it once it holds no session and reuses it while it holds one.
Another build is any difference in what the binary says about itself: the version, the commit stamp, or whether it was built uncommitted.

Replacing stops the old surface through its control token, waits a bounded time for the port to free, and starts this build's surface on it.
A port that never frees is refused as a port held by something else.
Whether a surface holds a session is never told to a caller without a token.

## Acceptance criteria

- [ ] With a surface of another build standing on the port and holding no session, `serve` hands back a page drawn by the build that ran it, and the old surface's process has exited.
- [ ] With a surface of another build standing on the port and holding a session, `serve` opens its page on that surface, and the session already there still answers.
- [ ] A surface of the same build is reused, as before.
- [ ] A surface of another build that does not let go of the port within the bound is a refusal naming the port, `surface.bind` and `--address`.
- [ ] No request without a token learns whether the surface holds a session.

## User stories

- 6, 7

## Blocked by

- 2-surface-found-by-asking-the-port.md
