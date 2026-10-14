# 61. The top level installs every namespace's skills

- Status: accepted
- Date: 2026-10-14
- Part of: one-command-install
- Amends: ADR 33

## Decision

`socrates skill` installs every skill the binary ships, by running each mounted namespace's own `skill` operation with that operation's defaults.
It takes no flags and no arguments.
A person who wants a skill somewhere other than its default runs the namespace's own `skill` with `--to`, as before.

ADR 33 kept the top level to routing, `version` and `help`, and put `skill` inside each namespace because a skill belongs to a tool rather than to the suite.
That reasoning still holds and this does not undo it: every skill is still written by its own tool's operation, from its own tool's source, to its own tool's default.
What the top level adds is the fan-out, and the fan-out is routing: the dispatcher finds a `skill` operation in each tree it mounts and runs it, and names no tool, no flag and no skill.
A namespace that gains a `skill` operation is installed by `socrates skill` without the dispatcher changing.

The person asked for it plainly: installing socrates should be one command, and that command should leave every skill in place.
Three invocations that a person has to know about, one per namespace, are three chances to install two of them.

## Rejected

- Leaving the top level as ADR 33 had it, and having the install script and `make install` each run the three namespace operations.
  That is two copies of the namespace list outside the tree, and a fourth tool with a skill would be installed by neither until somebody remembered both.
- Passing `--to` through from the top level to every namespace.
  `--to` is each namespace's own flag, the dispatcher names none, and `socrates implement skill` also writes an agent definition that `--to` does not move, so one directory for everything is a promise the top level cannot keep.
- A `socrates install` that also places the binary.
  The binary cannot install itself before it is on the machine, so placing it is the install script's work and the skills are the binary's.
