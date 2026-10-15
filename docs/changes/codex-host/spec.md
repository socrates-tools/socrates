# Codex is a host for every socrates skill

Decisions: ADR 66, ADR 67, ADR 68, ADR 69.

## Problem Statement

A person who works in Codex cannot use socrates.
socrates review can already spawn `codex` as its reviewer CLI, but the five skills that drive socrates (`socrates-review`, `socrates-review-handover`, `socrates-design`, `socrates-tickets` and `socrates-implement`) are written for Claude Code alone.
They install under `~/.claude`, mark themselves person-only with a front matter field only Claude Code reads, and describe subagents, background work and asking the person in Claude Code's terms.
`socrates-implement` also writes a Claude Code agent definition for its implementers, which Codex cannot read.
Under Codex's default sandbox, `socrates agora serve` cannot bind its loopback port, and an implementer cannot write its worktree under `~/.socrates/implement`.

## Solution

socrates supports two hosts, Claude Code and Codex, and every skill is rendered once per host from one source, per ADR 66.
Every install writes both hosts' renderings, Claude Code's under `~/.claude` and Codex's under `~/.agents/skills` and `$CODEX_HOME/agents`.

On Codex:

- A skill is person-only through `agents/openai.yaml` beside its `SKILL.md`, holding `policy.allow_implicit_invocation: false`, per ADR 67.
- A skill that serves a page tells the session to run `socrates agora serve` with escalated permissions, so Codex asks the person before the bind.
- The implementer is a Codex role file, `socrates-implementer.toml`, with a model and an effort of its own, per ADR 68.
- The overseer checks it can write `~/.socrates/implement` before a run and, where it cannot, names the line that grants it and starts nothing.
- `socrates review begin --host codex` makes `codex` the reviewer CLI where `reviewer_cli` is unset, per ADR 69.

## User Stories

1. As a person who works in Codex, I want `socrates skill` to install every socrates skill where Codex finds it, so that I can start any of them by typing `$socrates-design`, `$socrates-tickets`, `$socrates-implement`, `$socrates-review` or `$socrates-review-handover`.
2. As a person who works in Claude Code, I want everything I have today to keep working the same, so that Codex support costs me nothing but two directories I never look at.
3. As a person who uses both hosts, I want one install to leave both renderings in place, so that I can start a design in one host and implement it in the other.
4. As a person who installs Codex after socrates, I want the skills already there, so that I do not have to remember a second install.
5. As a person in Codex, I want no socrates skill offered to the model on its own, so that a session never starts a design, a review or an implementation I did not ask for.
6. As a person in Codex, I want each skill to tell Codex exactly what to do in Codex's own terms, so that the session does not have to translate Claude Code instructions and get it wrong.
7. As a person running a design or a review in Codex, I want Codex to ask me before it runs `socrates agora serve` outside the sandbox, so that the page opens without me changing my Codex sandbox settings and nothing runs unsandboxed without my say.
8. As a person starting `socrates-implement` in Codex whose sandbox cannot write `~/.socrates/implement`, I want the overseer to stop before handing anything out and tell me the one line that fixes it, so that a run does not fail halfway through at its first worktree.
9. As a person implementing in Codex, I want the implementers to run as Codex subagents, each in a worktree the overseer made, while the overseer stays free to talk to me, so that a run in Codex works as it does in Claude Code.
10. As a person implementing in Codex, I want to choose the implementers' Codex model and effort with `socrates implement configure`, separately from the Claude Code ones, so that each host runs on values it accepts.
11. As a person who never set a Codex model, I want the implementers to run on my Codex session's own model, so that I get something sensible without a setting.
12. As a person who reads `socrates implement config`, I want both hosts' implementer settings printed with the value in force and its default, so that I can see what the next implementer runs on in either host.
13. As a person running a review from Codex with no `reviewer_cli` set, I want `codex` to review the change, so that socrates review works on a machine with only Codex on it.
14. As a person who set `reviewer_cli`, I want that reviewer whichever host I start the review from, so that a choice I made once is never overridden.
15. As a person whose review runs several rounds, I want every round reviewed by the CLI that reviewed the first, so that a run never changes vendor halfway through.
16. As a person who installs a skill somewhere other than its default, I want `--to` to say which host's rendering it writes, and to be refused without `--host`, so that I never install a rendering a host cannot read.
17. As a person whose Claude Code reads its configuration from `CLAUDE_CONFIG_DIR`, I want every install to write Claude Code's skills and agent definition there, so that what I install is what my sessions load.
18. As a contributor, I want the suite to fail where a host's rendering names a command the binary does not offer, outgrows its ceiling, loses its person-only switch, or no longer matches its committed copy, so that the Codex renderings are held to everything the Claude Code ones are.

