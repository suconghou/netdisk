package tools

import (
	"fmt"
	"net/http"

	"github.com/suconghou/fastload/fastload"
	"github.com/suconghou/netdisk/util"
)

// HTTPProxy nginx like reverse proxy
func HTTPProxy(port int, url string, transport *http.Transport) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requ := fmt.Sprintf("%s%s", url, r.RequestURI)
		fastload.Pipe(w, r, requ, func(out http.Header, res *http.Response) int {
			origin := r.Header.Get("Origin")
			if origin == "" {
				out.Set("Access-Control-Allow-Origin", "*")
			} else {
				out.Set("Access-Control-Allow-Origin", origin)
				out.Set("Access-Control-Allow-Credentials", "true")
			}
			out.Set("Access-Control-Max-Age", "86400")
			out.Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,OPTIONS")
			out.Set("Access-Control-Allow-Headers", "X-OAUTH-TOKEN,Cache-Control,Pragma,reqid,nid,host,x-real-ip,x-forwarded-ip,event-type,event-id,accept,content-type")
			util.Log.Printf("%d %s %s", res.StatusCode, r.Method, requ)
			return res.StatusCode
		}, 60, transport)
	})
	util.Log.Printf("Starting up on port %d\nProxy %s", port, url)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
