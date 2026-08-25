## Context

socrates already decides which rounds a person reads and states that fact.
What follows that statement is a page, and today the page is made through `artifact-design`: rendered to a host outside the machine, read only, and hand-written in full by the session on every round.

Three things are wrong with that and only one of them is about design.
The document quotes the code it was written about and leaves the machine, which is the asymmetry `internal/corpus` already refuses to participate in.
The page collects nothing, so the person's decision walks back to the terminal through their memory and the session's paraphrase.
And the frame is rebuilt from nothing every round, in the session that is about to need its context for patching.

This change builds the thing that stands behind the page.
It is deliberately ignorant of socrates: it takes a document, serves it, waits, and prints what came back.
Rewiring the skill to it is `a-control-names-the-command-it-becomes`, which depends on this.

The prior art is lavish-axi, whose architecture is a local server plus a blocking poll and is the right shape.
What is not taken is its packaging, its annotation and whiteboard machinery, its export and its publishing to a third-party host.

## Goals / Non-Goals

**Goals:**

- A person answers where they are reading, and the answer reaches the caller as data rather than as recollection.
- The document never leaves the machine, and the page issues no request that does.
- What the page shows and what the terminal runs cannot drift, because they are the same string.
- A caller spends a few hundred tokens on an exchange rather than a few thousand.
- The tool is usable by any caller with a decision to put in front of a person, not only by socrates.

**Non-Goals:**

- Knowing anything about code review, findings, runs, rounds or verdicts. Those words do not appear in this binary.
- Running any command. The surface composes and returns; the caller executes.
- Annotation, whiteboards, export, sharing, or reaching any host.
- Persisting an exchange beyond its session, or being reachable from another machine under any configuration.
- Changing anything socrates does. No package under `internal/` that socrates reaches is touched except the convention check's scope.

## Decisions

**A second binary in this repository rather than its own.**
A separate repository means a second CI, a second release and a second version to keep aligned, for a tool with exactly one consumer today.
Two `main` packages in one Go module is ordinary.
The cost is that `TestSocratesReachesNothingNew` walks the module root and fails the moment `cmd/agora` imports `net/http`.
The rule it encodes is right and its subject is wrong: what must reach no network is the reviewing tool, not the git repository.
`NetworkReach` therefore takes the scope it is asked about, and the test names socrates' own packages.
Alternative considered: keeping agora in its own repository, which was the first proposal and which nobody was going to pay the overhead for.
Alternative rejected: a `socrates page` subcommand, which makes socrates render, serve and wait, contradicting the invariant that it does none of those.

**The caller writes a document; the surface renders the page.**
The alternative is lavish's model, where the caller hand-writes the HTML.
That is maximally flexible and costs the caller the frame, the palette, the layout and the diff machinery on every exchange, in exactly the context that is about to do the work.
It also makes fidelity a matter of care: nothing stops a hand-written page renaming an identifier to something friendlier.
Copying identifiers out of the caller's cited source makes that mechanical instead.
The escape hatch for what a document cannot express is one sanitised HTML fragment per decision, which is where a hand-authored illustration goes.

**The option carries its command.**
The obvious design returns which option was chosen and lets the caller work out what to do about it.
That is a second place where the mapping from option to action lives, and two copies drift.
Putting the command on the option means the page and the terminal are the same string, and it also satisfies, mechanically, the existing requirement that each decision says what answering it looks like in the terminal.
The surface never runs it, which is what keeps this tool ignorant and keeps the caller the one making the claim.

**Loopback only, with no way to widen it.**
The whole reason this exists is that the document quotes the code.
A flag that widens the bind is a flag somebody sets once and forgets, so there is none.
A wildcard address is refused outright rather than quietly narrowed, because quiet narrowing is how this repository's own environment ended up with an unreachable server and no error explaining it.
The host header is checked so a page on another origin cannot reach the API, and a per-session token stops another local process answering on the person's behalf.

