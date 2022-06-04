format:
	go fmt

test:
	go test -v ./... -cover
	go vet ./...
