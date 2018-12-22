package test

import (
	"testing"
	"../integrate/soaClient"
)

func Test_QueryToCouchDB(t *testing.T) {
	reply, err := soaClient.Call("GET",  "192.168.3.21:5984", "/registry/_all_docs?limit=5&skip=10", nil, nil)
	if nil != err {
		t.Error(err)
	}
	t.Log(reply)
	for _, v := range reply["rows"].([]interface{}) {
		id := (v.(map[string]interface{}))["id"]
		//t.Log(id.(string))
		r, _ := soaClient.Call("GET",  "192.168.3.21:5984", "/registry/" + id.(string), nil, nil)
		t.Log(r["name"])
	}
}