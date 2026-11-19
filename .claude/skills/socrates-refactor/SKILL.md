---
name: socrates-refactor
description: Run improve-codebase-architecture with its report and the pick of a candidate on a lavish page a person annotates, installing either skill where it is missing. Use when a codebase needs its deepening opportunities found and grilled.
disable-model-invocation: true
---

# socrates refactor

This skill composes two skills other people publish, so that an architecture review is read and answered on a page rather than in a terminal:

- `improve-codebase-architecture` finds deepening opportunities and grills through the one the person picks.
- `lavish` turns an HTML page into one the person annotates and sends back.

It writes nothing to the repository but what the grilling settles: glossary entries, ADRs, and the refactor the person asks for.

## The skills it composes

- `improve-codebase-architecture`, from `mattpocock/skills`, and the three it calls: `codebase-design`, `grilling` and `domain-modeling`.
- `lavish`, from `kunchenguid/lavish-axi`, which drives the `lavish-axi` CLI.

## Installing what is missing

Before anything else, look for each of the five as a directory holding a `SKILL.md`, first under the repository's own `.claude/skills` and then under `$CLAUDE_CONFIG_DIR/skills`, or `~/.claude/skills` where it is unset.
Where any of the four from `mattpocock/skills` is missing, install all four, and where `lavish` is missing, install it:

    npx -y skills add mattpocock/skills --skill improve-codebase-architecture codebase-design grilling domain-modeling -g -a claude-code -y
    npx -y skills add kunchenguid/lavish-axi --skill lavish -g -a claude-code -y

Say in one line what was installed, and read each from where it now is.
Where an install fails, stop, and hand the person the command that failed with what it printed.
Install nothing that is already there.

## Running it

`improve-codebase-architecture` is a person's to start, so no session can invoke it by name.
Read its `SKILL.md` and its `HTML-REPORT.md` from where it is installed, and follow its process with three changes.
Read `lavish`'s `SKILL.md` too, and take its current guidance from the CLI it names, `npx -y lavish-axi --help` and `npx -y lavish-axi design`, over anything written here.

1. The report is a lavish page.
   Write it as `HTML-REPORT.md` describes each candidate card and the top recommendation, styled as `npx -y lavish-axi design` directs, to the OS temp directory and never into the repository.
   Open it with `npx -y lavish-axi <file>` rather than `open`, and hand the person what it prints.
2. The pick is made on the page, never in the terminal.
   Where the skill says to ask "Which of these would you like to explore?", ask it on the page, wait for the person's annotations through lavish's review loop, and read the candidate they picked from what they sent.
3. The grilling stays in the terminal, one question at a time, as `grilling` says.
   Where a question is easier to see than to read, such as alternative interfaces for the deepened module or its before and after, add it to the same page and wait for their annotations again.

Answer the person's annotations with `npx -y lavish-axi reply` as its help describes, so the page shows what changed.

## Installing elsewhere

This document installs with every skill of `socrates design skill`, under `$CLAUDE_CONFIG_DIR/skills`, or `~/.claude/skills` where it is unset.
Write all three somewhere else instead with `socrates design skill --host claude --to <directory>`.
