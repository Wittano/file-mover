output = file-mover

build: test
	go build -o bin/$(output) src/*.go

test:
	go test ./test
