BINARY := bin/scrat

.PHONY: build test lint fmt

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
