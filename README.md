# socrates

An independent reviewer for the changes you just wrote and nobody has read.

Your session implements something, leaving unreviewed local changes.
`socrates review` freezes the working tree into a commit nothing can move, spawns a detached reviewer against that snapshot, and returns.
The reviewer authored none of the change and is told nothing about how it came to be.
It answers with structured findings, which you read and act on.

The split is deliberate.
The author has the context and does the patching.
The reviewer has independence and wrote none of it.

## What it needs

socrates is a front end for a reviewer it does not contain, so check for the reviewer first.
All of the following must be true before anything below works:

- **The `claude` CLI, 2.1.169 or later, on your `PATH`.**
  That floor is where `--safe-mode` arrived. It is the one mechanism keeping your machine's customisations out of a reviewer, and socrates has no fallback for a CLI without it.
- **An Anthropic account that CLI is signed in to.**
  Every round of every review is a model call, billed to your account against whatever plan or credit it carries.
  socrates itself is free; the reviews it runs are not.
- **Linux or macOS.**
  Two of socrates' guarantees rely on unix mechanisms. See Supported platforms below.
  On Windows, run socrates under WSL, which is a supported Linux.
- **`git`.**
  The snapshot a reviewer reads is a real commit in your repository.
- **Go 1.26 or later, if you install from source.**
  A release archive needs no toolchain.

Check the first two before going further.

## Installing

A release archive is the quickest route.
Each carries `socrates`, `scrat` and `agora` along with `LICENSE` and `NOTICE`, so extracting one gives you the same binaries as a source install.

```sh
tag=$(gh release view --repo ktei/socrates --json tagName --jq .tagName)
gh release download "$tag" --repo ktei/socrates --pattern "socrates_*_$(uname -s | tr 'A-Z' 'a-z')_*.tar.gz" --pattern checksums.txt
sha256sum --check --ignore-missing checksums.txt
tar -xzf socrates_*.tar.gz -C ~/.local/bin socrates scrat agora
socrates skill
```

Check the download against this repository rather than against whoever pushed the tag.
The checksums carry a keyless signature, and the release records build provenance saying which workflow and which commit produced each archive:

```sh
gh release download "$tag" --repo ktei/socrates --pattern 'checksums.txt.bundle'
cosign verify-blob checksums.txt \
  --bundle checksums.txt.bundle \
  --certificate-identity-regexp 'https://github.com/ktei/socrates/.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com
gh attestation verify socrates_*.tar.gz --repo ktei/socrates
```

GitHub stores build provenance for public repositories only, so the last of those three applies from the first release published from one.

`go install github.com/ktei/socrates/cmd/socrates@latest` works too, and gives you `socrates` alone: no `scrat`, no `agora`, and no version stamp, because nothing links one in.

`CHANGELOG.md` says what each release means to somebody who has the previous one.

## First run

Building from source needs Go 1.26 or later:

```sh
make build                     # bin/socrates
bin/socrates review --intent "add the parser I asked for"
```

`make install` puts it on your `PATH` under both of its names: `socrates`, and `scrat` for the short one.
They are one binary, so the two can never be at different commits, and `socrates version` says which release that is and which commit it was built from.

## agora, the second binary

Installing socrates installs a second program beside it, called `agora`.
It is a separate program, not another name for the first.

agora serves one document as one page on a loopback address, waits for a person to answer, and prints what they answered.
That is all it does.
It knows nothing about code review: no finding, verdict, run or round appears anywhere in it.
It renders a document a caller wrote and returns the decision a person made about it.

agora exists because a finding that asks a person is the one thing socrates cannot answer, and such a decision reads better on a page that can quote the code than in a terminal.
Nothing socrates does reaches agora yet, so today you run it yourself and it starts nothing on its own.

```sh
agora serve --document decision.json   # prints the address the page is served at
agora poll --session <session>         # blocks until somebody answers, then prints it
agora stop                             # shuts it down, taking every page with it
```

The two binaries share one home and one optional configuration file, and nothing else.

`review` prints the run identifier and returns while the review runs:

```
20260819T101500-ab12cd34
follow it with: socrates status 20260819T101500-ab12cd34
```

Then either poll it or wait for it:

```sh
socrates status                   # a few lines, and an exit code you can branch on
socrates findings                 # what it found, as socrates classified it
socrates findings --json          # the same, with every effective action resolved
socrates review --wait            # block until the verdict, and exit with its code
socrates cancel                   # stop a review and say that nothing was decided
```

