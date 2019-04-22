FROM golang:1.12.3-alpine3.9 as build-env

RUN apk update
RUN apk add git make

WORKDIR /
COPY / $GOPATH/github.com/rnidev/go-webscraper

COPY go.mod .
COPY go.sum .
# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /usr/bin/go-webscraper ./cmd/server
FROM alpine
COPY --from=build-env /usr/bin/go-webscraper /usr/bin/go-webscraper

CMD ["/usr/bin/go-webscraper","-redishost=redis:6379","-gatewayport=4000","-grpcport=3000"]
