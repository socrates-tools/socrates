# socrates implement

`socrates implement` builds a change's tickets through implementers, per ADR 48.
`socrates implement skill` installs `socrates-implement`, which turns the session that a person invokes it in into the overseer, and the `socrates-implementer` agent definition the overseer spawns.
It reads one configuration file of its own, `~/.socrates/implement/config.yaml`, `socrates implement config` reports every setting in it, and `socrates implement configure` changes one from a menu.
It knows nothing about review, agora or design.

This file is implement's layout and invariants.
`/AGENTS.md` holds the suite's, `CONTEXT.md` here holds implement's vocabulary, and `docs/adr/` here holds implement's decisions.

## Layout

- `implement/command` - the command tree: `config`, `configure`, `default-branch` and `skill`, which installs the overseer's skill and the implementer's agent definition.
- `implement/configure` - the menu over implement's settings: its labels and its refusals, the model aliases it offers, the default branch row inside a repository, and rewriting the installed agent definition after a model or an effort is confirmed, built on `internal/settingsmenu`.
- `implement/skill` - the source of the `socrates-implement` skill and of the implementer's agent definition and Codex role, installing them, and the record of where the definition was installed.
- `implement/config` - the one optional configuration file implement reads, every setting it holds, the default each falls back to, the default branch each repository's features start from and how a repository is named, by the top level of its main worktree, and the report `socrates implement config` prints.

## Invariants

Each is held by a test that fails when it stops holding.
None is a matter of care, and the reasoning for each is in this directory's `docs/adr/`.

**The implement skill**

- `socrates-implement` is rendered by `implement/skill` once per host, installed by `socrates implement skill` under each host's skills directory or under `--to` with `--host`, is person-only by each host's own switch through `internal/skill`, and is held to a ceiling of its own.
- The committed copies at `.claude/skills/socrates-implement/SKILL.md` and `.agents/skills/socrates-implement/SKILL.md` must match what the package renders, and `make skill` writes them.
- Every invocation it names exists in the binary with the flags it names, and it names only `socrates implement` ones.
- It names no `socrates review` invocation and no other skill, and tells the overseer to invoke no skill, per ADR 49.
- The overseer writes no code, edits only the main ticket's status table, and settles every coding question by her own judgement.
- She reads the main ticket, the spec and the ADRs, and the settings from `socrates implement config`, at the start of every run.
- A default branch that is not remembered is proposed from the branches named `main`, `master`, `dev` or `develop`, never one under a prefix such as `feature/` or `hotfix/`, asked in the terminal and saved with `socrates implement default-branch`, per ADR 51.
- The change's branch is cut from the current `HEAD` and starts with what the design and the tickets left, and anything else uncommitted stops the run.
- The change's branch is `<type>/<slug>`, its type one of `feature`, `fix`, `refactor`, `docs` or `chore` as the spec says the change is, each ticket's branch is `<type>/<slug>-<ticket>` beside it, and a resumed run keeps the branch it finds under any type.
- The frontier is every ticket whose blockers are done, at most `concurrency` implementers run, the longest chain goes first, then the most unblocked, then name order, and a ready ticket sharing a surface with a running one is held back without a blocker being written.
- Each implementer is briefed with pointers, runs in the background in a worktree and on a branch the overseer made with git, and gets feedback from the overseer with no cap on the exchanges.
- Before the first hand-out one subagent writes notes about the code the change touches under `~/.socrates/implement/notes/`, outside the repository, every brief points at them, and the overseer writes none herself.
- Review keeps what the ticket, the spec and the ADRs ask apart from the repository's own standards, quoting the line or citing the rule missed, and a ticket that needs a person is asked about before it is handed out.
- Landing is a fast-forward followed by the repository's own verification, a failure goes to the implementer whose ticket landed last, and where the only commits an implementer's branch lacks are the overseer's own status commits she merges them into its branch herself, while any other new commit goes back to the implementer to merge.
- In progress is committed before a hand-out and done once a ticket has landed and passed, and a resumed overseer continues each in-progress ticket from its own branch with a fresh implementer.
- What the person says mid-run is relayed at once to every implementer it concerns, and the person is reached in the terminal only about what coding cannot fix.
- The whole change is read against the spec from the default branch's fork point, story by story, every worktree is removed, and nothing is pushed.
- Every promise above holds in each host's rendering, and what differs is stated in the Codex rendering's `On Codex` section and in the three sentences that say how an implementer or the verification runs: Claude Code runs them in the background, as the `socrates-implementer` agent, and Codex starts an implementer with `spawn_agent` naming `agent_type` `socrates-implementer` and has a subagent run the verification.
- The Codex overseer sends feedback with `send_input`, collects reports with `wait_agent` between answering the person, closes each implementer with `close_agent` once its ticket has landed, and tells the person to run `socrates implement skill` where `socrates-implementer` is not a role it can spawn.
- The Codex overseer checks at the start of every run that it can write the directory beside the file `socrates implement config` prints on its first line, and where it cannot starts nothing and names `writable_roots` under `[sandbox_workspace_write]` and `--add-dir`, per ADR 68.
- Each host's rendering has a ceiling of its own, `OverseerCeilings`.
- The orchestration itself is the session's to perform from what the skill says: the suite asserts each promise as a statement in the rendered document and every command the skill names against the real binary.