Naming no run means the most recent run for the repository you are standing in.

## Supported platforms

socrates guarantees two things: the reviewer never writes, and everything the reviewer started is reaped with it.
Unix mechanisms hold both, so a platform is supported where both hold and where the suite runs, not merely where socrates compiles.

| platform | reaps by process group | the filesystem enforces the seal | the suite runs there |
| --- | --- | --- | --- |
| Linux | yes | yes | yes |
| macOS | yes | yes | yes |

On Windows, run socrates under WSL, which is a supported Linux and is how socrates itself is developed.

A review started where reaping by process group is unavailable is refused before a snapshot is taken and before a reviewer is spawned.
Reading a finished run is not refused there, because `socrates status`, `socrates findings` and `socrates logs` make neither guarantee.

## Reading a run

| command | behaviour | who it is for |
| --- | --- | --- |
| `socrates watch` | follows the run live, stops when it is over, exits with the verdict's code | you, in a second tab |
| `socrates logs` | prints what the run recorded and exits, paging where your terminal and `PAGER` say so | you, afterwards |
| `socrates status` | one shot, a handful of lines, no stream | the session waiting on the review |

The three are kept separate.
A coding session should call `status`: everything it prints lands in that session's context, and it prints the same few lines whether the reviewer was silent or verbose.
`watch` and `logs` print into a person's terminal.
`status` therefore has no follow mode, and `watch` has no single-read mode.

`socrates watch` stops when the run reaches a verdict, fails, is cancelled from another terminal, or parks the round, whether on a patch you owe it or a question somebody must answer.
It never waits for an event a finished run will not write.
Started between `socrates review` returning and the reviewer's first byte, it waits for the log rather than failing, and gives up if nothing is written.

## Driving it from a coding session

The sequence above is five commands with three branch points, which a session told once in prose will drift from.

```sh
socrates skill                    # writes the skill to ~/.claude/skills/socrates/SKILL.md
socrates skill --to <directory>   # writes it somewhere else instead
```

The skill teaches the loop: implement, `socrates review --intent`, then poll `socrates status` every few seconds until it reports a state a run rests in.
It tells the session not to bound that wait itself, because a reviewer that answers nothing before its deadline ends the run as `failed`.
On `clean`, say what was reviewed and stop.
On `failed`, report the change as unreviewed, because no review completed and nothing has been checked.
On `escalated`, hand the findings to you and stop, because socrates has nothing further it can do about them.
On `awaiting-response`, answer every finding that can be patched, bring every finding that asks a person to you and record what you decided, run `socrates continue`, and go back to polling.
It says the session does not answer a question about intent on its own behalf, and that it runs the next round without asking whether it is worth spending.
Where a run escalated at the round cap, it says that no third round exists and that looking again means a fresh `socrates review`, started only where you ask for one, whose reviewer is told nothing about the run before it.
On `cancelled`, say that nothing was decided, because somebody stopped the review on purpose.
Those are every state a run rests in. A test compares the list against `internal/runstate`, so a sixth cannot appear without the skill gaining a branch for it.

It teaches polling rather than following, because everything a session reads lands in its own context and it is about to need that context for the patching.
It carries the five response rules in the same words `socrates findings` prints them, in front of the session before it has read a single finding.
It also states what socrates will not do: it patches nothing, offers no way to explain the change to the reviewer, and refuses a rejection with no reason.

The text has one source in `internal/skill`, and the copy at `.claude/skills/socrates/SKILL.md` is generated from it by `make skill`.
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
`escalated` means read the findings.
`failed` means socrates could not do its job and nothing was reviewed.

A run escalates for two reasons and no others: a finding holds the change while its own action asks for nothing, so no answer addresses it, or the last permitted round still found something blocking.
A finding that asks a person is not one of them.
It parks the round and returns `4` until somebody answers it.
Ending the run there would leave the answer nowhere to go: you decide, you patch what you decided, and no round remains to judge either.
`socrates respond --decided --reason` records the answer, `--fixed` and `--rejected` are refused for such a finding, and the next round is told a person decided rather than that a patch was claimed.
After the cap, a further look is a fresh `socrates review` rather than a third round, and its reviewer knows nothing of the run before it.

## Effect on your repository

Nothing, other than adding the git objects the snapshot needs.

The snapshot is built through a temporary index, so your own index is untouched.
It includes untracked files, since a freshly implemented feature is mostly untracked files, and it honours `.gitignore`.
No ref names the snapshot commit, so ordinary garbage collection reclaims it once the worktree is gone.
`git status` before and after a review is byte-identical, and a test holds that.

