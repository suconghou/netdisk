package multiget

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"netdisk/request"
)

var (
	ua = []string{"netdisk;2.2.2;pc;pc-mac;10.14.5;macbaiduyunguanjia"}
)

// 获取可能的302后的地址,同时获取大小,调用者需检查大小,调用者需必要时关闭reader
func resURI(target string) (*url.URL, int64, io.ReadCloser, error) {
	res, _, err := request.Request(context.Background(), target, http.MethodGet, nil, 15, nil, "")
	if err != nil {
		return nil, 0, nil, err
	}
	var uri = res.Request.URL
	return uri, res.ContentLength, res.Body, nil
}
