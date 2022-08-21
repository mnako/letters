format:
	go fmt

lint:
	golangci-lint run

test:
	go test -v ./... -cover
	go vet ./...
