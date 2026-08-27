# The rules the reviewer was not given

## Context

See `proposal.md` for why.

Three constraints shape everything below.

The prompt is assembled from named sections in a fixed order, and a section with nothing to say is absent rather than empty.
So the question for each new rule is not what it says but which section owns it, because that decides when it is delivered and what it moves.

The test-quality rule is one block injected into every role that judges tests.
A reviewer-only instruction cannot go inside it without reaching roles that have nowhere to report.

Tests assert on the delivered prompt and never on socrates' source, and no test may assert what a model concluded from a prompt.
So every scenario here is a statement about text that reached the reviewer, and whether these rules find more is a development-time question answered by the corpus, not by the suite.

## Goals / Non-Goals

**Goals:**

- Three rules delivered, each in the section that owns when it is delivered.
- The shared test-quality block unchanged in every role that already receives it.
- No new input, no new field, no new configuration, and no change to what a finding is or what a verdict means.

**Non-Goals:**

- Repository-scoped review guidance. It is the fourth thing the audit found and it cannot be had without overturning `2026-08-22-no-customisations-reach-the-reviewer`. Recorded as an open question, not deferred here.
- Simplification findings. `no-mistakes` asks for them because its reviewer feeds its own fixer. Here every finding costs the author a round and may cost a person, and a blocking finding that asks for nothing escalates the run immediately. Left out deliberately rather than overlooked.
- The `review_scope` field and the pipeline-delivery clause that goes with it. They exist because that reviewer runs before push, pull request and continuous integration in the same run and kept reporting that they had not happened. socrates owns no delivery, so there is nothing to carve out.
- Anything about how a person answers. An intent contradiction is an `ask-user` finding, which parks the round through the rule that already exists.

## Decisions

### The reviewer's instruction is a second block, not a longer rule

`internal/testquality` exports one block today and will export two.
The shared one keeps its text exactly, because a role that writes tests receives it and its text is asserted to be the same in each.
The second is delivered only under the reviewer prompt's test-quality section.

The alternative was appending the instruction to `Rule`.
Rejected on what it would say to the other roles: a role being asked to write a test would be told to report tests, which is either ignored or answered by a role with no findings document to answer it in.
The existing requirement already says the rule is one block used everywhere; the way to add a duty for one role is another block, and the spec now says which is which.

The bounding half of the instruction is not decoration.
Without "do not turn this into a repository-wide cleanup", the rule that finds one bad test in the change finds forty in the repository, and the reviewer spends the round on tests nobody touched.
That is the anti-astronaut failure wearing a different hat, and `no-mistakes` bounds it in the same sentence for the same reason.

### The intent obligation lives in the intent section

It goes under `What the author says they were doing`, not under `How to review`.

That is what makes it disappear correctly.
The intent section is already absent when no intent was stated, so the obligation inherits that for free and no second conditional is written.
Putting it under the reviewing rules would deliver, to every review started without an intent, an obligation about an intent that is not there.

It also puts the obligation next to the sentence it qualifies.
The section already says the stated intent is what the person says they wanted, not a statement of what the change does, and that nobody has checked one against the other.
The obligation is what to do having checked.

### The obligation fires on contradiction, and silence is not contradiction

`no-mistakes` states this against an intent that is often written as acceptance criteria with explicit required and forbidden clauses.
socrates teaches a one-liner: `socrates review --intent "count the characters of a document"`.

A rule that obliged a finding wherever the change and the intent merely fail to line up would, against a one-liner, oblige a finding about everything the line did not mention.
Every one of those is an `ask-user`, every `ask-user` parks the round, and a person gets handed a decision per unmentioned sentence.
That is worse than the silence it replaces, because it teaches people to stop reading the page.

So the obligation names two shapes and only two: a behaviour the intent names and the change does not have, and a behaviour the intent rules out and the change has.
Both are things the person wrote down.
Everything else stays what it is today, something to report if the reviewer thinks it worth reporting.

