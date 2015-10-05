# Set an output prefix, which is the local directory if not specified
PREFIX?=$(shell pwd)

# Populate version variables
# Add to compile time flags
PKG := github.com/jfrazelle/weather
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

define buildpretty
mkdir -p ${PREFIX}/cross/$(1)/$(2);
GOOS=$(1) GOARCH=$(2) CGO_ENABLED=0 go build -o ${PREFIX}/cross/$(1)/$(2)/weather -a -tags "static_build netgo" -installsuffix netgo ${GO_LDFLAGS_STATIC} .;
endef

cross: *.go VERSION
	$(foreach GOARCH,$(GOARCHS),$(foreach GOOS,$(GOOSES),$(call buildpretty,$(GOOS),$(GOARCH))))

define buildrelease
mkdir -p ${PREFIX}/cross/$(1)/$(2);
GOOS=$(1) GOARCH=$(2) CGO_ENABLED=0 go build -o ${PREFIX}/cross/weather-$(1)-$(2) -a -tags "static_build netgo" -installsuffix netgo ${GO_LDFLAGS_STATIC} .;
endef

release: *.go VERSION
	$(foreach GOARCH,$(GOARCHS),$(foreach GOOS,$(GOOSES),$(call buildrelease,$(GOOS),$(GOARCH))))

clean:
	rm -rf ${PREFIX}/cross
