# socrates: one repository, one directory per tool, agora beside socrates-review

## Problem Statement

agora is a squatter.

It is a second program with its own binary, its own vocabulary and its own reason to exist, and it lives inside a repository named after a different tool.
The repository is `socrates-tools/socrates-review`, the module is `github.com/socrates-tools/socrates-review`, and `agora` is a directory inside it.

Nothing about that is visible from the outside, and it costs in three places.

**A second tool has nowhere to go.**
socrates-grill is coming, and it needs the same page agora already serves.
Under today's layout it would either live inside a repository named after socrates-review, which is absurd, or force agora out into its own repository, which is what this spec's first draft costed and rejected.

**agora's settings live in socrates-review's file.**
`~/.socrates/review/config.yaml` carries a `surface:` block that socrates-review never reads and agora always does.
`socrates-review config` prints those two keys as though they were its own, and `socrates-review configure` offers them for editing.
The package `internal/loopback` exists for no other reason: ADR 0024 records that `internal/config` needed agora's default bind address to report `surface.bind`, and importing agora's server package would have linked `net/http` into `cmd/socrates-review`.
A whole package and a whole convention check exist as scar tissue from one misplaced YAML block.

**The naming already anticipated this and then stopped short.**
`v0.1.0` shipped the reviewer as `socrates`, with `scrat` beside it.
The unreleased rename to `socrates-review` states its own reason: "`socrates` is the brand the organisation publishes under rather than this tool, so a second tool under it has a name to take."
The name was freed and nothing took it.

## Solution

The repository becomes the suite.

`socrates-tools/socrates`, module `github.com/socrates-tools/socrates`, one Go module, one directory per tool:

```
socrates/
  agora/
  socrates-review/
  internal/          the layer both tools share
  docs/adr/          decisions about the suite
  CONTEXT-MAP.md
  Makefile
```

`socrates-grill/` joins them when it exists, and not before.

Each tool owns its own home under the brand root and its own configuration file:

```
~/.socrates/review/config.yaml
~/.socrates/agora/config.yaml
```

`SOCRATES_HOME` still moves the whole root.
`AGORA_HOME` stays retired.

The two tools stay as independent at runtime as they are today.
socrates-review links none of agora's Go code even though a single module would make the import free.
A coding session runs `agora serve` and `agora poll` exactly as the skill already tells it to.

One version, one archive, one hand-written CHANGELOG.
A person downloads `socrates_<version>_<os>_<arch>.tar.gz` and gets every binary the suite publishes, plus `LICENSE` and `NOTICE`.

Because the rename to `socrates-review` has not been released, all of this folds into the same unreleased section.
Nobody lives through two renames.

## User Stories

**Somebody installing the tools**

1. As somebody installing these tools, I want one archive that carries every binary, so that extracting once gives me what installing from source gives me.
2. As somebody installing these tools, I want `make install` to put every tool on my path, so that I do not have to know how many binaries the suite publishes.
3. As somebody who downloaded the suite, I want every binary in it to report the same version, so that a version in a bug report identifies one set of behaviour rather than several.
4. As somebody who reads the repository name before installing, I want it to name what it holds, so that I am not surprised to find a second program inside it.
5. As somebody who found the project through agora, I want to understand immediately that agora is a component of a suite rather than a product on its own, so that I do not go looking for a use for it by itself.

**Somebody who already has socrates-review**

6. As somebody who already configured a bind address, I want the tool to tell me where that setting went rather than refusing my whole file with a generic parse error, so that I can fix it in one edit without reading the README again.
7. As somebody who already configured a bind address, I want the tool to stop rather than quietly run on defaults, so that I never discover mid-review that my configuration was ignored.
8. As somebody upgrading, I want one entry in the changelog explaining the move in terms of what I have installed, so that I learn what changed for me rather than what changed in the repository.
9. As somebody upgrading, I want the rename and the reorganisation to arrive in one release, so that I do not relearn the tool's name twice.
10. As somebody with a run in progress, I want my existing runs and their findings to still be readable after the upgrade, so that the move costs me no work I had already done.

**Somebody configuring a tool**

