package middleware

import (
	"fmt"
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
	_, err := fastload.Pipe(w, r, url, func(w http.ResponseWriter, header http.Header) {
		w.Header().Del("Cookie")
	})
	if err != nil {
		util.Log.Printf("pipe %s error:%s", url, err)
	}
}
