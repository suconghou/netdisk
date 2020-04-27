package middleware

import (
	"net/http"
	"regexp"

	"github.com/suconghou/netdisk/config"
	"github.com/suconghou/netdisk/layers/baidudisk"
	"github.com/suconghou/netdisk/util"
)

var netroute = []routeInfo{
	{regexp.MustCompile(`^ls/(.+)$`), ls},
	{regexp.MustCompile(`^info/(.+)$`), info},
}

// NetStreamAPI response json data
func NetStreamAPI(w http.ResponseWriter, r *http.Request, match []string) error {
	return dispatch(w, r, match, netroute, func(w http.ResponseWriter, r *http.Request, match []string) error {
		return get(w, r, match)
	})
}

func ls(w http.ResponseWriter, r *http.Request, match []string) error {
	file := match[1]
	url := baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).APILsURL(file)
	return util.ProxyURL(w, r, url, nil)
}

func info(w http.ResponseWriter, r *http.Request, match []string) error {
	file := match[1]
	url := baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).APIFileInfoURL(file)
	return util.ProxyURL(w, r, url, nil)
}

func get(w http.ResponseWriter, r *http.Request, match []string) error {
	file := match[1]
	url := baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).GetDownloadURL(file)
	return util.ProxyURL(w, r, url, nil)
}
