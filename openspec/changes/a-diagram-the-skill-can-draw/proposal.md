## Why

The bar the page is held to tells a session to use a diagram where a workflow is what makes a finding hard to follow.
The skill never says the document can carry one.

A session that reaches that line reads the document it is about to write, finds `evidence` holding source and a patch and no field for a picture, and concludes the surface cannot show one.
The slot exists and has since the surface did: `surface-document` requires that a document may carry one fragment per decision for an illustration the prose cannot carry, `surface-page` renders it, and `sanitise` reduces it to a fixed allowlist first.
The instruction and the mechanism are both already required, and nothing joins them.

The cost is not a missing feature, because nothing is missing.
It is a session that reads the bar, believes the surface cannot meet it, and either drops the diagram the bar asked for or reaches for a second tool and puts the person in front of two pages for one round.
A person deciding three findings should not be reading two surfaces, and a page that had to be assembled somewhere else is a page whose controls no longer name the commands they become.

## What Changes

- The skill says how an illustration reaches the page: one fragment per decision, beside the evidence rather than inside it.
- It says the fragment is inline SVG the session writes itself, because no diagram engine is embedded and the page renders what the caller supplied.
- It says the fragment is reduced to a fixed allowlist before it is embedded, so a session does not write script, an event handler or a remote reference and expect it to survive.
- It says this beside the bar's existing line about when to draw one, so that when and how are read together rather than a page apart.
- **The bar itself is untouched.** It states what the page owes a person, not which field carries it, and a bar that names a field is a bar that has to change when the field does.
- **socrates reaches nothing new.** The skill already names `agora serve`, `agora poll` and `--document`; this names one field of the document it already tells a session to write.

## Capabilities

### Modified Capabilities
- `agent-skill`: the skill says how an illustration reaches the page, what form it takes, and what is stripped from it before it is embedded.

## Impact

- `internal/skill`: the prose beside the bar, and the regenerated `.claude/skills/socrates/SKILL.md`.
- `internal/skill` tests: the surface prose gains the assertion that the skill names the fragment and says what becomes of it.
- No new command, no new flag, no new dependency, and nothing in `internal/document`, `internal/sanitise` or `internal/render` changes, because each already does what this teaches.