## Implementation Decisions

**The host**

- `internal/skill` states the supported hosts once: Claude Code and Codex, each with its skills directory under the home directory (`.claude/skills` and `.agents/skills`), how a person starts a skill (`/<name>` and `$<name>`), and its person-only switch.
- `skill.Skill` gains the host it is rendered for, and a skill's `Document` becomes a function of the host.
  Each namespace's `All` returns one rendering per host.
- `internal/skill.FrontMatter` takes the host.
  Claude Code's front matter carries `disable-model-invocation: true` as today.
  Codex's front matter carries `name` and `description` only.
- `internal/skill.Install` writes `agents/openai.yaml`, holding `policy:` and `allow_implicit_invocation: false`, beside every Codex `SKILL.md`.
- `internal/skill.UserDirectory` takes the host, and every namespace's `skill` writes every host's rendering to that host's directory.
- Each host's directory follows the variable that host reads for its own configuration: Claude Code's skills and agents go under `$CLAUDE_CONFIG_DIR` where it is set and `~/.claude` where it is not, as Codex's agents go under `$CODEX_HOME`.
- `--to` and `--agents-to` need `--host claude` or `--host codex`.
  Either flag without `--host` is refused, naming `--host`, before anything is written.
  `--host` without `--to` writes only that host's rendering to its default directory.
- `socrates skill` is unchanged and runs every namespace's `skill` with its defaults, so it writes every host.

**What a Codex rendering says differently**

- The "Installing elsewhere" section names Codex's directories and `--host codex`.
- Where a skill tells the person which skill comes next, it names it as `$<name>`.
- Every skill that runs `socrates agora serve` tells the Codex session to run it with `sandbox_permissions` set to `require_escalated` and a one-line justification naming the page.
  That covers `socrates-design`, `socrates-review` and `socrates-review-handover`.
  `socrates agora poll` and `socrates agora end` reach the surface over loopback as clients, and they are escalated the same way.
- Where a Claude Code rendering says to run work in the background, the Codex rendering says to start it with `spawn_agent` and collect it with `wait_agent`, keeping waits short enough to answer the person between them.
- Every Codex rendering keeps the sentence telling the session to invoke no skill, which is what holds ADR 49 on Codex.

**socrates-implement on Codex**

- `implement/skill` renders the implementer twice: `socrates-implementer.md` for Claude Code as today, and `socrates-implementer.toml` for Codex.
  The TOML carries `name`, `description`, `developer_instructions` (the same standing instructions), and `model` and `model_reasoning_effort` where set, and no other key, because Codex refuses a role file holding a key it does not know.
- The Codex agents directory is `$CODEX_HOME/agents`, `~/.codex/agents` where `CODEX_HOME` is unset, and it is created where it does not exist.
- The Codex overseer spawns each implementer with `spawn_agent` naming `agent_type` `socrates-implementer`, sends feedback to it with `send_input`, and closes it with `close_agent` once its ticket has landed, because an open agent counts against Codex's concurrency limit.
- Where `socrates-implementer` is not a role it can spawn, the Codex overseer tells the person to run `socrates implement skill` and starts nothing, as the Claude Code one does.
- At the start of every run the Codex overseer checks it can write `~/.socrates/implement`, the directory beside the file `socrates implement config` prints on its first line.
  Where it cannot, it starts nothing and tells the person to add that directory to `writable_roots` under `[sandbox_workspace_write]` in Codex's configuration, or to start Codex with `--add-dir <that directory>`.
- `implement/config` gains `roles.implementer.codex.model`, unset by default, and `roles.implementer.codex.effort`, `high` by default.
  Codex's efforts are Codex's own: `none`, `minimal`, `low`, `medium`, `high`, `xhigh`, `max` and `ultra`.
  An unset Codex model writes no `model` line, so the implementer runs on the session's model.
- `socrates implement config` reports both Codex settings beside the Claude Code ones.
- `socrates implement configure` offers a Codex model row, taking `use the default` or a typed identifier, and a Codex effort row offering Codex's efforts.
  A confirmed value rewrites the installed Codex role file at once.
