-- Initialize user (THIS LINE IS CORRECT!)
	
GRANT ALL PRIVILEGES ON *.* TO 'dvcsoftware'@'%' WITH GRANT OPTION;


-- Initialize development db
DROP DATABASE IF EXISTS dvc_api_server;

CREATE DATABASE dvc_api_server;

Use dvc_api_server;

DROP TABLE IF EXISTS name;

CREATE TABLE name (name VARCHAR(20));

-- Initialize test db

DROP DATABASE IF EXISTS dvc_api_server_test;

CREATE DATABASE dvc_api_server_test;

Use dvc_api_server_test;

DROP TABLE IF EXISTS name;

CREATE TABLE name (name VARCHAR(20));