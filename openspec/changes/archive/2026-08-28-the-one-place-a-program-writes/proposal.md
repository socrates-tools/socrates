## Why

socrates keeps what a person writes in one place and what it writes in another.
`config.yaml` is under the configuration directory, runs are under the state directory, and the two are resolved by different rules from three environment variables.
agora does the same again with its own two.

Nothing is wrong with the split in principle, and it buys nothing here.
Nobody is going to put runs on a separate volume, and the config is one small optional file.
What it costs is that the answer to "where does socrates keep things" is two paths and a resolution rule, and that answer is about to be given to strangers: the repository is going open source, and "where is your state" is the first question on every issue that is not reproducible.

The two rules have already drifted apart.
The configuration honours its variable only when it is absolute, matching the specification it follows; the state directory honours anything non-empty.
Same specification, two readings of it, and no test names the difference because no requirement states which reading is right.

One directory per program, one variable, and both questions have one answer.

## What Changes

- **BREAKING**: socrates keeps everything it reads and writes under one home directory, `~/.socrates` by default, holding `config.yaml` and `runs/`.
- **BREAKING**: agora does the same under `~/.agora`, holding what the surface writes.
- **BREAKING**: `SOCRATES_HOME` and `AGORA_HOME` replace `SOCRATES_STATE_HOME`, `AGORA_STATE_HOME`, `XDG_STATE_HOME` and `XDG_CONFIG_HOME`.
  One variable moves a program's whole home, which is strictly more control than the two it replaces and one fewer rule to state.
- A home stated by variable is required to be absolute, and a relative one is refused rather than ignored.
  Where two rules disagreed there is now one, and it refuses rather than guessing.
- The refusal a sealed worktree raises names the variable that moves the home, which is the variable that changed.
- **The repository stops naming one machine.** `docs/gotchas.md` names a path that exists on the author's computer and nowhere else; it says what kind of directory is wanted instead.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `cli-foundation`: `Configuration is one optional user file` states the one home the file lives in and the one variable that moves it, and gains that a stated home must be absolute.
- `run-store`: `A run is a directory whose name answers the common question` names socrates' home rather than a state directory distinct from where the configuration lives.
- `review-surface`: `Session state lives outside the repository being discussed` names agora's own home, resolved by the same rule under its own variable.

## Impact

- `internal/state`: resolution collapses to `SOCRATES_HOME` and `~/.socrates`, and a relative stated home is refused.
- `internal/config`: the configuration file moves under the same home, so `internal/config` stops resolving a directory of its own.
- `internal/surface`: the same collapse for `AGORA_HOME` and `~/.agora`, for both the session state and the configuration file.
- `internal/snapshot`: the seal refusal names `SOCRATES_HOME`.
- `internal/testing`, `internal/e2e`: one variable to set per program instead of two, so harness setup shortens.
- `README.md`, `docs/gotchas.md`: the paths a person is told, and the one path that names a machine.
- No new dependency, and nothing about what is written changes.
