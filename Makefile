BINARY := bin/scrat

.PHONY: build test lint fmt scenario-coverage fixtures fixtures-diff live-smoke

build:
	go build -o $(BINARY) ./cmd/scrat

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

scenario-coverage:
	./scripts/scenario-coverage

fixtures:
	go test -tags record -count=1 -v ./internal/fidelity/

fixtures-diff:
	go test -tags drift -count=1 ./internal/fidelity/

live-smoke:
	go test -tags livesmoke -count=1 -v -timeout 20m ./internal/fidelity/
