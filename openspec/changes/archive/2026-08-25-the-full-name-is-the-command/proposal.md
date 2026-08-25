## Why

The product is called socrates and the command is called `scrat`, and there is no reason for a person to have to know that.
The short name is worth keeping, because it is four characters in a loop somebody types all day, but it is the abbreviation rather than the name, and a tool whose only name is its abbreviation makes every document choose between the two.

Nothing about the review changes here.
This is what the tool is called.

## What Changes

- The binary is `socrates`, and `scrat` is installed beside it as the short name for the same binary.
  Both are commands a person can type, neither is a shell alias somebody has to write themselves, and there is one binary underneath so the two can never answer differently.
- Everything socrates prints names `socrates`: the failure prefix, every refusal, every next step it suggests, the help, and the line `scrat --help` prints while running under the short name.
  One name in the text, because a tool that reports itself by whichever name it was invoked as has no name.
- The skill is `socrates`, installed at `.claude/skills/socrates/`, invoked as `/socrates`, and its steps say `socrates review`, `socrates status` and the rest.
  The constant that names the skill is the same one the surface check parses invocations with, so the skill cannot name the command one thing while the check looks for another.
- `internal/skill`, `internal/report`, `internal/presentation` and the command layer all take the name from the same place they already took `scrat`.

Existing specifications that name `scrat <command>` are left alone, and that is deliberate rather than an omission.
`scrat` remains a command name, so those sentences stay true, and rewriting them would mean rewriting the scenario headings that tests are matched against and declaring a rename for each while two other changes are open against the same capabilities.

## Capabilities

### Modified Capabilities
- `cli-foundation`: the binary is `socrates`, installed under both names, and what it prints names the full one.
- `agent-skill`: the skill is named for the product, and one constant names it for the installer and for the check that reads its invocations.

## Impact

- `cmd/scrat` becomes `cmd/socrates`, which is what makes `go install` produce a binary of that name.
- `Makefile`: `install` puts the short name beside the full one, in whichever directory the toolchain installs to, rather than a directory guessed here.
- `internal/skill`: one constant, and the text it renders.
- `.claude/skills/scrat/` is replaced by `.claude/skills/socrates/`, regenerated rather than moved.
- The suite: what the harness builds and what the scripts invoke, and one journey that installs both names and reads the same answer from each.
- Nothing about a run, a round, a finding, a prompt or a verdict.
