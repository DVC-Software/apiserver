DB_IMAGE ?= mysql/mysql-server:latest
DB_CONTAINER ?= dvc_mysql
HIDE ?= @
DB_PORT ?= 3306
DB_HOSTPORT ?= 3306
DB_PASSWORD ?= dvcsoftware
NETWORK ?= bridge
VOLUME ?= mysqldata

.PHONY:

rm-volume:
	$(HIDE)docker volume rm $(VOLUME)

create-volume:
	$(HIDE)docker volume create --name $(VOLUME)

start-db:
	$(HIDE)docker-compose -f docker/docker-compose.yml up -d --build db 

init-db:
	$(HIDE)$(MAKE) create-volume
	$(HIDE)$(MAKE) start-db
	sleep 5
	$(HIDE)docker exec $(DB_CONTAINER) bash -c 'mysql -uroot -pdvcsoftware < /init/init.sql'

reset-db: stop-db
	$(HIDE)$(MAKE) rm-volume
	$(HIDE)$(MAKE) init-db

stop-db:
	$(HIDE)echo stopping $(DB_CONTAINER)...
	$(HIDE)docker stop $(DB_CONTAINER)
	$(HIDE)docker container rm $(DB_CONTAINER)