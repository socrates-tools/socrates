## Context

See `proposal.md` for why.

The constraint that shapes everything here is one sentence in `openspec/config.yaml`: **No trust boundary.**
It says that only the person whose machine this is can invoke socrates on their own repository, so there is no pinned trusted default-branch SHA and no trusted-versus-pushed field split.
That sentence is still true, and it is the reason this design is not a copy of `no-mistakes`' mechanism.

It is also the sentence most likely to be misread here, so it is worth being exact about what is being defended against.
There is no adversary.
Nobody is pushing a branch that rewrites the rules gating it, because there is no push and no branch under anybody else's control.

What there is instead is **capture**.
The session that wrote the change is a writer of files, and it is exactly the party the whole design keeps away from the reviewer's ear.
socrates already refuses that session's account of what it did, its reasoning, its summary and its confidence, and it refuses the repository's agent instructions because that session could have written them in the same breath.
A file read out of the working tree is that same channel, reopened, and it would not take malice to walk through it.
A coding agent that has just written a long function and is asked to make the review pass has an obvious next move, and it is one keystroke away from being in the guide.

So the anchor is not a defence against a stranger.
It is the same asymmetry the rest of socrates is built on, applied to one more input.

## Goals / Non-Goals

**Goals:**

- A team's conventions arrive with a clone and reach every review of that repository.
- The version of those conventions that reviews a change is one no part of that change authored.
- The channel cannot grow into a general one: one file at one fixed path, one section, and one ceiling the person sets rather than the repository.
- The `agent-seam` requirement and the `config.yaml` roadmap note stop disagreeing, and the reading that survives is written down.

**Non-Goals:**

- **Guidance that can switch a rule off.** See the first decision.
- **Glob-scoped guidance.** See the fourth decision.
- **User-level guidance keyed by repository.** It is a coherent feature and it answers a different question. Nothing here forecloses it: it would be a second source composing into the same section, and this change deliberately ships one source so that the question of how two would compose is asked only if somebody wants two.
- **A verification command.** Still deferred, still a separate trust question.
- **Guidance reaching any role but the reviewer.** No other role is asked to judge the project's code against the project's conventions.

## Decisions

### The guidance adds, and can remove nothing

This is the one place the design departs from the shape the change was first sketched in, which was that a project's rules override socrates'.

Under an override, `clean` means something different in every repository, and the tool stops being able to say what a clean verdict was.
The failure mode is not hypothetical or adversarial: the first rule written under a deadline is the one that removes the finding somebody did not want, and the person who writes it is the person who most wanted it gone.
An independent reviewer that the reviewed party can quiet is not one.

Additive guidance gives up very little of what a team actually wants.
"This repository writes no comments in source" and "a test that greps implementation source is worthless here" are additions.
The things an override would buy are exemptions, and an exemption is the thing worth not having.

If an off switch is ever genuinely needed, the shape is a named list of socrates rule identifiers a repository may waive, argued for on its own, and not free prose that outranks the prompt.
That is a different change and it starts from evidence that somebody missed the waivers.

The rule is enforced by what socrates does not build rather than by asking a model to comply: no code reads a word of the guidance, so no severity, action, classification, verdict, state or exit code can be derived from it.
What is asked of the model is only the framing sentence in the section, which says the guidance adds to the rules above and withdraws none of them.
That is the honest split, and the residual is stated below rather than covered.

### The base commit is the trust anchor

`no-mistakes` reads its configuration from the default branch.
socrates has no forge and no default branch, which is where the roadmap note stopped.

It has a base commit.
Every run is measured against one, every round is snapshotted afresh against the run's original base, and the base is by construction the last state no part of the change under review authored.
That is precisely the property "trusted copy" was standing in for, and socrates has held it since the first commit.

Two consequences, both wanted:

- **A guide the change writes governs the next run, not this one.** Writing the guide and reviewing the change that writes it is an ordinary thing to do, and it works: the review runs, and the guide takes effect once it is at the base.
- **The guidance cannot move mid-dialogue.** Rounds are snapshotted against the run's *original* base, so round two is reviewed against the same conventions round one was. A rebuttal cannot arrive alongside an edit to the rules that produced the finding.

Alternatives considered:

- **Read it from the working tree.** The cheapest, and the option the issue names so that it is rejected explicitly rather than by omission. It hands the thing under review a channel to the reviewer that gates it, for no benefit beyond one run's latency.
- **Read it from the snapshot.** The same thing with extra steps: the snapshot is the working tree frozen.
- **Ask git for a merge base against some branch.** socrates has no notion of which branch is authoritative, and inventing one would import the whole default-branch model the scope rule deleted.

