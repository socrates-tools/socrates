# Context map

The suite is one binary with one namespace per tool, and each tool has a vocabulary of its own.
A word means one thing inside the glossary that defines it, and the same word in another glossary is another word.

| Context | Glossary | Decisions | Instructions |
| --- | --- | --- | --- |
| The suite | [CONTEXT.md](CONTEXT.md) | [docs/adr/](docs/adr/) | [AGENTS.md](AGENTS.md) |
| review | [review/CONTEXT.md](review/CONTEXT.md) | [review/docs/adr/](review/docs/adr/) | [review/AGENTS.md](review/AGENTS.md) |
| agora | [agora/CONTEXT.md](agora/CONTEXT.md) | [agora/docs/adr/](agora/docs/adr/) | [agora/AGENTS.md](agora/AGENTS.md) |
| design | [design/CONTEXT.md](design/CONTEXT.md) | [design/docs/adr/](design/docs/adr/) | [design/AGENTS.md](design/AGENTS.md) |
| implement | [implement/CONTEXT.md](implement/CONTEXT.md) | [implement/docs/adr/](implement/docs/adr/) | [implement/AGENTS.md](implement/AGENTS.md) |

The root glossary holds the words every tool shares, and a tool's glossary holds that tool's own.
agora's glossary holds none of review's words, because agora does not know what a review is.
A word a tool shares with the suite is defined once, at the root.

## Decision records

Decision records share one number space across the root and every tool's directory, gaps and all.
A number names one document wherever it sits, so a reference to an ADR by number resolves without knowing where it was filed.
A record keeps its number when it moves, and a new record continues from the highest number across every directory.

A record sits with the tool it is about, and at the root where it is about the suite: the one binary, the command surface, one configuration file per tool, how a change and its tickets are filed, and who may start a skill.
