package utils

import (
	"demogo/config"
	"demogo/pkg/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	httpClient = &http.Client{}
)

func GetBaseURL(conf config.HTTPItemConfig) string {
	protocol := "http"
	if conf.TLS {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s", protocol, conf.Host)
}

func HTTPCall(method string, headers map[string]string, url string, body io.Reader, result interface{}) (err error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	// HTTP_DEBUG
	//byts, _ := httputil.DumpRequestOut(request, true)
	//log.Println(string(byts))

	response, err := httpClient.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	resp, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	logger.Log.Debugf("HTTP Response: %s", string(resp))

	return json.Unmarshal(resp, &result)

}

func FileSize(url string) (length int64, err error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return
	}

	return resp.ContentLength, nil
}
