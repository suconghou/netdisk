package middleware

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"netdisk/config"
	"netdisk/layers/baidudisk"
	"netdisk/util"

	"netdisk/multiget"
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
	url := baidudisk.NewClient(config.Cfg.Token, config.Cfg.Root).APILsURL(file)
	return util.ProxyURL(w, r, url, nil)
}

func info(w http.ResponseWriter, r *http.Request, match []string) error {
	file := match[1]
	url := baidudisk.NewClient(config.Cfg.Token, config.Cfg.Root).APIFileInfoURL(file)
	return util.ProxyURL(w, r, url, nil)
}

func get(w http.ResponseWriter, _ *http.Request, match []string) error {
	file := match[1]
	url := baidudisk.NewClient(config.Cfg.Token, config.Cfg.Root).GetDownloadURL(file)
	rr, size, err := multiget.Get(url, 0, 0, config.Cfg.Hosts)
	fmt.Println(size)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, rr)
	if err != nil {
		return err
	}
	return nil
}
