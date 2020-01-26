package handler

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	ErrorCode    int
	ErrorMessage string
}

func ErrorResponse(w http.ResponseWriter, errCode int, message string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(errCode)
	error := Error{ErrorCode: errCode, ErrorMessage: message}
	bytes, _ := json.Marshal(error)
	w.Write(bytes)
}
