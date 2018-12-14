package httpClient

import (
	"net/http"
	"net"
	"time"
	"fmt"
	"io"

	"../logger"
	"io/ioutil"
)

var (
	key string
	needToken bool
	daemonAddr string
	MaxIdleConns int
	MaxIdleConnsPerHost int
	IdleConnTimeout int
)

func init() {
	MaxIdleConns = 100
	MaxIdleConnsPerHost = 100
	IdleConnTimeout = 90
}

func generatorClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{ Timeout: 30 * time.Second,}).DialContext,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			IdleConnTimeout:	 time.Duration(IdleConnTimeout)* time.Second,
		},
		Timeout: 20 * time.Second,
	}
	return client
}

func Call(method, url string, body io.Reader, header map[string]string) string {
	remote, err := http.NewRequest(method, url, body)
	reply := ""
	for key, value := range header {
		remote.Header.Add(key, value)
	}
	response, err := generatorClient().Do(remote)
	if err != nil && response == nil {
		logger.Error("daemon", fmt.Sprintf("forward %+v", err))
	} else {
		defer response.Body.Close()
		logger.Logger("daemon", fmt.Sprintf("%3d | %-7s | %s", response.StatusCode, method, url))
		buf, _ := ioutil.ReadAll(response.Body)
		reply = string(buf)
	}
	return reply
}