Everything both binaries read and everything they write lives under one home, `~/.socrates` unless `SOCRATES_HOME` names an absolute path instead.
Runs are `~/.socrates/runs/`, one directory per run, legible with `ls` and `cat`.
What agora writes about a page is `~/.socrates/surface/`, beside them.

## Where your code goes

socrates sends nothing anywhere.
No telemetry, no usage counter, no crash report, no update check.
It opens no network connection at all outside a build tag the ordinary suite does not compile, and a test holds that by reading what `cmd/socrates` and every package it reaches import.

The reviewer CLI it drives is the other half, and it is not silent.
socrates spawns `claude`, and that CLI sends the change to its vendor, Anthropic, exactly as it does when you run it yourself.
Your diff, the file contents a reviewer reads and the one line of intent you chose to state go there, under the terms your own account carries.
That is how it works rather than a leak: a review is a model call, and the model is not on your machine.
If that is unacceptable for a repository, do not point socrates at it.

socrates is not affiliated with, endorsed by or sponsored by Anthropic.
"Claude" and "Anthropic" are trademarks of Anthropic PBC, named here only to say which CLI socrates drives.

## Configuration

One optional file, at `~/.socrates/config.yaml`, beside the runs, read by both binaries.
socrates never creates it, never reads configuration out of the repository under review, and refuses rather than quietly running on defaults when the file will not parse.
The one file socrates does read out of a repository is the review guide below, and it configures nothing.

```yaml
roles:
  reviewer:
    model: opus
    effort: high
diff_ceiling: 131072
guidance_ceiling: 16384
agent_deadline: 30m
log_ceiling: 1048576
run_retention: 336h
surface:
  bind: 127.0.0.1:18018
  reached_at: http://localhost:18018/
```

The `surface` block is what an operator states about the page a person reads: the loopback address it binds, and the address a person is handed where that is not the one it bound.

## The review guide

The reviewer arrives knowing nothing about your project, by design.
It loads no `CLAUDE.md`, no project or local settings and no skills, because all of those are written by the thing under review.

One file is the exception, and socrates reads it itself rather than letting the CLI load it:

```sh
mkdir -p .socrates
$EDITOR .socrates/review-guide.md
git add .socrates/review-guide.md && git commit -m "what a reviewer should know about this repository"
```

Write what a reviewer would otherwise get wrong about the repository: the conventions it holds itself to, what a test here is expected to prove, what a directory is for.
Keep it under 16 KiB, or state your own ceiling with `guidance_ceiling` in `~/.socrates/config.yaml`.
A guide over the ceiling is carried whole or not at all, so socrates carries none of it, runs the review anyway, and reports the size and the ceiling.
The ceiling lives in your own configuration, outside every repository socrates reviews, so the repository cannot raise it.

Three things are worth knowing before writing one:

- **It takes effect once committed.**
  socrates reads the guide at the commit your change is measured against, never from your working tree, so the session that wrote the change cannot also write the rules it is judged by.
  A review that finds a guide only in the change says so, and the next review reads it.
- **It adds and removes nothing.**
  It cannot switch off a rule socrates delivers, redefine a severity or an action, exempt a file, or ask for a finding to be withheld.
  `clean` means the same thing in every repository, and a guide cannot weaken it.
- **socrates reads no word of it.**
  It is carried into one section of the reviewer's prompt and nowhere else.
  No verdict, exit code or classification comes from anything it says.

## Reading the page from another machine

The page is served on `127.0.0.1:18018` and can be opened from that machine and no other.
The document quotes the code it was written about, so nothing widens the bind and there is no tunnel and no share command.

Where nothing states an address and something else already holds that port, the surface counts down through a window of ten ports and takes the first one free, rather than refusing over a port nobody chose.
Stating either field of the `surface` block stops it counting.
An operator who states one has aimed something at a port, so a port held by something else is refused rather than walked past.
Serving a page names the port it bound, which is the one to forward.

Where you run socrates in a devcontainer, over ssh, or anywhere else with a network namespace of its own, your browser is on another machine and the address you are handed will not open.
You get a connection that opens and carries nothing, which a browser reports as an empty reply and which looks exactly like a port nobody forwarded.

Forward the port once, and say where the page is then reached at.

In a devcontainer, that is one line of `devcontainer.json`:

```json
{ "forwardPorts": [18018] }
```

