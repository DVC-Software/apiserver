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
var dev_db_name string = os.Getenv("MYSQL_DATABASE")
var test_db_name string = dev_db_name + "_test"
var env string = os.Getenv("ENVIRONMENT")
var user string = os.Getenv("MYSQL_USER")
var password string = os.Getenv("MYSQL_PASSWORD")

func ConnectDB() *gorm.DB {
	var dbName string
	if env == "test" {
		dbName = test_db_name
	} else if env == "init" {
		user = "root"
	} else {
		dbName = dev_db_name
	}

	db, err := gorm.Open("mysql", user+":"+password+"@tcp(db:3306)/"+dbName+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("successfully connected to db")
	return db
}

func migrateTestDB(db *gorm.DB) {
	// Drop all tables and re-initialize
	db.DropTableIfExists(&model.Name{}, &model.Member{}, &model.Profile{}, &model.Position{})
	db.AutoMigrate(&model.Name{}, &model.Member{}, &model.Profile{}, &model.Position{})
}

func migrateDB(db *gorm.DB) {
	// Migrate all tables
	db.Exec("Use " + dev_db_name)
	db.AutoMigrate(&model.Name{}, &model.Member{}, &model.Profile{}, &model.Position{})
}

func initDB(db *gorm.DB) {
	db.Exec("GRANT ALL PRIVILEGES ON *.* TO 'dvcsoftware'@'%' WITH GRANT OPTION;")
	db.Exec("DROP DATABASE IF EXISTS " + dev_db_name + ";")
	db.Exec("CREATE DATABASE " + dev_db_name + ";")
	db.Exec("DROP DATABASE IF EXISTS " + test_db_name + ";")
	db.Exec("CREATE DATABASE " + test_db_name + ";")
}

func main() {
	db := ConnectDB()
	fmt.Println("migrating...")
	if env == "init" {
		initDB(db)
	}

	if env == "test" {
		migrateTestDB(db)
	} else {
		migrateDB(db)
	}
	fmt.Printf("successfully migrated for %s \n", env)
	defer db.Close()
}
