package middleware

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/suconghou/fastload/multiget"
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
	url := baidudisk.NewClient(config.Cfg.Token, config.Cfg.Root).APILsURL(file)
	return util.ProxyURL(w, r, url, nil)
}

func info(w http.ResponseWriter, r *http.Request, match []string) error {
	file := match[1]
	url := baidudisk.NewClient(config.Cfg.Token, config.Cfg.Root).APIFileInfoURL(file)
	return util.ProxyURL(w, r, url, nil)
}

func get(w http.ResponseWriter, r *http.Request, match []string) error {
	file := match[1]
	url := baidudisk.NewClient(config.Cfg.Token, config.Cfg.Root).GetDownloadURL(file)
	url = strings.ReplaceAll(url, "qdall01.baidupcs.com", "qdcu02.baidupcs.com")
	addr := map[string][]string{
		"qdall01.baidupcs.com:80": []string{
			"119.167.143.21:80",
			"119.167.143.20:80",
			"140.249.34.22:80",
			"119.167.143.12:80",
			"111.63.70.21:80",
		},
		"qdcu02.baidupcs.com": []string{
			"119.167.143.21:80",
			"119.167.143.20:80",
			"140.249.34.22:80",
			"119.167.143.12:80",
			"111.63.70.21:80",
		},
	}
	rr, size, err := multiget.NewDNS(url, 0, 40090163, addr)
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
