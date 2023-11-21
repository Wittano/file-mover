output = filebot

build:
	go build -o build/$(output) cmd/filebot/*.go

test: build
	go test ./...

clean:
	rm -r build

# TODO Install .service file
systemd:

# TODO Create installer for filebot
install:

# TODO Create uninstall for filebot
uninstall: clean


nix:
	nix build .
