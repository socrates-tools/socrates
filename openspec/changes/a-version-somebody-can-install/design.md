## Context

See `proposal.md` for why.

`Makefile` composes `STAMP` from `MODULE` and two variables, `commit` and `uncommitted`, set with `-X` into `internal/version`.
`SURFACE_STAMP` does the same into `cmd/agora`'s `main`.
`internal/version` has a boundary test asserting that what a binary says about itself was recorded when it was linked and is never read where it is run.

`make install` runs `go install` for both binaries, links `scrat` beside `socrates`, and runs `socrates skill`.
So the short name is a symlink the Makefile creates, and a person who runs `go install` gets `socrates` and no `scrat`.

`platform-support` states the supported set.
`the-licence-the-work-is-published-under` produced `LICENSE` and `NOTICE`.

## Goals / Non-Goals

**Goals:**

- A version somebody can cite in a bug report, pin in a script, and compare when deciding to upgrade.
- A binary that never claims a version it was not built from.
- A release a stranger can verify without trusting the publisher.
- One archive that gives a person what `make install` gives them.

**Non-Goals:**

- A Homebrew tap. It needs a second repository under the organisation and a formula that is updated on every release, and it is a distribution channel rather than a distribution. Raised as an issue.
- A curl-to-shell installer. It asks a person to pipe a network response into their shell, which is the opposite of what the signing in this change is for.
- Generating the changelog from commits. This repository's commits say what a change is, deliberately and by convention. A changelog says what it means to somebody using the tool, which is a different document with a different audience.
- Publishing the module to a proxy as a library. Everything is `internal/`, so there is nothing importable, and that is by design.

## Decisions

### The version is stamped, never a constant

`internal/version` gains a `version` variable set by `-X` from the tag, beside `commit` and `uncommitted`.

A constant in the source is a second place the version lives, and it is wrong for the whole window between the release and the next bump somebody remembers.
The existing boundary test already says that what a binary says about itself was recorded when it was linked and is never read where it is run; the version is the same kind of fact and gets the same treatment.

An untagged build reports that it is a development build rather than reporting the last tag.
Reporting the last tag is the failure this avoids: it makes a build from an arbitrary commit indistinguishable from the release whose name it borrowed, in exactly the bug report where that distinction matters.

### `scrat` is a real target, not only a symlink

The release archive contains `socrates`, `scrat` and `agora`.

Inside the archive `scrat` can be a copy or a link depending on what the archive format and the extracting platform support, and what matters is that extracting one archive gives a person both names.
`make install`'s symlink stays as it is for the from-source path.

The invariant that `socrates` and `scrat` are one program under two names is unaffected: they are produced from the same build in the same release, so they cannot be at different commits, which is what that invariant exists to guarantee.

### Verification is checksums, a signature, a bill of materials and provenance

Four artifacts, answering four different questions.

The checksum file answers whether the download is intact.
The signature over the checksum file answers whether the checksums are the ones this project published, and it is keyless, so there is no private key to lose or to be asked for.
The bill of materials answers what is inside, which for a binary embedding fonts and a client bundle is not obvious from the outside.
Build provenance answers which workflow, from which commit, produced it, which is the only one of the four that a compromised publisher cannot forge.

Together they let a stranger check the binary against this repository without extending any trust to whoever pushed the tag.

### The release runs the full verification first, and reaches no CLI

The release workflow depends on the verification matrix from `what-continuous-integration-actually-proves`.
A tag that has not passed it produces no artifacts.

`make fixtures` and `make live-smoke` reach the real CLI and stay out of it, as they stay out of continuous integration today.
That boundary is held by build tags rather than by workflow configuration, so the release inherits it, and a task in this change verifies it rather than assuming it.

### `v0.1.0` is the first tag

Pre-1.0 is the honest statement.
The invariants are strong and the command surface is settled, and the surface half is young enough that a breaking change to it is plausible within months.
Tagging 1.0 now would either be a promise the project has to keep or a version number that means nothing, and both are worse than saying where it actually is.

## Risks / Trade-offs

**A release pipeline is machinery that can rot.**
It runs on tags, so it runs rarely, and a workflow that runs rarely breaks quietly.
The mitigation is that its failure is loud and blocking rather than silent: no artifacts appear.

**Cross-compilation is not the same as testing.**
The release builds for the supported set from one runner; the verification matrix is what actually runs the suite on each.
The two read the same stated set, and `platform-support` requires them to agree in both directions, which is what keeps a cross-compiled target from being mistaken for a tested one.

**A hand-written changelog is a discipline, not a mechanism.**
Nothing fails when it is not updated.
Making it fail would mean gating a tag on prose, which trades a real cost for a checkbox.
It is stated as a requirement and left to the person tagging.
