# socrates: one command, and a workflow is a subcommand of it

## Problem Statement

`socrates-review` asserts a parallel that is not true.

It reads as a peer of `socrates`, when the truth is that review is one workflow among several the brand will publish.
A person who types it is told, by the shape of the name, that reviewing sits beside the brand rather than inside it.

Three things follow from that, and they are all visible from a command line.

**The workflow is spelled twice.**
`socrates-review review --wait` names reviewing in the binary and again in the verb.

**The second program is a stranger.**
`agora` is published in the same archive and installed by the same command, and nothing about its name says what published it.
A person who lists what the suite put on their path finds two names with nothing in common.

**The next workflow has nowhere to go that is not worse.**
socrates-grill would arrive as `socrates-grill`, and the suite would publish `socrates-review`, `socrates-grill` and `agora`: three binaries, two of which spell the brand and one of which does not.
Nobody publishes `aws-s3` and `aws-ecs`.
The convention every multi-service CLI settled on is one command and a namespace under it, and `docker-compose` becoming `docker compose` is the same move made by somebody who had already shipped the hyphen.

The window to fix this is open and closing.
`v0.1.0` is the only tag, it shipped the reviewer as `socrates`, and the rename to `socrates-review` has never been released.
Nobody has yet learned a name this would take away.

## Solution

`socrates` is the command.
A workflow is a subcommand of it.

```
socrates review begin --intent "..." --wait
socrates review status
socrates review findings --json
socrates review respond --finding ... --fixed --note "..."
socrates review continue --wait
socrates review watch
socrates review logs
socrates review cancel
socrates review config
socrates review configure
socrates review skill
socrates agora serve
socrates agora poll
socrates agora end
socrates agora stop
socrates version
```

The suite publishes one binary.
Installing puts one name on a path, and every workflow and component the suite holds is reached through it.

The top level holds routing and what is true of the whole suite, which is `version` and `help` and nothing else.
Everything a tool owns is reached through that tool's namespace, including its configuration and its skill.
A namespace typed bare prints help and starts nothing.

The command that begins a run is `begin`, which is the word the source has always used for it.
The stutter is gone because the workflow is named once, by the namespace.

## User Stories

**Somebody installing the suite**

1. As somebody installing the suite, I want one name on my path, so that I learn one command rather than an inventory.
2. As somebody installing the suite, I want every workflow reachable from that one name, so that discovering what the suite does is typing it rather than reading a list.
3. As somebody installing the suite, I want the archive to carry one binary, so that what I extract and what `make install` gives me are visibly the same thing.
4. As somebody who has used any modern multi-service CLI, I want socrates to follow the convention I already know, so that I can guess a command before I look it up.
5. As somebody who lists what a tool installed, I want nothing on my path whose name does not say what published it.

**Somebody typing commands**

6. As somebody starting a review, I want to name the workflow once, so that the command reads as an instruction rather than a stutter.
7. As somebody starting a review, I want the verb to say that a run begins here, so that I can tell it apart from the commands that read an existing run.
8. As somebody who types a namespace and nothing else, I want help rather than a side effect, so that exploring the command surface is safe.
9. As somebody scripting against socrates, I want an exit code to mean the same thing whichever workflow produced it, so that a script can branch on it without knowing which namespace ran.
10. As somebody scripting against socrates, I want `1` to mean a person is wanted, whether that is a review that escalated or a page nobody answered.
11. As somebody configuring a workflow, I want its settings reached through its own namespace, so that I cannot change one tool's behaviour from another tool's menu.
12. As somebody who wants to know what version I have, I want one answer for the whole suite, so that a version in a bug report identifies one set of behaviour.
13. As somebody who serves a page, I want the command to say which brand serves it, so that a process on my machine is attributable.

**A coding session driving a review**

14. As a coding session, I want the review loop to be the same loop with the same states, so that nothing about how a review is driven changes because the command was renamed.
15. As a coding session, I want the page still reached by serving and polling, so that the handover is unchanged in everything but its prefix.
16. As a coding session, I want the skill to write every invocation for me, so that the extra word costs me nothing I have to remember.
17. As a coding session, I want socrates to keep saying whether a round waits on a person and nothing more, so that the scope boundary I rely on is untouched.

**A person answering a round**

18. As a person answering a round, I want the page to look and behave exactly as it did, so that a command rename is invisible to me.

**A contributor**

