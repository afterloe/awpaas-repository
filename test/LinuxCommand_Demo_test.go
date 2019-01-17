package test

import (
	"testing"
	"os"
	"os/exec"
	"bufio"
	"io"
	"fmt"
)

func Test_Demo_demo(t *testing.T) {
	//errorNotFound := errors.New("executable file not found in $PATH")
	commandList := [3]string{"tar -xzvf risk-point-0.0.5.tar.gz", "docker build -t 127.0.0.1/ascs/risk-point:0.0.5",
	"docker push 127.0.0.1/ascs/risk-point:0.0.5"}
	dir := "/tmp/download/9E04A3DCA7954DC39488D15ECD1CBE36"
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
	stdout, err := cmd.StdoutPipe()
	if nil != err {
		t.Error(err)
		t.Error("exec cmd.sh failed.")
		return
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}
	cmd.Wait()
	t.Log("done")
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