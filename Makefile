COMMIT := $(shell git rev-parse --short=8 HEAD)
GOPATH := $(shell go env GOPATH)
TARGET := gather

.PHONY: all test install

all: 
	go build \
		-ldflags="-X 'main.version=${COMMIT}'" \
		-o \
		${TARGET} main.go

test:
	go test

install: all
	mv ${TARGET} ${GOPATH}/bin/