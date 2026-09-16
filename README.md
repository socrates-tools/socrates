# socrates

An independent reviewer for the changes you just wrote and nobody has read, and the page it hands a finding to when only a person can answer it.

Your session implements something, leaving unreviewed local changes.
`socrates review begin` freezes the working tree into a commit nothing can move, spawns a detached reviewer against that snapshot, and returns.
The reviewer authored none of the change and is told nothing about how it came to be.
It answers with structured findings, which you read and act on.

The split is deliberate.
The author has the context and does the patching.
The reviewer has independence and wrote none of it.

## What it needs

socrates review is a front end for a reviewer it does not contain, so check for the reviewer first.
All of the following must be true before anything below works:

- **One of the reviewer CLIs socrates review drives, on your `PATH`: `claude`, `codex` or `opencode`.**
  socrates review runs `claude` unless you name another in `reviewer_cli`. See Reviewer CLIs below for what each one asks of you.
  The `claude` floor is 2.1.169 or later, which is where `--safe-mode` arrived.
- **An account that CLI is signed in to.**
  Every round of every review is a model call, billed to your account with that CLI's vendor against whatever plan or credit it carries.
  socrates itself is free; the reviews it runs are not.
- **Linux or macOS.**
  Two of socrates review's guarantees rely on unix mechanisms. See Supported platforms below.
  On Windows, run socrates under WSL, which is a supported Linux.
- **`git`.**
  The snapshot a reviewer reads is a real commit in your repository.
- **Go 1.26 or later, if you install from source.**
  A release archive needs no toolchain.

Check the first two before going further.

## Installing

A release archive is the quickest route.
It carries `socrates` along with `LICENSE` and `NOTICE`, so extracting one gives you the same binary as a source install.

```sh
tag=$(gh release view --repo socrates-tools/socrates --json tagName --jq .tagName)
gh release download "$tag" --repo socrates-tools/socrates --pattern "socrates_*_$(uname -s | tr 'A-Z' 'a-z')_*.tar.gz" --pattern checksums.txt
sha256sum --check --ignore-missing checksums.txt
tar -xzf socrates_*.tar.gz -C ~/.local/bin socrates
socrates review skill
```

Check the download against this repository rather than against whoever pushed the tag.
The checksums carry a keyless signature, and the release records build provenance saying which workflow and which commit produced each archive:

```sh
gh release download "$tag" --repo socrates-tools/socrates --pattern 'checksums.txt.bundle'
cosign verify-blob checksums.txt \
  --bundle checksums.txt.bundle \
  --certificate-identity-regexp 'https://github.com/socrates-tools/socrates/.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com
gh attestation verify socrates_*.tar.gz --repo socrates-tools/socrates
```

GitHub stores build provenance for public repositories only, so the last of those three applies from the first release published from one.

`go install github.com/socrates-tools/socrates/cmd/socrates@latest` works too, and gives you `socrates` alone, with no version stamp, because nothing links one in.

`CHANGELOG.md` says what each release means to somebody who has the previous one.

## First run

Building from source needs Go 1.26 or later:

```sh
make build                     # bin/socrates
bin/socrates review begin --intent "add the parser I asked for"
```

`make install` puts it on your `PATH` under its one name, `socrates`.
Every workflow and component the suite holds is reached through that one name, so `socrates version` says which release this is and which commit it was built from, whichever of them you asked about.

`begin` prints the run identifier and returns while the review runs:

```
20260819T101500-ab12cd34
follow it with: socrates review status 20260819T101500-ab12cd34
```

Then either poll it or wait for it:

```sh
socrates review status                   # a few lines, and an exit code you can branch on
socrates review findings                 # what it found, as socrates review classified it
socrates review findings --json          # the same, with every effective action resolved
socrates review begin --wait             # block until the verdict, and exit with its code
socrates review cancel                   # stop a review and say that nothing was decided
```

Naming no run means the most recent run for the repository you are standing in.

A review of a working tree that holds nothing its base commit does not is refused rather than reported clean.
Nothing was read, so there is nothing to pass, and no run is written for a session to act on.
Where the answers to a round revert the whole change, `socrates review continue` ends that run `failed` for the same reason: `clean` is reported for no round a reviewer did not read.

