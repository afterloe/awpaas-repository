package test

import (
	"testing"
	"../integrate/soaClient"
	"../config"
	"../util"
	"net/http"
	"fmt"
	"os"
	"strings"
	"io"
)

func Test_HttpClient_demo(t *testing.T) {
	fsServiceName := config.GetByTarget(config.Get("custom"), "fsServiceName").(string)
	reqUrl := fmt.Sprintf("http://%s/%s", fsServiceName, "v1/download/731dfcb7dc1c6c34298abc319500074c")
	remote, err := http.NewRequest("GET", reqUrl, nil)
	if nil != err {
		t.Error(err)
		return
	}
	soaClient.Invoke(remote, "soa-client", func(response *http.Response) (map[string]interface{}, error) {
		if 200 != response.StatusCode {
			t.Error("download failed! ....")
			return nil, nil
		}
		savePath := "/tmp/download/" + util.GeneratorUUID()
		_, e := os.Stat(savePath)
		if nil != e {
			os.Mkdir(savePath, os.ModePerm)
		}
		head := response.Header.Get("Content-Disposition")
		filename := strings.Split(head, "attachment;filename=")[1]
		desFile, _ := os.Create(savePath + "/" + filename)
		io.Copy(desFile, response.Body)
		defer desFile.Close()
		defer response.Body.Close()
		return nil, nil
	})
}
	//condition := couchdb.Condition().Append("_id", "$eq", "768edfd797eb8200d853b13632000e63").
	//	Append("Status", "$eq", true)
	//	//Fields("Name", "_id", "Size")
	//
	//t.Log(condition.String())
	//reply, _ := couchdb.Find("file-system", condition)
	//t.Log(reply)
//}

//func Test_HttpClient_demo(t *testing.T) {
//	reply, err := couchdb.Read("demoi/_all_docs?include_docs=true", nil)
//	if nil != err {
//		t.Error(reply)
//	}
//	t.Log(reply)
//}

//func Test_HttpClient_demo(t *testing.T) {
//	reply, err := couchdb.Create("demoi", map[string]interface{}{
//		"name": "afterloe",
//		"age": 6,
//		"sex": "小男孩",
//	})
//	if nil != err {
//		t.Error(reply)
//	}
//	t.Log(reply)
//}

//func Test_HttpClient_demo(t *testing.T) {
//	reply, err := couchdb.CreateDB("demo")
//	if nil != err {
//		t.Error(reply)
//	}
//	t.Log(reply)
//}

//func Test_HttpClient_demo(t *testing.T) {
//	flag, err := couchdb.Login()
//	if nil != err {
//		t.Error(err)
//		return
//	}
//	t.Log(flag)
//	reply, err := couchdb.Read("_session", nil)
//	if nil != err {
//		t.Error(reply)
//	}
//	t.Log(reply)
//}