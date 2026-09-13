# 25. socrates-review waits for a person about its own settings

- Status: accepted
- Date: 2026-09-13
- Part of: socrates-review configure

## Decision

`CONTRIBUTING.md` and `AGENTS.md` each said socrates-review does not wait for a person.
`socrates-review configure` draws a menu and blocks on a keystroke, so one of the two had to give.

The sentence was about a review.
A review is started by a session, reports a verdict and parks on a person through a recorded response rather than through a prompt, and that is what stops socrates-review being a thing a pipeline cannot run.
Configuring socrates-review is not a review: nothing about a run is decided there, no reviewer is spawned, no snapshot is taken, and the command exists precisely because a person is at the keyboard.

So the scope statements are amended to say what they meant: socrates-review does not wait for a person about a review.
`configure` is the one command that waits, it waits about socrates-review's own settings, and it refuses rather than waits where there is nobody to wait for.

## Rejected

- leaving the scope statements alone and shipping a command that contradicts them, which makes the document the thing nobody can trust rather than the thing a contributor reads first.
- a non-interactive `configure --set <key> <value>`, so that nothing blocks. Reading settings without a terminal is `socrates-review config` and writing them without a terminal is the editor, and a flag form would be a second way to write the file that has to agree with the first about validation, removal and comments forever.
- making the menu a mode of `socrates-review config`, which puts reading and writing behind one name and makes the safe command the one that can rewrite the file.
