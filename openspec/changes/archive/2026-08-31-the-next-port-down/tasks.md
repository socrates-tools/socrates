**The suite must be green after every task, not only at the end of a group.**
The ordering below exists for that reason: the window is stated and probed before anything chooses from it, and the refusals split before `serve` picks between them.

**A held port is held by a real listener.**
No task proves a port is taken by stubbing the probe. Every test that needs a held port holds one, as `internal/e2e/address_test.go` already does, and releases it.

## 1. The window is stated in one place

- [x] 1.1 State the window in `internal/surface/bind.go` as its top port and its count, compose `DefaultAddress` from the top rather than writing the port twice, and verify `make test` passes unchanged
- [x] 1.2 Add `surface.WindowAddresses()` returning the window's addresses from the top down, and verify a unit test that it holds ten addresses, the first is `DefaultAddress`, each is loopback, and the order descends by one

## 2. The refusals split

- [x] 2.1 Add a message for a window with nothing free in it, naming the window and what an operator states to pin a port, beside today's `portIsHeldByAnother`, and verify the convention check on user-facing strings still passes
- [x] 2.2 Add `surface.FirstFreeInWindow()` returning the highest free address of the window, or the new refusal where none is free, and verify a unit test against a listener the test holds on the top port that it returns the next port down

## 3. The window is walked from the top down

- [x] 3.1 Add an end-to-end test that a surface started with nothing stated, and with the top port held by a listener the test holds, serves a page on the next port down and prints that address. Covers `A held default is the next port down rather than a refusal`
- [x] 3.2 Add an end-to-end test that with the top three ports held the surface binds the fourth, and that a surface started again after it stops binds that same one. Covers `The window is walked from the top down rather than at random`
- [x] 3.3 Add an end-to-end test that with every port of the window held the surface refuses, names the window rather than one port, and starts no server. Covers `A whole window held is a refusal that names the window`
- [x] 3.4 Resolve the address in `cmd/agora/serve.go` from the window where nothing is stated, keeping `--address` at its default so an empty address a caller passed is still refused, and reading whether a port was committed to from the flag and the configuration rather than from the address itself. Verify 3.1 through 3.3 pass and `make test` is green

## 4. A stated port is pinned

- [x] 4.1 Add an end-to-end test that a caller stating a loopback address whose port is held is refused, that the refusal names that address, and that no other port was bound. Covers `A stated address that is held is refused rather than walked`
- [x] 4.2 Add an end-to-end test that a configured `surface.bind` whose port is held is refused the same way and binds no other port. Covers `A configured bind that is held is refused rather than walked`
- [x] 4.3 Add an end-to-end test that a configured `surface.reached_at` with no `surface.bind` and the top port held refuses rather than binding the next port down. Covers `A stated reached-at address pins the port too`
- [x] 4.4 Make the choice between the walk and `CheckFree` follow from whether the flag was passed, `surface.bind` is set, or `surface.reached_at` is set, and verify 4.1 through 4.3 pass

## 5. What was already true stays true

- [x] 5.1 Update the existing default-bind test to hold no port of the window while it runs, and verify it still asserts the fixed top port and still passes. Covers `The default bind is a fixed loopback port`
- [x] 5.2 Verify the existing scenarios for a wildcard bind, a configured bind that is not loopback, a stated flag beating a configured bind, and an address a person could not open still pass
- [x] 5.3 Move `internal/fidelity/livesmoke_test.go` off asserting the surface bound `DefaultAddress` exactly, serving on the resolved address instead, and verify `make live-smoke` still compiles under its build tag

## 6. What a person and an operator read

- [x] 6.1 Verify the page-opens-on-this-machine-only line names the address bound rather than the top of the window when the surface walked, driving it through the real binary. Covers `A page served with nothing stated says where it can be opened from`
- [x] 6.2 Say in `README.md` that the port moves down through the window where nothing states one, and that forwarding `18018` and stating `surface.reached_at` is what pins it. Verify the check that the documented port is read from the surface rather than written a second time still passes. Covers `The port an operator is told to forward is the one the surface binds`

## 7. Verification

- [x] 7.1 Run `make fmt`, `make lint`, `make test`, `make scenario-coverage`, `make scenario-drop` and `make corpus-outcomes`, and verify every scenario of the modified requirement is named by a test
- [x] 7.2 Mutate the guard that a stated address is refused rather than walked, by letting the walk run for a stated address, and confirm the suite fails, then restore it
- [x] 7.3 Mutate the window's order to ascend and confirm the walk test fails, then restore it
