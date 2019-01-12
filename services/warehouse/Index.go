package warehouse

import (
	"../../integrate/couchdb"
	"../../exceptions"
	"../../integrate/soaClient"
	"../../config"
	"time"
)

var fsServiceName string

func init() {
	fsServiceName = config.GetByTarget(config.Get("custom"), "fsServiceName").(string)
}

/**
	补充文件信息
 */
func supplementFileStatus(w *warehouse) (*warehouse, error) {
	reply, _ := soaClient.Call("GET", fsServiceName, "/v1/file/" + w.Fid, nil, nil)
	if 200 != reply["code"].(float64) {
		return nil, &exceptions.Error{Msg: "no such this file", Code: 400}
	}
	w.PackInfo = reply["data"].(map[string]interface{})
	return w, nil
}

/**
	发送至远程
*/
func (this *warehouse) SaveToDB() (map[string]interface{}, error) {
	this, err := supplementFileStatus(this)
	if nil != err {
		return nil, err
	}
	return couchdb.Create(this)
}

func Default() *warehouse {
	return &warehouse{
		Status: true,
		UploadTime: time.Now().Unix(),
	}
}

/**
	获取包列表
*/
func GetList(begin, limit int) []interface{} {
	reply, _ := couchdb.Find(couchdb.Condition().Append("status", "$eq", true).
		Fields("_id", "name", "uploadTime", "group").
		Page(begin, limit))
	return reply
}

/**
	更行包信息
 */
func Update() (interface{}, error) {
	// TODO
	return map[string]interface{}{
		"flag": true,
	}, nil
}

/**
	查询包详细信息
*/
func GetOne(key string, fields ...string) (*warehouse, error) {
	condition := couchdb.Condition().Append("_id", "$eq", key).
		Append("status", "$eq", true)
	if 0 != len(fields) {
		condition = condition.Fields(fields...)
	}
	reply, _ := couchdb.Find(condition)
	if 0 != len(reply) {
		var target warehouse
		item := reply[0].(map[string]interface{})
		couchdb.Decode(item, &target)
		target.Id = item["_id"].(string)
		target.rev = item["_rev"].(string)
		return &target, nil
	} else {
		return nil, &exceptions.Error{Msg: "no such this package", Code: 404}
	}
}