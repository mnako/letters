TOOL_ENV := GOWORK=$(CURDIR)/tools/go.work
GOLANGCI_LINT := $(TOOL_ENV) go tool golangci-lint

devcontainer:
	docker run --rm -it -v .:/src -v ./.gopath:/go -w /src golang:1.26.5-trixie bash

format:
	$(GOLANGCI_LINT) fmt
	$(GOLANGCI_LINT) run --fix

test:
	go test -v ./... -cover
	go vet ./...
	go mod verify

lint:
	$(GOLANGCI_LINT) fmt --diff
	$(GOLANGCI_LINT) run

.PHONY: devcontainer format test lint
