package warehouse

import (
	"../dbConnect"
	"../../integrate/logger"
	"../../exceptions"
	"../../config"
	"../../util"
	"time"
	"fmt"
	"os/exec"
	"os"
	"database/sql"
)

var (
	fsServiceName, root string
	insertSQL = "INSERT INTO warehouse(name, \"group\", remarks, version, packInfo, uploadTime, modifyTime, status) VALUES(?,?,?,?,?,?,?,?)"
	updateSQL = "UPDATE warehouse SET name = ?, \"group\" = ?, remarks = ?, version = ? packInfo = ?, modifyTime = ?, status = ? WHERE id = ?"
)

func init() {
	cfg := config.Get("custom")
	fsServiceName = config.GetByTarget(cfg, "fsServiceName").(string)
	root = config.GetByTarget(cfg, "root").(string)
	registryType = [4]string{"code", "image", "tar", "soa-jvm"}
}

/**
	保存至远程
*/
func (this *warehouse) SaveToDB() (map[string]interface{}, error) {
	return dbConnect.WithTransaction(func(tx *sql.Tx) (map[string]interface{}, error) {
		stmt, err := tx.Prepare(insertSQL)
		if nil != err {
			return nil, &exceptions.Error{Msg: "db stmt open failed.", Code: 500}
		}
		stmt.Exec(this.Name, this.Group, this.Remarks, this.Version, this.PackInfo, this.UploadTime, this.ModifyTime, this.Status)
		logger.Logger("warehouse", "insert success")
		return map[string]interface{}{}, nil
	})
}

func (this *warehouse) Modify() (map[string]interface{}, error) {
	this.ModifyTime = time.Now().Unix()
	if 0 == this.Id {
		return nil, &exceptions.Error{Msg: "no such this id", Code: 400}
	}
	return dbConnect.WithTransaction(func(tx *sql.Tx) (map[string]interface{}, error) {
		stmt, err := tx.Prepare(updateSQL)
		if nil != err {
			return nil, &exceptions.Error{Msg: "db stmt open failed.", Code: 500}
		}
		stmt.Exec(this.Name, this.Group, this.Remarks, this.Version, this.PackInfo, this.ModifyTime, this.Status, this.Id)
		logger.Logger("warehouse", "update success")
		return map[string]interface{}{}, nil
	})
}

func GetRegistryType() interface{} {
	return registryType
}

func DefaultCmd(inputType string, content ...string) (*cmd, error) {
	for _, t := range registryType {
		if t == inputType {
			return &cmd{inputType, content, "", 0}, nil
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
func GetList(begin, limit int) []map[string]interface{} {
	reply, err := dbConnect.Select("warehouse").
		Fields("id", "name", "uploadTime", "group").
		AND("status = ?").Page(begin, limit).Query(true)
	if nil != err {
		return nil
	}
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
	if !flag {
		return nil, &exceptions.Error{Msg: "no change", Code: 400}
	}
	return old.Modify()
}

/**
	查询包详细信息
*/
func GetOne(key string, fields ...string) (*warehouse, error) {
	str := dbConnect.Select("file")
	if 0 == len(fields) {
		str.Fields("id, name, \"group\", remarks, version, packInfo, uploadTime, modifyTime, status")
	} else {
		str.Fields(fields...)
	}
	str.AND("id = ?", "status = ?")
	one, err := dbConnect.WithQuery(str.Preview(), func(rows *sql.Rows) (interface{}, error) {
		target := new(warehouse)
		for rows.Next() {
			rows.Scan(&target.Id, &target.Name, &target.Group, &target.Remarks, &target.Version, &target.PackInfo, &target.UploadTime, &target.ModifyTime, &target.Status)
		}
		return target, nil
	}, key, true)
	w := one.(warehouse)
	if 0 == w.Id {
		return nil, &exceptions.Error{Msg: "no such this package", Code: 404}
	}
	return &w, err
}

func GeneratorFsFile(name, contentType string, size int64) *fsFile {
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

/**
	软件构建

	1.查询源文件
	2.判断ci类型
	3.按照类型进行分发处理
 */
func Build(w *warehouse) (interface{}, error) {
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
