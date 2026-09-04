## Why

socrates is about to be public, and the seam it reaches a reviewer through states that it is for one vendor:

> The seam exists so that the test suite can substitute the agent and so that a per-role model choice has one place to live, not so that many vendors can plug in.

That was correct while there was one vendor and no audience.
It is wrong for a public tool.
The first fact `README.md` puts in front of a reader is that socrates does not work without the `claude` CLI and the paid account behind it, and for many readers that ends the evaluation.
A person who runs a different coding agent is turned away by a coupling rather than by anything socrates decided about review.

### What the reviewer contract actually requires

The obvious bar is that the CLI returns structured output against a JSON schema, because `clean` requires a findings document that parsed.
That bar is wrong, and it is worth stating why before it is written into a specification.

socrates already validates the document itself.
`Request.Honours` takes the parsed structured output and checks it against the schema, and a violation fails the call.
A vendor's schema flag constrains what the model may emit; it does not perform socrates' check, and socrates does not skip its check when the flag is present.

So a CLI without a schema flag can still deliver a findings document.
The reviewer is asked to end its message with the document, socrates extracts and parses it, and validates it against the same schema through the same function.
A document that does not parse or does not satisfy the schema fails, and `dialogue.After` already sends that round to `failed`.
What the vendor's flag buys is reliability rather than correctness: constrained decoding means the document conforms on the first attempt more often.

That is a difference of degree, and it is a tier rather than a gate.

The reshaping SHALL NOT be done by the main coding session.
That session is the author.
A findings document the author transcribed from the reviewer's prose is a document in which the author chose the severity, chose the action, and chose what to leave out.
The asymmetry the whole design rests on is that the author has context and does the patching while the reviewer has independence, and an author who edits the findings has taken both halves.
Whatever produces the document is inside socrates and downstream of nothing the author wrote.

What cannot be recovered downstream is a review that was corrupted at the input.
If the CLI loads the reviewed repository's own instruction files, the change under review is instructing the reviewer judging it, and no reshaping of the output repairs that.
Isolation is therefore the gate, and the schema is a tier.

### Measured against that

| | claude | codex | opencode | cursor-agent | copilot | pi |
|---|---|---|---|---|---|---|
| **Repository's own instructions can be suppressed** | `--setting-sources user` | `-c project_doc_max_bytes=0` | `OPENCODE_DISABLE_PROJECT_CONFIG=1` | **none found** | **none found** | **none found** |
| Terminal event carrying the account of a failure | yes | `turn.failed` at `error.message` | `error` at `error.data.message` | yes | undeclared | not documented |
| Model selectable per invocation | yes | yes | yes | yes | yes | yes |
| Schema constrained by the vendor | `--json-schema` | `--output-schema` | no | no | no | no |
| Tier | enforced | enforced | parsed | - | - | - |

The three mechanisms in the first row were measured rather than read, against `codex-cli 0.153.2` and `opencode 1.18.19`, and `design.md` says how.
Both new vendors cleared all three criteria, so both ship here.

For the other three CLIs no such mechanism was found in their documentation, which is weaker than a finding that none exists, and settling it is one recording each rather than a redesign.
Those three are recorded as an issue to return to rather than built on a guess.

## What Changes

- The seam SHALL admit more than one vendor, and the sentence stating it exists for one SHALL be withdrawn.
- socrates SHALL state in one place which reviewer CLIs are supported and the basis on which one qualifies, in the shape `internal/platform` uses for a platform.
- The bar SHALL be: the reviewed repository's own instructions can be suppressed; a terminal event reports whether the work happened and what the agent said where it did not; and a model is selectable per role.
- A vendor constraining output to a schema SHALL be a tier rather than a criterion, because socrates validates the document either way.
- socrates SHALL obtain the findings document in one of two ways per vendor: constrained by the vendor against the supplied schema, or extracted from the reviewer's final message and validated against that same schema.
- Both tiers SHALL reach the same validation, and a document that does not parse or does not satisfy the schema SHALL fail the call in both.
- Nothing outside socrates SHALL reshape a reviewer's answer into the findings document, and the main coding session SHALL NOT be asked to.
- A configured reviewer CLI socrates does not support SHALL be refused by name, naming the guarantee it does not offer, before a snapshot is taken and before a reviewer is spawned.
- A `codex` implementation SHALL be added on the enforced tier, and an `opencode` implementation on the parsed tier.
- The schema file `codex` requires SHALL be written under socrates' home, never in the repository under review.
- Where a vendor cannot report the session's tool inventory, the run SHALL record that it did not, and SHALL NOT record the session as having held only what socrates asked for.
- The vendor SHALL be chosen in one configuration field, SHALL default to `claude`, and SHALL be recorded on the run beside the model that answered.
- Each supported vendor SHALL have its own simulator binary speaking that vendor's wire contract, its own recorded shapes, its own drift check and its own live smoke check.
- The end-to-end scenario matrix SHALL run against every supported vendor.
- What a person reads SHALL name the supported set by reading it, SHALL state which tier each vendor is on, and SHALL state which capability an unsupported CLI lacks.

## Capabilities

### New Capabilities

- `agent-support`: which reviewer CLIs socrates drives, the basis on which one qualifies, the two tiers on which a findings document arrives, and what happens when a person names a CLI socrates does not support.

### Modified Capabilities

- `agent-seam`: the seam admits more than one vendor; the claude wire contract becomes one of three stated wire contracts; the simulator and its fidelity check become per-vendor.

## Impact

- `internal/agent`: the supported set and the qualifying basis, the two tiers, `codex` and `opencode` implementations beside the `claude` one, and construction from a stated vendor.
- `internal/findings`: extraction of a findings document from a final message, reaching the existing validation rather than a second one.
- `internal/config`: the vendor field, its default, and refusal of one that is not supported.
- `internal/runstore`: the vendor and the tier recorded on the run.
- `internal/claudesim`, `cmd/claudesim`: unchanged in behaviour, joined by `internal/codexsim`, `cmd/codexsim`, `internal/opencodesim` and `cmd/opencodesim`.
- `internal/testing/simulator`: installs the simulator for the vendor under test rather than always `claude`.
- `internal/e2e`: the scenario matrix parameterised over the supported set.
- `internal/fidelity`: recorded shapes, drift and live smoke per vendor, all behind the build tag.
- `internal/convention`: a check that no second copy of the supported set exists in the source.
- `internal/prompt`: the reporting section states how the document is to be delivered on the parsed tier, generated from `internal/findings` as it already is.
- `README.md`, `internal/skill`: the supported set read rather than restated.
- `cmd/socrates/supervise.go`: the one construction site selects the implementation.
- Nothing reaches the snapshot, the dialogue's rules, the verdict rule, the surface or the page.

## Dependencies

`the-door-a-stranger-arrives-at` states the reviewer CLI and its floor at the top of `README.md` and holds it with a test.
This change makes that statement plural and makes the test read the set, so it follows that change rather than racing it.

## Deferred

`cursor-agent`, `copilot` and `pi` are deferred to https://github.com/ktei/socrates/issues/20 rather than proposed here.
Each needs one recording establishing whether the reviewed repository's own instructions can be kept out of the session.
A vendor that clears that is a new entry in the supported set on the parsed tier, which this change builds.
