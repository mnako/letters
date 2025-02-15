GOLANGCI_LINT_VERSION := $(shell cat .golangci-version)
GOLANGCI_LINT_ALIAS := GOPROXY=direct go run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

devcontainer:
	docker run --rm -it -v .:/src -v ./.gopath:/go -w /src golang:1.24.0-bookworm bash

format:
	go fmt
	$(GOLANGCI_LINT_ALIAS) run --fix

test:
	go test -v ./... -cover
	go vet ./...
	go mod verify

lint:
	$(GOLANGCI_LINT_ALIAS) run

.PHONY: devcontainer format test lint
