# socrates

An independent reviewer for the changes you just wrote and nobody has read.

Your session implements something.
You now have unreviewed local changes.
You run `scrat review`.
socrates freezes the working tree into a commit nothing can move, spawns a detached reviewer against that snapshot, and returns.
The reviewer never authored anything, is told nothing about how the change came to be, and answers with structured findings.
You read them and decide.

The asymmetry is the whole point.
The author has the context and does the patching.
The reviewer has independence and has written none of it.

## The first run

socrates needs Go 1.25 or later, `git`, and the `claude` CLI on your `PATH`.

```sh
make build                     # bin/scrat
bin/scrat review --intent "add the parser I asked for"
```

`review` prints the run identifier and returns while the review runs:

```
20260819T101500-ab12cd34
follow it with: scrat status 20260819T101500-ab12cd34
```

Then either poll it or wait for it:

```sh
scrat status                   # a few lines, and an exit code you can branch on
scrat findings                 # what it found, as socrates classified it
scrat findings --json          # the same, with every effective action resolved
scrat review --wait            # block until the verdict, and exit with its code
scrat cancel                   # stop a review and say that nothing was decided
```

Naming no run means the most recent run for the repository you are standing in.

Exit codes are a contract:

| code | meaning |
| --- | --- |
| `0` | `clean`, or a command whose job was to start something that started |
| `1` | `escalated`: a reviewer found something and it is yours to decide |
| `2` | `failed`: no review completed, so nothing has been checked |
| `3` | a usage or precondition error |
| `4` | the run is not terminal yet |

`failed` is deliberately not `escalated`.
Escalated means read the findings.
Failed means socrates could not do its job and nothing was reviewed.

## What it does to your repository

Nothing, other than adding the git objects the snapshot needs.

The snapshot is built through a temporary index, so your own index is untouched.
It includes untracked files, because a freshly implemented feature is mostly untracked files, and it honours `.gitignore`.
No ref names the snapshot commit, so ordinary garbage collection reclaims it once the worktree is gone.
`git status` before and after a review is byte-identical, and a test holds that.

Run state lives under `${XDG_STATE_HOME:-~/.local/state}/socrates/runs/`, one directory per run, legible with `ls` and `cat`.

## Configuration

One optional file, at `${XDG_CONFIG_HOME:-~/.config}/socrates/config.yaml`.
socrates never creates it, never reads configuration out of the repository under review, and refuses rather than quietly running on defaults when the file will not parse.

```yaml
roles:
  reviewer:
    model: opus
    effort: high
diff_ceiling: 131072
agent_deadline: 30m
log_ceiling: 1048576
run_retention: 336h
```

## What it deliberately is not

- **It does not land your change.** No commit, no push, no pull request, no continuous integration. Those are yours, and a reviewer that could land things is not a reviewer.
- **It does not verify your change.** It runs no tests and no linters. `clean` means an independent reviewer found nothing blocking, not that the suite is green.
- **It does not write in your repository.** The reviewer's tools cannot write, and the round fails if the snapshot worktree differs by one byte afterwards.
- **It does not load your tool servers.** A tool server is a second source of tools the vendor's own tool restriction does not cover, and one of them writing into the working directory is enough to fail a round. Where your CLI will not let socrates ignore the servers your machine configures, the review runs anyway and the run says so.
- **It does not read your project's agent instructions.** They are authored by the thing under review, so the reviewer loads user-level settings only. The cost is that it does not know your conventions, and the prompt's "stay in lane" rule is what covers most of what that would otherwise produce.
- **It is not a daemon.** One detached process per round and a directory of files. Nothing to install, nothing to keep running, nothing to restart.
- **It does not tell the reviewer what you told it about how the change came to be.** Only the change, and optionally the one line of intent you chose to state.

## The layers

This repository is built in layers, specified in `openspec/` before they are implemented.

- `a-review-nobody-authored` - the core: freeze, review, findings, verdict, one round. This is what ships today.
- `a-rejection-is-a-claim` - rounds, rebuttals and the dialogue between author and reviewer.
- `one-window-three-needs` - `scrat watch` and `scrat logs`.
- `the-loop-an-agent-can-drive` - the generated skill that lets a session drive the loop itself.

`AGENTS.md` is the working document: layout, verification sequence and the invariants that hold the design up.
