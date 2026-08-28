## Context

See proposal.md for why.

What exists today, settled by `the-one-place-a-program-writes` and unchanged since.

`internal/state` resolves `SOCRATES_HOME` or `~/.socrates`, holding `config.yaml` and `runs/`.
`internal/surface` resolves `AGORA_HOME` or `~/.agora` by the same rule written a second time, holding `config.yaml`, `surface.json` and `sessions/<session>/`.
`internal/config` reads socrates' configuration file; `internal/surface` reads agora's, with a loader of its own that differs only in the shape it decodes into.

That change listed a shared home as an explicit non-goal, and `TestNeitherProgramWritesInTheOthersHome` holds the boundary.
This change reverses that one non-goal and nothing else.

## Goals / Non-Goals

**Goals:**

- One home, one variable and one configuration file for the whole tool.
- The resolution rule stated once, as it already is, and reached by both binaries.
- A person handed one address for a whole review rather than one per round.

**Non-Goals:**

- Merging the binaries. `socrates` and `agora` stay two programs built from one module, for the reasons `One binary, five commands, beside socrates` already gives.
- Changing how a session comes into existence. agora mints it, prints it, and `--session` names one it holds.
- socrates learning anything about a page. No command, no flag, no run-record field, no network reach.
- Teaching agora what a review is. It serves a document to a person and knows nothing of findings, rounds or verdicts, and `WordsOfTheOtherTool` still holds that.
- Migrating what a person already has under `~/.agora`.

## Decisions

### The surface's state is a directory under socrates' home, not a second home

`~/.socrates/surface/` holds `surface.json` and `sessions/<session>/`, beside `runs/`.
`internal/surface` asks `internal/state` where that is, so the rule stays stated once, which is what `the-one-place-a-program-writes` was for and what a second resolver quietly undid.

`~/.socrates/agora/` was considered and rejected.
A person opens that directory to find out what is going on, and `surface/` says what is in it to somebody who has never heard of agora, which is most people reading it.

agora keeps its own name in what it prints about itself, because a person running `agora version` is asking about the binary in their hand rather than about the tool it serves.

### The configuration is one file with a block, not two files

An operator states two things about the surface, `bind` and `reached_at`, and stating them in socrates' `config.yaml` under `surface:` is what "one file" means.
`Configuration is one optional user file` says the file holds a stated list "and nothing else", and that list grows by one block rather than the sentence going.

`internal/config` stays the one owner of the file and exposes the block; `internal/surface` reads it from there rather than parsing anything itself, so `KnownFields` and the refusal that names an unknown field are written once.

A second file at `~/.socrates/surface/config.yaml` was considered and rejected: it is the two-answers problem again, one directory deeper.

### The old variable is refused, not ignored

An `AGORA_HOME` in somebody's environment is a person who believes they have moved the surface's state.
Ignoring it is the failure `the-one-place-a-program-writes` refused to ship for a relative home, in a different costume: state that is not where the person was told to expect it, with nothing saying so.

So `agora` refuses at startup while it is set, names it, and names `SOCRATES_HOME`.
The refusal is cheap to remove later and expensive to have needed and not had.

### The session stays agora's to mint, and the driver's to remember

A review is two rounds and one person reading them, so the person should be handed one address and see the second round arrive on the page they already have.
Nothing about minting has to change for that: `agora serve` mints on the first round, `agora serve --session` writes the second round to the same page, and both already work.
What was missing is that nothing told the session driving the loop to keep the identifier, so it handed over a second address on the second round.
That is a line in the skill.

Considered and rejected: recording the session on the run so it survives the driver being replaced.
It buys a real guarantee, and it costs a flag on `socrates` and a field about a page on a run's record, both of which `socrates states the fact and does nothing further about the page` forbids in as many words.
The failure it prevents is a driver compacted or restarted mid-review, which costs the person one stale tab and one new address, and that is not worth reversing a stated rule for.
If it ever bites often enough to notice, it is its own change and this one does not stand in its way.

## Risks / Trade-offs

- **A person with state under `~/.agora` loses it** → there is no migration, exactly as there was none last time. What is lost is a surface record and any page still open, and stopping the old binary would have lost both anyway.
- **agora stops being usable on its own terms** → it still runs, still serves any document to any person, and still knows nothing of review. What it loses is a home of its own, which is a decision about this repository rather than about the program.
- **A driver replaced mid-review still hands over a second address** → accepted above, and the person's first page says the exchange is over rather than going silently stale.
- **One stated non-goal is reversed** → it was right about the rule and wrong about the count, and reversing it in a change that says so is cheaper than a repository whose specification argues with itself.

## Migration Plan

None.

`~/.agora` is left where it is and read by nothing.
A person who wants it gone removes it; a person who does not notice loses nothing they were using.
