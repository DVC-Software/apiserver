package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		name = "default"
	}
	w.Write([]byte(fmt.Sprintf("Hello, %s\n", name)))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/{name}", indexHandler).Methods("GET")
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
}
