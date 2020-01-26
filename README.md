# apiserver
DVC API server for database 

## How to use the current version?

make init-db 		This will initialize the database with pre-set data

	- If an access denied error occur when make init-db, please use docker exec -it dvc_mysql to run init/init.sql 

make start		This will start the api server


## How to stop containers?

make stop		This will stop the api server container

make stop-db		This will stop the mysql container


## How to test?

make test   This will initialize the test container and the api_server container using test db
            After running test the test container and api_server container using test db will be removed