The clause that a contradiction is reported even where the change is otherwise clean is the load-bearing half.
The failure it is aimed at is not a reviewer that cannot see the contradiction, it is one that saw it, found nothing else wrong, and returned an empty findings array because the change looked fine.

The clause that it is not the author's to resolve alone follows from what the actions already mean.
`ask-user` is defined as the finding that touches functional requirements, product behaviour, or anything the author decided on purpose, and what the change is for is exactly that.
Saying it in the obligation is saying the existing definition again where it is about to be applied, which is cheaper than finding out a reviewer classified it `auto-fix`.

`no-mistakes` splits intent into authoritative and inferred, because one of its two sources is a summary inferred from the author's own agent transcript.
socrates has one source, the person typing `--intent`, and no channel for the authoring session's account by requirement.
So there is nothing to split, and the obligation is unconditional on the intent being present.

### The durable-fix rule joins the reviewing rules

It is unconditional, it is about how to read code, and it belongs beside counterexample tracing.

The two are complements rather than duplicates, which is the reason to carry both.
Counterexample tracing asks whether the code in front of you is wrong.
This asks whether the thing it claims to have fixed is still reachable somewhere it is not showing you.
A symptom patch passes the first and fails the second, and a symptom patch is what a review of a bug fix exists to catch.

It is placed after the anti-astronaut clamp rather than before it, so that the clamp is read first.
The rule asks the reviewer to look at sibling paths and shared state, which is the exact neighbourhood where an architectural preference gets mistaken for a finding.
The spec states that the rule does not weaken the clamp, and the delivered text says the path has to be one the source shows.

Not blocking a stated containment is the other half and is not optional.
A reviewer that treats every containment as a defect gets the next one written without saying it is one, and then nobody can see it at all.

### The spec's copy of the test-quality block is brought back to what is delivered

`internal/testquality.Rule` has three paragraphs and the spec quotes two.
The third, that a bug-fix change's test reproduces the reported failure, has been delivered all along and is missing from the requirement.

It is restored in the modified requirement rather than left for later, because this change adds the rule about whether a bug is still reachable, and the two are the same subject read from either end: the test proves the reported failure is gone, the review asks whether the failure is still reachable by another path.
Leaving the spec quoting two paragraphs while the change reasons about the third is how the next reader concludes the third was added here.

## Risks / Trade-offs

**More findings reach a person, and the intent obligation is the one that does it.** → The obligation is bounded to two named shapes and silence is explicitly not one of them. What it actually costs is visible rather than assumed: `make corpus-outcomes` classifies every recorded answer through `dialogue.After` and prints the distribution, so a shift toward parked rounds shows up there. If it shifts and the findings are not ones a person wanted, the wording is where to look, not the rule.

**The durable-fix rule reads as licence to demand a refactor.** → It requires a path the source shows and it sits after the clamp, which refuses inference from code shape, duplication or architectural preference. The clause about not blocking a stated containment closes the other route to the same outcome. This is the risk most worth watching in the corpus, because the failure looks like a thoughtful finding.

**Three more blocks of prompt, and the diff ceiling is unchanged.** → All three are short and none scales with the change. The ceiling governs whether the diff is inlined, not the rules, and the rules were never what pushed a prompt over it.

**Nothing here proves a reviewer finds more.** → It cannot, and the suite is forbidden from claiming it does. What the suite proves is that the text was delivered. Whether it works is a development-time evaluation against the corpus and `make live-smoke`, and that is the honest division rather than a gap.

**A rule borrowed from another tool's incidents.** → Stated in the requirement rather than smoothed over: the last of the reviewing rules was paid for elsewhere. The alternative was to wait until socrates hit each of them itself, which is a policy of learning things twice.

## Migration Plan

None.
Every round already records the prompt it delivered, so runs recorded before this change keep the prompt they were given and stay readable.
The corpus is answers rather than prompts, so nothing in it is invalidated; its outcome distribution is expected to move and `scripts/corpus-outcomes -write` is how a deliberate move is recorded.
Simulator fixtures are shapes of what the real CLI answers and do not carry the prompt, so `make fixtures` is not required by this change.