The cost is that the first version of a guide takes one commit to take effect.
That is stated, recorded and shown to the person rather than left to be discovered, which is the third recorded outcome below.

### One fixed path

`.socrates/review-guide.md`, relative to the repository root.

A directory rather than a single dotfile, because if a second thing is ever read out of a reviewed repository it belongs beside this one rather than adding a second root-level file.
Nothing socrates writes goes there, in that repository or any other, and the name is the same as socrates' home for the same reason both binaries share one home: one tool, one name.

Nothing configures the path.
A configurable path is a way to point socrates at a file the repository's own tooling generates, which is the working-tree read wearing a hat, and one fixed name is also how somebody arriving at a repository knows where to look.
`internal/config` already decodes with unknown fields refused, so a configuration naming a path is refused by the mechanism that is already there rather than by a new check.

The ceiling is different, and the next decision says why.

### The ceiling defaults small, is the person's to raise, and over it the guidance is dropped rather than the review

The diff ceiling is 128 KiB and it bounds the change itself.
Guidance competes with the change for the reviewer's attention, so it defaults an order of magnitude lower: **16 KiB**.
That is already more convention than most projects state usefully inside a review prompt.

It is a default rather than a fixed number, settable as `guidance_ceiling` beside `diff_ceiling`.
The round cap is the precedent for a number nothing configures, and it does not apply here: the cap bounds a dialogue that spends money and can loop, while this bounds how much of a person's own reviewer they spend on their own conventions.
It was worth being firm about while an oversized guide refused the review, because then the number decided whether reviews happened at all. Once over the ceiling only drops the guidance, the number decides nothing but attention, and that is the person's to spend.

The ceiling a run is held to is resolved once, when the run begins, and recorded on the run.
A fixed base is not on its own enough to hold the guidance still: reading configuration afresh each round would let somebody raise the ceiling while a run is parked and hand round two guidance round one was never given, with the run's own record still saying the guide was over the ceiling.
Both halves of what the reviewer is given, the text and the bound on it, therefore come from the run rather than from whatever the configuration says at the moment a round starts.

What matters is which side states it. The file is `~/.socrates/config.yaml`, outside every repository, so raising the ceiling is a decision by whoever runs the review and never by the thing being reviewed. A repository stating one reaches nothing, exactly as it reaches nothing about the diff ceiling. The anchor and the additive rule are untouched by it.

Over the ceiling, the guide is not carried at all.
Truncating would deliver part of a maintainer's rules while telling the reviewer they are the rules, which is the silent-wrong-answer shape socrates refuses everywhere it can name one.
Whole or nothing.

What is dropped is the guidance, not the review.
Refusing the review was the first shape of this decision and it was wrong: the guide is always read at the base, so an oversized guide at the base refuses every review in that repository, including the review of the change that shrinks it, until somebody commits a smaller one unreviewed.
It would also let one person's oversized guide hold up everybody else's work in progress, which is a cost paid by people who did not write it.
The ceiling exists to bound what the reviewer's attention is spent on, and that job is done by not carrying the guide.
The run records `over-the-ceiling` distinctly from `none`, and the person is told the size, the ceiling and that the next review carries it, so nothing about it is quiet.

Measured on the size the base commit's tree entry records, before anything is read, so an enormous file at that path is refused without being loaded at all and the number the person is shown is the number on disk.

### Neutralised on the existing mechanism, not a new one

The guidance is text from outside socrates entering a prompt, which is the same class of input as a finding socrates re-embeds into the next round.
`internal/hygiene` already exists for exactly that, and it is applied here unchanged.

It matters most for one thing: it strips a leading heading marker, so a guide cannot open a section of the delivered prompt.
Sections are how the reviewer tells socrates' own instruction from everything socrates is quoting to it, and a guide that could write `## What to report` would be writing socrates' text rather than the project's.

The cost is that the guide's own Markdown structure is flattened: indentation goes, nested lists become flat ones, headings become plain lines.
That is acceptable and, for this input, close to correct.
The guide is a list of things a reviewer should know, not a document anybody renders, and a second neutraliser that preserved structure would be a second mechanism to keep honest for no gain socrates can name.

### The section goes after socrates' rules and before the reporting contract

