# socrates

[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/socrates-tools/socrates/badge)](https://scorecard.dev/viewer/?uri=github.com/socrates-tools/socrates)

socrates puts an independent reviewer in front of changes nobody has read yet.

Your coding session writes a change.
`socrates review begin` freezes it into a snapshot and starts a reviewer that did not write it and is told nothing about how it was written.
The reviewer answers with findings.
Your session fixes what it can, brings you what only you can decide, and asks for another round.

The same binary holds three more tools:

- `socrates agora` shows a person one page of decisions in the browser and hands back their answers.
- `socrates design` is a skill that interviews you over that page and leaves a spec behind.
- `socrates implement` is a skill that builds the tickets of that spec through implementers while you stay free to talk.

## What it needs

- **A reviewer CLI on your `PATH`: `claude`, `codex` or `opencode`.**
  socrates runs `claude` unless you set `reviewer_cli`.
  `claude` must be 2.1.169 or later.
- **An account that CLI is signed in to.**
  socrates is free, but every review round is a model call billed to that account.
- **Linux or macOS.**
  On Windows, use WSL.
- **`git`.**
- **Go 1.26 or later**, only if you build from source.

## Installing

```sh
curl -fsSL https://github.com/socrates-tools/socrates/releases/latest/download/install.sh | sh
```

That downloads the archive for your operating system and architecture from the latest release, checks it against the release's checksums, puts `socrates` in `~/.local/bin`, and installs every skill.
Where `cosign` is on your `PATH`, it also verifies the release's signature before trusting the checksums.
It refuses, and installs nothing, where a check fails or the release carries no archive for your machine.
It says so where `~/.local/bin` is not on your `PATH`, and how to add it.

`SOCRATES_INSTALL_DIR` puts the binary somewhere else, and `SOCRATES_VERSION` installs another release than the latest:

```sh
curl -fsSL https://github.com/socrates-tools/socrates/releases/latest/download/install.sh | SOCRATES_INSTALL_DIR=/usr/local/bin SOCRATES_VERSION=v0.2.0 sh
```

Running it again upgrades the binary and the skills.

### By hand

```sh
os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m | sed 's/x86_64/amd64/; s/aarch64/arm64/')
base=https://github.com/socrates-tools/socrates/releases/latest/download
curl -fsSLO "$base/checksums.txt"
archive=$(grep -o "socrates_[^ ]*_${os}_${arch}\.tar\.gz$" checksums.txt)
curl -fsSLO "$base/$archive"
sha256sum --check --ignore-missing checksums.txt
mkdir -p ~/.local/bin && tar -xzf "$archive" -C ~/.local/bin socrates
~/.local/bin/socrates skill
```

On a Mac without `sha256sum`, check with `shasum -a 256 --check --ignore-missing checksums.txt` instead.

To check the release was built by this repository, verify the signature and the build provenance:

```sh
curl -fsSLO "$base/checksums.txt.sigstore.json"
cosign verify-blob checksums.txt \
  --bundle checksums.txt.sigstore.json \
  --certificate-identity-regexp 'https://github.com/socrates-tools/socrates/.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com
gh attestation verify "$archive" --repo socrates-tools/socrates
```

The release also carries the same provenance as `socrates.intoto.jsonl`, for a verifier that reads it from the release rather than from GitHub.
Releases before v0.4.0 publish the signature as `checksums.txt.bundle` and carry no `socrates.intoto.jsonl`.

### From source

`make install` puts `socrates` on your `PATH` and installs every skill.
`go install github.com/socrates-tools/socrates/cmd/socrates@latest` also works, but that binary does not know its own version, and you run `socrates skill` yourself afterwards.

`socrates version` says which release and commit you have.
`CHANGELOG.md` says what changed in each release.

### The skills

`socrates skill` installs every skill your coding session reads, which is what the installer and `make install` run:

- `socrates-review` and `socrates-review-handover`, from `socrates review skill`.
- `socrates-design` and `socrates-tickets`, from `socrates design skill`.
- `socrates-implement`, and the `socrates-implementer` agent, from `socrates implement skill`.

They write under `~/.claude/skills/`, and the implementer's agent definition under `~/.claude/agents/`.
To put one tool's skills somewhere else, run its own command with `--to`, and the implementer's agent definition with `--agents-to`.

## Reviewing a change

```sh
socrates review begin --intent "add the parser I asked for"
```

`begin` prints a run identifier and returns while the review runs.
Every other command reads the most recent run in the current repository unless you name one.

```sh
socrates review status            # a few lines and an exit code, for a coding session
socrates review findings          # what the reviewer found
socrates review findings --json   # the same, as JSON
socrates review watch             # follow it live, for you in a second terminal
socrates review logs              # everything the run recorded, afterwards
socrates review cancel            # stop it
socrates review begin --wait      # start and block until it finishes or needs you
```

A coding session should poll `status`, because everything it reads fills its context.
`watch` and `logs` are for a person.

Answer each finding, then start the next round:

```sh
socrates review respond --finding f1 --fixed --note "guarded the nil map"
socrates review respond --finding f2 --rejected --reason "the caller already holds the lock"
socrates review respond --finding f3 --decided --reason "keep the old flag for one release"
socrates review continue
```

A rejection is a claim the next reviewer judges, not a way to hide a finding, so it needs a reason.
A finding that asks a question waits for a person, and `--decided` records what they said.
There are at most two rounds.
After that, looking again means a new `begin`.

A working tree with nothing to review is refused rather than reported `clean`.

### Run states and exit codes

| state | exit code | what to do |
| --- | --- | --- |
| `starting`, `reviewing` | `4` | wait |
| `awaiting-response` | `4` | answer the findings, then `continue` |
| `clean` | `0` | the reviewer found nothing blocking |
| `escalated` | `1` | read the findings: what is left is yours to decide |
| `failed` | `2` | nothing was reviewed |
| `cancelled` | `2` | somebody stopped it, and nothing was decided |

Exit code `3` means a usage error or a missing precondition.
A command that only starts something exits `0` once it has started.

### Reviewing a branch, or somebody else's change

```sh
socrates review begin --base origin/main                   # everything this branch added
socrates review begin --base origin/main --not-my-change   # a change you cannot patch
```

`--base` reviews everything since the branch forked from that ref, committed or not.
socrates never fetches or rebases, so fetch first if you want a fresh base.
`status` prints the base commit it used.

For somebody else's change, check it out in its own worktree so your uncommitted work stays out of it:

```sh
git worktree add --detach ../their-change origin/their-branch
```

Each finding is then forwarded, dropped or rejected:

```sh
socrates review respond --finding f1 --forwarded
socrates review respond --finding f2 --dropped --reason "true, but not worth raising here"
socrates review respond --finding f3 --rejected --reason "the reviewer misread the lock"
```

`continue` then ends the run `escalated` and hands the forwarded findings over.
socrates posts nothing.
`findings --json` says whether each finding's line is inside the change, so your session can post it as a line comment or a general one.

### Letting a coding session drive it

The `socrates-review` skill teaches a session the whole loop: begin, poll, answer, continue, and bring you the questions that are yours.
The `socrates-review-handover` skill does the same for somebody else's change, and posts only what you chose.

## agora

`socrates agora` serves one page of decisions on a loopback address, waits for a person to answer, and prints what they answered.
The skills use it to show you a review round or a design question.
You can also run it yourself:

```sh
socrates agora serve --document decision.json   # prints the page's address
socrates agora poll --session <session>         # blocks until somebody answers
socrates agora end --session <session>          # closes one page
socrates agora stop                             # shuts the server down
```

## design

`socrates design skill` installs two skills.
`socrates-design` interviews you over agora pages and writes a spec and its decisions.
`socrates-tickets` cuts a spec under `docs/changes/<slug>/` into tickets in `docs/changes/<slug>/tickets/`, with `README.md` as the main ticket that holds the order and status.
Nothing is published to a tracker, and nothing is committed for you.

## implement

`socrates implement skill` installs the `socrates-implement` skill and the `socrates-implementer` agent definition.
You invoke the skill once the tickets are cut and committed, and the session you invoked it in becomes the overseer.
She writes no code.
She hands each ticket whose blockers are done to an implementer in a worktree of its own, reviews what comes back, and lands it by fast-forward on a branch named `<type>/<slug>`, such as `feature/<slug>` or `fix/<slug>`, running the repository's own verification after every landing.
The implementers run in the background, so you can talk to her while they build, and she brings you only what coding cannot fix.
Once the last ticket lands she reads the whole change against the spec and tells you the branch is ready.
Nothing is pushed.

```sh
socrates implement config                    # every setting in force, and this repository's default branch
socrates implement configure                 # change one setting from a menu
socrates implement default-branch <branch>   # remember the branch this repository's features start from
```

`concurrency` is how many implementers run at once, 4 by default.
An implementer runs on Sonnet at high effort by default, and `roles.implementer.model` and `roles.implementer.effort` change that.
The settings live in `~/.socrates/implement/config.yaml`, and the overseer asks for the default branch once per repository and remembers it there.

## Reviewer CLIs

A CLI qualifies when:

- the reviewed repository's own instruction files can be kept out of the reviewer's session, so the change cannot instruct its own reviewer,
- it reports whether the work happened and why it failed, and
- a model can be chosen per call.

| reviewer CLI | how it produces the findings document | what keeps the repository's instructions out |
| --- | --- | --- |
| `claude` | the vendor constrains it against the schema | `--setting-sources user` |
| `codex` | the vendor constrains it against the schema | `-c project_doc_max_bytes=0` |
| `opencode` | socrates review extracts it from the final message | `OPENCODE_DISABLE_PROJECT_CONFIG=1` |

socrates validates the document either way.
The first two conform more often, so prefer them if you have the choice.

These are not supported yet:

| reviewer CLI | what is missing |
| --- | --- |
| `cursor-agent` | a way to keep the reviewed repository's own instructions out of the reviewer's session |
| `copilot` | a way to keep the reviewed repository's own instructions out of the reviewer's session, and a terminal event saying whether the work happened and what it says about why it did not |
| `pi` | a way to keep the reviewed repository's own instructions out of the reviewer's session, and a terminal event saying whether the work happened and what it says about why it did not |

That is what their documentation showed, not proof it cannot be done.
Supporting them is not planned.
Naming an unsupported CLI is refused before anything starts.

## Supported platforms

A review needs two guarantees: the reviewer cannot write, and everything it started is killed with it.

| platform | reaps by process group | the filesystem enforces the seal | the suite runs there |
| --- | --- | --- | --- |
| Linux | yes | yes | yes |
| macOS | yes | yes | yes |

On Windows, run socrates under WSL.
Elsewhere a review is refused, but reading a finished run still works.

## Configuration

socrates review reads one optional file, `~/.socrates/review/config.yaml`:

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
```

`socrates review config` prints every setting, its value and whether you set it.
`socrates review configure` changes one from a menu, and keeps your comments and ordering.
A `model` is passed to your CLI as written, so use that CLI's own model names.
A file that does not parse is refused rather than ignored.

socrates agora reads its own file, `~/.socrates/agora/config.yaml`, which holds only a `surface` block:

```yaml
surface:
  bind: 127.0.0.1:24672
  reached_at: http://localhost:24672/
```

A `surface` block left in review's file is refused.
Move it to agora's file.

Nothing is read from the repository under review except the review guide below.

## The review guide

The reviewer loads none of your project's `CLAUDE.md`, `AGENTS.md`, settings or skills, because the change under review could have written them.

To tell it what it would otherwise get wrong, commit `.socrates/review-guide.md`:

```sh
mkdir -p .socrates
$EDITOR .socrates/review-guide.md
git add .socrates/review-guide.md && git commit -m "what a reviewer should know about this repository"
```

- It is read at the base commit, so a guide added by the change only applies to the next review.
- It only adds context.
  It cannot turn off a rule, exempt a file or change what `clean` means.
- It must fit under `guidance_ceiling` (16 KiB by default), or it is left out and the review says so.

## Reading the page from another machine

The page is served on `127.0.0.1:24672` and only opens on that machine.
If something else is already listening on that port, socrates agora refuses and starts no page anywhere else.
To move the page, state an address with `surface.bind` in agora's file or with `--address`, and a taken port you state is refused the same way.

In a devcontainer, over ssh, or anywhere with its own network namespace, your browser cannot reach loopback, and it sees an empty reply.
Forward the port from inside that namespace:

- a devcontainer: `{ "forwardPorts": [24672] }` in `devcontainer.json`
- ssh: `ssh -L 24672:127.0.0.1:24672 <host>`
- or a reverse proxy running beside it

Publishing the port in a compose file does not work, because it forwards to the container's network interface, not to loopback.

Then set the address people should open:

```yaml
surface:
  reached_at: http://socrates-review.localhost/
```

## Your repository and your data

socrates review writes nothing in your repository except the git objects for its snapshot.
Your index is untouched, and `git status` is the same before and after.
Runs live in `~/.socrates/review/runs/`, one directory per run.
The pages socrates agora serves are written under `~/.socrates/agora/`, beside its own configuration file.
`SOCRATES_HOME` moves `~/.socrates`.

socrates review sends nothing anywhere: no telemetry, no update check, no network connection.
The reviewer CLI it runs sends the change to its vendor, exactly as it does when you use it yourself: `claude` to Anthropic, `codex` to OpenAI, `opencode` to whichever provider you set up.
That includes the diff, the files the reviewer reads and your `--intent`.
If that is not acceptable for a repository, do not review it with socrates.

socrates is not affiliated with, endorsed by or sponsored by any of these vendors.
"Claude" and "Anthropic" are trademarks of Anthropic PBC, and "Codex" and "OpenAI" of OpenAI.

## Out of scope

- **It does not land your change.** No commit, push, pull request or CI.
- **It does not fetch, rebase or post comments.** It talks to no forge.
- **It does not approve anything.**
- **It does not run your tests or linters.** `clean` means the reviewer found nothing blocking, not that the build is green.
- **It does not write in your repository.**
  The reviewer only gets read-only tools, the snapshot it reads is made unwritable, and the round fails if the snapshot changed by a single byte.
  A review is refused where the unwritable snapshot cannot be enforced, such as when running as root.
- **It does not load your customisations.**
  Plugins, memory files, skills and hooks are switched off with the vendor's own flag, and any tool the session reports anyway is recorded on the run.
  A hook your administrator's policy configures stays on.
  It runs inside the sealed snapshot, so a write there is still caught, but nothing stops it reaching outside the snapshot.
- **It does not tell the reviewer how the change came to be.** Only the change and your `--intent`.
- **It is not a daemon.** One process per round and a directory of files.

`SECURITY.md` says what socrates review defends against and how to report a problem.

## Licence

MIT, stated in `LICENSE`.
`NOTICE` lists everything a shipped binary contains that socrates did not write, with each licence in full.
`socrates licences` prints both from the binary itself, so a copy installed with `go install` carries them too.

## Contributing

Read `CONTRIBUTING.md`.
`AGENTS.md` maps the layout and invariants, `CONTEXT-MAP.md` the vocabulary, and `docs/adr/` the decisions, each at the root and in each tool's directory.
