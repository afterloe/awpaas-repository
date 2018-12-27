package test

import (
	"testing"
	"../integrate/couchdb"
)

func Test_HttpClient_demo(t *testing.T) {
	flag, err := couchdb.Login()
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(flag)
	couchdb.UserInfo()
}