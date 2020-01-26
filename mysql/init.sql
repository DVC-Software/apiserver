-- Only use this if there's an access denied error in make init-db
CREATE USER 'root'@'%' IDENTIFIED BY 'dvcsoftware';
GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' WITH GRANT OPTION;