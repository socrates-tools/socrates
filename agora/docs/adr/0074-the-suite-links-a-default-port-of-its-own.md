# 74. The suite links a default port of its own

- Status: accepted
- Date: 2026-10-16
- Part of: issue 76

## Decision

The binaries the suite built bound `127.0.0.1:24672` where nothing stated an address, the same port a real surface stands on.
A contributor using socrates while working on it, which is the ordinary case, holds that port with their own surface, and every test about the default port failed for it.
Those tests serialised on a machine-wide lock and skipped where the port was held, and a skip is a gate reporting a pass for something it never checked.

The default address is now a variable in `agora/surface` the linker can set, and every binary the suite builds is linked with a loopback port of its own in place of it.
Each run of the suite picks that port once, from a range below the ephemeral ranges the supported platforms hand out and apart from the range the suite picks free ports from, so two runs do not share it and no test binding port 0 is handed it.
The tests about the default bind it through the real path with nothing stated, serialise on it within their own run only, and fail rather than skip where something holds it.

The production default is still stated once, in `agora/surface`.
A unit test, compiled with nothing linked, pins it to `127.0.0.1:24672`, and a release stamps nothing outside `internal/version`, so no published binary binds any other port where nothing states one.

This amends ADR 57, which stands in everything else.

## Rejected

- An environment variable naming the default address.
  It would be a third place to state the address beside `surface.bind` and `--address`, which ADR 57 already rejected, and it would reach every build rather than only the suite's.
- Keeping the production port and skipping where it is held.
  A skip reports nothing wrong for a check that never ran, so the suite would pass on the machines of exactly the people using socrates.
- Driving the default-port tests through `surface.bind`.
  A stated bind is a different path from the one taken where nothing is stated, so the tests would stop checking the default at all.
- One fixed port for the suite.
  Two runs of the suite on one machine, in two worktrees, would contend for it as they contended for the production port.
