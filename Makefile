output = filebot

build:
	go build -o build/$(output) cmd/filebot/*.go

test: build
	go test -race ./...

clean:
	rm -r build

install: build
	cp build/filebot /usr/bin/filebot

uninstall: clean
	rm /usr/bin/filebot

nix:
	go mod vendor
	nix build .

# TODO Add win10 support
# TODO Add debian package
