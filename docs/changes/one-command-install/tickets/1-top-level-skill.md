# 1: socrates skill installs every namespace's skills

## What to build

`socrates skill` runs the `skill` operation of every namespace the dispatcher mounts, at its defaults, and stops at the first failure.
It takes no flags and refuses an argument.
`make install` runs it once in place of the three namespace operations.

## Acceptance criteria

- [x] `socrates skill` installs every skill each namespace's own `skill` installs, and the implementer's agent definition, where those operations put them, with what their sources generate.
- [x] What it prints is what each namespace's `skill` prints.
- [x] `socrates skill anything` is refused, naming the argument.
- [x] The dispatcher names no namespace, flag or skill to do it.
- [x] `make install` runs `socrates skill` and no namespace operation.

## User stories

- 14, 15, 18

## Blocked by

None
