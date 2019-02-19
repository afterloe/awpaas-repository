package ciTool

import (
	"../../exceptions"
	"fmt"
	"os"
	"time"
)

func init()  {
	registryType = [4]string{"code", "image", "tar", "soa-jvm"}
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

func AppendCI(w *warehouse, ci *cmd) (interface{}, error) {
	if nil == ci {
		return nil, &exceptions.Error{Msg: "cmd not found", Code: 400}
	}
	w.Cmd = *ci
	return w.Modify()
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
