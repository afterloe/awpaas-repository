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
	selectCISQL = "SELECT id, registryType, lastReport, lastCITime FROM ci WHERE status = ? AND id = ? ORDER BY createTime DESC"
	selectCMDSQL = "SELECT id,context FROM cmd WHERE status = ? AND cid = ? ORDER BY id AES"
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
	_, err := warehouse.GetOne(id, "id")
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
	_, err := warehouse.GetOne(id, "id")
	if nil != err {
		return nil, &exceptions.Error{Msg: "no such this package", Code: 404}
	}
	return dbConnect.Select("ci").Fields("id", "registryType", "lastReport", "lastCITime", "createTime").
		AND("warehouseId = ?", "status = ?").Query(id, true)
}

func GetOne(id int64) (*ci, error) {
	p, err := dbConnect.WithTransaction(func(tx *sql.Tx) (interface{}, error) {
		ciInfo, _ := tx.Prepare(selectCISQL)
		rows, _ := ciInfo.Query(true, id)
		target := new(ci)
		for rows.Next() {
			rows.Scan(&target.Id, &target.RegistryType, &target.LastReport, &target.LastCiTime)
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
		saveCmd, _ := tx.Prepare("INSERT INTO cmd(cid, context, createTime, status) VALUES (?, ?, ?)")
		for v := range cmds {
			saveCmd.Exec(r.LastInsertId(), v, time.Now().Unix(), true)
		}
		return nil, nil
	})
	return "APPEND SUCCESS", nil
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
	switch ci.RegistryType {
	case "tar":
		task := util.GeneratorUUID()
		context := tmpDIR + task
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
		task := util.GeneratorUUID()
		w, err := warehouse.GetOne(ci.WarehouseId, "fid", "name")
		if nil != err {
			return nil, err
		}
		f, err := borderSystem.GetOne(w.FId)
		if nil != err {
			return nil, err
		}
		context := tmpDIR + task
		_, err = borderSystem.Copy(f.Id, context)
		if nil != err {
			return nil, err
		}
		go func() {
			rep, _ := execShell(context, []string{
				fmt.Sprintf("tar -xzvf %s", f.Name),
				fmt.Sprintf("docker build -t %s/%s/%s:%s .", "127.0.0.1", w.Group, w.Name, w.Version),
				fmt.Sprintf("docker push %s/%s/%s:%s", "127.0.0.1", w.Group, w.Name, w.Version),
			}...)
			ci.LastReport = rep.(string)
			ci.Update()
			os.RemoveAll(context)
		}()
		return nil, nil
	default:
		return nil, &exceptions.Error{Msg: "can't supper this"}
	}
}

func execShell(dir string, args ...string) (interface{}, error) {
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