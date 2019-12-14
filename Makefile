# Go parameters
GOCMD=go
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOGEN=$(GOCMD) generate

# App parameters
GOPI=github.com/djthorpe/gopi
GOLDFLAGS += -X $(GOPI).GitTag=$(shell git describe --tags)
GOLDFLAGS += -X $(GOPI).GitBranch=$(shell git name-rev HEAD --name-only --always)
GOLDFLAGS += -X $(GOPI).GitHash=$(shell git rev-parse HEAD)
GOLDFLAGS += -X $(GOPI).GoBuildTime=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GOFLAGS = -ldflags "-s -w $(GOLDFLAGS)" 

# Prefix for installation
PREFIX=/opt/gaffer
SSLORG=mutablelogic.com

all: test install clean

install: rotel-service rotel-client rotel-ctrl
	install -m 775 -d $(PREFIX)
	install -m 775 -d $(PREFIX)/etc
	install -m 775 -d $(PREFIX)/sbin
	install -m 775 -d $(PREFIX)/bin
	install etc/rotel-service.service $(PREFIX)/etc
	install $(GOBIN)/rotel-service $(PREFIX)/sbin
	install $(GOBIN)/rotel-client $(PREFIX)/bin/rotel
	openssl req -x509 -nodes -newkey rsa:2048 \
		-keyout "${PREFIX}/etc/selfsigned.key" \
		-out "${PREFIX}/etc/selfsigned.crt" \
		-days 9999 -subj "/O=${SSLORG}"
	echo sudo useradd --system --user-group gopi
	echo sudo ln -s $(PREFIX)/etc/rotel-service.service /etc/systemd/system
	echo sudo systemctl enable rotel-service

test: protobuf
	$(GOTEST) ./...

protobuf:
	$(GOGEN) -x ./rpc/...

rotel-ctrl:
	$(GOINSTALL) $(GOFLAGS) ./cmd/rotel-ctrl/...

rotel-service: protobuf
	$(GOINSTALL) $(GOFLAGS) ./cmd/rotel-service/...

rotel-client: protobuf
	$(GOINSTALL) $(GOFLAGS) ./cmd/rotel-client/...

clean: 
	$(GOCLEAN)