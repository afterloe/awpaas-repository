package test

import (
	"testing"
	"../integrate/soaClient"
	"encoding/json"
	"fmt"
)

func FormatToStruct(chunk *string) (map[string]interface{}, error){
	rep := make(map[string]interface{})
	err := json.Unmarshal([]byte(*chunk), &rep) // 使用这种方式来转换复杂json
	if nil != err {
		return nil, err
	}
	return rep, nil
}

func Test_HttpClient_demo(t *testing.T) {
	reply, err := soaClient.Call("GET",  "192.168.1.3", "/v2/_catalog", nil, nil)
	t.Log(reply)
	if nil != err {
		t.Error(err)
	}
	for k, v := range reply {
		t.Log(fmt.Sprintf("key is %s -> value is %s", k, v))
	}
}