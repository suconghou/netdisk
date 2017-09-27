package route

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"regexp"

	"github.com/suconghou/netdisk/middleware"
)

var (
	apiv1, _ = url.Parse("http://api.suconghou.cn/")
)

// 路由定义
type routeInfo struct {
	Reg     *regexp.Regexp
	Handler func(http.ResponseWriter, *http.Request, []string)
}

// RoutePath defines all routes
var RoutePath = []routeInfo{
	{regexp.MustCompile(`^/apiv1/(.+)$`), serveHTTP(apiv1)},
	{regexp.MustCompile(`^/files/(.+)$`), middleware.Files},
	{regexp.MustCompile(`^/imgs/(.+)$`), middleware.Imgs},
	{regexp.MustCompile(`^/play/(.+)$`), middleware.MediaStreamAPI},
	{regexp.MustCompile(`^/net/(.+)$`), middleware.NetStreamAPI},
	{regexp.MustCompile(`^/(?:(https?):/?)?(?:[\w\-]+\.)+[\w\-]+(?::\d+)?/.*$`), middleware.Pipe},
}

func serveHTTP(target *url.URL) func(http.ResponseWriter, *http.Request, []string) {
	return func(w http.ResponseWriter, r *http.Request, match []string) {
		proxy := &httputil.ReverseProxy{Director: func(req *http.Request) {
			req.Host = target.Host
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.URL.Path = path.Join("/", target.Path, match[1])
		}}
		proxy.ServeHTTP(w, r)
	}
}
