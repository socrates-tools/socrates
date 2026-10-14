# 33. The command is socrates namespace operation

- Status: accepted
- Date: 2026-09-16
- Part of: one command
- Amended by: ADR 61, which adds a top-level `skill` that runs every namespace's own

## Decision

`socrates <namespace> <operation>` is the command surface for the whole suite.
`review` and `agora` are the namespaces, each naming a workflow or a component; `begin`, `serve`, `status` and the rest are operations inside one of them.
The top level holds routing and what is true of the whole suite, which is `version` and `help` and nothing else: `config`, `configure` and `skill` sit inside their namespace, because configuration and a skill both belong to a tool rather than to the suite.

A reader who finds the archived rename rejecting `socrates` as a binary name needs to know why a dispatcher is not what was rejected.
That rename considered leaving the binary named `socrates` while it was the reviewer, and rejected it on the second tool: two binaries published by the organisation, one of them called `socrates`, leaves a person unable to say which one they mean.
That reasoning is correct and this does not contradict it.
Under a dispatcher, `socrates` is not one of the tools.
It is the brand, it routes, and every tool is named by the word after it, so there is no pair to be confused between.

The trade-off recorded is a word added to every invocation: `socrates-review begin` becomes `socrates review begin`, and `agora serve` becomes `socrates agora serve`.
Nothing about what a review does, what a page renders, or what either exits with changes; a script that branched on an exit code keeps branching on the same code, and a session that read `socrates-review findings --json` reads `socrates review findings --json` instead.

`internal/dispatcher` is where this is built: it mounts `review/command`'s tree and `agora/command`'s tree as subcommands of its own root, and holds wiring and nothing else.
A namespace typed bare prints its own help and starts nothing, because its root command is the same command a person could execute on its own, just reached one word later.
The top level offers no workflow's own verb: `socrates begin` and `socrates serve` are refused as unknown commands, naming `socrates` rather than a namespace, because beginning a review and serving a page are `review`'s and `agora`'s to own.

## Rejected

- Keeping `socrates-review` and `agora` as the two published names, which is the shape the archived rename settled on and the shape a second workflow after review has nowhere better to go from: a third binary is either `socrates-grill` restating the brand for a third time, or a bare word with no claim on it at all.
- A suite-level `configure` reading every tool's settings from one menu, which is the suite-wide configuration file already rejected once, now reached by a different door.
- Naming the workflow verb `review` at the top level (`socrates review` as an operation rather than `socrates review begin` as a namespace and an operation), which was ruled out before this decision by keeping `begin` as the verb: the workflow is named once, by the namespace, and the verb says what happens inside it.
