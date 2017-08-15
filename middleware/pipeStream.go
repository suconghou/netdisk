package middleware

import (
	"fmt"
	"net"
	"net/http"
	"netdisk/util"
	"strings"

	"github.com/suconghou/fastload/fastload"
)

// Pipe response stream
func Pipe(w http.ResponseWriter, r *http.Request, match []string) {
	var url string
	if match[1] == "" {
		url = fmt.Sprintf("http:/%s", match[0])
	} else {
		url = strings.Replace(strings.TrimPrefix(match[0], "/"), ":/", "://", 1)
	}
	if r.URL.RawQuery != "" {
		url = url + "?" + r.URL.RawQuery
	}
	_, err := fastload.Pipe(w, r, url, usecachefilter, 3600, nil)
	if err != nil {
		util.Log.Printf("pipe %s error:%s", url, err)
	}
}

// Proxy is a http_proxy and just http_proxy server
func Proxy(w http.ResponseWriter, r *http.Request) (int64, error) {
	return fastload.Pipe(w, r, r.RequestURI, usecachefilter, 3600, nil)
}

// ProxySocks is a socks proxy server
func ProxySocks(c net.Conn) error {
	return nil
}