## Reviewing a branch, and reviewing a change that is not yours

The default base is the commit `HEAD` resolves to, so `socrates review begin` on its own reviews what you have not committed.
Two flags widen that.

```sh
socrates review begin --base origin/develop                    # everything this branch added
socrates review begin --base origin/main --not-my-change       # somebody else's change, which you cannot patch
```

`--base <ref>` resolves to the fork point, which is `merge-base(<ref>, HEAD)`, and reviews everything between that commit and your working tree.
On a long-running branch that is the whole of what the branch added, committed and uncommitted alike, and none of what the trunk has done since.
Resolving to the fork point rather than to the named commit is what stops every file the trunk changed being reported to the reviewer as a file your branch deleted.
Rebase onto the trunk first and the two readings coincide, which is why rebasing is the right instinct rather than a workaround.

socrates review neither fetches nor rebases.
Fetching would reach the network and rebasing would write in the repository under review, so your remote refs are as fresh as you last made them and your commits are as you left them.
It also says nothing about how far the named ref has moved past the fork point: `socrates review status` prints the base commit it resolved, as a bare fact, and you decide whether that is what you meant.
A ref that does not resolve, and a ref with no commit in common with `HEAD`, are refused before a snapshot is taken and before a reviewer is spawned.

`--not-my-change` says the change is not yours to patch.
Nothing about it is discovered, because discovering it means asking a forge who authored a commit, and socrates review talks to no forge.
The head of a review is still the working tree it is run in, so check their branch out in a worktree of its own and start the review inside it:

```sh
git worktree add --detach ../reviewing-their-change origin/their-branch
```

Otherwise what gets reviewed is their branch plus whatever you have not committed, and you would be commenting on your own mess in their pull request.

The findings of such a run take the resolutions that make sense for somebody else's code:

```sh
socrates review respond --finding f1 --forwarded
socrates review respond --finding f2 --dropped --reason "true, and not worth raising on a change this size"
socrates review respond --finding f3 --rejected --reason "the reviewer misread the lock, which the caller takes"
```

A forward needs no reason, because what gets raised is the finding's own recorded text.
A drop and a rejection each need one, and neither reason ever leaves your machine.
No second round follows, since nothing here was patched for a fresh reviewer to judge, so `socrates review continue` ends the run `escalated` and says the findings were handed to whoever wrote the change.

socrates review posts nothing.
`socrates review findings --json` reports, per finding, the resolution you recorded and whether the finding's location is inside the change it reviewed.
Your session reads that and posts with its own tooling: a finding inside the change can be a comment on that line, and the rest belong in one general comment.
Where the base resolves to the fork point, those are the same hunks a forge renders for a pull request opened against that branch.

## agora, the page a round is handed to

`socrates agora` is a second namespace inside the one binary, not a second program.
It knows nothing about code review: no finding, verdict, run or round appears anywhere in it.

socrates agora serves one document as one page on a loopback address, waits for a person to answer, and prints what they answered.
That is all it does.
It renders a document a caller wrote and returns the decision a person made about it.

agora exists because a finding that asks a person is the one thing socrates review cannot answer, and such a decision reads better on a page that can quote the code than in a terminal.
Nothing socrates review does reaches agora yet, so today you run it yourself and it starts nothing on its own.

```sh
socrates agora serve --document decision.json   # prints the address the page is served at
socrates agora poll --session <session>         # blocks until somebody answers, then prints it
socrates agora stop                             # shuts it down, taking every page with it
```

review and agora share one home and one optional configuration file, and nothing else.

## Reviewer CLIs

socrates review does not contain a reviewer.
It drives one of these, and which one is yours to state rather than socrates review's to guess from what is installed.

A CLI qualifies on three things socrates review cannot repair after the fact.
The reviewed repository's own instruction files can be kept out of the reviewer's session, because a review by a session the change instructed is defective in its judgement rather than in its format.
A terminal event says whether the work happened and what the agent said where it did not, because a run that fails silently is one you cannot act on.
A model is selectable per invocation, because per-role routing exists so the reviewer gets the strongest model available.

