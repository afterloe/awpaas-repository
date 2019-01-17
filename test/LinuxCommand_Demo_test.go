package test

import (
	"testing"
	"os"
	"os/exec"
)

func Test_Demo_demo(t *testing.T) {
	//errorNotFound := errors.New("executable file not found in $PATH")
	commandList := [3]string{"tar -xzvf risk-point-0.0.5.tar.gz", "docker build -t 127.0.0.1/ascs/risk-point:0.0.5",
	"docker push 127.0.0.1/ascs/risk-point:0.0.5"}
	dir := "/tmp/download/D51937E461CF4A699FE58E7E8BDF18F4"
	os.Remove(dir + "/cmd.sh")
	sh, err := os.Create(dir + "/cmd.sh")
	if nil != err {
		t.Error(err)
	}
	sh.WriteString("#!/bin/sh\n")
	for _, c := range commandList {
		sh.WriteString(c + "\n")
	}
	sh.Chmod(os.ModePerm)
	sh.Close()
	cmd := exec.Command("/bin/sh", "-c", "./cmd.sh 2>&1 | tee report.log")
	cmd.Dir = dir
	output, err := cmd.Output()
	if nil != err {
		t.Error(err)
	}
	t.Log(string(output))
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