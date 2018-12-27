package couchdb

import (
	"../../config"
	"../soaClient"
	"../logger"
	"strings"
	"fmt"
	"net/http"
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

func UserInfo() {
	reqUrl := fmt.Sprintf("http://%s%s", host, "/_session")
	remote, err := http.NewRequest("GET", reqUrl, nil)
	remote.AddCookie(&http.Cookie{Name: key, Value:value, HttpOnly: true})
	if nil != err {
		return
	}
	response, err := soaClient.GeneratorClient().Do(remote)
	if err != nil && response == nil {
		logger.Error("couchDB-sdk", fmt.Sprintf("forward %+v", err))
		return
	} else {
		defer response.Body.Close()
		logger.Logger("couchDB-sdk", fmt.Sprintf("%3d | %-7s | %s", response.StatusCode, "POST", reqUrl))
		reply, err := ioutil.ReadAll(response.Body)
		if nil != err {
			return
		}
		logger.Logger("couchDB-sdk", string(reply))
	}
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
	response, err := soaClient.GeneratorClient().Do(remote)
	if err != nil && response == nil {
		logger.Error("couchDB-sdk", fmt.Sprintf("forward %+v", err))
		return false, err
	} else {
		defer response.Body.Close()
		if 200 == response.StatusCode {
			content := response.Header.Get("Set-Cookie")
			item := strings.Split(content, "; ")
			cookieInfo := strings.Split(item[0], "=")
			key = cookieInfo[0]
			value = cookieInfo[1]
		}
		logger.Logger("couchDB-sdk", fmt.Sprintf("%3d | %-7s | %s", response.StatusCode, "POST", reqUrl))
		reply, err := ioutil.ReadAll(response.Body)
		if nil != err {
			return false, err
		}
		logger.Logger("couchDB-sdk", string(reply))
		return true, nil
	}
}
