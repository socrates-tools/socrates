# 24. The loopback window is stated outside the server

- Status: accepted
- Date: 2026-09-13
- Part of: socrates-review config

## Decision

`socrates-review config` reports every setting with the value in force, and each default it reports is read from the package that owns it so that a default has one home.
One of those settings is `surface.bind`, and its default is the top of agora's loopback window.

Reading it where it was stated would have made `internal/config` import `internal/surface`, and `internal/surface` imports `net` and `net/http`.
`cmd/socrates-review` reaches `internal/config`, so that import would have linked an HTTP server into socrates-review and falsified what `README.md` promises a person about where their code goes.
`convention.NetworkReach` measures `cmd/socrates-review` and every package it reaches, so it would have failed rather than let that happen quietly.

So the window moved: `internal/loopback` states the host, the ten ports counted down and the address bound where nothing states one, and reaches nothing.
`internal/surface` keeps the server and every guard that judges an address, `cmd/agora` reads the window from `internal/loopback`, and `internal/config` reads the same statement for its report.
A test in `internal/surface` holds each address of the window against `surface.CheckLoopback`, so what the two packages say about a loopback address cannot drift apart.

## Rejected

- `internal/config` stating the default itself with a test that reads `surface.DefaultAddress` and fails where the two disagree, which keeps the layout still and puts the address in the source twice; a gate holds it, but the compiler does not, and this repository states a thing once wherever one statement is possible.
- reporting `surface.bind` with no value in force, which is what `roles.reviewer.model` under `codex` honestly reads, and is a lie here because the surface does choose an address where nothing states one.
- widening what `convention.NetworkReach` allows, which trades a promise a person reads for the convenience of leaving one import alone.