11. As somebody configuring socrates-review, I want `socrates-review config` to print only settings socrates-review reads, so that the list is what governs my reviews and nothing else.
12. As somebody configuring socrates-review, I want `socrates-review configure` to offer only settings socrates-review reads, so that I cannot change a page's behaviour from a review tool's menu.
13. As somebody configuring the page, I want agora to have a configuration file of its own, so that I know which tool a setting belongs to from the file it is in.
14. As somebody who moves their state directory, I want one variable to still move everything, so that relocating the suite stays one decision rather than one per tool.

**A coding session driving a review**

15. As a coding session, I want the review loop to work exactly as it did before the reorganisation, so that a skill written against the old layout needs no changes.
16. As a coding session, I want `agora serve` and `agora poll` to remain the way a round reaches a person, so that nothing about the handover changes because directories moved.
17. As a coding session, I want socrates-review to keep saying whether a round waits on a person and nothing more, so that the scope boundary I rely on is unchanged.

**A person answering a round**

18. As a person answering a round, I want the page to look and behave exactly as it did, so that a reorganisation is invisible to me.
19. As a person answering a round, I want the page to stay ignorant of what a review is, so that the same page can serve a grill session later without learning two vocabularies.

**A contributor**

20. As a contributor, I want one verification command that runs every gate, so that passing locally and passing in continuous integration are one claim rather than two.
21. As a contributor changing one tool, I want the other tool's tests to run too, so that I cannot break a shared package and find out from somebody else.
22. As a contributor, I want each tool to have its own glossary, so that I read the vocabulary of the thing I am changing and not of the thing beside it.
23. As a contributor, I want each tool to have its own decision record directory, so that a decision about the page is not filed among decisions about reviews.
24. As a contributor, I want a decision record to keep its number when it moves, so that a cross-reference written last month still resolves.
25. As a contributor, I want the repository's own conventions to forbid one tool importing another, so that the boundary between tools is checked rather than remembered.
26. As a contributor reading the map file at the root, I want it to describe the suite, so that I can find the tool I am looking for without opening every directory.
27. As a contributor working inside one tool's directory, I want that tool's own instructions to be what I read, so that I am not held to invariants belonging to a tool I am not touching.

**A future socrates-grill**

28. As socrates-grill, I want agora to already be a peer directory rather than something I reach into another tool for, so that depending on it costs me nothing structural.
29. As socrates-grill, I want to reach agora the same way socrates-review does, so that nobody has to explain why two tools use the same component differently.
30. As socrates-grill, I want agora's vocabulary to still exclude review words, so that the page can carry my questions without renaming anything.
31. As somebody planning socrates-grill, I want the shared layer to grow only when a second tool actually needs something, so that I inherit no abstraction invented for a tool that did not exist yet.

**An agent implementing this**

32. As an agent implementing this, I want the reorganisation to arrive as a pure move with no behaviour change, so that a failing test after it means the move was wrong and nothing else.
33. As an agent implementing this, I want the module rename to happen before the directory move, so that imports are rewritten once rather than twice.
34. As an agent implementing this, I want each behaviour change to be a change of its own after the move, so that a bisect lands on a decision rather than on a relocation.
35. As an agent implementing this, I want the test inventory compared before and after the move, so that a test file lost in the move is caught rather than passing silently by absence.
36. As an agent implementing this, I want the checks that encode the old layout named up front, so that I restate them deliberately rather than deleting them when they fail.

## Implementation Decisions

**One module, not one per tool.**
`github.com/socrates-tools/socrates`.
Three modules with a workspace was considered and rejected: a workspace only fixes local development, so the shared layer would have to become a fourth published module with its own version, and every convention change would become a four-module release dance.
The cost accepted is that the tools release together, which is the coupling the first draft of this work was trying to remove.
This supersedes ADR 0013, which argued from a premise that no longer holds.

**The shared layer lives in root `internal/`.**
Go's own internal rule then makes each tool's own `internal/` invisible to the other at no cost.
Packages move up to the root only when a second tool genuinely needs them, never on speculation.
The two duplicated build stamps unify into one shared version package, since the suite now has one version; agora currently repeats the describe logic in its own `main`.

**socrates-review links none of agora's Go code.**
This is deliberate and is the decision most likely to be undone by accident, because in one module the import is free and would make the document schema compiler-checked.
It is refused because socrates-review does not render, serve, or wait for a person about a review, per ADR 25, and because no command, flag or run record of socrates-review may name a page.
Emitting a document skeleton from socrates-review was considered and rejected: the skill has a subagent draft the question text, the meaning, the illustration and the stylesheet, none of which socrates-review can write, so a skeleton would buy field names only, at the cost of a flag that names a page.
The drift protection arrives free instead, from the end-to-end tests that already drive both binaries against the real schema.
socrates-grill will reach agora the same way.

