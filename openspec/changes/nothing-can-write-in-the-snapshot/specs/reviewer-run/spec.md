## MODIFIED Requirements

### Requirement: The reviewer never writes, and socrates proves it

socrates SHALL hold this with three mechanisms, named here in the order they act, each weaker than the one after it:

1. The vendor's own tool restriction, which refuses a write the model asks for. It covers the tools the model may call and nothing else, so it does not cover a tool server, a hook, or any other process the vendor's session starts.
2. The seal on the snapshot worktree, which refuses a write from anything in that session, whatever asked for it. It is what makes the property true rather than requested.
3. The comparison of the snapshot worktree against the snapshot commit when the round ends. It is what catches a write that defeated the seal.

When a round ends, socrates SHALL compare the snapshot worktree against the snapshot commit.
Any difference SHALL fail the round, naming the paths that differ, and SHALL leave the worktree in place.

A reviewer with write access will eventually decide to just fix something, after which the next round's diff contains changes nobody attributes.
The first two mechanisms can be taken away from socrates by the machine it runs on, and were: a vendor CLI can refuse to stop loading a tool server, and a filesystem can decline to enforce permissions.
The comparison cannot be taken away, which is why it stays even though the seal makes it unreachable in ordinary use.

#### Scenario: A reviewer that wrote fails the round

- **WHEN** a file in the snapshot worktree has been modified by the time the round ends, whatever defeated the seal
- **THEN** the run ends as `failed`, names the modified path, keeps the worktree, and reports no findings as a result

#### Scenario: A reviewer that created a file fails the round

- **WHEN** a file has been created in the snapshot worktree by the time the round ends, whatever defeated the seal
- **THEN** the run ends as `failed` and names the created path

#### Scenario: A clean reviewer passes the check

- **WHEN** the reviewer reads freely and writes nothing
- **THEN** the check passes and the round's verdict is whatever the findings say

#### Scenario: A write the seal refused is not a finding against the reviewer

- **WHEN** something in the reviewer's session attempted a write, the seal refused it, and the reviewer returned findings anyway
- **THEN** the worktree still matches the snapshot commit, the comparison passes, and the verdict is whatever the findings say