Publishing the port in a compose file does not work: it forwards to the container's own interface, and the surface listens on loopback.
A forwarder must run inside the namespace the surface is in.
Your editor's port forwarding does, `ssh -L 18018:127.0.0.1:18018 <host>` does, and a reverse proxy running alongside it does.

Then state where a person reaches it, so the address printed is one that opens:

```yaml
surface:
  reached_at: http://socrates.localhost/
```

Until you do, serving a page says so.

## Out of scope

- **It does not land your change.** No commit, no push, no pull request, no continuous integration. Those are yours.
- **It does not verify your change.** It runs no tests and no linters. `clean` means an independent reviewer found nothing blocking, not that the suite is green.
- **It does not write in your repository.**
  Three mechanisms hold that, in the order they act.

  1. The reviewer's tools cannot write. It gets `Read`, `Grep`, `Glob`, `WebSearch` and `WebFetch`, and none of the ones that edit.
  2. The snapshot worktree it reads is sealed, so the filesystem refuses a write from anything in that session, whatever asked for it.
  3. The round fails if the worktree differs by one byte afterwards, which catches anything that defeated the seal.

  Your machine can take the first two away from socrates. It cannot take the third.
- **It does not load your own customisations.**
  Plugins and external tool sources put tools in the reviewer's session it was never given, memory files and skills can seat it, and a hook runs a command of its own.
  Any of those writing into the working directory is enough to fail a round.

  socrates asks the CLI to disable them with the vendor's own flag, which your machine accepts even where a managed configuration refuses the flags that ask to replace it.
  It then reads the session's tool list rather than trusting the request, and records anything it did not ask for, plus any external source the session reports connecting.
  A source is named before its tools arrive, and socrates asked for none.
  Where such a tool is present anyway, the review still runs and the run says so, which is worth reading: the seal bounds the snapshot worktree, and nothing socrates holds bounds a tool that reaches outside it.

  The flag leaves an administrator's policy-configured hook in force.
  The tool list will not catch one, because a hook contributes no tool, but the other two mechanisms will: the hook's command runs as a child of the reviewer's session, whose working directory is the sealed worktree, so a write it makes there is refused and the comparison afterwards catches one the seal did not.
  Nothing socrates holds bounds where that command reaches outside the worktree.
  On a managed machine, that residual is stated here rather than covered.
- **It will not review where the seal cannot hold.** socrates proves the seal with one write before a reviewer is spawned, and refuses the review if that write succeeds. That covers a filesystem that ignores permission bits, and running as root, which ignores them by design.
- **It does not read your project's agent instructions.**
  They are authored by the thing under review, so the reviewer restricts settings to user level, and disabling your customisations takes the memory files and skills with them.
  The cost is that it does not know your conventions, which the prompt's "stay in lane" rule covers most of.
  The benefit is that the `socrates` skill you installed cannot seat the reviewer either, on top of the reviewer's tools being read-only and a review inside a review being refused.
- **It is not a daemon.** One detached process per round and a directory of files. Nothing to install, nothing to keep running, nothing to restart.
- **It does not tell the reviewer what you told it about how the change came to be.** Only the change, and optionally the one line of intent you chose to state. The skill says the same thing to the session, because a session that believes such a channel exists will go looking for it.

## Licence

socrates is published under the MIT licence, and `LICENSE` at the root states it in full.

`NOTICE` beside it lists everything a shipped binary carries that socrates did not author, with the licence each is published under and that licence's full text.
That is the client bundle's third-party libraries and the font faces embedded with it, and the Go modules `cmd/socrates` and `cmd/agora` link in.
A test fails while anything embedded is not accounted for there, so a new asset cannot arrive without its notice.

The terms are stated in those two files and nowhere else.
No source file carries a licence header, since a header is a comment and this repository's source carries none.

## Architecture

This repository is built in layers, specified in `openspec/` before they are implemented.

- `a-review-nobody-authored` - the core: freeze, review, findings, verdict, one round.
- `a-rejection-is-a-claim` - rounds, rebuttals and the dialogue between author and reviewer.
- `what-a-reviewer-actually-answered` - the corpus of what real reviewers answered, replayed by name in the suite.
- `one-window-three-needs` - `socrates watch` and `socrates logs`.
- `the-loop-an-agent-can-drive` - the generated skill that lets a session drive the loop itself.
- `an-answered-question-is-a-response` - a question the reviewer asked parks the round, and the person's answer is recorded and judged.

All of them ship today.

`AGENTS.md` is the working document, covering layout, the verification sequence and the invariants the design rests on.
