**The suite must be green after every task.**
The harness sets the home before anything reads it, so the harness moves first and every package follows it.

**Nothing reads the old locations, at any point.**
A task that falls back to `XDG_STATE_HOME` while the rest catches up leaves a fallback nobody removes.

## 1. socrates has one home

- [ ] 1.1 Add a test that socrates reads its configuration and writes its runs under one stated home, failing for the right reason while two variables still resolve two directories. Covers `One variable moves everything`
- [ ] 1.2 Resolve socrates' home from `SOCRATES_HOME`, else a directory named for socrates under the person's own home, and verify 1.1 passes
- [ ] 1.3 Refuse a stated home that is not absolute, naming the variable and the path, and verify it does not fall back to the default. Covers `A relative home is refused`
- [ ] 1.4 Put `config.yaml` at the top of that home and `runs/` beside it, and verify a review writes its run there and reads its configuration from there. Covers `A run is found beside the configuration`
- [ ] 1.5 Have `internal/config` ask `internal/state` where the home is rather than resolving a directory of its own, and verify no second resolution rule remains in the package

## 2. agora has one home

- [ ] 2.1 Resolve agora's home from `AGORA_HOME`, else a directory named for agora under the person's own home, holding its configuration and its session state, and verify both are read and written there. Covers `One variable moves everything agora has`
- [ ] 2.2 Refuse a stated home that is not absolute, by the same rule and with the same shape of message as socrates'. Covers `A relative home is refused`
- [ ] 2.3 Verify a review with a surface served from inside it leaves socrates' home holding only socrates' state and agora's holding only agora's. Covers `Neither program writes in the other's home`

## 3. What a person is told

- [ ] 3.1 Name `SOCRATES_HOME` in the refusal a filesystem that does not enforce permissions raises, keeping the string in its package's `messages.go`, and verify the refusal names the variable that exists
- [ ] 3.2 Update `README.md` to the one path per program, and verify no user-facing string names a variable the binaries no longer read
- [ ] 3.3 Replace the path in `docs/gotchas.md` that exists on one computer with what kind of directory is wanted and how to recognise one, and verify no tracked file outside the archive names a person's home directory

## 4. The suite

- [ ] 4.1 Set one variable per program in the harness and the temporary repository helper instead of two, and verify `make test` is green
- [ ] 4.2 Mutate the refusal of a relative home and confirm the suite fails, then restore it
- [ ] 4.3 Verify a repository under review is byte-identical after a review and a surface, unchanged by the move. Covers `The repository is untouched` and `The repository is never written to`
- [ ] 4.4 Run `make fmt`, `make lint`, `make test`, `make scenario-coverage`, `make scenario-drop` and `make corpus-outcomes`, and verify every scenario this change states is named by a test
