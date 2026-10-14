# socrates

[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/socrates-tools/socrates/badge)](https://scorecard.dev/viewer/?uri=github.com/socrates-tools/socrates)

socrates takes a change from an idea to a reviewed branch.
In your coding session, run these skills in order:

1. `/socrates-design` interviews you in the browser and writes a spec.
2. `/socrates-tickets` cuts that spec into tickets.
3. `/socrates-implement` builds the tickets on a branch.
4. `/socrates-review` has an independent reviewer check the change before you push.

Each step works on its own, so you can review any change without a spec.

## What it needs

- **Claude Code**, which runs the skills.
- **For review, `claude`, `codex` or `opencode` on your `PATH`, signed in to an account.**
  Review uses `claude` unless you set `reviewer_cli`, and `claude` must be 2.1.169 or later.
  socrates is free, but every review is a model call billed to that account.
- **Linux or macOS.**
  On Windows, use WSL.
- **`git`.**
- **Go 1.26 or later**, only to build from source.

## Installing

```sh
curl -fsSL https://github.com/socrates-tools/socrates/releases/latest/download/install.sh | sh
```

That puts `socrates` in `~/.local/bin` and installs every skill and agent definition under `~/.claude/`.
It checks the archive against the release's checksums, and the checksums against the release's signature where `cosign` is on your `PATH`, and installs nothing if a check fails.
Running it again upgrades everything.

`SOCRATES_INSTALL_DIR` puts the binary somewhere else, and `SOCRATES_VERSION` installs another release than the latest:

```sh
curl -fsSL https://github.com/socrates-tools/socrates/releases/latest/download/install.sh | SOCRATES_INSTALL_DIR=/usr/local/bin SOCRATES_VERSION=v0.2.0 sh
```

From a clone, `make install` does the same from source.

`socrates version` says which release and commit you have, and `CHANGELOG.md` says what changed in each release.

## Design

Run `/socrates-design` and describe what you want.
It asks its questions on a page in your browser and writes the spec to `docs/changes/<slug>/spec.md`, with the decisions behind it.
Read it and commit it.

## Tickets

Run `/socrates-tickets` once the spec is committed.
It cuts the spec into tickets in `docs/changes/<slug>/tickets/`, with `README.md` holding their order and status.
Commit them.

## Implement

Run `/socrates-implement` once the tickets are committed.
It builds each ticket in a worktree of its own, reviews it, lands it on a branch such as `feature/<slug>`, and runs your repository's verification after every landing.
Keep talking to it while it works: it only asks you what code cannot settle.
Nothing is pushed.

## Review

Run `/socrates-review` before you push.
A reviewer that did not write the change reads it, and your session fixes what it can and brings you the rest.
For somebody else's change, check it out in its own worktree and run `/socrates-review-handover`, which posts only the findings you choose.

The reviewer reads none of your `CLAUDE.md`, `AGENTS.md`, settings or skills, because the change could have written them.
To tell it what it should know about your repository, commit `.socrates/review-guide.md`.

### Which CLI reviews

Set `reviewer_cli` to one of these:

| reviewer CLI | how it produces the findings document | what keeps the repository's instructions out |
| --- | --- | --- |
| `claude` | the vendor constrains it against the schema | `--setting-sources user` |
| `codex` | the vendor constrains it against the schema | `-c project_doc_max_bytes=0` |
| `opencode` | socrates review extracts it from the final message | `OPENCODE_DISABLE_PROJECT_CONFIG=1` |

socrates validates the document either way.
The first two conform more often, so prefer them if you have the choice.

## Configuration

```sh
socrates review configure      # reviewer CLI, model, effort and limits
socrates implement configure   # how many implementers run at once, and their model and effort
```

`socrates review config` and `socrates implement config` print every setting in force.
Settings live in `~/.socrates/<tool>/config.yaml`, and `SOCRATES_HOME` moves `~/.socrates`.

## Reading the page from another machine

The page is served on `127.0.0.1:24672`.
From a devcontainer or over ssh, forward that port:

- a devcontainer: `{ "forwardPorts": [24672] }` in `devcontainer.json`
- ssh: `ssh -L 24672:127.0.0.1:24672 <host>`

Behind a reverse proxy, set the address people should open in `~/.socrates/agora/config.yaml`:

```yaml
surface:
  reached_at: http://socrates-review.localhost/
```

## Supported platforms

| platform | supported |
| --- | --- |
| Linux | yes |
| macOS | yes |
| Windows | under WSL |

## Your code and your data

socrates review sends nothing anywhere: no telemetry, no update check, no network connection.
The reviewer CLI it runs sends the change to its vendor, exactly as it does when you use it yourself.
Runs and pages are kept under `~/.socrates`, and nothing is written in your repository except the git objects for a review's snapshot.

socrates is not affiliated with, endorsed by or sponsored by Anthropic or OpenAI.
"Claude" and "Anthropic" are trademarks of Anthropic PBC, and "Codex" and "OpenAI" of OpenAI.

## Out of scope

- **It does not push, open pull requests or run CI.**
- **A hook your administrator's policy configures stays on in the reviewer's session**, and nothing stops it reaching outside the snapshot.

`SECURITY.md` says what socrates review defends against and how to report a problem.

## Licence

MIT, stated in `LICENSE`.
`NOTICE` lists everything a shipped binary contains that socrates did not write, with each licence in full.
`socrates licences` prints both from the binary itself.

## Contributing

Read `CONTRIBUTING.md`.
`AGENTS.md` maps the layout and invariants, `CONTEXT-MAP.md` the vocabulary, and `docs/adr/` the decisions, each at the root and in each tool's directory.
