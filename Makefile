output = filebot

build:
	go build -o build/$(output) cmd/file_mover/*.go

build-nix:
	nix build .
