output = filebot

build:
	go build -o build/$(output) cmd/filebot/*.go

build-nix:
	nix build .
