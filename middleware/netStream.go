package middleware

import (
	"net/http"
	"netdisk/config"
	"netdisk/layers/baidudisk"
	"regexp"

	"github.com/suconghou/fastload/fastload"
)

var netroute = []routeInfo{
	{regexp.MustCompile(`^ls/(.+)$`), ls},
	{regexp.MustCompile(`^info/(.+)$`), info},
}

// NetStreamAPI response json data
func NetStreamAPI(w http.ResponseWriter, r *http.Request, match []string) {
	dispatch(w, r, match, netroute, func(w http.ResponseWriter, r *http.Request, match []string) {
		get(w, r, match)
	})
}

func ls(w http.ResponseWriter, r *http.Request, match []string) {
	file := match[1]
	url := baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).APILsURL(file)
	fastload.Pipe(w, r, url, usecachefilter, 30, nil)
}

func info(w http.ResponseWriter, r *http.Request, match []string) {
	file := match[1]
	url := baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).APIFileInfoURL(file)
	fastload.Pipe(w, r, url, usecachefilter, 30, nil)
}

func get(w http.ResponseWriter, r *http.Request, match []string) {
	file := match[1]
	url := baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).GetDownloadURL(file)
	fastload.Pipe(w, r, url, usecachefilter, 3600, nil)
}
