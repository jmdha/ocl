all:
	go mod tidy
	go fmt ./...
	go test ./...
	go build ./cmd/client
	go build ./cmd/server

run: all
	./server

watch:
	air .air.toml

test:
	go test ./...
