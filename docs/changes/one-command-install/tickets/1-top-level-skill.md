# 1: socrates skill installs every namespace's skills

## What to build

`socrates skill` runs the `skill` operation of every namespace the dispatcher mounts, at its defaults, in mount order, and stops at the first failure.
It takes no flags and refuses an argument.
`make install` runs it once in place of the three namespace operations.

## Acceptance criteria

- [ ] `socrates skill` installs every skill each namespace's own `skill` installs, and the implementer's agent definition, where those operations put them, with what their sources generate.
- [ ] What it prints is what each namespace's `skill` prints, in mount order.
- [ ] `socrates skill anything` is refused, naming the argument.
- [ ] The dispatcher names no namespace, flag or skill to do it.
- [ ] `make install` runs `socrates skill` and no namespace operation.

## User stories

- 14, 15, 18

## Blocked by

None
