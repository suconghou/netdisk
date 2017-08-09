package middleware

import (
	"net/http"
	"netdisk/config"
	"netdisk/layers/baidudisk"
	"regexp"

	"github.com/suconghou/fastload/fastload"
	"github.com/suconghou/utilgo"
)

type routeInfo struct {
	Reg     *regexp.Regexp
	Handler func(http.ResponseWriter, *http.Request, []string)
}

var route = []routeInfo{
	{regexp.MustCompile(`^ls/(.+)$`), ls},
	{regexp.MustCompile(`^info/(.+)$`), info},
}

var filter = func(out http.Header, res *http.Response) int {
	utilgo.UseHTTPCache(out, 259200)
	return res.StatusCode
}

// NetStreamAPI response json data
func NetStreamAPI(w http.ResponseWriter, r *http.Request, match []string) {
	uri := match[1]
	found := false
	for _, p := range route {
		if p.Reg.MatchString(uri) {
			found = true
			p.Handler(w, r, p.Reg.FindStringSubmatch(uri))
			break
		}
	}
	if !found {
		fallback(w, r, match)
	}
}

func ls(w http.ResponseWriter, r *http.Request, match []string) {
	file := match[1]
	url := baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).APILsURL(file)
	fastload.Pipe(w, r, url, filter, 30, nil)
}

func info(w http.ResponseWriter, r *http.Request, match []string) {
	file := match[1]
	url := baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).APIFileInfoURL(file)
	fastload.Pipe(w, r, url, filter, 30, nil)
}

func get(w http.ResponseWriter, r *http.Request, match []string) {
	file := match[1]
	url := baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).GetDownloadURL(file)
	fastload.Pipe(w, r, url, filter, 3600, nil)
}

func fallback(w http.ResponseWriter, r *http.Request, match []string) {
	get(w, r, match)
}
