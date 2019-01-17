package test

import (
	"testing"
	"os"
	"os/exec"
)

func Test_Demo_demo(t *testing.T) {
	//errorNotFound := errors.New("executable file not found in $PATH")
	commandList := [3]string{"ls", "pwd", "tar -xzvf risk-point-0.0.5.tar.gz"}
	dir := "/tmp/download/F1A0DA6DF9F34FF3B511900BA0FEA1A8"
	os.Remove(dir + "/cmd.sh")
	sh, err := os.Create(dir + "/cmd.sh")
	if nil != err {
		t.Error(err)
	}
	sh.WriteString("#!/bin/bash\n")
	for _, c := range commandList {
		sh.WriteString(c + "\n")
	}
	sh.Chmod(os.ModePerm)
	sh.Close()
	cmd := exec.Command("/bin/bash", "-c", "./cmd.sh > report.log")
	cmd.Dir = dir
	tpl, err := cmd.Output()
	if nil != err {
		t.Error(err)
		t.Error("exec cmd.sh failed.")
		return
	}
	t.Log(tpl)
}

/*

		cmd.Dir = dir
		cmd.Env = []string{"$PATH=" + dir}
		t.Log("Running command and waiting for it finish")
		tpl, err := cmd.Output()
		if nil != err {
			t.Error(err)
			t.Error("exec command " + c + " failed.")
			return
		}
		str := string(tpl)

	 */