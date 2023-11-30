.PHONY: build run test

build: clean
	go build ./...

clean:
	rm go-rssaggregator

run:
	go run ./...

test:
	go test ./...
