DB_IMAGE ?= mysql/mysql-server:latest
DB_CONTAINER ?= dvc_mysql
HIDE ?= @
DB_PORT ?= 3306
DB_HOSTPORT ?= 3306
DB_PASSWORD ?= dvcsoftware
NETWORK ?= bridge
VOLUME ?= mysqldata

.PHONY:

create-volume:
	$(HIDE)docker volume rm $(VOLUME)
	$(HIDE)docker volume create --name $(VOLUME)

start-db:
	$(HIDE)docker-compose -f docker/docker-compose.yml up -d --build db 

init-db: start-db
	docker exec $(DB_CONTAINER) bash -c 'mysql -uroot -pdvcsoftware < /init/init.sql'

stop-db:
	$(HIDE)echo stopping $(DB_CONTAINER)...
	$(HIDE)docker stop $(DB_CONTAINER)
	$(HIDE)docker container rm $(DB_CONTAINER)