| reviewer CLI | how it arrives at a findings document | what keeps the repository out of the session |
| --- | --- | --- |
| `claude` | the vendor constrains it against the schema | `--setting-sources user` |
| `codex` | the vendor constrains it against the schema | `-c project_doc_max_bytes=0` |
| `opencode` | socrates review extracts it from the reviewer's final message | `OPENCODE_DISABLE_PROJECT_CONFIG=1` |

Both ways reach the same check: socrates review validates the document against the schema itself, and a document that does not parse or does not satisfy it fails the round either way.
What the first column buys is reliability rather than correctness.
Where the vendor constrains the output, the document conforms on the first attempt more often; where socrates review extracts it, a reviewer that answers without one fails the round rather than producing a wrong verdict.
If you have accounts for more than one, prefer a CLI in the first two rows.

These three are not supported, and what each lacks is the first criterion:

| reviewer CLI | what socrates review found no way to get from it |
| --- | --- |
| `cursor-agent` | a way to keep the reviewed repository's own instructions out of the reviewer's session |
| `copilot` | a way to keep the reviewed repository's own instructions out of the reviewer's session, and a terminal event saying whether the work happened and what it says about why it did not |
| `pi` | a way to keep the reviewed repository's own instructions out of the reviewer's session, and a terminal event saying whether the work happened and what it says about why it did not |

No such mechanism was found in their documentation, which is weaker than a finding that none exists.
Settling each is one recording rather than a redesign, and https://github.com/socrates-tools/socrates/issues/20 is where that is tracked.
A CLI that clears the bar is a new row here rather than a new design.

Naming a CLI socrates review does not drive is refused by name, before a snapshot is taken and before a reviewer is spawned.

## Supported platforms

socrates review guarantees two things: the reviewer never writes, and everything the reviewer started is reaped with it.
Unix mechanisms hold both, so a platform is supported where both hold and where the suite runs, not merely where socrates compiles.

| platform | reaps by process group | the filesystem enforces the seal | the suite runs there |
| --- | --- | --- | --- |
| Linux | yes | yes | yes |
| macOS | yes | yes | yes |

On Windows, run socrates under WSL, which is a supported Linux and is how socrates itself is developed.

A review started where reaping by process group is unavailable is refused before a snapshot is taken and before a reviewer is spawned.
Reading a finished run is not refused there, because `socrates review status`, `socrates review findings` and `socrates review logs` make neither guarantee.

## Reading a run

| command | behaviour | who it is for |
| --- | --- | --- |
| `socrates review watch` | follows the run live, stops when it is over, exits with the verdict's code | you, in a second tab |
| `socrates review logs` | prints what the run recorded and exits, paging where your terminal and `PAGER` say so | you, afterwards |
| `socrates review status` | one shot, a handful of lines, no stream | the session waiting on the review |

The three are kept separate.
A coding session should call `status`: everything it prints lands in that session's context, and it prints the same few lines whether the reviewer was silent or verbose.
`watch` and `logs` print into a person's terminal.
`status` therefore has no follow mode, and `watch` has no single-read mode.

`socrates review watch` stops when the run reaches a verdict, fails, is cancelled from another terminal, or parks the round, whether on a patch you owe it or a question somebody must answer.
It never waits for an event a finished run will not write.
Started between `socrates review begin` returning and the reviewer's first byte, it waits for the log rather than failing, and gives up if nothing is written.

## Driving it from a coding session

The sequence above is five commands with three branch points, which a session told once in prose will drift from.

```sh
socrates review skill                    # writes both skills under ~/.claude/skills/
socrates review skill --to <directory>   # writes them somewhere else instead
```

There are two, and a session reads one of them.
`socrates-review` teaches reviewing the change that session just wrote.
`socrates-review-handover` teaches reviewing a change somebody else wrote: the throwaway worktree, the author's own description relayed verbatim and acted on in no way, the three resolutions, and posting what the person chose and nothing else.
They are separate documents rather than one with a branch in it, because a session reviewing its own work should not spend its context on the forwarding half.

