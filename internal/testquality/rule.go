package testquality

const Rule = `Never accept a test whose only evidence is that it opens, reads, greps, parses or snapshots implementation source and finds or omits particular strings, tokens, lines, commands, function names, prompt phrases, regex matches, abstract syntax shapes or incidental snapshots.
That does not prove behaviour: matching text can be dead or commented out, and a behaviour-preserving refactor can change it.
Execute a public or executable interface and assert observable behaviour, state, output, side effects and failure modes instead.

Reading a file is legitimate where the file is itself generated public output, a serialised protocol, persisted state, an intentional snapshot or another explicitly owned contract.
Name the contract, and do not use its contents as a proxy for unrelated code working.

Where a change claims to fix a bug, the test reproduces the reported failure: it fails before the fix and passes after it.`

const ReportingDuty = `Report every test this change adds whose only evidence is implementation source, as a finding like any other.
A test already in the repository that does the same is in scope where this change touches it.
Do not turn this change into a repository-wide cleanup of tests it never went near.`
