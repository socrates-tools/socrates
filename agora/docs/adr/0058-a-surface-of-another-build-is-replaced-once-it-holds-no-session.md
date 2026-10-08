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

What counts as another build is anything the binary says about itself, per the suite's release invariants: the version and the commit stamp, so two development builds from different commits are two builds.

## Rejected

- Always replacing it, which is lavish-axi's shape, and costs any open page its session.
- Always reusing it, as before, so the old build serves until it goes idle.
- Refusing, saying which build holds the port, and leaving the person to run `socrates agora stop`.
