package test

import (
	"testing"
	"../integrate/couchdb"
)

func Test_HttpClient_demo(t *testing.T) {
	reply, err := couchdb.Create("demoi", map[string]interface{}{
		"name": "afterloe",
		"age": 6,
		"sex": "小男孩",
	})
	if nil != err {
		t.Error(reply)
	}
	t.Log(reply)
}

//func Test_HttpClient_demo(t *testing.T) {
//	reply, err := couchdb.CreateDB("demo")
//	if nil != err {
//		t.Error(reply)
//	}
//	t.Log(reply)
//}

//func Test_HttpClient_demo(t *testing.T) {
//	flag, err := couchdb.Login()
//	if nil != err {
//		t.Error(err)
//		return
//	}
//	t.Log(flag)
//	reply, err := couchdb.Read("_session", nil)
//	if nil != err {
//		t.Error(reply)
//	}
//	t.Log(reply)
//}