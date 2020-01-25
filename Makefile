DOCKER_IMAGE ?= dvc-software/dvc_api_server
DOCKER_CONTAINER ?= dvc_api_server
HIDE ?= @
PORT ?= 8080
HOSTPORT ?= 8080
NETWORK ?= bridge

-include mysql/db.mk

.PHONY: test

build:
	$(HIDE)docker build -f Dockerfile -t $(DOCKER_IMAGE) $(PWD)

start:
	$(HIDE)docker-compose -f docker/docker-compose.yml up --build $(DOCKER_CONTAINER)

daemon:
	$(HIDE)docker-compose -f docker/docker-compose.yml up -d --build $(DOCKER_CONTAINER)

stop: 
	$(HIDE)docker stop $(DOCKER_CONTAINER)
	$(HIDE)docker container rm $(DOCKER_CONTAINER)

test:
	$(HIDE)docker exec $(DB_CONTAINER) bash -c 'mysql -uroot -pdvcsoftware < /init/init_test.sql'
	$(HIDE)docker-compose -f docker/docker-compose.yml up test
	$(HIDE)docker-compose -f docker/docker-compose.yml rm -f -s $(DOCKER_CONTAINER)_test test

rm:
	$(HIDE)docker rm $(docker ps -a -q)

rm-all:
	$(HIDE)docker system prune -f