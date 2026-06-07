all:
	go mod tidy
	go fmt ./...
	go test ./...
	go build .

run: all
	./ocl

watch:
	air .

test:
	go test ./...
