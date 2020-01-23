package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

// Golbal
var db *sql.DB

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
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			panic(err.Error())
		}
		w.Write([]byte(fmt.Sprintf("%s\n", name)))
	}
	defer rows.Close()
}

func dbWriteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		vars := mux.Vars(r)
		name := vars["name"]
		fmt.Println("EXECUTE TO THIS POINT")
		stmt, err := db.Prepare("INSERT INTO name(name) VALUES(?)")
		if err != nil {
			panic(err.Error())
		}
		stmt.Exec(name)
		w.Write([]byte(fmt.Sprintf("successfully insert, %s!\n", name)))
		defer stmt.Close()
	}
}

func connectDB() *sql.DB {
	db, err := sql.Open("mysql", "dvcsoftware:dvcsoftware@tcp(db:3306)/dvc_api_server")
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

func main() {
	// Connect to DB
	db = connectDB()

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/{name}", indexHandler).Methods("GET")
	r.HandleFunc("/db/post/{name}", dbWriteHandler).Methods("POST")
	r.HandleFunc("/db/get", dbReadHandler).Methods("GET")
	http.Handle("/", r)
	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8080",
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
