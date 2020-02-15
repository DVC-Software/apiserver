package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/DVC-Software/apiserver/model"
)

func TrainingSessionCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" || r.Header.Get("Content-Type") != "application/json" {
		ErrorResponse(w, 500, "Invalid method or headers")
		return
	}
	data, _ := ioutil.ReadAll(r.Body)
	var info model.TrainingSessionInfo
	err := json.Unmarshal(data, &info)
	if err != nil {
		ErrorResponse(w, 500, err.Error())
		return
	}
	fmt.Println(string(data))
	if !info.ValidateInfo() {
		ErrorResponse(w, 500, "Invalid request content")
		return
	}
	db := ConnectDB()
	success, sessionCreated, errmsg := model.CreateTrainingSession(db, info)
	if !success {
		ErrorResponse(w, 500, errmsg)
		return
	}
	bytes, _ := json.Marshal(sessionCreated)
	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)
	defer db.Close()
	return
}
