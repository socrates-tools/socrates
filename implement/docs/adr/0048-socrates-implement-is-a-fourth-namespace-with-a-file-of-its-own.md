# 48. socrates implement is a fourth namespace with a file of its own

- Status: accepted
- Date: 2026-10-01
- Part of: socrates-implement

## Decision

`socrates implement` is a fourth namespace of the one binary, holding `skill`, `config`, `configure` and `default-branch`.
It reads one optional configuration file, `~/.socrates/implement/config.yaml`, in its own directory under the brand's root, per ADR 42.

What forces a namespace is the setting.
An implementation needs a cap on how many implementers run at once, a model and an effort for them, and a remembered default branch per repository, and a person changes each the way they change a reviewer's model.
design reads and writes no configuration, which is what ADR 34 rests on, and review lands nothing, which is what its scope rests on.
So neither can hold these settings without losing the property it was defined by.

A consequence taken deliberately: the binary now carries three menus' worth of the same machinery, so the menu leaves `review/configure` for a package both build on rather than being copied.

## Rejected

- A third skill of `socrates design`, installed beside the other two, with design gaining a configuration file and `config` and `configure`.
  Declined without a word.
- A third skill of `socrates review`, with the cap a key in review's file under `socrates review configure`.
  Declined without a word.
- A fourth namespace, and one top-level `socrates configure` that is a menu over every tool's file.
  Declined without a word.
