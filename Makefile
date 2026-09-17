BINARY := bin/socrates
MODULE := github.com/socrates-tools/socrates
LINTER := github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.13.2
RELEASER := github.com/goreleaser/goreleaser/v2@v2.17.1

INSTALL_DIR = $(firstword $(shell go env GOBIN) $(shell go env GOPATH)/bin)

COMMIT = $(shell git rev-parse --short=7 HEAD 2>/dev/null)
UNCOMMITTED = $(if $(shell git status --porcelain 2>/dev/null),yes,)
VERSION = $(shell git describe --tags --exact-match 2>/dev/null)
STAMP = -X $(MODULE)/internal/version.version=$(VERSION) -X $(MODULE)/internal/version.commit=$(COMMIT) -X $(MODULE)/internal/version.uncommitted=$(UNCOMMITTED)

TARGETS := build install release release-snapshot verify test lint fmt help skill assets platforms corpus-outcomes fixtures fixtures-diff live-smoke forge-smoke

.PHONY: $(TARGETS)

build:
	go build -ldflags '$(STAMP)' -o $(BINARY) ./cmd/socrates

install:
	go install -ldflags '$(STAMP)' ./cmd/socrates
	$(INSTALL_DIR)/socrates review skill
	$(INSTALL_DIR)/socrates design skill

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
	go run ./cmd/socrates review skill --to .claude/skills
	go run ./cmd/socrates design skill --to .claude/skills

assets:
	cd agora/web && pnpm install --frozen-lockfile && node build.mjs

platforms:
	go run ./review/cmd/platforms

corpus-outcomes:
	./scripts/corpus-outcomes

fixtures:
	go test -tags record -count=1 -v ./review/fidelity/

fixtures-diff:
	go test -tags drift -count=1 ./review/fidelity/

live-smoke:
	go test -tags livesmoke -count=1 -v -timeout 20m ./review/fidelity/

forge-smoke:
	go test -tags forgesmoke -count=1 -v -timeout 10m -run TestAForge ./review/fidelity/
