dev-container:
	docker run --rm -it -v .:/src -v ./.gopath:/go golang:1.21-bookworm bash

format:
	go fmt

lint:
	golangci-lint run

test:
	go test -v ./... -cover
	go vet ./...
