## Which spec does this implement

<!-- The spec under docs/changes/ this implements, and what it asked for. A change with no spec behind it will be turned away: see CONTRIBUTING.md. -->

- Spec:
- What it asked for:

Where this is a bug fix that restores stated behaviour and needs no spec, say which behaviour it restores instead.

## What the change is

<!-- What it does, in the vocabulary the tool's CONTEXT.md defines. Not why: that is the spec's job. -->

## Which scenario holds it

<!-- The scenario each new behaviour is named by, and the test that names it. A task is not done until a test named for the scenario failed for the right reason and then passed. -->

## Which decision it rests on or contradicts

<!-- The ADR under a docs/adr/ directory this change assumes, or the one it contradicts. Contradicting one is allowed and saying nothing about it is not. -->

## Which guard was mutated

<!-- A check nobody has broken on purpose is a check nobody knows works. Say what you broke to prove the check catches it, and that you restored it. -->

## Checks

- [ ] `make fmt` and `make verify` are green
- [ ] The scenarios this change claims are covered by tests named for them
- [ ] No comments were added to the source
- [ ] No em dash anywhere, and every user-facing string lives in its package's `messages.go`
- [ ] The spec and the source agree, and what this says it does is what it does
