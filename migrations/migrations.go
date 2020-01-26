package main

import (
	"fmt"
	"github.com/DVC-Software/apiserver/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

// Golbal
var dev_db_name string = "dvc_api_server"
var test_db_name string = "dvc_api_server_test"
var env string = os.Getenv("ENVIRONMENT")

func ConnectDB() *gorm.DB {
	var dbName string
	if env == "test" {
		dbName = test_db_name
	} else {
		dbName = dev_db_name
	}
	db, err := gorm.Open("mysql", "root:dvcsoftware@tcp(db:3306)/"+dbName+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("successfully connected to db")
	return db
}

func migrateTestDB(db *gorm.DB) {
	// Drop all tables and re-initialize
	db.DropTableIfExists(&model.Name{})
	db.AutoMigrate(&model.Name{})
}

func migrateDB(db *gorm.DB) {
	// Migrate all tables
	db.Exec("Use " + dev_db_name)
	db.AutoMigrate(&model.Name{})
}

func initDB(db *gorm.DB) {
	db.Exec("GRANT ALL PRIVILEGES ON *.* TO 'dvcsoftware'@'%' WITH GRANT OPTION;")
	db.Exec("DROP DATABASE IF EXISTS dvc_api_server;")
	db.Exec("CREATE DATABASE dvc_api_server;")
	db.Exec("DROP DATABASE IF EXISTS dvc_api_server_test;")
	db.Exec("CREATE DATABASE dvc_api_server_test;")
}

func main() {
	db := ConnectDB()
	if env == "init" {
		initDB(db)
	}

	if env == "test" {
		migrateTestDB(db)
	} else {
		migrateDB(db)
	}
	defer db.Close()
}
