VERSION = v0.3.0
IMAGE_REGISTRY ?= docker.io
IMAGE_ORG ?= criticalstack
GCFLAGS = -gcflags "all=-trimpath=$(PWD)" -asmflags "all=-trimpath=$(PWD)"
GO_BUILD_ENV_VARS := GO111MODULE=on CGO_ENABLED=0

bootstrap-server:
	@cd ../../../; $(GO_BUILD_ENV_VARS) go build -o bin/bootstrap-server $(GCFLAGS)

build:
	@cd ../../../; docker build . -f build/images/bootstrap-server/Dockerfile -t $(IMAGE_REGISTRY)/$(IMAGE_ORG)/bootstrap-server:$(VERSION) --build-arg GOPROXY --build-arg GOSUMDB

push:
	@docker push $(IMAGE_REGISTRY)/$(IMAGE_ORG)/bootstrap-server:$(VERSION)
