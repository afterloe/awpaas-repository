package warehouse

import (
	"../../integrate/couchdb"
	"../../exceptions"
	"../../integrate/soaClient"
	"../../config"
	"../../util"
	"time"
	"fmt"
	"os/exec"
	"bytes"
)

var fsServiceName string

func init() {
	fsServiceName = config.GetByTarget(config.Get("custom"), "fsServiceName").(string)
	registryType = [3]string{"code", "image", "tar"}
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
	保存至远程
*/
func (this *warehouse) SaveToDB() (map[string]interface{}, error) {
	this, err := supplementFileStatus(this)
	if nil != err {
		return nil, err
	}
	return couchdb.Create(this)
}

func (this *warehouse) Modify() (map[string]interface{}, error) {
	this, err := supplementFileStatus(this)
	if nil != err {
		return nil, err
	}
	this.ModifyTime = time.Now().Unix()
	jsonStr, _ := util.FormatToString(*this)
	m, _ := util.FormatToMap(jsonStr)
	m["_rev"] = this.rev
	return couchdb.Update(this.Id, m)
}

func GetRegistryType() interface{} {
	return registryType
}

func DefaultCmd(inputType string, content ...string) (*cmd, error) {
	for _, t := range registryType {
		if t == inputType {
			return &cmd{inputType, content}, nil
		}
	}
	return nil, &exceptions.Error{Msg: "no such this type", Code: 400}
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

func AppendCI(w *warehouse, ci *cmd) (interface{}, error) {
	if nil == ci {
		return nil, &exceptions.Error{Msg: "cmd not found", Code: 400}
	}
	w.Cmd = *ci
	return w.Modify()
}


/**
	更行包信息
 */
func Update(args, old *warehouse) (interface{}, error) {
	flag := false
	if "" != args.Name {
		old.Name = args.Name
		flag = true
	}
	if "" != args.Group {
		old.Group = args.Group
		flag = true
	}
	if "" != args.Remarks {
		old.Remarks = args.Remarks
		flag = true
	}
	if "" != args.Version {
		old.Version = args.Version
		flag = true
	}
	if "" != args.Fid && old.Fid != args.Fid {
		old.Fid = args.Fid
		flag = true
	}
	if !flag {
		return nil, &exceptions.Error{Msg: "no change", Code: 400}
	}
	return old.Modify()
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
		if nil != item["_id"] {
			target.Id = item["_id"].(string)
		}
		if nil != item["_rev"] {
			target.rev = item["_rev"].(string)
		}
		return &target, nil
	} else {
		return nil, &exceptions.Error{Msg: "no such this package", Code: 404}
	}
}

/**
	软件构建

	1.查询源文件
	2.判断ci类型
	3.按照类型进行分发处理
 */
func Build(w *warehouse) (interface{}, error) {
	cmd := w.Cmd
	switch cmd.RegistryType {
		case "code":
			context := "/tmp/download/" + util.GeneratorUUID()
			_, err := soaClient.DownloadFile(fmt.Sprintf("http://%s/v1/download/%s", fsServiceName, w.Fid),
				context)
			if nil != err {
				return nil, err
			}
			reply := make([]string, 0)
			for _, c := range cmd.Content {
				var out bytes.Buffer
				e := exec.Command(context, c)
				e.Stdout = &out
				err := e.Run()
				if nil != err {
					return nil, err
				}
				reply = append(reply, out.String())
			}
			return reply, nil
		case "image":
			return nil, nil
		case "tar":
			return nil, nil
		default:
			return nil, &exceptions.Error{Msg: "can't supper this"}
	}
}