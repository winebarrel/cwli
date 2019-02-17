SHELL      := /bin/bash
PROGRAM    := cwli
VERSION    := v0.1.0
GOOS       := $(shell go env GOOS)
GOARCH     := $(shell go env GOARCH)
ENTRYPOINT := cmd/cwli/main.go
SRC        := $(wildcard *.go) $(wildcard */*.go)

.PHONY: all
	all: $(PROGRAM)

$(PROGRAM): $(SRC) lint vet
	go build \
	-ldflags "-X github.com/winebarrel/cwli/cli.version=$(VERSION)" \
	-o pkg/$(PROGRAM) $(ENTRYPOINT)

.PHONY: package
package: clean $(PROGRAM)
	gzip -c pkg/$(PROGRAM) > pkg/$(PROGRAM)-$(VERSION)-$(GOOS)-$(GOARCH).gz
	rm pkg/$(PROGRAM)

.PHONY: lint
lint:
	golint -set_exit_status

.PHONY: vet
vet:
	go vet

.PHONY: clean
clean:
	rm -f pkg/*
