# go-webscraper

[![Go Report Card](https://goreportcard.com/badge/github.com/rnidev/go-webscraper)](https://goreportcard.com/report/github.com/rnidev/go-webscraper)
[![Build Status](https://travis-ci.org/rnidev/go-webscraper.svg?branch=master)](https://travis-ci.org/rnidev/go-webscraper)

## Live Demo

Heroku Link: https://go-webscraper.herokuapp.com/v1/amazon/product/asin/{ASIN}

Table View Sample Link: https://go-webscraper.herokuapp.com/v1/amazon/product/asin/B01644OCVS

Bullet View Sample Link: https://go-webscraper.herokuapp.com/v1/amazon/product/asin/B004QWYCVG

## Development
REST API endpoint:
```
GET /v1/amazon/product/asin/{asin}
```
Default config info
```
-redispassord=""
-redishost=:6379
-grpcport=3000
-gatewayport=4000
```

## Examples
Expected Product Info in JSON
```JSON
{
  "product": {
    "asin": "B002QYW8LW",
    "name": "Baby Banana Infant Training Toothbrush and Teether",
    "categories": [
      {
        "name": "Baby Products",
        "level": "1"
      },
      {
        "name": "Baby Care",
        "level": "2"
      },
      {
        "name": "Pacifiers, Teethers & Teething Relief",
        "level": "3"
      },
      {
        "name": "Teethers",
        "level": "4"
      }
    ],
    "ranks": [
      {
        "rank_info": "#24 in Baby ",
        "level": "1"
      },
      {
        "rank_info": "#1 in Baby Health Care Products",
        "level": "2"
      },
      {
        "rank_info": "#2 in Baby Teether Toys",
        "level": "3"
      }
    ],
    "dimensions": [
      "4.3 x 0.4 x 7.9 inches"
    ],
    "created_at": "2019-04-23T20:18:01.155542Z"
  }
}
```

Error message returned in JSON for the ASIN not found

```JSON
{
  "error": "Not Found",
  "message": "Not Found",
  "code": 2
}
```

Test REST API with curl
```
curl -H "Content-Type: application/json" -v https://go-webscraper.herokuapp.com/v1/amazon/product/asin/B004QWYCVG
```

Result

```
*   Trying 54.208.229.218...
* TCP_NODELAY set
* Connected to go-webscraper.herokuapp.com (54.208.229.218) port 443 (#0)
* ALPN, offering h2
* ALPN, offering http/1.1
* Cipher selection: ALL:!EXPORT:!EXPORT40:!EXPORT56:!aNULL:!LOW:!RC4:@STRENGTH
* successfully set certificate verify locations:
*   CAfile: /etc/ssl/cert.pem
  CApath: none
* TLSv1.2 (OUT), TLS handshake, Client hello (1):
* TLSv1.2 (IN), TLS handshake, Server hello (2):
* TLSv1.2 (IN), TLS handshake, Certificate (11):
* TLSv1.2 (IN), TLS handshake, Server key exchange (12):
* TLSv1.2 (IN), TLS handshake, Server finished (14):
* TLSv1.2 (OUT), TLS handshake, Client key exchange (16):
* TLSv1.2 (OUT), TLS change cipher, Client hello (1):
* TLSv1.2 (OUT), TLS handshake, Finished (20):
* TLSv1.2 (IN), TLS change cipher, Client hello (1):
* TLSv1.2 (IN), TLS handshake, Finished (20):
* SSL connection using TLSv1.2 / ECDHE-RSA-AES128-GCM-SHA256
* ALPN, server did not agree to a protocol
* Server certificate:
*  subject: C=US; ST=California; L=San Francisco; O=Heroku, Inc.; CN=*.herokuapp.com
*  start date: Apr 19 00:00:00 2017 GMT
*  expire date: Jun 22 12:00:00 2020 GMT
*  subjectAltName: host "go-webscraper.herokuapp.com" matched cert's "*.herokuapp.com"
*  issuer: C=US; O=DigiCert Inc; OU=www.digicert.com; CN=DigiCert SHA2 High Assurance Server CA
*  SSL certificate verify ok.
> GET /v1/amazon/product/asin/B004QWYCVG HTTP/1.1
> Host: go-webscraper.herokuapp.com
> User-Agent: curl/7.54.0
> Accept: */*
> Content-Type: application/json
>
< HTTP/1.1 200 OK
< Server: Cowboy
< Connection: keep-alive
< Content-Type: application/json
< Grpc-Metadata-Content-Type: application/grpc
< Date: Tue, 23 Apr 2019 20:45:44 GMT
< Content-Length: 590
< Via: 1.1 vegur
<
* Connection #0 to host go-webscraper.herokuapp.com left intact
{"product":{"asin":"B004QWYCVG","name":"Native Unisex Kid's Jefferson Slip-On Sneaker","categories":[{"name":"Clothing, Shoes \u0026 Jewelry","level":"1"},{"name":"Girls","level":"2"},{"name":"Shoes","level":"3"},{"name":"Flats","level":"4"}],"ranks":[{"rank_info":"#5 in Clothing, Shoes \u0026 Jewelry","level":"1"},{"rank_info":"#1 in Baby Boys' Oxfords \u0026 Loafers","level":"2"},{"rank_info":"#1 in Men's Fashion Sneakers","level":"3"},{"rank_info":"#1 in Women's Fashion Sneakers","level":"4"}],"dimensions":["8.3 x 2.8 x 2 inches"],"created_at":"2019-04-23T20:44:20.953757697Z"}}%
```

## Installation
```
  go get github.com/rnidev/go-webscaper
```

## Docker
Dockerhub link: https://hub.docker.com/r/devrni/go-webscraper
- Build from local
```
docker build -t go-webscraper .
```
- Pull from Dockerhub
```
docker pull devrni/go-webscraper
```
- Docker Compose
```
docker-compose -f ./deployments/docker-compose.yml up
```

## Make
- Build binary to ./build/
```
make build
```
- Build and Run binary locally
```
make run
```
- Run tests
```
make test
```
- Clean up tests and binary files
```
make clean
```
- Build For Linux
```
make build-linux
```

## Tech Used
[go-redis](https://github.com/go-redis/redis): Type-safe Redis client for Golang

[Colly](https://github.com/gocolly/colly): Scraper and Crawler Framework for Golang http://go-colly.org/

[grpc-go](https://github.com/grpc/grpc-go): The Go language implementation of gRPC. HTTP/2 based RPC

[grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway): gRPC to JSON proxy generator following the gRPC HTTP spec

[zap](https://github.com/uber-go/zap): Blazing fast, structured, leveled logging in Go. https://godoc.org/go.uber.org/zap

[testify](https://github.com/stretchr/testify): A toolkit with common assertions and mocks for testing

## Go version
```1.12.3```

## License
This project is licensed under the terms of the MIT license.

## References & Readings:
[Standard Project Layout for Go](https://github.com/golang-standards/project-layout), by Kyle Quest

[REST is not the Best for Micro-Services GRPC and Docker makes a compelling case](https://hackernoon.com/rest-in-peace-grpc-for-micro-service-and-grpc-for-the-web-a-how-to-908cc05e1083), by Alex Punnen

[REST vs. gRPC: Battle of the APIs](https://code.tutsplus.com/tutorials/rest-vs-grpc-battle-of-the-apis--cms-30711), by Gigi Sayfan

[Scraping the Web in Golang with Colly and Goquery](https://benjamincongdon.me/blog/2018/03/01/Scraping-the-Web-in-Golang-with-Colly-and-Goquery), by Benjamin Congdon

