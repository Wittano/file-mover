output = filebot

build:
	go build -o build/$(output) cmd/filebot/*.go

test: build
	go test ./...

clean:
	rm -r build

build-nix:
	nix build .
