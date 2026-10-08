# Context

The suite is one binary, `socrates`, with one namespace per tool.
Each tool has a glossary of its own, and `CONTEXT-MAP.md` says where each one is.

This file holds the words the whole suite shares, and a tool's glossary holds that tool's own.
Use these words rather than synonyms, in a spec, an issue title, a test name or a commit message.
`docs/adr/` holds the decisions about the suite, each tool keeps its own decisions beside it, and `AGENTS.md` holds the invariants.

## The suite

- **change**: one unit of work, from the spec that opens it to the diff that lands it. Named by a slug, and `docs/changes/<slug>/` holds its spec and its tickets until it lands, when the folder is removed and git history keeps it.
- **spec**: what a design leaves behind, holding the problem, the solution, the user stories and the decisions. Committed to the change's own folder, and published nowhere else.
- **setting**: one thing a person can configure, named by its YAML key, carrying the value in force, whether that value is theirs or a default, and the default it falls back to. What a tool's `config` reports and what its `configure` offers are the same set.
- **use the default**: removing a setting's line, rather than writing today's default into the file, so that the file holds only what a person decided.
- **scenario**: the behaviour a test is named for. A task is not done until a test named for the scenario failed for the right reason and then passed.
