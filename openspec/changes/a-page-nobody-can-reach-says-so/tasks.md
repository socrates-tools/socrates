**The suite must be green after every task, not only at the end of a group.**

**socrates is untouched.**
No task adds a command, a flag, a run-record field or a network reach to `cmd/socrates` or any package it reaches, and no task changes `internal/skill`.

## 1. What serving a page says

- [x] 1.1 Say, where a page is served and nothing states the address the surface is reached at, that the page can be opened from the machine serving it and from nowhere else, naming the field that changes the address, on the stream a person reads. Covers `A page served with nothing stated says where it can be opened from`
- [x] 1.2 Verify against a live invocation that the stream a caller parses holds the address and nothing else. Covers `What a caller reads is the address and nothing else`
- [x] 1.3 Say nothing about where the page can be opened from where an operator has stated the address the surface is reached at. Covers `A stated address is an operator who has already answered`
- [x] 1.4 Verify the page is served and answerable either way, so what is said is a statement and never a refusal, and that the one condition on saying it is the stated one rather than any mark of the machine

## 2. The step that makes it work

- [x] 2.1 Give `README.md` the install step: forward the fixed port once, state the address the surface is reached at, and hand over an address naming socrates rather than the binary serving it
- [x] 2.2 Verify the port the README tells an operator to forward is the fixed default the surface binds, read from the surface. Covers `The port an operator is told to forward is the one the surface binds`
- [x] 2.3 Verify nothing in the documented address a person opens names the binary serving it

## 3. Verification

- [x] 3.1 Mutate the rule that a stated reached-at address silences the sentence and confirm the suite fails, then restore it
- [x] 3.2 Mutate the stream the sentence is written to and confirm the suite fails, then restore it
- [x] 3.3 Run `make fmt`, `make lint`, `make test`, `make scenario-coverage`, `make scenario-drop` and `make corpus-outcomes`, and verify every scenario this change states is named by a test
