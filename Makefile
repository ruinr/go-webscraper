# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=go-webscraper
BINARY_UNIX=$(BINARY_NAME)_unix

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f ./build/$(BINARY_NAME)

.PHONY: build
build:
	$(GOBUILD) -o ./build/$(BINARY_NAME) -v ./cmd/server/

.PHONY: test
test:
	$(GOTEST) ./... -v

run:
	$(GOBUILD) -o ./build/$(BINARY_NAME) -v ./cmd/server/
	./build/$(BINARY_NAME) -redishost=:6379 -gatewayport=4000 -grpcport=3000

deps:
	$(GOMOD) download
	$(GOMOD) verify

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./build/$(BINARY_UNIX) ./cmd/server
