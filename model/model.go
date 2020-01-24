package model

type Name struct {
	Name string `json:"name"`
}

type Response struct {
	Status bool `json:"status"`
	Name   Name `json:"name"`
}
