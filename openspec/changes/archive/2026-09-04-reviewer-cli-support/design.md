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

## What was measured

Group 1 ran against `codex-cli 0.153.2` and `opencode 1.18.19` on 2026-09-04.
Both cleared all three criteria, so both ship in this change and neither moves to the deferred issue.

The recordings are in `internal/fidelity/criteria_test.go` behind the `record` build tag, and what each observed is committed under `internal/fidelity/testdata/criteria`.
Each states the question it asked, what it saw, and whether that establishes the criterion, so a reader can tell a measurement from a belief.

Where a criterion is about what reaches the session rather than about what the model answered, it was measured on the wire.
`codex` was pointed at a local capturing provider through `model_providers`, and the request body it sent was read.
That is a stronger observation than reading the reply, because a reply that omits an instruction does not establish that the instruction was absent from the input.

**codex keeps the repository's instructions out of the session.**
A temporary repository carried an `AGENTS.md` demanding a distinctive token at the start of every reply.
Without the override the token was in the request body codex sent.
With `-c project_doc_max_bytes=0` it was absent, and the rest of the request was otherwise the same.
The inference in openai/codex#5983 is settled for this version: zero disables the project document.

**opencode keeps the repository's instructions out of the session.**
The same repository under the same prompt.
The token was in two of three control replies and in none of three under `OPENCODE_DISABLE_PROJECT_CONFIG=1`.
The variable is documented as skipping project discovery rather than the global file, and what was measured is the half socrates depends on.

Where the observation is a model's reply rather than a request body it is sampled rather than taken once, because a model that was given an instruction does not always follow it.
A control that never read the repository establishes nothing, and the recording fails rather than reporting the criterion as met.

**opencode reads the directory it is given rather than the one it runs in.**
This was not on the list, and it is the reason the first opencode recording failed.
An invocation whose child process has its working directory set to the repository, and which states no directory, read no repository instructions at all: nought of three.
The same invocation under `--dir` read them.

That matters more than a test detail.
socrates sets the child's working directory to the sealed snapshot worktree and nothing else.
Under opencode that would have left the reviewer resolving the project from the `PWD` it inherited, which is the repository under review, so the reviewer would have read the live working tree rather than the frozen copy of it.
The opencode invocation therefore states `--dir` explicitly, and the end-to-end test that a repository's own instructions do not reach the reviewer is what holds it.

**codex reports why it failed, and its exit status does not.**
A turn that could not authenticate ended with `{"type":"turn.failed","error":{"message":...}}`, preceded by `{"type":"error","message":...}` events carrying each attempt.
The process exited 0.
So the terminal event is the account and the exit status is not, and the parser reads the stream rather than the status.

codex also substitutes its own message for some upstream failures.
A 500 from the capturing provider was reported as "We're currently experiencing high demand, which may cause temporary errors."
socrates quotes what the agent said and never matches on it, so this changes nothing socrates decides.

**opencode reports why it failed on `error` rather than on `step_finish`.**
A run naming a model that does not exist emitted `{"type":"error","error":{"name":...,"data":{"message":...}}}`, emitted no `step_finish` at all, and exited 1.
`step_finish` arrived on every successful run measured, carrying `reason`, `tokens` and `cost`.
anomalyco/opencode#26855 and #31435 did not reproduce here.
The design below is corrected: the account of a failure is on the `error` event, and the absence of the terminal event is itself the signal.

**codex has no structured output field, and a constrained answer is the whole final message.**
A successful turn was measured by answering codex from a local provider that speaks the Responses API, so the shapes come from the real binary rather than from documentation.
It writes `thread.started`, `turn.started`, `item.completed` and `turn.completed`, and under `--output-schema` the constrained JSON arrives as the `agent_message` item's `text`, whole and unfenced.
Nothing in the stream is a structured output field, and `--output-last-message` writes that same text to a file.

So the enforced tier means something slightly different under codex from under claude, and the specification's wording holds for both.
Under claude the document is a field beside the message.
Under codex the document is the message.
Either way socrates takes what the vendor constrained and never goes looking inside prose for a block, which is the property the tier exists for.

