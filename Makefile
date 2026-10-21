BINARY := bin/socrates
SITE := bin/site
MODULE := github.com/socrates-tools/socrates
LINTER := github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.13.2
RELEASER := github.com/goreleaser/goreleaser/v2@v2.17.1

INSTALL_DIR = $(firstword $(shell go env GOBIN) $(shell go env GOPATH)/bin)

COMMIT = $(shell git rev-parse --short=7 HEAD 2>/dev/null)
UNCOMMITTED = $(if $(shell git status --porcelain 2>/dev/null),yes,)
VERSION = $(shell git describe --tags --exact-match 2>/dev/null)
STAMP = -X $(MODULE)/internal/version.version=$(VERSION) -X $(MODULE)/internal/version.commit=$(COMMIT) -X $(MODULE)/internal/version.uncommitted=$(UNCOMMITTED)

TARGETS := build install release release-snapshot verify test lint fmt help skill assets dev platforms site corpus-outcomes fixtures fixtures-diff live-smoke forge-smoke

.PHONY: $(TARGETS)

build:
	go build -ldflags '$(STAMP)' -o $(BINARY) ./cmd/socrates

install:
	go install -ldflags '$(STAMP)' ./cmd/socrates
	$(INSTALL_DIR)/socrates skill

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
	go run ./cmd/socrates review skill --host claude --to .claude/skills
	go run ./cmd/socrates review skill --host codex --to .agents/skills
	go run ./cmd/socrates design skill --host claude --to .claude/skills
	go run ./cmd/socrates design skill --host codex --to .agents/skills
	home=$$(mktemp -d); trap 'rm -rf "$$home"' EXIT; SOCRATES_HOME=$$home go run ./cmd/socrates implement skill --host claude --to .claude/skills --agents-to .claude/agents; SOCRATES_HOME=$$home go run ./cmd/socrates implement skill --host codex --to .agents/skills --agents-to .codex/agents

assets:
	cd agora/web && pnpm install --frozen-lockfile && node build.mjs

dev:
	go build -tags agoradev -ldflags '$(STAMP)' -o $(BINARY) ./cmd/socrates
	cd agora/web && pnpm install --frozen-lockfile && node build.mjs --watch

platforms:
	go run ./review/cmd/platforms

site:
	rm -rf $(SITE)
	mkdir -p $(SITE)/fonts
	cp site/index.html site/site.generated.css site/site.js $(SITE)/
	cp agora/web/packages/ui/tokens/tokens.css agora/web/packages/ui/tokens/fonts.css $(SITE)/
	cp agora/web/packages/ui/tokens/fonts/*.woff2 $(SITE)/fonts/
	cp agora/web/packages/ui/tokens/brand/socrates-mark.svg NOTICE $(SITE)/
	for asset in tokens.css site.generated.css site.js; do \
		fingerprint=$$(git hash-object $(SITE)/$$asset | cut -c1-12); \
		sed "s|\"$$asset\"|\"$$asset?v=$$fingerprint\"|" $(SITE)/index.html > $(SITE)/index.html.next && mv $(SITE)/index.html.next $(SITE)/index.html; \
	done

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
