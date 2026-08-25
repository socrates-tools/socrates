## Why

A round a person reads becomes a page, and the page is a dead end.
It is rendered through `artifact-design`, which uploads it to a host outside the machine, and a findings document quotes the code it was written about.
`internal/corpus` refuses to record such a document for exactly that reason, and the skill tells a session to publish one, which is the asymmetry recorded as issue 1.

The page also collects nothing, and that was correct while nothing stood behind it.
A person reads it, closes the tab, walks to the terminal and restates in prose a decision they already made with their eyes.
The session then writes that decision down in words the person did not choose.

Nothing about socrates has to change for either problem to go away.
What is missing is a thing behind the page: a local surface that serves what the session wrote, waits for the person, and hands their answer back.
This change builds that surface and nothing else.
It knows nothing about code review, findings or socrates, and rewiring the skill to it is a separate change that depends on this one.

## What Changes

- A second binary, `agora`, built from `cmd/agora` in this repository.
  It serves one document as one page, waits for a person, and prints what they answered.
  Two programs in one module, one release and one `make test`, because a second repository for a tool with one consumer is overhead with no payer.
- The binary offers `serve`, `poll`, `end`, `stop` and `version`, and nothing else.
  `serve` returns as soon as the page is reachable, `poll` blocks until the person sends, and the shape is lavish's because the shape is right.
- **The surface binds a loopback address and only a loopback address.**
  No flag widens it, no tunnel, no share command and no third-party host.
  Every asset is embedded, so the page issues no network request at all, and the document never leaves the machine.
- The session writes a document rather than a page: prose, options, and an optional fragment per decision.
  Identifiers, severities and audiences are copied out of the caller's own JSON rather than retyped, so a page cannot rename a decision to something more readable.
- **Each option carries the command it becomes.**
  `agora poll` hands the caller that command rather than a paraphrase to re-derive, so what the page shows and what the terminal runs cannot drift.
  agora composes it and never runs it: the surface writes nowhere but its own session file.
- The page holds two interactions and no others.
  A decision card, single-choice or multiple-choice, with a free-text field it always carries; and a composer with two sends, one that expects the caller to answer and re-render, and one that ends the exchange.
  There is no annotation, no whiteboard, no export and no publishing.
- Selected is shown separately from queued, because a person who cannot tell the two apart does not know what the caller will receive.
- The page renders its own evidence: syntax highlighting, a split and unified diff, and an inline SVG fragment the caller supplies.
  One vendored dependency covers the first two, and no diagram engine is embedded.
- A warm light design system with a stated floor on interface text size, and a warm dark variant beside it.
- **`TestSocratesReachesNothingNew` narrows from the module to socrates' own packages.**
  The rule does not change and the subject does. What must be true is that the reviewing tool reaches no network, not that this git repository holds no networking code.
- The suite drives the real binary against a real browserless client over a real loopback socket, because a surface whose loop nobody executes is prose.

## Capabilities

### New Capabilities
- `surface-document`: what a caller writes for the surface to render, what the surface copies rather than accepts, the command an option carries, the fragment a caller may supply and how it is neutralised, and what a document that does not validate produces.
- `review-surface`: the binary and its commands, the loopback-only bind and the guards on it, the session file and where it lives, the blocking poll and what it prints, the deadline on waiting for a person, and the idle shutdown.
- `surface-page`: what the rendered page holds, the two interactions and nothing beyond them, selected against queued, the evidence surfaces, live reload, and the floor the presentation is held to.

### Modified Capabilities
- `findings-page`: the scenario that measures what socrates reaches narrows its subject from the whole module to socrates' own binary and packages, because a second binary in the same module is not socrates gaining reach.

## Impact

- `cmd/agora`: the one new binary, and the second `main` in a repository that deliberately had one. The invariant it does not disturb is that `socrates` and `scrat` remain one program under two names.
- `internal/document`: the document contract and its validation. Imports nothing outward.
- `internal/render`: document to one HTML page. Pure, and reaches no socket.
- `internal/sanitise`: the caller's fragment neutralised against a fixed allowlist. The same job `internal/hygiene` does for a prompt, against a different threat.
- `internal/surface`: the loopback server, the session file, the long poll, the live-reload stream, the idle timeout.
- `internal/assets`: `go:embed` of a built bundle, committed and checked against its source the way the generated skill already is.
- `internal/convention`: `NetworkReach` gains the scope it is asked about, and the test that calls it names socrates' own packages.
- `web/`: the real client source and the fonts, built once by `make assets`. `go install` needs no Node; only a release does.
- `internal/e2e`: a journey that starts the real binary, answers the page over the real socket, and asserts what `poll` printed.
- **Nothing in `internal/reviewer`, `internal/runstore`, `internal/dialogue`, `internal/findings` or `internal/skill` changes.** This change adds a tool beside socrates and rewires nothing.

## Dependencies

`a-round-a-person-reads-becomes-a-page` is complete and unarchived, and `findings-page` therefore exists in that change rather than in `openspec/specs`.
The delta above has nothing to compare against until it is archived, which is the state `make scenario-drop` exists to catch.
That archive is the first thing this change waits on, and it is not this change's to perform.
