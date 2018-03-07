# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GODEP=dep

BINARY_NAME=ovirt-cloudprovider

IMAGE=rgolangh/ovirt-cloudprovider
REGISTRY=rgolangh
VERSION?=$(shell git describe --tags --always| cut -d "-" -f1)
RELEASE?=$(shell git describe --tags --always| cut -d "-" -f2- | sed 's/-/_/')
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

COMMON_ENV=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
COMMON_GO_BUILD_FLAGS=-a -ldflags '-extldflags "-static"'

TARBALL=$(BINARY_NAME)-$(VERSION)$(if $(RELEASE),_$(RELEASE)).tar.gz

all: clean deps build test

build:
	go env
	$(COMMON_ENV) $(GOBUILD) \
	$(COMMON_GO_BUILD_FLAGS) \
	-o $(BINARY_NAME) \
	-v cmd/$(BINARY_NAME)/*.go

container: \
	build
	quick-container

quick-container:
	docker build -t $(REGISTRY)/$(BINARY_NAME):$(VERSION) . -f deployment/container/Dockerfile
	docker tag $(REGISTRY)/$(BINARY_NAME):$(VERSION)

push:
	@docker login -u rgolangh -p ${DOCKER_BUILDER_API_KEY}
	docker push $(REGISTRY)/$(BINARY_NAME):$(VERSION)

test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)ARTIFACT_DIR
run: \
	build \
	./$(BINARY_NAME)
deps:
	go env
	# glide install --strip-vendor
	/home/rgolan/go/bin/dep ensure

rpm: 
	/bin/git archive --format=tar.gz HEAD > $(TARBALL)
	rpmbuild -tb $(TARBALL) \
		--define "debug_package %{nil}" \
		--define "_rpmdir ." \
		--define "_version $(VERSION)" \
		--define "_release $(RELEASE)"

.PHONY: build container quick-container push test clean run dep rpm
