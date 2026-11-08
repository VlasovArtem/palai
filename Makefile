BINARY  := bin/palai
PKG     := ./...
GOFILES := $(shell find . -name '*.go' -not -path './vendor/*')

.DEFAULT_GOAL := build
.PHONY: build run test lint fmt tidy

build: ## Build the server binary into bin/
	go build -o $(BINARY) ./cmd/palai

run: ## Run the server locally
	go run ./cmd/palai

test: ## Run all tests with the race detector
	go test -race $(PKG)

lint: ## gofmt check, go vet, golangci-lint
	@unformatted="$$(gofmt -l $(GOFILES))"; \
	if [ -n "$$unformatted" ]; then \
		echo "gofmt needed on:"; echo "$$unformatted"; exit 1; \
	fi
	go vet $(PKG)
	golangci-lint run $(PKG)

fmt: ## Format all Go files
	gofmt -w $(GOFILES)

tidy: ## Tidy go.mod / go.sum
	go mod tidy
