BINARY := bin/socrates-review
SURFACE := bin/agora
MODULE := github.com/socrates-tools/socrates-review
LINTER := github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.13.2
RELEASER := github.com/goreleaser/goreleaser/v2@v2.17.1

INSTALL_DIR = $(firstword $(shell go env GOBIN) $(shell go env GOPATH)/bin)

COMMIT = $(shell git rev-parse --short=7 HEAD 2>/dev/null)
UNCOMMITTED = $(if $(shell git status --porcelain 2>/dev/null),yes,)
VERSION = $(shell git describe --tags --exact-match 2>/dev/null)
STAMP = -X $(MODULE)/internal/version.version=$(VERSION) -X $(MODULE)/internal/version.commit=$(COMMIT) -X $(MODULE)/internal/version.uncommitted=$(UNCOMMITTED)
SURFACE_STAMP = -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.uncommitted=$(UNCOMMITTED)

TARGETS := build install release release-snapshot verify test lint fmt help skill assets platforms corpus-outcomes fixtures fixtures-diff live-smoke

.PHONY: $(TARGETS)

build:
	go build -ldflags '$(STAMP)' -o $(BINARY) ./cmd/socrates-review
	go build -ldflags '$(SURFACE_STAMP)' -o $(SURFACE) ./cmd/agora

install:
	go install -ldflags '$(STAMP)' ./cmd/socrates-review
	$(INSTALL_DIR)/socrates-review skill
	go install -ldflags '$(SURFACE_STAMP)' ./cmd/agora

release:
	go run $(RELEASER) release --clean

release-snapshot:
	go run $(RELEASER) release --snapshot --clean --skip=publish,sign,sbom

verify:
	$(MAKE) lint
	$(MAKE) test
	$(MAKE) corpus-outcomes
	$(MAKE) fixtures-diff

test:
	go test -race -timeout 30m ./...

lint:
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo 'gofmt would rewrite:'; echo "$$unformatted"; exit 1; \
	fi
	go vet ./...
	go run $(LINTER) run

fmt:
	gofmt -w .

help:
	@printf 'socrates targets:\n'
	@printf '  make %s\n' $(TARGETS)

skill:
	go run ./cmd/socrates-review skill --to .claude/skills

assets:
	cd web && pnpm install --frozen-lockfile && node build.mjs

platforms:
	go run ./cmd/platforms

corpus-outcomes:
	./scripts/corpus-outcomes

fixtures:
	go test -tags record -count=1 -v ./internal/fidelity/

fixtures-diff:
	go test -tags drift -count=1 ./internal/fidelity/

live-smoke:
	go test -tags livesmoke -count=1 -v -timeout 20m ./internal/fidelity/
