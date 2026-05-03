package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"netdisk/util"

	"github.com/suconghou/fastload/fastload"
)

// Pipe response stream
func Pipe(w http.ResponseWriter, r *http.Request, match []string) error {
	var url string
	if match[1] == "" {
		url = fmt.Sprintf("http:/%s", match[0])
	} else {
		url = strings.Replace(strings.TrimPrefix(match[0], "/"), ":/", "://", 1)
	}
	if r.URL.RawQuery != "" {
		url = url + "?" + r.URL.RawQuery
	}
	return util.ProxyURL(w, r, url, nil)
}

// Proxy is a http_proxy and just http_proxy server
func Proxy(w http.ResponseWriter, r *http.Request) error {
	return fastload.HTTPProxy(w, r)
}
