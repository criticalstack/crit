BIN_DIR := bin
CRIT    := $(BIN_DIR)/crit
CINDER  := $(BIN_DIR)/cinder

.PHONY: $(CRIT) $(CINDER)

$(CRIT): # Build crit
	$(MAKE) -s -C ../../../ crit BIN_DIR=$(PWD)/$(BIN_DIR)

$(CINDER): # Build cinder
	$(MAKE) -s -C ../../../ cinder BIN_DIR=$(PWD)/$(BIN_DIR)

clean:
	@rm -f bin/*

.PHONY: build

TAG            ?= v0.2.17
IMAGE          ?= cscr.io/criticalstack/cindinder:$(TAG)
CINDER_VERSION ?= v0.2.18
CINDER_IMAGE   ?= docker.io/criticalstack/cinder:$(CINDER_VERSION)
BASE_IMAGE     ?= docker:dind
BUILD_DIR      := $(PWD)/build


build: clean $(CINDER) pre-pull docker-build build-clean

docker-build:
	docker build . -t $(IMAGE)

pre-pull:
	mkdir -p $(BUILD_DIR)
	docker save -o $(BUILD_DIR)/cinder.tar $(CINDER_IMAGE)

build-clean:
	rm -f $(BUILD_DIR)/*
