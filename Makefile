PACKAGE  = go-lang-web-app

init:
	go get ./...
	go get -u github.com/stretchr/testify/assert

clean:
	go clean

test:
	go test ./...

build:
	go build -o run-app /Users/jamesmillner/Developer/go-wahoo-cloud-api/cmd/main

default: build