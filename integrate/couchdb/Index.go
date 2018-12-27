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
	addr, port string
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
}

func Get(dbName string, params map[string]interface{}) (map[string]interface{}, error) {
	reqUrl := fmt.Sprintf("http://%s/%s?%s", host, dbName, soaClient.Encode(params))
	remote, err := http.NewRequest("GET", reqUrl, nil)
	remote.AddCookie(&http.Cookie{Name: key, Value:value, HttpOnly: true})
	if nil != err {
		return nil, err
	}
	return soaClient.Invoke(remote, "couchDB-sdk", nil)
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
