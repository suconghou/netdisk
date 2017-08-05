package middleware

import (
	"net/http"
	"netdisk/config"
	"netdisk/layers/baidudisk"
	"regexp"

	"github.com/suconghou/fastload/fastload"
)

type routeInfo struct {
	Reg     *regexp.Regexp
	Handler func(http.ResponseWriter, *http.Request, []string)
}

var route = []routeInfo{
	{regexp.MustCompile(`^ls/(.+)$`), ls},
	{regexp.MustCompile(`^info/(.+)$`), info},
}

// NetStream response stream data
func NetStream(w http.ResponseWriter, r *http.Request, match []string) {
	file := match[1]
	url := baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).GetDownloadURL(file)
	fastload.Pipe(w, r, url, nil)
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
		fallback(w, r)
	}
}

func ls(w http.ResponseWriter, r *http.Request, match []string) {
	file := match[1]
	url := baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).APILsURL(file)
	fastload.Pipe(w, r, url, nil)
}

func info(w http.ResponseWriter, r *http.Request, match []string) {
	file := match[1]
	url := baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).APIFileInfoURL(file)
	fastload.Pipe(w, r, url, nil)
}

func fallback(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