The skill teaches the loop: implement, `socrates review begin --intent`, then poll `socrates review status` every few seconds until it reports a state a run rests in.
It tells the session not to bound that wait itself, because a reviewer that answers nothing before its deadline ends the run as `failed`.
On `clean`, say what was reviewed and stop.
On `failed`, report the change as unreviewed, because no review completed and nothing has been checked.
On `escalated`, hand the findings to you and stop, because socrates review has nothing further it can do about them.
On `awaiting-response`, answer every finding that can be patched, bring every finding that asks a person to you and record what you decided, run `socrates review continue`, and go back to polling.
It says the session does not answer a question about intent on its own behalf, and that it runs the next round without asking whether it is worth spending.
Where a run escalated at the round cap, it says that no third round exists and that looking again means a fresh run of `socrates review begin`, started only where you ask for one, whose reviewer is told nothing about the run before it.
On `cancelled`, say that nothing was decided, because somebody stopped the review on purpose.
Those are every state a run rests in. A test compares the list against `review/runstate`, so a sixth cannot appear without the skill gaining a branch for it.

It teaches polling rather than following, because everything a session reads lands in its own context and it is about to need that context for the patching.
It carries the five response rules in the same words `socrates review findings` prints them, in front of the session before it has read a single finding.
It also states what socrates review will not do: it patches nothing, offers no way to explain the change to the reviewer, and refuses a rejection with no reason.

Both texts have one source in `review/skill`, and the copies under `.claude/skills/` are generated from it by `make skill`.
A test fails while a copy and its source differ, and another checks that every command and flag either skill names exists on the built binary.

Exit codes are a contract, shared by every namespace the suite holds:

| code | meaning |
| --- | --- |
| `0` | `clean`, or a command whose job was to start something that started |
| `1` | `escalated`: socrates review is finished and what is left is yours to decide |
| `2` | `failed`: no review completed, so nothing has been checked |
| `3` | a usage or precondition error |
| `4` | the run is not terminal yet |

`failed` is deliberately not `escalated`.
`escalated` means read the findings.
`failed` means socrates review could not do its job and nothing was reviewed.

A run escalates for two reasons and no others: a finding holds the change while its own action asks for nothing, so no answer addresses it, or the last permitted round still found something blocking.
A finding that asks a person is not one of them.
It parks the round and returns `4` until somebody answers it.
Ending the run there would leave the answer nowhere to go: you decide, you patch what you decided, and no round remains to judge either.
`socrates review respond --decided --reason` records the answer, `--fixed` and `--rejected` are refused for such a finding, and the next round is told a person decided rather than that a patch was claimed.
After the cap, a further look is a fresh run of `socrates review begin` rather than a third round, and its reviewer knows nothing of the run before it.

## Effect on your repository

Nothing, other than adding the git objects the snapshot needs.

The snapshot is built through a temporary index, so your own index is untouched.
It includes untracked files, since a freshly implemented feature is mostly untracked files, and it honours `.gitignore`.
No ref names the snapshot commit, so ordinary garbage collection reclaims it once the worktree is gone.
`git status` before and after a review is byte-identical, and a test holds that.

Everything review and agora read and everything they write lives under one home, `~/.socrates/review`.
`~/.socrates` is the root the `socrates` brand keeps, and each tool published under it owns one directory there, so `SOCRATES_HOME` moves the root and this tool follows it.
Runs are `~/.socrates/review/runs/`, one directory per run, legible with `ls` and `cat`.
What agora writes about a page is `~/.socrates/review/surface/`, beside them.

## Where your code goes

socrates review sends nothing anywhere.
No telemetry, no usage counter, no crash report, no update check.
It opens no network connection at all outside a build tag the ordinary suite does not compile, and a test holds that by reading what `review/command` and every package it reaches import.

