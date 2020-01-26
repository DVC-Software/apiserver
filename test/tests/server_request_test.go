package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DVC-Software/apiserver/model"
	"io/ioutil"
	"net/http"
	"testing"
)

var endpointPrefix string = "http://dvc_api_server_test:8070"
var content = []string{"Hello", "dvc", "is", "awesome"}

func TestDBPost(t *testing.T) {
	for _, item := range content {
		var n model.Name
		n.Name = item
		t.Run(item, func(t *testing.T) {
			body, _ := json.Marshal(n)
			resp, err := http.Post(endpointPrefix+"/db/post", "application/json", bytes.NewBuffer(body))
			if err != nil {
				t.Error(err.Error())
			}
			fmt.Println(resp)
			if resp == nil {
				t.Error("No Response")
			}
			defer resp.Body.Close()

			data, err := ioutil.ReadAll(resp.Body)
			json.Unmarshal(data, &n)

			// Check response json
			if n.Name != item {
				t.Errorf("TestDBPost: expectd %s but got %s", item, n.Name)
			}
		})
	}
}

func TestDBGet(t *testing.T) {
	resp, err := http.Get(endpointPrefix + "/db/get")
	if err != nil {
		t.Error(err.Error())
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	// Deserialize response body
	var list []model.Name
	json.Unmarshal(data, &list)
	for index, item := range list {
		if item.Name != content[index] {
			t.Errorf("TestGetRequest: expected %s but got %s", content[index], string(item.Name))
		}
	}
}
