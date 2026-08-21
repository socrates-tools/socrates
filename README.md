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

## Three ways to read one run, for three readers

| command | behaviour | who it is for |
| --- | --- | --- |
| `scrat watch` | follows the run live, stops when it is over, exits with the verdict's code | you, in a second tab |
| `scrat logs` | prints what the run recorded and exits, paging where your terminal and `PAGER` say so | you, afterwards |
| `scrat status` | one shot, a handful of lines, no stream | the session waiting on the review |

They stay distinct on purpose.
`status` is what a coding session should call, because everything it prints lands in that session's context, and it prints the same few lines for a silent reviewer and a loud one.
`watch` and `logs` print into a person's terminal, where being told what happened is the point.
So `status` has no way to follow, and `watch` has no way to return after one read.

`scrat watch` stops when the run reaches a verdict, when it fails, when you cancel it from another terminal, and when the round parks, whether on a patch you owe it or on a question somebody has to answer.
It never waits for an event a finished run will not write.
Started in the moment between `scrat review` returning and the reviewer's first byte, it waits for the log rather than failing, and gives up if nothing is ever written.

## Letting a session drive it

The sequence above is five commands with three branch points, and a coding session told it once in prose gets it right for a while and then drifts.

```sh
scrat skill                    # writes the skill to ~/.claude/skills/scrat/SKILL.md
scrat skill --to <directory>   # writes it somewhere else instead
```

The skill teaches the loop: implement, `scrat review --intent`, then poll `scrat status` every few seconds until it reports a state a run rests in.
It tells the session not to bound that wait itself, because a reviewer that answers nothing before its deadline ends the run as `failed`.
On `clean`, say what was reviewed and stop.
On `failed`, report the change as unreviewed, because no review completed and nothing has been checked.
On `escalated`, hand the findings to you and stop, because socrates has nothing further it can do about them.
On `awaiting-response`, answer every finding that can be patched, bring every finding that asks a person to you and record what you decided, run `scrat continue`, and go back to polling.
It says the session does not answer a question about intent on its own behalf, and that the round the findings earned is spent without asking you whether it is worth spending.
Where a run escalated at the round cap, it says that no third round exists and that looking again means a fresh `scrat review`, started only where you ask for one, whose reviewer is told nothing about the run before it.
On `cancelled`, say that nothing was decided, because somebody stopped the review on purpose.
Those are every state a run rests in, and a test compares the list against `internal/runstate` so a sixth could not appear without the skill gaining a branch for it.

It teaches polling rather than following, because everything a session reads lands in its own context and it is about to need that context for the patching.
It carries the five response rules in the same words `scrat findings` prints them, in front of the session before it has read a single finding.
And it says what socrates will not do: it patches nothing, it offers no way to explain the change to the reviewer, and it refuses a rejection with no reason.

The text has one source in `internal/skill`, and the copy at `.claude/skills/scrat/SKILL.md` is generated from it by `make skill`.
A test fails while those two differ, and another checks that every command and flag the skill names exists on the built binary.

Exit codes are a contract:

| code | meaning |
| --- | --- |
| `0` | `clean`, or a command whose job was to start something that started |
| `1` | `escalated`: socrates is finished and what is left is yours to decide |
| `2` | `failed`: no review completed, so nothing has been checked |
| `3` | a usage or precondition error |
| `4` | the run is not terminal yet |

`failed` is deliberately not `escalated`.
Escalated means read the findings.
Failed means socrates could not do its job and nothing was reviewed.

A run escalates for two reasons and no others: a finding holds the change while its own action asks for nothing, so no answer addresses it, or the last permitted round still found something blocking.
A finding that asks a person is not one of them.
It parks the round and returns `4` until somebody answers it, because ending the run there is what leaves the answer nowhere to go: you decide, you patch what you decided, and no round remains to judge either.
`scrat respond --decided --reason` records the answer, `--fixed` and `--rejected` are refused for such a finding, and the next round is told a person decided rather than that a patch was claimed.
After the cap, a further look is a fresh `scrat review` rather than a third round, and its reviewer knows nothing of the run before it.

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
- **It does not write in your repository.** Three things hold that, in the order they act. The reviewer's tools cannot write. The snapshot worktree it reads is sealed, so the filesystem refuses a write from anything in that session, whether the model asked for it or a tool server did. And the round fails if the worktree differs by one byte afterwards, which is what catches anything that defeated the seal. Your machine can take the first two away from socrates; it cannot take the third.
- **It does not load your tool servers.** A tool server is a second source of tools the vendor's own tool restriction does not cover, and one of them writing into the working directory is enough to fail a round. Where your CLI will not let socrates ignore the servers your machine configures, the review runs anyway and the run says so, because the seal holds the reviewer read-only whatever is loaded.
- **It will not review where the seal cannot hold.** socrates proves the seal with one write before a reviewer is spawned, and refuses the review if that write succeeds. That covers a filesystem that ignores permission bits and it covers running as root, which ignores them by design. A verdict it cannot stand behind is worse than no verdict.
- **It does not read your project's agent instructions.** They are authored by the thing under review, so the reviewer loads user-level settings only. The cost is that it does not know your conventions, and the prompt's "stay in lane" rule is what covers most of what that would otherwise produce. One consequence of loading user-level settings is that a reviewer can see the `scrat` skill you installed there. It cannot act on it: a reviewer's tools are read-only, and a review inside a review is refused.
- **It is not a daemon.** One detached process per round and a directory of files. Nothing to install, nothing to keep running, nothing to restart.
- **It does not tell the reviewer what you told it about how the change came to be.** Only the change, and optionally the one line of intent you chose to state. The skill says the same thing to the session, because a session that believes such a channel exists will go looking for it.

## The layers

This repository is built in layers, specified in `openspec/` before they are implemented.

- `a-review-nobody-authored` - the core: freeze, review, findings, verdict, one round.
- `a-rejection-is-a-claim` - rounds, rebuttals and the dialogue between author and reviewer.
- `what-a-reviewer-actually-answered` - the corpus of what real reviewers answered, replayed by name in the suite.
- `one-window-three-needs` - `scrat watch` and `scrat logs`.
- `the-loop-an-agent-can-drive` - the generated skill that lets a session drive the loop itself.
- `an-answered-question-is-a-response` - a question the reviewer asked parks the round, and the person's answer is recorded and judged.

Every one of them ships today.

`AGENTS.md` is the working document: layout, verification sequence and the invariants that hold the design up.
