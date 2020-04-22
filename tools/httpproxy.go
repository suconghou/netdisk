package tools

import (
	"fmt"
	"net/http"

	"github.com/suconghou/fastload/fastload"
	"github.com/suconghou/netdisk/util"
	"github.com/suconghou/utilgo"
)

// HTTPProxy nginx like reverse proxy
func HTTPProxy(port int, url string, transport *http.Transport, str string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			utilgo.CrossShare(w.Header(), r.Header, str)
			http.Error(w, "ok", http.StatusOK)
			return
		}
		requ := fmt.Sprintf("%s%s", url, r.RequestURI)
		fastload.Pipe(w, r, requ, func(out *http.Header, res *http.Header, status int) int {
			w.Header().Set("Cache-Control", "public, max-age=604800")
			utilgo.CrossShare(w.Header(), r.Header, str)
			util.Log.Printf("%d %s %s", status, r.Method, requ)
			return status
		}, 60, transport)
	})
	util.Log.Printf("Starting up on port %d\nProxy %s", port, url)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
