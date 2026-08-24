BINARY := bin/scrat
MODULE := github.com/rui-ktei/socrates

COMMIT = $(shell git rev-parse --short=7 HEAD 2>/dev/null)
UNCOMMITTED = $(if $(shell git status --porcelain 2>/dev/null),yes,)
STAMP = -X $(MODULE)/internal/version.commit=$(COMMIT) -X $(MODULE)/internal/version.uncommitted=$(UNCOMMITTED)

.PHONY: build install test lint fmt skill scenario-coverage scenario-drop corpus-outcomes fixtures fixtures-diff live-smoke

build:
	go build -ldflags '$(STAMP)' -o $(BINARY) ./cmd/scrat

install:
	go install -ldflags '$(STAMP)' ./cmd/scrat

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
	go run ./cmd/scrat skill --to .claude/skills

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