**Assets are embedded, and the bundle is generated and committed.**
Highlighting and diffing come from one vendored dependency rather than two, and neither is fetched at run time, so the page renders with the network off.
`web/` holds the real client source and `make assets` builds it into `internal/assets`, which is committed and compared against its source by a test.
This is exactly the pattern `make skill` already establishes, for the same reason: a generated file somebody edits by hand and nobody regenerates fails in the way that does not fail.
`go install` needs no Node; only cutting a release does.

**Waiting is bounded, and the bound is a person's, not a process's.**
`poll` waits far longer than anything waiting on a machine would, because a person may be at lunch.
It is still bounded, and returns saying nobody answered, because an unattended caller that blocks forever has to be killed and reports nothing useful when it is.
The server idles out on its own so a forgotten surface does not hold a port.

**The page half is driven by a real browser, and the rest by a real socket.**
Half of what this change specifies is behaviour that only exists once a browser has run the client: choosing without sending, committing, taking a commitment back, the two sends, and the floor on presentation.
Asserting on the served HTML would test the renderer and claim to have tested the page, which is the shape of test this repository already refuses elsewhere.
So the end-to-end suite starts the real binary, drives a real headless browser against the real loopback address, and reads what `poll` printed on the other side.
Everything that is not about the client, the bind guards, the token, the delivery-once rule and the bounds, is driven with a plain HTTP client, because a browser there would be slower and would prove less.
The browser is a test-only dependency and appears in no build of the binary.
Alternative rejected: asserting on rendered markup, which cannot distinguish a page whose send button does nothing from one that works.

**No diagram engine.**
Embedding one is a megabyte and surrenders position, size and emphasis to a layout algorithm.
An illustration is worth drawing only where a workflow is what makes something hard to follow, which is rare, and where it is worth drawing it is worth drawing by hand.
The fragment channel carries it.

## Risks / Trade-offs

**The options frame the person's decision.** → The options are the caller's reading of the question, not the question. Three buttons make a question feel closed when it is not. The free-text field is mandatory on every decision, in both modes, and a document cannot remove it. Answering in words alone is a complete answer. This mitigates and does not eliminate; it is the real cost of the feature.

**A second binary is a second thing to keep current.** → Its document contract and any caller's source format have to agree. A contract test in this change reads a checked-in sample of the source shape it copies from. A live cross-tool test is not worth building for one consumer.

**Nobody opens the page.** → The bound on `poll` and the idle shutdown, above. The failure mode being designed out is a caller hanging forever on a tab nobody clicked.

**Bundle weight.** → Highlighting with a useful language set plus the fonts is several megabytes in the binary. Acceptable for a developer tool, and worth measuring before it is fixed. The lever, if it is needed, is trimming the language set.

**A caller's fragment is model output placed in a page that model also composed.** → Reduced to a fixed allowlist before embedding, and the page carries a strict policy that permits no outbound request. This is the same shape `internal/hygiene` exists for on the prompt side, and it is treated with the same suspicion.

**The convention check's scope narrows.** → This is the one place where a guard gets weaker, so it is done deliberately and visibly rather than by adding an exemption. The test still fails if `socrates` or anything it reaches imports the network. A scenario asserts that another binary importing the network leaves what socrates reaches unchanged, so the narrowing is itself covered rather than assumed.

## Migration Plan

Nothing to migrate. The binary is new, and no existing behaviour changes.

The one ordering constraint is external to this change: `a-round-a-person-reads-becomes-a-page` is complete and unarchived, so `findings-page` lives in that change rather than in `openspec/specs`.
The delta here has nothing to compare against until it is archived, which is the state `make scenario-drop` exists to catch.
Archiving it is the first thing this change waits on and is not this change's to perform.

## Open Questions

- The name of the fragment field, and whether one fragment per decision is enough or a document-level one is also wanted. Deferred until a real exchange wants the second.
- Whether the general-purpose mode needs a documented source-format contract for callers other than socrates, or whether "copy these fields if you cite a source" is enough. Nothing outside socrates uses it yet, so this stays unanswered until something does.
