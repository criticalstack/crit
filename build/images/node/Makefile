.DEFAULT_GOAL:=help

SHELL:=/bin/bash
.ONESHELL:

##@ Building

BIN_DIR := bin
CRIT    := $(BIN_DIR)/crit
CINDER  := $(BIN_DIR)/cinder

.PHONY: $(CRIT) $(CINDER)

$(CRIT): # Build crit
	$(MAKE) -s -C ../../../ crit BIN_DIR=$(PWD)/$(BIN_DIR)

$(CINDER): # Build cinder
	$(MAKE) -s -C ../../../ cinder BIN_DIR=$(PWD)/$(BIN_DIR)


TAG   ?= latest
IMAGE ?= docker.io/criticalstack/cinder:$(TAG)

BUILD_DIR             := $(PWD)/build
BUILD_CONTAINER_NAME  := cinder-build
BASE_IMAGE            := docker.io/kindest/base:v20200709-8b878ace
KUBERNETES_VERSION    := 1.18.5

# charts
CHARTS_DIR                          := $(PWD)/charts
CILIUM_VERSION                      := 1.8.1
LOCAL_PATH_PROVISIONER_VERSION      := 0.0.12
MACHINE_API_VERSION                 := 1.0.1
MACHINE_API_PROVIDER_DOCKER_VERSION := 1.0.2

charts/*.tgz:
	mkdir -p $(CHARTS_DIR)
	helm repo add stable https://kubernetes-charts.storage.googleapis.com
	helm repo add incubator https://kubernetes-charts-incubator.storage.googleapis.com/
	helm repo add criticalstack https://charts.cscr.io/criticalstack
	helm repo add traefik https://containous.github.io/traefik-helm-chart
	helm repo add jetstack https://charts.jetstack.io
	helm repo add cilium https://helm.cilium.io/
	helm repo update
	helm pull cilium/cilium --version $(CILIUM_VERSION) -d $(CHARTS_DIR)
	helm pull criticalstack/local-path-provisioner --version $(LOCAL_PATH_PROVISIONER_VERSION) -d $(CHARTS_DIR)
	mv $(CHARTS_DIR)/cilium-$(CILIUM_VERSION).tgz $(CHARTS_DIR)/cilium.tgz
	mv $(CHARTS_DIR)/local-path-provisioner-$(LOCAL_PATH_PROVISIONER_VERSION).tgz $(CHARTS_DIR)/local-path-provisioner.tgz

build/*.tar:
	mkdir -p $(BUILD_DIR)
	for image in $$($(CINDER) get images)
	do
	    docker pull $$image
	    filename=$$(echo $$image | sed -e 's/[^A-Za-z0-9._-]/_/g').tar
	    docker save -o build/$$filename $$image
	done

.PHONY: build

build: clean $(CRIT) $(CINDER) charts/*.tgz build/*.tar ## build cinder docker image
	docker run -d -v $(BUILD_DIR):/build --entrypoint=sleep --name=$(BUILD_CONTAINER_NAME) $(BASE_IMAGE) infinity
	docker exec $(BUILD_CONTAINER_NAME) bash -c 'mkdir -p /cinder'
	docker cp ./charts $(BUILD_CONTAINER_NAME):/cinder/charts
	docker cp ./scripts $(BUILD_CONTAINER_NAME):/cinder/scripts
	docker exec \
		-e KUBERNETES_VERSION=$(KUBERNETES_VERSION) \
		-e MACHINE_API_VERSION=$(MACHINE_API_VERSION) \
		-e MACHINE_API_PROVIDER_DOCKER_VERSION=$(MACHINE_API_PROVIDER_DOCKER_VERSION) \
		$(BUILD_CONTAINER_NAME) bash -c /cinder/scripts/build.sh
	for f in $$(find files/ -type f)
	do
	    docker cp $$f $(BUILD_CONTAINER_NAME):$${f##*files/}
	done
	docker cp ./bin/crit $(BUILD_CONTAINER_NAME):/usr/bin/
	docker cp ./bin/cinder $(BUILD_CONTAINER_NAME):/usr/bin/
	docker cp ../../../templates/audit-policy.yaml $(BUILD_CONTAINER_NAME):/etc/kubernetes/audit-policy.yaml
	docker exec $(BUILD_CONTAINER_NAME) bash -c 'nohup containerd > /dev/null 2>&1 &'
	for f in $$(find build/ -name '*.tar')
	do
	    docker exec $(BUILD_CONTAINER_NAME) ctr --namespace=k8s.io images import --all-platforms --no-unpack /build/$${f##*build/}
	done
	docker exec $(BUILD_CONTAINER_NAME) pkill containerd
	docker commit --change 'ENTRYPOINT [ "/usr/local/bin/entrypoint", "/sbin/init" ]' $(BUILD_CONTAINER_NAME) $(IMAGE)
	@$(MAKE) -s -C . clean

clean:
	@echo "Cleaning up ..."
	-@docker rm -f -v $(BUILD_CONTAINER_NAME) > /dev/null 2>&1
	@rm -f $(BIN_DIR)/*

clean-all: clean
	@rm -f $(BUILD_DIR)/*
	@rm -f $(CHARTS_DIR)/*

##@ Helpers

.PHONY: help

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