**One configuration file per tool.**
The `surface` block leaves socrates-review's configuration entirely, taking its two keys with it, and `config` and `configure` lose those rows.
This dissolves ADR 0024 rather than working around it: `internal/loopback` collapses back into agora's server package, and the convention check that exists only because a second binary in the module imports the network goes with it.
A suite-wide configuration file with a block per tool was considered and rejected, because it keeps ADR 0024's problem and multiplies it by the number of tools.

**agora's state moves to its own directory under the brand root.**
It stops resolving a path inside socrates-review's home.
`SOCRATES_HOME` still moves the root, so the invariant that one variable moves everything survives with a wider name.
`AGORA_HOME` stays retired and its existing refusal stands.

**A configuration file still carrying a `surface` block is refused by name.**
The loader rejects unknown keys, so such a file currently produces a generic unparseable-file refusal that does not say where the keys went.
socrates-review recognises the block specifically and names the file the keys now belong in.
This follows the house pattern already used to recognise the retired `AGORA_HOME` and point at its replacement.
Tolerating and ignoring the block for one release was rejected, because socrates-review states that it will not quietly run on defaults while a configuration file says otherwise.

**One version, one archive, one changelog.**
The archive carries every binary the suite builds plus the licence and notice files.
Per-tool archives were rejected: agora serves one document and waits, which is useless without something writing documents, so it is a component and not a product.
agora's reported version jumps to the suite's, which the changelog's existing note about the surface half changing under a zero major version already covers.
The notice file stays one file and keeps its existing two-part structure, embedded assets and linked modules, since one archive ships both.

**Documentation splits by tool.**
A context map at the root, a glossary per tool, decision records per tool alongside a root directory for decisions about the suite.
agora's defining invariant is that it does not know what a review is, enforced by an existing convention check; a tool with an enforced vocabulary boundary is a bounded context, and sharing a glossary would put two vocabularies in one file that a test forbids mixing.
The instruction file splits the same way: the root keeps the map, the verification sequence and suite-wide invariants, and each tool takes its own layout and invariants.

**Decision records keep their numbers when they move.**
One number space across the root and every tool directory, gaps and all, forever.
A given number therefore names one document no matter which directory holds it, and every existing cross-reference and commit message still resolves.
New records continue from the high-water mark.

**One verification entry point, unchanged.**
Every contributor needs a browser, including one changing a single message string in socrates-review.
Per-tool verification targets were rejected because they create a category of mistake that does not exist today, running the wrong gate and learning about it from continuous integration.
Skipping the browser tests when no browser is present is already forbidden: an absent browser fails rather than skips.

**The checks that encode the old layout, restated rather than deleted.**
The check that the repository is named for the tool becomes a check that it is named for the suite and that every tool it holds has a binary.
The check that the release builds the tool under one name and a second binary becomes a check that the release builds every tool the suite states.
The check listing the directories the page's vocabulary applies to shrinks to agora's directory.
The check that a second binary in the module is not socrates-review reaching further is deleted outright, along with the package it was guarding.

**One genuinely new rule: no tool imports another tool.**
Go's internal rule already keeps each tool's `internal/` private, but says nothing about non-internal packages, so one tool's public package could import another's and compile.
The suite forbids it, checked over the import graph.

## Testing Decisions

**What a good test is here.**
It asserts what a person or a session observes: what a binary prints, what a file refuses, what an archive contains, what the import graph permits.
It does not assert that a file sits at a path or that a package has a name, except in the convention package, whose entire job is to make structural claims executable.
Tests continue to name their scenario, and a test that legitimately reads Go source is registered in the existing list of contract-reading tests.

**The move itself is proved by the existing suite passing unchanged.**
That is the whole point of doing it as a pure relocation first.
One hazard has to be closed explicitly: a move can lose a test file and still go green, because nothing counts tests.
So the full inventory of package and test names is captured before the move and diffed after it, and the diff must be empty.
This is a one-off check for the relocation change, not a permanent gate.

