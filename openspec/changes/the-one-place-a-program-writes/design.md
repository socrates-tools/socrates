## Context

See proposal.md for why.

What exists today, in `internal/state`: `SOCRATES_STATE_HOME` if non-empty, else `XDG_STATE_HOME` joined with `socrates` if non-empty, else `~/.local/state/socrates`.
Runs are `runs/` under it.

In `internal/config`, and again in `internal/surface` for agora: `XDG_CONFIG_HOME` if `filepath.IsAbs`, else `os.UserConfigDir()`, else `~/.config`.
The file is `<directory>/<program>/config.yaml`.

Go's standard library offers `os.UserConfigDir()` and `os.UserCacheDir()` and no equivalent for state, which is why one of these resolves through the standard library and the other does not.
`internal/snapshot` names `SOCRATES_STATE_HOME` in the refusal a filesystem that does not enforce permissions raises.

## Goals / Non-Goals

**Goals:**

- One directory per program, and one variable that moves it.
- The same resolution rule for everything a program reads and everything it writes.
- A stated home that does not take effect is a refusal rather than a silence.

**Non-Goals:**

- Moving anything a person already has. There is no migration.
- A shared home for the two programs. They are separate, and the boundary is an invariant.
- Anything about what is written, or the shape of a run directory, or the shape of the state a surface keeps.

## Decisions

### A dotted directory under the person's home, and the XDG variables go

`~/socrates` is not available: it is exactly where somebody clones the repository, and a tool that claims that path cannot coexist with its own source tree.
That matters more the week the repository is public than it did before.
So `~/.socrates`, and `~/.agora`.

The XDG variables are dropped rather than honoured where present.
Honouring them makes "everything is in one place" conditional, which puts the second path back into every answer to the question this change exists to make answerable in one line.
`SOCRATES_HOME` is strictly more control than the two variables it replaces, because it moves the configuration as well as the state, and it is one thing to say rather than two things and a precedence rule.

The precedent is immediate: `~/.claude` holds settings, projects and state in one dotted directory for the very tool socrates drives, and nobody has been confused by it.

Considered and rejected: keeping the split and reconciling the two readings of the XDG variables.
It fixes the drift and leaves the cost that produced it, which is two rules where one will do.

### A relative stated home is refused, not ignored

The specification these variables came from says a relative value must be ignored, and `internal/config` implements that while `internal/state` does not.
Neither reading survives the variables going away, so the rule is chosen fresh rather than inherited.

Refusing is the better of the two.
A variable that silently does nothing sends a person to look for their runs in a directory that has none, and the failure surfaces as "socrates lost my review" rather than as "that path is relative".
Ignoring made sense for a specification describing a shared convention across many programs, where refusing would mean one misconfigured variable breaking everything on the system.
That argument does not carry to a variable named for one program.

### The configuration file sits at the top of the home, not in a subdirectory

`~/.socrates/config.yaml` beside `~/.socrates/runs/`.
The one file a person writes is the first thing they see in the directory, and burying it a level down to keep symmetry with the runs would be symmetry bought with the property this change is for.

### Every path stays behind the package that owns it

`internal/state` remains the one resolver of socrates' home and `internal/surface` of agora's.
`internal/config` stops resolving a directory at all and asks `internal/state` where the home is, which removes the second implementation of a resolution rule rather than editing it to match the first.

That is the whole reason the drift existed: two functions, in two packages, that both meant "where does this program keep things".

### The refusal that names a variable names the new one

`internal/snapshot` tells a person to set `SOCRATES_STATE_HOME` when their filesystem does not enforce permissions.
It becomes `SOCRATES_HOME`, and the message stays in its package's `messages.go` as every user-facing string does.

## Risks / Trade-offs

- **The configuration leaves `~/.config`, which is what dotfile managers target** → a person who version-controls their configuration now has one file to pick out of a directory holding run state and sealed worktrees, rather than a file in a tree their tool already manages. A symlink handles it, and dotfile managers handle symlinked single files routinely. This is the one real cost of the collapse, and it is paid by the people most likely to have an opinion about it.
- **A person with existing runs loses sight of them** → nothing migrates, and `~/.local/state/socrates` is left where it is rather than deleted. Runs are swept at fourteen days by design, so what is lost is at most a fortnight of finished reviews, and moving the directory by hand is the whole of the recovery.
- **`~` gains a dotted directory where the convention says it should not** → true, and it is the convention that says so rather than an outcome anybody suffers. The convention exists to stop tools each inventing a visible `~/something`; a hidden directory per program is what almost every tool a person already has does.
- **Dropping the XDG variables ignores people who set them deliberately** → they set one variable instead of two and get the configuration moved as well, which is what they wanted from both. Nobody who cared is worse off; the rule they knew is simply no longer the one to know.

## Migration Plan

None, deliberately.

Nothing is moved and nothing is read from the old locations, so a person with existing state finds socrates behaving as though it had never run.
`~/.local/state/socrates` and `~/.config/socrates` are left in place rather than removed, because a tool that deletes a directory it has stopped using is a tool deleting a directory it no longer understands.
