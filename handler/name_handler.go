package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/DVC-Software/apiserver/model"
)

func ReadHandler(w http.ResponseWriter, r *http.Request) {
	db := ConnectDB()
	var list []model.Name
	db.Find(&list)
	resp, _ := json.Marshal(list)
	w.Write([]byte(resp))
	defer db.Close()
}

func WriteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data, _ := ioutil.ReadAll(r.Body)
		db := ConnectDB()
		var n model.Name
		json.Unmarshal(data, &n)
		db.Create(&n)
		fmt.Println("EXECUTE TO THIS POINT")
		// Response JSON
		bytes, _ := json.Marshal(n)
		w.Write(bytes)
		defer db.Close()
	}
}
