# Copyright (c) 2018 Platform9 Systems, Inc.
#
# Usage:
# make                 # builds the artifact
# make container-build # build artifact on a Linux based container using golang 1.10

SHELL := /usr/bin/env bash
BIN := nodeadm
PACKAGE_GOPATH := /go/src/github.com/platform9/$(BIN)
LDFLAGS := $(shell source ./version.sh ; KUBE_ROOT=. ; KUBE_GIT_VERSION=${VERSION_OVERRIDE} ; kube::version::ldflags)
GIT_STORAGE_MOUNT := $(shell source ./git_utils.sh; container_git_storage_mount)

.PHONY: container-build default clean $(BIN)

default: $(BIN)

container-build:
	docker run --rm -e VERSION_OVERRIDE=${VERSION_OVERRIDE} -v $(PWD):$(PACKAGE_GOPATH) $(GIT_STORAGE_MOUNT) -w $(PACKAGE_GOPATH) golang:1.10 make

clean:
	rm -f $(BIN)

$(BIN):
	go build -v -ldflags "$(LDFLAGS)"
