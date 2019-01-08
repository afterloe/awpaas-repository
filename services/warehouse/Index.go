package warehouse

import (
	"../../integrate/couchdb"
	"../../exceptions"
	"../../integrate/soaClient"
	"time"
)

var fsServiceName string

func init() {
	fsServiceName = ""
}

/**
	补充文件信息
 */
func supplementFileStatus(w *warehouse) (*warehouse, error) {
	reply, _ := soaClient.Call("GET", fsServiceName, "/v1/file/" + w.Fid, nil, nil)
	// TODO
	if 200 != reply["code"].(int) {
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
		Fields("_id", "name", "uploadTime", "group", "").
		Page(begin, limit))
	return reply
}