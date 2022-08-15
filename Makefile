COMMIT := $(shell git rev-parse --short=8 HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_DATE := $(shell date '+%Y-%m-%d %H:%M:%S')
GOPATH := $(shell go env GOPATH)
TARGET := gather

.PHONY: all test install

all: 
	go build \
		-ldflags="-X 'main.commit=${COMMIT}' -X 'main.built=${BUILD_DATE}' -X 'main.branch=${BRANCH}'" \
		-o \
		${TARGET} main.go

test:
	go test

install: all
	mv ${TARGET} ${GOPATH}/bin/