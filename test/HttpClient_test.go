package test

import (
	"testing"
	"../integrate/soaClient"
	"fmt"
)

func Test_HttpClient_demo(t *testing.T) {
	reply, err := soaClient.Call("GET",  "192.168.3.21:180", "/v2/_catalog", nil, nil)
	if nil != err {
		t.Error(err)
	}
	t.Log(reply)
	for k, v := range  reply {
		t.Log(fmt.Sprintf("key is %s -> value is %s", k, v))
	}
}