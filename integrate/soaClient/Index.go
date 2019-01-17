package soaClient

import (
	"../logger"
	"net/http"
	"net"
	"time"
	"fmt"
	"io"
	"io/ioutil"
	"../../exceptions"
	"encoding/json"
	"net/url"
	"strings"
	"os"
)

var (
	maxIdleConn,
	maxIdleConnPerHost,
	idleConnTimeout int
)

func init() {
	maxIdleConn = 100
	maxIdleConnPerHost = 100
	idleConnTimeout = 90
}

func JsonToObject(chunk string) (map[string]interface{}, error){
	rep := make(map[string]interface{})
	err := json.Unmarshal([]byte(chunk), &rep)
	if nil != err {
		return nil, err
	}
	return rep, nil
}

func Encode(params map[string]interface{}) string {
	context := url.Values{}
	for key, value := range params {
		context.Add(key, value.(string))
	}
	return context.Encode()
}

func Invoke(remote *http.Request, module string, callback func(response *http.Response) (map[string]interface{}, error)) (map[string]interface{}, error) {
	response, err := GeneratorClient().Do(remote)
	if err != nil && response == nil {
		logger.Error(module, fmt.Sprintf("forward %+v", err))
		return nil, err
	} else {
		defer response.Body.Close()
		logger.Logger(module, fmt.Sprintf("%3d | %-7s | %s", response.StatusCode, remote.Method,
			remote.URL))
		if nil == callback {
			return defaultFunc(response)
		} else {
			return callback(response)
		}
	}
}

func defaultFunc(response *http.Response) (map[string]interface{}, error) {
	reply, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return nil, err
	}
	return JsonToObject(string(reply))
}

func Call(method, serviceName, url string, body io.Reader, header map[string]string) (map[string]interface{}, error) {
	reqUrl := fmt.Sprintf("http://%s%s", serviceName, url)
	remote, err := http.NewRequest(method, reqUrl, body)
	for key, value := range header {
		remote.Header.Add(key, value)
	}
	if nil != err {
		return nil, err
	}
	return Invoke(remote, "soa-client", defaultFunc)
}

func GeneratorClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{ Timeout: 30 * time.Second,}).DialContext,
			MaxIdleConns:        maxIdleConn,
			MaxIdleConnsPerHost: maxIdleConnPerHost,
			IdleConnTimeout:	 time.Duration(idleConnTimeout)* time.Second,
		},
		Timeout: 30 * time.Second,
	}
	return client
}

func GeneratorBody(vol interface{}) io.Reader {
	buf, err := json.Marshal(vol)
	if nil != err {
		return nil
	}
	return strings.NewReader(string(buf))
}

func DownloadFile(url,savePath string) (map[string]interface{}, error) {
	remote, err := http.NewRequest("GET", url, nil)
	if nil != err {
		return nil, err
	}
	return Invoke(remote, "soa-client", func(response *http.Response) (map[string]interface{}, error) {
		if 200 != response.StatusCode {
			return nil, &exceptions.Error{Msg: "download failed.", Code: 500}
		}
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
		return map[string]interface{} {
			"savePath": savePath,
			"fileName": filename,
		}, nil
	})
}