The prompt is assembled from named sections in a fixed order, and order is what says which text answers to which.
The reviewer reads socrates' rules, which are the same for every repository, then what this project adds to them, then the schema it reports in.

Putting the project's text first would let it read as the frame the rest sits inside, which is the override this design refused.
Putting it after the reporting contract would leave the last word with the repository.

### Reading a file at a commit stays in `internal/snapshot`

Every `git` invocation socrates makes is in one package.
Reading a blob at a commit is another one, so it is exported from there rather than shelled out from a new package, and `internal/guidance` owns the path, the ceiling, the neutralisation and the three outcomes.

### Three recorded outcomes, and the third is the point

The run records that guidance was carried, that there was none, or that **the change carries a guide the base commit does not**.

Without the third, the base-commit anchor is a silent surprise: somebody writes the file, runs a review, and the reviewer behaves exactly as before with nothing said.
With it, the run says what happened, a person reading the run is told, and the rule teaches itself the first time it applies.
It does not fail the round, because there is nothing wrong with what happened.

### Glob scoping is not borrowed

`no-mistakes` scopes each block of guidance to a glob and selects against the complete changed-file set.
The second half of that exists to defeat a specific attack: a contributor adding a rule's glob to `ignore_patterns` to delete that rule from their own review.

socrates has no ignore patterns, so the attack is absent and the mechanism guards nothing here.
What is left is attention management, and the 16 KiB ceiling already does that job with one moving part instead of a matcher, a selection rule and a rendering of which files matched.

If a repository's guidance is large enough that scoping would help, the ceiling has already refused it.

### The diff ceiling's requirement is untouched

`reviewer-prompt` requires that the configured diff ceiling applies to every repository and every run, with no per-repository override.
Nothing here touches it, and no amendment is needed, because a repository does not configure socrates in this design.

The distinction is worth stating because the issue asked: a configuration knob changes how socrates behaves, and a repository supplying one would be the repository steering the tool.
Guidance is content carried into one section on socrates' terms, and every knob remains user-level.

## Risks / Trade-offs

**A project tells its reviewer something wrong, and the reviewer believes it.** → Not covered, and stated rather than covered.
The guidance is additive, it decides nothing mechanically, and the section says so, but a model that is told "in this repository we do not test parsers" will weigh it.
The person who wrote that sentence owns the repository and is the person the review is for.
This is the same shape as the policy-hook residual `agent-seam` already accepts: where somebody entitled to set a rule sets one, socrates runs the review rather than not at all.

**Guidance becomes a dumping ground and crowds out the change.** → The 16 KiB ceiling and one file, both fixed. A project that outgrows it is told so rather than silently spending the reviewer's attention.

**The base-commit rule surprises somebody.** → The third recorded outcome and the line a person reads for it. The one case where the surprise persists is a guide edited in the same change as the code it is meant to govern, where the run says guidance was carried and does not say the change edits it. That is judged not worth a fourth outcome: the guidance carried is still the anchored one, which is the property that matters.

**A guide is written to defeat the additive rule in prose.** → Nothing socrates does reads it, so nothing mechanical yields. What is left is the model, and the framing text is what is offered to it. Nothing stronger is available without socrates parsing the guidance, which would be a decision derived from untrusted prose and is refused for the reason `agent-seam` refuses the same thing about vendor errors.

**`.socrates/` in a repository is mistaken for state socrates writes.** → Nothing is written there, a scenario asserts `git status` is byte-identical across a review, and the invariant is stated in `CLAUDE.md` beside the ones about the repository under review.

## Migration Plan

None required.

A repository with no `.socrates/review-guide.md` gets a delivered prompt byte-identical to today's, so every existing run, recorded corpus entry and fixture is unaffected.

The roadmap note in `openspec/config.yaml` is removed as part of this change: it is no longer deferred, and leaving it would leave the contradiction the change exists to settle.

socrates gets a guide of its own as the last task, which is the honest test of whether the channel earns its keep.
It must not restate what `internal/convention`, `make lint` or the coverage gates already enforce, because the prompt already tells the reviewer that deterministic tools own those.
What is left is what no gate here catches, and if that turns out to be nothing, that is worth knowing before other projects are told to write one.

## Open Questions

- Whether a second source of guidance, keyed under the user's own configuration rather than in the repository, is ever wanted beside this one. It would compose into the same section, and how two sources order and attribute themselves is a question worth asking only if somebody wants two.
- Whether the guidance should ever reach a role other than the reviewer. No role socrates has today needs it.
