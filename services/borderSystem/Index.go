package borderSystem

import (
	"database/sql"
	"../dbConnect"
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
	insertSQL = "INSERT INTO file(name, savePath, contentType, key, uploadTime, size, status, modifyTime, rev) values(?, ?, ?, ? ,? ,? ,? ,? ,?)"
	updateSQL = "UPDATE file SET modifyTime = ?, status = ? WHERE id = ?"
	deleteSQL = "DELETE FROM file WHERE id = ?"
)

func init() {
	cfg := config.Get("custom")
	root = config.GetByTarget(cfg, "root").(string)
}

func (this *fsFile) GeneratorSavePath() string {
	return fmt.Sprintf("%s/%s", this.SavePath, this.Key)
}

func (this *fsFile) SaveToDB(rev ...bool) (map[string]interface{}, error){
	if 0 != len(rev) {
		if 0 == this.Id {
			return nil, &exceptions.Error{Msg: "id can't be empty", Code: 400}
		}
		return dbConnect.WithTransaction(func(tx *sql.Tx) (map[string]interface{}, error) {
			stmt, err := tx.Prepare(updateSQL)
			if nil != err {
				return nil, &exceptions.Error{Msg: "db stmt open failed.", Code: 500}
			}
			stmt.Exec(this.ModifyTime, this.Status, this.Id)
			logger.Logger("borderSystem", "update success")
			return map[string]interface{}{}, nil
		})
	}
	return dbConnect.WithTransaction(func (tx *sql.Tx)(map[string]interface{}, error) {
		stmt, err := tx.Prepare(insertSQL)
		if nil != err {
			return nil, &exceptions.Error{Msg: "db stmt open failed.", Code: 500}
		}
		stmt.Exec(this.Name, this.SavePath, this.ContentType, this.Key, this.UploadTime, this.Size, this.Status, this.ModifyTime, this.rev)
		logger.Logger("borderSystem", "insert success")
		return map[string]interface{}{}, nil
	})
}

func (this *fsFile) Del(f ...bool) error {
	if 0 != len(f) { // 强制删除
		logger.Logger("borderSystem", "强制删除文件")
		dbConnect.WithPrepare(deleteSQL, func(stmt *sql.Stmt) (map[string]interface{}, error) {
			result, _ := stmt.Exec(this.Id)
			affectNum, _ := result.RowsAffected()
			if 0 != affectNum {
				logger.Logger("borderSystem", "delete success")
			}
			return nil, nil
		})
		err := os.Remove(this.GeneratorSavePath())
		if nil != err {
			return &exceptions.Error{Msg: "file has been deleted.", Code: 400}
		}
		return nil
	} else { // 逻辑删除
		this.Status = false
		this.ModifyTime = time.Now().Unix()
		_, err := this.SaveToDB(true) // 更新
		logger.Logger("borderSystem", "逻辑删除文件")
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

func GetAll(begin, limit int) []map[string]interface{} {
	reply, err := dbConnect.Select("file").Fields("id", "name", "uploadTime", "size").Page(begin, limit).Query()
	if nil != err {
		return nil
	}
	return reply
}

func GetList(begin, limit int) []map[string]interface{} {
	reply, err := dbConnect.Select("file").
		Fields("id", "name", "uploadTime", "size").
		AND("status = ?").Page(begin, limit).Query(true)
	if nil != err {
		return nil
	}
	return reply
}

func GetOne(key int, fields ...string) (*fsFile, error) {
	str := dbConnect.Select("file")
	if 0 == len(fields) {
		str.Fields("id,name, savePath, contentType, key, uploadTime, size, status, modifyTime")
	} else {
		str.Fields(fields...)
	}
	str.AND("id = ?", "status = ?")
	one, err := dbConnect.WithQuery(str.Preview(), func(rows *sql.Rows) (interface{}, error) {
		target := new(fsFile)
		for rows.Next() {
			rows.Scan(&target.Id, &target.Name, &target.SavePath,&target.ContentType,&target.Key,&target.UploadTime,&target.Size,&target.Status,&target.ModifyTime)
		}
		return target, nil
	}, key, true)
	f := one.(fsFile)
	return &f, err
}