19. As a contributor, I want each tool's directory named for the tool and nothing else, so that no path in the repository spells the brand twice.
20. As a contributor, I want the repository's own glossary to distinguish a workflow from a component, so that I can tell which rules apply to the thing I am changing.
21. As a contributor, I want the dispatcher to hold wiring and no behaviour, so that no logic can accumulate in the one place that imports everything.
22. As a contributor, I want the rule that no tool imports another tool to survive a single binary, so that the boundary between tools is still checked rather than remembered.
23. As a contributor, I want the network-reach check to measure the review code rather than the binary that contains it, so that linking a server beside it cannot quietly weaken the claim.
24. As a contributor, I want each change in this sequence to be readable on its own, so that a behaviour change is never buried inside a substitution that touches every file.
25. As a contributor, I want the command surface decision recorded, so that somebody who finds the archived rename rejecting `socrates` as a binary name learns why that rejection does not apply here.
26. As a contributor, I want the decision to fold the page into one binary recorded with its cost, so that nobody undoes it by accident and nobody adopts it without knowing what it gave up.
27. As a contributor reading the security document, I want it to say that the binary running reviews now contains a server, so that I am not surprised by what the process can do.

**A future socrates-grill**

28. As socrates-grill, I want a namespace waiting for me rather than a binary name to invent, so that arriving costs me nothing structural.
29. As socrates-grill, I want the page reached the same way review reaches it, so that nobody has to explain why two workflows use one component differently.
30. As socrates-grill, I want the exit codes to already mean something general, so that I inherit a contract rather than inventing a fourth table.
31. As socrates-grill, I want the top level to hold nothing of review's, so that I am not a guest in another workflow's namespace.

**An agent implementing this**

32. As an agent implementing this, I want the directory move to be a pure relocation, so that a failing test after it means the move was wrong and nothing else.
33. As an agent implementing this, I want the verb renamed before the namespace exists, so that the change that rewrites every invocation is a prefix substitution with no verb moving inside it.
34. As an agent implementing this, I want the command trees extracted into packages while two binaries still stand, so that the extraction is proved by the existing suite passing unedited.
35. As an agent implementing this, I want each change to restate the checks it breaks, so that no check is left failing for a later change to notice.
36. As an agent implementing this, I want the test inventory compared before and after the move, so that a test file lost in the move is caught rather than passing silently by absence.
37. As an agent implementing this, I want the changelog written once at the end, so that the same paragraph is not drafted five times.

## Implementation Decisions

**The command surface is `socrates <namespace> <operation>`.**
The namespace names a workflow or a component.
The operation names what is done.
This supersedes the archived rename's conclusion only in part: that change rejected leaving the binary named `socrates` while it was the reviewer, and that rejection stands.
A dispatcher is a different thing, because `socrates` then names the brand and routes, and is nobody's workflow.

**The top level holds routing and what is true of the whole suite.**
That is `version` and `help`.
`config`, `configure` and `skill` sit inside their namespace, because configuration belongs to a tool and each workflow will have a skill of its own.
A suite-level `configure` would have to grow a tool argument the first time a second tool had settings, which is the suite-wide configuration file that has already been rejected once.

**A namespace typed bare prints help and starts nothing.**

**The command that begins a run is `begin`.**
The reviewer package is already described as beginning, supervising, responding, advancing and cancelling, and the function that starts a run is already named for it, so the command surface and the source agree without anybody deciding twice.
`start` is a synonym with no claim on the domain.
`run` is refused outright: a run is a noun in the glossary, meaning one review from beginning to terminal state, so the command would read as an operation on a run rather than the creation of one.

**agora folds in and keeps its name.**
`socrates agora serve` and `socrates agora poll`.
The name has equity and is in the glossary, and a namespace name is not a word inside agora's own tree, so the check that forbids review vocabulary there is untouched.

**The suite publishes exactly one binary.**
Each tool exposes its command tree as a package of its own.
The dispatcher lives at the root, belongs to the suite and to no tool, holds wiring and nothing else, and sits beside the shared layer for the same reason the shared layer is there.
The simulator and gate binaries are unaffected: they are never installed, and the one-binary rule is about what the suite publishes.

**The network-reach claim moves from the binary to the review tree, and gets stronger.**
The check today computes its closure from the review binary's `main`, which measures the review code only because that binary happens to contain nothing else.
The scope becomes the review command tree's package, so the claim is about the review code itself and stays true no matter what else is linked beside it.
The measured claim is unchanged in meaning: the review code reaches no network outside a build tag.

