# socrates: design, a second workflow, and the page it holds a person on

Decisions: ADR 34, ADR 35, ADR 36, ADR 37, ADR 38, ADR 39.

## Problem Statement

A change here starts as a spec, and nothing in the suite helps anybody write one.

The loop that produces a spec today is a coding session interviewing a person in the terminal.
That loop has three faults, and all three are visible in this repository's own history.

**The questions arrive as prose, and the answers arrive as prose about prose.**
A session that wants six decisions from a person writes six numbered paragraphs into a scrollback, and the person answers by quoting numbers back.
The person is reading a wall of text to find the two questions that are theirs.
This is exactly what `socrates review` refuses to do with findings, and the refusal is already written down: a person who has to search for their own decision has been handed the terminal output again.

**A recommendation costs the same as a rejection.**
The session says which answer it would take, and the person types a reply to agree.
Agreeing with six recommendations is six sentences of typing, so the cheapest path through an interview is the one where the person stops reading and says yes to everything at once.

**What was rejected is lost by the time it is worth recording.**
The alternatives are discussed in the interview and never written down.
When the ADR is written, the alternatives are reconstructed from memory, or the Rejected section is omitted, which is why the common ADR format treats it as optional.
The next person then reopens a decision that was already settled, because nothing beside it says what it beat.

The suite already owns the machinery that fixes all three.
`socrates agora` serves one document as one page, marks the option the caller would take so agreeing costs one press, and hands back what the person chose in their own words.
Nothing uses it except review.

## Solution

`socrates design` is a second workflow, and a skill rather than a tool.

A coding session reads the skill, interviews the person over agora, and leaves a spec and its decisions behind.

    socrates design skill          install the skill
    socrates agora serve           put a round of questions in front of a person
    socrates agora poll            read what they answered

The interview is a design tree walked in rounds.
The frontier is every question whose prerequisites are settled.
A round is the whole frontier, served as one page, with the session's recommended answer already selected on each card.
The person answers the page, the session reads the answers, the tree reshapes, and the next round is served to the same page.

The session finds facts.
No question reaches the page that reading the repository would have answered.

What it leaves behind:

- `docs/changes/<slug>/spec.md`, the spec, on this repository's template.
- `docs/adr/NNNN-<sentence>.md`, one per decision that is hard to reverse, surprising without context, and the result of a real trade-off.
- `CONTEXT.md`, maintained as terms settle.

The ADR's Rejected section is written from the options the person declined.
Every alternative was offered on a card, and the reason it lost came back in the person's own words beside the choice.
A design session does not reconstruct what was considered, because it is what offered it.

## User Stories

**Somebody being interviewed**

1. As somebody being designed with, I want the questions on a page rather than in a scrollback, so that I read what is being asked of me rather than searching for it.
2. As somebody being designed with, I want the whole round at once, so that I answer an argument rather than a queue.
3. As somebody being designed with, I want the answer the session would take already selected, so that agreeing with a round costs one press.
4. As somebody being designed with, I want to take none of the offered answers and write my own, so that a question whose options are all wrong is not a trap.
5. As somebody being designed with, I want what I already settled shown above what is still open, so that I can see the argument I am standing in.
6. As somebody being designed with, I want to reopen something I settled two rounds ago, so that a design is not a ratchet.
7. As somebody being designed with, I want one address for the whole design, so that I keep one tab open rather than collecting links.
8. As somebody being designed with, I want to stop at any round and still get a spec, so that four rounds of work are not lost to a change of plan.
9. As somebody being designed with, I want to see the path a spec will be written to before it is written, so that I approve a location rather than discover one.
10. As somebody being designed with, I want to be asked nothing the session could have found by reading the repository, so that my attention is spent on decisions.
11. As somebody being designed with, I want a question about a word to be a question, so that the glossary says what I meant rather than what a model assumed.
12. As somebody being designed with, I want what a choice costs shown where seeing it is the point, so that a question about a layout is not a paragraph about a layout.

**Somebody reading the repository later**

13. As somebody reading a decision a year later, I want the alternatives beside it, so that I do not reopen a question that was already answered.
14. As somebody reading a decision a year later, I want the rejection in the words of the person who rejected it, so that I am reading a record rather than a rationalisation.
15. As somebody reading a decision a year later, I want one directory holding every decision the repository has made, so that finding one is a listing rather than a search.
16. As somebody reading a decision a year later, I want the decision to name the change it came from, so that I can read the argument around it.
17. As somebody picking up a change, I want its spec filed under the change's own name, so that finding it is knowing which change I am picking up.

**A coding session running the interview**

18. As a coding session, I want one skill that carries the whole loop, so that a design does not end in one session and resume in another.
19. As a coding session, I want the skill to name only commands the binary offers, so that an invocation I write is one that runs.
20. As a coding session, I want to compose each round fresh, so that losing a page costs an address and no work.
21. As a coding session, I want to know which rounds a person is still owed, so that I stop asking when the frontier is empty rather than when I run out of ideas.
22. As a coding session, I want to write the spec in the session that produced the understanding, so that nothing is handed to a second skill that was not there for the argument.

**Somebody installing the suite**

23. As somebody installing the suite, I want the design skill installed by the binary that runs it, so that the skill and the commands it names cannot drift apart.
24. As somebody installing the suite, I want one binary still, so that a second workflow costs me nothing on my path.

## Implementation Decisions

**`socrates design` is a skill, not a tool.**
The workflow spawns no process, freezes nothing and has nothing to supervise, so it gets no run store, no state machine and no retention policy.
`design/skill` holds the skill's source beside `review/skill`, generated, committed, and held to the same gates: a ceiling, a committed copy checked against its source, and the check that it names nothing the binary offers.
The namespace holds one operation, `socrates design skill`, which installs it.
A namespace with one operation is the honest shape, because design's machinery is agora and the session.

