package borderSystem

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"../../exceptions"
	"../../util"
	"../../config"
	"fmt"
	"time"
	"../../integrate/logger"
	"os"
)

var (
	root string
	dbPath string
)

func init() {
	cfg := config.Get("custom")
	root = config.GetByTarget(cfg, "root").(string)
}

func (this *fsFile) SaveToDB(rev ...bool) (map[string]interface{}, error){
	if 0 != len(rev) {
		jsonStr, _ := util.FormatToString(*this)
		m, _ := util.FormatToMap(jsonStr)
		m["_rev"] = this.rev
		// TODO
		return nil, nil
	}
	db, err := sql.Open("sqlite3", dbPath)
	defer db.Close()
	if nil != err {
		return nil, &exceptions.Error{Msg: "db open failed.", Code: 500}
	}
	tx, err := db.Begin()
	if nil != err {
		return nil, &exceptions.Error{Msg: "db transaction open failed.", Code: 500}
	}
	stmt, err := tx.Prepare("insert into file(name, savePath, contentType, key, uploadTime, size, status, modifyTime, rev) values(?, ?, ?, ? ,? ,? ,? ,? ,?)")
	stmt.Exec(this.Name, this.SavePath, this.ContentType, this.Key, this.UploadTime, this.Size, this.Status, this.ModifyTime, this.rev)
	if nil != err {
		return nil, &exceptions.Error{Msg: "db stmt open failed.", Code: 500}
	}
	defer stmt.Close()
	tx.Commit()
	return map[string]interface{}{}, nil
}

func (this *fsFile) Del(f ...bool) error {
	if 0 != len(f) { // 强制删除
		logger.Logger("borderSystem", "强制删除")
		couchdb.Delete(couchdb.GeneratorDelObj(this.Id, this.rev))
		err := os.Remove(this.GeneratorSavePath())
		if nil != err {
			return &exceptions.Error{Msg: "file has been deleted.", Code: 400}
		}
		return nil
	} else { // 逻辑删除
		this.Status = false
		this.ModifyTime = time.Now().Unix()
		_, err := this.SaveToDB(true) // 强制更新
		logger.Logger("borderSystem", "逻辑删除")
		return err
	}
}

func Del(id string, f ...bool) error {
	file, err := GetOne(id)
	if nil != err {
		return err
	}
	return file.Del(f...)
}

func Default(name, contentType string, size int64) *fsFile {
	return &fsFile{
		SavePath: root,
		Key: util.GeneratorUUID(),
		UploadTime: time.Now().Unix(),
		Name: name,
		ContentType: contentType,
		Size: size,
		Status: true,
	}
}

func (this *fsFile) GeneratorSavePath() string {
	return fmt.Sprintf("%s/%s", this.SavePath, this.Key)
}

func GetAll(begin, limit int) []interface{} {
	reply, _ := couchdb.Find(couchdb.Condition().Fields("_id", "name", "uploadTime", "size").
		Page(begin, limit))
	return reply
}

func GetList(begin, limit int) []interface{} {
	condition := couchdb.Condition().Append("status", "$eq", true).
		Fields("name", "uploadTime", "_id").
		Page(begin, limit)
	reply, _ := couchdb.Find(condition)
	return reply
}

func GetOne(key string, fields ...string) (*fsFile, error) {
	condition := couchdb.Condition().Append("_id", "$eq", key).
		Append("status", "$eq", true)
	if 0 != len(fields) {
		condition = condition.Fields(fields...)
	}
	reply, _ := couchdb.Find(condition)
	if 0 != len(reply) {
		var target fsFile
		item := reply[0].(map[string]interface{})
		couchdb.Decode(item, &target)
		target.Id = item["_id"].(string)
		target.rev = item["_rev"].(string)
		return &target, nil
	} else {
		return nil, &exceptions.Error{Msg: "no such this file", Code: 404}
	}
}