**The rule that no tool imports another tool survives untouched.**
The review tree still imports nothing of agora's, and the only thing importing both is the dispatcher, which is not a tool.
This restates, rather than abandons, the monorepo spec's decision that socrates-review links none of agora's Go code.
That decision's real content was never that the binary omits an HTTP server; it was that no review code reaches the page and no review command names one.
Both are checkable over the import graph and over the command surface, and both are kept.

**The residual is stated rather than hidden.**
The binary that runs reviews now contains an HTTP server, where today it provably does not.
The security document says so.

**One exit table belongs to the suite.**
The two programs arrived independently at the same four slots, so this codifies an agreement rather than inventing one.

| code | meaning |
| --- | --- |
| 0 | nothing is asked of you |
| 1 | something is asked of you |
| 2 | the command failed |
| 3 | the command was not usable |
| 4 | the work is still running |

`clean` and an answered page are both `0`.
`escalated` and an unanswered page are both `1`.
The table moves to the shared layer and the review tree keeps the function that turns a run state into one of its codes.
No code changes value, so no script changes behaviour.

**The review tool's directory is named for the tool and not for the brand.**
It becomes `review`, beside `agora`, matching the directory the tool has written under since the last rename.
The argument that carried the command surface carries here too: a name that reads as a peer of the brand on a command line reads as a peer of the brand in a path.

**The constant that names the repository splits.**
It is currently doing four jobs: the installed binary's name, the directory in the repository, one of the names a code-hosting URL may carry, and the string every refusal prints.
The suite's own name takes the first, third and fourth, and the per-tool directory is read from the stated tool set rather than from a constant.

**A name is brand-qualified exactly where it lands in a flat global namespace with no dispatcher.**
The agent skill keeps the name `socrates-review`, because the directory skills are installed into is such a namespace.
A directory inside this repository is not, and neither is a subcommand.
This is the rule, and the skill is its one case rather than an exception to it.

**Nothing is carried forward.**
No shim, no alias, no deprecation, no migration.
`socrates-review` as a binary name has never been released.
`agora` as a binary name shipped in the only tag, and goes anyway, by the precedent already recorded for the last rename: code that migrates a layout has to be correct about a state that exists on one machine for as long as it takes to run one command, and then stays in the tool forever, because nothing tells you when the last machine has moved.

**Unchanged, and confirmed rather than decided.**
The environment variable that moves the brand's root keeps its name and its meaning.
The guide read out of a repository under review keeps its path.
Both were already shaped for this.

**The changelog is written once, at the end.**
None of the changes below touch it.
The issue that renames the repository already exists to write one coherent unreleased section, and writing it five times produces five drafts of one paragraph.

**Five changes, in this order.**

1. The review tool's directory is renamed. A pure relocation, imports rewritten, nothing a person types changes.
2. The verb becomes `begin`, while the binary is still named for the workflow.
3. The exit codes become the suite's.
4. Each tool's command tree becomes a package, with both binaries still standing and every invocation unchanged.
5. The dispatcher arrives, the suite publishes one binary, and every invocation is rewritten.

Step 2 before step 5 is the load-bearing ordering.
Renaming the verb while the binary is still `socrates-review` means the change that touches every invocation changes a prefix and nothing else, which is a diff a person can read.
Step 4 before step 5 means the extraction is proved by the existing suite passing unedited, so the dispatcher's change contains no refactoring.

A change restates the checks it breaks.
The change that introduces the dispatcher restates the installed-name check, the archive check, the reach scope and the constant split, rather than leaving them for the later issue that restates the layout checks.

**Two decision records, landing with the dispatcher.**

- The page folds into one binary.
  It supersedes a decision written down specifically to stop itself being undone by accident, restates it as an import-graph rule, and records the residual.
  The one exit table is recorded here as a consequence and does not get a record of its own.
- The command is `socrates <namespace> <operation>`.
  A reader who finds the archived rename rejecting `socrates` as a binary name needs the document that says why a dispatcher is not what was rejected.
  The trade-off recorded is a word added to every invocation.

Not warranted: the directory rename, the exit table, the `begin` verb, and where the dispatcher's wiring sits.
Each follows from a decision already recorded.

**Three issues already open are made stale by this.**
The one that splits documentation per tool names a directory this renames.
The one that restates the layout checks asserts that every tool the suite holds has a binary, which this makes false, and adds a rule about cross-tool imports that now has to say something about a dispatcher that imports every tool and is not one.
The one that renames the repository writes a changelog that now has a different shape to describe.
All three lose the label that marks them ready for an agent immediately, and their bodies are amended when this spec is published.

## Testing Decisions

