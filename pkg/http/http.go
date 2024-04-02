package infrastructures

import (
	"bytes"
	"demogo/config"
	"demogo/pkg/constant"
	"demogo/pkg/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	httpContentType   = "Content-type"
	httpAuthorization = "Authorization"

	contentTypeJSON           = "application/json"
	contentTypeXML            = "application/xml"
	contentTypeFormURLEncoded = "application/x-www-form-urlencoded"
	contentTypeTextPlain      = "text/plain"
)

type HTTPCall struct {
	Conf   *config.SectionHTTP
	client *http.Client
	once   *sync.Once
}

// NewHTTPCall init HTTPCall
func NewHTTPCall(conf *config.SectionHTTP) *HTTPCall {
	timeout := time.Duration(conf.Timeout) * time.Second

	var transport http.RoundTripper = &http.Transport{
		DisableKeepAlives: conf.DisableKeepAlive,
		MaxIdleConns:      10,
		IdleConnTimeout:   10 * time.Second,
	}
	client := &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}

	return &HTTPCall{
		Conf:   conf,
		client: client,
	}
}

// CallService is a function to call another service
func (h *HTTPCall) CallService(method, url string, requestBody []byte) (string, *constant.ApplicationError) {
	if h.client == nil {
		timeout := time.Duration(h.Conf.Timeout) * time.Second

		var transport http.RoundTripper = &http.Transport{
			DisableKeepAlives: h.Conf.DisableKeepAlive,
			MaxIdleConns:      10,
			IdleConnTimeout:   10 * time.Second,
		}
		client := &http.Client{
			Timeout:   timeout,
			Transport: transport,
		}

		h.client = client
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", utils.GetErrorResponse(err.Error(), http.StatusBadRequest, constant.StatusBadRequest)
	}

	// legacy code
	request.Header.Set(httpContentType, contentTypeJSON)

	response, errResponse := h.client.Do(request)
	if errResponse != nil {
		if os.IsTimeout(errResponse) {
			errTimeout := fmt.Sprintf("request error timeout: %s", errResponse.Error())
			return "", utils.GetErrorResponse(errTimeout, http.StatusBadRequest, constant.StatusBadRequest)
		}
		return "", utils.GetErrorResponse(errResponse.Error(), http.StatusBadRequest, constant.StatusBadRequest)
	}

	defer response.Body.Close()
	body, errBody := io.ReadAll(response.Body)
	if errBody != nil {
		return "", utils.GetErrorResponse(errBody.Error(), http.StatusExpectationFailed, constant.StatusExpectationFailed)
	}

	if response.StatusCode != http.StatusOK {
		var errResponse constant.ApplicationError
		if err := json.Unmarshal(body, &errResponse); err != nil {
			return "", utils.GetErrorResponse(err.Error(), http.StatusExpectationFailed, constant.StatusExpectationFailed)
		}
		return "", &errResponse
	}

	return string(body), nil
}

// CallServiceByte is a function to call another service return byte data
func (h *HTTPCall) CallServiceByte(method, url string, requestBody []byte) ([]byte, *constant.ApplicationError) {
	if h.client == nil {
		timeout := time.Duration(h.Conf.Timeout) * time.Second

		var transport http.RoundTripper = &http.Transport{
			DisableKeepAlives: false,
			MaxIdleConns:      10,
			IdleConnTimeout:   10 * time.Second,
		}
		client := &http.Client{
			Timeout:   timeout,
			Transport: transport,
		}

		h.client = client
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, utils.GetErrorResponse(err.Error(), http.StatusBadRequest, constant.StatusBadRequest)
	}

	// legacy code
	request.Header.Set(httpContentType, contentTypeJSON)

	response, errResponse := h.client.Do(request)
	if errResponse != nil {
		if os.IsTimeout(errResponse) {
			errTimeout := fmt.Sprintf("request error timeout: %s", errResponse.Error())
			return nil, utils.GetErrorResponse(errTimeout, http.StatusBadRequest, constant.StatusBadRequest)
		}
		return nil, utils.GetErrorResponse(errResponse.Error(), http.StatusBadRequest, constant.StatusBadRequest)
	}

	if response != nil {
		// defer io.Copy(io.Discard, response.Body)
		defer response.Body.Close()
	}

	// var buf bytes.Buffer
	// buf := bytes.NewBuffer(nil)
	// n, err := io.Copy(buf, response.Body)
	// io.Copy(buf, response.Body)
	// body, errBody := io.ReadAll(response.Body)
	// io.Copy(io.Discard, response.Body)

	buf := bytes.NewBuffer(nil)
	_, errBody := io.Copy(buf, response.Body)
	if errBody != nil {
		return nil, utils.GetErrorResponse(errBody.Error(), http.StatusExpectationFailed, constant.StatusExpectationFailed)
	}

	if response.StatusCode != http.StatusOK {
		var errResponse constant.ApplicationError
		if err := json.Unmarshal(buf.Bytes(), &errResponse); err != nil {
			return nil, utils.GetErrorResponse(err.Error(), http.StatusExpectationFailed, constant.StatusExpectationFailed)
		}
		return nil, &errResponse
	}

	return buf.Bytes(), nil
}
