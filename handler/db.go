package handler

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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
	db, err := gorm.Open("mysql", "dvcsoftware:dvcsoftware@tcp(db:3306)/"+dbName+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("successfully connected to db")
	return db
}
