PROJECT_NAME := goss
PROJECT := github.com/lzakharov/$(PROJECT_NAME)
CONTAINER_NAME := $(PROJECT_NAME)
VERSION := $(shell cat version)
PKG_LIST := $(shell go list ./... | grep -v /vendor/)
APP_HTTP_ADDRESS ?= 127.0.0.1:8080

LDFLAGS = "-s -w -X $(PROJECT)/internal/version.Version=$(VERSION)"

build:
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -o ./bin/$(PROJECT_NAME) ./cmd/$(PROJECT_NAME)

container:
	@docker build --pull -t $(PROJECT_NAME):$(VERSION) .

run-dev:
	@docker-compose -f docker-compose.dev.yml up -d

run:
	@docker-compose up -d

test:
	@go test -v -cover --race $(PKG_LIST)

test-container:
	@docker build -f Dockerfile.test -t $(PROJECT):$(VERSION)-test . && docker run --rm test

lint:
	@golangci-lint run -v

check-health:
	curl -v $(APP_HTTP_ADDRESS)/v1/health

clean:
	@rm -rf ./bin