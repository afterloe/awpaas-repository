package test

import (
	"testing"
	"../integrate/soaClient"
)

var (
	host, dbName string
	username, password string
)

func init() {
	host = "mine:5984"
	dbName = "file-system"
	username = "ascs"
	password = "ascs.tech"
}

func Test_QueryToCouchDB(t *testing.T) {
	reply, err := soaClient.Call("GET",  host, "/hyoertable/_all_docs?include_docs=true&limit=10&skip=0", nil, nil)
	if nil != err {
		t.Error(err)
	}
	t.Log(reply)
	for _, v := range reply["rows"].([]interface{}) {
		doc := ((v.(map[string]interface{}))["doc"]).(map[string]interface{})
		t.Log(doc)
	}
}