**What a good test is here.**
It asserts what a person or a session observes: what a binary prints, what installing puts on a path, what an archive carries, what a command exits with, what the import graph permits.
It does not assert that a package has a name or that a file sits at a path, except in the convention package, whose whole job is to make structural claims executable.
Tests name their scenario, and a test that legitimately reads Go source is registered in the existing list of contract-reading tests.

**No new seam.**
Everything lands on the two that already exist.

The end-to-end package takes everything a person or a session observes, against real binaries, a real repository, real processes, each vendor's simulator and a real browser.
The six journey tests that drive both programs through the skill loop are the regression proof that this changed no behaviour, and they need nothing beyond the new invocations.

The convention package takes every structural claim: what the release builds, what the archive carries, the import graph, the network closure, and the wording checks.
The new reach scope is not new apparatus.
The existing check already parses the tree, computes a closure from a named scope and reports what it reached, so this is a second scope for the same machinery, exactly as the cross-tool import rule already planned to be.

The exit table keeps a package-level test for the meaning of each code, which is where its existing one lives.

**What proves each change.**

1. The directory rename is proved by the existing suite passing unchanged, plus a one-off inventory of package and test names captured before and diffed after, which must be empty.
   This is the hazard the last relocation closed explicitly, because a move can lose a test file and still go green.
2. The verb is proved end to end: a run begins under the new word, and the old word is refused as an unknown command with the usage code.
3. The exit codes are proved twice: a package test that each code means one thing, and end-to-end assertions that a clean run, an escalated run, an answered page and an unanswered page exit with the codes the table states.
4. The extraction is proved by the existing suite passing with no test edited at all.
   A test that needed editing would mean the extraction changed behaviour.
5. The dispatcher is proved end to end: installing puts one binary on a path and nothing beside it, every namespace and operation runs, a namespace typed bare prints help and starts nothing, the top level offers nothing a workflow owns, and the archive carries one binary with the licence and notice files.
   The reach check is proved by its restated scope reporting the review tree and reaching no network.

**Prior art to follow.**
The installed-name test already installs through the make target and reads what landed on the path.
The archive test already extracts a release archive and asserts what it carries.
The reach check already computes a closure from a named scope.
The refusal for an unknown command already exists and already names the tool, which is what the verb rename and the retired binary names reuse.

**Mutate a guard.**
Three are worth breaking on purpose before this is called done.
Make the review tree import the page's tree and confirm the cross-tool check fails.
Add a second published binary and confirm the one-binary check fails.
Remove a test file during the relocation and confirm the inventory diff catches it.

## Out of Scope

**socrates-grill.**
No namespace, no directory, no package.
It arrives as its own spec.
Nothing here is built for it beyond leaving the room it will need.

**Any change to the review loop, the dialogue, the verdict rule, the round cap, the snapshot, or the reviewer seam.**
This changes what a person types and what is published, and nothing about what a review does.

**Any change to what the page looks like, what it can render, or what it knows.**
Its vocabulary exclusion is untouched, and folding it into one binary does not teach it what a review is.

**The page's document model.**
Whether a grill question fits it is grill's question to raise.

**The configuration split and the page's own state directory.**
Both are already specified and sequenced ahead of this.

**The documentation split, the layout checks, and the repository rename.**
All three are already open, all three are made stale by this, and all three stay where they are in the sequence.

**Shell completions.**
A dispatcher makes them more valuable and nothing here depends on them.

**Migrating anybody's machine.**

## Further Notes

**Where this sits in the sequence.**
It follows the configuration split and the page's state directory, and precedes the documentation split and the layout checks.
That ordering is what stops the per-tool documents and the executable layout claims being written against a shape that is about to change, and then written again.

**Timing.**
One tag exists and this is unreleased.
The last spec observed that this was the last moment the suite could be named without a user living through two renames, and that window is still open because none of what it described has shipped.
This rides the same unreleased section.

**Why the archived rejection does not apply.**
The archived rename considered leaving the binary named `socrates` and renaming only the repository, and rejected it on the second tool: two binaries published by the organisation, one of them called `socrates`, leaves a person unable to say which one they mean.
That reasoning is correct and this does not contradict it.
Under a dispatcher `socrates` is not one of the tools, so there is no pair to be confused between.
It is the brand, it routes, and every tool is named by the word after it.

**The exit tables already agreed.**
Two programs written separately, with no shared code between their command layers, arrived at the same four slots with the same meanings.
Codifying that is recognition rather than design.

**One existing invariant gets stronger.**
The check that one name is installed and nothing is put beside it currently passes because the page's binary is a different file.
Under one binary it becomes literally what it says.
