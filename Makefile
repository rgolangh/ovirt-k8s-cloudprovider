# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GODEP=dep

BINARY_NAME=ovirt-cloudprovider

IMAGE=rgolangh/ovirt-cloudprovider
VERSION?=$(shell git describe --tags --always)
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

COMMON_ENV=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
COMMON_GO_BUILD_FLAGS=-a -ldflags '-extldflags "-static"'

all: clean deps build test

build:
	$(COMMON_ENV) $(GOBUILD) \
	$(COMMON_GO_BUILD_FLAGS) \
	-o $(BINARY_NAME) \
	-v cmd/$(BINARY_NAME)/*.go

container: \
	build
	quick-container

quick-container:
	cp $(BINARY_NAME) deployment/container
	docker build -t $(IMAGE):$(VERSION) deployment/container/

push:
    # don't forget docker login. TODO official registry
	docker push $(IMAGE):$(VERSION)

test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run: \
	build \
	./$(BINARY_NAME)
deps:
	glide install --strip-vendor

.PHONY: build