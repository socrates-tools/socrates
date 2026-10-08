# 28. The file holds only what a person set

- Status: accepted
- Date: 2026-09-13
- Part of: socrates-review configure

## Decision

`use the default` removes a setting's line from the configuration file.
It does not write the current default in its place, and where removing a setting empties the block that held it the empty block goes too.
Where it empties the file, an empty file is left rather than a file naming ten defaults.

A file naming a default is pinned to the default of the day it was written.
A later release that raises a ceiling, lengthens a deadline or changes what a reviewer CLI is routed to would reach everybody except the person who once ran the command, and it would reach them least where they had been most deliberate about it.
An absent field is what keeps a person current, so an absent field is what `use the default` produces.

The same rule makes the file readable: after setting one thing it holds one thing, and what a person reads back is what they decided rather than what socrates-review thinks about everything else.
`socrates-review config` is where the whole picture is read, and it reports the source of each value so a default is never mistaken for a decision.

Writing follows from this.
The write is a `yaml.Node` round trip rather than a marshal of the configuration struct, because a marshal would write every field and delete every comment.
Blank lines are carried through the round trip as a marker comment, so a file that nothing changed comes back byte for byte and a file that one key changed differs by that key alone.

## Rejected

- writing the current default in, which makes the file complete and makes every later release invisible to whoever ran the command.
- writing the whole struct on every save, which is one line of code and deletes the notes a person left themselves.
- leaving an emptied block behind, so that `roles:` with nothing under it survives its last setting, which the loader accepts and a person reads as something still being configured.
- creating the file on the first run of `configure` so that there is always something to edit, which puts socrates-review back in the business of writing a file nobody asked for.
