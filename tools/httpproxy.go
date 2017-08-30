package tools

import (
	"fmt"
	"net/http"

	"github.com/suconghou/fastload/fastload"
	"github.com/suconghou/netdisk/util"
	"github.com/suconghou/utilgo"
)

// HTTPProxy nginx like reverse proxy
func HTTPProxy(port int, url string, transport *http.Transport) error {
	str := "Range,X-OAUTH-TOKEN,Cache-Control,Pragma,reqid,nid,host,x-real-ip,x-forwarded-ip,event-type,event-id,accept,content-type"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			utilgo.CrossShare(w.Header(), r.Header, str)
			http.Error(w, "ok", http.StatusOK)
			return
		}
		requ := fmt.Sprintf("%s%s", url, r.RequestURI)
		fastload.Pipe(w, r, requ, func(out http.Header, res *http.Response) int {
			utilgo.CrossShare(out, r.Header, str)
			util.Log.Printf("%d %s %s", res.StatusCode, r.Method, requ)
			return res.StatusCode
		}, 60, transport)
	})
	util.Log.Printf("Starting up on port %d\nProxy %s", port, url)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
