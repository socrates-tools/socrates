# 60. The reviewer's model is never inherited from the session that called it

- Status: accepted
- Date: 2026-10-14
- Part of: issue 31

## Decision

The reviewer's model and effort are decided by the vendor's stated defaults, overridden field by field by `roles.reviewer` in the configuration file, and by nothing else.
What the calling session is running is not read, not asked for and not passed, whether or not the session and the reviewer are the same vendor.

A review is worth running because the reviewer is independent of what wrote the change.
Session isolation already removes the shared context: the reviewer inherits no session, loads none of the repository's instructions and none of the machine's customisations, and wrote none of the code.
A shared model puts back what isolation took out, as shared blind spots, and a blind spot is the one defect neither side reports.
Matching the reviewer to the author makes the reviewer more like the author, which is the opposite of what a review is for.

A person who wants the reviewer to run a particular model, including the one their session runs, says so in `roles.reviewer.model`.
That is a choice they make once, can read back with `socrates review config`, and own, rather than one made for them on every run by whatever happened to call.

## Rejected

- a flag on `review` by which the session states its own model, which the skill tells the session to pass. The value is a claim rather than a fact socrates review reads, and it puts the reviewer's routing under the control of the thing being reviewed, which is the shape socrates review refuses everywhere else.
- reading the variables a vendor's session exports, such as `CLAUDECODE` and `CLAUDE_EFFORT`. It names vendors and their variables above the seam `review/agent` exists to hold, it depends on variables no vendor documents or promises to keep, and the one measured vendor exports the effort but not the model, which is the field the question was about.
- inheriting the effort alone, because it is the half that can be read. It carries every cost of the variable reading for the field that matters less, and splits one routing into two sources that disagree about where they come from.
