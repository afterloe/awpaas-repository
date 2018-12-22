package remote

import (
	"../../integrate/soaClient"
	"../../config"
)

var (
	addr, listUrl  string
)

func init() {
	warehouse := config.GetByTarget(config.Get("services"), "warehouse")
	addr = config.GetByTarget(warehouse, "addr").(string)
	listUrl = "/v2/_catalog"
}

func GetList() []interface{} {
	reply, err := soaClient.Call("GET", addr, listUrl, nil, nil)
	if nil != err {
		return make([]interface{}, 0)
	}
	return  reply["repositories"].([]interface{})
}