**The dispatcher mounts a third tree.**
`internal/dispatcher` gains `design` beside `review` and `agora`, and stays what it is: wiring, naming no vendor and nothing either tree does.

**A design document states no severity.**
The surface orders cards by severity and breaks ties in the order the document listed them, so a document that states no severity anywhere is one tie, and every card renders where the session put it.
That is what puts the settled block first, then the questions the person skipped, then the new ones.
A design question is not more or less severe than another one, it is earlier or later in the argument.

**An option may become no command.**
agora refuses an option carrying no command today, because what the page shows and what the terminal runs must be the same string.
A design answer runs nothing, so the rule gains a clause: an option may become no command, and a decision whose options become commands offers one for every option of it.
The guarantee that survives is the one that matters, that a person choosing between options is choosing between things of the same kind, and it is a property of one decision rather than of a document.

**Nothing carries across a revision, and the settled block is written rather than carried.**
Each round is a fresh document.
A settled question is a card that states the question, the answer the person gave and their words, and offers no options, which is the card review already writes for a finding it has resolved.
Identifiers are stable for the life of the design, so q2 is q2 in every round.

**The preamble carries the argument.**
One line per settled question: what it was and what was agreed.
The person reads the whole argument without scrolling, and the settled cards below it are there for the one they want to reopen.

**The page provides its own way out.**
A single-choice decision already renders a control below its options, labelled *None of these*, which clears the selection and moves the cursor into the words field.
The skill writes the real options and nothing else.

**`CONTEXT.md` is maintained during the interview.**
A term that needs writing down is written down as it settles.
A term that conflicts with an entry already there, or that renames something, is a card of its own, because the value of the glossary is that the person agreed to the word.

**A change is a folder and a decision is a file.**
`docs/changes/<slug>/` holds the spec: what is only ever about this change.
The tickets are issues, because a ticket carries a state, an assignee and a blocking edge that a committed file cannot.
`docs/adr/NNNN-<sentence>.md` holds the decisions, numbered across the repository, each carrying `Part of: <slug>`.
An ADR outlives the change that produced it, is read by people who never heard of it, and amends ADRs from other changes, so its number is repository-wide and its home is one directory.
`docs/one-command/` and `docs/monorepo/` move under `docs/changes/`.

**The Rejected section is written from what was declined.**
An option nobody was offered never appears in it.
The section is a record of the page, so it is bounded by what was actually put to the person.

**There is no round cap.**
Review caps rounds at two because a round spawns a model process against a snapshot.
A design round costs the person's attention and nothing else, so the frontier decides when to stop, and the person can end it at any round.

**A session ends on a card.**
When the frontier empties, the last round is one card: the shape of what was agreed, what will be written and the path it will be written to.
The person writes it, keeps going, or stops without writing.
Where the person ends the exchange early, the spec is written from what is settled and every unsettled question goes to Further Notes under a heading saying the design was ended early.

**Publishing is the last decision and not a second workflow.**
Where the repository states an issue tracker, the final card offers publishing the spec to it.
Cutting tickets is a later workflow and stays out.

**A lost page costs an address.**
agora holds a session while its page is open and collects it after 45 minutes with no reader, and a poll waits six hours.
Because no design state lives outside the session, a collected session is recovered by serving the accumulated document again, which mints a new address and loses no work.
The skill says so rather than implying the address is permanent.

## Testing Decisions

The seam is the delivered document and the delivered skill, which is where review's own tests sit.

**The skill is asserted, never read as source.**
`design/skill` is tested the way `review/skill` is: the rendered document is the subject, its ceiling is enforced, the committed copy must match what the package renders, and every invocation it names must exist in the binary's command surface with the flags it names.
That last check is the whole reason the skill ships from the binary, and it is `skill.Unavailable` over the built command tree, extended to the third namespace.

**The prohibitions are asserted as prohibitions.**
The skill must not write an option labelled *None of these*, must not state a severity on any decision, and must not name a command the binary does not offer.
Each is a test that fails when the skill stops holding it, following `review/skill/prohibitions_test.go`.

**The option rule is tested in the document layer, and mutated.**
agora's new clause is a `document.Validate` case: a decision mixing an option that carries a command with one that does not is refused, a decision where none carries one is accepted, and review's own documents still require one everywhere.
The guard is broken on purpose once, as the conventions require.

**The ordering claim is tested where it is made.**
That a document stating no severity renders in document order is a surface test, and it is stated as the test that the settled block stays contiguous, because that is what the claim is for.

**The loop is driven end to end through the real binary.**
`internal/e2e` drives a design: serve a round, answer it in a real browser, poll, serve the next round to the same session, end it, and assert the spec and the ADRs on disk.
Review's loop is driven this way already and the same harness serves, which is why the skill is not read to prove behaviour.

**The page half is driven by a real browser.**
Asserting on served markup would test the renderer, which is review's rule and is unchanged.

## Out of Scope

- Cutting tickets from a spec. They are a later workflow, and they land on the tracker.
- A design session store. Nothing is recorded but the spec, the ADRs and the glossary.
- A transcript of the interview. What survives a design is what it decided, not what was said on the way.
- Implementing what a spec describes. Design designs.
- Anything about review. The two workflows share the binary, the home and agora, and nothing else.
- A suite-level `socrates skill`. Each namespace installs its own.

## Further Notes

The three faults in the problem statement are each answered by something agora already does, which is why this change is mostly a skill and a clause.

The improvement over the interview loops this replaces is concentrated in one place.
Those loops treat the Rejected section of an ADR as optional because reconstructing it is expensive.
Here it is a byproduct: the alternatives were cards, the choice came back, and the reason came back beside it.
