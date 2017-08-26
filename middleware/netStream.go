package middleware

import (
	"net/http"
	"regexp"

	"github.com/suconghou/fastload/fastload"
	"github.com/suconghou/netdisk/config"
	"github.com/suconghou/netdisk/layers/baidudisk"
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
	fastload.Pipe(w, cleanHeader(r, xheaders), url, usecachefilter, 30, nil)
}

func info(w http.ResponseWriter, r *http.Request, match []string) {
	file := match[1]
	url := baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).APIFileInfoURL(file)
	fastload.Pipe(w, cleanHeader(r, xheaders), url, usecachefilter, 30, nil)
}

func get(w http.ResponseWriter, r *http.Request, match []string) {
	file := match[1]
	url := baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).GetDownloadURL(file)
	fastload.FastPipe(w, cleanHeader(r, xheaders), url, 2, 262144, nil, usecachefilter, nil)
}
