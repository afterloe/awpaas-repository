package soaClient

import (
	"../logger"
	"net/http"
	"net"
	"time"
	"fmt"
	"io"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"strings"
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

func ObjectToJson(vol interface{}) string {
	buf, err := json.Marshal(vol)
	if nil != err {
		return ""
	}
	return string(buf)
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

func Call(method, serviceName, url string, body io.Reader, header map[string]string) (map[string]interface{}, error) {
	reqUrl := fmt.Sprintf("http://%s%s", serviceName, url)
	remote, err := http.NewRequest(method, reqUrl, body)
	for key, value := range header {
		remote.Header.Add(key, value)
	}
	if nil != err {
	}
	response, err := GeneratorClient().Do(remote)
	if err != nil && response == nil {
		logger.Error("soa-client", fmt.Sprintf("forward %+v", err))
		return nil, err
	} else {
		defer response.Body.Close()
		logger.Logger("soa-client", fmt.Sprintf("%3d | %-7s | %s", response.StatusCode, method, reqUrl))
		reply, err := ioutil.ReadAll(response.Body)
		if nil != err {
			return nil, err
		}
		return JsonToObject(string(reply))
	}
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

func GeneratorPostHeader() map[string]string {
	return map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
}

func GeneratorBody(vol interface{}) io.Reader {
	buf, err := json.Marshal(vol)
	if nil != err {
		return nil
	}
	return strings.NewReader(string(buf))
}