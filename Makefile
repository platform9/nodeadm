# Copyright (c) 2018 Platform9 Systems, Inc.
#
# Usage:
# make                 # builds the artifact
# make container-build # build artifact on a Linux based container using golang 1.10

SHELL := /usr/bin/env bash
BIN := nodeadm
PACKAGE_GOPATH := /go/src/github.com/platform9/$(BIN)

.PHONY: container-build default clean

default: $(BIN)

container-build:
	docker run --rm -v $(PWD):$(PACKAGE_GOPATH) -w $(PACKAGE_GOPATH) golang:1.10 make

clean:
	rm -f $(BIN)

$(BIN):
	go build
