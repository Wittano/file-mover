output = filebot

build:
	go build -o build/$(output) cmd/filebot/*.go

test: build
	go test ./...

clean:
	rm -r build

systemd: install
	cp systemd/filebot.service /etc/systemd/system/filebot.service

install: build
	cp build/filebot /usr/bin/filebot

uninstall: clean
	rm /usr/bin/filebot

nix:
	nix build .
