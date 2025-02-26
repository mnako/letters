GOLANGCI_LINT_VERSION := $(shell cat .golangci-version)
GOLANGCI_LINT_ALIAS := GOPROXY=direct go run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
GOLINES_ALIAS := GOPROXY=direct go run github.com/segmentio/golines@v0.12.2

devcontainer:
	docker run --rm -it -v .:/src -v ./.gopath:/go -w /src golang:1.24.0-bookworm bash

format:
	go fmt
	$(GOLINES_ALIAS) -w -m80 --no-chain-split-dots **.go
	$(GOLANGCI_LINT_ALIAS) run --fix

test:
	go test -v ./... -cover
	go vet ./...
	go mod verify

lint:
	$(GOLANGCI_LINT_ALIAS) run

.PHONY: devcontainer format test lint
