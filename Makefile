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
	$(GORUN) $(PACKAGE_NAME)/cmd

push:
	docker push $(IMAGE):$(BUILDVERSION)

latest:
	docker tag "$(IMAGE):$(BUILDVERSION)" "$(IMAGE):latest"
	docker push "$(IMAGE):latest"

build:
	docker build  \
		--pull -t "$(IMAGE):$(BUILDVERSION)" \
		--file Dockerfile .
run:
	docker run  --name "${CONTAINER_NAME}-$(BUILDVERSION)" -p 3000:3000 -d "$(IMAGE):$(BUILDVERSION)"

runlive:
	docker run -v ~/.aws/credentials:/root/.aws/credentials "$(IMAGE):$(BUILDVERSION)"

stop:
	docker stop "${CONTAINER_NAME}-$(BUILDVERSION)"
	docker rm "${CONTAINER_NAME}-$(BUILDVERSION)"
