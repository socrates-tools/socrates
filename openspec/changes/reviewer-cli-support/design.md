## Context

See `proposal.md` for why.

`internal/agent` is the only package that starts a model process, and `cmd/socrates/supervise.go` is its only construction site.
The interface takes a role, a prompt, a working directory and a schema, and answers structured output, text, the model, the effort, the session, the usage, the tools and the tool sources.
Nothing above the seam knows what a vendor is.

Three things below the seam are coupled to one vendor.
`agent.DefaultBinary` is the string `claude`.
`claudeAgent.arguments` composes that CLI's flags.
`parseStream` decodes that CLI's event shapes, and its fields reach the run's record by name.

The test suite is coupled harder than the code.
`internal/claudesim` simulates the real CLI's wire contract and is symlinked onto `PATH` as `claude`, and the specification forbids substituting at the Go interface instead.
Every end-to-end scenario runs through it.
`internal/fidelity` records what the real CLI answers, diffs the simulator against those recordings, and runs one live review, all behind a build tag.

## Goals / Non-Goals

**Goals:**

- socrates works for a person who runs `codex` or `opencode` and has no Claude subscription.
- Which CLIs are supported follows from a stated bar rather than from a maintained list.
- A person who names an unsupported CLI is told which guarantee it does not offer.
- The verdict rule, the findings contract and the dialogue are untouched by the arrival of a second and third vendor.
- Scenario coverage remains a claim about socrates rather than about one vendor's adapter.

**Non-Goals:**

- Asking the main coding session, or anything else outside socrates, to shape a reviewer's answer into findings.
- A plugin interface, a vendor registry, or configuration describing an arbitrary CLI's flags. A vendor is a file in `internal/agent` whose documentation somebody read and whose shapes somebody recorded.
- Choosing the vendor from what is on `PATH`. Which reviewer reads a person's change is theirs to state.
- Identical run records across vendors. A vendor that cannot report a tool inventory produces runs that say so.

## Decisions

### Isolation is the criterion, and the schema is a tier

The bar a CLI has to clear is what socrates cannot repair after the fact.

A findings document that arrived in the wrong shape is repairable, because socrates holds the schema and validates against it either way.
A review conducted by a session that read the reviewed repository's own instruction files is not repairable, because the defect is in the judgement rather than in the format.
The same holds for a failure the CLI does not report, since a run that fails silently is one a person cannot act on, and for a model that cannot be selected, since per-role routing exists so the reviewer gets the strongest model available.

```go
type Vendor struct {
    Name                      string
    Binary                    string
    Tier                      Tier
    IsolatedFromTheRepository bool
    ReportsWhyItFailed        bool
    ModelSelectablePerRole    bool
    ReportsItsToolInventory   bool
}

func (v Vendor) Qualifies() bool {
    return v.IsolatedFromTheRepository &&
        v.ReportsWhyItFailed &&
        v.ModelSelectablePerRole
}
```

`Tier` is outside the conjunction because both tiers reach a validated document.
`ReportsItsToolInventory` is outside it because socrates records what it observes about the tool surface and never fails a round for it, so a vendor reporting no inventory produces runs that record that.
Putting either inside would turn a recorded observation into a gate.

This lives in `internal/agent` rather than a package of its own.
`internal/platform` is its own package because a release, a workflow matrix, a refusal and a generated file read it from outside.
The vendor set is read by the agent constructor, the config validation, the refusal and two documents, and the first of those is `internal/agent`.

### The two tiers reach one validation

`Request.Honours` is already the check: it takes the structured output and validates it against the supplied schema, and a violation fails the call.
The vendor's schema flag constrains what the model may emit and does not perform that check, and socrates does not skip the check when the flag is present.

On the **enforced** tier the document is the vendor's structured output field, as today.

