## Context

`a-round-a-person-reads-becomes-a-page` decided that some rounds are handed to a person as a page, stated the bar that page is held to, and settled where it is drafted.
It also decided the page is read only, and gave the reason plainly: there is nothing behind it, so a control would take a decision that reaches no run.

`a-page-that-answers-back` builds the thing behind it.
This change is the rewiring, and it is confined to one package.

Two things are being fixed at once and they are not the same thing.
The first is that `artifact-design` puts the document on a host outside the machine, which contradicts `internal/corpus` refusing to record the same document for the same reason.
The second is that a person's decision currently travels back to the run through their memory and the session's paraphrase.

The second is why the read-only rule inverts, and it inverts on its own stated premise rather than in spite of it.

## Goals / Non-Goals

**Goals:**

- The findings document never leaves the machine, and issue 1 is closed by the code rather than restated.
- A person answers where they are reading, and what reaches the run is their words rather than a summary of them.
- The identifier match between page and terminal becomes mechanical.
- The check that the skill names only real commands learns that there are now two binaries to check against.

**Non-Goals:**

- Changing which rounds a person reads, what a round leads to, or what `socrates findings` reports. All three are already right and none is touched.
- Changing what a decision means. `socrates respond --decided` is unchanged, a decision is still a claim, and the next round's reviewer still judges it.
- Letting the surface record anything. The session runs the command; the page names it.
- Adding any command, flag, dependency or network reach to `socrates`.

## Decisions

**The rule inverts on its own premise.**
The spec said the page carries no control *because* an answer would reach no run.
The honest edit is not to delete that sentence but to keep the reasoning and move the conclusion: a control is permitted exactly where it names the command it becomes.
This matters because the failure the old rule prevented is still live in a different shape. A control whose command the page cannot name puts the person back to deciding something whose destination they cannot see, and it is refused for the original reason.
Alternative rejected: "the page may now carry controls", which throws away the reasoning and leaves nothing to refuse a bad one with.

**The session records, and the surface does not.**
Letting `agora` run `socrates respond` itself would be fewer moving parts and would put the person's words in the record with nothing in between.
It would also make a page a writer into the run and would require the surface to know socrates exists, which is precisely what `a-page-that-answers-back` is built not to know.
The claim stays the session's to make and the next reviewer's to judge, which is the invariant this repository is least willing to weaken.

**The identifiers are copied, not instructed.**
The existing requirement says every finding keeps the identifier, severity and audience socrates gave it. Today that is an instruction to the session and holds by its care.
The surface reads `socrates findings --json` and copies them, so a session that wrote something friendlier loses.
The requirement does not change; the mechanism under it does, which is the cheap half of this change and the one most likely to have silently failed otherwise.

**The surface check has to learn to count to two.**
`skill.Named` extracts a command by finding the product name and taking the word after it, so today it sees only what follows `socrates`.
Left alone it would silently ignore every `agora` invocation the skill names, and the check whose whole purpose is catching a plausible-but-wrong invocation would pass on the newest source of them.
So extraction attributes each name to a binary and the check resolves it against that binary.
A name belonging to no binary fails, rather than being skipped as unrecognised, because skipping is how this would go quiet again.

**Delegation at the parked branch survives, and its reason is re-derived rather than assumed.**
The original reason was that drafting a whole HTML page in the driving context costs tokens that are re-sent on every subsequent turn.
Writing a document is a fraction of that markup, so the reason had to be checked rather than carried.
It survives because the expensive part was never the markup: it is reading every finding and writing plain English about what each means, which is the same work and lands in the same context.
What is newly specified is that the delegate returns the document and the driving session serves and waits, because a delegate that served the surface would exit owning it.

**Where the branches now diverge is smaller than it was.**
Both branches draft a document, serve it, and wait. Only who drafts differs. The skill says the difference and its reason once, and neither branch restates the bar.

## Risks / Trade-offs

**A session could answer a person's finding on its own behalf through the page.** → It could before, by writing `--decided` with an invented reason. Nothing here makes that easier and nothing makes it impossible. The mitigation is unchanged and is stated rather than enforced: the skill says the session does not answer one on its own behalf, a test asserts the sentence is there, and the recorded decision is a claim a fresh reviewer judges. A scenario here keeps that sentence true across the new surface.

**The skill grows a second tool to drive, and every sentence added is read on every turn.** → The loop gains three invocations, not a section. The bar, the branches and the response rules are unchanged in length, and the polling step is untouched, which is the step whose cost is multiplied.

**The two binaries can drift.** → The surface check now fails on any command the skill names that no binary offers, which is the drift that actually bites. The document contract is covered on the other side by `a-page-that-answers-back`.

**Reviewing prose that was rewritten from a stale copy.** → `make scenario-drop` reads deltas and not prose, and this change rewrites requirement text that a previous change wrote. The requirements here were copied from the current source rather than retyped, and the reviewer of this change is what catches a sentence that drifted.

## Migration Plan

No data, no state and no run record changes shape, so there is nothing to migrate and nothing to roll back beyond reverting the skill.

Ordering is the whole of the plan:

1. `a-round-a-person-reads-becomes-a-page` is archived, promoting `findings-page` into `openspec/specs`. Until then the deltas in both changes have nothing to compare against.
2. `a-page-that-answers-back` is implemented, so the binary this skill names exists. Implementing this change first would fail its own surface check.
3. This change.

## Open Questions

- Whether the skill should say anything about what a person does when they close the page without answering. The surface returns saying nobody answered, so the branch exists; whether it earns a sentence in the skill is best decided after the first real escalation goes through it.
- Whether `socrates findings` should stop printing the response rules where a round parks, now that the page states each decision's command. Left alone here: it is a separate question about a command's output, and the rules are stated once and rendered twice already.
