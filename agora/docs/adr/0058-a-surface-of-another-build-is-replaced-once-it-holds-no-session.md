# 58. A surface of another build is replaced once it holds no session

- Status: accepted
- Date: 2026-10-08
- Part of: agora-port

## Decision

Where the port holds this root's surface but it was linked from another build of socrates, `serve` reuses it while it holds a session and replaces it once it holds none.
Replacing it stops the old surface through its control token and binds the port again from this build.

After `make install` or an upgrade, the surface still standing serves the old page and the old routes until it goes idle.
Replacing it at once would take a page away from a person in the middle of answering another design or review, and reusing it forever would leave the new build's page unseen for as long as the old one stays busy.
The next quiet moment is the first point at which the new build can take over without costing anybody an answer.

The surface keeps a session in its holdings until it has been idle for 45 minutes, and stops 45 minutes after the last one is gone, so counting every session in the holdings as held would leave the old build standing for most of an hour after the last page was answered.
A surface holds no session when every one of them has been ended, by `end` or by a person sending and ending, and replacing it takes those ended pages with it.
A `serve` that names a session never replaces the surface holding it.
The check and the stop are one request, so a session opened between them is never lost, and a surface being replaced takes no further session.

What counts as another build is anything the binary says about itself, per the suite's release invariants: the version and the commit stamp, so two development builds from different commits are two builds.

## Rejected

- Always replacing it, which is lavish-axi's shape, and costs any open page its session.
- Always reusing it, as before, so the old build serves until it goes idle.
- Refusing, saying which build holds the port, and leaving the person to run `socrates agora stop`.
