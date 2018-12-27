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

func Read(dbName string, params map[string]interface{}) (map[string]interface{}, error) {
	reqUrl := fmt.Sprintf("http://%s/%s?%s", host, dbName, soaClient.Encode(params))
	remote, err := http.NewRequest("GET", reqUrl, nil)
	remote.AddCookie(&http.Cookie{Name: key, Value:value, HttpOnly: true})
	if nil != err {
		return nil, err
	}
	return soaClient.Invoke(remote, "couchDB-sdk", nil)
}

func CreateDB(dbName string) (bool, error) {
	reqUrl := fmt.Sprintf("http://%s/%s", host, dbName)
	reTry:
	remote, err := http.NewRequest("PUT", reqUrl, nil)
	remote.AddCookie(&http.Cookie{Name: key, Value:value, HttpOnly: true})
	remote.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if nil != err {
		return false, err
	}
	_, err = soaClient.Invoke(remote, "couchDB-sdk", func(response *http.Response) (map[string]interface{}, error) {
		reply, _ := ioutil.ReadAll(response.Body)
		r, _ := soaClient.JsonToObject(string(reply))
		if 201 == response.StatusCode {
			return nil, nil
		} else {
			return nil, &exceptions.Error{Code: response.StatusCode, Msg: r["reason"].(string)}
		}
	})
	if nil == err {
		return true, err
	}
	if 401 == (err).(*exceptions.Error).Code {
		Login()
		goto reTry
	}
	return false, err
}

func Create(dbName string, vol interface{}) (map[string]interface{}, error) {
	reply, _ := soaClient.Call("GET", host, "/_uuids?count=1", nil, nil)
	id := reply["uuids"].([]interface{})[0]
	reqUrl := fmt.Sprintf("http://%s/%s/%v", host, dbName, id)
	reTry:
	remote, err := http.NewRequest("PUT", reqUrl, soaClient.GeneratorBody(vol))
	remote.AddCookie(&http.Cookie{Name: key, Value:value, HttpOnly: true})
	remote.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if nil != err {
		return nil, err
	}
	reply, err = soaClient.Invoke(remote, "couchDB-sdk", nil)
	if "not_found" == reply["error"]{
		CreateDB(dbName)
		goto reTry
	}
	return reply, nil
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
