package handler

import (
	"encoding/json"
	"fmt"
	"github.com/DVC-Software/apiserver/model"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func ReadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		db := ConnectDB()
		var list []model.Name
		db.Find(&list)
		resp, _ := json.Marshal(list)
		w.Write([]byte(resp))
		defer db.Close()
	} else {
		ErrorResponse(w, 500, "Invalid method or headers")
		return
	}
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" && r.Header.Get("Content-Type") == "application/json" {
		data, _ := ioutil.ReadAll(r.Body)
		db := ConnectDB()
		var n model.Name
		json.Unmarshal(data, &n)
		db.Create(&n)
		// Response JSON
		bytes, _ := json.Marshal(n)
		w.Header().Add("Content-Type", "application/json")
		w.Write(bytes)
		defer db.Close()
	} else {
		ErrorResponse(w, 500, "Invalid method or headers")
		return
	}
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" && r.Header.Get("Content-Type") == "application/json" {
		// Get id
		id := mux.Vars(r)["id"]
		db := ConnectDB()
		defer db.Close()
		var n, record model.Name
		// Check if record exists
		if err := db.Model(&record).Where("id = ?", id).Find(&record).Error; err != nil {
			ErrorResponse(w, 404, "Record does not exist")
			fmt.Println(err.Error())
			return
		}
		// Get object from request
		data, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(data, &n)
		// Update record
		db.Model(&record).UpdateColumns(n)
		// Get updated item from DB and create json response
		db.Model(&record).Where("id = ?", id).Find(&n)
		bytes, _ := json.Marshal(n)
		w.Header().Add("Content-Type", "application/json")
		w.Write(bytes)
	} else {
		ErrorResponse(w, 500, "Invalid method or headers")
		return
	}
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		// Get id
		id := mux.Vars(r)["id"]
		db := ConnectDB()
		// Attempt to delete record
		var record model.Name
		if err := db.Where("id = ?", id).Find(&record).Error; err != nil {
			ErrorResponse(w, 404, "Record does not exist")
			return
		}
		if err := db.Unscoped().Where("id = ?", id).Delete(model.Name{}).Error; err != nil {
			ErrorResponse(w, 500, err.Error())
			fmt.Println(err.Error())
			return
		}
		var result = Delete{Status: true, Id: id}
		bytes, _ := json.Marshal(result)
		fmt.Println("Successfully deleted table name, id " + id)
		w.Header().Add("Content-Type", "application/json")
		fmt.Println(string(bytes))
		w.Write(bytes)
	} else {
		ErrorResponse(w, 500, "Invalid method or headers")
	}
}
