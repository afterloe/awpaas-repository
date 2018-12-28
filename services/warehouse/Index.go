package warehouse

import (
	"../../integrate/couchdb"
	"../../config"
	"../../exceptions"
	"reflect"
)

var (
	addr,port  string
	dbName string
)

func init() {
	db := config.GetByTarget(config.Get("services"), "db")
	addr = config.GetByTarget(db, "addr").(string)
	port = config.GetByTarget(db, "port").(string)
	addr += ":" + port
	dbName = "registry"
}

type Module struct {
	Name, Group, Remarks, Version, Fid string
	PackInfo map[string]interface{}
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
	补充文件信息
 */
func supplementFileStatus(module *Module) (*Module, error) {
	reply, _ := couchdb.Read("file-system/" + module.Fid, nil)
	if "not_found" == reply["error"]{
		return nil, &exceptions.Error{"pack info not found", 404}
	}
	delete(reply, "_id")
	delete(reply, "_rev")
	module.PackInfo = reply
	return module, nil
}

/**
	发送至远程
*/
func (this *Module) AppendToRemote() (map[string]interface{}, error) {
	this, err := supplementFileStatus(this)
	if nil != err {
		return nil, err
	}
	return couchdb.Create(dbName, this)
}

/**
	获取包列表
*/
func GetList(skip, limit string) []interface{} {
	reply, _ :=couchdb.Read(dbName + "/_all_docs", map[string]interface{}{
		"skip": skip,
		"limit": limit,
		"include_docs": "true",
	})
	var list = make([]interface{}, 0)
	if "not_found" == reply["error"]{
		return list
	}
	for _, r := range (reply["rows"].([]interface{})) {
		doc := (r.(map[string]interface{}))["doc"].(map[string]interface{})
		delete(doc, "_rev")
		delete(doc, "PackInfo")
		list = append(list, doc)
	}
	return list
}