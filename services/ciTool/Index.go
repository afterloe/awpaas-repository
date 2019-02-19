package ciTool

import (
	"../../exceptions"
	"../dbConnect"
	"../warehouse"
	"fmt"
	"os"
	"time"
	"os/exec"
	"database/sql"
)

const (
	selectCISQL = "SELECT id, registryType, lastReport, lastCITime FROM ci WHERE status = ? AND warehouseId = ?"
	selectCMDSQL = "SELECT id,context FROM cmd WHERE status = ? AND cid = ? ORDER BY id AES"
)

func init()  {
	registryType = [4]string{"code", "image", "tar", "soa-jvm"}
}

func GetRegistryType() interface{} {
	return registryType
}

func DefaultCmd(id int64, inputType string, content ...string) (*cmd, error) {
	_, err := warehouse.GetOne(id, "id")
	if nil != err {
		return nil, &exceptions.Error{Msg: "no such this package", Code: 404}
	}
	for _, t := range registryType {
		if t == inputType {
			return &cmd{id, inputType, content, "", 0}, nil
		}
	}
	return nil, &exceptions.Error{Msg: "no such this type", Code: 400}
}

func FindCIInfo(id int64) (*ci, error) {
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

func AppendCI(ci *cmd) (interface{}, error) {
	if nil == ci {
		return nil, &exceptions.Error{Msg: "cmd not found", Code: 400}
	}
	// TODO
}

/**
	软件构建

	1.查询源文件
	2.判断ci类型
	3.按照类型进行分发处理
 */
func Build(warehouseId int64) (interface{}, error) {
	_, err := warehouse.GetOne(warehouseId, "id")
	if nil != err {
		return nil, &exceptions.Error{Msg: "no such this package", Code: 404}
	}
	// TODO
	cmd := w.Cmd
	cmd.LastCiTime = time.Now().Unix()
	switch cmd.RegistryType {
	case "tar":
		task := util.GeneratorUUID()
		context := "/tmp/download/" + task
		go func() {
			rep, _ := execShell(context, cmd.Content...)
			cmd.LastReport = rep.(string)
			w.Cmd = cmd
			w.Modify()
		}()
		return task, nil
	case "image":
		return nil, nil
	case "code":
		return nil, nil
	case "soa-jvm":
		task := util.GeneratorUUID()
		packageInfo := w.PackInfo
		context := packageInfo.GeneratorSavePath() + "/" + task
		go func() {
			rep, _ :=execShell(context, []string{
				fmt.Sprintf("cp ../%s ./", packageInfo.Name),
				fmt.Sprintf("tar -xzvf %s", packageInfo.Name),
				fmt.Sprintf("docker build -t %s/%s/%s:%s .", "127.0.0.1", w.Group, w.Name, w.Version),
				fmt.Sprintf("docker push %s/%s/%s:%s", "127.0.0.1", w.Group, w.Name, w.Version),
			}...)
			cmd.LastReport = rep.(string)
			w.Cmd = cmd
			w.Modify()
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