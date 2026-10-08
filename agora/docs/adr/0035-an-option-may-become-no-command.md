# 35. An option may become no command

- Status: accepted
- Date: 2026-09-17
- Part of: socrates-design

## Decision

agora refused an option that carried no command, because what the page shows and what the terminal runs could not otherwise be the same string.
Every option review writes is a `socrates review respond` invocation, so the rule cost review nothing.

A design answer runs nothing.
The session reads what the person chose and writes a spec, and there is no command for an option to become.

So an option may become no command, and a decision whose options become commands offers one for every option of it.
A decision mixing the two is refused.

The guarantee the original rule protected survives intact, because it was never about documents: a person choosing between the options of one decision is choosing between things of the same kind.

## Rejected

- Declaring it per document, so that a document states once whether its options become commands.
  This buys review a little safety it does not need, since review's documents are generated from one place, and it costs a field on every document written for the surface including those that would never carry a command.
- Making the command optional per option with no further rule.
  This drops the guarantee rather than narrowing it, and allows a decision offering one option that runs something beside one that does not.
- Writing a command into a design option that nothing would run.
  The page hands back the option's own string, so a string nobody runs is the page lying about what it composed.
