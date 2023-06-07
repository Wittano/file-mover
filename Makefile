output = file-mover

build: test
	go build -o bin/$(output) src/*.go

linux-test: test
	mkdir -p /tmp/test2/nested
	touch /tmp/test2/nested/test1.mp4 /tmp/test /tmp/test2/test.mp4
	$(MAKE) tests

tests:
	go test ./test
