## Which proposal does this implement

<!-- The change under openspec/changes/, and the tasks it closes. A change with no proposal behind it will be turned away: see CONTRIBUTING.md. -->

- Change:
- Tasks closed:

Where this is a bug fix that restores stated behaviour and needs no proposal, say which requirement it restores instead.

## What the change is

<!-- What it does, in the terms the specification uses. Not why: that is the proposal's job. -->

## Which scenario holds it

<!-- The scenario each new behaviour is named by, and the test that names it. A task is not done until a test named for the scenario failed for the right reason and then passed. -->

## Which guard was mutated

<!-- A check nobody has broken on purpose is a check nobody knows works. Say what you broke to prove the check catches it, and that you restored it. -->

## Checks

- [ ] `make fmt` and `make verify` are green
- [ ] The scenarios this change claims are covered by tests named for them
- [ ] No comments were added to the source
- [ ] No em dash anywhere, and every user-facing string lives in its package's `messages.go`
- [ ] `openspec/` and the source agree, and the task boxes checked are the ones actually done
