.PHONY: build run clean

build:
	go build -o bin/notifier ./cmd/notifier

run: build
	./bin/notifier

clean:
	rm -rf bin/
