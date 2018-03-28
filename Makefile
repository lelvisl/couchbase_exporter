PROJECT?=github.com/lelvisl/couchbase_exporter
APP?=couchbase_exporter

RELEASE?=$(shell git describe --tags 2>/dev/null || echo 0.0.1)
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date '+%Y-%m-%d_%H:%M:%S')
CONTAINER_IMAGE?=docker.io/webdeva/${APP}

GOOS?=$(shell uname | tr 'A-Z' 'a-z')
GOARCH?=amd64

clean:
	rm -f ${APP}

build: clean
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
		-ldflags "-s -w \
		-X ${PROJECT}/version.Release=${RELEASE} \
		-X ${PROJECT}/version.Commit=${COMMIT} \
		-X ${PROJECT}/version.BuildTime=${BUILD_TIME}" \
		-o ${APP}

run:
	./${APP} -node.url="http://localhost:8091" -node.auth="admin:admin"

.DEFAULT_GOAL := build