- `~/.socrates/implement/installed.yaml` records one agents directory per host.
  A record in today's shape, holding one `agents_directory`, is read as Claude Code's.

**socrates review**

- `socrates review begin` takes `--host`, `claude` or `codex`.
  Where `reviewer_cli` is unset, the host's vendor is the reviewer CLI, and without `--host` the default stays `claude`.
  `review/agent` states which vendor each host maps to, so nothing above the seam names a vendor.
- The vendor is resolved once at `begin` and recorded on the run, so `continue` uses it rather than resolving again.
- `socrates review config` still reports `reviewer_cli` as the file states it, and says that an unset one follows the host a review begins from.
- Each rendering of `socrates-review` and `socrates-review-handover` passes its own `--host` to `begin`.

**Documents**

- `README.md` names Claude Code or Codex as what runs the skills, where each host's skills are installed, and the Codex writable root `socrates-implement` needs.
- The root `CLAUDE.md`, `review/CLAUDE.md` and `implement/CLAUDE.md` state the new invariants: every skill per host, the person-only switch per host, the Codex role file and its settings, and `--host` on `begin`.
- `make skill` writes the Codex renderings beside the Claude Code ones: `.agents/skills/<name>/` for the skills and `.codex/agents/socrates-implementer.toml` for the role.
- ADR 49 and ADR 60 each gain a line naming the ADR that amends them.

## Testing Decisions

- Every skill test that reads a rendered document runs once per host: the ceiling, the commands it names against the real binary, the prohibitions, and the committed copy matching what the package renders.
  This follows `design/skill`, `implement/skill` and `review/skill` as they are, with the host as a parameter rather than a second copy of each test.
- The person-only check in `internal/skill` reads each host's own switch, over every committed skill and over everything `socrates <namespace> skill` installs, per ADR 49 as amended.
  It is mutated twice: once by dropping `disable-model-invocation` from a Claude Code rendering, and once by dropping `agents/openai.yaml` from a Codex one.
- A test reads the rendered Codex role file as TOML and fails on any key outside the five Codex accepts.
- End-to-end tests run `socrates skill` under a temporary home with no `codex` on `PATH`, and assert both hosts' trees are written, including the `openai.yaml` files and the role file under `$CODEX_HOME/agents`.
- End-to-end tests assert that `--to` and `--agents-to` without `--host` are refused before anything is written, and that `--host codex --to <dir>` writes the Codex rendering there.
- `socrates implement configure` is driven in a real pseudo-terminal: confirming a Codex model and a Codex effort rewrites the installed role file, and the Claude Code definition is left alone.
- An `installed.yaml` in today's shape is read as Claude Code's record.
- socrates review end-to-end tests run against each simulator on `PATH`.
  `begin --host codex` with `reviewer_cli` unset spawns `codexsim`.
  `begin --host codex` with `reviewer_cli: claude` spawns `claudesim`.
  A `continue` after `begin --host codex` spawns `codexsim` again, whatever `--host` it is given.
- What a Codex session does with a rendering is the session's to perform, as for Claude Code: the suite asserts each promise as a statement in the rendered document, and the orchestration is not simulated.
- One live run per skill under a real `codex` is recorded once the change is built: a design over agora, a review, and an implementation of one ticket.

## Out of Scope

- Any host beyond Claude Code and Codex.
- Codex plugins and the Codex marketplace.
- Installing into a repository's own `.agents/skills` or `.codex/agents`, beyond the copies `make skill` commits.
- Reading what model the calling session runs, which ADR 60 still rules out.

## Further Notes

- Codex was not installed in the container this design ran in.
  The loopback bind being denied under Codex's sandbox, the escalation prompt, subagents inheriting `--add-dir`, and a role naming a model the account cannot use were all read from Codex's source and documentation rather than seen.
  The live runs under Testing Decisions are where each is first seen.
- Codex's documentation moved from developers.openai.com/codex to learn.chatgpt.com/docs.
  Codex still reads `$CODEX_HOME/skills` but calls it deprecated, which is why Codex's rendering goes to `~/.agents/skills`.
- While the first round was being served, the page was not reachable from a Windows browser through the devcontainer's reverse proxy.
  agora refuses a request whose `Host` header is not a loopback address on the port it bound, so a proxy must rewrite the header as well as the address that `surface.reached_at` names.
  The README's section on reading the page from another machine does not say so, which is a documentation fix of its own rather than part of this change.
