# Context: design

`socrates design` interviews a person over agora and leaves a spec and its decisions behind, and cuts that spec into tickets.
It is a skill rather than a tool.
These are design's words, on top of the suite's shared ones in `/CONTEXT.md`.
`docs/adr/` here holds design's decisions and `AGENTS.md` its invariants.

## The interview

- **design**: one interview about one change, from its first round to the spec and the decision records it writes. The interview is the session itself, so nothing runs between one round and the next.
- **design round**: one page of questions put to a person, inside a design.
- **frontier**: every question whose prerequisites are settled. It is what one design round asks, and a design ends when it is empty.
- **question**: one thing a design needs decided, carrying the answers it offers and the one the session would take. Identified for the life of the design, so it keeps its name across rounds.
- **prototype**: throwaway code a design builds outside the repository to answer one question a card cannot settle in words. Removed once the question settles, and what the spec keeps of it is the snippet that states the decision.
- **settled**: a question the person has answered. It stays on every later page in the document's settled list, stating the answer and offering nothing.

## The tickets

- **ticket**: one tracer-bullet slice of a change, a Markdown file committed under `docs/changes/<slug>/tickets/` and named `{N}-{task}.md`. It names the tickets that block it and carries no status of its own.
- **needs a person**: what building a ticket needs that coding cannot do, such as a credential, an account, a permission or a manual step. Stated by the ticket where it has one, and asked for before the ticket is handed out.
- **blocker**: a ticket that must be done before another can start. Named by the ticket it blocks, and the only dependency a ticket states.
- **wave**: the `N` in a ticket's name, one more than the highest wave among its blockers, and 1 for a ticket with none. A wave holding one ticket names it `N`, and a wave holding several names them `N.1`, `N.2`. A wave orders the folder and gates nothing: a ticket starts once its own blockers are done. A name is fixed once the tickets are written, so abandoning a ticket renames nothing.
- **main ticket**: `README.md` in a change's tickets folder. It introduces the change to whoever implements it, points at the spec rather than repeating it, lists every ticket with its blockers, and is the one place a ticket's status is kept: `pending`, `in progress`, `done` or `abandoned`.
- **abandoned**: a ticket implementation found should not be built. Its row stays in the main ticket carrying a reason, so the table still says what was planned, and every ticket it blocked takes its blockers in its place.
