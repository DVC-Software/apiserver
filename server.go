package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/DVC-Software/apiserver/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Golbal
var db *sql.DB
var dev_db_name string = "dvc_api_server"
var test_db_name string = "dvc_api_server_test"
var dev_port string = "0.0.0.0:8080"
var test_port string = "0.0.0.0:8070"

// struct will be removed

func indexHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		name = "default"
	}
	w.Write([]byte(fmt.Sprintf("Hello, %s\n", name)))
}

func dbReadHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM name")
	if err != nil {
		panic(err.Error())
	}
	var list []model.Name
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			panic(err.Error())
		}
		// w.Write([]byte(fmt.Sprintf("%s\n", name)))
		var n model.Name
		n.Name = name
		list = append(list, n)
	}
	resp, _ := json.Marshal(list)
	w.Write([]byte(resp))
	defer rows.Close()
}

func dbWriteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data, _ := ioutil.ReadAll(r.Body)
		var n model.Name
		json.Unmarshal(data, &n)
		name := n.Name
		if name == "" {
			panic("Empty body")
		}
		fmt.Println("EXECUTE TO THIS POINT")
		stmt, err := db.Prepare("INSERT INTO name(name) VALUES(?)")
		if err != nil {
			panic(err.Error())
		}
		stmt.Exec(name)
		var resp model.Response
		resp.Status = true
		resp.Name = n
		// Response JSON
		bytes, _ := json.Marshal(resp)
		w.Write(bytes)
		defer stmt.Close()
	}
}

func connectDB(env string) *sql.DB {
	var db_name string
	if env == "development" {
		db_name = dev_db_name
	} else if env == "test" {
		db_name = test_db_name
	} else {
		db_name = dev_db_name
	}
	db, err := sql.Open("mysql", "dvcsoftware:dvcsoftware@tcp(db:3306)/"+db_name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("successfully connected to db")
	err = db.Ping()
	if err != nil {
		fmt.Println("error ping db, bad connection")
		panic(err.Error())
	}
	return db
}

func getPortFromEnv() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "development" {
		fmt.Print("development")
		return dev_port
	} else if env == "test" {
		fmt.Print("test")
		return test_port
	} else {
		fmt.Print("development")
		return dev_port
	}
}

func main() {
	// Check env
	env := os.Getenv("ENVIRONMENT")

	// Connect to DB
	db = connectDB(env)

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/{name}", indexHandler).Methods("GET")
	r.HandleFunc("/db/post", dbWriteHandler).Methods("POST")
	r.HandleFunc("/db/get", dbReadHandler).Methods("GET")
	http.Handle("/", r)
	port := getPortFromEnv()
	srv := &http.Server{
		Handler: r,
		Addr:    port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	fmt.Println("starting server...")

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
	defer srv.Close()
	defer db.Close()
}
