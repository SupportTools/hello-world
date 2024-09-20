.PHONY: all test build helm-package

# Variables
VERSION := $(shell date +%s)
GIT_COMMIT := $(shell git rev-parse HEAD)
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
DOCKER_IMAGE := supporttools/hello-world
CHART_VERSION := v$(shell git rev-list --count HEAD)

all: test build deploy 

test:
	@echo "Running tests and static analysis..."
	golint ./... && staticcheck ./... && go vet ./... && gosec ./... && go fmt ./...

build:
	@echo "Building backend Docker image..."
	docker buildx build \
	  --platform linux/amd64 \
	  --pull \
	  --build-arg VERSION=$(VERSION) \
	  --build-arg GIT_COMMIT=$(GIT_COMMIT) \
	  --build-arg BUILD_DATE=$(BUILD_DATE) \
	  --cache-from $(DOCKER_IMAGE):latest \
	  -t $(DOCKER_IMAGE):$(VERSION) \
	  -t $(DOCKER_IMAGE):latest \
	  --push \
	  -f Dockerfile .
