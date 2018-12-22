package warehouse

import (
	"../../integrate/soaClient"
	"../../config"
)

var (
	addr,port  string
)

func init() {
	warehouse := config.GetByTarget(config.Get("services"), "db")
	addr = config.GetByTarget(warehouse, "addr").(string)
	port = config.GetByTarget(warehouse, "port").(string)
	addr += ":" + port
}

func GetList(skip, limit string) []interface{} {
	params := soaClient.Encode(map[string]interface{}{
		"skip": skip,
		"limit": limit,
		"include_docs": "true",
	})
	reply, err := soaClient.Call("GET", addr, "/registry/_all_docs?" + params, nil, nil)
	if nil != err {
		return make([]interface{}, 0)
	}
	return  reply["rows"].([]interface{})
}