On the **parsed** tier the reviewer is asked to end its final message with the document, and socrates extracts it from the text and hands it to the same `Honours`.
Extraction lives in `internal/findings` beside the contract it serves, so that what shapes a document and what recognises one stay in one package.
It reads the last fenced JSON block in the final message, because a reviewer that reasons before answering will have written earlier blocks quoting the code it is discussing.

Both tiers fail the same way.
No block, a block that is not JSON, or JSON that violates the schema is a failed call, and `dialogue.After` already sends that round to `failed`.
The verdict rule is untouched: `clean` still requires a document that parsed.

What differs is how often the first attempt succeeds, which is a property to measure rather than assert.
The parsed tier's live smoke check records how many attempts a real review took.

### Nothing outside socrates shapes the document

The main coding session is the author.
A findings document the author transcribed from the reviewer's prose is one in which the author chose the severity, chose the action and chose what to omit.
`findings.EffectiveAction` defaulting a missing action to `ask-user` rather than `auto-fix` exists because an unclassified finding must never be silently patched, and an author who writes the actions has made that default unreachable.

So extraction is socrates', in Go, from the reviewer's own final message.
It is not a second model call, and it is not the session that will do the patching.

A second model call converting prose to JSON was considered and rejected.
It would be a model on the parsed tier producing unconstrained output, so it needs the same extraction and the same validation, and it adds a call that can drop a finding while removing nothing.

### The codex invocation

`codex exec` with:

- `--json` for the event stream.
- `--sandbox read-only`. This is also the `exec` default, and it is stated anyway, because a default can change in a release.
- `--output-schema <path>`, placing it on the enforced tier.
- `--ignore-user-config`, the nearest equivalent to `--safe-mode`. It skips `$CODEX_HOME/config.toml`, where MCP servers, plugins and hooks are configured.
- `-c project_doc_max_bytes=0`, to keep the repository's `AGENTS.md` out of the session.
- `-m <model>` and `-c model_reasoning_effort=<effort>` for the role's routing.

`--skip-git-repo-check` is not passed, because the snapshot worktree is a git worktree.

The schema is a file rather than a string, and it goes in the round's directory under socrates' home.
The reviewer's working directory is the sealed snapshot worktree, so writing it there is refused by the filesystem and would be wrong regardless: nothing writes in the repository under review.

Two of those flags are inferences, and group 1 of `tasks.md` settles both by recording.
`project_doc_max_bytes` is the documented cap on how much of `AGENTS.md` is read, and zero disabling it is an inference; openai/codex#5983 is an open request for a dedicated flag.
`--ignore-user-config` is documented as skipping the config file, and whether authentication survives it depends on where codex keeps credentials.

openai/codex#15451 reports `--output-schema` being silently ignored while tools are active.
socrates would meet that as a schema violation and fail the call, which is correct behaviour and a useless review.
Recorded before the adapter is trusted.

### The opencode invocation

`opencode run` with `--format json`, `--model <provider/model>`, and `OPENCODE_DISABLE_PROJECT_CONFIG=1` in the environment.

It is on the parsed tier, so the prompt carries the delivery instruction and socrates extracts from the final message.

Three things are recorded before the adapter is trusted.

`OPENCODE_DISABLE_PROJECT_CONFIG=1` is documented as skipping project `AGENTS.md` discovery and not the global file, so the recording establishes what a repository's own instructions do under it.
This is the criterion, so a negative answer takes opencode out of the supported set.

The terminal `step_finish` event is reported as dropped in containerised environments, in anomalyco/opencode#26855 and #31435.
That event is where the account of a failure lives, so a vendor that drops it does not clear `ReportsWhyItFailed`.
The recording establishes whether it arrives under socrates' invocation.

Reasoning effort per role has no documented flag on `run`, so the recording establishes whether the model variant covers it.
If it does not, opencode resolves a model and no effort, which the routing already permits, since an empty effort is not passed.

### The tool surface

claude names the session's tools and its external tool sources in the opening event, and socrates records everything it did not ask for.

