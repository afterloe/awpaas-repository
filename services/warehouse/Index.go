package warehouse

import (
	"../../integrate/soaClient"
	"../../config"
)

var (
	addr,port  string
)

func init() {
	db := config.GetByTarget(config.Get("services"), "db")
	addr = config.GetByTarget(db, "addr").(string)
	port = config.GetByTarget(db, "port").(string)
	addr += ":" + port
}
/**
	获取包列表
*/
func GetList(skip, limit string) []interface{} {
	params := soaClient.Encode(map[string]interface{}{
		"skip": skip,
		"limit": limit,
		"include_docs": "true",
	})
	reply, err := soaClient.Call("GET", addr, "/registry/_all_docs?" + params, nil, nil)
	var list = make([]interface{}, 0)
	if nil != err {
		return list
	}
	for _, r := range (reply["rows"].([]interface{})) {
		doc := (r.(map[string]interface{}))["doc"].(map[string]interface{})
		delete(doc, "_rev")
		list = append(list, doc)
	}
	return list
}