package ciTool

import (
	"../../exceptions"
	"../dbConnect"
	"../warehouse"
	"../borderSystem"
	"../../util"
	"fmt"
	"os"
	"time"
	"os/exec"
	"database/sql"
)

const (
	selectCISQL = "SELECT id, registryType, lastReport, lastCITime, status FROM ci WHERE status = ? AND id = ? ORDER BY createTime DESC"
	selectCMDSQL = "SELECT id,context FROM cmd WHERE status = ? AND cid = ? ORDER BY createTime DESC"
	updateCISQL = "UPDATE ci SET registryType = ?, lastReport = ?, lastCITime = ?, status = ? WHERE id = ?"
	tmpDIR = "/tmp/download/"
)

func init()  {
	registryType = [4]string{"code", "image", "tar", "soa-jvm"}
}

func (this *ci) Update() {
	dbConnect.WithPrepare(updateCISQL, func(stmt *sql.Stmt) (interface{}, error) {
		stmt.Exec(this.RegistryType, this.LastReport, this.LastCiTime, this.Status, this.Id)
		return nil, nil
	})
}

func GetRegistryType() interface{} {
	return registryType
}

func DefaultCmd(id int64, inputType string, content ...string) ([]*cmd, error) {
	_, err := warehouse.GetOne(id)
	if nil != err {
		return nil, &exceptions.Error{Msg: "no such this package", Code: 404}
	}
	flag := false
	for _, t := range registryType {
		if t == inputType {
			flag = true
		}
	}
	if !flag {
		return nil, &exceptions.Error{Msg: "no such this type", Code: 400}
	}
	context := make([]*cmd, 0)
	for _, v := range content {
		context = append(context, &cmd{Context: v, CreateTime: time.Now().Unix(), Status: true})
	}
	return context, nil
}

func CIList(id int64) (interface{}, error) {
	_, err := warehouse.GetOne(id)
	if nil != err {
		return nil, &exceptions.Error{Msg: "no such this package", Code: 404}
	}
	return dbConnect.Select("ci").Fields("id", "registryType", "lastReport", "lastCITime", "createTime").
		AND("warehouseId = ?", "status = ?").OrderBy("createTime desc").Query(id, true)
}

func GetOne(id int64) (*ci, error) {
	p, err := dbConnect.WithTransaction(func(tx *sql.Tx) (interface{}, error) {
		ciInfo, _ := tx.Prepare(selectCISQL)
		rows, _ := ciInfo.Query(true, id)
		target := new(ci)
		flag := new(int64)
		for rows.Next() {
			rows.Scan(&target.Id, &target.RegistryType, &target.LastReport, &target.LastCiTime, &flag)
			if 1 == *flag {
				target.Status = true
			}
		}
		if 0 == target.Id {
			return target, &exceptions.Error{Msg: "no such this warehouse", Code: 404}
		}
		cmds, _ := tx.Prepare(selectCMDSQL)
		rows, _ = cmds.Query(true, target.Id)
		context := make([]*cmd, 0)
		for rows.Next() {
			c := new(cmd)
			rows.Scan(&c.Id, &c.Context)
			context = append(context, c)
		}
		target.Content = context
		return target, nil
	})
	if nil != err {
		return nil, err
	}
	return p.(*ci), nil
}

func AppendCI(warehouseId int64, fileType string, cmds []*cmd) (interface{}, error) {
	dbConnect.WithTransaction(func(tx *sql.Tx) (interface{}, error) {
		saveCI, _ := tx.Prepare("INSERT INTO ci(warehouseId, registryType, createTime, status) VALUES (?,?,?,?)")
		r, _ := saveCI.Exec(warehouseId, fileType, time.Now().Unix(), true)
		saveCmd, _ := tx.Prepare("INSERT INTO cmd(cid, context, createTime, status) VALUES (?, ?, ?, ?)")
		lastId, _ := r.LastInsertId()
		for _, v := range cmds {
			context := v.Context
			_, err := saveCmd.Exec(lastId, context, time.Now().Unix(), true)
			if nil != err {
				fmt.Println(err)
				tx.Rollback()
			}
		}
		return nil, nil
	})
	return "APPEND SUCCESS", nil
}

func executeHistory(cid int64, taskId, context string) {
	dbConnect.WithPrepare("INSERT INTO ci_history(cid, taskId, context, createTime) VALUES (?, ?, ?, ?)", func(stmt *sql.Stmt) (interface{}, error) {
		stmt.Exec(cid, taskId, context, time.Now().Unix())
		return nil, nil
	})
}

