DOCKER_IMAGE ?= dvc-software/dvc_api_server
DOCKER_CONTAINER ?= dvc_api_server
HIDE ?= @
PORT ?= 8080
HOSTPORT ?= 8080
NETWORK ?= bridge

.PHONY: 

build:
	$(HIDE)docker build -f Dockerfile -t $(DOCKER_IMAGE) $(PWD)

start:
	$(HIDE)docker run -d -p $(PORT):$(HOSTPORT) --network=$(NETWORK) --name $(DOCKER_CONTAINER) $(DOCKER_IMAGE)

stop: 
	$(HIDE)docker stop $(DOCKER_CONTAINER)
	$(HIDE)docker container rm $(DOCKER_CONTAINER)