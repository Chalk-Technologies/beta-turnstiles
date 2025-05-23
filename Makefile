# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run
CONTAINER_NAME=beta-turnstiles
IMAGE = europe-west3-docker.pkg.dev/beta-291013/beta/turnstiles
PACKAGE_NAME=beta-turnstiles
BUILDVERSION=v0.0.1

.PHONY: build test start push run stop latest test_all mysql
all: test build run ssh runlive memc

test:
	go test ./... -short

test_all:
	go test ./...

start:
	$(GOBUILD) -v ./...
	$(GORUN) $(PACKAGE_NAME)