func GetDetail(key string) (map[string]interface{}, error){
	p, err := dbConnect.WithQuery("SELECT context FROM ci_history WHERE taskId = ?", func(rows *sql.Rows) (interface{}, error) {
		context := new(string)
		rows.Next()
		rows.Scan(&context)
		return context, nil
	}, key)
	if nil != err {
		return nil, err
	}
	context := p.(*string)
	path := *context + "/report.log"
	stat, err := os.Stat(path)
	if nil != err {
		return nil, &exceptions.Error{Msg: "file has delete", Code: 404}
	}
	return map[string]interface{}{
		"name": stat.Name(),
		"size": stat.Size(),
		"path": path,
	}, nil
}

func CIHistoryDetail(warehouseId int64) (interface{}, error) {
	return dbConnect.WithQuery("SELECT taskId, ci_history.createTime FROM ci_history JOIN ci ON ci.id = ci_history.cid AND ci.status = ? JOIN warehouse ON warehouse.id = ci.warehouseId AND ci.warehouseId = ?", func(rows *sql.Rows) (interface{}, error) {
		return dbConnect.ToMap(rows)
	}, true, warehouseId)
}

func History(begin, limit int) (interface{}, error) {
	return dbConnect.WithQuery("SELECT warehouse.name, warehouse.\"group\", warehouse.version, ci.id AS cid, ci_history.id AS hid, ci_history.taskId, ci_history.createTime FROM warehouse JOIN ci ON ci.warehouseId = warehouse.id AND ci.status = ? JOIN ci_history ON ci_history.cid = ci.id  LIMIT ? OFFSET ?", func(rows *sql.Rows) (interface{}, error) {
		return dbConnect.ToMap(rows)
	}, true, limit, begin)
}

/**
	软件构建

	1.查询源文件
	2.判断ci类型
	3.按照类型进行分发处理
 */
func Run(ciId int64) (interface{}, error) {
	ci, err := GetOne(ciId)
	if nil != err {
		return nil, &exceptions.Error{Msg: "no such this plain", Code: 404}
	}
	ci.LastCiTime = time.Now().Unix()
	task := util.GeneratorUUID()
	context := tmpDIR + task
	executeHistory(ci.Id, task, context)
	switch ci.RegistryType {
	case "tar":
		go func() {
			shell := make([]string, 0)
			for _, v := range ci.Content {
				shell = append(shell, v.Context)
			}
			rep, _ := execShell(context, shell...)
			ci.LastReport = rep.(string)
			ci.Update()
		}()
		return task, nil
	case "image":
		return nil, nil
	case "code":
		return nil, nil
	case "soa-jvm":
		w, err := warehouse.GetOne(ci.WarehouseId, "fid", "name")
		if nil != err {
			return nil, err
		}
		f, err := borderSystem.GetOne(w.FId)
		if nil != err {
			return nil, err
		}
		_, err = borderSystem.Copy(f.Id, context)
		if nil != err {
			return nil, err
		}
		go func() {
			rep, _ := execShell(context, []string{
				fmt.Sprintf("tar -xzvf %s", f.Name),
				fmt.Sprintf("rm -rf %s", f.Name),
				fmt.Sprintf("docker build -t %s/%s/%s:%s .", "127.0.0.1", w.Group, w.Name, w.Version),
				fmt.Sprintf("docker push %s/%s/%s:%s", "127.0.0.1", w.Group, w.Name, w.Version),
			}...)
			ci.LastReport = rep.(string)
			ci.Update()
		}()
		return nil, nil
	default:
		return nil, &exceptions.Error{Msg: "can't supper this"}
	}
}

func execShell(dir string, args ...string) (interface{}, error) {
	os.MkdirAll(dir, os.ModePerm)
	sh, err := os.Create(dir + "/cmd.sh")
	if nil != err {
		return nil, &exceptions.Error{Msg: "create file error", Code: 500}
	}
	sh.WriteString("#!/bin/sh\n")
	for _, c := range args {
		sh.WriteString(c + "\n")
	}
	sh.Chmod(os.ModePerm)
	sh.Close()
	cmd := exec.Command("/bin/sh", "-c", "./cmd.sh 2>&1 | tee report.log")
	cmd.Dir = dir
	tpl, err := cmd.Output()
	if nil != err {
		report, _ := os.Open(dir + "/report.log")
		report.WriteString(err.Error())
		return nil, &exceptions.Error{Msg: err.Error(), Code: 500}
	}
	os.Remove(dir + "/cmd.sh")
	return string(tpl), nil
}