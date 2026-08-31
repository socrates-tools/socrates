## Context

See `proposal.md` for why.

`internal/process` splits on `//go:build unix` and `//go:build !unix`.
The unix file sets `Setpgid` on spawn and kills the negated pid to reap the group.
The other file has `OwnGroup` and `Detach` as no-ops, `KillGroup` returning a stated error, and `Kill` and `Alive` falling back to `os.FindProcess`.

`internal/snapshot/seal.go` clears `0o222` across the worktree and `ProveTheSeal` opens a probe file with `O_CREATE|O_EXCL`.
A permission error is the seal holding; a successful create is the seal not holding, and that is a refusal.

`internal/reviewer` calls `ProveTheSeal` before spawning.
That is the existing place where socrates decides it cannot stand behind a verdict and stops.

## Goals / Non-Goals

**Goals:**

- One place states which platforms are supported, and the refusal, the test and the release matrix all read it.
- A person on an unsupported platform meets a refusal that names the missing guarantee and the way forward, before anything is spawned and before any repository is touched.
- Adding a platform to a release cannot happen without stating that it is supported.

**Non-Goals:**

- Implementing Windows support. Job objects for reaping and ACLs for the seal are a real change with a real design, and this one is about being honest rather than about being universal. It is raised as an issue rather than folded in here.
- Weakening either guarantee so more platforms qualify. The guarantees are why the tool is worth running.
- Deciding anything about the surface. `agora` binds loopback and serves a page; it makes neither of these two promises, and nothing here reaches it.

## Decisions

### The refusal is at the same seam as the seal's

`ProveTheSeal` is already the point where socrates refuses to review because a guarantee cannot be held, and its reasoning is stated in the specification: a verdict it cannot stand behind is worse than no verdict.

The platform check joins it rather than sitting somewhere earlier.
The alternative is refusing at startup, in `cmd/socrates`, which would refuse `socrates version`, `socrates logs` and `socrates skill` too.
Those commands make neither promise: reading a finished run's log needs no process group, and a person on an unsupported platform reading logs written on a supported one is a reasonable thing to be doing.

What is refused is the review.

### The unsupported path stops being ignorable

Today `KillGroup` returns an error, which a caller may log and continue past.
That is the wrong shape for a guarantee: it makes the absence of reaping a runtime condition somebody handles rather than a fact about the build.

The supported set becomes a compile-time-known value that `internal/reviewer` reads before spawning, so the refusal happens once, at a point where nothing has been started yet, rather than at the end of a round when there is a process group to fail to kill and nothing to be done about it.

`KillGroup` keeps its error for the case where the group exists and the kill fails, which is a genuine runtime failure.

### The supported set is stated once and read three times

A single place names the supported platforms.
The refusal reads it to decide whether to refuse.
A test reads it and compares it against what the release builds.
`README.md`'s supported-platforms section is checked against it, in the same way the port the README tells a person to forward is already checked against the surface rather than written a second time.

Writing the set three times is the failure this avoids, and it is the failure that already happened once in this repository, where `README.md` says Go 1.25 and `go.mod` says 1.26.

### Supported means both guarantees hold, tested

The basis stated is not "it compiles".
A platform is supported when reaping by process group is available, the seal is enforced by the filesystem, and the suite runs there.

Linux and macOS satisfy all three.
The release matrix in `a-version-somebody-can-install` builds those and no others, and the test in this change is what makes adding a fourth impossible without saying so.

## Risks / Trade-offs

**A refusal is worse than a working tool for somebody on Windows.**
It is, and it is better than a verdict that quietly means less than it says.
The remedy named in the refusal is WSL, which is a supported Linux and is how the tool is developed today.

**The seal is only probed, never platform-checked.**
`ProveTheSeal` answers the question empirically on whatever filesystem it finds, which is stronger than a platform list and stays as it is.
The platform list covers reaping, which cannot be probed without spawning something.

**Stating a supported set invites requests to extend it.**
That is the correct pressure.
An issue asking for Windows support is a better artifact than a silent degradation, and it is raised deliberately as part of this work.