**The convention package takes nearly every structural assertion**, which is where they already live.
Restated: the repository is named for the suite and every tool it names has a binary; the release builds every tool the suite states and the matrix equals the generated platform list; the archive carries every binary plus the licence and notice files; the page's vocabulary exclusion applies to agora's directory.
Unchanged and now spanning three trees: no em dash reaches a person, every user-facing string lives in its package's message registry, no source file carries commentary, no licence header opens a file, no message names an environment variable the binaries do not read, no document names a code-hosting account other than the organisation's.
New: no tool imports another tool, checked over the import graph using the same scope-and-closure machinery the existing network-reach check already uses.
Deleted: the check that a second binary in the module is not socrates-review reaching further.

**The end-to-end package takes the behaviour changes**, against real binaries and a real repository.
A configuration file still carrying a `surface` block produces a refusal that names the file the keys moved to, and socrates-review does not run.
One variable still moves everything both binaries write, with each tool's directory now under the brand root.
agora writes under its own directory and nothing lands in socrates-review's.
The archive carries every binary the suite builds.
The six journey tests that drive both binaries through the skill loop are the regression proof that the reorganisation changed no behaviour, and they need no edits beyond paths.

**Prior art to follow.**
The import-graph check follows the existing network-reach check in the same package, which already parses the tree, computes a closure from a named scope, and reports what it reached.
The refusal that names a move follows the existing refusal for the retired home variable, which recognises a specific retired input and points at its replacement.
The configuration tests follow the existing ones for keys refused outside their block, which are themselves scar tissue from the last time these two keys moved.
The archive test follows the existing one that extracts a release archive and asserts what it carries.

**Mutate a guard.**
Two are worth breaking on purpose before this is called done: make one tool import another and confirm the new check fails, and remove a test file during the move and confirm the inventory diff catches it.

## Out of Scope

**socrates-grill.**
No directory, no stub, no package.
It arrives as its own spec.
Nothing in this change is designed for it beyond leaving the room it will need.

**Generalising agora's document model.**
Its vocabulary is review-shaped in places, carrying severity, audience and location.
Whether a grill question fits that is a real question and it is grill's to raise, in agora's own directory, with its own tests.
An extraction that also redesigns gives no way to tell which half broke the caller.

**Any change to the review loop, the dialogue, the verdict rule, the snapshot, or the reviewer seam.**
This change moves directories and splits a configuration file.

**Any change to what the page looks like or what it can render.**

**Independent release cadence for the tools.**
This was the goal the first draft pursued and it is deliberately given up.
The tools ship together, at one version, in one archive.
If that becomes intolerable, the way back is a multi-module split, and the cost of that is already written down.

**Migrating anybody's existing configuration file automatically.**
The `surface` block is refused with an explanation, not rewritten.

## Further Notes

**The sequence.**
One spec, several changes under it, in this order:

1. #47 Rename the module path, as one mechanical pass with no directory moved.
2. #41 Move the directories, with no behaviour change, and diff the test inventory.
3. #42 Split the configuration: the `surface` block leaves, `internal/loopback` collapses into agora, and the refusal that names the move arrives.
4. #43 Move agora's state to its own directory under the brand root.
5. #44 Split the documentation: context map, per-tool glossaries, per-tool decision directories, per-tool instruction files.
6. #45 Restate the convention checks that encode the old layout, and add the rule that no tool imports another.
7. #46 Rename the repository, and fold the whole thing into the unreleased changelog section.

The module rename goes first so that imports are rewritten once.
The repository rename goes last because it is the only step that needs anything outside this clone.

**Decision records this warrants.**
Three, by the usual bar of hard to reverse, surprising without context, and the result of a real trade-off:

- One module rather than one per tool. Supersedes ADR 0013.
- socrates-review links none of agora's Go code, in a monorepo where the import would be free. This is the one most likely to be undone by a well-meaning future change.
- One configuration file per tool. Borderline on its own, and written mainly so that ADR 0024 has something to point at when it is superseded.

Not warranted: the repository name, the single archive, or requiring a browser of every contributor.
The first two follow from decisions already recorded, and the third belongs in the contributing guide.

**Timing.**
`v0.1.0` is the only tag, and the rename from `socrates` to `socrates-review` is unreleased.
This is the last moment the suite can be named without a user experiencing two renames.
