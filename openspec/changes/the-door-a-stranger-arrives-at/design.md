## Context

See `proposal.md` for why.

`README.md` is 15 KB and well organised.
It opens with what socrates is and states the requirements at line 18.
It then works through the first run, the reading commands, the exit-code contract, the effect on the repository, configuration, reading the page from another machine, what is out of scope, and the architecture.

`AGENTS.md` is the working document, symlinked as `CLAUDE.md`.
It is a map for an agent: layout, verification sequence, invariants, conventions.
It is not a contributor guide and reading it as one gives a person the invariants without telling them how to propose a change.

`openspec/` holds 27 archived changes and 17 promoted capabilities, and `openspec/config.yaml` carries the normative context.
Nothing outside `openspec/` explains that a change starts as a proposal.

`internal/convention` already holds checks like "every user-facing string lives in its package's `messages.go`" and "the documented port is read from the surface rather than written a second time".

## Goals / Non-Goals

**Goals:**

- A reader learns whether they can use socrates at all before they learn how it works.
- A person who finds a vulnerability has somewhere private to send it.
- A contributor learns how a change is proposed here before they write code that will be turned away.
- The facts a document states about the build are checked rather than maintained.

**Non-Goals:**

- Rewriting the README's body. The parts that explain socrates are good and this change adds to the top and around them.
- Softening the "Out of scope" section. See the decision below.
- A documentation site. It is a channel rather than a document, and it is raised as an issue.
- Turning `AGENTS.md` into `CONTRIBUTING.md`. They have different readers, and the map an agent needs is not the path a person needs.

## Decisions

### The requirements come before the explanation

The README gains a short block, before anything else, naming what has to be true before socrates works.
That is a supported platform, git, Go for a source install, and the `claude` CLI at its floor with the account it needs.

The argument against is that leading with a dependency on a paid third-party product loses readers.
It loses exactly the readers who would otherwise spend ten minutes reading and then discover it, and lose that time to a project that knew and did not say.
A reader who continues past it is a reader who can actually use the tool, and that is a better first hundred users than a larger number who cannot.

### The "Out of scope" section is untouched

That section names the residual risk on a managed machine, states that `clean` does not mean the tests pass, and says the reviewer does not know the project's conventions.

It is the single most trust-building thing in the repository, and the pressure a public release creates is exactly the pressure to soften it.
It is stated here as a non-goal so that a later change proposing to trim it is recognised as trimming it.

### A telemetry statement, because silence is read as unknown

socrates reaches no network outside a build tag, and a person reading the README today cannot learn that.

For a tool that reads the repository somebody is working in and hands parts of it to a model, "we send nothing" is not a boast, it is the fact a careful person will otherwise assume the worst about.
It goes in the README and in `SECURITY.md`, where the reader who cares most is already looking.

What it says has a real boundary: socrates sends nothing, and the reviewer CLI it drives sends the change to its vendor, which is the whole mechanism.
The statement names both rather than only the flattering half.

### The security policy carries a threat model, not just an address

An address alone tells a reporter where to send something and not whether it is worth sending.

The policy states what socrates defends against, which is a reviewer that writes, a session carrying the person's customisations, a review inside a review, and a page reachable from another machine.
It states what it does not, which is a hook an administrator's policy configures and whatever that hook's command reaches outside the sealed worktree.
The README already says that last one, which means the honest version of the policy is mostly a pointer to prose that exists.

### `CONTRIBUTING.md` leads with the proposal, not the toolchain

The first thing it says is that a change starts as a proposal in `openspec/`.

A contributor who learns the toolchain first writes code first, and the pull request that arrives is one that has to be turned away for a reason nothing told them.
The toolchain, the one verification command, and the conventions follow, and they are short because `make verify` and `make help` do the work a long list would otherwise do.

### The Go floor is checked, not maintained

`internal/convention` gains a check that the floor the README states and the floor `go.mod` requires agree.

This is the second instance of a shape this repository already uses: the documented port is read from the surface rather than written twice.
The existing disagreement between README and `go.mod` is what the check exists to prevent recurring, and fixing it without the check would leave the same drift free to happen again.

## Risks / Trade-offs

**More documents is more to keep current.**
The mitigation is that the facts most likely to drift are checked: the Go floor, the supported platforms, and the documented port.
What is left is prose about intent, which drifts slowly.

**A code of conduct with a contact address is a commitment.**
It commits the maintainer to reading and acting on reports.
It is written anyway, because a project that invites contribution and has no stated conduct is making an implicit commitment with none of the clarity.

**A stated security policy attracts low-quality reports.**
It also attracts the one that matters.
The scope section is what makes triage cheap, and it is written for that purpose rather than as decoration.
