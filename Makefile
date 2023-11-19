output = filebot

build:
	go build -o build/$(output) cmd/filebot/*.go

test: build
	go test ./...

build-nix:
	nix build .
