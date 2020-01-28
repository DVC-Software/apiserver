package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/DVC-Software/apiserver/model"
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
		if initMember.Name == "" || initMember.CreatedFrom == "" || (initMember.SlackUserID == "" && initMember.DiscordUserID == "") {
			ErrorResponse(w, 500, "Invaid request content")
			return
		}
		db := ConnectDB()
		success, member := model.CreateMember(db, initMember)
		if !success {
			ErrorResponse(w, 500, "Error creating member")
			return
		}
		bytes, _ := json.Marshal(member)
		w.Header().Add("Content-Type", "application/json")
		w.Write(bytes)
		defer db.Close()
	} else {
		ErrorResponse(w, 500, "Invalid method or headers")
		return
	}
}
