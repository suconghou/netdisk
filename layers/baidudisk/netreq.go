package baidudisk

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func httpGetResp(url string) (*http.Response, error) {
	response, err := http.Get(url)
	return response, err
}

func httpGet(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return body, err
	}
	return body, nil
}

func httpPost(url string, contentType string, body io.Reader) ([]byte, error) {
	response, err := http.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bodyStr, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return bodyStr, err
	}
	return bodyStr, nil
}

func httpPostWait(url string, contentType string, body io.Reader) ([]byte, error) {
	response, err := http.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	time.Sleep(time.Second)
	bodyStr, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return bodyStr, err
	}
	return bodyStr, nil
}
