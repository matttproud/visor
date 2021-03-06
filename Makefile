VERSION=$$(cat VERSION)

GOPATH?=$(PWD)
GOBIN?=$(GOPATH)/bin
GOFLAGS=-v -x -ldflags "-X main.VERSION_STRING $(VERSION)"

compile:
	GOPATH=$(GOPATH) go get $(GOFLAGS) -d ./cmd/visor
	GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFLAGS) ./cmd/visor

install:
	GOBIN=$(DESTDIR)/usr/local/bin $(MAKE)

########## packaging
DEB_NAME=visor
DEB_URL=http://github.com/soundcloud/visor
DEB_VERSION=$(VERSION)
DEB_DESCRIPTION=A command line interface for visor
DEB_MAINTAINER=Daniel Bornkessel <daniel@soundcloud.com>

include deb.mk

debroot:
	GOBIN=$(DEB_ROOT)/usr/bin $(MAKE)

##########
build: clean debroot debbuild

clean: debclean
	rm -rf bin src pkg $(DEB_ROOT)
