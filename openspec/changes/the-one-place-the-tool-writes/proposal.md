## Why

`the-one-place-a-program-writes` gave each program one home and one variable, and made a shared home an explicit non-goal.
socrates keeps everything under `~/.socrates`, moved by `SOCRATES_HOME`; agora keeps everything under `~/.agora`, moved by `AGORA_HOME`, with a configuration file of its own beside socrates' one.

That was right about the rule and wrong about the count.
It reasoned about two programs, and this repository does not hold two programs anybody uses separately: it holds one tool for having somebody who did not write a change read it, and a surface that puts one round of that review in front of a person.
Nothing else calls agora and nothing else is going to.

So the question the earlier change existed to make answerable in one line, "where is your state?", still has two answers.
A person who wants to move it moves it twice.
A person configuring it opens two files.
A stranger reading an issue that will not reproduce is told about a second directory belonging to a program they have never heard of.

One tool, one home.

## What Changes

- **BREAKING**: `AGORA_HOME` is removed.
  `SOCRATES_HOME` moves everything both binaries read and write, and where none is stated everything is under `~/.socrates`.
  The surface's state moves from `~/.agora` to `~/.socrates/surface/`.
- **BREAKING**: agora's configuration file is removed as a file of its own.
  What an operator states about the bind and the address the surface is reached at moves into a `surface` block in socrates' one configuration file.
- `agora` refuses at startup while `AGORA_HOME` is still set, naming it and naming what replaced it, rather than ignoring it.
- The skill keeps the session `agora serve` printed and writes the second round of a review to that same page, so a person is handed one address for the whole review.

Nothing about how a session comes into existence changes.
`agora serve` given no session still mints one and prints an address holding it, `--session` still names one the surface already holds, and socrates still names no page, records no page and reaches no network.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `review-surface`: `Session state lives outside the repository being discussed` puts the surface's state under socrates' home and drops agora's own variable and configuration file.
- `cli-foundation`: `Configuration is one optional user file` gains the surface block, so the one file is the one file for both binaries.
- `agent-skill`: `The skill renders a round a person reads as a page` says to keep the session across the rounds of one review rather than handing over a second address.

## Impact

- `internal/state`: gains the surface's directory beside `runs/`, and stays the one resolver.
- `internal/surface`: `Home` resolves through `internal/state` rather than a variable of its own, and its configuration loader goes.
- `internal/config`: the one configuration file gains a surface block, and `internal/surface` reads it from there.
- `cmd/agora`: refuses while the old variable is set.
- `internal/skill`: the session survives from round one to round two.
- `cmd/socrates` and every package it reaches: untouched. No command, no flag, no run-record field, no network.
- No new dependency.
