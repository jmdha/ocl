all:
	go build ./cmd/client
	go build ./cmd/server

air:
	air \
	--build.cmd "go build ./cmd/server" \
	--build.entrypoint "./server"
