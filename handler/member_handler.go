package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/DVC-Software/apiserver/model"
	"github.com/gorilla/mux"
)

func MemberCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" && r.Header.Get("Content-Type") == "application/json" {
		data, _ := ioutil.ReadAll(r.Body)
		var initMember model.MemberInit
		err := json.Unmarshal(data, &initMember)
		if err != nil {
			ErrorResponse(w, 500, err.Error())
			return
		}
		fmt.Println(string(data))
		if initMember.Name == "" || initMember.CreatedFrom == "" || (initMember.SlackUserID == "" && initMember.DiscordUserID == "") {
			fmt.Println(initMember.Name)
			fmt.Println(initMember.CreatedFrom)
			fmt.Println(initMember.DiscordUserID)
			ErrorResponse(w, 500, "Invaid request content")
			return
		}
		// need to be move to model package
		db := ConnectDB()
		success, memberInfo, errmsg := model.CreateMember(db, initMember)
		if !success {
			ErrorResponse(w, 500, errmsg)
			return
		}
		bytes, _ := json.Marshal(memberInfo)
		w.Header().Add("Content-Type", "application/json")
		w.Write(bytes)
		defer db.Close()
		return
	} else {
		ErrorResponse(w, 500, "Invalid method or headers")
		return
	}
}

func MemberInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		ErrorResponse(w, 500, "Invalid method or headers")
		return
	}
	id := mux.Vars(r)["id"]
	// need to be move to model package
	db := ConnectDB()
	success, memberInfo, errmsg := model.IdentifyMember(db, id)
	if !success {
		ErrorResponse(w, 500, errmsg)
		return
	}
	bytes, _ := json.Marshal(memberInfo)
	fmt.Println(string(bytes))
	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)
	defer db.Close()
}