`turn.completed` carries `input_tokens`, `cached_input_tokens`, `cache_write_input_tokens`, `output_tokens` and `reasoning_output_tokens`.
The stream names no model at all, so a codex run records the model socrates asked for rather than one the vendor reported.

An `error` item inside `item.completed` is informational: a turn that carried one about missing model metadata still completed.
The terminal event decides, and `error` items and `error` events on the way do not.

Neither is the exit status the decision.
A turn that failed on a missing environment variable exited 1, and a turn that failed after exhausting retries against a 401 exited 0.

**codex sends the schema, with tools active.**
`--output-schema <file>` reached the request as `text.format` with `type: json_schema`, `strict: true`, the name `codex_output_schema`, and the supplied schema verbatim, while the tool list in the same request was non-empty.
openai/codex#15451 does not reproduce at the request level in this version.
Whether the model's answer conforms was not measured, because that needs an account.
socrates validates the document either way, so an answer that does not conform fails the call.

**`--ignore-user-config` keeps the machine's customisations out of the session.**
A `CODEX_HOME` holding a config that named a model and an MCP server was used twice.
Without the flag the request carried that model and the MCP tool family the server contributes.
With the flag it carried the built-in default model and neither.
The binary's own help states that auth still uses `CODEX_HOME`.
That half was not measured, because it needs an account.

**opencode can state a reasoning effort.**
`opencode run --variant <effort>` is documented as the provider-specific reasoning effort, and a run under `--variant high` completed.
The fallback of resolving a model and passing no effort is therefore not needed for opencode, though the routing still permits it.

**Two halves need an account and were not measured.**
Whether codex authenticates under `--ignore-user-config`, and whether a real model's answer conforms to the schema codex sends.
Both are covered by the live smoke check rather than left to inference.

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

Every one of those flags was measured against the request codex sends.
`project_doc_max_bytes=0` removes the repository's document, `--ignore-user-config` removes the machine's model and its configured tool sources, and `--output-schema` arrives as a strict `json_schema` format with tools active.
`-m` and `-c model_reasoning_effort` arrive as `model` and `reasoning.effort`.

### The opencode invocation

`opencode run` with `--format json`, `--dir <the snapshot worktree>`, `--model <provider/model>`, `--variant <effort>`, and `OPENCODE_DISABLE_PROJECT_CONFIG=1` in the environment.

`--dir` is stated because opencode does not take the repository from the process's working directory, and socrates has nowhere else to say which copy of the change is under review.

It is on the parsed tier, so the prompt carries the delivery instruction and socrates extracts from the final message.

A completed run ends with `step_finish`, carrying the stop reason and the token counts.
A failed run emits no `step_finish` and reports the account on the `error` event, at `error.data.message` beside `error.name`.
So the parser reads the terminal event where there is one, and treats its absence beside an `error` event as the failure.

`--variant` is the reasoning effort, so a role resolves to a model and an effort as it does under the other vendors.
The routing still passes no effort where one is empty.

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
What stays per-vendor is which event and which fields carry the account: codex reports `turn.failed` at `error.message` and `error` at `message`, opencode reports `error` at `error.data.message`.

Neither vendor's exit status is the account.
codex exits 0 on a turn it reported as failed, so a parser reading the status alone would call that round clean.
The stream is read in both, and the status is one more thing that can say a round did not finish.

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

**What was measured was measured on one version each, on one day.**
`codex-cli 0.153.2` and `opencode 1.18.19`, on 2026-09-04.
Both are moving fast, and a release that reinstates the repository's own instructions would take a supported vendor below the criterion silently.
The recordings and the drift check are what notice, and neither runs in continuous integration.

**Two halves rest on the vendor's own documentation rather than on a measurement.**
Whether codex authenticates under `--ignore-user-config`, which its help asserts, and whether a real model's answer conforms to the schema codex sends.
The live smoke check reaches both, and nobody has an account to run it here.

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
