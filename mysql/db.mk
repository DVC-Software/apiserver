DB_IMAGE ?= mysql/mysql-server:latest
DB_CONTAINER ?= dvc_mysql
HIDE ?= @
DB_PORT ?= 3306
DB_HOSTPORT ?= 3306
DB_PASSWORD ?= dvcsoftware
NETWORK ?= bridge
VOLUME ?= mysqldata

.PHONY:

pull-db:
	$(HIDE)docker pull mysql:latest

create-volume:
	$(HIDE)docker container rm $(VOLUME)
	$(HIDE)docker create -v /home/user/mysql:/var/lib/mysql --name $(VOLUME) $(DB_IMAGE)

start-db:
	$(HIDE)echo starting $(DB_CONTAINER) from $(DB_IMAGE)...
	$(HIDE)docker run -d -p $(DB_PORT):$(DB_HOSTPORT) --name $(DB_CONTAINER) -e MYSQL_ROOT_PASSWORD=$(DB_PASSWORD) -e MYSQL_ROOT_HOST=% $(DB_IMAGE)

stop-db:
	$(HIDE)echo stopping $(DB_CONTAINER)...
	$(HIDE)docker stop $(DB_CONTAINER)
	$(HIDE)docker container rm $(DB_CONTAINER)