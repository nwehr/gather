COMMIT := $(shell git rev-parse --short=8 HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_DATE := $(shell date '+%Y-%m-%d %H:%M:%S')

.PHONY: all install

all: 
	go build \
		-ldflags="-X 'main.commit=${COMMIT}' -X 'main.built=${BUILD_DATE}' -X 'main.branch=${BRANCH}'" \
		-o \
		gather main.go

install: all
	mv gather $(shell go env GOPATH)/bin/