The reviewer CLI it drives is the other half, and it is not silent.
socrates spawns the CLI you named, and that CLI sends the change to its own vendor, exactly as it does when you run it yourself.
`claude` sends it to Anthropic, `codex` to OpenAI, `opencode` to whichever provider that CLI is configured against.
Your diff, the file contents a reviewer reads and the one line of intent you chose to state go there, under the terms your own account carries.
That is how it works rather than a leak: a review is a model call, and the model is not on your machine.
If that is unacceptable for a repository, do not point socrates review at it.

socrates is not affiliated with, endorsed by or sponsored by any of these vendors.
"Claude" and "Anthropic" are trademarks of Anthropic PBC, "Codex" and "OpenAI" of OpenAI, and each is named here only to say which CLIs socrates review drives.

## Configuration

One optional file, at `~/.socrates/review/config.yaml`, beside the runs, read by review and agora alike.
socrates review never creates it, never reads configuration out of the repository under review, and refuses rather than quietly running on defaults when the file will not parse.
The one file socrates review does read out of a repository is the review guide below, and it configures nothing.

```yaml
reviewer_cli: claude
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

`socrates review config` prints every setting, the value in force on this machine and whether you set it or socrates review defaulted it, so the sample above is a sample and that command is the report.

`socrates review configure` changes one, from a menu, without you opening the file.
It lists the same settings `socrates review config` reports, asks for a value, writes the file as soon as you confirm, and comes back for the next one.
Every setting also offers `use the default`, which deletes its line rather than writing today's default in, so a later release still reaches you.
Your comments, your key order and your blank lines survive the write, and only the key you chose is touched.
It draws a menu and waits for you, so it refuses where nothing is a terminal; `socrates review config` is the way to read settings from a script.

`reviewer_cli` names which of the reviewer CLIs above socrates review drives.
An absent field runs `claude`, which is what socrates review shipped with, so an existing configuration keeps working.
Naming one socrates review does not drive is refused.

Each CLI names its models its own way, so a `model` you state under `roles` is passed to the CLI you named and nothing translates it.
Leave it out to let that CLI choose its own.

The `surface` block is what an operator states about the page a person reads: the loopback address it binds, and the address a person is handed where that is not the one it bound.

## The review guide

The reviewer arrives knowing nothing about your project, by design.
It loads no `CLAUDE.md`, no `AGENTS.md`, no project or local settings and no skills, because all of those are written by the thing under review.

One file is the exception, and socrates review reads it itself rather than letting the CLI load it:

```sh
mkdir -p .socrates
$EDITOR .socrates/review-guide.md
git add .socrates/review-guide.md && git commit -m "what a reviewer should know about this repository"
```

Write what a reviewer would otherwise get wrong about the repository: the conventions it holds itself to, what a test here is expected to prove, what a directory is for.
Keep it under 16 KiB, or state your own ceiling with `guidance_ceiling` in `~/.socrates/review/config.yaml`.
A guide over the ceiling is carried whole or not at all, so socrates review carries none of it, runs the review anyway, and reports the size and the ceiling.
The ceiling lives in your own configuration, outside every repository socrates review reviews, so the repository cannot raise it.

Three things are worth knowing before writing one:

- **It takes effect once committed.**
  socrates review reads the guide at the commit your change is measured against, never from your working tree, so the session that wrote the change cannot also write the rules it is judged by.
  A review that finds a guide only in the change says so, and the next review reads it.
- **It adds and removes nothing.**
  It cannot switch off a rule socrates review delivers, redefine a severity or an action, exempt a file, or ask for a finding to be withheld.
  `clean` means the same thing in every repository, and a guide cannot weaken it.
- **socrates review reads no word of it.**
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
  reached_at: http://socrates-review.localhost/
```

Until you do, serving a page says so.

## Out of scope

