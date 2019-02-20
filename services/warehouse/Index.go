package warehouse

import (
	"../dbConnect"
	"../../integrate/logger"
	"../../exceptions"
	"../../config"
	"time"
	"database/sql"
)

var (
	fsServiceName, root string
	insertSQL = "INSERT INTO warehouse(fid, name, \"group\", remarks, version, uploadTime, modifyTime, status) VALUES(?,?,?,?,?,?,?,?)"
	updateSQL = "UPDATE warehouse SET name = ?, \"group\" = ?, remarks = ?, version = ?, modifyTime = ?, status = ? WHERE id = ?"
)

func init() {
	cfg := config.Get("custom")
	fsServiceName = config.GetByTarget(cfg, "fsServiceName").(string)
	root = config.GetByTarget(cfg, "root").(string)
}

/**
	保存至远程
*/
func (this *warehouse) SaveToDB() (interface{}, error) {
	return dbConnect.WithTransaction(func(tx *sql.Tx) (interface{}, error) {
		stmt, err := tx.Prepare(insertSQL)
		if nil != err {
			return nil, &exceptions.Error{Msg: "db stmt open failed.", Code: 500}
		}
		result, _ := stmt.Exec(this.FId, this.Name, this.Group, this.Remarks, this.Version, this.UploadTime, this.ModifyTime, this.Status)
		id, _ := result.LastInsertId()
		this.Id = id
		logger.Logger("warehouse", "insert success")
		return this, nil
	})
}

func (this *warehouse) Modify() (interface{}, error) {
	this.ModifyTime = time.Now().Unix()
	if 0 == this.Id {
		return nil, &exceptions.Error{Msg: "no such this id", Code: 400}
	}
	return dbConnect.WithTransaction(func(tx *sql.Tx) (interface{}, error) {
		stmt, err := tx.Prepare(updateSQL)
		if nil != err {
			return nil, &exceptions.Error{Msg: "db stmt open failed.", Code: 500}
		}
		stmt.Exec(this.Name, this.Group, this.Remarks, this.Version, this.ModifyTime, this.Status, this.Id)
		logger.Logger("warehouse", "update success")
		return map[string]interface{}{}, nil
	})
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
func GetList(begin, limit int) []map[string]interface{} {
	reply, err := dbConnect.Select("warehouse").
		Fields("id", "name", "uploadTime", "\"group\"").
		AND("status = ?").Page(begin, limit).Query(true)
	if nil != err {
		return nil
	}
	return reply
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
	if !flag {
		return nil, &exceptions.Error{Msg: "no change", Code: 400}
	}
	return old.Modify()
}

/**
	查询包详细信息
*/
func GetOne(key int64, fields ...string) (*warehouse, error) {
	str := dbConnect.Select("file")
	if 0 == len(fields) {
		str.Fields("id, name, \"group\", remarks, version, uploadTime, modifyTime, status")
	} else {
		str.Fields(fields...)
	}
	str.AND("id = ?", "status = ?")
	one, err := dbConnect.WithQuery(str.Preview(), func(rows *sql.Rows) (interface{}, error) {
		target := new(warehouse)
		for rows.Next() {
			rows.Scan(&target.Id, &target.Name, &target.Group, &target.Remarks, &target.Version, &target.UploadTime, &target.ModifyTime, &target.Status)
		}
		return target, nil
	}, key, true)
	w := one.(warehouse)
	if 0 == w.Id {
		return nil, &exceptions.Error{Msg: "no such this package", Code: 404}
	}
	return &w, err
}
