**The suite must be green after every task, not only at the end of a group.**
The ordering below exists for that reason: the check is written and failing before the rename it describes, so that the rename is what makes it pass.

**This change carries nothing else.**
No task in it edits behaviour, adds a feature, or fixes something noticed along the way. The evidence that it changed no behaviour is `make verify` passing identically before and after, and one unrelated edit destroys that evidence.

**The history operations happen once and in order.**
Authorship is rewritten before the module is renamed, and the transfer happens before the first public push. Rewriting after a push means force-pushing over something somebody may already have cloned.

## 1. The check exists before the rename

- [x] 1.1 Add a convention check that no source file or live specification names a code-hosting account other than the organisation's, skipping `openspec/changes/archive` as `internal/coverage` already skips it, and verify it fails against the tree as it stands, naming files
- [x] 1.2 Add a unit test that the check fails for a file naming another account and names it. Covers `A source file naming another account fails the check`
- [x] 1.3 Add a unit test that an archived change naming the former account passes. Covers `An archived change naming the former account passes`

## 2. The rename

- [x] 2.1 Change the module path in `go.mod` to `github.com/socrates-devtool/socrates` and rewrite all 448 import references across the 181 files that carry them, and verify `make lint test` is green
- [x] 2.2 Change `MODULE` in the `Makefile`, and verify a built binary reports the same commit stamp it did before. Covers `The build reads the module path rather than repeating it`
- [x] 2.3 Verify the check from 1.1 now passes, and add an end-to-end test that every import of this project's own packages resolves under the organisation. Covers `Every import resolves under the organisation`
- [x] 2.4 Update the module path wherever `README.md` names it, and verify `make test` is green

## 3. The normative context

- [x] 3.1 Rewrite `openspec/config.yaml` so the binary it names is `socrates` with `scrat` as a second name, and the command directory is `cmd/socrates`. Covers `The context names the binary the repository builds`
- [x] 3.2 Replace the private artifact reference with a statement that `openspec/` is the sole normative source, and verify no unreachable citation remains. Covers `The context cites nothing a reader cannot reach`

## 4. What only happens once

- [x] 4.1 Verify `make verify` passes identically to the run recorded before task 2.1 began, and that the diff of this change touches paths, the context file and the new check and nothing else
- [ ] 4.2 Rewrite the commit history so every commit is authored from the personal address, before any public push, and record that it was done
- [ ] 4.3 Re-create every local clone and worktree from the rewritten history, and verify a fresh clone builds and `make verify` passes in it
- [ ] 4.4 Transfer the repository to the organisation so its issues and its URL redirects follow it, and verify the archived issue link still resolves
