**Every assertion about what a caller sees is on the real binary**, driven from `internal/e2e`, never on a handler called in process.
What `serve` prints is the whole of this change's observable behaviour, and a test that reads it from anywhere but the binary's output is testing the resolution rather than the command.

**The configuration is reached by pointing `XDG_CONFIG_HOME` at a temporary directory**, as socrates' own configuration already is in its tests.
Nothing writes in the operator's real one.

**Each guard is mutated on purpose before it is trusted:** every assertion is run against the code as it stands first and expected to fail there, and the guard is put back and taken out again where the task says so.
The host guard in particular passes today for the wrong reason, so an assertion that it turns away the stated address must be shown to fail once the guard is widened.

**`make lint test` after every task, not at the end of a group.**

## 1. agora reads a file it did not have

- [ ] 1.1 Assert that with no configuration present, `agora serve` binds a loopback address the operating system chose and prints it, exactly as it does today. Establishes the baseline this change must not move
- [ ] 1.2 Assert that a configuration stating a loopback address to bind causes `agora serve`, given no `--address`, to bind that address, and confirm the assertion fails against the code as it stands. Establishes the gap
- [ ] 1.3 Read `$XDG_CONFIG_HOME/agora/config.yaml` in `internal/surface`, resolving the directory the way `internal/config` resolves socrates', returning the zero value when the file is absent and refusing a file that will not parse or that carries a field it does not know
- [ ] 1.4 Supply the `--address` default in `cmd/agora/serve.go` from what was read, falling back to `surface.DefaultAddress`, and confirm 1.1 and 1.2 both pass
- [ ] 1.5 Take the default back out and confirm 1.2 fails while 1.1 still passes

## 2. The configured bind reaches the same refusal the flag does

- [ ] 2.1 Assert that a configuration stating a bind address that is not loopback makes `agora serve` refuse, name the address, and start no server. Covers `A configured bind that is not loopback is refused too`
- [ ] 2.2 Route the configured address through `CheckLoopback` before any server starts, and confirm 2.1 passes
- [ ] 2.3 Assert that a caller passing `--address` with a different loopback address binds what the caller passed rather than what was configured. Covers `A stated flag beats a configured bind`
- [ ] 2.4 Assert the wildcard refusal a stated flag already gets is unchanged. Covers `A wildcard bind is refused rather than quietly narrowed`
- [ ] 2.5 Let the configured address bypass `CheckLoopback` and confirm 2.1 fails while 2.4 still passes, so the two paths are shown to be separately covered

## 3. serve prints the address the surface is reached at

- [ ] 3.1 Assert that a configuration stating the address the surface is reached at makes `agora serve` print that address rather than the one it bound, and confirm the assertion fails against the code as it stands. Establishes the gap against `serve prints the address stated over the address bound`
- [ ] 3.2 Resolve what to print once in `cmd/agora/serve.go`, and use it at the exit that started a server and at the exit that revised a live one, so a caller cannot tell the two apart
- [ ] 3.3 Assert the exit that finds a live surface and revises the document on it prints the same address as the exit that started one, by serving twice and comparing what each printed. Covers `serve prints the address stated over the address bound` on the second path
- [ ] 3.4 Assert that with nothing stated, both exits print the address bound. Covers `serve returns before the person has read anything`
- [ ] 3.5 Change only the exit that starts a server and confirm 3.3 fails while 3.1 passes

## 4. An address a person could not open is refused for its shape

- [ ] 4.1 Assert that a stated address which is not an absolute `http` or `https` URL makes `agora serve` refuse, name the address, and start no server. Covers `An address a person could not open is refused`
- [ ] 4.2 Check the stated address's shape where the bind address is checked, before any server starts, and confirm 4.1 passes
- [ ] 4.3 Assert that a stated address which is well shaped and which nothing serves still starts the surface and is still printed, so the check is on shape alone. Covers `An address nothing answers on is still served`
- [ ] 4.4 Add a reachability probe of the stated address, confirm 4.3 fails, and take the probe back out

## 5. The host guard is unmoved

- [ ] 5.1 Assert that with an address stated, a request naming that address is turned away before any handler runs, driven against a surface started by the real binary. Covers `The address the surface is reached at is not a host it admits`
- [ ] 5.2 Assert the existing rejection of a host the surface did not bind is unchanged. Covers `A request claiming another host is rejected`
- [ ] 5.3 Admit the stated address in `hostIsTheSurface`, confirm 5.1 fails, and take the widening back out. The guard passes 5.1 today for the wrong reason, so this is the only way 5.1 is shown to be asserting anything

## 6. What the change does not touch

- [ ] 6.1 Confirm `internal/skill` is byte-identical and `make skill` regenerates nothing, so no machine's address reaches the generated document
- [ ] 6.2 Confirm socrates' own configuration is unchanged and that `internal/config` gained no field, so the two programs stay separate. Covers `socrates keeps its two names and gains no third`
- [ ] 6.3 Confirm the surface's command list is unchanged. Covers `The surface offers nothing beyond the five`
- [ ] 6.4 Confirm the page still reaches nothing, so a stated address changed nothing about what the page requests. Covers `The page reaches nothing`

## 7. Verification

- [ ] 7.1 Run `make fmt lint test`, and confirm the suite passes with no configuration present anywhere it runs
- [ ] 7.2 Run `make scenario-coverage` and confirm every scenario this change added is named by a test
- [ ] 7.3 Run `make scenario-drop` and confirm nothing was taken out of the promoted specifications
- [ ] 7.4 Serve a document with a stated address and a bind address behind something that rewrites the host to the bound address, open the printed address in a real browser, answer, and confirm `agora poll` returns the answer. Records that the arrangement this change exists for works end to end
