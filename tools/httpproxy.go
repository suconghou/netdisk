package tools

import (
	"fmt"
	"net/http"

	"netdisk/util"

	"github.com/suconghou/fastload/fastload"
	"github.com/suconghou/utilgo"
)

// HTTPProxy nginx like reverse proxy
func HTTPProxy(port int, url string, str string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			utilgo.CrossShare(w.Header(), r.Header, str)
			http.Error(w, "ok", http.StatusOK)
			return
		}
		requ := fmt.Sprintf("%s%s", url, r.RequestURI)
		_, err := fastload.Pipe(w, r, requ, func(out *http.Header, res *http.Header, status int) int {
			w.Header().Set("Cache-Control", "public, max-age=604800")
			utilgo.CrossShare(w.Header(), r.Header, str)
			util.Log.Printf("%d %s %s", status, r.Method, requ)
			return status
		}, 60)
		if err != nil {
			util.Log.Print(err)
		}
	})
	util.Log.Printf("Starting up on port %d\nProxy %s", port, url)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
