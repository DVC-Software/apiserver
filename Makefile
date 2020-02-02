DOCKER_IMAGE ?= dvc_api_server
DOCKER_CONTAINER ?= dvc_api_server
HIDE ?= @
PORT ?= 8080
HOSTPORT ?= 8080
NETWORK ?= bridge
SERVICES = $(DOCKER_IMAGE) $(DOCKER_IMAGE)_test db

-include mysql/db.mk
.PHONY: test

build:
	$(HIDE)docker-compose -f docker/docker-compose.yml build $(SERVICES)

start:
	$(HIDE)docker-compose -f docker/docker-compose.yml up --build $(DOCKER_CONTAINER)

daemon:
	$(HIDE)docker-compose -f docker/docker-compose.yml up -d --build $(DOCKER_CONTAINER)

stop: 
	$(HIDE)docker stop $(DOCKER_CONTAINER)
	$(HIDE)docker container rm $(DOCKER_CONTAINER)

test:
	$(HIDE)docker-compose -f docker/docker-compose.yml run --rm -e ENVIRONMENT=test migrations
	$(HIDE)docker-compose -f docker/docker-compose.yml build $(DOCKER_CONTAINER)_test
	$(HIDE)docker-compose -f docker/docker-compose.yml up test
	$(HIDE)docker-compose -f docker/docker-compose.yml rm -f -s $(DOCKER_CONTAINER)_test test

rm:
	$(HIDE)docker rm $(docker ps -a -q)

rm-all:
	$(HIDE)docker system prune -a -f