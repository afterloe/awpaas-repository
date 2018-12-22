package test

import (
	"testing"
	"../integrate/soaClient"
)

func Test_QueryToCouchDB(t *testing.T) {
	reply, err := soaClient.Call("GET",  "192.168.3.21:5984", "/registry/_all_docs?include_docs=true&limit=10&skip=10", nil, nil)
	if nil != err {
		t.Error(err)
	}
	t.Log(reply)
	for _, v := range reply["rows"].([]interface{}) {
		doc := ((v.(map[string]interface{}))["doc"]).(map[string]interface{})
		t.Log(doc["name"])
	}
}