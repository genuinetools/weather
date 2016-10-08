# Set an output prefix, which is the local directory if not specified
PREFIX?=$(shell pwd)
BUILDTAGS=

.PHONY: clean cross all fmt vet lint build test install static

# Populate version variables
# Add to compile time flags
PKG := github.com/jessfraz/weather
VERSION := $(shell cat VERSION)
GITCOMMIT := $(shell git rev-parse --short HEAD)
GITUNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no)
ifneq ($(GITUNTRACKEDCHANGES),)
	GITCOMMIT := $(GITCOMMIT)-dirty
endif
CTIMEVAR=-X $(PKG)/version.GITCOMMIT=$(GITCOMMIT) -X $(PKG)/version.VERSION=$(VERSION)
GO_LDFLAGS=-ldflags "-w $(CTIMEVAR)"
GO_LDFLAGS_STATIC=-ldflags "-w $(CTIMEVAR) -extldflags -static"
GOOSES = darwin freebsd linux windows
GOARCHS = amd64 386

all: clean build fmt lint test vet install

define buildpretty
mkdir -p ${PREFIX}/cross/$(1)/$(2);
GOOS=$(1) GOARCH=$(2) CGO_ENABLED=0 go build -o ${PREFIX}/cross/$(1)/$(2)/weather -a -tags "static_build netgo" -installsuffix netgo ${GO_LDFLAGS_STATIC} .;
endef

cross: *.go VERSION
	@echo "+ $@"
	$(foreach GOARCH,$(GOARCHS),$(foreach GOOS,$(GOOSES),$(call buildpretty,$(GOOS),$(GOARCH))))

define buildrelease
GOOS=$(1) GOARCH=$(2) CGO_ENABLED=0 go build -o ${PREFIX}/cross/weather-$(1)-$(2) -a -tags "static_build netgo" -installsuffix netgo ${GO_LDFLAGS_STATIC} .;
endef

release: *.go VERSION
	@echo "+ $@"
	$(foreach GOARCH,$(GOARCHS),$(foreach GOOS,$(GOOSES),$(call buildrelease,$(GOOS),$(GOARCH))))

clean:
	@echo "+ $@"
	@rm -rf weather
	@rm -rf ${PREFIX}/cross

build:
	@echo "+ $@"
	@go build -tags "$(BUILDTAGS)" .

static:
	@echo "+ $@"
	CGO_ENABLED=0 go build -tags "$(BUILDTAGS) static_build" -ldflags "-w -extldflags -static" -o magneto .

fmt:
	@echo "+ $@"
	@gofmt -s -l . | grep -v vendor | tee /dev/stderr

lint:
	@echo "+ $@"
	@golint ./... | grep -v vendor | tee /dev/stderr

test: fmt lint vet
	@echo "+ $@"
	@go test -v -tags "$(BUILDTAGS) cgo" $(shell go list ./... | grep -v vendor)

vet:
	@echo "+ $@"
	@go vet $(shell go list ./... | grep -v vendor)

install:
	@echo "+ $@"
	@go install .
