BINARY := bin/socrates
SURFACE := bin/agora
MODULE := github.com/socrates-devtool/socrates
SHORT_NAME := scrat

INSTALL_DIR = $(firstword $(shell go env GOBIN) $(shell go env GOPATH)/bin)

COMMIT = $(shell git rev-parse --short=7 HEAD 2>/dev/null)
UNCOMMITTED = $(if $(shell git status --porcelain 2>/dev/null),yes,)
STAMP = -X $(MODULE)/internal/version.commit=$(COMMIT) -X $(MODULE)/internal/version.uncommitted=$(UNCOMMITTED)
SURFACE_STAMP = -X main.commit=$(COMMIT) -X main.uncommitted=$(UNCOMMITTED)

.PHONY: build install verify test lint fmt skill assets platforms scenario-coverage scenario-drop corpus-outcomes fixtures fixtures-diff live-smoke

build:
	go build -ldflags '$(STAMP)' -o $(BINARY) ./cmd/socrates
	go build -ldflags '$(SURFACE_STAMP)' -o $(SURFACE) ./cmd/agora

install:
	go install -ldflags '$(STAMP)' ./cmd/socrates
	ln -sf socrates $(INSTALL_DIR)/$(SHORT_NAME)
	$(INSTALL_DIR)/socrates skill
	go install -ldflags '$(SURFACE_STAMP)' ./cmd/agora

verify:
	$(MAKE) lint
	$(MAKE) test
	$(MAKE) scenario-coverage
	$(MAKE) scenario-drop
	$(MAKE) corpus-outcomes
	$(MAKE) fixtures-diff

test:
	go test -race ./...

lint:
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo 'gofmt would rewrite:'; echo "$$unformatted"; exit 1; \
	fi
	go vet ./...

fmt:
	gofmt -w .

skill:
	go run ./cmd/socrates skill --to .claude/skills

assets:
	cd web && pnpm install --frozen-lockfile && node build.mjs

platforms:
	go run ./cmd/platforms

scenario-coverage:
	./scripts/scenario-coverage

scenario-drop:
	./scripts/scenario-drop

corpus-outcomes:
	./scripts/corpus-outcomes

fixtures:
	go test -tags record -count=1 -v ./internal/fidelity/

fixtures-diff:
	go test -tags drift -count=1 ./internal/fidelity/

live-smoke:
	go test -tags livesmoke -count=1 -v -timeout 20m ./internal/fidelity/
