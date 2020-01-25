# apiserver
DVC API server for database 

## How to use the current version?

make init-db 		This will initialize the database with pre-set data

make start		This will start the api server


## How to stop containers?

make stop		This will stop the api server container

make stop-db		This will stop the mysql container


## How to test?

make test   This will initialize the test container and the api_server container using test db
            After running test the test container and api_server container using test db will be removed
