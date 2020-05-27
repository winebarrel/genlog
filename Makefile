SHELL   := /bin/bash
VERSION := v1.0.0
GOOS      := $(shell go env GOOS)
GOARCH    := $(shell go env GOARCH)

.PHONY: all
all: build

.PHONY: build
build:
	go build -ldflags "-X main.version=$(VERSION)" ./cmd/genlog

.PHONY: package
package: clean build
	gzip genlog -c > genlog_$(VERSION)_$(GOOS)_$(GOARCH).gz
	sha1sum genlog_$(VERSION)_$(GOOS)_$(GOARCH).gz > genlog_$(VERSION)_$(GOOS)_$(GOARCH).gz.sha1sum

.PHONY: clean
clean:
	rm -f genlog