**The implementer's agent definition**

- `implement/skill` renders `socrates-implementer.md` once, and the committed copy at `.claude/agents/socrates-implementer.md` must match what it renders at the defaults.
- Its `model` and `effort` are the settings in force, and `background` is true.
- It names no `socrates` command, read against every namespace in `convention.Tools`.
- It never uses `git stash`, and where the repository states no way of testing it writes each test before the code and watches it fail.
- It is held to a ceiling of its own.
- Where it was written is recorded in `~/.socrates/implement/installed.yaml`, one directory per host, a second install elsewhere for a host replaces that host's entry only, a record holding one `agents_directory` is read as Claude Code's, and `config.yaml` never holds it.

**The implementer's Codex role**

- `implement/skill` renders `socrates-implementer.toml` beside the definition, and the committed copy at `.codex/agents/socrates-implementer.toml` must match what it renders at the defaults.
- It holds `name`, `description`, `developer_instructions`, `model_reasoning_effort` and, where a model is set, `model`, and no other key, read as TOML, because Codex refuses a role file holding a key it does not know.
- Its instructions are the definition's, it names no `socrates` command, and it is held to a ceiling of its own.
- It installs under `$CODEX_HOME/agents`, `~/.codex/agents` where `CODEX_HOME` is unset, creating the directory, and `--host codex --agents-to <directory>` writes it there.

**The default branch**

- A repository is keyed by the top level of its main worktree, so every worktree of one clone reads one entry and a second clone at another path reads none.
- A branch the repository does not hold, local or remote-tracking, is refused by `socrates implement default-branch` and the file is untouched.
- Outside a repository and in a bare one, `default-branch` is refused naming what it needs, and `config` reports no `default_branch`.
- A save keeps a person's comments, key order and blank lines, through the same machinery as the menu, and a repository path holding a dot is one key.

**The settings**

- `socrates implement configure` draws the menu `internal/settingsmenu` draws for every tool, and refuses where nothing is a terminal, naming `socrates implement config`.
- `socrates implement configure` offers `sonnet`, `opus`, `haiku`, `fable` and `inherit` for `roles.implementer.model` and also takes a full identifier typed.
- A confirmed model or effort rewrites the installed agent definition at once, and with none installed the menu names `socrates implement skill` as what writes it.
- Inside a repository the menu offers that repository's `default_branch`, a branch the repository does not hold is refused before the file is touched, and `use the default` removes that repository's entry and no other.
