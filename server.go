package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/DVC-Software/apiserver/handler"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Golbal
var dev_port string = "0.0.0.0:8080"
var test_port string = "0.0.0.0:8070"
var url_prefix string = "/"

func indexHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		name = "default"
	}
	w.Write([]byte(fmt.Sprintf("Hello, %s\n", name)))
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
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/{name}", indexHandler).Methods("GET")
	r.HandleFunc("/name/create", handler.NameCreateHandler).Methods("POST")
	r.HandleFunc("/name/update/{id}", handler.NameUpdateHandler).Methods("PUT")
	r.HandleFunc("/name/delete/{id}", handler.NameDeleteHandler).Methods("DELETE")
	r.HandleFunc("/name/show", handler.NameReadHandler).Methods("GET")
	r.HandleFunc("/member/create", handler.MemberCreateHandler).Methods("POST")
	r.HandleFunc("/member/info/{id}", handler.MemberInfoHandler).Methods("GET")
	r.HandleFunc("/training/create", handler.TrainingSessionCreateHandler).Methods("POST")
	http.Handle(url_prefix, r)
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
}
