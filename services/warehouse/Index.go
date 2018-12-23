package warehouse

import (
	"../../integrate/soaClient"
	"../../config"
	"../../exceptions"
	"fmt"
	"reflect"
)

var (
	addr,port  string
	db string
)

func init() {
	db := config.GetByTarget(config.Get("services"), "db")
	addr = config.GetByTarget(db, "addr").(string)
	port = config.GetByTarget(db, "port").(string)
	addr += ":" + port
	db = "registry"
}

type Module struct {
	Name, Group, Remarks, Version, Fid string
	_id string
}

/**
	参数检测
*/
func (this *Module) Check(args ...string) error {
	value := reflect.ValueOf(*this)
	for _, arg := range args {
		v := value.FieldByName(arg)
		if !v.IsValid() {
			break
		}
		if "" == v.Interface() {
			return &exceptions.Error{Msg: "lack param " + arg, Code: 400}
		}
	}

	return nil
}

/**
	发送至远程
*/
func (this *Module) AppendToRemote() (map[string]interface{}, error) {
	reply, _ := soaClient.Call("GET", addr, "/_uuids?count=1", nil, nil)
	id := reply["uuids"].([]interface{})[0]
	this._id = id.(string)
	save:
	reply, _ = soaClient.Call("PUT", addr, fmt.Sprintf("/%s/%v", addr, id),
		soaClient.GeneratorBody(this), soaClient.GeneratorPostHeader())
	// 如果数据库不存在，则创建
	if "not_found" == reply["error"]{
		soaClient.Call("PUT", addr, db, nil, nil)
		goto save
	}
	return reply, nil
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
	reply, err := soaClient.Call("GET", addr, fmt.Sprintf("/%s/%s?%s", db, "_all_docs", params), nil, nil)
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