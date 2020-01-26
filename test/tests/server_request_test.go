package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/DVC-Software/apiserver/model"
)

var endpointPrefix string = "http://dvc_api_server_test:8070"
var content = []string{"Hello", "dvc", "is", "awesome"}
var updated = []string{"1", "2", "3", "4"}

type Delete struct {
	Status bool   `json:"successful"`
	Id     string `json:"id"`
}

func TestDBCreate(t *testing.T) {
	for _, item := range content {
		var n model.Name
		n.Name = item
		t.Run(item, func(t *testing.T) {
			body, _ := json.Marshal(n)
			resp, err := http.Post(endpointPrefix+"/name/create", "application/json", bytes.NewBuffer(body))
			if err != nil {
				t.Error(err.Error())
			}
			if resp.StatusCode != 200 {
				t.Errorf("TestDBPost: response code is not 200, error: %d", resp.StatusCode)
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

func TestDBCreateInvalidHeader(t *testing.T) {
	n := model.Name{Name: "Invalid"}
	body, _ := json.Marshal(n)
	request, _ := http.NewRequest("POST", endpointPrefix+"/name/create", bytes.NewReader(body))
	request.Header.Set("Content-Type", "invalid")
	client := &http.Client{}
	resp, _ := client.Do(request)
	if resp.StatusCode != 500 {
		t.Errorf("TestDBUpdate: response code is not 500, error: %d", resp.StatusCode)
	}
}

func TestDBRead(t *testing.T) {
	resp, err := http.Get(endpointPrefix + "/name/show")
	if err != nil {
		t.Error(err.Error())
	}
	if resp.StatusCode != 200 {
		t.Errorf("TestDBRead: response code is not 200, error: %d", resp.StatusCode)
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
			t.Errorf("TestDBRead: expected %s but got %s", content[index], string(item.Name))
		}
	}
}

func TestDBUpdate(t *testing.T) {
	for _, item := range updated {
		t.Run(item, func(*testing.T) {
			n := model.Name{Name: item}
			body, _ := json.Marshal(n)
			request, _ := http.NewRequest("PUT", endpointPrefix+"/name/update/"+item, bytes.NewReader(body))
			request.Header.Set("Content-Type", "application/json")
			client := &http.Client{}
			resp, err := client.Do(request)
			if err != nil {
				t.Error(err.Error())
			}
			if resp.StatusCode != 200 {
				t.Errorf("TestDBUpdate: response code is not 200, error: %d", resp.StatusCode)
			}
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Error(err.Error())
			}
			json.Unmarshal(data, &n)
			if n.Name != item {
				t.Errorf("TestDBUpdate: expected %s but got %s", item, n.Name)
			}
		})
	}
}

func TestDBUpdateInvalidID(t *testing.T) {
	n := model.Name{Name: "Invalid"}
	body, _ := json.Marshal(n)
	request, _ := http.NewRequest("PUT", endpointPrefix+"/name/update/300000", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(request)
	if resp.StatusCode != 404 {
		t.Errorf("TestDBUpdate: response code is not 404, error: %d", resp.StatusCode)
	}
}

func TestDBUpdateInvalidHeader(t *testing.T) {
	n := model.Name{Name: "Invalid"}
	body, _ := json.Marshal(n)
	request, _ := http.NewRequest("PUT", endpointPrefix+"/name/update/1", bytes.NewReader(body))
	request.Header.Set("Content-Type", "invalid")
	client := &http.Client{}
	resp, _ := client.Do(request)
	if resp.StatusCode != 500 {
		t.Errorf("TestDBUpdate: response code is not 500, error: %d", resp.StatusCode)
	}
}

func TestDBDelete(t *testing.T) {
	for _, item := range updated {
		t.Run(item, func(*testing.T) {
			request, _ := http.NewRequest("DELETE", endpointPrefix+"/name/delete/"+item, nil)
			client := &http.Client{}
			resp, err := client.Do(request)
			if err != nil {
				t.Error(err.Error())
			}
			if resp.StatusCode != 200 {
				t.Errorf("TestDBDelete: response code is not 200, error: %d", resp.StatusCode)
			}
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Error(err.Error())
			}
			var result Delete
			json.Unmarshal(data, &result)
			if result.Status != true || result.Id != item {
				t.Errorf("TestDBDelete: status is %t and id is %s not %s", result.Status, result.Id, item)
			}
		})
	}
}

func TestDBDeleteInvalidID(t *testing.T) {
	request, _ := http.NewRequest("DELETE", endpointPrefix+"/name/delete/300000", nil)
	client := &http.Client{}
	resp, _ := client.Do(request)
	if resp.StatusCode != 404 {
		t.Errorf("TestDBDelete: response code is not 404, error: %d", resp.StatusCode)
	}
}
