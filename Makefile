.PHONY: clean build test run

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=go-webscraper
BINARY_UNIX=$(BINARY_NAME)_unix


clean:
	$(GOCLEAN)
	rm -f ./build/$(BINARY_NAME)

build:
	GO111MODULE=on $(GOMOD) download
	GO111MODULE=on $(GOMOD) verify
	$(GOBUILD) -o ./build/$(BINARY_NAME) -v ./cmd/server/

test:
	GO111MODULE=on $(GOMOD) download
	GO111MODULE=on $(GOMOD) verify
	$(GOTEST) ./... -v

run:
	GO111MODULE=on $(GOMOD) download
	GO111MODULE=on $(GOMOD) verify
	$(GOBUILD) -o ./build/$(BINARY_NAME) -v ./cmd/server/
	./build/$(BINARY_NAME) -redishost=:6379 -gatewayport=4000 -grpcport=3000

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./build/$(BINARY_UNIX) ./cmd/server
