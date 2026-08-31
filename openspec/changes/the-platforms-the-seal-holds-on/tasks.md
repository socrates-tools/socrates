**The suite must be green after every task, not only at the end of a group.**
The ordering below exists for that reason: the set is stated before anything reads it, and the refusal is proved before the message it carries is written.

**The refusal is proved end to end, not by calling the predicate.**
A test that asks the supported set whether it is supported proves nothing about whether a review refuses. Every scenario about refusing drives the real binary, as `internal/e2e` already drives it.

**This change builds no Windows support.**
A task that reaches for job objects or ACLs is in the wrong change. What is unsupported is stated and refused here, and adding support is raised as an issue.

## 1. The set is stated once

- [ ] 1.1 State the supported platforms in one place in the source, with the basis on which a platform qualifies, and verify `make lint test` is green
- [ ] 1.2 Add a unit test that the stated set names the platforms on which reaping by process group is available, and that nothing else reads a second copy of the set. Covers `The supported set is stated once`

## 2. A review refuses where it cannot keep its promise

- [ ] 2.1 Add an end-to-end test that a review on a platform where reaping by process group is unavailable refuses, spawns no reviewer and takes no snapshot, and verify it fails for the right reason today. Covers `An unsupported platform refuses the review before anything starts`
- [ ] 2.2 Add the refusal message, naming the missing guarantee and naming the way forward, and verify the convention check that every user-facing string lives in its package's `messages.go` and carries no em dash still passes
- [ ] 2.3 Raise the refusal in `internal/reviewer`, at the seam where `ProveTheSeal` is already raised, before the snapshot and before the spawn, and verify 2.1 passes
- [ ] 2.4 Add an end-to-end test that a command which neither spawns nor reaps still runs on an unsupported platform, and that a finished run's record is readable there. Covers `Reading a run is not refused on an unsupported platform`
- [ ] 2.5 Mutate the guard: make the supported set claim an unsupported platform is supported, and verify 2.1 fails, then restore

## 3. A release and the statement cannot drift

- [ ] 3.1 Add a test that fails when a platform is built by a release and not named by the stated set, and verify it fails for a target added without a statement. Covers `A platform added to a release without being stated fails the suite`
- [ ] 3.2 Extend that test to fail in the other direction, for a platform stated and built by nothing. Covers `A platform stated but never built fails the suite`

## 4. What a person reads

- [ ] 4.1 Add a supported-platforms section to `README.md` naming the platforms and pointing Windows users at WSL, reading the set from the source rather than writing it a second time
- [ ] 4.2 Add a test that the documented platforms and the stated set agree, in the same shape as the existing check that the documented port is read from the surface rather than written twice. Covers `What a person reads matches what the source states`
