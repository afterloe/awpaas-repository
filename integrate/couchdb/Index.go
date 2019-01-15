package couchdb

import (
	"../../config"
	"../soaClient"
	"../logger"
	"strings"
	"fmt"
	"net/http"
	"../../exceptions"
	"io/ioutil"
)

var (
	addr, port, dbName string
	host, username, password string
	key, value string
)

func init() {
	db := config.GetByTarget(config.Get("services"), "db")
	addr = config.GetByTarget(db, "addr").(string)
	port = config.GetByTarget(db, "port").(string)
	host = addr + ":" + port
	username = config.GetByTarget(db, "username").(string)
	password = config.GetByTarget(db, "password").(string)
	dbName = config.GetByTarget(db, "database").(string)
	// 如果开启检测
	flg := config.GetByTarget(db, "ping")
	if nil != flg {
		if flg.(bool) {
			ping()
		}
	}
}

func ping() {
	remote, err := http.NewRequest("GET", fmt.Sprintf("http://%s", host), nil)
	if nil != err {
		logger.Error(err)
		return
	}
	soaClient.Invoke(remote, "couchDB-sdk", func(response *http.Response) (map[string]interface{}, error) {
		reply, err := ioutil.ReadAll(response.Body)
		if nil != err {
			return nil, nil
		}
		logger.Logger("couchDB-sdk", string(reply))
		return nil, nil
	})
}

func autoCfg(response *http.Response) (map[string]interface{}, error) {
	reply, _ := ioutil.ReadAll(response.Body)
	r, _ := soaClient.JsonToObject(string(reply))
	if 401 == response.StatusCode {
		return map[string]interface{}{"needLogin": true}, nil
	}
	return r, nil
}

func Delete(o ...*obj) (map[string]interface{}, error) {
	reqUrl := fmt.Sprintf("http://%s/%s/_bulk_docs", host, dbName)
	mapResult := make(map[string]interface{})
	mapResult["docs"] = o
	reTry:
	remote, err := http.NewRequest("POST", reqUrl, soaClient.GeneratorBody(mapResult))
	remote.AddCookie(&http.Cookie{Name: key, Value:value, HttpOnly: true})
	remote.Header.Add("Content-Type", "application/json")
	if nil != err {
		return nil, err
	}
	reply, err := soaClient.Invoke(remote, "couchDB-sdk", autoCfg)
	if nil != err {
		return nil, err
	}
	if nil != reply["needLogin"] {
		Login()
		goto reTry
	}
	return reply, err
}

func Find(conditions *condition) ([]interface{}, error) {
	reqUrl := fmt.Sprintf("http://%s/%s/_find", host, dbName)
	reTry:
	remote, err := http.NewRequest("POST", reqUrl, strings.NewReader(conditions.String()))
	remote.AddCookie(&http.Cookie{Name: key, Value:value, HttpOnly: true})
	remote.Header.Add("Content-Type", "application/json")
	if nil != err {
		return nil, err
	}
	reply, err := soaClient.Invoke(remote, "couchDB-sdk", autoCfg)
	if nil != err {
		return nil, err
	}
	if nil != reply["needLogin"] {
		Login()
		goto reTry
	}
	if nil == reply["docs"] {
		return []interface{}{}, err
	}
	return reply["docs"].([]interface{}), err
}

func ReadAll(params map[string]interface{}) (map[string]interface{}, error) {
	var reqUrl string
	if nil != params {
		reqUrl = fmt.Sprintf("http://%s/%s/_all_docs?%s", host, dbName, soaClient.Encode(params))
	} else {
		reqUrl = fmt.Sprintf("http://%s/%s/_all_docs", host, dbName)
	}
	reTry:
	remote, err := http.NewRequest("GET", reqUrl, nil)
	remote.AddCookie(&http.Cookie{Name: key, Value:value, HttpOnly: true})
	if nil != err {
		return nil, err
	}
	reply, err := soaClient.Invoke(remote, "couchDB-sdk", autoCfg)
	if nil != err {
		return nil, err
	}
	if nil != reply["needLogin"] {
		Login()
		goto reTry
	}
	return reply, nil
}

func CreateDB() (bool, error) {
	reqUrl := fmt.Sprintf("http://%s/%s", host, dbName)
	reTry:
	remote, err := http.NewRequest("PUT", reqUrl, nil)
	remote.AddCookie(&http.Cookie{Name: key, Value:value, HttpOnly: true})
	remote.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if nil != err {
		return false, err
	}
	reply, err := soaClient.Invoke(remote, "couchDB-sdk", autoCfg)
	if nil != err {
		return false, err
	}
	if nil != reply["needLogin"] {
		Login()
		goto reTry
	}
	return true, err
}

func getUUID(count int) (interface{}, error){
	remote, err := http.NewRequest("GET", fmt.Sprintf("http://%s/_uuids?count=%d", host, count), nil)
	if nil != err {
		return nil, err
	}
	reply, err := soaClient.Invoke(remote, "couchDB-sdk", nil)
	if nil != err {
		return nil, err
	}
	id := reply["uuids"].([]interface{})[0]
	return id,nil
}

func Update(id, vol interface{}) (map[string]interface{}, error) {
	reqUrl := fmt.Sprintf("http://%s/%s/%v", host, dbName, id)
	reTry:
	remote, err := http.NewRequest("PUT", reqUrl, soaClient.GeneratorBody(vol))
	remote.AddCookie(&http.Cookie{Name: key, Value:value, HttpOnly: true})
	remote.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if nil != err {
		return nil, err
	}
	reply, err := soaClient.Invoke(remote, "couchDB-sdk", nil)
	if "not_found" == reply["error"] {
		CreateDB()
		goto reTry
	}
	return reply, nil
}

func Create(vol interface{}) (map[string]interface{}, error) {
	id, _ := getUUID(1)
	return Update(id, vol)
}

func Login() (bool, error) {
	content := soaClient.Encode(map[string]interface{}{
		"name": username,
		"password": password,
	})
	reqUrl := fmt.Sprintf("http://%s%s", host, "/_session")
	remote, err := http.NewRequest("POST", reqUrl, strings.NewReader(content))
	remote.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if nil != err {
		return false, err
	}
	reply, err := soaClient.Invoke(remote, "couchDB-sdk", func(response *http.Response) (map[string]interface{}, error){
		if 200 == response.StatusCode {
			content := response.Header.Get("Set-Cookie")
			item := strings.Split(content, "; ")
			cookieInfo := strings.Split(item[0], "=")
			key = cookieInfo[0]
			value = cookieInfo[1]
		}
		reply, err := ioutil.ReadAll(response.Body)
		if nil != err {
			return map[string]interface{}{}, err
		}
		logger.Logger("couchDB-sdk", string(reply))
		return soaClient.JsonToObject(string(reply))
	})
	if nil != reply["error"] {
		return false, &exceptions.Error{Msg: reply["reason"].(string), Code: 400}
	} else {
		return true, nil
	}
}
