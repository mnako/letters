dev-container:
	docker run --rm -it -v .:/src -v ./.gopath:/go golang:1.21-bookworm bash

format:
	go fmt ./...

lint:
	golangci-lint run ./...

test:
	go vet ./...
	go test -v ./... -cover -coverprofile=coverage.out

cover:
	go tool cover -html=coverage.out -o cover.html

clean:
	rm cover.html coverage.out