- **It does not land your change.** No commit, no push, no pull request, no continuous integration. Those are yours.
- **It does not fetch, and it does not rebase.** You fetch, and you rebase. socrates review resolves a base out of what your repository already holds.
- **It does not post a comment, and it knows nothing of pull requests.** No forge, no authentication, no network. It records what you decided about a finding and reports whether the finding is inside the change; your session posts.
- **It does not approve anything.** Approving outlasts the review and is a governance act. A review that finds nothing serves a page with no options and offers no approval.
- **It does not verify your change.** It runs no tests and no linters. `clean` means an independent reviewer found nothing blocking, not that the suite is green.
- **It does not write in your repository.**
  Three mechanisms hold that, in the order they act.

  1. The reviewer's tools cannot write. It gets `Read`, `Grep`, `Glob`, `WebSearch` and `WebFetch`, and none of the ones that edit.
  2. The snapshot worktree it reads is sealed, so the filesystem refuses a write from anything in that session, whatever asked for it.
  3. The round fails if the worktree differs by one byte afterwards, which catches anything that defeated the seal.

  Your machine can take the first two away from socrates review. It cannot take the third.
- **It does not load your own customisations.**
  Plugins and external tool sources put tools in the reviewer's session it was never given, memory files and skills can seat it, and a hook runs a command of its own.
  Any of those writing into the working directory is enough to fail a round.

  socrates review asks the CLI to disable them with the vendor's own flag, which your machine accepts even where a managed configuration refuses the flags that ask to replace it.
  It then reads the session's tool list rather than trusting the request, and records anything it did not ask for, plus any external source the session reports connecting.
  A source is named before its tools arrive, and socrates review asked for none.
  Where such a tool is present anyway, the review still runs and the run says so, which is worth reading: the seal bounds the snapshot worktree, and nothing socrates review holds bounds a tool that reaches outside it.

  The flag leaves an administrator's policy-configured hook in force.
  The tool list will not catch one, because a hook contributes no tool, but the other two mechanisms will: the hook's command runs as a child of the reviewer's session, whose working directory is the sealed worktree, so a write it makes there is refused and the comparison afterwards catches one the seal did not.
  Nothing socrates review holds bounds where that command reaches outside the worktree.
  On a managed machine, that residual is stated here rather than covered.
- **It will not review where the seal cannot hold.** socrates review proves the seal with one write before a reviewer is spawned, and refuses the review if that write succeeds. That covers a filesystem that ignores permission bits, and running as root, which ignores them by design.
- **It does not read your project's agent instructions.**
  They are authored by the thing under review, so the reviewer restricts settings to user level, and disabling your customisations takes the memory files and skills with them.
  The cost is that it does not know your conventions, which the prompt's "stay in lane" rule covers most of.
  The benefit is that the `socrates-review` skill you installed cannot seat the reviewer either, on top of the reviewer's tools being read-only and a review inside a review being refused.
- **It is not a daemon.** One detached process per round and a directory of files. Nothing to install, nothing to keep running, nothing to restart.
- **It does not tell the reviewer what you told it about how the change came to be.** Only the change, and optionally the one line of intent you chose to state. The skill says the same thing to the session, because a session that believes such a channel exists will go looking for it.

## Licence

socrates is published under the MIT licence, and `LICENSE` at the root states it in full.

`NOTICE` beside it lists everything a shipped binary carries that socrates did not author, with the licence each is published under and that licence's full text.
That is the client bundle's third-party libraries and the font faces embedded with it, and the Go modules `cmd/socrates` links in.
A test fails while anything embedded is not accounted for there, so a new asset cannot arrive without its notice.

The terms are stated in those two files and nowhere else.
No source file carries a licence header, since a header is a comment and this repository's source carries none.

## Architecture

This repository is built in layers, each specified before it was implemented.

- The core: freeze, review, findings, verdict, one round.
- Rounds, rebuttals and the dialogue between author and reviewer.
- The corpus of what real reviewers answered, replayed by name in the suite.
- `socrates review watch` and `socrates review logs`.
- The generated skill that lets a session drive the loop itself.
- A question the reviewer asked parks the round, and the person's answer is recorded and judged.
- The page a round is handed over on, served by agora.
- A base you name, and a change that is not yours to patch.
- One binary, one dispatcher, and a namespace per workflow or component.

All of them ship today.

`AGENTS.md` is the working document, covering layout, the verification sequence and the invariants the design rests on.
`CONTEXT.md` holds the vocabulary, and `docs/adr/` the decisions behind it.