Neither codex nor opencode carries such an inventory in its event stream.
codex's read-only sandbox is a stronger mechanism than a tool allowlist for the guarantee the allowlist defends, because it does not depend on the model choosing not to call a tool, but it answers a different question from the one the inventory answers.

socrates already has a rule for a session that did not report one, written for a claude release that stops emitting the field:

> Where a round completes and the stream reported no tool list, the run SHALL record that it did not.

That rule is reused rather than restated, and every completed round under either vendor records it.
An exemption stating that the sandbox makes the observation unnecessary would record a guarantee nobody measured.

### The failure account

claude's `result` event carries `is_error`, `subtype` and `result`.
The specification is detailed about preferring the account to the error stream, about not reporting the subtype where an account exists, and about saying the agent said nothing rather than quoting an empty string.

Those rules are about failures rather than about claude, so they move to the seam.
What stays per-vendor is which event and which fields carry the account: codex reports `turn.failed` and `error`, opencode reports it on `step_finish`.

### One simulator per vendor

The specification forbids a Go-level fake as the primary mechanism, because argument construction, stream parsing, process handling and reaping are where the defects are.
That reasoning applies per vendor.
A shared simulator with a vendor switch would exercise one argument construction and one parser, which are the things being tested.

`cmd/codexsim` and `cmd/opencodesim` are second and third simulators, installed on `PATH` under their vendors' names, speaking their event shapes and reading the same scenario files.
Scenario files stay vendor-neutral: a scenario is a truncated stream, a schema violation or a corpus entry, and every vendor must be able to produce each.
On the parsed tier a schema violation is a final message whose block violates the schema, and a missing document is a final message with no block, so the vocabulary covers both tiers without a per-tier scenario.

`internal/testing/simulator` gains the vendor to install and the name to install it under, and the e2e harness runs its table over `agent.Supported()`.
A vendor added to the supported set without a simulator fails the suite.

### The corpus is not re-recorded per vendor

`internal/corpus` holds what real reviewers answered, verbatim, classified by `dialogue.After`.
A recorded document replayed through a simulator reaches no network.

The corpus exists so that the scenarios stop resting on their author's beliefs about what a reviewer says, and that holds once for the tool rather than once per vendor.
`make corpus-outcomes` states the distribution and never fails for one, so adding entries from another vendor later moves a recorded number and nothing else.

## Risks / Trade-offs

**The parsed tier fails more often than the enforced tier.**
A reviewer that answers without the block produces a failed round rather than a wrong verdict, which is the correct direction.
How often is measured by the live smoke check rather than assumed, and the number is recorded.

**codex or opencode may not clear the isolation criterion.**
`project_doc_max_bytes=0` and `OPENCODE_DISABLE_PROJECT_CONFIG=1` are both established by recording in group 1, before the work is spent.
A negative answer takes that vendor out of the supported set and into the deferred issue.

**opencode may drop the event carrying the account of a failure.**
anomalyco/opencode#26855 and #31435.
That is the `ReportsWhyItFailed` criterion, and a negative recording takes opencode out of the supported set.

**Three live smoke checks need three accounts.**
None runs in continuous integration and none compiles without the build tag, so the cost falls on whoever runs `make live-smoke`.

**The scenario matrix triples the end-to-end suite's wall time.**
Accepted.
End to end is the primary mechanism, and a vendor whose scenarios are not run is a vendor nobody has tested.

**A field in a table is a weaker claim than a check.**
The criteria beside each entry are facts somebody read documentation for on a date, and CLIs move.
The live smoke checks keep the supported vendors honest.
Nothing automated watches the deferred ones.

## Migration

A person with an existing `~/.socrates/config.yaml` gets `claude` by default and nothing changes for them.
The vendor is a new optional field, so an absent one resolves the way an absent role does.
Run records written before this change name no vendor, and what reads them treats an absent vendor as `claude`.
